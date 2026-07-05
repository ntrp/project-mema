package httpapi

import (
	"context"
	"sort"
	"strings"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func (s *Server) searchMetadataProviderPeople(
	ctx context.Context,
	provider storage.MetadataProvider,
	query string,
) ([]metadata.PersonSearchResult, error) {
	cacheKey := strings.ToLower(strings.Join(strings.Fields(query), " "))
	cached := []metadata.PersonSearchResult{}
	found, err := s.settings.GetMetadataSearchCache(ctx, provider.ID, "person", cacheKey, nil, &cached)
	if err != nil {
		return nil, err
	}
	if found {
		s.recordMetadataSearchHistory(ctx, provider, "person", cacheKey, nil, true, cached, nil)
		return cached, nil
	}

	results, err := s.metadata.SearchPeople(ctx, metadataProviderConfig(provider), query)
	if err != nil {
		s.recordMetadataSearchHistory(ctx, provider, "person", cacheKey, nil, false, nil, err)
		return nil, err
	}
	expiresAt := s.now().Add(24 * time.Hour)
	s.recordMetadataSearchHistory(ctx, provider, "person", cacheKey, nil, false, results, nil)
	if err := s.settings.SetMetadataSearchCache(ctx, provider.ID, "person", cacheKey, nil, results, expiresAt); err != nil {
		return nil, err
	}
	s.publishMetadataCacheUpdated(ctx, provider, "person", cacheKey, nil, results, expiresAt)
	return results, nil
}

func (s *Server) providerPersonSearchGroups(
	ctx context.Context,
	request groupedMediaSearchRequest,
	limit int,
) ([]MediaSearchGroup, error) {
	if request.providerIDsProvided && len(request.providerIDs) == 0 {
		return []MediaSearchGroup{}, nil
	}
	providerRequest := request
	if len(providerRequest.mediaTypes) == 0 {
		providerRequest.mediaTypes = []string{"movie", "serie"}
	}
	providers, err := s.searchableMetadataProviders(ctx, providerRequest)
	if err != nil {
		return nil, err
	}
	groups := make([]MediaSearchGroup, 0, len(providers))
	for _, provider := range providers {
		results, err := s.searchMetadataProviderPeople(ctx, provider, request.query)
		if err != nil || len(results) == 0 {
			continue
		}
		sort.SliceStable(results, func(i, j int) bool {
			return personPopularityValue(results[i]) > personPopularityValue(results[j])
		})
		providerID := openapi_types.UUID(provider.ID)
		people := make([]PersonSearchResult, 0, min(limit, len(results)))
		for _, result := range results {
			if len(people) >= limit {
				break
			}
			people = append(people, personSearchResultResponse(result))
		}
		groups = append(groups, MediaSearchGroup{
			SourceType: "provider",
			SourceName: provider.Name + " People",
			ProviderId: &providerID,
			Results:    []MediaSearchResult{},
			People:     &people,
		})
	}
	return groups, nil
}

func personSearchResultResponse(result metadata.PersonSearchResult) PersonSearchResult {
	return PersonSearchResult{
		Name:             result.Name,
		ExternalProvider: result.ExternalProvider,
		ExternalId:       result.ExternalID,
		ProfilePath:      result.ProfilePath,
		Popularity:       result.Popularity,
		KnownFor:         optionalResponseStrings(result.KnownFor),
	}
}

func personPopularityValue(result metadata.PersonSearchResult) float64 {
	if result.Popularity == nil {
		return 0
	}
	return *result.Popularity
}
