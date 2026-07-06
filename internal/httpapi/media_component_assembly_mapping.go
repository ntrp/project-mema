package httpapi

import (
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func mediaComponentAssemblyInput(request MediaComponentAssemblyRequest) storage.MediaComponentAssemblyRunInput {
	artifactIDs := make([]openapi_types.UUID, 0, len(request.ArtifactIds))
	artifactIDs = append(artifactIDs, request.ArtifactIds...)
	return storage.MediaComponentAssemblyRunInput{
		BaseSourceID: uuidFromOpenAPI(request.BaseSourceId),
		ArtifactIDs:  uuidSliceFromOpenAPI(artifactIDs),
	}
}

func mediaComponentAssemblyRunResponses(values []storage.MediaComponentAssemblyRun) []MediaComponentAssemblyRun {
	items := make([]MediaComponentAssemblyRun, 0, len(values))
	for _, value := range values {
		items = append(items, mediaComponentAssemblyRunResponse(value))
	}
	return items
}

func mediaComponentAssemblyRunResponse(value storage.MediaComponentAssemblyRun) MediaComponentAssemblyRun {
	return MediaComponentAssemblyRun{
		Id:           openapi_types.UUID(value.ID),
		MediaItemId:  openapi_types.UUID(value.MediaItemID),
		BaseSourceId: openapi_types.UUID(value.BaseSourceID),
		OutputPath:   value.OutputPath,
		Status:       MediaComponentAssemblyStatus(value.Status),
		ToolName:     value.ToolName,
		ToolSummary:  value.ToolSummary,
		ErrorMessage: value.ErrorMessage,
		JobId:        value.JobID,
		SizeBytes:    value.SizeBytes,
		CreatedAt:    value.CreatedAt,
		UpdatedAt:    value.UpdatedAt,
		CompletedAt:  value.CompletedAt,
		Inputs:       mediaComponentAssemblyInputResponses(value.Inputs),
	}
}

func mediaComponentAssemblyInputResponses(values []storage.MediaComponentAssemblyInput) []MediaComponentAssemblyInput {
	items := make([]MediaComponentAssemblyInput, 0, len(values))
	for _, value := range values {
		items = append(items, MediaComponentAssemblyInput{
			Id:         openapi_types.UUID(value.ID),
			RunId:      openapi_types.UUID(value.RunID),
			SourceId:   optionalOpenAPIUUID(value.SourceID),
			ArtifactId: optionalOpenAPIUUID(value.ArtifactID),
			StreamType: MediaComponentAssemblyStreamType(value.StreamType),
			InputPath:  value.InputPath,
			Provenance: value.Provenance,
			CreatedAt:  value.CreatedAt,
		})
	}
	return items
}

func uuidFromOpenAPI(value openapi_types.UUID) uuid.UUID {
	return uuid.UUID(value)
}

func uuidSliceFromOpenAPI(values []openapi_types.UUID) []uuid.UUID {
	items := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		items = append(items, uuid.UUID(value))
	}
	return items
}
