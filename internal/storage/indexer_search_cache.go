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

type IndexerSearchHistoryInput struct {
	IndexerID       uuid.UUID
	IndexerName     string
	IndexerProtocol string
	MediaType       string
	Query           string
	CacheHit        bool
	Success         bool
	ResultCount     int32
	Error           *string
	Response        any
}

func (s *SettingsStore) GetIndexerSearchSettings(ctx context.Context) (IndexerSearchSettings, error) {
	if err := s.ensureIndexerSearchSettings(ctx); err != nil {
		return IndexerSearchSettings{}, err
	}
	row, err := storagegen.New(s.pool).GetIndexerSearchSettings(ctx)
	return indexerSearchSettingsFromRow(row), err
}

func (s *SettingsStore) SaveIndexerSearchSettings(ctx context.Context, input IndexerSearchSettings) (IndexerSearchSettings, error) {
	if input.CacheDurationMinutes < 0 || input.CacheDurationMinutes > 43200 {
		return IndexerSearchSettings{}, ErrInvalidInput
	}
	if input.HistoryRetentionDays < 1 || input.HistoryRetentionDays > 365 {
		return IndexerSearchSettings{}, ErrInvalidInput
	}
	if input.AutomaticBlocklistExpiryDays < 1 || input.AutomaticBlocklistExpiryDays > 365 {
		return IndexerSearchSettings{}, ErrInvalidInput
	}
	err := storagegen.New(s.pool).SaveIndexerSearchSettings(ctx, storagegen.SaveIndexerSearchSettingsParams{
		CacheDurationMinutes:         input.CacheDurationMinutes,
		HistoryRetentionDays:         input.HistoryRetentionDays,
		AutomaticBlocklistExpiryDays: input.AutomaticBlocklistExpiryDays,
	})
	if err != nil {
		return IndexerSearchSettings{}, err
	}
	return s.GetIndexerSearchSettings(ctx)
}

func (s *SettingsStore) GetIndexerSearchCache(ctx context.Context, indexerID uuid.UUID, mediaType string, query string, target any) (bool, error) {
	raw, err := storagegen.New(s.pool).GetIndexerSearchCacheResponse(ctx, storagegen.GetIndexerSearchCacheResponseParams{
		IndexerID: indexerID,
		MediaType: mediaType,
		Query:     query,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal(raw, target)
}

func (s *SettingsStore) SetIndexerSearchCache(ctx context.Context, indexerID uuid.UUID, mediaType string, query string, response any, resultCount int32, expiresAt time.Time) (IndexerSearchCacheEntry, error) {
	raw, err := json.Marshal(response)
	if err != nil {
		return IndexerSearchCacheEntry{}, err
	}
	err = storagegen.New(s.pool).SetIndexerSearchCache(ctx, storagegen.SetIndexerSearchCacheParams{
		IndexerID:   indexerID,
		MediaType:   mediaType,
		Query:       query,
		Response:    raw,
		ResultCount: resultCount,
		ExpiresAt:   expiresAt,
	})
	if err != nil {
		return IndexerSearchCacheEntry{}, err
	}
	return s.GetIndexerSearchCacheEntry(ctx, indexerID, mediaType, query)
}

func (s *SettingsStore) RecordIndexerSearchHistory(ctx context.Context, input IndexerSearchHistoryInput) (IndexerSearchHistoryEntry, error) {
	raw, err := json.Marshal(input.Response)
	if err != nil {
		return IndexerSearchHistoryEntry{}, err
	}
	indexerID := input.IndexerID
	row, err := storagegen.New(s.pool).RecordIndexerSearchHistory(ctx, storagegen.RecordIndexerSearchHistoryParams{
		ID:              uuid.New(),
		IndexerID:       &indexerID,
		IndexerName:     input.IndexerName,
		IndexerProtocol: input.IndexerProtocol,
		MediaType:       input.MediaType,
		Query:           input.Query,
		CacheHit:        input.CacheHit,
		Success:         input.Success,
		ResultCount:     input.ResultCount,
		Error:           textValue(input.Error),
		Response:        raw,
	})
	return indexerSearchHistoryEntryFromRecordRow(row), err
}

func (s *SettingsStore) CleanupIndexerSearchHistory(ctx context.Context, retentionDays int32) (int32, error) {
	rows, err := storagegen.New(s.pool).CleanupIndexerSearchHistory(ctx, retentionDays)
	if err != nil {
		return 0, err
	}
	return int32(rows), nil
}

func (s *SettingsStore) ClearIndexerSearchHistory(ctx context.Context) (int32, error) {
	rows, err := storagegen.New(s.pool).ClearIndexerSearchHistory(ctx)
	if err != nil {
		return 0, err
	}
	return int32(rows), nil
}

func (s *SettingsStore) ClearIndexerSearchCache(ctx context.Context) (int32, error) {
	rows, err := storagegen.New(s.pool).ClearIndexerSearchCache(ctx)
	if err != nil {
		return 0, err
	}
	return int32(rows), nil
}

func (s *SettingsStore) ClearIndexerSearchCacheByPattern(ctx context.Context, pattern string) (int32, error) {
	rows, err := storagegen.New(s.pool).ClearIndexerSearchCacheByPattern(ctx, pattern)
	if err != nil {
		return 0, err
	}
	return int32(rows), nil
}

func (s *SettingsStore) DeleteIndexerSearchCacheEntry(ctx context.Context, indexerID uuid.UUID, mediaType string, query string) (int32, error) {
	rows, err := storagegen.New(s.pool).DeleteIndexerSearchCacheEntry(ctx, storagegen.DeleteIndexerSearchCacheEntryParams{
		IndexerID: indexerID,
		MediaType: mediaType,
		Query:     query,
	})
	if err != nil {
		return 0, err
	}
	return int32(rows), nil
}

func (s *SettingsStore) ensureIndexerSearchSettings(ctx context.Context) error {
	return storagegen.New(s.pool).EnsureIndexerSearchSettings(ctx)
}

func indexerSearchSettingsFromRow(row storagegen.GetIndexerSearchSettingsRow) IndexerSearchSettings {
	return IndexerSearchSettings{
		CacheDurationMinutes:         row.CacheDurationMinutes,
		HistoryRetentionDays:         row.HistoryRetentionDays,
		AutomaticBlocklistExpiryDays: row.AutomaticBlocklistExpiryDays,
	}
}

func indexerSearchHistoryEntryFromRecordRow(row storagegen.RecordIndexerSearchHistoryRow) IndexerSearchHistoryEntry {
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
