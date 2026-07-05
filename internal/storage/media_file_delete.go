package storage

import (
	"context"
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
	result := s.applyFileDeletePolicy(ctx, item, target)
	if err := s.recordFileDeletePolicy(ctx, item.ID, result); err != nil {
		return MediaItem{}, err
	}
	if result.Status == "failed" {
		return MediaItem{}, ErrInvalidInput
	}
	return s.RescanMediaItemFiles(ctx, id)
}

func (s *SettingsStore) DeleteLibraryFolderFileForMedia(ctx context.Context, mediaItemID uuid.UUID, libraryFolderID uuid.UUID, relativePath string) error {
	item, err := s.GetMediaItem(ctx, mediaItemID)
	if err != nil {
		return err
	}
	folder, err := s.GetLibraryFolder(ctx, libraryFolderID)
	if err != nil {
		return err
	}
	source, err := safePathUnderRoot(folder.Path, filepath.Clean(relativePath), false)
	if err != nil {
		return err
	}
	result := s.applyFileDeletePolicy(ctx, item, source)
	if err := s.recordFileDeletePolicy(ctx, item.ID, result); err != nil {
		return err
	}
	if result.Status == "failed" {
		return ErrInvalidInput
	}
	return nil
}

func mediaItemFileTarget(item MediaItem, filePath string) (string, error) {
	if item.MediaFolderPath == nil || strings.TrimSpace(*item.MediaFolderPath) == "" {
		return "", ErrInvalidInput
	}
	return safePathUnderRoot(*item.MediaFolderPath, filePath, false)
}
