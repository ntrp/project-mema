package httpapi

import (
	"context"
	"strings"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func (s *Server) metadataProviderDetails(ctx context.Context, provider storage.MetadataProvider, request metadata.DetailsRequest) (metadata.Details, error) {
	details, err := s.cachedMetadataProviderDetails(ctx, provider, request)
	if err != nil || provider.Type != "anidb" {
		return details, err
	}
	return enrichAniDBRelatedMedia(ctx, provider, details, s.cachedMetadataProviderDetails), nil
}

func enrichAniDBRelatedMedia(
	ctx context.Context,
	provider storage.MetadataProvider,
	resource metadata.MediaResource,
	load metadataResourceLoader,
) metadata.MediaResource {
	loaded := map[string]metadata.MediaResource{}
	attempted := map[string]bool{}
	enrich := func(results []metadata.SearchResult) []metadata.SearchResult {
		for index, result := range results {
			externalID := strings.TrimSpace(result.ExternalID)
			if externalID == "" {
				continue
			}
			details, ok := loaded[externalID]
			if !ok && !attempted[externalID] {
				attempted[externalID] = true
				var err error
				details, err = load(ctx, provider, metadata.DetailsRequest{
					MediaType: result.Type, ExternalID: externalID,
				})
				if err == nil {
					loaded[externalID] = details
					ok = true
				}
			}
			if ok {
				results[index] = enrichedAniDBRelatedResult(result, details)
			}
		}
		return results
	}
	resource.Recommendations = enrich(resource.Recommendations)
	resource.Similar = enrich(resource.Similar)
	return resource
}

func enrichedAniDBRelatedResult(result metadata.SearchResult, details metadata.MediaResource) metadata.SearchResult {
	if strings.TrimSpace(details.Title) != "" {
		result.Title = details.Title
	}
	if strings.TrimSpace(details.Type) != "" {
		result.Type = details.Type
	}
	result.Year = details.Year
	result.ExternalURL = details.ExternalURL
	result.Overview = details.Overview
	result.PosterPath = details.PosterPath
	result.ReleaseDate = details.ReleaseDate
	result.RuntimeMinutes = details.RuntimeMinutes
	result.VoteAverage = details.VoteAverage
	result.OriginalLanguage = details.OriginalLanguage
	result.Genres = append([]string(nil), details.Genres...)
	result.Keywords = append([]string(nil), details.Keywords...)
	result.BackdropPath = details.BackdropPath
	return result
}
