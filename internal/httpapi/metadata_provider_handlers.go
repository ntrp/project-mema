package httpapi

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func (s *Server) ListMetadataProviders(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	providers, err := s.settings.ListMetadataProviders(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_list_failed", "Could not list metadata providers")
		return
	}

	response := MetadataProviderListResponse{Providers: make([]MetadataProvider, 0, len(providers))}
	for _, provider := range providers {
		response.Providers = append(response.Providers, metadataProviderResponse(provider))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreateMetadataProvider(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body MetadataProviderRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := metadataProviderInput(w, body)
	if !ok {
		return
	}

	provider, err := s.settings.CreateMetadataProvider(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_create_failed", "Could not create metadata provider")
		return
	}
	writeJSON(w, http.StatusCreated, metadataProviderResponse(provider))
}

func (s *Server) UpdateMetadataProvider(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body MetadataProviderRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := metadataProviderInput(w, body)
	if !ok {
		return
	}

	provider, err := s.settings.UpdateMetadataProvider(r.Context(), uuid.UUID(id), input)
	if err != nil {
		writeSettingsError(w, err, "Could not update metadata provider")
		return
	}
	writeJSON(w, http.StatusOK, metadataProviderResponse(provider))
}

func (s *Server) DeleteMetadataProvider(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	if err := s.settings.DeleteMetadataProvider(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete metadata provider")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) TestMetadataProvider(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	provider, err := s.settings.GetMetadataProvider(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find metadata provider")
		return
	}

	result := s.metadata.Test(r.Context(), metadataProviderConfig(provider))
	writeJSON(w, http.StatusOK, metadataProviderTestResponse(s.now(), result))
}

func (s *Server) GetMetadataCache(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	stats, err := s.settings.MetadataCacheStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_cache_stats_failed", "Could not load metadata cache stats")
		return
	}
	entries, err := s.settings.ListMetadataCacheEntries(r.Context(), 100)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_cache_entries_failed", "Could not load metadata cache entries")
		return
	}

	response := MetadataCacheResponse{
		Stats:   metadataCacheStatsResponse(stats),
		Entries: make([]MetadataCacheEntry, 0, len(entries)),
	}
	for _, entry := range entries {
		response.Entries = append(response.Entries, metadataCacheEntryResponse(entry))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) ClearMetadataCache(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	count, err := s.settings.ClearMetadataCache(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_cache_clear_failed", "Could not clear metadata cache")
		return
	}
	writeJSON(w, http.StatusOK, MetadataCacheClearResponse{DeletedCount: count})
}

func (s *Server) ClearMetadataCacheByPattern(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body MetadataCacheClearRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	pattern := strings.TrimSpace(body.Pattern)
	if pattern == "" {
		writeError(w, http.StatusBadRequest, "invalid_pattern", "Cache reset pattern is required")
		return
	}
	if _, err := regexp.Compile(pattern); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_pattern", "Cache reset pattern is not a valid regex")
		return
	}

	count, err := s.settings.ClearMetadataCacheByPattern(r.Context(), pattern)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_cache_clear_failed", "Could not clear metadata cache")
		return
	}
	writeJSON(w, http.StatusOK, MetadataCacheClearResponse{DeletedCount: count})
}
