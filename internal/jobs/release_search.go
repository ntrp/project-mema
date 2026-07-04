package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

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
	profile, formats, languages := releaseDecisionContext(ctx, settings, item)
	criteria := decisions.SearchCriteriaForQuery(item, query)
	branches := releaseSearchBranches(item, criteria, query)
	branchCount := 0
	for _, branch := range branches {
		branchCount += len(branch.queries)
	}
	publishReleaseSearchProgress(progress, "Searching %d indexer(s) with %d query branch(es)", len(configs), branchCount)
	limiter := newIndexerRateLimiter()
	cacheSettings, err := settings.GetIndexerSearchSettings(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("load indexer search settings: %w", err)
	}
	if _, err := settings.CleanupIndexerSearchHistory(ctx, cacheSettings.HistoryRetentionDays); err != nil {
		slog.Error("indexer search history cleanup failed", "error", err)
	}
	for _, branch := range branches {
		branchReleases, branchErrors, err := searchReleaseQueries(ctx, releaseQuerySearch{
			settings:       settings,
			indexerService: indexerService,
			limiter:        limiter,
			configs:        configs,
			item:           item,
			criteria:       branch.criteria,
			queries:        branch.queries,
			cacheSettings:  cacheSettings,
			eventBroker:    eventBroker,
			manual:         manual,
			progress:       progress,
		})
		if err != nil {
			return nil, branchErrors, err
		}
		releases = append(releases, branchReleases...)
		searchErrors = append(searchErrors, branchErrors...)
	}
	if shouldFallbackEpisodeToSeason(item, criteria, profile, formats, languages, releases) {
		publishReleaseSearchProgress(progress, "No matching episode release found; searching the whole season")
		seasonCriteria := seasonFallbackCriteria(criteria)
		seasonQueries := decisions.SearchQueriesForCriteria(seasonCriteria, "")
		seasonReleases, seasonErrors, err := searchReleaseQueries(ctx, releaseQuerySearch{
			settings:       settings,
			indexerService: indexerService,
			limiter:        limiter,
			configs:        configs,
			item:           item,
			criteria:       seasonCriteria,
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
	releases = dedupeReleaseCandidates(item, profile, formats, languages, releases)
	sortReleaseCandidates(item, profile, formats, languages, releases)
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

func indexerConfig(indexer storage.Indexer) indexers.Config {
	return indexers.Config{
		ID:             indexer.ID.String(),
		DefinitionID:   indexer.DefinitionID,
		Name:           indexer.Name,
		Implementation: indexer.Implementation,
		Protocol:       indexer.Protocol,
		BaseURL:        indexer.BaseURL,
		APIKey:         indexer.APIKey,
		Categories:     indexer.Categories,
		Fields:         append([]byte(nil), indexer.Fields...),
		Redirect:       indexer.Redirect,
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
		IndexerProtocol:  release.IndexerProtocol,
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
