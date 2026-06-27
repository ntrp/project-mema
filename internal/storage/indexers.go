package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Indexer struct {
	ID         uuid.UUID
	Name       string
	Type       string
	BaseURL    string
	APIKey     *string
	Categories []int32
	Enabled    bool
	Priority   int32
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
		select id, name, type, base_url, api_key, categories, enabled, priority, created_at, updated_at
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

func (s *SettingsStore) CreateIndexer(ctx context.Context, input IndexerInput) (Indexer, error) {
	id := uuid.New()
	return scanIndexerRow(s.pool.QueryRow(ctx, `
		insert into app.indexers (
			id, name, type, base_url, api_key, categories, enabled, priority
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		returning id, name, type, base_url, api_key, categories, enabled, priority, created_at, updated_at
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
			updated_at = now()
		where id = $1
		returning id, name, type, base_url, api_key, categories, enabled, priority, created_at, updated_at
	`,
		id, input.Name, input.Type, input.BaseURL, input.APIKey, input.Categories, input.Enabled, input.Priority,
	))
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
		&indexer.CreatedAt,
		&indexer.UpdatedAt,
	)
	return indexer, err
}
