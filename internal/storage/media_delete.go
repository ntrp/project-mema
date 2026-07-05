package storage

import (
	"context"
	"fmt"
	"strings"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

func (s *SettingsStore) DeleteMediaItem(ctx context.Context, id uuid.UUID, keepFiles bool) error {
	item, err := s.GetMediaItem(ctx, id)
	if err != nil {
		return err
	}
	if !keepFiles {
		if err := s.removeMediaFolder(ctx, item); err != nil {
			return err
		}
	}
	rows, err := storagegen.New(s.pool).DeleteMediaItemRecord(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) removeMediaFolder(ctx context.Context, item MediaItem) error {
	path, ok, err := mediaFolderDeletePath(item)
	if err != nil || !ok {
		return err
	}
	result := s.applyFolderDeletePolicy(ctx, item, path)
	if err := s.recordFileDeletePolicy(ctx, item.ID, result); err != nil {
		return err
	}
	if result.Status == "failed" {
		return fmt.Errorf("delete media folder: %s", result.Failure)
	}
	return nil
}

func mediaFolderDeletePath(item MediaItem) (string, bool, error) {
	if item.MediaFolderPath == nil || strings.TrimSpace(*item.MediaFolderPath) == "" {
		return "", false, nil
	}
	path, err := safeAbsRoot(*item.MediaFolderPath)
	if err != nil {
		return "", false, fmt.Errorf("unsafe media folder path: %s", path)
	}
	if item.LibraryFolderPath == nil || strings.TrimSpace(*item.LibraryFolderPath) == "" {
		return path, true, nil
	}
	if err := validatePathInRoot(*item.LibraryFolderPath, path, false); err != nil {
		return "", false, fmt.Errorf("media folder is outside library root: %s", path)
	}
	return path, true, nil
}
