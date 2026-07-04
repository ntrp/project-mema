package storage

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DiscoverBlacklistItem struct {
	ID               uuid.UUID
	Type             string
	Title            string
	Year             *int32
	ExternalProvider *string
	ExternalID       *string
	Overview         *string
	PosterPath       *string
	CreatedAt        time.Time
}

type DiscoverBlacklistInput struct {
	Type             string
	Title            string
	Year             *int32
	ExternalProvider *string
	ExternalID       *string
	Overview         *string
	PosterPath       *string
}

func (s *SettingsStore) ListDiscoverBlacklist(ctx context.Context) ([]DiscoverBlacklistItem, error) {
	if err := s.ensureDiscoverBlacklistSchema(ctx); err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		select id, media_type, title, year, external_provider, external_id, overview, poster_path, created_at
		from app.discover_blacklist
		order by created_at desc, lower(title)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []DiscoverBlacklistItem{}
	for rows.Next() {
		item, err := scanDiscoverBlacklistItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *SettingsStore) SaveDiscoverBlacklistItem(
	ctx context.Context,
	input DiscoverBlacklistInput,
) (DiscoverBlacklistItem, error) {
	if err := s.ensureDiscoverBlacklistSchema(ctx); err != nil {
		return DiscoverBlacklistItem{}, err
	}
	input.Title = strings.Join(strings.Fields(input.Title), " ")
	if input.Title == "" || (input.Type != "movie" && input.Type != "serie") {
		return DiscoverBlacklistItem{}, ErrInvalidInput
	}

	item, err := scanDiscoverBlacklistItem(s.pool.QueryRow(ctx, `
		insert into app.discover_blacklist (
			id, media_type, title, year, external_provider, external_id, overview, poster_path
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		on conflict (media_type, external_provider, external_id)
			where external_provider is not null and external_id is not null
		do update
		set title = excluded.title,
			year = excluded.year,
			overview = excluded.overview,
			poster_path = excluded.poster_path
		returning id, media_type, title, year, external_provider, external_id, overview, poster_path, created_at
	`,
		uuid.New(),
		input.Type,
		input.Title,
		input.Year,
		normalizedOptionalString(input.ExternalProvider),
		normalizedOptionalString(input.ExternalID),
		normalizedOptionalString(input.Overview),
		normalizedOptionalString(input.PosterPath),
	))
	if err == nil {
		return item, nil
	}
	if !isDiscoverBlacklistFallbackConflict(err) {
		return DiscoverBlacklistItem{}, normalizeDiscoverBlacklistWriteError(err)
	}
	return s.saveDiscoverBlacklistTitleFallback(ctx, input)
}

func (s *SettingsStore) DeleteDiscoverBlacklistItem(ctx context.Context, id uuid.UUID) error {
	if err := s.ensureDiscoverBlacklistSchema(ctx); err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `delete from app.discover_blacklist where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) saveDiscoverBlacklistTitleFallback(
	ctx context.Context,
	input DiscoverBlacklistInput,
) (DiscoverBlacklistItem, error) {
	item, err := scanDiscoverBlacklistItem(s.pool.QueryRow(ctx, `
		insert into app.discover_blacklist (
			id, media_type, title, year, external_provider, external_id, overview, poster_path
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		on conflict (media_type, lower(title), coalesce(year, 0))
			where external_provider is null or external_id is null
		do update
		set overview = excluded.overview,
			poster_path = excluded.poster_path
		returning id, media_type, title, year, external_provider, external_id, overview, poster_path, created_at
	`,
		uuid.New(),
		input.Type,
		input.Title,
		input.Year,
		nil,
		nil,
		normalizedOptionalString(input.Overview),
		normalizedOptionalString(input.PosterPath),
	))
	return item, normalizeDiscoverBlacklistWriteError(err)
}

func (s *SettingsStore) ensureDiscoverBlacklistSchema(ctx context.Context) error {
	if _, err := s.pool.Exec(ctx, `
		create table if not exists app.discover_blacklist (
			id uuid primary key,
			media_type text not null check (media_type in ('movie', 'series')),
			title text not null,
			year integer,
			external_provider text,
			external_id text,
			overview text,
			poster_path text,
			created_at timestamptz not null default now()
		)
	`); err != nil {
		return err
	}
	if _, err := s.pool.Exec(ctx, `
		create unique index if not exists idx_discover_blacklist_external
			on app.discover_blacklist (media_type, external_provider, external_id)
			where external_provider is not null and external_id is not null
	`); err != nil {
		return err
	}
	_, err := s.pool.Exec(ctx, `
		create unique index if not exists idx_discover_blacklist_title
			on app.discover_blacklist (media_type, lower(title), coalesce(year, 0))
			where external_provider is null or external_id is null
	`)
	return err
}

func scanDiscoverBlacklistItem(row pgx.Row) (DiscoverBlacklistItem, error) {
	var item DiscoverBlacklistItem
	err := row.Scan(
		&item.ID,
		&item.Type,
		&item.Title,
		&item.Year,
		&item.ExternalProvider,
		&item.ExternalID,
		&item.Overview,
		&item.PosterPath,
		&item.CreatedAt,
	)
	return item, err
}

func normalizedOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil
	}
	return &clean
}

func isDiscoverBlacklistFallbackConflict(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" &&
		pgErr.ConstraintName == "idx_discover_blacklist_title"
}

func normalizeDiscoverBlacklistWriteError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23505" || pgErr.Code == "23514") {
		return ErrInvalidInput
	}
	return err
}
