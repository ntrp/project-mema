package httpapi

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) DeleteMediaItem(w http.ResponseWriter, r *http.Request, id ResourceId, params DeleteMediaItemParams) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	if err := s.settings.DeleteMediaItem(r.Context(), uuid.UUID(id), boolDefault(params.KeepFiles, false)); err != nil {
		writeSettingsError(w, err, "Could not remove media item")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
