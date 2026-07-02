package jobs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/decisions"
	"media-manager/internal/events"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

func searchReleases(
	ctx context.Context,
	settings *storage.SettingsStore,
	indexerService *indexers.Service,
	item storage.MediaItem,
	query string,
	eventBroker *events.Broker,
	manual bool,
) ([]storage.ReleaseCandidateInput, []string, error) {
	configs, err := settings.ListEnabledIndexers(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("list enabled indexers: %w", err)
	}
	if len(configs) == 0 {
		slog.Debug("release search skipped because no indexers are enabled", "mediaItemId", item.ID, "title", item.Title)
		return nil, []string{"No enabled indexer is configured"}, nil
	}

	releases := []storage.ReleaseCandidateInput{}
	searchErrors := []string{}
	criteria := decisions.SearchCriteriaForQuery(item, query)
	queries := decisions.SearchQueriesForCriteria(criteria, query)
	limiter := newIndexerRateLimiter()
	cacheSettings, err := settings.GetIndexerSearchSettings(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("load indexer search settings: %w", err)
	}
	if _, err := settings.CleanupIndexerSearchHistory(ctx, cacheSettings.HistoryRetentionDays); err != nil {
		slog.Error("indexer search history cleanup failed", "error", err)
	}
	for _, config := range configs {
		for _, searchQuery := range queries {
			found, cacheHit, err := executeIndexerSearch(ctx, settings, indexerService, limiter, config, item.Type, searchQuery, cacheSettings, eventBroker)
			if err != nil {
				if manual && isIndexer429(err) {
					publishManualIndexerRateLimitEvent(ctx, settings, eventBroker, config, item, searchQuery, err)
					return nil, []string{fmt.Sprintf("%s (%s): %s", config.Name, searchQuery, err.Error())}, err
				}
				if !errors.Is(err, errIndexerBackoffActive) {
					recordIndexerSearchFailure(ctx, settings, eventBroker, config, err)
					limiter.recordError(config.ID, err)
				}
				slog.Error("indexer release search failed", "mediaItemId", item.ID, "title", item.Title, "indexerName", config.Name, "query", searchQuery, "error", err)
				searchErrors = append(searchErrors, fmt.Sprintf("%s (%s): %s", config.Name, searchQuery, err.Error()))
				continue
			}
			if !cacheHit {
				recordIndexerSearchSuccess(ctx, settings, config)
			}
			slog.Debug("indexer release search finished", "mediaItemId", item.ID, "title", item.Title, "indexerName", config.Name, "query", searchQuery, "cacheHit", cacheHit, "releaseCount", len(found))
			for _, release := range found {
				releases = append(releases, releaseCandidateInput(item.ID, release, criteria))
			}
		}
	}
	releases = dedupeReleaseCandidates(item, releases)
	sortReleaseCandidates(releases)
	if len(releases) == 0 && len(searchErrors) == 0 {
		searchErrors = append(searchErrors, "No releases found")
	}
	return releases, searchErrors, nil
}

func isIndexer429(err error) bool {
	statusCode := indexers.StatusCode(err)
	return statusCode != nil && *statusCode == http.StatusTooManyRequests
}

func publishManualIndexerRateLimitEvent(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	config storage.Indexer,
	item storage.MediaItem,
	query string,
	err error,
) {
	publishSystemEvent(ctx, settings, eventBroker, jobEventWarning, "indexers", "Manual search failed because an indexer rate limit was reached", map[string]any{
		"indexerId":   config.ID.String(),
		"indexerName": config.Name,
		"mediaItemId": item.ID.String(),
		"title":       item.Title,
		"query":       query,
		"error":       err.Error(),
	})
}

func recordIndexerSearchSuccess(ctx context.Context, settings *storage.SettingsStore, config storage.Indexer) {
	if _, err := settings.RecordIndexerSuccess(ctx, config.ID); err != nil {
		slog.Error("indexer success state update failed", "indexerName", config.Name, "error", err)
	}
}

