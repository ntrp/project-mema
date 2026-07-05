package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *SettingsStore) GetMetadataSearchCache(ctx context.Context, providerID uuid.UUID, mediaType string, query string, year *int32, target any) (bool, error) {
	raw, err := storagegen.New(s.pool).GetMetadataSearchCacheResults(ctx, storagegen.GetMetadataSearchCacheResultsParams{
		ProviderID: providerID,
		MediaType:  mediaType,
		Query:      query,
		Year:       cacheYear(year),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return false, err
	}
	return true, nil
}

func (s *SettingsStore) SetMetadataSearchCache(ctx context.Context, providerID uuid.UUID, mediaType string, query string, year *int32, value any, expiresAt time.Time) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return storagegen.New(s.pool).SetMetadataSearchCache(ctx, storagegen.SetMetadataSearchCacheParams{
		ProviderID: providerID,
		MediaType:  mediaType,
		Query:      query,
		Year:       cacheYear(year),
		Results:    raw,
		ExpiresAt:  expiresAt,
	})
}

func (s *SettingsStore) RecordMetadataSearchHistory(ctx context.Context, input MetadataSearchHistoryInput) (MetadataSearchHistoryEntry, error) {
	raw, err := json.Marshal(input.Response)
	if err != nil {
		return MetadataSearchHistoryEntry{}, err
	}
	providerID := input.ProviderID
	row, err := storagegen.New(s.pool).RecordMetadataSearchHistory(ctx, storagegen.RecordMetadataSearchHistoryParams{
		ID:           uuid.New(),
		ProviderID:   &providerID,
		ProviderName: input.ProviderName,
		ProviderType: input.ProviderType,
		MediaType:    input.MediaType,
		Query:        input.Query,
		Year:         cacheYear(input.Year),
		CacheHit:     input.CacheHit,
		Success:      input.Success,
		ItemCount:    input.ItemCount,
		Error:        textValue(input.Error),
		Response:     raw,
	})
	if err != nil {
		return MetadataSearchHistoryEntry{}, err
	}
	return metadataSearchHistoryEntryFromRecordRow(row), nil
}

func (s *SettingsStore) MetadataCacheStats(ctx context.Context) (MetadataCacheStats, error) {
	row, err := storagegen.New(s.pool).MetadataCacheStats(ctx)
	return metadataCacheStatsFromRow(row), err
}

func (s *SettingsStore) ListMetadataSearchHistoryEntries(ctx context.Context, limit int32) ([]MetadataSearchHistoryEntry, error) {
	limit = inspectionLimit(limit)
	rows, err := storagegen.New(s.pool).ListMetadataSearchHistoryEntries(ctx, limit)
	if err != nil {
		return nil, err
	}
	entries := make([]MetadataSearchHistoryEntry, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, metadataSearchHistoryEntryFromListRow(row))
	}
	return entries, nil
}

func (s *SettingsStore) ListMetadataCacheEntries(ctx context.Context, limit int32) ([]MetadataCacheEntry, error) {
	limit = inspectionLimit(limit)
	rows, err := storagegen.New(s.pool).ListMetadataCacheEntries(ctx, limit)
	if err != nil {
		return nil, err
	}
	entries := make([]MetadataCacheEntry, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, metadataCacheEntryFromListRow(row))
	}
	return entries, nil
}

func (s *SettingsStore) MetadataSearchHistoryCount(ctx context.Context) (int32, error) {
	return storagegen.New(s.pool).MetadataSearchHistoryCount(ctx)
}

func (s *SettingsStore) MetadataSearchHistoryStats(ctx context.Context) (QueryHistoryStats, error) {
	row, err := storagegen.New(s.pool).MetadataSearchHistoryStats(ctx)
	return metadataSearchHistoryStatsFromRow(row), err
}

func (s *SettingsStore) ClearMetadataCache(ctx context.Context) (int32, error) {
	rows, err := storagegen.New(s.pool).ClearMetadataCache(ctx)
	if err != nil {
		return 0, err
	}
	return int32(rows), nil
}

