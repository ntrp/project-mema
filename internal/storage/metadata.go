package storage

import (
	"context"
	"errors"
	"time"

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
	id := uuid.New()
	_, err := s.pool.Exec(ctx, `
		insert into app.metadata_providers (
			id, name, type, base_url, enabled, priority
		)
		select $1, $2, $3, $4, $5, $6
		where not exists (
			select 1 from app.metadata_providers where type = $3
		)
	`, id, input.Name, input.Type, input.BaseURL, input.Enabled, input.Priority)
	return err
}

func (s *SettingsStore) ListMetadataProviders(ctx context.Context) ([]MetadataProvider, error) {
	rows, err := s.pool.Query(ctx, `
		select id, name, type, base_url, api_key, pin, access_token, session_token,
			session_token_expires_at, enabled, priority, created_at, updated_at
		from app.metadata_providers
		order by priority asc, name asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providers := []MetadataProvider{}
	for rows.Next() {
		provider, err := scanMetadataProvider(rows)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	return providers, rows.Err()
}

func (s *SettingsStore) ListEnabledMetadataProviders(ctx context.Context, mediaType string) ([]MetadataProvider, error) {
	rows, err := s.pool.Query(ctx, `
		select id, name, type, base_url, api_key, pin, access_token, session_token,
			session_token_expires_at, enabled, priority, created_at, updated_at
		from app.metadata_providers
		where enabled = true
			and (($1 = 'movie' and type in ('tmdb', 'tvdb')) or ($1 = 'series' and type in ('tmdb', 'tvdb')))
		order by priority asc, name asc
	`, mediaType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providers := []MetadataProvider{}
	for rows.Next() {
		provider, err := scanMetadataProvider(rows)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	return providers, rows.Err()
}

func (s *SettingsStore) GetMetadataProvider(ctx context.Context, id uuid.UUID) (MetadataProvider, error) {
	provider, err := scanMetadataProvider(s.pool.QueryRow(ctx, `
		select id, name, type, base_url, api_key, pin, access_token, session_token,
			session_token_expires_at, enabled, priority, created_at, updated_at
		from app.metadata_providers
		where id = $1
	`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return MetadataProvider{}, ErrNotFound
	}
	return provider, err
}

func (s *SettingsStore) CreateMetadataProvider(ctx context.Context, input MetadataProviderInput) (MetadataProvider, error) {
	id := uuid.New()
	return scanMetadataProvider(s.pool.QueryRow(ctx, `
		insert into app.metadata_providers (
			id, name, type, base_url, api_key, pin, access_token, enabled, priority
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		returning id, name, type, base_url, api_key, pin, access_token, session_token,
			session_token_expires_at, enabled, priority, created_at, updated_at
	`, id, input.Name, input.Type, input.BaseURL, input.APIKey, input.PIN, input.AccessToken, input.Enabled, input.Priority))
}

func (s *SettingsStore) UpdateMetadataProvider(ctx context.Context, id uuid.UUID, input MetadataProviderInput) (MetadataProvider, error) {
	provider, err := scanMetadataProvider(s.pool.QueryRow(ctx, `
		update app.metadata_providers
		set name = $2,
			type = $3,
			base_url = $4,
			api_key = $5,
			pin = $6,
			access_token = $7,
			session_token = null,
			session_token_expires_at = null,
			enabled = $8,
			priority = $9,
			updated_at = now()
		where id = $1
		returning id, name, type, base_url, api_key, pin, access_token, session_token,
			session_token_expires_at, enabled, priority, created_at, updated_at
	`, id, input.Name, input.Type, input.BaseURL, input.APIKey, input.PIN, input.AccessToken, input.Enabled, input.Priority))
	if errors.Is(err, pgx.ErrNoRows) {
		return MetadataProvider{}, ErrNotFound
	}
	return provider, err
}

func (s *SettingsStore) DeleteMetadataProvider(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `delete from app.metadata_providers where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) UpdateMetadataProviderSessionToken(ctx context.Context, id uuid.UUID, token string, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx, `
		update app.metadata_providers
		set session_token = $2, session_token_expires_at = $3, updated_at = now()
		where id = $1
	`, id, token, expiresAt)
	return err
}

func scanMetadataProvider(row pgx.Row) (MetadataProvider, error) {
	var provider MetadataProvider
	err := row.Scan(
		&provider.ID,
		&provider.Name,
		&provider.Type,
		&provider.BaseURL,
		&provider.APIKey,
		&provider.PIN,
		&provider.AccessToken,
		&provider.SessionToken,
		&provider.SessionTokenExpiresAt,
		&provider.Enabled,
		&provider.Priority,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)
	return provider, err
}
