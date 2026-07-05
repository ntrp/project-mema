package storage

import (
	"context"
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
	result := s.applyFileDeletePolicy(ctx, item, target)
	if err := s.recordFileDeletePolicy(ctx, item.ID, result); err != nil {
		return MediaItem{}, err
	}
	if result.Status == "failed" {
		return MediaItem{}, ErrInvalidInput
	}
	return s.RescanMediaItemFiles(ctx, id)
}

func mediaItemFileTarget(item MediaItem, filePath string) (string, error) {
	if item.MediaFolderPath == nil || strings.TrimSpace(*item.MediaFolderPath) == "" {
		return "", ErrInvalidInput
	}
	return safePathUnderRoot(*item.MediaFolderPath, filePath, false)
}
