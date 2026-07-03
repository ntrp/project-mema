package metadata

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (s *Service) DiscoverMovies(ctx context.Context, config Config, request DiscoverMovieRequest) ([]SearchResult, error) {
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), "discover", "movie")
	if err != nil {
		return nil, err
	}
	genres := map[string]string{}
	if len(request.Genres) > 0 || len(request.WithoutGenres) > 0 {
		var err error
		genres, err = s.tmdbMovieGenreMap(ctx, config)
		if err != nil {
			return nil, err
		}
	}
	request.Studios = s.resolveFacetIDs(ctx, config, "company", request.Studios)
	request.Keywords = s.resolveFacetIDs(ctx, config, "keyword", request.Keywords)
	request.WithoutKeywords = s.resolveFacetIDs(ctx, config, "keyword", request.WithoutKeywords)
	values := tmdbDiscoverMovieValues(request, genres)
	endpoint = endpoint + "?" + values.Encode()

	var payload tmdbSearchResponse
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return nil, err
	}
	return tmdbDiscoverMovieResults(payload.Results, genres), nil
}

func (s *Service) resolveFacetIDs(ctx context.Context, config Config, facet string, values []string) []string {
	resolved := []string{}
	for _, value := range values {
		cleaned := strings.TrimSpace(value)
		if cleaned == "" {
			continue
		}
		if id, _, ok := strings.Cut(cleaned, ":"); ok {
			cleaned = id
		}
		if _, err := strconv.ParseInt(cleaned, 10, 64); err == nil {
			resolved = append(resolved, cleaned)
			continue
		}
		options, err := s.tmdbSearchFacet(ctx, config, facet, cleaned)
		if err == nil && len(options) > 0 {
			resolved = append(resolved, options[0].ID)
		}
	}
	return resolved
}

func (s *Service) DiscoverMovieFacet(ctx context.Context, config Config, facet string, query string) ([]FacetOption, error) {
	switch facet {
	case "genres":
		genres, err := s.tmdbMovieGenres(ctx, config)
		if err != nil {
			return nil, err
		}
		return genreFacetOptions(genres, query), nil
	case "studios":
		return s.tmdbSearchFacet(ctx, config, "company", query)
	case "keywords":
		return s.tmdbSearchFacet(ctx, config, "keyword", query)
	default:
		return []FacetOption{}, nil
	}
}

func tmdbDiscoverMovieValues(request DiscoverMovieRequest, genres map[string]string) url.Values {
	values := url.Values{}
	values.Set("include_adult", "false")
	values.Set("page", strconv.Itoa(discoverPage(request.Page)))
	values.Set("sort_by", tmdbMovieSort(request.Sort))
	setOptional(values, "primary_release_date.gte", request.ReleaseDateFrom)
	setOptional(values, "primary_release_date.lte", request.ReleaseDateTo)
	setInt(values, "with_runtime.gte", request.RuntimeMin)
	setInt(values, "with_runtime.lte", request.RuntimeMax)
	setFloat(values, "vote_average.gte", request.ScoreMin)
	setFloat(values, "vote_average.lte", request.ScoreMax)
	setInt(values, "vote_count.gte", request.MinVoteCount)
	setJoined(values, "with_companies", request.Studios, "|")
	setJoined(values, "with_keywords", request.Keywords, "|")
	setJoined(values, "without_keywords", request.WithoutKeywords, ",")
	setJoined(values, "certification", request.ContentRatings, "|")
	if len(request.ContentRatings) > 0 {
		values.Set("certification_country", "US")
	}
	if ids := genreIDs(request.Genres, genres); len(ids) > 0 {
		values.Set("with_genres", strings.Join(ids, "|"))
	}
	if ids := genreIDs(request.WithoutGenres, genres); len(ids) > 0 {
		values.Set("without_genres", strings.Join(ids, ","))
	}
	if len(request.OriginalLanguages) > 0 {
		values.Set("with_original_language", request.OriginalLanguages[0])
	}
	return values
}

func tmdbMovieSort(value string) string {
	switch value {
	case "release_date.desc":
		return "primary_release_date.desc"
	case "release_date.asc":
		return "primary_release_date.asc"
	case "vote_average.desc", "vote_average.asc", "title.asc", "title.desc", "popularity.asc":
		return value
	default:
		return "popularity.desc"
	}
}

func setOptional(values url.Values, key string, value *string) {
	if value != nil && strings.TrimSpace(*value) != "" {
		values.Set(key, strings.TrimSpace(*value))
	}
}

func setInt(values url.Values, key string, value *int32) {
	if value != nil {
		values.Set(key, strconv.Itoa(int(*value)))
	}
}

func setFloat(values url.Values, key string, value *float64) {
	if value != nil {
		values.Set(key, strconv.FormatFloat(*value, 'f', 1, 64))
	}
}

func setJoined(values url.Values, key string, input []string, separator string) {
	if ids := numericOrClean(input); len(ids) > 0 {
		values.Set(key, strings.Join(ids, separator))
	}
}
