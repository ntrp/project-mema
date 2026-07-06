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
	artifacts := mediaComponentArtifactResponses(value.Artifacts)
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
		Artifacts:       &artifacts,
	}
}

func mediaComponentArtifactInput(request MediaComponentExtractionRequest) storage.MediaComponentArtifactInput {
	return storage.MediaComponentArtifactInput{
		StreamID:   request.StreamId,
		StreamType: string(request.StreamType),
		Language:   optionalTrimmedString(request.Language),
	}
}

func mediaComponentArtifactResponses(values []storage.MediaComponentArtifact) []MediaComponentArtifact {
	items := make([]MediaComponentArtifact, 0, len(values))
	for _, value := range values {
		items = append(items, mediaComponentArtifactResponse(value))
	}
	return items
}

func mediaComponentArtifactResponse(value storage.MediaComponentArtifact) MediaComponentArtifact {
	return MediaComponentArtifact{
		Id:           openapi_types.UUID(value.ID),
		MediaItemId:  openapi_types.UUID(value.MediaItemID),
		SourceId:     openapi_types.UUID(value.SourceID),
		StreamId:     value.StreamID,
		StreamType:   MediaComponentArtifactStreamType(value.StreamType),
		Language:     value.Language,
		OutputPath:   value.OutputPath,
		Status:       MediaComponentArtifactStatus(value.Status),
		ToolName:     value.ToolName,
		ToolSummary:  value.ToolSummary,
		ErrorMessage: value.ErrorMessage,
		JobId:        value.JobID,
		SizeBytes:    value.SizeBytes,
		CreatedAt:    value.CreatedAt,
		UpdatedAt:    value.UpdatedAt,
		CompletedAt:  value.CompletedAt,
	}
}
