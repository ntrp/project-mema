package storage

import (
	"context"
	"errors"
	"time"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SubtitleProvider struct {
	ID            uuid.UUID
	Name          string
	Type          string
	BaseURL       string
	Username      *string
	Password      *string
	APIKey        *string
	Enabled       bool
	Priority      int32
	MockSubtitles []MockSubtitleProviderRow
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type SubtitleProviderInput struct {
	Name          string
	Type          string
	BaseURL       string
	Username      *string
	Password      *string
	APIKey        *string
	Enabled       bool
	Priority      int32
	MockSubtitles []MockSubtitleProviderRowInput
}

func (s *SettingsStore) ListSubtitleProviders(ctx context.Context) ([]SubtitleProvider, error) {
	rows, err := storagegen.New(s.pool).ListSubtitleProviders(ctx)
	if err != nil {
		return nil, err
	}
	providers := make([]SubtitleProvider, 0, len(rows))
	for _, row := range rows {
		provider, err := subtitleProviderWithRows(ctx, s.pool, row)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	return providers, nil
}

func (s *SettingsStore) GetSubtitleProvider(ctx context.Context, id uuid.UUID) (SubtitleProvider, error) {
	row, err := storagegen.New(s.pool).GetSubtitleProvider(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return SubtitleProvider{}, ErrNotFound
	}
	if err != nil {
		return SubtitleProvider{}, err
	}
	return subtitleProviderWithRows(ctx, s.pool, row)
}

func (s *SettingsStore) CreateSubtitleProvider(ctx context.Context, input SubtitleProviderInput) (SubtitleProvider, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return SubtitleProvider{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	q := storagegen.New(tx)
	row, err := q.CreateSubtitleProvider(ctx, subtitleProviderParams(uuid.New(), input))
	if err != nil {
		return SubtitleProvider{}, err
	}
	provider, err := subtitleProviderWithRows(ctx, tx, row)
	if err != nil {
		return SubtitleProvider{}, err
	}
	if provider.MockSubtitles, err = replaceMockSubtitleProviderRows(ctx, q, row.ID, input.MockSubtitles); err != nil {
		return SubtitleProvider{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return SubtitleProvider{}, err
	}
	return provider, nil
}

func (s *SettingsStore) UpdateSubtitleProvider(ctx context.Context, id uuid.UUID, input SubtitleProviderInput) (SubtitleProvider, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return SubtitleProvider{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	q := storagegen.New(tx)
	row, err := q.UpdateSubtitleProvider(ctx, subtitleProviderUpdateParams(id, input))
	if errors.Is(err, pgx.ErrNoRows) {
		return SubtitleProvider{}, ErrNotFound
	}
	if err != nil {
		return SubtitleProvider{}, err
	}
	mockRows, err := replaceMockSubtitleProviderRows(ctx, q, id, input.MockSubtitles)
	if err != nil {
		return SubtitleProvider{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return SubtitleProvider{}, err
	}
	provider := subtitleProviderFromRow(row)
	provider.MockSubtitles = mockRows
	return provider, nil
}

func (s *SettingsStore) DeleteSubtitleProvider(ctx context.Context, id uuid.UUID) error {
	rows, err := storagegen.New(s.pool).DeleteSubtitleProvider(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func subtitleProviderParams(id uuid.UUID, input SubtitleProviderInput) storagegen.CreateSubtitleProviderParams {
	return storagegen.CreateSubtitleProviderParams{
		ID:       id,
		Name:     input.Name,
		Type:     input.Type,
		BaseUrl:  input.BaseURL,
		Username: textValue(input.Username),
		Password: textValue(input.Password),
		ApiKey:   textValue(input.APIKey),
		Enabled:  input.Enabled,
		Priority: input.Priority,
	}
}

func subtitleProviderUpdateParams(id uuid.UUID, input SubtitleProviderInput) storagegen.UpdateSubtitleProviderParams {
	return storagegen.UpdateSubtitleProviderParams{
		ID:       id,
		Name:     input.Name,
		Type:     input.Type,
		BaseUrl:  input.BaseURL,
		Username: textValue(input.Username),
		Password: textValue(input.Password),
		ApiKey:   textValue(input.APIKey),
		Enabled:  input.Enabled,
		Priority: input.Priority,
	}
}

func subtitleProviderFromRow(row storagegen.AppSubtitleProvider) SubtitleProvider {
	return SubtitleProvider{
		ID:            row.ID,
		Name:          row.Name,
		Type:          row.Type,
		BaseURL:       row.BaseUrl,
		Username:      textPtr(row.Username),
		Password:      textPtr(row.Password),
		APIKey:        textPtr(row.ApiKey),
		Enabled:       row.Enabled,
		Priority:      row.Priority,
		MockSubtitles: []MockSubtitleProviderRow{},
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}
