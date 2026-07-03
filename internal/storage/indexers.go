package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

func (s *SettingsStore) ListIndexers(ctx context.Context) ([]Indexer, error) {
	rows, err := s.pool.Query(ctx, `
		select `+indexerColumns+`
		from app.indexers
		order by priority asc, name asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indexers := []Indexer{}
	for rows.Next() {
		indexer, err := scanIndexer(rows)
		if err != nil {
			return nil, err
		}
		indexers = append(indexers, indexer)
	}
	return indexers, rows.Err()
}

func (s *SettingsStore) ListEnabledIndexers(ctx context.Context) ([]Indexer, error) {
	rows, err := s.pool.Query(ctx, `
		select `+indexerColumns+`
		from app.indexers
		where enabled = true
			and health_status <> 'disabled'
			and (next_check_at is null or next_check_at <= now())
		order by priority asc, name asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indexers := []Indexer{}
	for rows.Next() {
		indexer, err := scanIndexer(rows)
		if err != nil {
			return nil, err
		}
		indexers = append(indexers, indexer)
	}
	return indexers, rows.Err()
}

func (s *SettingsStore) GetIndexer(ctx context.Context, id uuid.UUID) (Indexer, error) {
	return scanIndexerRow(s.pool.QueryRow(ctx, `
		select `+indexerColumns+`
		from app.indexers
		where id = $1
	`, id))
}

func (s *SettingsStore) CreateIndexer(ctx context.Context, input IndexerInput) (Indexer, error) {
	input = normalizeIndexerInput(input)
	id := uuid.New()
	return scanIndexerRow(s.pool.QueryRow(ctx, `
		insert into app.indexers (
			id, definition_id, name, implementation, implementation_name, protocol, privacy,
			language, encoding, description, indexer_urls, legacy_urls, base_url, api_key,
			categories, fields, capabilities, redirect, app_profile_id, minimum_seeders,
			seed_ratio, seed_time, pack_seed_time, prefer_magnet_url, supports_rss,
			supports_search, supports_redirect, supports_pagination, enabled, priority
		)
		values (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
			$29, $30
		)
		returning `+indexerColumns+`
	`,
		id, input.DefinitionID, input.Name, input.Implementation, input.ImplementationName, input.Protocol,
		input.Privacy, input.Language, input.Encoding, input.Description, input.IndexerURLs, input.LegacyURLs,
		input.BaseURL, input.APIKey, input.Categories, input.Fields, input.Capabilities, input.Redirect,
		input.AppProfileID, input.MinimumSeeders, input.SeedRatio, input.SeedTime, input.PackSeedTime,
		input.PreferMagnetURL, input.SupportsRSS, input.SupportsSearch, input.SupportsRedirect,
		input.SupportsPagination, input.Enabled, input.Priority,
	))
}

func (s *SettingsStore) UpdateIndexer(ctx context.Context, id uuid.UUID, input IndexerInput) (Indexer, error) {
	input = normalizeIndexerInput(input)
	return scanIndexerRow(s.pool.QueryRow(ctx, `
		update app.indexers
		set definition_id = $2,
			name = $3,
			implementation = $4,
			implementation_name = $5,
			protocol = $6,
			privacy = $7,
			language = $8,
			encoding = $9,
			description = $10,
			indexer_urls = $11,
			legacy_urls = $12,
			base_url = $13,
			api_key = $14,
			categories = $15,
			fields = $16,
			capabilities = $17,
			redirect = $18,
			app_profile_id = $19,
			minimum_seeders = $20,
			seed_ratio = $21,
			seed_time = $22,
			pack_seed_time = $23,
			prefer_magnet_url = $24,
			supports_rss = $25,
			supports_search = $26,
			supports_redirect = $27,
			supports_pagination = $28,
			enabled = $29,
			priority = $30,
			health_status = 'healthy',
			last_query_at = null,
			last_success_at = null,
			last_failure_at = null,
			next_check_at = null,
			last_status_code = null,
			last_error = null,
			failure_count = 0,
			updated_at = now()
		where id = $1
		returning `+indexerColumns+`
	`,
		id, input.DefinitionID, input.Name, input.Implementation, input.ImplementationName, input.Protocol,
		input.Privacy, input.Language, input.Encoding, input.Description, input.IndexerURLs, input.LegacyURLs,
		input.BaseURL, input.APIKey, input.Categories, input.Fields, input.Capabilities, input.Redirect,
		input.AppProfileID, input.MinimumSeeders, input.SeedRatio, input.SeedTime, input.PackSeedTime,
		input.PreferMagnetURL, input.SupportsRSS, input.SupportsSearch, input.SupportsRedirect,
		input.SupportsPagination, input.Enabled, input.Priority,
	))
}

func normalizeIndexerInput(input IndexerInput) IndexerInput {
	if input.DefinitionID == "" {
		input.DefinitionID = "generic-torznab"
	}
	if input.Implementation == "" {
		input.Implementation = "Cardigann"
	}
	if input.ImplementationName == "" {
		input.ImplementationName = input.Name
	}
	if input.Protocol == "" {
		input.Protocol = "torrent"
	}
	if input.Privacy == "" {
		input.Privacy = "private"
	}
	if input.Language == "" {
		input.Language = "en-US"
	}
	if input.IndexerURLs == nil {
		input.IndexerURLs = []string{}
	}
	if input.LegacyURLs == nil {
		input.LegacyURLs = []string{}
	}
	if input.Categories == nil {
		input.Categories = []int32{}
	}
	if len(input.Fields) == 0 {
		input.Fields = json.RawMessage("[]")
	}
	if len(input.Capabilities) == 0 {
		input.Capabilities = json.RawMessage(`{"categories":[],"supportsRawSearch":true,"searchParams":["q"],"tvSearchParams":["q","season","ep"],"movieSearchParams":["q","imdbid"]}`)
	}
	if input.AppProfileID == "" {
		input.AppProfileID = "default"
	}
	if !input.SupportsRSS && !input.SupportsSearch && !input.SupportsRedirect && !input.SupportsPagination {
		input.SupportsRSS = true
		input.SupportsSearch = true
		input.SupportsRedirect = true
		input.SupportsPagination = true
	}
	return input
}

func (s *SettingsStore) RecordIndexerSuccess(ctx context.Context, id uuid.UUID) (Indexer, error) {
	return scanIndexerRow(s.pool.QueryRow(ctx, `
		update app.indexers
		set health_status = 'healthy',
			last_query_at = now(),
			last_success_at = now(),
			last_failure_at = null,
			next_check_at = null,
			last_status_code = null,
			last_error = null,
			failure_count = 0,
			updated_at = now()
		where id = $1
		returning `+indexerColumns+`
	`, id))
}

func (s *SettingsStore) RecordIndexerFailure(
	ctx context.Context,
	id uuid.UUID,
	statusCode *int32,
	message string,
	permanent bool,
	retryUntil *time.Time,
) (Indexer, error) {
	return scanIndexerRow(s.pool.QueryRow(ctx, `
		update app.indexers
		set health_status = case
				when $4 then 'disabled'
				when failure_count >= 5 then 'disabled'
				else 'temporary_disabled'
			end,
			last_query_at = now(),
			last_failure_at = now(),
			last_status_code = $2,
			last_error = $3,
			failure_count = failure_count + 1,
			next_check_at = case
				when $4 then null
				when failure_count >= 5 then null
				when $5::timestamptz is not null then $5
				when failure_count = 0 then now() + interval '1 minute'
				when failure_count = 1 then now() + interval '5 minutes'
				when failure_count = 2 then now() + interval '15 minutes'
				when failure_count = 3 then now() + interval '30 minutes'
				when failure_count = 4 then now() + interval '1 hour'
				else null
			end,
			updated_at = now()
		where id = $1
		returning `+indexerColumns+`
	`, id, statusCode, message, permanent, retryUntil))
}

func (s *SettingsStore) DeleteIndexer(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `delete from app.indexers where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
