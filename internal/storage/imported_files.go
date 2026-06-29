package storage

import (
	"context"
	"errors"
	"path/filepath"

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

	var exists bool
	if err := tx.QueryRow(ctx, `
		select exists(
			select 1
			from app.library_scan_items
			where media_item_id = $1 and path = $2
		)
	`, item.ID, filePath).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return tx.Commit(ctx)
	}

	scanID := uuid.New()
	if _, err := tx.Exec(ctx, `
		insert into app.library_scans (
			id, library_folder_id, status, total_files, auto_matched_count, manual_count, completed_at
		)
		values ($1, $2, 'completed', 1, 1, 0, now())
	`, scanID, *item.LibraryFolderID); err != nil {
		return err
	}

	kind, err := mediaItemKind(item.Type)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		insert into app.library_scan_items (
			id, scan_id, path, file_name, detected_title, detected_year, detected_media_kind,
			status, matched_title, matched_year, matched_media_kind, media_item_id
		)
		values ($1, $2, $3, $4, $5, $6, $7, 'auto_added', $5, $6, $7, $8)
	`, uuid.New(), scanID, filePath, filepath.Base(filePath), item.Title, item.Year, kind, item.ID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func mediaItemKind(mediaType string) (string, error) {
	switch mediaType {
	case "movie":
		return "movie", nil
	case "series":
		return "series", nil
	default:
		return "", errors.New("unsupported media type")
	}
}
