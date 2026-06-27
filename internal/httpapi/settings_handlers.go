package httpapi

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) ListDownloadClients(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
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
	if _, ok := s.requireSession(w, r); !ok {
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
	if _, ok := s.requireSession(w, r); !ok {
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
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	if err := s.settings.DeleteDownloadClient(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete download client")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) ListIndexers(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	indexers, err := s.settings.ListIndexers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_list_failed", "Could not list indexers")
		return
	}

	response := IndexerListResponse{Indexers: make([]Indexer, 0, len(indexers))}
	for _, indexer := range indexers {
		response.Indexers = append(response.Indexers, indexerResponse(indexer))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreateIndexer(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	var body IndexerRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := indexerInput(w, body)
	if !ok {
		return
	}

	indexer, err := s.settings.CreateIndexer(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_create_failed", "Could not create indexer")
		return
	}
	writeJSON(w, http.StatusCreated, indexerResponse(indexer))
}

func (s *Server) UpdateIndexer(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	var body IndexerRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := indexerInput(w, body)
	if !ok {
		return
	}

	indexer, err := s.settings.UpdateIndexer(r.Context(), uuid.UUID(id), input)
	if err != nil {
		writeSettingsError(w, err, "Could not update indexer")
		return
	}
	writeJSON(w, http.StatusOK, indexerResponse(indexer))
}

func (s *Server) DeleteIndexer(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	if err := s.settings.DeleteIndexer(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete indexer")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
