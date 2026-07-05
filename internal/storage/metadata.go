package storage

import (
	"context"
	"errors"
	"time"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MetadataProvider struct {
	ID                    uuid.UUID
	Name                  string
	Type                  string
	BaseURL               string
	APIKey                *string
	PIN                   *string
	AccessToken           *string
	SessionToken          *string
	SessionTokenExpiresAt *time.Time
	Enabled               bool
	Priority              int32
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type MetadataProviderInput struct {
	Name        string
	Type        string
	BaseURL     string
	APIKey      *string
	PIN         *string
	AccessToken *string
	Enabled     bool
	Priority    int32
}

type MetadataCacheStats struct {
	TotalEntries   int32
	ActiveEntries  int32
	ExpiredEntries int32
	ProviderCount  int32
}

type MetadataCacheEntry struct {
	ProviderID   uuid.UUID
	ProviderName string
	ProviderType string
	MediaType    string
	Query        string
	Year         int32
	ItemCount    int32
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Expired      bool
}

type MetadataSearchHistoryEntry struct {
	ProviderName string
	ProviderType string
	MediaType    string
	Query        string
	CacheKind    string
	Year         int32
	CacheHit     bool
	Success      bool
	ItemCount    int32
	Error        *string
	Response     string
	CreatedAt    time.Time
}

func (s *SettingsStore) EnsureDefaultMetadataProviders(ctx context.Context) error {
	defaults := []MetadataProviderInput{
		{
			Name:     "TMDB",
			Type:     "tmdb",
			BaseURL:  "https://api.themoviedb.org/3",
			Enabled:  true,
			Priority: 100,
		},
		{
			Name:     "TVDB",
			Type:     "tvdb",
			BaseURL:  "https://api4.thetvdb.com/v4",
			Enabled:  true,
			Priority: 110,
		},
	}
	for _, input := range defaults {
		if err := s.ensureMetadataProvider(ctx, input); err != nil {
			return err
		}
	}
	return nil
}

func (s *SettingsStore) ensureMetadataProvider(ctx context.Context, input MetadataProviderInput) error {
	return storagegen.New(s.pool).EnsureMetadataProvider(ctx, storagegen.EnsureMetadataProviderParams{
		ID:       uuid.New(),
		Name:     input.Name,
		Type:     input.Type,
		BaseUrl:  input.BaseURL,
		Enabled:  input.Enabled,
		Priority: input.Priority,
	})
}

func (s *SettingsStore) ListMetadataProviders(ctx context.Context) ([]MetadataProvider, error) {
	rows, err := storagegen.New(s.pool).ListMetadataProviders(ctx)
	if err != nil {
		return nil, err
	}
	providers := make([]MetadataProvider, 0, len(rows))
	for _, row := range rows {
		providers = append(providers, metadataProviderFromRow(row))
	}
	return providers, nil
}

func (s *SettingsStore) ListEnabledMetadataProviders(ctx context.Context, mediaType string) ([]MetadataProvider, error) {
	rows, err := storagegen.New(s.pool).ListEnabledMetadataProviders(ctx, mediaType)
	if err != nil {
		return nil, err
	}
	providers := make([]MetadataProvider, 0, len(rows))
	for _, row := range rows {
		providers = append(providers, metadataProviderFromRow(row))
	}
	return providers, nil
}

func (s *SettingsStore) GetMetadataProvider(ctx context.Context, id uuid.UUID) (MetadataProvider, error) {
	row, err := storagegen.New(s.pool).GetMetadataProvider(ctx, id)
	return metadataProviderResult(row, err)
}

func (s *SettingsStore) CreateMetadataProvider(ctx context.Context, input MetadataProviderInput) (MetadataProvider, error) {
	row, err := storagegen.New(s.pool).CreateMetadataProvider(ctx, metadataProviderCreateParams(uuid.New(), input))
	return metadataProviderResult(row, err)
}

func (s *SettingsStore) UpdateMetadataProvider(ctx context.Context, id uuid.UUID, input MetadataProviderInput) (MetadataProvider, error) {
	row, err := storagegen.New(s.pool).UpdateMetadataProvider(ctx, storagegen.UpdateMetadataProviderParams(metadataProviderCreateParams(id, input)))
	return metadataProviderResult(row, err)
}

func (s *SettingsStore) DeleteMetadataProvider(ctx context.Context, id uuid.UUID) error {
	rows, err := storagegen.New(s.pool).DeleteMetadataProvider(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) UpdateMetadataProviderSessionToken(ctx context.Context, id uuid.UUID, token string, expiresAt time.Time) error {
	return storagegen.New(s.pool).UpdateMetadataProviderSessionToken(ctx, storagegen.UpdateMetadataProviderSessionTokenParams{
		ID:                    id,
		SessionToken:          textValue(&token),
		SessionTokenExpiresAt: &expiresAt,
	})
}

func metadataProviderResult(row storagegen.AppMetadataProvider, err error) (MetadataProvider, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return MetadataProvider{}, ErrNotFound
	}
	return metadataProviderFromRow(row), err
}

func metadataProviderCreateParams(id uuid.UUID, input MetadataProviderInput) storagegen.CreateMetadataProviderParams {
	return storagegen.CreateMetadataProviderParams{
		ID:          id,
		Name:        input.Name,
		Type:        input.Type,
		BaseUrl:     input.BaseURL,
		ApiKey:      textValue(input.APIKey),
		Pin:         textValue(input.PIN),
		AccessToken: textValue(input.AccessToken),
		Enabled:     input.Enabled,
		Priority:    input.Priority,
	}
}

func metadataProviderFromRow(row storagegen.AppMetadataProvider) MetadataProvider {
	return MetadataProvider{
		ID:                    row.ID,
		Name:                  row.Name,
		Type:                  row.Type,
		BaseURL:               row.BaseUrl,
		APIKey:                textPtr(row.ApiKey),
		PIN:                   textPtr(row.Pin),
		AccessToken:           textPtr(row.AccessToken),
		SessionToken:          textPtr(row.SessionToken),
		SessionTokenExpiresAt: row.SessionTokenExpiresAt,
		Enabled:               row.Enabled,
		Priority:              row.Priority,
		CreatedAt:             row.CreatedAt,
		UpdatedAt:             row.UpdatedAt,
	}
}
