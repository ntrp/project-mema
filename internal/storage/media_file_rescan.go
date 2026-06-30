package storage

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/uuid"
)

var mediaFileExtensions = map[string]struct{}{
	".avi":  {},
	".m4v":  {},
	".mkv":  {},
	".mov":  {},
	".mp4":  {},
	".mpeg": {},
	".mpg":  {},
	".ts":   {},
	".webm": {},
	".wmv":  {},
}

func (s *SettingsStore) RescanMediaItemFiles(ctx context.Context, id uuid.UUID) (MediaItem, error) {
	item, err := s.GetMediaItem(ctx, id)
	if err != nil {
		return MediaItem{}, err
	}
	if item.LibraryFolderID == nil || item.MediaFolderPath == nil || strings.TrimSpace(*item.MediaFolderPath) == "" {
		return MediaItem{}, ErrInvalidInput
	}
	root := filepath.Clean(strings.TrimSpace(*item.MediaFolderPath))
	files, err := mediaFilesInRoot(root)
	if err != nil {
		return MediaItem{}, ErrInvalidInput
	}
	kind, err := mediaItemKind(item.Type)
	if err != nil {
		return MediaItem{}, ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return MediaItem{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, `delete from app.library_scan_items where media_item_id = $1`, item.ID); err != nil {
		return MediaItem{}, err
	}
	scanID := uuid.New()
	if _, err := tx.Exec(ctx, `
		insert into app.library_scans (
			id, library_folder_id, status, total_files, auto_matched_count, manual_count, completed_at
		)
		values ($1, $2, 'completed', $3, $3, 0, now())
	`, scanID, *item.LibraryFolderID, int32(len(files))); err != nil {
		return MediaItem{}, err
	}
	for _, path := range files {
		if _, err := tx.Exec(ctx, `
			insert into app.library_scan_items (
				id, scan_id, path, file_name, detected_title, detected_year, detected_media_kind,
				status, matched_title, matched_year, matched_media_kind, media_item_id
			)
			values ($1, $2, $3, $4, $5, $6, $7, 'auto_added', $5, $6, $7, $8)
		`, uuid.New(), scanID, path, filepath.Base(path), item.Title, item.Year, kind, item.ID); err != nil {
			return MediaItem{}, err
		}
	}
	if _, err := tx.Exec(ctx, `update app.media_items set updated_at = now() where id = $1`, item.ID); err != nil {
		return MediaItem{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return MediaItem{}, err
	}
	return s.GetMediaItem(ctx, id)
}

func mediaFilesInRoot(root string) ([]string, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, ErrInvalidInput
	}
	files := []string{}
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if strings.HasPrefix(entry.Name(), ".") && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		if _, ok := mediaFileExtensions[strings.ToLower(filepath.Ext(entry.Name()))]; !ok {
			return nil
		}
		files = append(files, path)
		return nil
	})
	sort.Strings(files)
	return files, err
}
