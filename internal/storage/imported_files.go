package storage

import (
	"context"
	"errors"
	"path/filepath"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

func (s *SettingsStore) RecordImportedMediaFile(ctx context.Context, item MediaItem, filePath string) error {
	if item.LibraryFolderID == nil {
		return ErrNotFound
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	queries := storagegen.New(s.pool).WithTx(tx)
	mediaItemID := item.ID
	exists, err := queries.ImportedMediaFileExists(ctx, storagegen.ImportedMediaFileExistsParams{
		MediaItemID: &mediaItemID,
		Path:        filePath,
	})
	if err != nil {
		return err
	}
	if exists {
		return tx.Commit(ctx)
	}

	scanID := uuid.New()
	if err := queries.CreateImportedFileLibraryScan(ctx, storagegen.CreateImportedFileLibraryScanParams{
		ID:              scanID,
		LibraryFolderID: *item.LibraryFolderID,
	}); err != nil {
		return err
	}

	kind, err := mediaItemKind(item.Type)
	if err != nil {
		return err
	}
	if err := queries.CreateImportedFileLibraryScanItem(ctx, storagegen.CreateImportedFileLibraryScanItemParams{
		ID:                uuid.New(),
		ScanID:            scanID,
		Path:              filePath,
		FileName:          filepath.Base(filePath),
		DetectedTitle:     item.Title,
		DetectedYear:      int4Value(item.Year),
		DetectedMediaKind: kind,
		MediaItemID:       &mediaItemID,
	}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func mediaItemKind(mediaType string) (string, error) {
	switch mediaType {
	case "movie":
		return "movie", nil
	case "serie":
		return "series", nil
	default:
		return "", errors.New("unsupported media type")
	}
}
