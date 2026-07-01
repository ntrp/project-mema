package httpapi

import "net/http"

func (s *Server) GetFileNamingSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	settings, err := s.settings.GetFileNamingSettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_load_failed", "Could not load file naming settings")
		return
	}
	writeJSON(w, http.StatusOK, fileNamingSettingsResponse(settings))
}

func (s *Server) UpdateFileNamingSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body FileNamingSettingsRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := fileNamingSettingsInput(w, body)
	if !ok {
		return
	}

	settings, err := s.settings.SaveFileNamingSettings(r.Context(), input)
	if err != nil {
		writeSettingsError(w, err, "Could not update file naming settings")
		return
	}
	writeJSON(w, http.StatusOK, fileNamingSettingsResponse(settings))
}
