package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	storagegen "media-manager/internal/storage/generated"
)

func (s *SettingsStore) CreateLibraryFolder(ctx context.Context, path string, kind string) (LibraryFolder, error) {
	return s.CreateLibraryFolderWithManual(ctx, path, kind, false)
}

func (s *SettingsStore) CreateLibraryFolderWithManual(ctx context.Context, path string, kind string, manual bool) (LibraryFolder, error) {
	path, err := absoluteCleanPath(path)
	if err != nil {
		return LibraryFolder{}, err
	}
	row, err := storagegen.New(s.pool).UpsertLibraryFolder(ctx, storagegen.UpsertLibraryFolderParams{
		ID: uuid.New(), Path: path, Kind: kind, Manual: manual,
	})
	if err != nil {
		return LibraryFolder{}, err
	}
	return libraryFolderFromRow(row), nil
}

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

func (s *SettingsStore) UpdateLibraryFolderManual(ctx context.Context, id uuid.UUID, manual bool) (LibraryFolder, error) {
	row, err := storagegen.New(s.pool).UpdateLibraryFolderManual(ctx, storagegen.UpdateLibraryFolderManualParams{
		ID: id, Manual: manual,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return LibraryFolder{}, ErrNotFound
	}
	if err != nil {
		return LibraryFolder{}, err
	}
	return libraryFolderFromRow(row), nil
}

func (s *SettingsStore) MediaItemAllowsAutomation(ctx context.Context, item MediaItem, path string) (bool, error) {
	folders, err := s.ListLibraryFolders(ctx)
	if err != nil {
		return false, err
	}
	return mediaItemAllowsAutomation(folders, item, path), nil
}

func (s *SettingsStore) AutomatedMediaFileFacts(ctx context.Context, item MediaItem) ([]MediaFileFact, error) {
	folders, err := s.ListLibraryFolders(ctx)
	if err != nil {
		return nil, err
	}
	facts := make([]MediaFileFact, 0, len(item.FileFacts))
	for _, fact := range item.FileFacts {
		if mediaItemAllowsAutomation(folders, item, fact.FilePath) {
			facts = append(facts, fact)
		}
	}
	return facts, nil
}

func mediaItemAllowsAutomation(folders []LibraryFolder, item MediaItem, path string) bool {
	if path == "" && item.MediaFolderPath != nil {
		path = *item.MediaFolderPath
	}
	for _, folder := range folders {
		if folder.Manual && path != "" && pathWithinOrEqual(folder.Path, path) {
			return false
		}
	}
	if path != "" || item.LibraryFolderID == nil {
		return true
	}
	for _, folder := range folders {
		if folder.Manual && *item.LibraryFolderID == folder.ID {
			return false
		}
	}
	return true
}
