package jobs

import (
	"context"
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
	return searchReleasesWithProgress(ctx, settings, indexerService, item, query, eventBroker, manual, nil)
}

func searchReleasesWithProgress(
	ctx context.Context,
	settings *storage.SettingsStore,
	indexerService *indexers.Service,
	item storage.MediaItem,
	query string,
	eventBroker *events.Broker,
	manual bool,
	progress ReleaseSearchProgress,
) ([]storage.ReleaseCandidateInput, []string, error) {
	configs, err := settings.ListEnabledIndexers(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("list enabled indexers: %w", err)
	}
	if len(configs) == 0 {
		slog.Debug("release search skipped because no indexers are enabled", "mediaItemId", item.ID, "title", item.Title)
		publishReleaseSearchProgress(progress, "No enabled indexer is configured")
		return nil, []string{"No enabled indexer is configured"}, nil
	}

	releases := []storage.ReleaseCandidateInput{}
	searchErrors := []string{}
	criteria := decisions.SearchCriteriaForQuery(item, query)
	queries := decisions.SearchQueriesForCriteria(criteria, query)
	publishReleaseSearchProgress(progress, "Searching %d indexer(s) with %d query branch(es)", len(configs), len(queries))
	limiter := newIndexerRateLimiter()
	cacheSettings, err := settings.GetIndexerSearchSettings(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("load indexer search settings: %w", err)
	}
	if _, err := settings.CleanupIndexerSearchHistory(ctx, cacheSettings.HistoryRetentionDays); err != nil {
		slog.Error("indexer search history cleanup failed", "error", err)
	}
	releases, searchErrors, err = searchReleaseQueries(ctx, releaseQuerySearch{
		settings:       settings,
		indexerService: indexerService,
		limiter:        limiter,
		configs:        configs,
		item:           item,
		criteria:       criteria,
		queries:        queries,
		cacheSettings:  cacheSettings,
		eventBroker:    eventBroker,
		manual:         manual,
		progress:       progress,
	})
	if err != nil {
		return nil, searchErrors, err
	}
	if len(releases) == 0 && criteria.Kind == "episode" && criteria.SeasonNumber != nil {
		publishReleaseSearchProgress(progress, "No episode releases found; searching the whole season")
		seasonCriteria := decisions.ReleaseSearchCriteria{
			Kind:         "season",
			Title:        criteria.Title,
			Year:         criteria.Year,
			SeasonNumber: criteria.SeasonNumber,
		}
		seasonQueries := decisions.SearchQueriesForCriteria(seasonCriteria, "")
		seasonReleases, seasonErrors, err := searchReleaseQueries(ctx, releaseQuerySearch{
			settings:       settings,
			indexerService: indexerService,
			limiter:        limiter,
			configs:        configs,
			item:           item,
			criteria:       criteria,
			queries:        seasonQueries,
			cacheSettings:  cacheSettings,
			eventBroker:    eventBroker,
			manual:         manual,
			progress:       progress,
		})
		if err != nil {
			return nil, seasonErrors, err
		}
		releases = append(releases, seasonReleases...)
		searchErrors = append(searchErrors, seasonErrors...)
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
