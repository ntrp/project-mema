package storage

import (
	"context"

	"github.com/google/uuid"
)

func (s *SettingsStore) IndexerSearchCacheStats(ctx context.Context) (IndexerSearchCacheStats, error) {
	var stats IndexerSearchCacheStats
	err := s.pool.QueryRow(ctx, `
		select
			count(*)::int,
			count(*) filter (where expires_at > now())::int,
			count(*) filter (where expires_at <= now())::int,
			count(distinct indexer_id)::int
		from app.indexer_search_cache
	`).Scan(&stats.TotalEntries, &stats.ActiveEntries, &stats.ExpiredEntries, &stats.IndexerCount)
	return stats, err
}

func (s *SettingsStore) ListIndexerSearchCacheEntries(ctx context.Context, limit int32) ([]IndexerSearchCacheEntry, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	rows, err := s.pool.Query(ctx, `
		select i.name,
			i.type,
			c.media_type,
			c.query,
			c.result_count,
			c.expires_at,
			c.created_at,
			c.updated_at,
			c.expires_at <= now()
		from app.indexer_search_cache c
		join app.indexers i on i.id = c.indexer_id
		order by c.updated_at desc
		limit $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []IndexerSearchCacheEntry{}
	for rows.Next() {
		var entry IndexerSearchCacheEntry
		if err := rows.Scan(
			&entry.IndexerName,
			&entry.IndexerType,
			&entry.MediaType,
			&entry.Query,
			&entry.ResultCount,
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

func (s *SettingsStore) GetIndexerSearchCacheEntry(
	ctx context.Context,
	indexerID uuid.UUID,
	mediaType string,
	query string,
) (IndexerSearchCacheEntry, error) {
	var entry IndexerSearchCacheEntry
	err := s.pool.QueryRow(ctx, `
		select i.name,
			i.type,
			c.media_type,
			c.query,
			c.result_count,
			c.expires_at,
			c.created_at,
			c.updated_at,
			c.expires_at <= now()
		from app.indexer_search_cache c
		join app.indexers i on i.id = c.indexer_id
		where c.indexer_id = $1 and c.media_type = $2 and c.query = $3
	`, indexerID, mediaType, query).Scan(
		&entry.IndexerName,
		&entry.IndexerType,
		&entry.MediaType,
		&entry.Query,
		&entry.ResultCount,
		&entry.ExpiresAt,
		&entry.CreatedAt,
		&entry.UpdatedAt,
		&entry.Expired,
	)
	return entry, err
}

func (s *SettingsStore) ListIndexerSearchHistoryEntries(ctx context.Context, limit int32) ([]IndexerSearchHistoryEntry, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	rows, err := s.pool.Query(ctx, `
		select indexer_name, indexer_type, media_type, query, cache_hit, success,
			result_count, error, response::text, created_at
		from app.indexer_search_history
		order by created_at desc
		limit $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []IndexerSearchHistoryEntry{}
	for rows.Next() {
		var entry IndexerSearchHistoryEntry
		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}
