package metadata

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (s *Service) tmdbMovieGenres(ctx context.Context, config Config) ([]tmdbIDName, error) {
	return s.tmdbGenres(ctx, config, "movie")
}

func (s *Service) tmdbSeriesGenres(ctx context.Context, config Config) ([]tmdbIDName, error) {
	return s.tmdbGenres(ctx, config, "tv")
}

func (s *Service) tmdbGenres(ctx context.Context, config Config, mediaPath string) ([]tmdbIDName, error) {
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), "genre", mediaPath, "list")
	if err != nil {
		return nil, err
	}
	var payload tmdbGenreList
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return nil, err
	}
	return payload.Genres, nil
}

func (s *Service) tmdbMovieGenreMap(ctx context.Context, config Config) (map[string]string, error) {
	genres, err := s.tmdbMovieGenres(ctx, config)
	if err != nil {
		return nil, err
	}
	return genreMap(genres), nil
}

func (s *Service) tmdbSeriesGenreMap(ctx context.Context, config Config) (map[string]string, error) {
	genres, err := s.tmdbSeriesGenres(ctx, config)
	if err != nil {
		return nil, err
	}
	return genreMap(genres), nil
}

func genreMap(genres []tmdbIDName) map[string]string {
	values := map[string]string{}
	for _, genre := range genres {
		id := strconv.FormatInt(genre.ID, 10)
		values[id] = strings.TrimSpace(genre.Name)
		values[cleanLookup(genre.Name)] = id
	}
	return values
}

func (s *Service) tmdbSearchFacet(ctx context.Context, config Config, facet string, query string) ([]FacetOption, error) {
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return []FacetOption{}, nil
	}
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), "search", facet)
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("query", query)
	endpoint = endpoint + "?" + values.Encode()
	var payload tmdbFacetSearchResponse
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return nil, err
	}
	return facetOptions(payload.Results, 10), nil
}

func genreFacetOptions(genres []tmdbIDName, query string) []FacetOption {
	query = cleanLookup(query)
	options := []FacetOption{}
	for _, genre := range genres {
		name := strings.TrimSpace(genre.Name)
		if name == "" || (query != "" && !strings.Contains(cleanLookup(name), query)) {
			continue
		}
		options = append(options, FacetOption{ID: strconv.FormatInt(genre.ID, 10), Name: name})
	}
	return options
}

func facetOptions(values []tmdbIDName, limit int) []FacetOption {
	options := []FacetOption{}
	for _, value := range values {
		name := strings.TrimSpace(value.Name)
		if name == "" {
			continue
		}
		options = append(options, FacetOption{ID: strconv.FormatInt(value.ID, 10), Name: name})
		if len(options) >= limit {
			break
		}
	}
	return options
}

func genreIDs(input []string, genres map[string]string) []string {
	ids := []string{}
	for _, value := range input {
		cleaned := strings.TrimSpace(value)
		if cleaned == "" {
			continue
		}
		if _, err := strconv.ParseInt(cleaned, 10, 64); err == nil {
			ids = append(ids, cleaned)
			continue
		}
		if id := genres[cleanLookup(cleaned)]; id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

func numericOrClean(input []string) []string {
	values := []string{}
	for _, value := range input {
		cleaned := strings.TrimSpace(value)
		if cleaned == "" {
			continue
		}
		if id, _, ok := strings.Cut(cleaned, ":"); ok {
			cleaned = id
		}
		if _, err := strconv.ParseInt(cleaned, 10, 64); err == nil {
			values = append(values, cleaned)
		}
	}
	return values
}

func cleanLookup(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(value), " "))
}
