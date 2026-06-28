package httpapi

import (
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func mediaItemInput(request MediaItemRequest) (storage.MediaItemInput, bool) {
	title := strings.TrimSpace(request.Title)
	if title == "" || !request.Type.Valid() {
		return storage.MediaItemInput{}, false
	}
	return storage.MediaItemInput{
		Type:      string(request.Type),
		Title:     title,
		Year:      request.Year,
		Monitored: request.Monitored,
	}, true
}

func mediaItemResponse(item storage.MediaItem) MediaItem {
	return MediaItem{
		Id:        openapi_types.UUID(item.ID),
		Type:      MediaType(item.Type),
		Title:     item.Title,
		Year:      item.Year,
		Monitored: item.Monitored,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
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
