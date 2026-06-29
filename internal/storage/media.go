package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MediaItem struct {
	ID               uuid.UUID
	Type             string
	Title            string
	Year             *int32
	Monitored        bool
	ExternalProvider *string
	ExternalID       *string
	Overview         *string
	PosterPath       *string
	QualityProfileID *string
	LibraryFolderID  *uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type MediaItemInput struct {
	Type             string
	Title            string
	Year             *int32
	Monitored        bool
	ExternalProvider *string
	ExternalID       *string
	Overview         *string
	PosterPath       *string
	QualityProfileID *string
	LibraryFolderID  *uuid.UUID
}

type DownloadActivity struct {
	ID                 uuid.UUID
	MediaItemID        uuid.UUID
	MediaTitle         string
	MediaType          string
	ReleaseTitle       string
	IndexerName        string
	DownloadClientName string
	DownloadURL        string
	Status             string
	Error              *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type DownloadActivityInput struct {
	MediaItemID        uuid.UUID
	ReleaseTitle       string
	IndexerName        string
	DownloadClientName string
	DownloadURL        string
	Status             string
	Error              *string
}

type ReleaseCandidate struct {
	ID          uuid.UUID
	MediaItemID uuid.UUID
	IndexerID   *uuid.UUID
	IndexerName string
	IndexerType string
	Title       string
	DownloadURL string
	InfoURL     *string
	GUID        *string
	SizeBytes   int64
	Seeders     *int32
	Peers       *int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ReleaseCandidateInput struct {
	MediaItemID uuid.UUID
	IndexerID   *uuid.UUID
	IndexerName string
	IndexerType string
	Title       string
	DownloadURL string
	InfoURL     *string
	GUID        *string
	SizeBytes   int64
	Seeders     *int32
	Peers       *int32
}

type ReleaseSearchSnapshot struct {
	Releases []ReleaseCandidate
	Errors   []string
}

func (s *SettingsStore) ListMediaItems(ctx context.Context) ([]MediaItem, error) {
	rows, err := s.pool.Query(ctx, `
		select id, media_type, title, year, monitored, external_provider, external_id, overview, poster_path, quality_profile_id, library_folder_id, created_at, updated_at
		from app.media_items
		order by created_at desc, title asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []MediaItem{}
	for rows.Next() {
		item, err := scanMediaItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *SettingsStore) SearchMediaItems(ctx context.Context, query string, mediaType *string, limit int) ([]MediaItem, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	rows, err := s.pool.Query(ctx, `
		select id, media_type, title, year, monitored, external_provider, external_id, overview, poster_path, quality_profile_id, library_folder_id, created_at, updated_at
		from app.media_items
		where title ilike '%' || $1 || '%'
			and ($2::text is null or media_type = $2)
		order by
			case when lower(title) = lower($1) then 0 else 1 end,
			title asc
		limit $3
	`, query, mediaType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []MediaItem{}
	for rows.Next() {
		item, err := scanMediaItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *SettingsStore) GetMediaItem(ctx context.Context, id uuid.UUID) (MediaItem, error) {
	return scanMediaItemRow(s.pool.QueryRow(ctx, `
		select id, media_type, title, year, monitored, external_provider, external_id, overview, poster_path, quality_profile_id, library_folder_id, created_at, updated_at
		from app.media_items
		where id = $1
	`, id))
}

func (s *SettingsStore) CreateMediaItem(ctx context.Context, input MediaItemInput) (MediaItem, error) {
	id := uuid.New()
	return scanMediaItemRow(s.pool.QueryRow(ctx, `
		insert into app.media_items (
			id, media_type, title, year, monitored, external_provider, external_id, overview, poster_path, quality_profile_id, library_folder_id
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		returning id, media_type, title, year, monitored, external_provider, external_id, overview, poster_path, quality_profile_id, library_folder_id, created_at, updated_at
	`, id, input.Type, input.Title, input.Year, input.Monitored, input.ExternalProvider, input.ExternalID, input.Overview, input.PosterPath, input.QualityProfileID, input.LibraryFolderID))
}

func (s *SettingsStore) DeleteMediaItem(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `delete from app.media_items where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) ListEnabledIndexers(ctx context.Context) ([]Indexer, error) {
	rows, err := s.pool.Query(ctx, `
		select id, name, type, base_url, api_key, categories, enabled, priority, created_at, updated_at
		from app.indexers
		where enabled = true
		order by priority asc, name asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indexers := []Indexer{}
	for rows.Next() {
		indexer, err := scanIndexer(rows)
		if err != nil {
			return nil, err
		}
		indexers = append(indexers, indexer)
	}
	return indexers, rows.Err()
}

func (s *SettingsStore) ListEnabledDownloadClients(ctx context.Context) ([]DownloadClient, error) {
	rows, err := s.pool.Query(ctx, `
		select id, name, type, base_url, username, password, api_key, category, enabled, priority, created_at, updated_at
		from app.download_clients
		where enabled = true
		order by priority asc, name asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := []DownloadClient{}
	for rows.Next() {
		client, err := scanDownloadClient(rows)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, rows.Err()
}

func (s *SettingsStore) CreateDownloadActivity(ctx context.Context, input DownloadActivityInput) (DownloadActivity, error) {
	id := uuid.New()
	return scanDownloadActivityRow(s.pool.QueryRow(ctx, `
		insert into app.download_activity (
			id, media_item_id, release_title, indexer_name, download_client_name, download_url, status, error
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		returning id, media_item_id, release_title, indexer_name, download_client_name, download_url, status, error, created_at, updated_at
	`, id, input.MediaItemID, input.ReleaseTitle, input.IndexerName, input.DownloadClientName, input.DownloadURL, input.Status, input.Error))
}

func (s *SettingsStore) UpdateDownloadActivityStatus(ctx context.Context, id uuid.UUID, status string, activityError *string) (DownloadActivity, error) {
	return scanDownloadActivityRow(s.pool.QueryRow(ctx, `
		update app.download_activity
		set status = $2, error = $3, updated_at = now()
		where id = $1
		returning id, media_item_id, release_title, indexer_name, download_client_name, download_url, status, error, created_at, updated_at
	`, id, status, activityError))
}

func (s *SettingsStore) ReplaceReleaseSearchResults(ctx context.Context, mediaItemID uuid.UUID, releases []ReleaseCandidateInput, searchErrors []string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, `delete from app.media_release_candidates where media_item_id = $1`, mediaItemID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `delete from app.media_release_search_errors where media_item_id = $1`, mediaItemID); err != nil {
		return err
	}
	for _, release := range releases {
		if _, err := tx.Exec(ctx, `
			insert into app.media_release_candidates (
				id, media_item_id, indexer_id, indexer_name, indexer_type, title, download_url,
				info_url, guid, size_bytes, seeders, peers
			)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, uuid.New(), mediaItemID, release.IndexerID, release.IndexerName, release.IndexerType, release.Title,
			release.DownloadURL, release.InfoURL, release.GUID, release.SizeBytes, release.Seeders, release.Peers); err != nil {
			return err
		}
	}
	for _, message := range searchErrors {
		if _, err := tx.Exec(ctx, `
			insert into app.media_release_search_errors (id, media_item_id, message)
			values ($1, $2, $3)
		`, uuid.New(), mediaItemID, message); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (s *SettingsStore) GetReleaseCandidate(ctx context.Context, id uuid.UUID, mediaItemID uuid.UUID) (ReleaseCandidate, error) {
	release, err := scanReleaseCandidate(s.pool.QueryRow(ctx, `
		select id, media_item_id, indexer_id, indexer_name, indexer_type, title, download_url,
			info_url, guid, size_bytes, seeders, peers, created_at, updated_at
		from app.media_release_candidates
		where id = $1 and media_item_id = $2
	`, id, mediaItemID))
	if errors.Is(err, pgx.ErrNoRows) {
		return ReleaseCandidate{}, ErrNotFound
	}
	return release, err
}

func (s *SettingsStore) ListReleaseSearchResults(ctx context.Context, mediaItemID uuid.UUID) (ReleaseSearchSnapshot, error) {
	releaseRows, err := s.pool.Query(ctx, `
		select id, media_item_id, indexer_id, indexer_name, indexer_type, title, download_url,
			info_url, guid, size_bytes, seeders, peers, created_at, updated_at
		from app.media_release_candidates
		where media_item_id = $1
		order by coalesce(seeders, -1) desc, size_bytes desc, created_at desc
	`, mediaItemID)
	if err != nil {
		return ReleaseSearchSnapshot{}, err
	}
	defer releaseRows.Close()

	snapshot := ReleaseSearchSnapshot{Releases: []ReleaseCandidate{}, Errors: []string{}}
	for releaseRows.Next() {
		release, err := scanReleaseCandidate(releaseRows)
		if err != nil {
			return ReleaseSearchSnapshot{}, err
		}
		snapshot.Releases = append(snapshot.Releases, release)
	}
	if err := releaseRows.Err(); err != nil {
		return ReleaseSearchSnapshot{}, err
	}

	errorRows, err := s.pool.Query(ctx, `
		select message
		from app.media_release_search_errors
		where media_item_id = $1
		order by created_at asc
	`, mediaItemID)
	if err != nil {
		return ReleaseSearchSnapshot{}, err
	}
	defer errorRows.Close()

	for errorRows.Next() {
		var message string
		if err := errorRows.Scan(&message); err != nil {
			return ReleaseSearchSnapshot{}, err
		}
		snapshot.Errors = append(snapshot.Errors, message)
	}
	return snapshot, errorRows.Err()
}

func (s *SettingsStore) ListDownloadActivity(ctx context.Context) ([]DownloadActivity, error) {
	rows, err := s.pool.Query(ctx, `
		select
			a.id,
			a.media_item_id,
			m.title,
			m.media_type,
			a.release_title,
			a.indexer_name,
			a.download_client_name,
			a.download_url,
			a.status,
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

	activities := []DownloadActivity{}
	for rows.Next() {
		var activity DownloadActivity
		if err := rows.Scan(
			&activity.ID,
			&activity.MediaItemID,
			&activity.MediaTitle,
			&activity.MediaType,
			&activity.ReleaseTitle,
			&activity.IndexerName,
			&activity.DownloadClientName,
			&activity.DownloadURL,
			&activity.Status,
			&activity.Error,
			&activity.CreatedAt,
			&activity.UpdatedAt,
		); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, rows.Err()
}

func scanMediaItemRow(row pgx.Row) (MediaItem, error) {
	item, err := scanMediaItem(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaItem{}, ErrNotFound
	}
	return item, err
}

func scanMediaItem(row pgx.Row) (MediaItem, error) {
	var item MediaItem
	err := row.Scan(
		&item.ID,
		&item.Type,
		&item.Title,
		&item.Year,
		&item.Monitored,
		&item.ExternalProvider,
		&item.ExternalID,
		&item.Overview,
		&item.PosterPath,
		&item.QualityProfileID,
		&item.LibraryFolderID,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	return item, err
}

func scanReleaseCandidate(row pgx.Row) (ReleaseCandidate, error) {
	var release ReleaseCandidate
	err := row.Scan(
		&release.ID,
		&release.MediaItemID,
		&release.IndexerID,
		&release.IndexerName,
		&release.IndexerType,
		&release.Title,
		&release.DownloadURL,
		&release.InfoURL,
		&release.GUID,
		&release.SizeBytes,
		&release.Seeders,
		&release.Peers,
		&release.CreatedAt,
		&release.UpdatedAt,
	)
	return release, err
}

func scanDownloadActivityRow(row pgx.Row) (DownloadActivity, error) {
	var activity DownloadActivity
	err := row.Scan(
		&activity.ID,
		&activity.MediaItemID,
		&activity.ReleaseTitle,
		&activity.IndexerName,
		&activity.DownloadClientName,
		&activity.DownloadURL,
		&activity.Status,
		&activity.Error,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return DownloadActivity{}, ErrNotFound
	}
	return activity, err
}
