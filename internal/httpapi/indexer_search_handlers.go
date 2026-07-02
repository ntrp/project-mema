package httpapi

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func (s *Server) GetIndexerSearch(w http.ResponseWriter, r *http.Request, params GetIndexerSearchParams) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	response, err := s.indexerSearchResponse(r, params.CacheLimit, params.HistoryLimit)
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
	response, err := s.indexerSearchResponse(r, nil, nil)
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

func (s *Server) DeleteIndexerSearchCacheEntry(w http.ResponseWriter, r *http.Request, params DeleteIndexerSearchCacheEntryParams) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	if strings.TrimSpace(params.Query) == "" {
		writeError(w, http.StatusBadRequest, "invalid_cache_entry", "Cache entry query is required")
		return
	}
	count, err := s.settings.DeleteIndexerSearchCacheEntry(
		r.Context(),
		uuid.UUID(params.IndexerId),
		string(params.MediaType),
		params.Query,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "indexer_cache_delete_failed", "Could not delete indexer search cache entry")
		return
	}
	writeJSON(w, http.StatusOK, MetadataCacheClearResponse{DeletedCount: count})
}

func (s *Server) ClearIndexerSearchHistory(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	count, err := s.settings.ClearIndexerSearchHistory(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "indexer_history_clear_failed", "Could not clear indexer query history")
		return
	}
	writeJSON(w, http.StatusOK, MetadataCacheClearResponse{DeletedCount: count})
}

func (s *Server) indexerSearchResponse(r *http.Request, cacheLimit *int32, historyLimit *int32) (IndexerSearchResponse, error) {
	settings, err := s.settings.GetIndexerSearchSettings(r.Context())
	if err != nil {
		return IndexerSearchResponse{}, err
	}
	stats, err := s.settings.IndexerSearchCacheStats(r.Context())
	if err != nil {
		return IndexerSearchResponse{}, err
	}
	cacheEntries, err := s.settings.ListIndexerSearchCacheEntries(r.Context(), optionalLimit(cacheLimit))
	if err != nil {
		return IndexerSearchResponse{}, err
	}
	historyEntries, err := s.settings.ListIndexerSearchHistoryEntries(r.Context(), optionalLimit(historyLimit))
	if err != nil {
		return IndexerSearchResponse{}, err
	}
	historyStats, err := s.settings.IndexerSearchHistoryStats(r.Context())
	if err != nil {
		return IndexerSearchResponse{}, err
	}
	return indexerSearchResponse(settings, stats, cacheEntries, historyEntries, historyStats), nil
}

func optionalLimit(limit *int32) int32 {
	if limit == nil {
		return 0
	}
	return *limit
}
