package storage

import (
	"context"

	"github.com/google/uuid"
)

func (s *SettingsStore) MediaItemFilePath(ctx context.Context, id uuid.UUID, filePath string) (string, error) {
	item, err := s.GetMediaItem(ctx, id)
	if err != nil {
		return "", err
	}
	return mediaItemFileTarget(item, filePath)
}
