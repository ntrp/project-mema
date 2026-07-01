package httpapi

import "net/http"

func (s *Server) ListQualitySizeSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	settings, err := s.settings.ListQualitySizeSettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_list_failed", "Could not list quality sizes")
		return
	}
	writeJSON(w, http.StatusOK, qualitySizeSettingsResponse(settings))
}

func (s *Server) UpdateQualitySizeSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body QualitySizeSettingsUpdateRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := qualitySizeSettingsInput(w, body)
	if !ok {
		return
	}

	settings, err := s.settings.SaveQualitySizeSettings(r.Context(), input)
	if err != nil {
		writeSettingsError(w, err, "Could not update quality sizes")
		return
	}
	writeJSON(w, http.StatusOK, qualitySizeSettingsResponse(settings))
}