func dedupeReleaseCandidates(
	item storage.MediaItem,
	releases []storage.ReleaseCandidateInput,
) []storage.ReleaseCandidateInput {
	byKey := map[string]storage.ReleaseCandidateInput{}
	for _, release := range releases {
		key := releaseDedupeKey(release)
		if key == "" {
			byKey[uuid.NewString()] = release
			continue
		}
		if existing, ok := byKey[key]; !ok || betterCandidate(item, release, existing) {
			byKey[key] = release
		}
	}
	deduped := make([]storage.ReleaseCandidateInput, 0, len(byKey))
	for _, release := range byKey {
		deduped = append(deduped, release)
	}
	return deduped
}

func betterCandidate(
	item storage.MediaItem,
	left storage.ReleaseCandidateInput,
	right storage.ReleaseCandidateInput,
) bool {
	leftMatch := decisions.EvaluateReleaseCandidateInputMatch(item, left)
	rightMatch := decisions.EvaluateReleaseCandidateInputMatch(item, right)
	if leftMatch.Severity != rightMatch.Severity {
		return severityRank(leftMatch.Severity) > severityRank(rightMatch.Severity)
	}
	if left.Seeders != nil && right.Seeders != nil && *left.Seeders != *right.Seeders {
		return *left.Seeders > *right.Seeders
	}
	if left.SizeBytes != right.SizeBytes {
		return left.SizeBytes > right.SizeBytes
	}
	return strings.ToLower(left.Title) < strings.ToLower(right.Title)
}

func sortReleaseCandidates(releases []storage.ReleaseCandidateInput) {
	sort.SliceStable(releases, func(i, j int) bool {
		left := releases[i]
		right := releases[j]
		if left.Seeders != nil && right.Seeders != nil && *left.Seeders != *right.Seeders {
			return *left.Seeders > *right.Seeders
		}
		if left.SizeBytes != right.SizeBytes {
			return left.SizeBytes > right.SizeBytes
		}
		return strings.ToLower(left.Title) < strings.ToLower(right.Title)
	})
}

func releaseDedupeKey(release storage.ReleaseCandidateInput) string {
	for _, value := range []*string{release.GUID, release.InfoURL, &release.DownloadURL} {
		if value != nil && strings.TrimSpace(*value) != "" {
			return strings.ToLower(strings.TrimSpace(*value))
		}
	}
	return ""
}

func severityRank(severity string) int {
	switch severity {
	case "info":
		return 3
	case "warning":
		return 2
	default:
		return 1
	}
}

func indexerConfig(indexer storage.Indexer) indexers.Config {
	return indexers.Config{
		ID:         indexer.ID.String(),
		Name:       indexer.Name,
		Type:       indexer.Type,
		BaseURL:    indexer.BaseURL,
		APIKey:     indexer.APIKey,
		Categories: indexer.Categories,
	}
}

func releaseCandidateInput(
	mediaItemID uuid.UUID,
	release indexers.Release,
	criteria decisions.ReleaseSearchCriteria,
) storage.ReleaseCandidateInput {
	var indexerID *uuid.UUID
	if parsed, err := uuid.Parse(release.IndexerID); err == nil {
		indexerID = &parsed
	}
	return storage.ReleaseCandidateInput{
		MediaItemID:      mediaItemID,
		IndexerID:        indexerID,
		IndexerName:      release.IndexerName,
		IndexerType:      release.IndexerType,
		Title:            release.Title,
		DownloadURL:      release.DownloadURL,
		InfoURL:          optionalString(release.InfoURL),
		GUID:             optionalString(release.GUID),
		SizeBytes:        release.SizeBytes,
		Seeders:          release.Seeders,
		Peers:            release.Peers,
		PublishedAt:      release.PublishedAt,
		SearchKind:       criteria.Kind,
		RequestedSeason:  criteria.SeasonNumber,
		RequestedEpisode: criteria.EpisodeNumber,
	}
}
