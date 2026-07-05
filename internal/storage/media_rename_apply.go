package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

func (s *SettingsStore) ApplyMediaItemRename(ctx context.Context, id uuid.UUID) (MediaRenameApplyResult, error) {
	item, err := s.GetMediaItem(ctx, id)
	if err != nil {
		return MediaRenameApplyResult{}, err
	}
	preview, err := s.PreviewMediaItemRename(ctx, id)
	if err != nil {
		return MediaRenameApplyResult{}, err
	}
	result := MediaRenameApplyResult{Rows: make([]MediaRenamePreviewRow, 0, len(preview.Rows))}
	for _, row := range preview.Rows {
		applied := s.applyMediaRenameRow(ctx, item, row)
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
	if item.MediaFolderPath == nil || item.LibraryFolderPath == nil {
		return ErrInvalidInput
	}
	if _, err := safePathUnderRoot(*item.MediaFolderPath, row.CurrentPath, false); err != nil {
		return err
	}
	if _, err := safePathUnderRoot(*item.LibraryFolderPath, row.ProposedPath, false); err != nil {
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
