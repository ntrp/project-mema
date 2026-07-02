package jobs

import (
	"context"
	"log/slog"
	"time"

	"media-manager/internal/events"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

type indexerSearchResponse struct {
	Releases []indexers.Release `json:"releases,omitempty"`
	Error    *string            `json:"error,omitempty"`
}

func executeIndexerSearch(
	ctx context.Context,
	settings *storage.SettingsStore,
	indexerService *indexers.Service,
	limiter *indexerRateLimiter,
	config storage.Indexer,
	mediaType string,
	query string,
	cacheSettings storage.IndexerSearchSettings,
	eventBroker *events.Broker,
) ([]indexers.Release, bool, error) {
	if cacheSettings.CacheDurationMinutes > 0 {
		cached := indexerSearchResponse{}
		found, err := settings.GetIndexerSearchCache(ctx, config.ID, mediaType, query, &cached)
		if err != nil {
			slog.Error("indexer cache read failed", "indexerName", config.Name, "query", query, "error", err)
		}
		if found {
			recordIndexerHistory(ctx, settings, eventBroker, config, mediaType, query, true, cached.Releases, nil)
			return cached.Releases, true, nil
		}
	}

	if err := limiter.wait(ctx, config.ID); err != nil {
		recordIndexerHistory(ctx, settings, eventBroker, config, mediaType, query, false, nil, err)
		return nil, false, err
	}
	found, err := indexerService.Search(ctx, indexerConfig(config), query, mediaType)
	if err != nil {
		recordIndexerHistory(ctx, settings, eventBroker, config, mediaType, query, false, nil, err)
		return nil, false, err
	}
	if cacheSettings.CacheDurationMinutes > 0 {
		response := indexerSearchResponse{Releases: found}
		expiresAt := time.Now().Add(time.Duration(cacheSettings.CacheDurationMinutes) * time.Minute)
		entry, err := settings.SetIndexerSearchCache(ctx, config.ID, mediaType, query, response, int32(len(found)), expiresAt)
		if err != nil {
			slog.Error("indexer cache write failed", "indexerName", config.Name, "query", query, "error", err)
		} else {
			publishIndexerSearchCacheUpdated(ctx, settings, eventBroker, entry)
		}
	}
	recordIndexerHistory(ctx, settings, eventBroker, config, mediaType, query, false, found, nil)
	return found, false, nil
}

func recordIndexerHistory(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	config storage.Indexer,
	mediaType string,
	query string,
	cacheHit bool,
	releases []indexers.Release,
	searchErr error,
) {
	var message *string
	response := indexerSearchResponse{Releases: releases}
	if searchErr != nil {
		text := searchErr.Error()
		message = &text
		response.Error = &text
	}
	entry, err := settings.RecordIndexerSearchHistory(ctx, storage.IndexerSearchHistoryInput{
		IndexerID:   config.ID,
		IndexerName: config.Name,
		IndexerType: config.Type,
		MediaType:   mediaType,
		Query:       query,
		CacheHit:    cacheHit,
		Success:     searchErr == nil,
		ResultCount: int32(len(releases)),
		Error:       message,
		Response:    response,
	})
	if err != nil {
		slog.Error("indexer history write failed", "indexerName", config.Name, "query", query, "error", err)
		return
	}
	publishIndexerSearchHistoryCreated(eventBroker, entry)
}
