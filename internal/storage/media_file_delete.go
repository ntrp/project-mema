package storage

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func (s *SettingsStore) DeleteMediaItemFile(ctx context.Context, id uuid.UUID, filePath string) (MediaItem, error) {
	item, err := s.GetMediaItem(ctx, id)
	if err != nil {
		return MediaItem{}, err
	}
	target, err := mediaItemFileTarget(item, filePath)
	if err != nil {
		return MediaItem{}, err
	}
	if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
		return MediaItem{}, err
	}
	return s.RescanMediaItemFiles(ctx, id)
}

func mediaItemFileTarget(item MediaItem, filePath string) (string, error) {
	if item.MediaFolderPath == nil || strings.TrimSpace(*item.MediaFolderPath) == "" {
		return "", ErrInvalidInput
	}
	root := filepath.Clean(strings.TrimSpace(*item.MediaFolderPath))
	value := strings.TrimSpace(filePath)
	if value == "" {
		return "", ErrInvalidInput
	}
	target := filepath.Clean(value)
	if !filepath.IsAbs(target) {
		target = filepath.Join(root, target)
	}
	rel, err := filepath.Rel(root, target)
	if err != nil || rel == "." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return "", ErrInvalidInput
	}
	return target, nil
}
