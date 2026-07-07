package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

func (s *SettingsStore) ApplyMediaItemRename(ctx context.Context, id uuid.UUID) (MediaRenameApplyResult, error) {
	return s.ApplySelectedMediaItemRename(ctx, id, nil)
}

func (s *SettingsStore) ApplySelectedMediaItemRename(
	ctx context.Context,
	id uuid.UUID,
	currentPaths []string,
) (MediaRenameApplyResult, error) {
	item, err := s.GetMediaItem(ctx, id)
	if err != nil {
		return MediaRenameApplyResult{}, err
	}
	preview, err := s.PreviewMediaItemRename(ctx, id)
	if err != nil {
		return MediaRenameApplyResult{}, err
	}
	selected := renamePathSelection(currentPaths)
	result := MediaRenameApplyResult{Rows: make([]MediaRenamePreviewRow, 0, len(preview.Rows))}
	for _, row := range preview.Rows {
		applied := s.applySelectedMediaRenameRow(ctx, item, row, selected)
		result.Rows = append(result.Rows, applied)
		switch applied.Status {
		case "applied":
			result.AppliedCount++
		case "failed":
			result.FailedCount++
		default:
			result.SkippedCount++
		}
	}
	return result, nil
}

func (s *SettingsStore) applySelectedMediaRenameRow(
	ctx context.Context,
	item MediaItem,
	row MediaRenamePreviewRow,
	selected map[string]struct{},
) MediaRenamePreviewRow {
	if selected != nil {
		if _, ok := selected[row.CurrentPath]; !ok {
			row.Status = "skipped"
			row.Messages = append(row.Messages, "Skipped because it was not selected.")
			return row
		}
	}
	return s.applyMediaRenameRow(ctx, item, row)
}

func renamePathSelection(paths []string) map[string]struct{} {
	if paths == nil {
		return nil
	}
	selected := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		selected[path] = struct{}{}
	}
	return selected
}

func (s *SettingsStore) applyMediaRenameRow(ctx context.Context, item MediaItem, row MediaRenamePreviewRow) MediaRenamePreviewRow {
	if row.Status != "safe" {
		row.Status = "skipped"
		row.Messages = append(row.Messages, "Skipped because the latest preview is not safe to apply.")
		return row
	}
	if err := validateRenameApplyPaths(item, row); err != nil {
		return failedRenameRow(ctx, s, item.ID, row, err)
	}
	if _, err := os.Stat(row.CurrentPath); err != nil {
		return failedRenameRow(ctx, s, item.ID, row, fmt.Errorf("source file is not available: %w", err))
	}
	if _, err := os.Stat(row.ProposedPath); err == nil {
		return failedRenameRow(ctx, s, item.ID, row, fmt.Errorf("destination already exists"))
	} else if !os.IsNotExist(err) {
		return failedRenameRow(ctx, s, item.ID, row, fmt.Errorf("destination cannot be checked: %w", err))
	}
	if err := moveFile(row.CurrentPath, row.ProposedPath); err != nil {
		return failedRenameRow(ctx, s, item.ID, row, err)
	}
	if err := s.commitAppliedRename(ctx, item.ID, row.CurrentPath, row.ProposedPath); err != nil {
		rollbackErr := moveFile(row.ProposedPath, row.CurrentPath)
		if rollbackErr != nil {
			err = fmt.Errorf("%w; rollback failed: %v", err, rollbackErr)
		}
		return failedRenameRow(ctx, s, item.ID, row, err)
	}
	row.Status = "applied"
	row.Messages = append(row.Messages, "File renamed.")
	return row
}

func validateRenameApplyPaths(item MediaItem, row MediaRenamePreviewRow) error {
	if item.MediaFolderPath == nil {
		return ErrInvalidInput
	}
	if _, err := safePathUnderRoot(*item.MediaFolderPath, row.CurrentPath, false); err != nil {
		return err
	}
	if _, err := safePathUnderRoot(*item.MediaFolderPath, row.ProposedPath, false); err != nil {
		return err
	}
	return nil
}

func (s *SettingsStore) commitAppliedRename(ctx context.Context, mediaItemID uuid.UUID, source string, destination string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	queries := storagegen.New(s.pool).WithTx(tx)
	updated, err := queries.RenameMediaFileRecord(ctx, storagegen.RenameMediaFileRecordParams{
		MediaItemID:     &mediaItemID,
		SourcePath:      source,
		DestinationPath: destination,
		FileName:        filepath.Base(destination),
	})
	if err != nil {
		return err
	}
	if updated == 0 {
		updated, err = renameRelativeMediaFileRecord(ctx, queries, mediaItemID, source, destination)
		if err != nil {
			return err
		}
	}
	if updated == 0 {
		return ErrInvalidInput
	}
	if _, err := createMediaFileHistory(ctx, tx, MediaFileHistoryInput{
		MediaItemID:     &mediaItemID,
		FilePath:        destination,
		SourcePath:      &source,
		DestinationPath: &destination,
		Operation:       "renamed",
		Status:          "succeeded",
		ActorType:       "user",
	}); err != nil {
		return err
	}
	if err := queries.TouchMediaItem(ctx, mediaItemID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func renameRelativeMediaFileRecord(
	ctx context.Context,
	queries *storagegen.Queries,
	mediaItemID uuid.UUID,
	source string,
	destination string,
) (int64, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return 0, nil
	}
	relativeSource, err := filepath.Rel(cwd, source)
	if err != nil || relativeSource == "." || strings.HasPrefix(relativeSource, ".."+string(filepath.Separator)) {
		return 0, nil
	}
	relativeDestination, err := filepath.Rel(cwd, destination)
	if err != nil || relativeDestination == "." || strings.HasPrefix(relativeDestination, ".."+string(filepath.Separator)) {
		return 0, nil
	}
	return queries.RenameMediaFileRecord(ctx, storagegen.RenameMediaFileRecordParams{
		MediaItemID:     &mediaItemID,
		SourcePath:      relativeSource,
		DestinationPath: relativeDestination,
		FileName:        filepath.Base(destination),
	})
}

func failedRenameRow(ctx context.Context, s *SettingsStore, mediaItemID uuid.UUID, row MediaRenamePreviewRow, err error) MediaRenamePreviewRow {
	message := err.Error()
	row.Status = "failed"
	row.Messages = append(row.Messages, message)
	_, _ = s.CreateMediaFileHistory(ctx, MediaFileHistoryInput{
		MediaItemID:     &mediaItemID,
		FilePath:        row.CurrentPath,
		SourcePath:      &row.CurrentPath,
		DestinationPath: &row.ProposedPath,
		Operation:       "renamed",
		Status:          "failed",
		ActorType:       "user",
		FailureDetails:  &message,
	})
	return row
}
