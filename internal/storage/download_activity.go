package storage

import (
	"context"
	"errors"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *SettingsStore) CreateDownloadActivity(ctx context.Context, input DownloadActivityInput) (DownloadActivity, error) {
	id := uuid.New()
	row, err := storagegen.New(s.pool).CreateDownloadActivity(ctx, storagegen.CreateDownloadActivityParams{
		ID:                 id,
		MediaItemID:        input.MediaItemID,
		ReleaseTitle:       input.ReleaseTitle,
		IndexerName:        input.IndexerName,
		DownloadClientName: input.DownloadClientName,
		DownloadID:         textValue(input.DownloadID),
		DownloadUrl:        input.DownloadURL,
		Status:             input.Status,
		Error:              textValue(input.Error),
		FailureType:        textValue(input.FailureType),
	})
	return downloadActivityFromRow(row), err
}

func (s *SettingsStore) UpdateDownloadActivityStatus(ctx context.Context, id uuid.UUID, status string, activityError *string) (DownloadActivity, error) {
	return s.UpdateDownloadActivityProgress(ctx, id, status, nil, activityError)
}

func (s *SettingsStore) FailDownloadActivity(ctx context.Context, id uuid.UUID, activityError *string, failureType string) (DownloadActivity, error) {
	row, err := storagegen.New(s.pool).FailDownloadActivity(ctx, storagegen.FailDownloadActivityParams{
		Error:       textValue(activityError),
		FailureType: textValue(&failureType),
		ID:          id,
	})
	return downloadActivityResult(row, err)
}

func (s *SettingsStore) UpdateDownloadActivityProgress(ctx context.Context, id uuid.UUID, status string, progressPercent *int, activityError *string) (DownloadActivity, error) {
	return s.updateDownloadActivity(ctx, id, status, nil, progressPercent, activityError)
}

func (s *SettingsStore) UpdateDownloadActivityClientState(ctx context.Context, id uuid.UUID, status string, downloadID *string, activityError *string) (DownloadActivity, error) {
	return s.updateDownloadActivity(ctx, id, status, downloadID, nil, activityError)
}

func (s *SettingsStore) updateDownloadActivity(ctx context.Context, id uuid.UUID, status string, downloadID *string, progressPercent *int, activityError *string) (DownloadActivity, error) {
	row, err := storagegen.New(s.pool).UpdateDownloadActivity(ctx, storagegen.UpdateDownloadActivityParams{
		Status:          status,
		DownloadID:      textValue(downloadID),
		ProgressPercent: intValue(progressPercent),
		Error:           textValue(activityError),
		ID:              id,
	})
	return downloadActivityResult(row, err)
}

func (s *SettingsStore) ListDownloadActivity(ctx context.Context) ([]DownloadActivity, error) {
	rows, err := storagegen.New(s.pool).ListDownloadActivity(ctx)
	if err != nil {
		return nil, err
	}
	activities := make([]DownloadActivity, 0, len(rows))
	for _, row := range rows {
		activities = append(activities, downloadActivityFromListRow(row))
	}
	return activities, nil
}

func (s *SettingsStore) GetDownloadActivity(ctx context.Context, id uuid.UUID) (DownloadActivity, error) {
	row, err := storagegen.New(s.pool).GetDownloadActivity(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return DownloadActivity{}, ErrNotFound
	}
	return downloadActivityFromGetRow(row), err
}

func (s *SettingsStore) CancelDownloadActivity(ctx context.Context, id uuid.UUID) (DownloadActivity, error) {
	row, err := storagegen.New(s.pool).CancelDownloadActivity(ctx, id)
	return downloadActivityResult(row, err)
}

func (s *SettingsStore) DeleteDownloadActivity(ctx context.Context, id uuid.UUID) error {
	rows, err := storagegen.New(s.pool).DeleteDownloadActivity(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) ListActiveDownloadActivity(ctx context.Context) ([]DownloadActivity, error) {
	rows, err := storagegen.New(s.pool).ListActiveDownloadActivity(ctx)
	if err != nil {
		return nil, err
	}
	activities := make([]DownloadActivity, 0, len(rows))
	for _, row := range rows {
		activities = append(activities, downloadActivityFromActiveRow(row))
	}
	return activities, nil
}

func downloadActivityResult(row storagegen.AppDownloadActivity, err error) (DownloadActivity, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return DownloadActivity{}, ErrNotFound
	}
	return downloadActivityFromRow(row), err
}

func downloadActivityFromRow(row storagegen.AppDownloadActivity) DownloadActivity {
	return DownloadActivity{
		ID:                 row.ID,
		MediaItemID:        row.MediaItemID,
		ReleaseTitle:       row.ReleaseTitle,
		IndexerName:        row.IndexerName,
		DownloadClientName: row.DownloadClientName,
		DownloadID:         textPtr(row.DownloadID),
		DownloadURL:        row.DownloadUrl,
		Status:             row.Status,
		ProgressPercent:    intPtr(row.ProgressPercent),
		Error:              textPtr(row.Error),
		FailureType:        textPtr(row.FailureType),
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
}

func downloadActivityFromListRow(row storagegen.ListDownloadActivityRow) DownloadActivity {
	return DownloadActivity{
		ID:                 row.ID,
		MediaItemID:        row.MediaItemID,
		MediaTitle:         row.MediaTitle,
		MediaType:          row.MediaType,
		MediaYear:          int4Ptr(row.MediaYear),
		ReleaseTitle:       row.ReleaseTitle,
		IndexerName:        row.IndexerName,
		DownloadClientName: row.DownloadClientName,
		DownloadID:         textPtr(row.DownloadID),
		DownloadURL:        row.DownloadUrl,
		Status:             row.Status,
		ProgressPercent:    intPtr(row.ProgressPercent),
		Error:              textPtr(row.Error),
		FailureType:        textPtr(row.FailureType),
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
}

func downloadActivityFromGetRow(row storagegen.GetDownloadActivityRow) DownloadActivity {
	return DownloadActivity{
		ID:                 row.ID,
		MediaItemID:        row.MediaItemID,
		MediaTitle:         row.MediaTitle,
		MediaType:          row.MediaType,
		MediaYear:          int4Ptr(row.MediaYear),
		ReleaseTitle:       row.ReleaseTitle,
		IndexerName:        row.IndexerName,
		DownloadClientName: row.DownloadClientName,
		DownloadID:         textPtr(row.DownloadID),
		DownloadURL:        row.DownloadUrl,
		Status:             row.Status,
		ProgressPercent:    intPtr(row.ProgressPercent),
		Error:              textPtr(row.Error),
		FailureType:        textPtr(row.FailureType),
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
}

func downloadActivityFromActiveRow(row storagegen.ListActiveDownloadActivityRow) DownloadActivity {
	return DownloadActivity{
		ID:                 row.ID,
		MediaItemID:        row.MediaItemID,
		MediaTitle:         row.MediaTitle,
		MediaType:          row.MediaType,
		MediaYear:          int4Ptr(row.MediaYear),
		ReleaseTitle:       row.ReleaseTitle,
		IndexerName:        row.IndexerName,
		DownloadClientName: row.DownloadClientName,
		DownloadID:         textPtr(row.DownloadID),
		DownloadURL:        row.DownloadUrl,
		Status:             row.Status,
		ProgressPercent:    intPtr(row.ProgressPercent),
		Error:              textPtr(row.Error),
		FailureType:        textPtr(row.FailureType),
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
}
