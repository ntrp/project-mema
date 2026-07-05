package storage

import (
	"context"
	"os"
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
		_ = s.recordMediaFileDelete(ctx, item, target, "failed", err.Error())
		return MediaItem{}, err
	}
	if err := s.recordMediaFileDelete(ctx, item, target, "succeeded", ""); err != nil {
		return MediaItem{}, err
	}
	return s.RescanMediaItemFiles(ctx, id)
}

func mediaItemFileTarget(item MediaItem, filePath string) (string, error) {
	if item.MediaFolderPath == nil || strings.TrimSpace(*item.MediaFolderPath) == "" {
		return "", ErrInvalidInput
	}
	return safePathUnderRoot(*item.MediaFolderPath, filePath, false)
}

func (s *SettingsStore) recordMediaFileDelete(
	ctx context.Context,
	item MediaItem,
	path string,
	status string,
	failure string,
) error {
	mediaItemID := item.ID
	_, err := s.CreateMediaFileHistory(ctx, MediaFileHistoryInput{
		MediaItemID:    &mediaItemID,
		FilePath:       path,
		SourcePath:     optionalHistoryString(path),
		Operation:      "deleted",
		Status:         status,
		ActorType:      "user",
		FailureDetails: optionalHistoryString(failure),
	})
	return err
}
