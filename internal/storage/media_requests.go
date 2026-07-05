package storage

import (
	"context"
	"errors"
	"time"

	storagegen "media-manager/internal/storage/generated"

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
	MonitorMode         string
	SeriesType          *string
	MinimumAvailability string
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
	RequestedByUserID   uuid.UUID
	Type                string
	Title               string
	Year                *int32
	ExternalProvider    *string
	ExternalID          *string
	Overview            *string
	PosterPath          *string
	MonitorMode         string
	SeriesType          *string
	MinimumAvailability string
	Tags                []string
}

type MediaRequestApprovalInput struct {
	QualityProfileID string
	LibraryFolderID  uuid.UUID
	MediaInput       *MediaItemInput
}

func (s *SettingsStore) ListMediaRequests(ctx context.Context, userID uuid.UUID, includeAll bool) ([]MediaRequest, error) {
	rows, err := storagegen.New(s.pool).ListMediaRequests(ctx, storagegen.ListMediaRequestsParams{
		IncludeAll: includeAll,
		UserID:     userID,
	})
	if err != nil {
		return nil, err
	}

	requests := make([]MediaRequest, 0, len(rows))
	for _, row := range rows {
		requests = append(requests, mediaRequestFromListRow(row))
	}
	return requests, nil
}

func (s *SettingsStore) GetMediaRequest(ctx context.Context, id uuid.UUID, userID uuid.UUID, includeAll bool) (MediaRequest, error) {
	row, err := storagegen.New(s.pool).GetMediaRequestForUser(ctx, storagegen.GetMediaRequestForUserParams{
		ID:         id,
		IncludeAll: includeAll,
		UserID:     userID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaRequest{}, ErrNotFound
	}
	return mediaRequestFromGetForUserRow(row), err
}

func getMediaRequest(ctx context.Context, q mediaItemQuerier, id uuid.UUID) (MediaRequest, error) {
	row, err := storagegen.New(q).GetMediaRequest(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaRequest{}, ErrNotFound
	}
	return mediaRequestFromGetRow(row), err
}

func (s *SettingsStore) CreateMediaRequest(ctx context.Context, input MediaRequestInput) (MediaRequest, error) {
	input = normalizeMediaRequestOptions(input)
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return MediaRequest{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	id := uuid.New()
	requestID, err := storagegen.New(tx).CreateMediaRequest(ctx, mediaRequestCreateParams(id, input))
	if err != nil {
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

	if _, err := storagegen.New(tx).GetLibraryFolder(ctx, input.LibraryFolderID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return MediaRequest{}, MediaItem{}, ErrNotFound
		}
		return MediaRequest{}, MediaItem{}, err
	}

	row, err := storagegen.New(tx).GetMediaRequestForUpdate(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaRequest{}, MediaItem{}, ErrNotFound
	}
	if err != nil {
		return MediaRequest{}, MediaItem{}, err
	}
	request := mediaRequestFromUpdateRow(row)
	if request.Status != "pending" {
		return MediaRequest{}, MediaItem{}, ErrRequestClosed
	}

	qualityProfileID := input.QualityProfileID
	libraryFolderID := input.LibraryFolderID
	mediaInput := MediaItemInput{
		Type:                request.Type,
		Title:               request.Title,
		Year:                request.Year,
		Monitored:           true,
		ExternalProvider:    request.ExternalProvider,
		ExternalID:          request.ExternalID,
		Overview:            request.Overview,
		PosterPath:          request.PosterPath,
		MonitorMode:         request.MonitorMode,
		SeriesType:          request.SeriesType,
		MinimumAvailability: request.MinimumAvailability,
		Tags:                request.Tags,
		QualityProfileID:    &qualityProfileID,
		LibraryFolderID:     &libraryFolderID,
	}
	if input.MediaInput != nil {
		mediaInput = *input.MediaInput
		mediaInput.QualityProfileID = &qualityProfileID
		mediaInput.LibraryFolderID = &libraryFolderID
		mediaInput.Tags = request.Tags
		mediaInput.SeriesType = request.SeriesType
	}
	item, err := createMediaItemIfMissing(ctx, tx, mediaInput)
	if err != nil {
		return MediaRequest{}, MediaItem{}, err
	}

	updatedRow, err := storagegen.New(tx).ApproveMediaRequest(ctx, storagegen.ApproveMediaRequestParams{
		QualityProfileID: textValue(&input.QualityProfileID),
		LibraryFolderID:  &input.LibraryFolderID,
		MediaItemID:      &item.ID,
		ID:               id,
	})
	if err != nil {
		return MediaRequest{}, MediaItem{}, err
	}
	updated := mediaRequestFromApproveRow(updatedRow)
	if err := tx.Commit(ctx); err != nil {
		return MediaRequest{}, MediaItem{}, err
	}
	return updated, item, nil
}

func mediaRequestCreateParams(id uuid.UUID, input MediaRequestInput) storagegen.CreateMediaRequestParams {
	return storagegen.CreateMediaRequestParams{
		ID:                  id,
		RequestedByUserID:   input.RequestedByUserID,
		MediaType:           input.Type,
		Title:               input.Title,
		Year:                int4Value(input.Year),
		ExternalProvider:    textValue(input.ExternalProvider),
		ExternalID:          textValue(input.ExternalID),
		Overview:            textValue(input.Overview),
		PosterPath:          textValue(input.PosterPath),
		MonitorMode:         input.MonitorMode,
		SeriesType:          textValue(input.SeriesType),
		MinimumAvailability: input.MinimumAvailability,
	}
}

func mediaRequestFromListRow(row storagegen.ListMediaRequestsRow) MediaRequest {
	return mediaRequestFromGetRow(storagegen.GetMediaRequestRow(row))
}

func mediaRequestFromGetForUserRow(row storagegen.GetMediaRequestForUserRow) MediaRequest {
	return mediaRequestFromGetRow(storagegen.GetMediaRequestRow(row))
}

func mediaRequestFromGetRow(row storagegen.GetMediaRequestRow) MediaRequest {
	return MediaRequest{
		ID:                  row.ID,
		RequestedByUserID:   row.RequestedByUserID,
		RequestedByUsername: row.RequestedByUsername,
		Type:                row.MediaType,
		Title:               row.Title,
		Year:                int4Ptr(row.Year),
		ExternalProvider:    textPtr(row.ExternalProvider),
		ExternalID:          textPtr(row.ExternalID),
		Overview:            textPtr(row.Overview),
		PosterPath:          textPtr(row.PosterPath),
		Status:              row.Status,
		MonitorMode:         row.MonitorMode,
		SeriesType:          textPtr(row.SeriesType),
		MinimumAvailability: row.MinimumAvailability,
		QualityProfileID:    textPtr(row.QualityProfileID),
		LibraryFolderID:     row.LibraryFolderID,
		MediaItemID:         row.MediaItemID,
		DecidedAt:           row.DecidedAt,
		Tags:                row.Tags,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
	}
}

func mediaRequestFromUpdateRow(row storagegen.GetMediaRequestForUpdateRow) MediaRequest {
	return mediaRequestFromGetRow(storagegen.GetMediaRequestRow(row))
}

func mediaRequestFromApproveRow(row storagegen.ApproveMediaRequestRow) MediaRequest {
	return mediaRequestFromGetRow(storagegen.GetMediaRequestRow(row))
}
