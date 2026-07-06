package httpapi

import (
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *Server) EvaluateMediaComponentCompatibility(
	w http.ResponseWriter,
	r *http.Request,
	id ResourceId,
	sourceID openapi_types.UUID,
) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body MediaComponentCompatibilityEvaluateRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	decision, err := s.settings.EvaluateMediaComponentCompatibility(
		r.Context(),
		uuid.UUID(id),
		uuid.UUID(sourceID),
		uuid.UUID(body.BaseSourceId),
	)
	if err != nil {
		writeSettingsError(w, err, "Could not evaluate component compatibility")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "media", "Component compatibility evaluated", map[string]any{
		"mediaItemId":     uuid.UUID(id).String(),
		"sourceId":        uuid.UUID(sourceID).String(),
		"baseSourceId":    uuid.UUID(body.BaseSourceId).String(),
		"confidenceState": decision.ConfidenceState,
		"automationState": decision.AutomationState,
	})
	writeJSON(w, http.StatusOK, mediaComponentCompatibilityResponse(decision))
}

func (s *Server) ReviewMediaComponentCompatibility(
	w http.ResponseWriter,
	r *http.Request,
	id ResourceId,
	sourceID openapi_types.UUID,
	decisionID openapi_types.UUID,
) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body MediaComponentCompatibilityReviewRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	decision, err := s.settings.ReviewMediaComponentCompatibility(
		r.Context(),
		uuid.UUID(id),
		uuid.UUID(sourceID),
		uuid.UUID(decisionID),
		mediaComponentCompatibilityReviewInput(body),
	)
	if err != nil {
		writeSettingsError(w, err, "Could not review component compatibility")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "media", "Component compatibility reviewed", map[string]any{
		"mediaItemId":     uuid.UUID(id).String(),
		"sourceId":        uuid.UUID(sourceID).String(),
		"decisionId":      uuid.UUID(decisionID).String(),
		"reviewState":     decision.ReviewState,
		"automationState": decision.AutomationState,
	})
	writeJSON(w, http.StatusOK, mediaComponentCompatibilityResponse(decision))
}
