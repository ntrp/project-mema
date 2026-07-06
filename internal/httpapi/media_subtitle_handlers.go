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
