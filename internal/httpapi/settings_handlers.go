package httpapi

import (
	"net/http"
	"regexp"
	"strings"

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

func (s *Server) ListQualitySizeSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	settings, err := s.settings.ListQualitySizeSettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_list_failed", "Could not list quality sizes")
		return
	}
	writeJSON(w, http.StatusOK, qualitySizeSettingsResponse(settings))
}

func (s *Server) UpdateQualitySizeSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body QualitySizeSettingsUpdateRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := qualitySizeSettingsInput(w, body)
	if !ok {
		return
	}

	settings, err := s.settings.SaveQualitySizeSettings(r.Context(), input)
	if err != nil {
		writeSettingsError(w, err, "Could not update quality sizes")
		return
	}
	writeJSON(w, http.StatusOK, qualitySizeSettingsResponse(settings))
}

func (s *Server) GetFileNamingSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	settings, err := s.settings.GetFileNamingSettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_load_failed", "Could not load file naming settings")
		return
	}
	writeJSON(w, http.StatusOK, fileNamingSettingsResponse(settings))
}

func (s *Server) UpdateFileNamingSettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body FileNamingSettingsRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := fileNamingSettingsInput(w, body)
	if !ok {
		return
	}

	settings, err := s.settings.SaveFileNamingSettings(r.Context(), input)
	if err != nil {
		writeSettingsError(w, err, "Could not update file naming settings")
		return
	}
	writeJSON(w, http.StatusOK, fileNamingSettingsResponse(settings))
}

func (s *Server) ListMediaProfiles(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	profiles, err := s.settings.ListMediaProfiles(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "settings_list_failed", "Could not list media profiles")
		return
	}
	writeJSON(w, http.StatusOK, mediaProfileListResponse(profiles))
}

func (s *Server) CreateMediaProfile(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body MediaProfileRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := mediaProfileInput(w, body)
	if !ok {
		return
	}

	profile, err := s.settings.CreateMediaProfile(r.Context(), input)
	if err != nil {
		writeSettingsError(w, err, "Could not create media profile")
		return
	}
	writeJSON(w, http.StatusCreated, mediaProfileResponse(profile))
}

func (s *Server) UpdateMediaProfile(w http.ResponseWriter, r *http.Request, id ProfileId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body MediaProfileRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := mediaProfileInput(w, body)
	if !ok {
		return
	}

	profile, err := s.settings.UpdateMediaProfile(r.Context(), string(id), input)
	if err != nil {
		writeSettingsError(w, err, "Could not update media profile")
		return
	}
	writeJSON(w, http.StatusOK, mediaProfileResponse(profile))
}

func (s *Server) DeleteMediaProfile(w http.ResponseWriter, r *http.Request, id ProfileId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	if err := s.settings.DeleteMediaProfile(r.Context(), string(id)); err != nil {
		writeSettingsError(w, err, "Could not delete media profile")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

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
