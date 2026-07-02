package jobs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"media-manager/internal/decisions"
	"media-manager/internal/events"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

type releaseQuerySearch struct {
	settings       *storage.SettingsStore
	indexerService *indexers.Service
	limiter        *indexerRateLimiter
	configs        []storage.Indexer
	item           storage.MediaItem
	criteria       decisions.ReleaseSearchCriteria
	queries        []string
	cacheSettings  storage.IndexerSearchSettings
	eventBroker    *events.Broker
	manual         bool
}

func searchReleaseQueries(
	ctx context.Context,
	input releaseQuerySearch,
) ([]storage.ReleaseCandidateInput, []string, error) {
	releases := []storage.ReleaseCandidateInput{}
	searchErrors := []string{}
	for _, config := range input.configs {
		for _, searchQuery := range input.queries {
			found, cacheHit, err := executeIndexerSearch(ctx, input.settings, input.indexerService, input.limiter, config, input.item.Type, searchQuery, input.cacheSettings, input.eventBroker)
			if err != nil {
				if input.manual && isIndexer429(err) {
					publishManualIndexerRateLimitEvent(ctx, input.settings, input.eventBroker, config, input.item, searchQuery, err)
					searchErrors = append(searchErrors, branchError(config, searchQuery, err))
					continue
				}
				if !errors.Is(err, errIndexerBackoffActive) {
					recordIndexerSearchFailure(ctx, input.settings, input.eventBroker, config, err)
					input.limiter.recordError(config.ID, err)
				}
				if input.manual {
					publishManualIndexerBranchFailureEvent(ctx, input.settings, input.eventBroker, config, input.item, searchQuery, err)
				}
				slog.Error("indexer release search failed", "mediaItemId", input.item.ID, "title", input.item.Title, "indexerName", config.Name, "query", searchQuery, "error", err)
				searchErrors = append(searchErrors, branchError(config, searchQuery, err))
				continue
			}
			if !cacheHit {
				recordIndexerSearchSuccess(ctx, input.settings, config)
			}
			slog.Debug("indexer release search finished", "mediaItemId", input.item.ID, "title", input.item.Title, "indexerName", config.Name, "query", searchQuery, "cacheHit", cacheHit, "releaseCount", len(found))
			for _, release := range found {
				releases = append(releases, releaseCandidateInput(input.item.ID, release, input.criteria))
			}
		}
	}
	return releases, searchErrors, nil
}

func branchError(config storage.Indexer, query string, err error) string {
	return fmt.Sprintf("%s (%s): %s", config.Name, query, err.Error())
}

func publishManualIndexerBranchFailureEvent(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	config storage.Indexer,
	item storage.MediaItem,
	query string,
	err error,
) {
	publishSystemEvent(ctx, settings, eventBroker, jobEventError, "indexers", "Manual search branch failed", map[string]any{
		"indexerId":   config.ID.String(),
		"indexerName": config.Name,
		"mediaItemId": item.ID.String(),
		"title":       item.Title,
		"query":       query,
		"error":       err.Error(),
	})
}
