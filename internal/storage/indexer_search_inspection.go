package storage

import (
	"context"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

func (s *SettingsStore) IndexerSearchCacheStats(ctx context.Context) (IndexerSearchCacheStats, error) {
	row, err := storagegen.New(s.pool).IndexerSearchCacheStats(ctx)
	return indexerSearchCacheStatsFromRow(row), err
}

func (s *SettingsStore) ListIndexerSearchCacheEntries(ctx context.Context, limit int32) ([]IndexerSearchCacheEntry, error) {
	limit = inspectionLimit(limit)
	rows, err := storagegen.New(s.pool).ListIndexerSearchCacheEntries(ctx, limit)
	if err != nil {
		return nil, err
	}
	entries := make([]IndexerSearchCacheEntry, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, indexerSearchCacheEntryFromListRow(row))
	}
	return entries, nil
}

func (s *SettingsStore) GetIndexerSearchCacheEntry(
	ctx context.Context,
	indexerID uuid.UUID,
	mediaType string,
	query string,
) (IndexerSearchCacheEntry, error) {
	row, err := storagegen.New(s.pool).GetIndexerSearchCacheEntry(ctx, storagegen.GetIndexerSearchCacheEntryParams{
		IndexerID: indexerID,
		MediaType: mediaType,
		Query:     query,
	})
	return indexerSearchCacheEntryFromGetRow(row), err
}

func (s *SettingsStore) ListIndexerSearchHistoryEntries(ctx context.Context, limit int32) ([]IndexerSearchHistoryEntry, error) {
	limit = inspectionLimit(limit)
	rows, err := storagegen.New(s.pool).ListIndexerSearchHistoryEntries(ctx, limit)
	if err != nil {
		return nil, err
	}
	entries := make([]IndexerSearchHistoryEntry, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, indexerSearchHistoryEntryFromListRow(row))
	}
	return entries, nil
}

func (s *SettingsStore) IndexerSearchHistoryCount(ctx context.Context) (int32, error) {
	return storagegen.New(s.pool).IndexerSearchHistoryCount(ctx)
}

func (s *SettingsStore) IndexerSearchHistoryStats(ctx context.Context) (QueryHistoryStats, error) {
	row, err := storagegen.New(s.pool).IndexerSearchHistoryStats(ctx)
	return indexerSearchHistoryStatsFromRow(row), err
}

func inspectionLimit(limit int32) int32 {
	if limit <= 0 {
		return 10
	}
	if limit > 500 {
		return 500
	}
	return limit
}

func indexerSearchCacheStatsFromRow(row storagegen.IndexerSearchCacheStatsRow) IndexerSearchCacheStats {
	return IndexerSearchCacheStats{
		TotalEntries:   row.TotalEntries,
		ActiveEntries:  row.ActiveEntries,
		ExpiredEntries: row.ExpiredEntries,
		IndexerCount:   row.IndexerCount,
	}
}

func indexerSearchCacheEntryFromListRow(row storagegen.ListIndexerSearchCacheEntriesRow) IndexerSearchCacheEntry {
	return IndexerSearchCacheEntry{
		IndexerID:       row.IndexerID,
		IndexerName:     row.IndexerName,
		IndexerProtocol: row.IndexerProtocol,
		MediaType:       row.MediaType,
		Query:           row.Query,
		ResultCount:     row.ResultCount,
		ExpiresAt:       row.ExpiresAt,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		Expired:         row.Expired,
	}
}

func indexerSearchCacheEntryFromGetRow(row storagegen.GetIndexerSearchCacheEntryRow) IndexerSearchCacheEntry {
	return IndexerSearchCacheEntry{
		IndexerID:       row.IndexerID,
		IndexerName:     row.IndexerName,
		IndexerProtocol: row.IndexerProtocol,
		MediaType:       row.MediaType,
		Query:           row.Query,
		ResultCount:     row.ResultCount,
		ExpiresAt:       row.ExpiresAt,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		Expired:         row.Expired,
	}
}

func indexerSearchHistoryEntryFromListRow(row storagegen.ListIndexerSearchHistoryEntriesRow) IndexerSearchHistoryEntry {
	return IndexerSearchHistoryEntry{
		IndexerName:     row.IndexerName,
		IndexerProtocol: row.IndexerProtocol,
		MediaType:       row.MediaType,
		Query:           row.Query,
		CacheHit:        row.CacheHit,
		Success:         row.Success,
		ResultCount:     row.ResultCount,
		Error:           textPtr(row.Error),
		Response:        row.Response,
		CreatedAt:       row.CreatedAt,
	}
}

func indexerSearchHistoryStatsFromRow(row storagegen.IndexerSearchHistoryStatsRow) QueryHistoryStats {
	return QueryHistoryStats{
		TotalEntries: row.TotalEntries,
		CacheHits:    row.CacheHits,
		CacheMisses:  row.CacheMisses,
		Failures:     row.Failures,
	}
}
