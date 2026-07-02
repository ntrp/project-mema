package httpapi

import (
	"net/http"
	"regexp"
	"strings"

	"media-manager/internal/storage"
)

func (s *Server) GetIndexerSearch(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	response, err := s.indexerSearchResponse(r)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "indexer_search_load_failed", "Could not load indexer search cache")
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) UpdateIndexerSearchSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body IndexerSearchSettings
	if !decodeJSON(w, r, &body) {
		return
	}
	input := storage.IndexerSearchSettings{
		CacheDurationMinutes: body.CacheDurationMinutes,
		HistoryRetentionDays: body.HistoryRetentionDays,
	}
	if _, err := s.settings.SaveIndexerSearchSettings(r.Context(), input); err != nil {
		writeSettingsError(w, err, "Could not update indexer search settings")
		return
	}
	response, err := s.indexerSearchResponse(r)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "indexer_search_load_failed", "Could not load indexer search cache")
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) ClearIndexerSearchCache(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	count, err := s.settings.ClearIndexerSearchCache(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "indexer_cache_clear_failed", "Could not clear indexer search cache")
		return
	}
	writeJSON(w, http.StatusOK, MetadataCacheClearResponse{DeletedCount: count})
}

func (s *Server) ClearIndexerSearchCacheByPattern(w http.ResponseWriter, r *http.Request) {
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
	count, err := s.settings.ClearIndexerSearchCacheByPattern(r.Context(), pattern)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "indexer_cache_clear_failed", "Could not clear matching indexer search cache entries")
		return
	}
	writeJSON(w, http.StatusOK, MetadataCacheClearResponse{DeletedCount: count})
}

func (s *Server) indexerSearchResponse(r *http.Request) (IndexerSearchResponse, error) {
	settings, err := s.settings.GetIndexerSearchSettings(r.Context())
	if err != nil {
		return IndexerSearchResponse{}, err
	}
	stats, err := s.settings.IndexerSearchCacheStats(r.Context())
	if err != nil {
		return IndexerSearchResponse{}, err
	}
	cacheEntries, err := s.settings.ListIndexerSearchCacheEntries(r.Context(), 100)
	if err != nil {
		return IndexerSearchResponse{}, err
	}
	historyEntries, err := s.settings.ListIndexerSearchHistoryEntries(r.Context(), 100)
	if err != nil {
		return IndexerSearchResponse{}, err
	}
	return indexerSearchResponse(settings, stats, cacheEntries, historyEntries), nil
}
