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

func (s *SettingsStore) MediaItemFilePathByID(ctx context.Context, mediaID, fileID uuid.UUID) (string, error) {
	item, err := s.GetMediaItem(ctx, mediaID)
	if err != nil {
		return "", err
	}
	for _, path := range item.FilePaths {
		if MediaItemFileID(item, path) == fileID {
			return mediaItemFileTarget(item, path)
		}
	}
	return "", ErrNotFound
}

func MediaItemFileID(item MediaItem, path string) uuid.UUID {
	for _, fact := range item.FileFacts {
		if fact.FilePath == path {
			return fact.ID
		}
	}
	return uuid.NewSHA1(item.ID, []byte(path))
}
