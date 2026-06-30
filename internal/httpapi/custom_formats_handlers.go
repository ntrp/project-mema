package httpapi

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) ListCustomFormats(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	formats, err := s.settings.ListCustomFormats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_list_failed", "Could not list custom formats")
		return
	}
	writeJSON(w, http.StatusOK, customFormatListResponse(formats))
}

func (s *Server) CreateCustomFormat(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body CustomFormatRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := customFormatInput(w, body)
	if !ok {
		return
	}

	format, err := s.settings.CreateCustomFormat(r.Context(), input)
	if err != nil {
		writeSettingsError(w, err, "Could not create custom format")
		return
	}
	writeJSON(w, http.StatusCreated, customFormatResponse(format))
}

func (s *Server) UpdateCustomFormat(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body CustomFormatRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := customFormatInput(w, body)
	if !ok {
		return
	}

	format, err := s.settings.UpdateCustomFormat(r.Context(), uuid.UUID(id), input)
	if err != nil {
		writeSettingsError(w, err, "Could not update custom format")
		return
	}
	writeJSON(w, http.StatusOK, customFormatResponse(format))
}

func (s *Server) DeleteCustomFormat(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	if err := s.settings.DeleteCustomFormat(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete custom format")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
