package storage

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	storagegen "media-manager/internal/storage/generated"
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
	rows, err := storagegen.New(s.pool).ListPathMappings(ctx)
	if err != nil {
		return nil, err
	}

	mappings := make([]PathMapping, 0, len(rows))
	for _, row := range rows {
		mappings = append(mappings, pathMappingFromRow(row))
	}
	return mappings, nil
}

func (s *SettingsStore) CreatePathMapping(ctx context.Context, input PathMappingInput) (PathMapping, error) {
	row, err := storagegen.New(s.pool).UpsertPathMapping(ctx, storagegen.UpsertPathMappingParams{
		ID:         uuid.New(),
		ClientPath: normalizePathMappingPath(input.ClientPath),
		AppPath:    normalizePathMappingPath(input.AppPath),
	})
	if err != nil {
		return PathMapping{}, err
	}
	return pathMappingFromRow(row), nil
}

func (s *SettingsStore) DeletePathMapping(ctx context.Context, id uuid.UUID) error {
	rowsAffected, err := storagegen.New(s.pool).DeletePathMapping(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func pathMappingFromRow(row storagegen.AppPathMapping) PathMapping {
	return PathMapping{
		ID:         row.ID,
		ClientPath: row.ClientPath,
		AppPath:    row.AppPath,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}

func normalizePathMappingPath(value string) string {
	value = strings.TrimSpace(value)
	return strings.TrimRight(value, "/")
}
