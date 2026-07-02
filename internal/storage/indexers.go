package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Indexer struct {
	ID             uuid.UUID
	Name           string
	Type           string
	BaseURL        string
	APIKey         *string
	Categories     []int32
	Enabled        bool
	Priority       int32
	HealthStatus   string
	LastQueryAt    *time.Time
	LastSuccessAt  *time.Time
	LastFailureAt  *time.Time
	NextCheckAt    *time.Time
	LastStatusCode *int32
	LastError      *string
	FailureCount   int32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type IndexerInput struct {
	Name       string
	Type       string
	BaseURL    string
	APIKey     *string
	Categories []int32
	Enabled    bool
	Priority   int32
}

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
	id := uuid.New()
	return scanIndexerRow(s.pool.QueryRow(ctx, `
		insert into app.indexers (
			id, name, type, base_url, api_key, categories, enabled, priority
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		returning `+indexerColumns+`
	`,
		id, input.Name, input.Type, input.BaseURL, input.APIKey, input.Categories, input.Enabled, input.Priority,
	))
}

func (s *SettingsStore) UpdateIndexer(ctx context.Context, id uuid.UUID, input IndexerInput) (Indexer, error) {
	return scanIndexerRow(s.pool.QueryRow(ctx, `
		update app.indexers
		set name = $2,
			type = $3,
			base_url = $4,
			api_key = $5,
			categories = $6,
			enabled = $7,
			priority = $8,
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
		id, input.Name, input.Type, input.BaseURL, input.APIKey, input.Categories, input.Enabled, input.Priority,
	))
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

func scanIndexerRow(row pgx.Row) (Indexer, error) {
	indexer, err := scanIndexer(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Indexer{}, ErrNotFound
	}
	return indexer, err
}

func scanIndexer(row pgx.Row) (Indexer, error) {
	var indexer Indexer
	err := row.Scan(
		&indexer.ID,
		&indexer.Name,
		&indexer.Type,
		&indexer.BaseURL,
		&indexer.APIKey,
		&indexer.Categories,
		&indexer.Enabled,
		&indexer.Priority,
		&indexer.HealthStatus,
		&indexer.LastQueryAt,
		&indexer.LastSuccessAt,
		&indexer.LastFailureAt,
		&indexer.NextCheckAt,
		&indexer.LastStatusCode,
		&indexer.LastError,
		&indexer.FailureCount,
		&indexer.CreatedAt,
		&indexer.UpdatedAt,
	)
	return indexer, err
}

const indexerColumns = `
	id, name, type, base_url, api_key, categories, enabled, priority,
	health_status, last_query_at, last_success_at, last_failure_at, next_check_at,
	last_status_code, last_error, failure_count, created_at, updated_at
`
