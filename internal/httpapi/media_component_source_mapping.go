package httpapi

import (
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func mediaComponentSourceInput(request MediaComponentSourceRetainRequest) storage.MediaComponentSourceInput {
	return storage.MediaComponentSourceInput{
		SourceRole:      string(request.SourceRole),
		SourceFilePath:  request.SourceFilePath,
		ReleaseTitle:    optionalTrimmedString(request.ReleaseTitle),
		SourceMetadata:  optionalTrimmedString(request.SourceMetadata),
		StreamInventory: stringValue(request.StreamInventory),
		Checksum:        optionalTrimmedString(request.Checksum),
	}
}

func mediaComponentSourceResponses(values []storage.MediaComponentSource) []MediaComponentSource {
	items := make([]MediaComponentSource, 0, len(values))
	for _, value := range values {
		items = append(items, mediaComponentSourceResponse(value))
	}
	return items
}

func mediaComponentSourceResponse(value storage.MediaComponentSource) MediaComponentSource {
	return MediaComponentSource{
		Id:              openapi_types.UUID(value.ID),
		MediaItemId:     openapi_types.UUID(value.MediaItemID),
		SourceRole:      MediaComponentSourceRole(value.SourceRole),
		SourceFilePath:  value.SourceFilePath,
		RetainedPath:    value.RetainedPath,
		ReleaseTitle:    value.ReleaseTitle,
		SourceMetadata:  value.SourceMetadata,
		StreamInventory: value.StreamInventory,
		Checksum:        value.Checksum,
		SizeBytes:       value.SizeBytes,
		RetentionState:  MediaComponentSourceRetentionState(value.RetentionState),
		RetainedAt:      value.RetainedAt,
		ReleasedAt:      value.ReleasedAt,
		CreatedAt:       value.CreatedAt,
		UpdatedAt:       value.UpdatedAt,
	}
}
