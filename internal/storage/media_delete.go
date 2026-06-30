package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func (s *SettingsStore) DeleteMediaItem(ctx context.Context, id uuid.UUID, keepFiles bool) error {
	item, err := s.GetMediaItem(ctx, id)
	if err != nil {
		return err
	}
	if !keepFiles {
		err = removeMediaFolder(item)
	}
	if err != nil {
		return err
	}
	tag, err := s.pool.Exec(ctx, `delete from app.media_items where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func removeMediaFolder(item MediaItem) error {
	path, ok, err := mediaFolderDeletePath(item)
	if err != nil || !ok {
		return err
	}
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("delete media folder: %w", err)
	}
	return nil
}

func mediaFolderDeletePath(item MediaItem) (string, bool, error) {
	if item.MediaFolderPath == nil || strings.TrimSpace(*item.MediaFolderPath) == "" {
		return "", false, nil
	}
	path := filepath.Clean(strings.TrimSpace(*item.MediaFolderPath))
	if !filepath.IsAbs(path) || path == string(os.PathSeparator) {
		return "", false, fmt.Errorf("unsafe media folder path: %s", path)
	}
	if item.LibraryFolderPath == nil || strings.TrimSpace(*item.LibraryFolderPath) == "" {
		return path, true, nil
	}
	root := filepath.Clean(strings.TrimSpace(*item.LibraryFolderPath))
	if path == root {
		return "", false, fmt.Errorf("refusing to delete library root: %s", path)
	}
	if !strings.HasPrefix(path, root+string(os.PathSeparator)) {
		return "", false, fmt.Errorf("media folder is outside library root: %s", path)
	}
	return path, true, nil
}
