package httpapi

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) UpdateMediaItemMode(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body MediaItemModeRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	if !body.Mode.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_media_mode", "Media mode must be manual or automatic")
		return
	}

	item, err := s.settings.UpdateMediaItemManual(r.Context(), uuid.UUID(id), body.Mode == Manual)
	if err != nil {
		writeSettingsError(w, err, "Could not update media mode")
		return
	}
	writeJSON(w, http.StatusOK, mediaItemResponse(item))
}
