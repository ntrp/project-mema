package httpapi

import "net/http"

func (s *Server) ListMediaProfiles(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	profiles, err := s.settings.ListMediaProfiles(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_list_failed", "Could not list media profiles")
		return
	}
	writeJSON(w, http.StatusOK, mediaProfileListResponse(profiles))
}

func (s *Server) CreateMediaProfile(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body MediaProfileRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := mediaProfileInput(w, body)
	if !ok {
		return
	}

	profile, err := s.settings.CreateMediaProfile(r.Context(), input)
	if err != nil {
		writeSettingsError(w, err, "Could not create media profile")
		return
	}
	writeJSON(w, http.StatusCreated, mediaProfileResponse(profile))
}

func (s *Server) UpdateMediaProfile(w http.ResponseWriter, r *http.Request, id ProfileId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body MediaProfileRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := mediaProfileInput(w, body)
	if !ok {
		return
	}

	profile, err := s.settings.UpdateMediaProfile(r.Context(), string(id), input)
	if err != nil {
		writeSettingsError(w, err, "Could not update media profile")
		return
	}
	writeJSON(w, http.StatusOK, mediaProfileResponse(profile))
}

func (s *Server) DeleteMediaProfile(w http.ResponseWriter, r *http.Request, id ProfileId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	if err := s.settings.DeleteMediaProfile(r.Context(), string(id)); err != nil {
		writeSettingsError(w, err, "Could not delete media profile")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
