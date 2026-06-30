package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *SettingsStore) GetLibraryFolder(ctx context.Context, id uuid.UUID) (LibraryFolder, error) {
	folder, err := scanLibraryFolder(s.pool.QueryRow(ctx, `
		select id, path, created_at, updated_at
		from app.library_folders
		where id = $1
	`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return LibraryFolder{}, ErrNotFound
	}
	return folder, err
}
