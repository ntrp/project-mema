package storage

import (
	"context"
	"errors"
	"strings"
	"time"

	storagegen "media-manager/internal/storage/generated"

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
	rows, err := storagegen.New(s.pool).ListDiscoverBlacklist(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]DiscoverBlacklistItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, discoverBlacklistItemFromRow(row))
	}
	return items, nil
}

func (s *SettingsStore) SaveDiscoverBlacklistItem(
	ctx context.Context,
	input DiscoverBlacklistInput,
) (DiscoverBlacklistItem, error) {
	input.Title = strings.Join(strings.Fields(input.Title), " ")
	if input.Title == "" || (input.Type != "movie" && input.Type != "serie") {
		return DiscoverBlacklistItem{}, ErrInvalidInput
	}

	item, err := storagegen.New(s.pool).SaveDiscoverBlacklistByExternalID(ctx, discoverBlacklistExternalParams(input))
	if err == nil {
		return discoverBlacklistItemFromRow(item), nil
	}
	if !isDiscoverBlacklistFallbackConflict(err) {
		return DiscoverBlacklistItem{}, normalizeDiscoverBlacklistWriteError(err)
	}
	return s.saveDiscoverBlacklistTitleFallback(ctx, input)
}

func (s *SettingsStore) DeleteDiscoverBlacklistItem(ctx context.Context, id uuid.UUID) error {
	rows, err := storagegen.New(s.pool).DeleteDiscoverBlacklistItem(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) saveDiscoverBlacklistTitleFallback(
	ctx context.Context,
	input DiscoverBlacklistInput,
) (DiscoverBlacklistItem, error) {
	item, err := storagegen.New(s.pool).SaveDiscoverBlacklistByTitle(ctx, discoverBlacklistTitleParams(input))
	return discoverBlacklistItemFromRow(item), normalizeDiscoverBlacklistWriteError(err)
}

func discoverBlacklistExternalParams(input DiscoverBlacklistInput) storagegen.SaveDiscoverBlacklistByExternalIDParams {
	return storagegen.SaveDiscoverBlacklistByExternalIDParams{
		ID:               uuid.New(),
		MediaType:        input.Type,
		Title:            input.Title,
		Year:             int4Value(input.Year),
		ExternalProvider: textValue(normalizedOptionalString(input.ExternalProvider)),
		ExternalID:       textValue(normalizedOptionalString(input.ExternalID)),
		Overview:         textValue(normalizedOptionalString(input.Overview)),
		PosterPath:       textValue(normalizedOptionalString(input.PosterPath)),
	}
}

func discoverBlacklistTitleParams(input DiscoverBlacklistInput) storagegen.SaveDiscoverBlacklistByTitleParams {
	return storagegen.SaveDiscoverBlacklistByTitleParams{
		ID:         uuid.New(),
		MediaType:  input.Type,
		Title:      input.Title,
		Year:       int4Value(input.Year),
		Overview:   textValue(normalizedOptionalString(input.Overview)),
		PosterPath: textValue(normalizedOptionalString(input.PosterPath)),
	}
}

func discoverBlacklistItemFromRow(row storagegen.AppDiscoverBlacklist) DiscoverBlacklistItem {
	return DiscoverBlacklistItem{
		ID:               row.ID,
		Type:             row.MediaType,
		Title:            row.Title,
		Year:             int4Ptr(row.Year),
		ExternalProvider: textPtr(row.ExternalProvider),
		ExternalID:       textPtr(row.ExternalID),
		Overview:         textPtr(row.Overview),
		PosterPath:       textPtr(row.PosterPath),
		CreatedAt:        row.CreatedAt,
	}
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
