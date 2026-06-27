package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

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

func (s *SettingsStore) GetDownloadClient(ctx context.Context, id uuid.UUID) (DownloadClient, error) {
	return scanDownloadClientRow(s.pool.QueryRow(ctx, `
		select id, name, type, base_url, username, password, api_key, category, enabled, priority, created_at, updated_at
		from app.download_clients
		where id = $1
	`, id))
}

func (s *SettingsStore) CreateDownloadClient(ctx context.Context, input DownloadClientInput) (DownloadClient, error) {
	id := uuid.New()
	return scanDownloadClientRow(s.pool.QueryRow(ctx, `
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
	return scanDownloadClientRow(s.pool.QueryRow(ctx, `
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

func scanDownloadClientRow(row pgx.Row) (DownloadClient, error) {
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