func (s *SettingsStore) ClearMetadataCacheByPattern(ctx context.Context, pattern string) (int32, error) {
	rows, err := storagegen.New(s.pool).ClearMetadataCacheByPattern(ctx, pattern)
	if err != nil {
		return 0, err
	}
	return int32(rows), nil
}

func (s *SettingsStore) DeleteMetadataCacheEntry(ctx context.Context, providerID uuid.UUID, mediaType string, query string, year int32) (int32, error) {
	rows, err := storagegen.New(s.pool).DeleteMetadataCacheEntry(ctx, storagegen.DeleteMetadataCacheEntryParams{
		ProviderID: providerID,
		MediaType:  mediaType,
		Query:      query,
		Year:       year,
	})
	if err != nil {
		return 0, err
	}
	return int32(rows), nil
}

func (s *SettingsStore) ClearMetadataSearchHistory(ctx context.Context) (int32, error) {
	rows, err := storagegen.New(s.pool).ClearMetadataSearchHistory(ctx)
	if err != nil {
		return 0, err
	}
	return int32(rows), nil
}

func cacheYear(year *int32) int32 {
	if year == nil {
		return 0
	}
	return *year
}

func metadataCacheKind(query string) string {
	switch {
	case len(query) >= 9 && query[:9] == "discover:":
		return "discover"
	case len(query) >= 8 && query[:8] == "details:":
		return "details"
	default:
		return "search"
	}
}

func metadataSearchHistoryEntryFromRecordRow(row storagegen.RecordMetadataSearchHistoryRow) MetadataSearchHistoryEntry {
	entry := MetadataSearchHistoryEntry{
		ProviderName: row.ProviderName,
		ProviderType: row.ProviderType,
		MediaType:    row.MediaType,
		Query:        row.Query,
		Year:         row.Year,
		CacheHit:     row.CacheHit,
		Success:      row.Success,
		ItemCount:    row.ItemCount,
		Error:        textPtr(row.Error),
		Response:     row.Response,
		CreatedAt:    row.CreatedAt,
	}
	entry.CacheKind = metadataCacheKind(entry.Query)
	return entry
}

func metadataCacheStatsFromRow(row storagegen.MetadataCacheStatsRow) MetadataCacheStats {
	return MetadataCacheStats{
		TotalEntries:   row.TotalEntries,
		ActiveEntries:  row.ActiveEntries,
		ExpiredEntries: row.ExpiredEntries,
		ProviderCount:  row.ProviderCount,
	}
}

func metadataSearchHistoryEntryFromListRow(row storagegen.ListMetadataSearchHistoryEntriesRow) MetadataSearchHistoryEntry {
	entry := MetadataSearchHistoryEntry{
		ProviderName: row.ProviderName,
		ProviderType: row.ProviderType,
		MediaType:    row.MediaType,
		Query:        row.Query,
		Year:         row.Year,
		CacheHit:     row.CacheHit,
		Success:      row.Success,
		ItemCount:    row.ItemCount,
		Error:        textPtr(row.Error),
		Response:     row.Response,
		CreatedAt:    row.CreatedAt,
	}
	entry.CacheKind = metadataCacheKind(entry.Query)
	return entry
}

func metadataCacheEntryFromListRow(row storagegen.ListMetadataCacheEntriesRow) MetadataCacheEntry {
	return MetadataCacheEntry{
		ProviderID:   row.ProviderID,
		ProviderName: row.ProviderName,
		ProviderType: row.ProviderType,
		MediaType:    row.MediaType,
		Query:        row.Query,
		Year:         row.Year,
		ItemCount:    row.ItemCount,
		ExpiresAt:    row.ExpiresAt,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		Expired:      row.Expired,
	}
}

func metadataSearchHistoryStatsFromRow(row storagegen.MetadataSearchHistoryStatsRow) QueryHistoryStats {
	return QueryHistoryStats{
		TotalEntries: row.TotalEntries,
		CacheHits:    row.CacheHits,
		CacheMisses:  row.CacheMisses,
		Failures:     row.Failures,
	}
}
