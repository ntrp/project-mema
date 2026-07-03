package httpapi

import (
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func (s *Server) SearchDiscoverMovies(w http.ResponseWriter, r *http.Request, params SearchDiscoverMoviesParams) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	provider, ok := s.discoverMovieProvider(w, r)
	if !ok {
		return
	}
	if provider.ID.String() == "00000000-0000-0000-0000-000000000000" {
		writeJSON(w, http.StatusOK, DiscoverMovieSearchResponse{Results: []MediaSearchResult{}, HasMore: false})
		return
	}
	results, err := s.discoverMoviesProvider(r.Context(), provider, discoverMovieRequest(params))
	if err != nil {
		writeMetadataDetailsError(w, err)
		return
	}
	response := DiscoverMovieSearchResponse{
		Results: make([]MediaSearchResult, 0, len(results)),
		HasMore: len(results) >= 20,
	}
	for _, result := range results {
		response.Results = append(response.Results, metadataSearchResultResponse(result))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) AutocompleteDiscoverMovieFacet(w http.ResponseWriter, r *http.Request, facet AutocompleteDiscoverMovieFacetParamsFacet, params AutocompleteDiscoverMovieFacetParams) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	provider, ok := s.discoverMovieProvider(w, r)
	if !ok {
		return
	}
	if provider.ID.String() == "00000000-0000-0000-0000-000000000000" {
		writeJSON(w, http.StatusOK, DiscoverMovieFacetResponse{Options: []DiscoverMovieFacetOption{}})
		return
	}
	query := ""
	if params.Query != nil {
		query = *params.Query
	}
	options, err := s.metadata.DiscoverMovieFacet(r.Context(), metadataProviderConfig(provider), string(facet), query)
	if err != nil {
		writeMetadataDetailsError(w, err)
		return
	}
	response := DiscoverMovieFacetResponse{Options: make([]DiscoverMovieFacetOption, 0, len(options))}
	for _, option := range options {
		response.Options = append(response.Options, DiscoverMovieFacetOption{Id: option.ID, Name: option.Name})
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) discoverMovieProvider(w http.ResponseWriter, r *http.Request) (storage.MetadataProvider, bool) {
	providers, err := s.settings.ListMetadataProviders(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_provider_list_failed", "Could not list metadata providers")
		return storage.MetadataProvider{}, false
	}
	provider, ok := discoverProvider(providers)
	if !ok {
		return storage.MetadataProvider{}, true
	}
	return provider, true
}

func discoverMovieRequest(params SearchDiscoverMoviesParams) metadata.DiscoverMovieRequest {
	return metadata.DiscoverMovieRequest{
		Sort:              valueOrDefault(params.Sort, "popularity.desc"),
		Page:              int(int32Value(params.Page, 1)),
		ReleaseDateFrom:   openAPIDateString(params.ReleaseDateFrom),
		ReleaseDateTo:     openAPIDateString(params.ReleaseDateTo),
		Studios:           stringListValue(params.Studios),
		Genres:            stringListValue(params.Genres),
		Keywords:          stringListValue(params.Keywords),
		WithoutGenres:     stringListValue(params.WithoutGenres),
		WithoutKeywords:   stringListValue(params.WithoutKeywords),
		OriginalLanguages: stringListValue(params.OriginalLanguages),
		ContentRatings:    stringListValue(params.ContentRatings),
		RuntimeMin:        params.RuntimeMin,
		RuntimeMax:        params.RuntimeMax,
		ScoreMin:          params.ScoreMin,
		ScoreMax:          params.ScoreMax,
		MinVoteCount:      params.MinVoteCount,
	}
}

func valueOrDefault(value *string, fallback string) string {
	if value == nil || *value == "" {
		return fallback
	}
	return *value
}

func int32Value(value *int32, fallback int32) int32 {
	if value == nil {
		return fallback
	}
	return *value
}

func stringListValue(values *[]string) []string {
	if values == nil {
		return nil
	}
	return append([]string(nil), (*values)...)
}

func openAPIDateString(value *openapi_types.Date) *string {
	if value == nil {
		return nil
	}
	formatted := value.Time.Format("2006-01-02")
	return &formatted
}
