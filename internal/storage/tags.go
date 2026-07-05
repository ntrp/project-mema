package storage

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	storagegen "media-manager/internal/storage/generated"
)

type Tag struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *SettingsStore) ListTags(ctx context.Context) ([]Tag, error) {
	rows, err := storagegen.New(s.pool).ListTags(ctx)
	if err != nil {
		return nil, err
	}

	tags := make([]Tag, 0, len(rows))
	for _, row := range rows {
		tags = append(tags, tagFromRow(row))
	}
	return tags, nil
}

func (s *SettingsStore) SaveTag(ctx context.Context, id *uuid.UUID, name string) (Tag, error) {
	cleanName := normalizeTagName(name)
	if cleanName == "" {
		return Tag{}, ErrInvalidInput
	}
	if id == nil {
		row, err := storagegen.New(s.pool).UpsertTagByName(ctx, storagegen.UpsertTagByNameParams{
			ID:   uuid.New(),
			Name: cleanName,
		})
		if err != nil {
			return Tag{}, normalizeTagWriteError(err)
		}
		return tagFromRow(row), nil
	}
	row, err := storagegen.New(s.pool).UpdateTag(ctx, storagegen.UpdateTagParams{
		ID:   *id,
		Name: cleanName,
	})
	if err != nil {
		return Tag{}, normalizeTagWriteError(err)
	}
	return tagFromRow(row), nil
}

func (s *SettingsStore) DeleteTag(ctx context.Context, id uuid.UUID) error {
	rowsAffected, err := storagegen.New(s.pool).DeleteTag(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func assignMediaItemTags(ctx context.Context, q mediaItemQuerier, mediaItemID uuid.UUID, tagNames []string) error {
	queries := storagegen.New(q)
	if err := queries.DeleteMediaItemTags(ctx, mediaItemID); err != nil {
		return err
	}
	for _, name := range normalizeTagNames(tagNames) {
		tagID, err := ensureTag(ctx, queries, name)
		if err != nil {
			return err
		}
		if err := queries.AddMediaItemTag(ctx, storagegen.AddMediaItemTagParams{
			MediaItemID: mediaItemID,
			TagID:       tagID,
		}); err != nil {
			return err
		}
	}
	return nil
}

func assignMediaRequestTags(ctx context.Context, q mediaItemQuerier, requestID uuid.UUID, tagNames []string) error {
	queries := storagegen.New(q)
	if err := queries.DeleteMediaRequestTags(ctx, requestID); err != nil {
		return err
	}
	for _, name := range normalizeTagNames(tagNames) {
		tagID, err := ensureTag(ctx, queries, name)
		if err != nil {
			return err
		}
		if err := queries.AddMediaRequestTag(ctx, storagegen.AddMediaRequestTagParams{
			MediaRequestID: requestID,
			TagID:          tagID,
		}); err != nil {
			return err
		}
	}
	return nil
}

func ensureTag(ctx context.Context, queries *storagegen.Queries, name string) (uuid.UUID, error) {
	return queries.EnsureTag(ctx, storagegen.EnsureTagParams{
		ID:   uuid.New(),
		Name: name,
	})
}

func normalizeTagNames(values []string) []string {
	seen := map[string]struct{}{}
	tags := []string{}
	for _, value := range values {
		name := normalizeTagName(value)
		if name == "" {
			continue
		}
		key := strings.ToLower(name)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		tags = append(tags, name)
	}
	return tags
}

func normalizeTagName(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func tagFromRow(row storagegen.AppTag) Tag {
	return Tag{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func normalizeTagWriteError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrInvalidInput
	}
	return err
}
