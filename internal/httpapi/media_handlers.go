package httpapi

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/metadata"
)

func (s *Server) SearchMedia(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	var body MediaSearchRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	query := strings.TrimSpace(body.Query)
	if query == "" {
		writeError(w, http.StatusBadRequest, "invalid_query", "Search query is required")
		return
	}
	if !body.Type.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_type", "Media type is not supported")
		return
	}

	providers, err := s.settings.ListEnabledMetadataProviders(r.Context(), string(body.Type))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_provider_list_failed", "Could not list metadata providers")
		return
	}
	if len(providers) == 0 {
		writeJSON(w, http.StatusOK, MediaSearchResponse{
			Results: []MediaSearchResult{
				{
					Title: query,
					Type:  body.Type,
					Year:  body.Year,
				},
			},
		})
		return
	}

	response := MediaSearchResponse{Results: []MediaSearchResult{}}
	for _, provider := range providers {
		results, err := s.searchMetadataProvider(r.Context(), provider, metadata.SearchRequest{
			Query:     query,
			MediaType: string(body.Type),
			Year:      body.Year,
		})
		if err != nil {
			continue
		}
		for _, result := range results {
			response.Results = append(response.Results, metadataSearchResultResponse(result))
		}
		if len(response.Results) > 0 {
			break
		}
	}
	if len(response.Results) == 0 {
		response.Results = append(response.Results, MediaSearchResult{
			Title: query,
			Type:  body.Type,
			Year:  body.Year,
		})
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GetMediaDiscover(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	providers, err := s.settings.ListMetadataProviders(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_provider_list_failed", "Could not list metadata providers")
		return
	}
	blacklist, err := s.settings.ListDiscoverBlacklist(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "discover_blacklist_list_failed", "Could not list discover blacklist")
		return
	}

	response := MediaDiscoverResponse{Sections: make([]MediaDiscoverSection, 0, len(discoverSections))}
	for _, section := range discoverSections {
		response.Sections = append(response.Sections, s.discoverSectionResponse(r.Context(), providers, section, 20, 1, blacklist))
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GetMediaDiscoverSection(w http.ResponseWriter, r *http.Request, sectionId string, params GetMediaDiscoverSectionParams) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	section, ok := discoverSectionByID(sectionId)
	if !ok {
		writeError(w, http.StatusNotFound, "discover_section_not_found", "Discovery section was not found")
		return
	}
	providers, err := s.settings.ListMetadataProviders(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_provider_list_failed", "Could not list metadata providers")
		return
	}
	blacklist, err := s.settings.ListDiscoverBlacklist(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "discover_blacklist_list_failed", "Could not list discover blacklist")
		return
	}
	page := int32(1)
	if params.Page != nil {
		page = *params.Page
	}
	limit := int32(20)
	if params.Limit != nil {
		limit = *params.Limit
	}
	writeJSON(w, http.StatusOK, s.discoverSectionResponse(r.Context(), providers, section, int(limit), int(page), blacklist))
}

func (s *Server) AutocompleteMedia(w http.ResponseWriter, r *http.Request, params AutocompleteMediaParams) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	query := strings.TrimSpace(params.Query)
	if len(query) < 2 {
		writeError(w, http.StatusBadRequest, "invalid_query", "Search query must contain at least 2 characters")
		return
	}

	groups, err := s.groupedMediaSearch(r.Context(), groupedMediaSearchRequest{
		query:            query,
		mediaTypes:       []string{"movie", "series"},
		limit:            5,
		includeLibrary:   boolDefault(params.IncludeLibrary, true),
		includeProviders: boolDefault(params.IncludeProviders, true),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "media_autocomplete_failed", "Could not search media")
		return
	}
	writeJSON(w, http.StatusOK, MediaGroupedSearchResponse{Groups: groups})
}

func (s *Server) AdvancedSearchMedia(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	var body MediaAdvancedSearchRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	query := strings.TrimSpace(valueOrEmpty(body.Query))
	if query == "" {
		writeError(w, http.StatusBadRequest, "invalid_query", "Search query is required")
		return
	}
	mediaTypes := []string{"movie", "series"}
	if body.Type != nil {
		if !body.Type.Valid() {
			writeError(w, http.StatusBadRequest, "invalid_type", "Media type is not supported")
			return
		}
		mediaTypes = []string{string(*body.Type)}
	}
	limit := int32(20)
	if body.Limit != nil {
		limit = *body.Limit
	}

	providerIDs := map[uuid.UUID]struct{}{}
	if body.ProviderIds != nil {
		for _, id := range *body.ProviderIds {
			providerIDs[uuid.UUID(id)] = struct{}{}
		}
	}

	groups, err := s.groupedMediaSearch(r.Context(), groupedMediaSearchRequest{
		query:               query,
		mediaTypes:          mediaTypes,
		year:                body.Year,
		providerIDs:         providerIDs,
		providerIDsProvided: body.ProviderIds != nil,
		limit:               int(limit),
		includeLibrary:      true,
		includeProviders:    true,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "media_advanced_search_failed", "Could not search media")
		return
	}
	writeJSON(w, http.StatusOK, MediaGroupedSearchResponse{Groups: groups})
}

func (s *Server) GetMediaMetadataDetails(w http.ResponseWriter, r *http.Request, providerType MetadataProviderType, mediaType MediaType, externalID string) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	if !providerType.Valid() || !mediaType.Valid() || strings.TrimSpace(externalID) == "" {
		writeError(w, http.StatusBadRequest, "invalid_metadata_request", "Metadata provider, media type, and external id are required")
		return
	}

	providers, err := s.settings.ListMetadataProviders(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_provider_list_failed", "Could not list metadata providers")
		return
	}
	provider, ok := metadataProviderByType(providers, string(providerType))
	if !ok {
		writeError(w, http.StatusNotFound, "metadata_provider_not_found", "Metadata provider is not configured")
		return
	}

	details, err := s.metadataProviderDetails(r.Context(), provider, metadata.DetailsRequest{
		MediaType:  string(mediaType),
		ExternalID: externalID,
	})
	if err != nil {
		writeMetadataDetailsError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, metadataDetailsResponse(details))
}
