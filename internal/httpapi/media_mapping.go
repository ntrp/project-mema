package httpapi

import (
	"strings"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func mediaItemInput(request MediaItemRequest) (storage.MediaItemInput, bool) {
	title := strings.TrimSpace(request.Title)
	if title == "" || !request.Type.Valid() {
		return storage.MediaItemInput{}, false
	}
	return storage.MediaItemInput{
		Type:             string(request.Type),
		Title:            title,
		Year:             request.Year,
		Monitored:        request.Monitored,
		ExternalProvider: optionalTrimmedString(request.ExternalProvider),
		ExternalID:       optionalTrimmedString(request.ExternalId),
		Overview:         optionalTrimmedString(request.Overview),
		PosterPath:       optionalTrimmedString(request.PosterPath),
		QualityProfileID: optionalTrimmedString(request.QualityProfileId),
		LibraryFolderID:  optionalUUID(request.LibraryFolderId),
		Tags:             optionalStringSlice(request.Tags),
	}, true
}

func mediaItemResponse(item storage.MediaItem) MediaItem {
	return MediaItem{
		Id:               openapi_types.UUID(item.ID),
		Type:             MediaType(item.Type),
		Title:            item.Title,
		Year:             item.Year,
		Monitored:        item.Monitored,
		ExternalProvider: item.ExternalProvider,
		ExternalId:       item.ExternalID,
		Overview:         item.Overview,
		PosterPath:       item.PosterPath,
		QualityProfileId: item.QualityProfileID,
		LibraryFolderId:  optionalOpenAPIUUID(item.LibraryFolderID),
		Tags:             &item.Tags,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
	}
}

func mediaRequestInput(request MediaRequestCreateRequest, requestedByUserID uuid.UUID) (storage.MediaRequestInput, bool) {
	title := strings.TrimSpace(request.Title)
	if title == "" || !request.Type.Valid() {
		return storage.MediaRequestInput{}, false
	}
	return storage.MediaRequestInput{
		RequestedByUserID: requestedByUserID,
		Type:              string(request.Type),
		Title:             title,
		Year:              request.Year,
		ExternalProvider:  optionalTrimmedString(request.ExternalProvider),
		ExternalID:        optionalTrimmedString(request.ExternalId),
		Overview:          optionalTrimmedString(request.Overview),
		PosterPath:        optionalTrimmedString(request.PosterPath),
		Tags:              optionalStringSlice(request.Tags),
	}, true
}

func mediaRequestResponse(request storage.MediaRequest) MediaRequest {
	return MediaRequest{
		Id:                  openapi_types.UUID(request.ID),
		RequestedByUserId:   openapi_types.UUID(request.RequestedByUserID),
		RequestedByUsername: request.RequestedByUsername,
		Type:                MediaType(request.Type),
		Title:               request.Title,
		Year:                request.Year,
		ExternalProvider:    request.ExternalProvider,
		ExternalId:          request.ExternalID,
		Overview:            request.Overview,
		PosterPath:          request.PosterPath,
		Tags:                &request.Tags,
		Status:              MediaRequestStatus(request.Status),
		QualityProfileId:    request.QualityProfileID,
		LibraryFolderId:     optionalOpenAPIUUID(request.LibraryFolderID),
		MediaItemId:         optionalOpenAPIUUID(request.MediaItemID),
		DecidedAt:           request.DecidedAt,
		CreatedAt:           request.CreatedAt,
		UpdatedAt:           request.UpdatedAt,
	}
}

func optionalStringSlice(values *[]string) []string {
	if values == nil {
		return nil
	}
	return append([]string(nil), (*values)...)
}

func releaseCandidateResponse(release storage.ReleaseCandidate) ReleaseCandidate {
	var indexerID *openapi_types.UUID
	if release.IndexerID != nil {
		value := openapi_types.UUID(*release.IndexerID)
		indexerID = &value
	}
	return ReleaseCandidate{
		Id:          openapi_types.UUID(release.ID),
		IndexerId:   indexerID,
		IndexerName: release.IndexerName,
		IndexerType: IndexerType(release.IndexerType),
		Title:       release.Title,
		InfoUrl:     release.InfoURL,
		Guid:        release.GUID,
		SizeBytes:   release.SizeBytes,
		Seeders:     release.Seeders,
		Peers:       release.Peers,
	}
}

func downloadActivityResponse(activity storage.DownloadActivity) DownloadActivity {
	return DownloadActivity{
		Id:                 openapi_types.UUID(activity.ID),
		MediaItemId:        openapi_types.UUID(activity.MediaItemID),
		MediaTitle:         activity.MediaTitle,
		MediaType:          MediaType(activity.MediaType),
		ReleaseTitle:       activity.ReleaseTitle,
		IndexerName:        activity.IndexerName,
		DownloadClientName: activity.DownloadClientName,
		DownloadUrl:        activity.DownloadURL,
		Status:             DownloadActivityStatus(activity.Status),
		Error:              activity.Error,
		CreatedAt:          activity.CreatedAt,
		UpdatedAt:          activity.UpdatedAt,
	}
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func optionalUUID(value *openapi_types.UUID) *uuid.UUID {
	if value == nil {
		return nil
	}
	converted := uuid.UUID(*value)
	return &converted
}

func optionalOpenAPIUUID(value *uuid.UUID) *openapi_types.UUID {
	if value == nil {
		return nil
	}
	converted := openapi_types.UUID(*value)
	return &converted
}
