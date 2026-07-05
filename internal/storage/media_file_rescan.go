package storage

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	storagegen "media-manager/internal/storage/generated"

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

	queries := storagegen.New(s.pool).WithTx(tx)
	mediaItemID := item.ID
	if err := queries.DeleteLibraryScanItemsForMediaItem(ctx, &mediaItemID); err != nil {
		return MediaItem{}, err
	}
	scanID := uuid.New()
	if err := queries.CreateMediaFileRescanLibraryScan(ctx, storagegen.CreateMediaFileRescanLibraryScanParams{
		ID:              scanID,
		LibraryFolderID: *item.LibraryFolderID,
		TotalFiles:      int32(len(files)),
	}); err != nil {
		return MediaItem{}, err
	}
	for _, path := range files {
		if err := queries.CreateImportedFileLibraryScanItem(ctx, storagegen.CreateImportedFileLibraryScanItemParams{
			ID:                uuid.New(),
			ScanID:            scanID,
			Path:              path,
			FileName:          filepath.Base(path),
			DetectedTitle:     item.Title,
			DetectedYear:      int4Value(item.Year),
			DetectedMediaKind: kind,
			MediaItemID:       &mediaItemID,
		}); err != nil {
			return MediaItem{}, err
		}
	}
	if err := queries.TouchMediaItem(ctx, item.ID); err != nil {
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
