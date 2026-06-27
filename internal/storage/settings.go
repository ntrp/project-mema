package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("resource not found")

type DownloadClient struct {
	ID        uuid.UUID
	Name      string
	Type      string
	BaseURL   string
	Username  *string
	Password  *string
	APIKey    *string
	Category  *string
	Enabled   bool
	Priority  int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DownloadClientInput struct {
	Name     string
	Type     string
	BaseURL  string
	Username *string
	Password *string
	APIKey   *string
	Category *string
	Enabled  bool
	Priority int32
}

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

type SettingsStore struct {
	pool *pgxpool.Pool
}

func NewSettingsStore(pool *pgxpool.Pool) *SettingsStore {
	return &SettingsStore{pool: pool}
}

func (s *SettingsStore) ListDownloadClients(ctx context.Context) ([]DownloadClient, error) {
	rows, err := s.pool.Query(ctx, `
		select id, name, type, base_url, username, password, api_key, category, enabled, priority, created_at, updated_at
		from app.download_clients
		order by priority asc, name asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := []DownloadClient{}
	for rows.Next() {
		client, err := scanDownloadClient(rows)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, rows.Err()
}

func (s *SettingsStore) CreateDownloadClient(ctx context.Context, input DownloadClientInput) (DownloadClient, error) {
	id := uuid.New()
	return s.scanDownloadClientRow(s.pool.QueryRow(ctx, `
		insert into app.download_clients (
			id, name, type, base_url, username, password, api_key, category, enabled, priority
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		returning id, name, type, base_url, username, password, api_key, category, enabled, priority, created_at, updated_at
	`,
		id, input.Name, input.Type, input.BaseURL, input.Username, input.Password, input.APIKey, input.Category, input.Enabled, input.Priority,
	))
}

func (s *SettingsStore) UpdateDownloadClient(ctx context.Context, id uuid.UUID, input DownloadClientInput) (DownloadClient, error) {
	return s.scanDownloadClientRow(s.pool.QueryRow(ctx, `
		update app.download_clients
		set name = $2,
			type = $3,
			base_url = $4,
			username = $5,
			password = $6,
			api_key = $7,
			category = $8,
			enabled = $9,
			priority = $10,
			updated_at = now()
		where id = $1
		returning id, name, type, base_url, username, password, api_key, category, enabled, priority, created_at, updated_at
	`,
		id, input.Name, input.Type, input.BaseURL, input.Username, input.Password, input.APIKey, input.Category, input.Enabled, input.Priority,
	))
}

func (s *SettingsStore) DeleteDownloadClient(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `delete from app.download_clients where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
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
	return s.scanIndexerRow(s.pool.QueryRow(ctx, `
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
	return s.scanIndexerRow(s.pool.QueryRow(ctx, `
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

func (s *SettingsStore) scanDownloadClientRow(row pgx.Row) (DownloadClient, error) {
	client, err := scanDownloadClient(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return DownloadClient{}, ErrNotFound
	}
	return client, err
}

func scanDownloadClient(row pgx.Row) (DownloadClient, error) {
	var client DownloadClient
	err := row.Scan(
		&client.ID,
		&client.Name,
		&client.Type,
		&client.BaseURL,
		&client.Username,
		&client.Password,
		&client.APIKey,
		&client.Category,
		&client.Enabled,
		&client.Priority,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	return client, err
}

func (s *SettingsStore) scanIndexerRow(row pgx.Row) (Indexer, error) {
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
