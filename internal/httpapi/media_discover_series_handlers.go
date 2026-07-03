package httpapi

import (
	"net/http"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func (s *Server) SearchDiscoverSeries(w http.ResponseWriter, r *http.Request, params SearchDiscoverSeriesParams) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	provider, ok := s.discoverSeriesProvider(w, r)
	if !ok {
		return
	}
	if provider.ID.String() == "00000000-0000-0000-0000-000000000000" {
		writeJSON(w, http.StatusOK, DiscoverMovieSearchResponse{Results: []MediaSearchResult{}, HasMore: false})
		return
	}
	results, err := s.discoverSeriesProviderResults(r.Context(), provider, discoverSeriesRequest(params))
	if err != nil {
		writeMetadataDetailsError(w, err)
		return
	}
	response := DiscoverMovieSearchResponse{Results: make([]MediaSearchResult, 0, len(results)), HasMore: len(results) >= 20}
	for _, result := range results {
		response.Results = append(response.Results, metadataSearchResultResponse(result))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) AutocompleteDiscoverSeriesFacet(w http.ResponseWriter, r *http.Request, facet AutocompleteDiscoverSeriesFacetParamsFacet, params AutocompleteDiscoverSeriesFacetParams) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	provider, ok := s.discoverSeriesProvider(w, r)
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
	options, err := s.metadata.DiscoverSeriesFacet(r.Context(), metadataProviderConfig(provider), string(facet), query)
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

func (s *Server) discoverSeriesProvider(w http.ResponseWriter, r *http.Request) (storage.MetadataProvider, bool) {
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

func discoverSeriesRequest(params SearchDiscoverSeriesParams) metadata.DiscoverSeriesRequest {
	return metadata.DiscoverSeriesRequest{
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
		Status:            stringListValue(params.Status),
		RuntimeMin:        params.RuntimeMin,
		RuntimeMax:        params.RuntimeMax,
		ScoreMin:          params.ScoreMin,
		ScoreMax:          params.ScoreMax,
		MinVoteCount:      params.MinVoteCount,
	}
}
