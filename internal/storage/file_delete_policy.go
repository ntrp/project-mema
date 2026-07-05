package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type fileDeletePolicyResult struct {
	SourcePath      string
	DestinationPath string
	Mode            string
	Status          string
	Failure         string
}

func (s *SettingsStore) applyFileDeletePolicy(
	ctx context.Context,
	item MediaItem,
	source string,
) fileDeletePolicyResult {
	settings, err := s.GetFileDeleteSettings(ctx)
	if err != nil {
		return failedFileDeleteResult(source, FileDeleteModePermanent, err)
	}
	result := fileDeletePolicyResult{
		SourcePath: source,
		Mode:       settings.Mode,
		Status:     "succeeded",
	}
	switch settings.Mode {
	case FileDeleteModeKeep:
		result.Status = "skipped"
		return result
	case FileDeleteModeRecycle:
		destination, err := recycleDestination(item, source, settings)
		if err != nil {
			return failedFileDeleteResult(source, settings.Mode, err)
		}
		if err := moveFile(source, destination); err != nil {
			return failedFileDeleteResult(source, settings.Mode, err)
		}
		result.DestinationPath = destination
		return result
	default:
		if err := os.Remove(source); err != nil && !os.IsNotExist(err) {
			return failedFileDeleteResult(source, settings.Mode, err)
		}
		return result
	}
}

func (s *SettingsStore) applyFolderDeletePolicy(ctx context.Context, item MediaItem, source string) fileDeletePolicyResult {
	settings, err := s.GetFileDeleteSettings(ctx)
	if err != nil {
		return failedFileDeleteResult(source, FileDeleteModePermanent, err)
	}
	result := fileDeletePolicyResult{SourcePath: source, Mode: settings.Mode, Status: "succeeded"}
	switch settings.Mode {
	case FileDeleteModeKeep:
		result.Status = "skipped"
		return result
	case FileDeleteModeRecycle:
		destination, err := recycleDestination(item, source, settings)
		if err != nil {
			return failedFileDeleteResult(source, settings.Mode, err)
		}
		if err := moveFile(source, destination); err != nil {
			return failedFileDeleteResult(source, settings.Mode, err)
		}
		result.DestinationPath = destination
		return result
	default:
		if err := os.RemoveAll(source); err != nil {
			return failedFileDeleteResult(source, settings.Mode, err)
		}
		return result
	}
}

func recycleDestination(item MediaItem, source string, settings FileDeleteSettings) (string, error) {
	if item.LibraryFolderPath == nil {
		return "", ErrInvalidInput
	}
	root, err := safeAbsRoot(*item.LibraryFolderPath)
	if err != nil {
		return "", err
	}
	source, err = safePathUnderRoot(root, source, false)
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(root, source)
	if err != nil || relative == "." {
		return "", ErrInvalidInput
	}
	destination, err := safePathUnderRoot(root, filepath.Join(settings.RecycleFolder, relative), false)
	if err != nil {
		return "", err
	}
	if destination == source {
		return "", ErrInvalidInput
	}
	return uniqueRecycleDestination(destination), nil
}

func uniqueRecycleDestination(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}
	ext := filepath.Ext(path)
	base := path[:len(path)-len(ext)]
	suffix := time.Now().UTC().Format("20060102T150405")
	for index := 1; ; index++ {
		candidate := fmt.Sprintf("%s.%s.%d%s", base, suffix, index, ext)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}

func failedFileDeleteResult(source string, mode string, err error) fileDeletePolicyResult {
	return fileDeletePolicyResult{
		SourcePath: source,
		Mode:       mode,
		Status:     "failed",
		Failure:    err.Error(),
	}
}

func (s *SettingsStore) recordFileDeletePolicy(
	ctx context.Context,
	mediaItemID uuid.UUID,
	result fileDeletePolicyResult,
) error {
	details := map[string]any{"deleteMode": result.Mode}
	if result.DestinationPath != "" {
		details["recycledPath"] = result.DestinationPath
	}
	_, err := s.CreateMediaFileHistory(ctx, MediaFileHistoryInput{
		MediaItemID:     &mediaItemID,
		FilePath:        result.SourcePath,
		SourcePath:      optionalHistoryString(result.SourcePath),
		DestinationPath: optionalHistoryString(result.DestinationPath),
		Operation:       "deleted",
		Status:          result.Status,
		ActorType:       "user",
		Details:         details,
		FailureDetails:  optionalHistoryString(result.Failure),
	})
	return err
}
