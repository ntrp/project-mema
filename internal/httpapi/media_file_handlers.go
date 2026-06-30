package httpapi

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) RescanMediaItemFiles(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	item, err := s.settings.RescanMediaItemFiles(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not rescan media folder")
		return
	}
	writeJSON(w, http.StatusOK, mediaItemResponse(item))
}
