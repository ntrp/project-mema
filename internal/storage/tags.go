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

type Tag struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *SettingsStore) ListTags(ctx context.Context) ([]Tag, error) {
	rows, err := s.pool.Query(ctx, `
		select id, name, created_at, updated_at
		from app.tags
		order by lower(name)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []Tag{}
	for rows.Next() {
		tag, err := scanTag(rows)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

func (s *SettingsStore) SaveTag(ctx context.Context, id *uuid.UUID, name string) (Tag, error) {
	cleanName := normalizeTagName(name)
	if cleanName == "" {
		return Tag{}, ErrInvalidInput
	}
	if id == nil {
		tag, err := scanTag(s.pool.QueryRow(ctx, `
			insert into app.tags (id, name)
			values ($1, $2)
			on conflict (lower(name)) do update
			set name = excluded.name, updated_at = now()
			returning id, name, created_at, updated_at
		`, uuid.New(), cleanName))
		return tag, normalizeTagWriteError(err)
	}
	tag, err := scanTag(s.pool.QueryRow(ctx, `
		update app.tags
		set name = $2, updated_at = now()
		where id = $1
		returning id, name, created_at, updated_at
	`, *id, cleanName))
	return tag, normalizeTagWriteError(err)
}

func (s *SettingsStore) DeleteTag(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `delete from app.tags where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func assignMediaItemTags(ctx context.Context, q mediaItemQuerier, mediaItemID uuid.UUID, tagNames []string) error {
	if _, err := q.Exec(ctx, `delete from app.media_item_tags where media_item_id = $1`, mediaItemID); err != nil {
		return err
	}
	for _, name := range normalizeTagNames(tagNames) {
		tagID, err := ensureTag(ctx, q, name)
		if err != nil {
			return err
		}
		if _, err := q.Exec(ctx, `
			insert into app.media_item_tags (media_item_id, tag_id)
			values ($1, $2)
			on conflict do nothing
		`, mediaItemID, tagID); err != nil {
			return err
		}
	}
	return nil
}

func assignMediaRequestTags(ctx context.Context, q mediaItemQuerier, requestID uuid.UUID, tagNames []string) error {
	if _, err := q.Exec(ctx, `delete from app.media_request_tags where media_request_id = $1`, requestID); err != nil {
		return err
	}
	for _, name := range normalizeTagNames(tagNames) {
		tagID, err := ensureTag(ctx, q, name)
		if err != nil {
			return err
		}
		if _, err := q.Exec(ctx, `
			insert into app.media_request_tags (media_request_id, tag_id)
			values ($1, $2)
			on conflict do nothing
		`, requestID, tagID); err != nil {
			return err
		}
	}
	return nil
}

func ensureTag(ctx context.Context, q mediaItemQuerier, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := q.QueryRow(ctx, `
		insert into app.tags (id, name)
		values ($1, $2)
		on conflict (lower(name)) do update
		set name = excluded.name, updated_at = now()
		returning id
	`, uuid.New(), name).Scan(&id)
	return id, err
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

func scanTag(row pgx.Row) (Tag, error) {
	var tag Tag
	err := row.Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	return tag, err
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
