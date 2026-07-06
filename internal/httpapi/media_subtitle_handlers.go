package httpapi

import (
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *Server) ListMediaItemSubtitles(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	item, err := s.settings.GetMediaItem(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not load media subtitles")
		return
	}
	writeJSON(w, http.StatusOK, MediaItemSubtitleListResponse{
		Subtitles: mediaSubtitleResponses(item.ExternalSubtitles),
	})
}

func (s *Server) DeleteMediaItemSubtitle(
	w http.ResponseWriter,
	r *http.Request,
	id ResourceId,
	subtitleID openapi_types.UUID,
) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	item, err := s.settings.DeleteMediaItemSubtitle(r.Context(), uuid.UUID(id), uuid.UUID(subtitleID))
	if err != nil {
		writeSettingsError(w, err, "Could not delete subtitle")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "subtitles", "Subtitle deleted", map[string]any{
		"mediaItemId": uuid.UUID(id).String(),
		"subtitleId":  uuid.UUID(subtitleID).String(),
	})
	writeJSON(w, http.StatusOK, mediaItemResponse(item))
}

func (s *Server) UpdateMediaItemSubtitle(
	w http.ResponseWriter,
	r *http.Request,
	id ResourceId,
	subtitleID openapi_types.UUID,
) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body MediaItemSubtitleSelectionRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	item, err := s.settings.UpdateMediaItemSubtitleSelection(
		r.Context(),
		uuid.UUID(id),
		uuid.UUID(subtitleID),
		mediaSubtitleSelectionInput(body),
	)
	if err != nil {
		writeSettingsError(w, err, "Could not update subtitle")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "subtitles", "Subtitle selection updated", map[string]any{
		"mediaItemId":   uuid.UUID(id).String(),
		"subtitleId":    uuid.UUID(subtitleID).String(),
		"selected":      body.Selected,
		"retentionMode": body.RetentionMode,
	})
	writeJSON(w, http.StatusOK, mediaItemResponse(item))
}
