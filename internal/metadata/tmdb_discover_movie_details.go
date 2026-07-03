package metadata

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func tmdbDiscoverMovieResults(items []tmdbMedia, genres map[string]string) []SearchResult {
	results := make([]SearchResult, 0, min(len(items), 20))
	for _, item := range items {
		result := tmdbDiscoverMovieResult(item, genres)
		if result.Title == "" {
			continue
		}
		results = append(results, result)
		if len(results) >= 20 {
			break
		}
	}
	return results
}

func tmdbDiscoverMovieResult(item tmdbMedia, genres map[string]string) SearchResult {
	return SearchResult{
		Title:            strings.TrimSpace(item.Title),
		Type:             "movie",
		Year:             yearFromDate(item.ReleaseDate),
		ExternalProvider: "tmdb",
		ExternalID:       strconv.FormatInt(item.ID, 10),
		Overview:         optionalString(item.Overview),
		PosterPath:       optionalString(item.PosterPath),
		BackdropPath:     optionalString(item.BackdropPath),
		Popularity:       optionalFloat64(item.Popularity),
		ReleaseDate:      optionalString(item.ReleaseDate),
		VoteAverage:      optionalFloat64(item.VoteAverage),
		VoteCount:        optionalInt32(item.VoteCount),
		OriginalLanguage: optionalString(item.Language),
		Genres:           genreNames(item.GenreIDs, genres),
	}
}

func (s *Service) tmdbMovieDiscoverDetails(ctx context.Context, config Config, id int64) (tmdbDetails, error) {
	endpoint, err := url.JoinPath(strings.TrimRight(config.BaseURL, "/"), "movie", strconv.FormatInt(id, 10))
	if err != nil {
		return tmdbDetails{}, err
	}
	values := url.Values{}
	values.Set("append_to_response", "keywords,release_dates")
	endpoint = endpoint + "?" + values.Encode()
	var details tmdbDetails
	if err := s.doJSON(ctx, config, http.MethodGet, endpoint, nil, &details); err != nil {
		return tmdbDetails{}, err
	}
	return details, nil
}

func mergeMovieDiscoverDetails(result *SearchResult, details tmdbDetails) {
	if details.Runtime > 0 {
		result.RuntimeMinutes = &details.Runtime
	}
	if len(details.Production) > 0 {
		result.Studios = tmdbNames(details.Production)
	}
	if len(details.Keywords.Keywords) > 0 || len(details.Keywords.Results) > 0 {
		result.Keywords = tmdbKeywordNames(details.Keywords)
	}
	if rating := tmdbReleaseCertification(details.ReleaseDates); rating != "" {
		result.ContentRating = &rating
	}
}

func genreNames(ids []int64, genres map[string]string) []string {
	names := []string{}
	for _, id := range ids {
		if name := genres[strconv.FormatInt(id, 10)]; name != "" {
			names = append(names, name)
		}
	}
	return names
}

func optionalInt32(value int32) *int32 {
	if value == 0 {
		return nil
	}
	return &value
}
