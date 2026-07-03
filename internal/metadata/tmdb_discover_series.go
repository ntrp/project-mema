package metadata

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (s *Service) DiscoverSeries(ctx context.Context, config Config, request DiscoverSeriesRequest) ([]SearchResult, error) {
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), "discover", "tv")
	if err != nil {
		return nil, err
	}
	genres := map[string]string{}
	if len(request.Genres) > 0 || len(request.WithoutGenres) > 0 {
		genres, err = s.tmdbSeriesGenreMap(ctx, config)
		if err != nil {
			return nil, err
		}
	}
	request.Studios = s.resolveFacetIDs(ctx, config, "company", request.Studios)
	request.Keywords = s.resolveFacetIDs(ctx, config, "keyword", request.Keywords)
	request.WithoutKeywords = s.resolveFacetIDs(ctx, config, "keyword", request.WithoutKeywords)
	values := tmdbDiscoverSeriesValues(request, genres)
	endpoint = endpoint + "?" + values.Encode()

	var payload tmdbSearchResponse
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &payload); err != nil {
		return nil, err
	}
	return tmdbDiscoverSeriesResults(payload.Results, genres), nil
}

func (s *Service) DiscoverSeriesFacet(ctx context.Context, config Config, facet string, query string) ([]FacetOption, error) {
	switch facet {
	case "genres":
		genres, err := s.tmdbSeriesGenres(ctx, config)
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

func tmdbDiscoverSeriesValues(request DiscoverSeriesRequest, genres map[string]string) url.Values {
	values := url.Values{}
	values.Set("include_adult", "false")
	values.Set("page", strconv.Itoa(discoverPage(request.Page)))
	values.Set("sort_by", tmdbSeriesSort(request.Sort))
	setOptional(values, "first_air_date.gte", request.ReleaseDateFrom)
	setOptional(values, "first_air_date.lte", request.ReleaseDateTo)
	setInt(values, "with_runtime.gte", request.RuntimeMin)
	setInt(values, "with_runtime.lte", request.RuntimeMax)
	setFloat(values, "vote_average.gte", request.ScoreMin)
	setFloat(values, "vote_average.lte", request.ScoreMax)
	setInt(values, "vote_count.gte", request.MinVoteCount)
	setJoined(values, "with_companies", request.Studios, "|")
	setJoined(values, "with_keywords", request.Keywords, "|")
	setJoined(values, "without_keywords", request.WithoutKeywords, ",")
	setJoined(values, "certification", request.ContentRatings, "|")
	setJoined(values, "with_status", seriesStatusIDs(request.Status), "|")
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

func tmdbSeriesSort(value string) string {
	switch value {
	case "first_air_date.desc", "first_air_date.asc":
		return value
	case "vote_average.desc", "vote_average.asc", "name.asc", "name.desc", "popularity.asc":
		return value
	default:
		return "popularity.desc"
	}
}
