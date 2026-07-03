package httpapi

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func (s *Server) discoverSeriesProviderResults(
	ctx context.Context,
	provider storage.MetadataProvider,
	request metadata.DiscoverSeriesRequest,
) ([]metadata.SearchResult, error) {
	cacheKey := discoverSeriesCacheKey(request)
	cached := []metadata.SearchResult{}
	found, err := s.settings.GetMetadataSearchCache(ctx, provider.ID, "series", cacheKey, nil, &cached)
	if err != nil {
		return nil, err
	}
	if found {
		s.recordMetadataSearchHistory(ctx, provider, "series", cacheKey, nil, true, cached, nil)
		return cached, nil
	}

	results, err := s.metadata.DiscoverSeries(ctx, metadataProviderConfig(provider), request)
	if err != nil {
		s.recordMetadataSearchHistory(ctx, provider, "series", cacheKey, nil, false, nil, err)
		return nil, err
	}
	expiresAt := s.now().Add(24 * time.Hour)
	s.recordMetadataSearchHistory(ctx, provider, "series", cacheKey, nil, false, results, nil)
	if err := s.settings.SetMetadataSearchCache(ctx, provider.ID, "series", cacheKey, nil, results, expiresAt); err != nil {
		return nil, err
	}
	s.publishMetadataCacheUpdated(ctx, provider, "series", cacheKey, nil, results, expiresAt)
	return results, nil
}

func discoverSeriesCacheKey(request metadata.DiscoverSeriesRequest) string {
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
	setQueryStrings(values, "status", request.Status)
	return "discover:series:" + values.Encode()
}
