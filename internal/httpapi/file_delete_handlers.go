package httpapi

import (
	"net/http"

	"media-manager/internal/storage"
)

func (s *Server) GetFileDeleteSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	settings, err := s.settings.GetFileDeleteSettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_load_failed", "Could not load file delete settings")
		return
	}
	writeJSON(w, http.StatusOK, fileDeleteSettingsResponse(settings))
}

func (s *Server) UpdateFileDeleteSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body FileDeleteSettingsRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	settings, err := s.settings.UpdateFileDeleteSettings(r.Context(), storage.FileDeleteSettingsInput{
		Mode:          string(body.Mode),
		RecycleFolder: body.RecycleFolder,
	})
	if err != nil {
		writeSettingsError(w, err, "Could not update file delete settings")
		return
	}
	writeJSON(w, http.StatusOK, fileDeleteSettingsResponse(settings))
}
