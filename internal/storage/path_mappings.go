package storage

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PathMapping struct {
	ID         uuid.UUID
	ClientPath string
	AppPath    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type PathMappingInput struct {
	ClientPath string
	AppPath    string
}

func (s *SettingsStore) ListPathMappings(ctx context.Context) ([]PathMapping, error) {
	rows, err := s.pool.Query(ctx, `
		select id, client_path, app_path, created_at, updated_at
		from app.path_mappings
		order by client_path asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mappings := []PathMapping{}
	for rows.Next() {
		mapping, err := scanPathMapping(rows)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, mapping)
	}
	return mappings, rows.Err()
}

func (s *SettingsStore) CreatePathMapping(ctx context.Context, input PathMappingInput) (PathMapping, error) {
	id := uuid.New()
	return scanPathMapping(s.pool.QueryRow(ctx, `
		insert into app.path_mappings (id, client_path, app_path)
		values ($1, $2, $3)
		on conflict (client_path) do update
		set app_path = excluded.app_path, updated_at = now()
		returning id, client_path, app_path, created_at, updated_at
	`, id, normalizePathMappingPath(input.ClientPath), normalizePathMappingPath(input.AppPath)))
}

func (s *SettingsStore) DeletePathMapping(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `delete from app.path_mappings where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func scanPathMapping(row pgx.Row) (PathMapping, error) {
	var mapping PathMapping
	err := row.Scan(&mapping.ID, &mapping.ClientPath, &mapping.AppPath, &mapping.CreatedAt, &mapping.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return PathMapping{}, ErrNotFound
	}
	return mapping, err
}

func normalizePathMappingPath(value string) string {
	value = strings.TrimSpace(value)
	return strings.TrimRight(value, "/")
}
