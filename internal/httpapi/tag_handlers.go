package httpapi

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) ListTags(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	tags, err := s.settings.ListTags(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_list_failed", "Could not list tags")
		return
	}
	response := TagListResponse{Tags: make([]Tag, 0, len(tags))}
	for _, tag := range tags {
		response.Tags = append(response.Tags, tagResponse(tag))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreateTag(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body TagRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	name, ok := tagInput(w, body)
	if !ok {
		return
	}

	tag, err := s.settings.SaveTag(r.Context(), nil, name)
	if err != nil {
		writeSettingsError(w, err, "Could not create tag")
		return
	}
	writeJSON(w, http.StatusCreated, tagResponse(tag))
}

func (s *Server) UpdateTag(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body TagRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	name, ok := tagInput(w, body)
	if !ok {
		return
	}

	tagID := uuid.UUID(id)
	tag, err := s.settings.SaveTag(r.Context(), &tagID, name)
	if err != nil {
		writeSettingsError(w, err, "Could not update tag")
		return
	}
	writeJSON(w, http.StatusOK, tagResponse(tag))
}

func (s *Server) DeleteTag(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	if err := s.settings.DeleteTag(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete tag")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
