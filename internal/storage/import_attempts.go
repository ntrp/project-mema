package storage

import (
	"context"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

func (s *SettingsStore) CreateImportAttempt(ctx context.Context, input ImportAttemptInput) (ImportAttempt, error) {
	mode := input.ImportMode
	if mode == "" {
		mode = "hardlink"
	}
	row, err := storagegen.New(s.pool).CreateImportAttempt(ctx, storagegen.CreateImportAttemptParams{
		ID:                     uuid.New(),
		ActivityID:             input.ActivityID,
		MediaItemID:            input.MediaItemID,
		SourcePath:             textValue(input.SourcePath),
		TargetPath:             textValue(input.TargetPath),
		ImportMode:             mode,
		Status:                 input.Status,
		FailureStage:           textValue(input.FailureStage),
		ErrorMessage:           textValue(input.ErrorMessage),
		CreatedTargets:         stringSliceValue(input.CreatedTargets),
		InsertedMediaFilePaths: stringSliceValue(input.InsertedMediaFilePaths),
	})
	return importAttemptFromRow(row), err
}

func (s *SettingsStore) ListImportAttemptsForActivity(ctx context.Context, activityID uuid.UUID) ([]ImportAttempt, error) {
	rows, err := storagegen.New(s.pool).ListImportAttemptsForActivity(ctx, activityID)
	if err != nil {
		return nil, err
	}
	attempts := make([]ImportAttempt, 0, len(rows))
	for _, row := range rows {
		attempts = append(attempts, importAttemptFromRow(row))
	}
	return attempts, nil
}

func importAttemptFromRow(row storagegen.AppImportAttempt) ImportAttempt {
	return ImportAttempt{
		ID:                     row.ID,
		ActivityID:             row.ActivityID,
		MediaItemID:            row.MediaItemID,
		SourcePath:             textPtr(row.SourcePath),
		TargetPath:             textPtr(row.TargetPath),
		ImportMode:             row.ImportMode,
		Status:                 row.Status,
		FailureStage:           textPtr(row.FailureStage),
		ErrorMessage:           textPtr(row.ErrorMessage),
		CreatedTargets:         row.CreatedTargets,
		InsertedMediaFilePaths: row.InsertedMediaFilePaths,
		CreatedAt:              row.CreatedAt,
		UpdatedAt:              row.UpdatedAt,
	}
}

func stringSliceValue(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}
