package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *SettingsStore) CreateDownloadActivity(ctx context.Context, input DownloadActivityInput) (DownloadActivity, error) {
	id := uuid.New()
	return scanDownloadActivityRow(s.pool.QueryRow(ctx, `
		insert into app.download_activity (
			id, media_item_id, release_title, indexer_name, download_client_name, download_id, download_url, status, error
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		returning id, media_item_id, release_title, indexer_name, download_client_name, download_id, download_url, status, progress_percent, error, created_at, updated_at
	`, id, input.MediaItemID, input.ReleaseTitle, input.IndexerName, input.DownloadClientName, input.DownloadID, input.DownloadURL, input.Status, input.Error))
}

func (s *SettingsStore) UpdateDownloadActivityStatus(ctx context.Context, id uuid.UUID, status string, activityError *string) (DownloadActivity, error) {
	return s.UpdateDownloadActivityProgress(ctx, id, status, nil, activityError)
}

func (s *SettingsStore) UpdateDownloadActivityProgress(ctx context.Context, id uuid.UUID, status string, progressPercent *int, activityError *string) (DownloadActivity, error) {
	return s.updateDownloadActivity(ctx, id, status, nil, progressPercent, activityError)
}

func (s *SettingsStore) UpdateDownloadActivityClientState(ctx context.Context, id uuid.UUID, status string, downloadID *string, activityError *string) (DownloadActivity, error) {
	return s.updateDownloadActivity(ctx, id, status, downloadID, nil, activityError)
}

func (s *SettingsStore) updateDownloadActivity(ctx context.Context, id uuid.UUID, status string, downloadID *string, progressPercent *int, activityError *string) (DownloadActivity, error) {
	return scanDownloadActivityRow(s.pool.QueryRow(ctx, `
		update app.download_activity
		set status = $2,
			download_id = coalesce($3, download_id),
			progress_percent = $4,
			error = $5,
			updated_at = now()
		where id = $1
		returning id, media_item_id, release_title, indexer_name, download_client_name, download_id, download_url, status, progress_percent, error, created_at, updated_at
	`, id, status, downloadID, progressPercent, activityError))
}

func (s *SettingsStore) ListDownloadActivity(ctx context.Context) ([]DownloadActivity, error) {
	rows, err := s.pool.Query(ctx, `
		select
			a.id,
			a.media_item_id,
			m.title,
			m.media_type,
			m.year,
			a.release_title,
			a.indexer_name,
			a.download_client_name,
			a.download_id,
			a.download_url,
			a.status,
			a.progress_percent,
			a.error,
			a.created_at,
			a.updated_at
		from app.download_activity a
		join app.media_items m on m.id = a.media_item_id
		order by a.created_at desc
		limit 100
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanDownloadActivities(rows)
}

func (s *SettingsStore) GetDownloadActivity(ctx context.Context, id uuid.UUID) (DownloadActivity, error) {
	activity, err := scanDownloadActivityWithMediaRow(s.pool.QueryRow(ctx, `
		select
			a.id,
			a.media_item_id,
			m.title,
			m.media_type,
			m.year,
			a.release_title,
			a.indexer_name,
			a.download_client_name,
			a.download_id,
			a.download_url,
			a.status,
			a.progress_percent,
			a.error,
			a.created_at,
			a.updated_at
		from app.download_activity a
		join app.media_items m on m.id = a.media_item_id
		where a.id = $1
	`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return DownloadActivity{}, ErrNotFound
	}
	return activity, err
}

func (s *SettingsStore) CancelDownloadActivity(ctx context.Context, id uuid.UUID) (DownloadActivity, error) {
	return scanDownloadActivityRow(s.pool.QueryRow(ctx, `
		update app.download_activity
		set status = 'cancelled', progress_percent = null, error = null, updated_at = now()
		where id = $1
			and status in ('queued', 'grabbed', 'downloading')
		returning id, media_item_id, release_title, indexer_name, download_client_name, download_id, download_url, status, progress_percent, error, created_at, updated_at
	`, id))
}

func (s *SettingsStore) ListActiveDownloadActivity(ctx context.Context) ([]DownloadActivity, error) {
	rows, err := s.pool.Query(ctx, `
		select
			a.id,
			a.media_item_id,
			m.title,
			m.media_type,
			m.year,
			a.release_title,
			a.indexer_name,
			a.download_client_name,
			a.download_id,
			a.download_url,
			a.status,
			a.progress_percent,
			a.error,
			a.created_at,
			a.updated_at
		from app.download_activity a
		join app.media_items m on m.id = a.media_item_id
		where a.status in ('queued', 'grabbed', 'downloading')
			and a.download_id is not null
		order by a.updated_at asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanDownloadActivities(rows)
}

func scanDownloadActivities(rows pgx.Rows) ([]DownloadActivity, error) {
	activities := []DownloadActivity{}
	for rows.Next() {
		activity, err := scanDownloadActivityWithMediaRow(rows)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, rows.Err()
}

func scanDownloadActivityWithMediaRow(row pgx.Row) (DownloadActivity, error) {
	var activity DownloadActivity
	err := row.Scan(
		&activity.ID,
		&activity.MediaItemID,
		&activity.MediaTitle,
		&activity.MediaType,
		&activity.MediaYear,
		&activity.ReleaseTitle,
		&activity.IndexerName,
		&activity.DownloadClientName,
		&activity.DownloadID,
		&activity.DownloadURL,
		&activity.Status,
		&activity.ProgressPercent,
		&activity.Error,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)
	return activity, err
}

func scanDownloadActivityRow(row pgx.Row) (DownloadActivity, error) {
	var activity DownloadActivity
	err := row.Scan(
		&activity.ID,
		&activity.MediaItemID,
		&activity.ReleaseTitle,
		&activity.IndexerName,
		&activity.DownloadClientName,
		&activity.DownloadID,
		&activity.DownloadURL,
		&activity.Status,
		&activity.ProgressPercent,
		&activity.Error,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return DownloadActivity{}, ErrNotFound
	}
	return activity, err
}
