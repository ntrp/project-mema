package httpapi

import "net/http"

func (s *Server) ListLanguages(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	languages, err := s.settings.ListLanguages(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_list_failed", "Could not list languages")
		return
	}
	response := LanguageListResponse{Languages: make([]Language, 0, len(languages))}
	for _, language := range languages {
		response.Languages = append(response.Languages, languageResponse(language))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreateLanguage(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body LanguageRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := languageInput(w, body)
	if !ok {
		return
	}
	language, err := s.settings.SaveLanguage(r.Context(), "", input)
	if err != nil {
		writeSettingsError(w, err, "Could not create language")
		return
	}
	writeJSON(w, http.StatusCreated, languageResponse(language))
}

func (s *Server) UpdateLanguage(w http.ResponseWriter, r *http.Request, code string) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body LanguageUpdateRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := languageUpdateInput(w, body)
	if !ok {
		return
	}
	language, err := s.settings.SaveLanguage(r.Context(), code, input)
	if err != nil {
		writeSettingsError(w, err, "Could not update language")
		return
	}
	writeJSON(w, http.StatusOK, languageResponse(language))
}

func (s *Server) DeleteLanguage(w http.ResponseWriter, r *http.Request, code string) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	if err := s.settings.DeleteLanguage(r.Context(), code); err != nil {
		writeSettingsError(w, err, "Could not delete language")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
