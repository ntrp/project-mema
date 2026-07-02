package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *SettingsStore) GetMetadataSearchCache(ctx context.Context, providerID uuid.UUID, mediaType string, query string, year *int32, target any) (bool, error) {
	var raw []byte
	err := s.pool.QueryRow(ctx, `
		select results
		from app.metadata_search_cache
		where provider_id = $1 and media_type = $2 and query = $3 and year = $4 and expires_at > now()
	`, providerID, mediaType, query, cacheYear(year)).Scan(&raw)
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
	_, err = s.pool.Exec(ctx, `
		insert into app.metadata_search_cache (
			provider_id, media_type, query, year, results, expires_at
		)
		values ($1, $2, $3, $4, $5, $6)
		on conflict (provider_id, media_type, query, year) do update
		set results = excluded.results, expires_at = excluded.expires_at, updated_at = now()
	`, providerID, mediaType, query, cacheYear(year), raw, expiresAt)
	return err
}

func (s *SettingsStore) RecordMetadataSearchHistory(ctx context.Context, input MetadataSearchHistoryInput) (MetadataSearchHistoryEntry, error) {
	raw, err := json.Marshal(input.Response)
	if err != nil {
		return MetadataSearchHistoryEntry{}, err
	}
	var entry MetadataSearchHistoryEntry
	err = s.pool.QueryRow(ctx, `
		insert into app.metadata_search_history (
			id, provider_id, provider_name, provider_type, media_type, query, year,
			cache_hit, success, item_count, error, response
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		returning provider_name, provider_type, media_type, query, year, cache_hit,
			success, item_count, error, response::text, created_at
	`, uuid.New(), input.ProviderID, input.ProviderName, input.ProviderType, input.MediaType, input.Query,
		cacheYear(input.Year), input.CacheHit, input.Success, input.ItemCount, input.Error, raw).Scan(
		&entry.ProviderName,
		&entry.ProviderType,
		&entry.MediaType,
		&entry.Query,
		&entry.Year,
		&entry.CacheHit,
		&entry.Success,
		&entry.ItemCount,
		&entry.Error,
		&entry.Response,
		&entry.CreatedAt,
	)
	if err != nil {
		return MetadataSearchHistoryEntry{}, err
	}
	entry.CacheKind = metadataCacheKind(entry.Query)
	return entry, nil
}

func (s *SettingsStore) MetadataCacheStats(ctx context.Context) (MetadataCacheStats, error) {
	var stats MetadataCacheStats
	err := s.pool.QueryRow(ctx, `
		select
			count(*)::int,
			count(*) filter (where expires_at > now())::int,
			count(*) filter (where expires_at <= now())::int,
			count(distinct provider_id)::int
		from app.metadata_search_cache
	`).Scan(&stats.TotalEntries, &stats.ActiveEntries, &stats.ExpiredEntries, &stats.ProviderCount)
	return stats, err
}

func (s *SettingsStore) ListMetadataSearchHistoryEntries(ctx context.Context, limit int32) ([]MetadataSearchHistoryEntry, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	rows, err := s.pool.Query(ctx, `
		select provider_name, provider_type, media_type, query, year, cache_hit, success,
			item_count, error, response::text, created_at
		from app.metadata_search_history
		order by created_at desc
		limit $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []MetadataSearchHistoryEntry{}
	for rows.Next() {
		var entry MetadataSearchHistoryEntry
		if err := rows.Scan(
			&entry.ProviderName,
			&entry.ProviderType,
			&entry.MediaType,
			&entry.Query,
			&entry.Year,
			&entry.CacheHit,
			&entry.Success,
			&entry.ItemCount,
			&entry.Error,
			&entry.Response,
			&entry.CreatedAt,
		); err != nil {
			return nil, err
		}
		entry.CacheKind = metadataCacheKind(entry.Query)
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (s *SettingsStore) ListMetadataCacheEntries(ctx context.Context, limit int32) ([]MetadataCacheEntry, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	rows, err := s.pool.Query(ctx, `
		select p.name,
			p.type,
			c.media_type,
			c.query,
			c.year,
			case
				when jsonb_typeof(c.results) = 'array' then jsonb_array_length(c.results)
				else 1
			end::int,
			c.expires_at,
			c.created_at,
			c.updated_at,
			c.expires_at <= now()
		from app.metadata_search_cache c
		join app.metadata_providers p on p.id = c.provider_id
		order by c.updated_at desc
		limit $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []MetadataCacheEntry{}
	for rows.Next() {
		var entry MetadataCacheEntry
		if err := rows.Scan(
			&entry.ProviderName,
			&entry.ProviderType,
			&entry.MediaType,
			&entry.Query,
			&entry.Year,
			&entry.ItemCount,
			&entry.ExpiresAt,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&entry.Expired,
		); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (s *SettingsStore) ClearMetadataCache(ctx context.Context) (int32, error) {
	tag, err := s.pool.Exec(ctx, `delete from app.metadata_search_cache`)
	if err != nil {
		return 0, err
	}
	return int32(tag.RowsAffected()), nil
}

func (s *SettingsStore) ClearMetadataCacheByPattern(ctx context.Context, pattern string) (int32, error) {
	tag, err := s.pool.Exec(ctx, `
		delete from app.metadata_search_cache
		where query ~* $1
	`, pattern)
	if err != nil {
		return 0, err
	}
	return int32(tag.RowsAffected()), nil
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
