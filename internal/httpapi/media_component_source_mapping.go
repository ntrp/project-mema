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
		ReleaseGroup:    optionalTrimmedString(request.ReleaseGroup),
		ReleaseName:     optionalTrimmedString(request.ReleaseName),
		ReleaseID:       optionalTrimmedString(request.ReleaseId),
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
	compatibility := mediaComponentCompatibilityResponses(value.Compatibility)
	return MediaComponentSource{
		Id:              openapi_types.UUID(value.ID),
		MediaItemId:     openapi_types.UUID(value.MediaItemID),
		SourceRole:      MediaComponentSourceRole(value.SourceRole),
		SourceFilePath:  value.SourceFilePath,
		RetainedPath:    value.RetainedPath,
		ReleaseTitle:    value.ReleaseTitle,
		ReleaseGroup:    value.ReleaseGroup,
		ReleaseName:     value.ReleaseName,
		ReleaseId:       value.ReleaseID,
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
		Compatibility:   &compatibility,
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

func mediaComponentCompatibilityReviewInput(
	request MediaComponentCompatibilityReviewRequest,
) storage.MediaComponentCompatibilityReviewInput {
	return storage.MediaComponentCompatibilityReviewInput{
		ReviewState: string(request.ReviewState),
		Reason:      optionalTrimmedString(request.Reason),
	}
}

func mediaComponentCompatibilityResponses(
	values []storage.MediaComponentCompatibilityDecision,
) []MediaComponentCompatibilityDecision {
	items := make([]MediaComponentCompatibilityDecision, 0, len(values))
	for _, value := range values {
		items = append(items, mediaComponentCompatibilityResponse(value))
	}
	return items
}

func mediaComponentCompatibilityResponse(
	value storage.MediaComponentCompatibilityDecision,
) MediaComponentCompatibilityDecision {
	return MediaComponentCompatibilityDecision{
		Id:                openapi_types.UUID(value.ID),
		MediaItemId:       openapi_types.UUID(value.MediaItemID),
		BaseSourceId:      openapi_types.UUID(value.BaseSourceID),
		ComponentSourceId: openapi_types.UUID(value.ComponentSourceID),
		ConfidenceState:   MediaComponentCompatibilityConfidenceState(value.ConfidenceState),
		AutomationState:   MediaComponentCompatibilityAutomationState(value.AutomationState),
		ReviewState:       MediaComponentCompatibilityReviewState(value.ReviewState),
		Reason:            value.Reason,
		RuntimeDeltaMs:    value.RuntimeDeltaMs,
		Evidence:          value.Evidence,
		ReviewReason:      value.ReviewReason,
		ReviewedAt:        value.ReviewedAt,
		CreatedAt:         value.CreatedAt,
		UpdatedAt:         value.UpdatedAt,
	}
}
