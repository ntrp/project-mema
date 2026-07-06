package httpapi

import (
	"strings"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func downloadActivityResponse(activity storage.DownloadActivity) DownloadActivity {
	return DownloadActivity{
		Id:                 openapi_types.UUID(activity.ID),
		MediaItemId:        openapi_types.UUID(activity.MediaItemID),
		MediaTitle:         activity.MediaTitle,
		MediaType:          MediaType(activity.MediaType),
		MediaYear:          activity.MediaYear,
		ReleaseTitle:       activity.ReleaseTitle,
		IndexerName:        activity.IndexerName,
		DownloadClientName: activity.DownloadClientName,
		DownloadId:         activity.DownloadID,
		DownloadUrl:        activity.DownloadURL,
		Status:             DownloadActivityStatus(activity.Status),
		ProgressPercent:    activity.ProgressPercent,
		Error:              activity.Error,
		FailureType:        downloadActivityFailureType(activity.FailureType),
		CreatedAt:          activity.CreatedAt,
		UpdatedAt:          activity.UpdatedAt,
	}
}

func downloadActivityFailureType(value *string) *DownloadActivityFailureType {
	if value == nil {
		return nil
	}
	failureType := DownloadActivityFailureType(*value)
	return &failureType
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
