package httpapi

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *Server) EnqueueMediaComponentExtraction(
	w http.ResponseWriter,
	r *http.Request,
	id ResourceId,
	sourceID openapi_types.UUID,
) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body MediaComponentExtractionRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	artifact, err := s.settings.CreateMediaComponentArtifact(
		r.Context(),
		uuid.UUID(id),
		uuid.UUID(sourceID),
		mediaComponentArtifactInput(body),
	)
	if err != nil {
		writeSettingsError(w, err, "Could not create component extraction")
		return
	}
	jobID, err := s.jobs.EnqueueMediaComponentExtraction(r.Context(), artifact.ID)
	if err != nil {
		_, _ = s.settings.FailMediaComponentArtifact(r.Context(), artifact.ID, "", err.Error())
		s.recordEvent(r.Context(), eventSeverityError, "media", "Component extraction enqueue failed", map[string]any{
			"mediaItemId": uuid.UUID(id).String(),
			"sourceId":    uuid.UUID(sourceID).String(),
			"artifactId":  artifact.ID.String(),
			"error":       err.Error(),
		})
		writeError(w, http.StatusInternalServerError, "component_extraction_enqueue_failed", "Could not enqueue component extraction")
		return
	}
	artifact, err = s.settings.AssignMediaComponentArtifactJob(r.Context(), artifact.ID, strconv.FormatInt(jobID, 10))
	if err != nil {
		writeSettingsError(w, err, "Could not update component extraction")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "media", "Component extraction queued", map[string]any{
		"mediaItemId": uuid.UUID(id).String(),
		"sourceId":    uuid.UUID(sourceID).String(),
		"artifactId":  artifact.ID.String(),
		"jobId":       jobID,
	})
	writeJSON(w, http.StatusAccepted, MediaComponentExtractionEnqueueResponse{
		JobId:    jobID,
		Message:  "Component extraction queued",
		Artifact: mediaComponentArtifactResponse(artifact),
	})
}
