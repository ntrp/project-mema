package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MediaRequest struct {
	ID                  uuid.UUID
	RequestedByUserID   uuid.UUID
	RequestedByUsername string
	Type                string
	Title               string
	Year                *int32
	ExternalProvider    *string
	ExternalID          *string
	Overview            *string
	PosterPath          *string
	Tags                []string
	Status              string
	QualityProfileID    *string
	LibraryFolderID     *uuid.UUID
	MediaItemID         *uuid.UUID
	DecidedAt           *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type MediaRequestInput struct {
	RequestedByUserID uuid.UUID
	Type              string
	Title             string
	Year              *int32
	ExternalProvider  *string
	ExternalID        *string
	Overview          *string
	PosterPath        *string
	Tags              []string
}

type MediaRequestApprovalInput struct {
	QualityProfileID string
	LibraryFolderID  uuid.UUID
}

func (s *SettingsStore) ListMediaRequests(ctx context.Context, userID uuid.UUID, includeAll bool) ([]MediaRequest, error) {
	rows, err := s.pool.Query(ctx, `
		select r.id, r.requested_by_user_id, u.username, r.media_type, r.title, r.year,
			r.external_provider, r.external_id, r.overview, r.poster_path, r.status,
			r.quality_profile_id, r.library_folder_id, r.media_item_id, r.decided_at,
			coalesce(array(
				select t.name
				from app.media_request_tags mrt
				join app.tags t on t.id = mrt.tag_id
				where mrt.media_request_id = r.id
				order by lower(t.name)
			), '{}') as tags,
			r.created_at, r.updated_at
		from app.media_requests r
		join app.users u on u.id = r.requested_by_user_id
		where $2::boolean = true or r.requested_by_user_id = $1
		order by
			case r.status when 'pending' then 0 else 1 end,
			r.created_at desc
	`, userID, includeAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []MediaRequest{}
	for rows.Next() {
		request, err := scanMediaRequest(rows)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, rows.Err()
}

func (s *SettingsStore) GetMediaRequest(ctx context.Context, id uuid.UUID, userID uuid.UUID, includeAll bool) (MediaRequest, error) {
	request, err := scanMediaRequest(s.pool.QueryRow(ctx, `
		select r.id, r.requested_by_user_id, u.username, r.media_type, r.title, r.year,
			r.external_provider, r.external_id, r.overview, r.poster_path, r.status,
			r.quality_profile_id, r.library_folder_id, r.media_item_id, r.decided_at,
			coalesce(array(
				select t.name
				from app.media_request_tags mrt
				join app.tags t on t.id = mrt.tag_id
				where mrt.media_request_id = r.id
				order by lower(t.name)
			), '{}') as tags,
			r.created_at, r.updated_at
		from app.media_requests r
		join app.users u on u.id = r.requested_by_user_id
		where r.id = $1 and ($3::boolean = true or r.requested_by_user_id = $2)
	`, id, userID, includeAll))
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaRequest{}, ErrNotFound
	}
	return request, err
}

func getMediaRequest(ctx context.Context, q mediaItemQuerier, id uuid.UUID) (MediaRequest, error) {
	request, err := scanMediaRequest(q.QueryRow(ctx, `
		select r.id, r.requested_by_user_id, u.username, r.media_type, r.title, r.year,
			r.external_provider, r.external_id, r.overview, r.poster_path, r.status,
			r.quality_profile_id, r.library_folder_id, r.media_item_id, r.decided_at,
			coalesce(array(
				select t.name
				from app.media_request_tags mrt
				join app.tags t on t.id = mrt.tag_id
				where mrt.media_request_id = r.id
				order by lower(t.name)
			), '{}') as tags,
			r.created_at, r.updated_at
		from app.media_requests r
		join app.users u on u.id = r.requested_by_user_id
		where r.id = $1
	`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaRequest{}, ErrNotFound
	}
	return request, err
}

func (s *SettingsStore) CreateMediaRequest(ctx context.Context, input MediaRequestInput) (MediaRequest, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return MediaRequest{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	id := uuid.New()
	var requestID uuid.UUID
	if err := tx.QueryRow(ctx, `
		insert into app.media_requests (
			id, requested_by_user_id, media_type, title, year, external_provider, external_id, overview, poster_path
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		returning id
	`, id, input.RequestedByUserID, input.Type, input.Title, input.Year, input.ExternalProvider, input.ExternalID, input.Overview, input.PosterPath).Scan(&requestID); err != nil {
		return MediaRequest{}, err
	}
	if err := assignMediaRequestTags(ctx, tx, requestID, input.Tags); err != nil {
		return MediaRequest{}, err
	}
	request, err := getMediaRequest(ctx, tx, requestID)
	if err != nil {
		return MediaRequest{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return MediaRequest{}, err
	}
	return request, nil
}

func (s *SettingsStore) ApproveMediaRequest(ctx context.Context, id uuid.UUID, input MediaRequestApprovalInput) (MediaRequest, MediaItem, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return MediaRequest{}, MediaItem{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var folderExists bool
	if err := tx.QueryRow(ctx, `select exists(select 1 from app.library_folders where id = $1)`, input.LibraryFolderID).Scan(&folderExists); err != nil {
		return MediaRequest{}, MediaItem{}, err
	}
	if !folderExists {
		return MediaRequest{}, MediaItem{}, ErrNotFound
	}

	request, err := scanMediaRequest(tx.QueryRow(ctx, `
		select r.id, r.requested_by_user_id, u.username, r.media_type, r.title, r.year,
			r.external_provider, r.external_id, r.overview, r.poster_path, r.status,
			r.quality_profile_id, r.library_folder_id, r.media_item_id, r.decided_at,
			coalesce(array(
				select t.name
				from app.media_request_tags mrt
				join app.tags t on t.id = mrt.tag_id
				where mrt.media_request_id = r.id
				order by lower(t.name)
			), '{}') as tags,
			r.created_at, r.updated_at
		from app.media_requests r
		join app.users u on u.id = r.requested_by_user_id
		where r.id = $1
		for update
	`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaRequest{}, MediaItem{}, ErrNotFound
	}
	if err != nil {
		return MediaRequest{}, MediaItem{}, err
	}
	if request.Status != "pending" {
		return MediaRequest{}, MediaItem{}, ErrRequestClosed
	}

	qualityProfileID := input.QualityProfileID
	libraryFolderID := input.LibraryFolderID
	item, err := createMediaItemIfMissing(ctx, tx, MediaItemInput{
		Type:             request.Type,
		Title:            request.Title,
		Year:             request.Year,
		Monitored:        true,
		ExternalProvider: request.ExternalProvider,
		ExternalID:       request.ExternalID,
		Overview:         request.Overview,
		PosterPath:       request.PosterPath,
		Tags:             request.Tags,
		QualityProfileID: &qualityProfileID,
		LibraryFolderID:  &libraryFolderID,
	})
	if err != nil {
		return MediaRequest{}, MediaItem{}, err
	}

	updated, err := scanMediaRequest(tx.QueryRow(ctx, `
		update app.media_requests
		set status = 'approved',
			quality_profile_id = $2,
			library_folder_id = $3,
			media_item_id = $4,
			decided_at = now(),
			updated_at = now()
		where id = $1
		returning id, requested_by_user_id, (
				select username from app.users where id = requested_by_user_id
			), media_type, title, year, external_provider, external_id, overview, poster_path,
			status, quality_profile_id, library_folder_id, media_item_id, decided_at,
			coalesce(array(
				select t.name
				from app.media_request_tags mrt
				join app.tags t on t.id = mrt.tag_id
				where mrt.media_request_id = $1
				order by lower(t.name)
			), '{}') as tags,
			created_at, updated_at
	`, id, input.QualityProfileID, input.LibraryFolderID, item.ID))
	if err != nil {
		return MediaRequest{}, MediaItem{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return MediaRequest{}, MediaItem{}, err
	}
	return updated, item, nil
}

func scanMediaRequest(row pgx.Row) (MediaRequest, error) {
	var request MediaRequest
	err := row.Scan(
		&request.ID,
		&request.RequestedByUserID,
		&request.RequestedByUsername,
		&request.Type,
		&request.Title,
		&request.Year,
		&request.ExternalProvider,
		&request.ExternalID,
		&request.Overview,
		&request.PosterPath,
		&request.Status,
		&request.QualityProfileID,
		&request.LibraryFolderID,
		&request.MediaItemID,
		&request.DecidedAt,
		&request.Tags,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	return request, err
}
