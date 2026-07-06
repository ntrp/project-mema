package httpapi

import (
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *Server) ListMediaComponentSources(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	item, err := s.settings.GetMediaItem(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not load component sources")
		return
	}
	writeJSON(w, http.StatusOK, MediaComponentSourceListResponse{
		Sources: mediaComponentSourceResponses(item.ComponentSources),
	})
}

func (s *Server) RetainMediaComponentSource(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body MediaComponentSourceRetainRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	source, err := s.settings.RetainMediaComponentSource(
		r.Context(),
		uuid.UUID(id),
		mediaComponentSourceInput(body),
	)
	if err != nil {
		writeSettingsError(w, err, "Could not retain component source")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "media", "Component source retained", map[string]any{
		"mediaItemId": uuid.UUID(id).String(),
		"sourceId":    source.ID.String(),
		"sourceRole":  source.SourceRole,
	})
	writeJSON(w, http.StatusCreated, mediaComponentSourceResponse(source))
}

func (s *Server) GetMediaComponentSource(
	w http.ResponseWriter,
	r *http.Request,
	id ResourceId,
	sourceID openapi_types.UUID,
) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	source, err := s.settings.GetMediaComponentSource(r.Context(), uuid.UUID(id), uuid.UUID(sourceID))
	if err != nil {
		writeSettingsError(w, err, "Could not load component source")
		return
	}
	writeJSON(w, http.StatusOK, mediaComponentSourceResponse(source))
}

func (s *Server) ReleaseMediaComponentSource(
	w http.ResponseWriter,
	r *http.Request,
	id ResourceId,
	sourceID openapi_types.UUID,
) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	source, err := s.settings.ReleaseMediaComponentSource(r.Context(), uuid.UUID(id), uuid.UUID(sourceID))
	if err != nil {
		writeSettingsError(w, err, "Could not release component source")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "media", "Component source released", map[string]any{
		"mediaItemId": uuid.UUID(id).String(),
		"sourceId":    source.ID.String(),
		"sourceRole":  source.SourceRole,
	})
	writeJSON(w, http.StatusOK, mediaComponentSourceResponse(source))
}
