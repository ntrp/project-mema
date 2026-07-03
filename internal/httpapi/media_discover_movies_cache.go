package httpapi

import (
	"context"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func (s *Server) discoverMoviesProvider(
	ctx context.Context,
	provider storage.MetadataProvider,
	request metadata.DiscoverMovieRequest,
) ([]metadata.SearchResult, error) {
	cacheKey := discoverMoviesCacheKey(request)
	cached := []metadata.SearchResult{}
	found, err := s.settings.GetMetadataSearchCache(ctx, provider.ID, "movie", cacheKey, nil, &cached)
	if err != nil {
		return nil, err
	}
	if found {
		s.recordMetadataSearchHistory(ctx, provider, "movie", cacheKey, nil, true, cached, nil)
		return cached, nil
	}

	results, err := s.metadata.DiscoverMovies(ctx, metadataProviderConfig(provider), request)
	if err != nil {
		s.recordMetadataSearchHistory(ctx, provider, "movie", cacheKey, nil, false, nil, err)
		return nil, err
	}
	expiresAt := s.now().Add(24 * time.Hour)
	s.recordMetadataSearchHistory(ctx, provider, "movie", cacheKey, nil, false, results, nil)
	if err := s.settings.SetMetadataSearchCache(ctx, provider.ID, "movie", cacheKey, nil, results, expiresAt); err != nil {
		return nil, err
	}
	s.publishMetadataCacheUpdated(ctx, provider, "movie", cacheKey, nil, results, expiresAt)
	return results, nil
}

func discoverMoviesCacheKey(request metadata.DiscoverMovieRequest) string {
	values := url.Values{}
	values.Set("sort", valueOrDefaultString(request.Sort, "popularity.desc"))
	values.Set("page", strconv.Itoa(discoverMoviePage(request.Page)))
	setQueryString(values, "releaseDateFrom", request.ReleaseDateFrom)
	setQueryString(values, "releaseDateTo", request.ReleaseDateTo)
	setQueryInt(values, "runtimeMin", request.RuntimeMin)
	setQueryInt(values, "runtimeMax", request.RuntimeMax)
	setQueryFloat(values, "scoreMin", request.ScoreMin)
	setQueryFloat(values, "scoreMax", request.ScoreMax)
	setQueryInt(values, "minVoteCount", request.MinVoteCount)
	setQueryStrings(values, "studios", request.Studios)
	setQueryStrings(values, "genres", request.Genres)
	setQueryStrings(values, "keywords", request.Keywords)
	setQueryStrings(values, "withoutGenres", request.WithoutGenres)
	setQueryStrings(values, "withoutKeywords", request.WithoutKeywords)
	setQueryStrings(values, "originalLanguages", request.OriginalLanguages)
	setQueryStrings(values, "contentRatings", request.ContentRatings)
	return "discover:movies:" + values.Encode()
}

func discoverMoviePage(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

func valueOrDefaultString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func setQueryString(values url.Values, key string, value *string) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return
	}
	values.Set(key, strings.TrimSpace(*value))
}

func setQueryInt(values url.Values, key string, value *int32) {
	if value != nil {
		values.Set(key, strconv.Itoa(int(*value)))
	}
}

func setQueryFloat(values url.Values, key string, value *float64) {
	if value != nil {
		values.Set(key, strconv.FormatFloat(*value, 'f', 1, 64))
	}
}

func setQueryStrings(values url.Values, key string, input []string) {
	items := cleanedSortedStrings(input)
	for _, item := range items {
		values.Add(key, item)
	}
}

func cleanedSortedStrings(input []string) []string {
	items := []string{}
	for _, value := range input {
		if cleaned := strings.TrimSpace(value); cleaned != "" {
			items = append(items, cleaned)
		}
	}
	slices.Sort(items)
	return items
}
