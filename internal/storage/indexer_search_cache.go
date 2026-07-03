package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type IndexerSearchHistoryInput struct {
	IndexerID   uuid.UUID
	IndexerName string
	IndexerType string
	MediaType   string
	Query       string
	CacheHit    bool
	Success     bool
	ResultCount int32
	Error       *string
	Response    any
}

func (s *SettingsStore) GetIndexerSearchSettings(ctx context.Context) (IndexerSearchSettings, error) {
	if err := s.ensureIndexerSearchSettings(ctx); err != nil {
		return IndexerSearchSettings{}, err
	}
	var settings IndexerSearchSettings
	err := s.pool.QueryRow(ctx, `
		select cache_duration_minutes, history_retention_days, automatic_blocklist_expiry_days
		from app.indexer_search_settings
		where id = true
	`).Scan(
		&settings.CacheDurationMinutes,
		&settings.HistoryRetentionDays,
		&settings.AutomaticBlocklistExpiryDays,
	)
	return settings, err
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
	_, err := s.pool.Exec(ctx, `
		insert into app.indexer_search_settings (
			id, cache_duration_minutes, history_retention_days, automatic_blocklist_expiry_days
		)
		values (true, $1, $2, $3)
		on conflict (id) do update
		set cache_duration_minutes = excluded.cache_duration_minutes,
			history_retention_days = excluded.history_retention_days,
			automatic_blocklist_expiry_days = excluded.automatic_blocklist_expiry_days,
			updated_at = now()
	`, input.CacheDurationMinutes, input.HistoryRetentionDays, input.AutomaticBlocklistExpiryDays)
	if err != nil {
		return IndexerSearchSettings{}, err
	}
	return s.GetIndexerSearchSettings(ctx)
}

func (s *SettingsStore) GetIndexerSearchCache(ctx context.Context, indexerID uuid.UUID, mediaType string, query string, target any) (bool, error) {
	var raw []byte
	err := s.pool.QueryRow(ctx, `
		select response
		from app.indexer_search_cache
		where indexer_id = $1 and media_type = $2 and query = $3 and expires_at > now()
	`, indexerID, mediaType, query).Scan(&raw)
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
	_, err = s.pool.Exec(ctx, `
		insert into app.indexer_search_cache (
			indexer_id, media_type, query, response, result_count, expires_at
		)
		values ($1, $2, $3, $4, $5, $6)
		on conflict (indexer_id, media_type, query) do update
		set response = excluded.response,
			result_count = excluded.result_count,
			expires_at = excluded.expires_at,
			updated_at = now()
	`, indexerID, mediaType, query, raw, resultCount, expiresAt)
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
	var entry IndexerSearchHistoryEntry
	err = s.pool.QueryRow(ctx, `
		insert into app.indexer_search_history (
			id, indexer_id, indexer_name, indexer_type, media_type, query, cache_hit,
			success, result_count, error, response
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		returning indexer_name, indexer_type, media_type, query, cache_hit, success,
			result_count, error, response::text, created_at
	`, uuid.New(), input.IndexerID, input.IndexerName, input.IndexerType, input.MediaType, input.Query,
		input.CacheHit, input.Success, input.ResultCount, input.Error, raw).Scan(
		&entry.IndexerName,
		&entry.IndexerType,
		&entry.MediaType,
		&entry.Query,
		&entry.CacheHit,
		&entry.Success,
		&entry.ResultCount,
		&entry.Error,
		&entry.Response,
		&entry.CreatedAt,
	)
	return entry, err
}

func (s *SettingsStore) CleanupIndexerSearchHistory(ctx context.Context, retentionDays int32) (int32, error) {
	tag, err := s.pool.Exec(ctx, `
		delete from app.indexer_search_history
		where created_at < now() - make_interval(days => $1::int)
	`, retentionDays)
	if err != nil {
		return 0, err
	}
	return int32(tag.RowsAffected()), nil
}

func (s *SettingsStore) ClearIndexerSearchHistory(ctx context.Context) (int32, error) {
	tag, err := s.pool.Exec(ctx, `delete from app.indexer_search_history`)
	if err != nil {
		return 0, err
	}
	return int32(tag.RowsAffected()), nil
}

func (s *SettingsStore) ClearIndexerSearchCache(ctx context.Context) (int32, error) {
	tag, err := s.pool.Exec(ctx, `delete from app.indexer_search_cache`)
	if err != nil {
		return 0, err
	}
	return int32(tag.RowsAffected()), nil
}

func (s *SettingsStore) ClearIndexerSearchCacheByPattern(ctx context.Context, pattern string) (int32, error) {
	tag, err := s.pool.Exec(ctx, `
		delete from app.indexer_search_cache
		where query ~* $1
	`, pattern)
	if err != nil {
		return 0, err
	}
	return int32(tag.RowsAffected()), nil
}

func (s *SettingsStore) DeleteIndexerSearchCacheEntry(ctx context.Context, indexerID uuid.UUID, mediaType string, query string) (int32, error) {
	tag, err := s.pool.Exec(ctx, `
		delete from app.indexer_search_cache
		where indexer_id = $1 and media_type = $2 and query = $3
	`, indexerID, mediaType, query)
	if err != nil {
		return 0, err
	}
	return int32(tag.RowsAffected()), nil
}

func (s *SettingsStore) ensureIndexerSearchSettings(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `
		insert into app.indexer_search_settings (
			id, cache_duration_minutes, history_retention_days, automatic_blocklist_expiry_days
		)
		values (true, 1440, 7, 7)
		on conflict (id) do nothing
	`)
	return err
}
