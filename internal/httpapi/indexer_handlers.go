package httpapi

import (
	"net/http"

	"github.com/google/uuid"

	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

func (s *Server) ListIndexers(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
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

func (s *Server) GetIndexer(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	indexer, err := s.settings.GetIndexer(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find indexer")
		return
	}
	writeJSON(w, http.StatusOK, indexerResponse(indexer))
}

func (s *Server) ListIndexerCatalog(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	writeJSON(w, http.StatusOK, indexerCatalogResponse(indexers.Catalog()))
}

func (s *Server) GetIndexerCatalogDefinition(w http.ResponseWriter, r *http.Request, definitionId string) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	entry, ok := indexers.CatalogEntryByID(definitionId)
	if !ok {
		writeError(w, http.StatusNotFound, "indexer_definition_not_found", "Indexer definition not found")
		return
	}
	writeJSON(w, http.StatusOK, catalogEntryResponse(entry))
}

func (s *Server) CreateIndexer(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
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
	if _, ok := s.requireAdmin(w, r); !ok {
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

func (s *Server) BulkUpdateIndexers(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body IndexerBulkUpdateRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	ids := make([]uuid.UUID, 0, len(body.Ids))
	for _, id := range body.Ids {
		ids = append(ids, uuid.UUID(id))
	}
	updated, err := s.settings.BulkUpdateIndexers(r.Context(), storage.IndexerBulkUpdateInput{
		IDs:             ids,
		Enabled:         body.Enabled,
		AppProfileID:    body.AppProfileId,
		Priority:        body.Priority,
		MinimumSeeders:  body.MinimumSeeders,
		SeedRatio:       body.SeedRatio,
		SeedTime:        body.SeedTime,
		PackSeedTime:    body.PackSeedTime,
		PreferMagnetURL: body.PreferMagnetUrl,
	})
	if err != nil {
		writeSettingsError(w, err, "Could not update indexers")
		return
	}
	response := IndexerListResponse{Indexers: make([]Indexer, 0, len(updated))}
	for _, indexer := range updated {
		response.Indexers = append(response.Indexers, indexerResponse(indexer))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) DeleteIndexer(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	if err := s.settings.DeleteIndexer(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete indexer")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) TestIndexer(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	indexer, err := s.settings.GetIndexer(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find indexer")
		return
	}

	result := s.indexers.Test(r.Context(), indexerConfig(indexer))
	s.recordIndexerTestResult(r.Context(), indexer, result)
	writeJSON(w, http.StatusOK, indexerTestResponse(s.now(), result))
}
