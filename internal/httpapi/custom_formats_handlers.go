package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
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

func (s *Server) TestCustomFormatParsing(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body CustomFormatParsingRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	fileName := strings.TrimSpace(body.FileName)
	if fileName == "" {
		writeError(w, http.StatusBadRequest, "invalid_file_name", "File name is required")
		return
	}

	formats, err := s.settings.ListCustomFormats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_list_failed", "Could not list custom formats")
		return
	}
	parsed := decisions.ParseReleaseFileName(fileName)
	matches := decisions.MatchCustomFormats(parsed, formats)
	profile, err := s.customFormatParsingProfile(r, parsed)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "media_match_failed", "Could not match media profile")
		return
	}
	writeJSON(w, http.StatusOK, customFormatParsingResponse(parsed, matches, profile))
}

func (s *Server) customFormatParsingProfile(
	r *http.Request,
	parsed decisions.ParsedRelease,
) (*storage.MediaProfile, error) {
	item, err := s.settings.FindMonitoredMediaMatch(r.Context(), parsed.MovieTitle, parsed.Year)
	if errors.Is(err, storage.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if item.QualityProfileID == nil {
		return nil, nil
	}
	profile, err := s.settings.GetMediaProfile(r.Context(), *item.QualityProfileID)
	if errors.Is(err, storage.ErrNotFound) {
		return nil, nil
	}
	return &profile, err
}
