package jobs

import (
	"context"
	"time"

	"media-manager/internal/events"
	"media-manager/internal/storage"
)

type indexerSearchCacheEntryEvent struct {
	IndexerName string    `json:"indexerName"`
	IndexerType string    `json:"indexerType"`
	MediaType   string    `json:"mediaType"`
	Query       string    `json:"query"`
	ResultCount int32     `json:"resultCount"`
	ExpiresAt   time.Time `json:"expiresAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Expired     bool      `json:"expired"`
}

type indexerSearchCacheUpdatedEvent struct {
	Entry indexerSearchCacheEntryEvent `json:"entry"`
	Stats indexerSearchStatsEvent      `json:"stats"`
}

type indexerSearchHistoryEntryEvent struct {
	IndexerName string    `json:"indexerName"`
	IndexerType string    `json:"indexerType"`
	MediaType   string    `json:"mediaType"`
	Query       string    `json:"query"`
	CacheHit    bool      `json:"cacheHit"`
	Success     bool      `json:"success"`
	ResultCount int32     `json:"resultCount"`
	Error       *string   `json:"error"`
	Response    string    `json:"response"`
	CreatedAt   time.Time `json:"createdAt"`
}

type indexerSearchStatsEvent struct {
	TotalEntries   int32 `json:"totalEntries"`
	ActiveEntries  int32 `json:"activeEntries"`
	ExpiredEntries int32 `json:"expiredEntries"`
	IndexerCount   int32 `json:"indexerCount"`
}

func publishIndexerSearchCacheUpdated(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	entry storage.IndexerSearchCacheEntry,
) {
	stats, err := settings.IndexerSearchCacheStats(ctx)
	if err != nil {
		return
	}
	eventBroker.Publish("indexer.search.cache.updated", indexerSearchCacheUpdatedEvent{
		Entry: indexerSearchCacheEntryPayload(entry),
		Stats: indexerSearchStatsEvent{
			TotalEntries:   stats.TotalEntries,
			ActiveEntries:  stats.ActiveEntries,
			ExpiredEntries: stats.ExpiredEntries,
			IndexerCount:   stats.IndexerCount,
		},
	})
}

func publishIndexerSearchHistoryCreated(eventBroker *events.Broker, entry storage.IndexerSearchHistoryEntry) {
	eventBroker.Publish("indexer.search.history.created", indexerSearchHistoryEntryEvent{
		IndexerName: entry.IndexerName,
		IndexerType: entry.IndexerType,
		MediaType:   entry.MediaType,
		Query:       entry.Query,
		CacheHit:    entry.CacheHit,
		Success:     entry.Success,
		ResultCount: entry.ResultCount,
		Error:       entry.Error,
		Response:    entry.Response,
		CreatedAt:   entry.CreatedAt,
	})
}

func indexerSearchCacheEntryPayload(entry storage.IndexerSearchCacheEntry) indexerSearchCacheEntryEvent {
	return indexerSearchCacheEntryEvent{
		IndexerName: entry.IndexerName,
		IndexerType: entry.IndexerType,
		MediaType:   entry.MediaType,
		Query:       entry.Query,
		ResultCount: entry.ResultCount,
		ExpiresAt:   entry.ExpiresAt,
		CreatedAt:   entry.CreatedAt,
		UpdatedAt:   entry.UpdatedAt,
		Expired:     entry.Expired,
	}
}
