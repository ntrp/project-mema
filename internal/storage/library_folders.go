package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	storagegen "media-manager/internal/storage/generated"
)

func (s *SettingsStore) GetLibraryFolder(ctx context.Context, id uuid.UUID) (LibraryFolder, error) {
	row, err := storagegen.New(s.pool).GetLibraryFolder(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return LibraryFolder{}, ErrNotFound
	}
	if err != nil {
		return LibraryFolder{}, err
	}
	return libraryFolderFromRow(row), nil
}
