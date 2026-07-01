package httpapi

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) ListDownloadClients(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	clients, err := s.settings.ListDownloadClients(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_list_failed", "Could not list download clients")
		return
	}

	response := DownloadClientListResponse{Clients: make([]DownloadClient, 0, len(clients))}
	for _, client := range clients {
		response.Clients = append(response.Clients, downloadClientResponse(client))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreateDownloadClient(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body DownloadClientRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := downloadClientInput(w, body)
	if !ok {
		return
	}

	client, err := s.settings.CreateDownloadClient(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_create_failed", "Could not create download client")
		return
	}
	writeJSON(w, http.StatusCreated, downloadClientResponse(client))
}

func (s *Server) UpdateDownloadClient(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body DownloadClientRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := downloadClientInput(w, body)
	if !ok {
		return
	}

	client, err := s.settings.UpdateDownloadClient(r.Context(), uuid.UUID(id), input)
	if err != nil {
		writeSettingsError(w, err, "Could not update download client")
		return
	}
	writeJSON(w, http.StatusOK, downloadClientResponse(client))
}

func (s *Server) DeleteDownloadClient(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	if err := s.settings.DeleteDownloadClient(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete download client")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) TestDownloadClient(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	client, err := s.settings.GetDownloadClient(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find download client")
		return
	}

	result := s.downloadClients.Test(r.Context(), downloadClientConfig(client))
	writeJSON(w, http.StatusOK, downloadClientTestResponse(s.now(), result))
}

func (s *Server) TestDownloadClientConfig(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body DownloadClientRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := downloadClientInput(w, body)
	if !ok {
		return
	}

	result := s.downloadClients.Test(r.Context(), downloadClientInputConfig(input))
	writeJSON(w, http.StatusOK, downloadClientTestResponse(s.now(), result))
}
