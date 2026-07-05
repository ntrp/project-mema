package storage

import (
	"context"
	"errors"
	"time"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type DownloadClient struct {
	ID        uuid.UUID
	Name      string
	Type      string
	Protocol  string
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
	Protocol string
	BaseURL  string
	Username *string
	Password *string
	APIKey   *string
	Category *string
	Enabled  bool
	Priority int32
}

func (s *SettingsStore) ListDownloadClients(ctx context.Context) ([]DownloadClient, error) {
	rows, err := storagegen.New(s.pool).ListDownloadClients(ctx)
	if err != nil {
		return nil, err
	}
	clients := make([]DownloadClient, 0, len(rows))
	for _, row := range rows {
		clients = append(clients, downloadClientFromRow(row))
	}
	return clients, nil
}

func (s *SettingsStore) ListEnabledDownloadClients(ctx context.Context) ([]DownloadClient, error) {
	rows, err := storagegen.New(s.pool).ListEnabledDownloadClients(ctx)
	if err != nil {
		return nil, err
	}
	clients := make([]DownloadClient, 0, len(rows))
	for _, row := range rows {
		clients = append(clients, downloadClientFromRow(row))
	}
	return clients, nil
}

func (s *SettingsStore) GetDownloadClient(ctx context.Context, id uuid.UUID) (DownloadClient, error) {
	row, err := storagegen.New(s.pool).GetDownloadClient(ctx, id)
	return downloadClientResult(row, err)
}

func (s *SettingsStore) CreateDownloadClient(ctx context.Context, input DownloadClientInput) (DownloadClient, error) {
	row, err := storagegen.New(s.pool).CreateDownloadClient(ctx, downloadClientCreateParams(uuid.New(), input))
	return downloadClientResult(row, err)
}

func (s *SettingsStore) UpdateDownloadClient(ctx context.Context, id uuid.UUID, input DownloadClientInput) (DownloadClient, error) {
	row, err := storagegen.New(s.pool).UpdateDownloadClient(ctx, storagegen.UpdateDownloadClientParams(downloadClientCreateParams(id, input)))
	return downloadClientResult(row, err)
}

func (s *SettingsStore) DeleteDownloadClient(ctx context.Context, id uuid.UUID) error {
	rows, err := storagegen.New(s.pool).DeleteDownloadClient(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func downloadClientResult(row storagegen.AppDownloadClient, err error) (DownloadClient, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return DownloadClient{}, ErrNotFound
	}
	return downloadClientFromRow(row), err
}

func downloadClientCreateParams(id uuid.UUID, input DownloadClientInput) storagegen.CreateDownloadClientParams {
	return storagegen.CreateDownloadClientParams{
		ID:       id,
		Name:     input.Name,
		Type:     input.Type,
		Protocol: input.Protocol,
		BaseUrl:  input.BaseURL,
		Username: textValue(input.Username),
		Password: textValue(input.Password),
		ApiKey:   textValue(input.APIKey),
		Category: textValue(input.Category),
		Enabled:  input.Enabled,
		Priority: input.Priority,
	}
}

func downloadClientFromRow(row storagegen.AppDownloadClient) DownloadClient {
	return DownloadClient{
		ID:        row.ID,
		Name:      row.Name,
		Type:      row.Type,
		Protocol:  row.Protocol,
		BaseURL:   row.BaseUrl,
		Username:  textPtr(row.Username),
		Password:  textPtr(row.Password),
		APIKey:    textPtr(row.ApiKey),
		Category:  textPtr(row.Category),
		Enabled:   row.Enabled,
		Priority:  row.Priority,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
