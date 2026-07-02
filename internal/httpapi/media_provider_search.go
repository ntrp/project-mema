package httpapi

import (
	"context"
	"sort"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func (s *Server) providerMediaSearchGroups(
	ctx context.Context,
	request groupedMediaSearchRequest,
	limit int,
) ([]MediaSearchGroup, error) {
	if request.providerIDsProvided && len(request.providerIDs) == 0 {
		return []MediaSearchGroup{}, nil
	}
	providers, err := s.searchableMetadataProviders(ctx, request)
	if err != nil {
		return nil, err
	}
	groups := make([]MediaSearchGroup, 0, len(providers))
	for _, provider := range providers {
		results := s.searchProviderMediaTypes(ctx, provider, request)
		if len(results) == 0 {
			continue
		}
		if request.sortByPopularity {
			sort.SliceStable(results, func(i, j int) bool {
				return popularityValue(results[i]) > popularityValue(results[j])
			})
		}
		providerID := openapi_types.UUID(provider.ID)
		group := MediaSearchGroup{
			SourceType: "provider",
			SourceName: provider.Name,
			ProviderId: &providerID,
			Results:    make([]MediaSearchResult, 0, min(limit, len(results))),
		}
		for _, result := range results {
			if len(group.Results) >= limit {
				break
			}
			group.Results = append(group.Results, metadataSearchResultResponse(result))
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (s *Server) searchableMetadataProviders(
	ctx context.Context,
	request groupedMediaSearchRequest,
) ([]storage.MetadataProvider, error) {
	seen := map[string]struct{}{}
	providers := []storage.MetadataProvider{}
	for _, mediaType := range request.mediaTypes {
		found, err := s.settings.ListEnabledMetadataProviders(ctx, mediaType)
		if err != nil {
			return nil, err
		}
		for _, provider := range found {
			if len(request.providerIDs) > 0 {
				if _, ok := request.providerIDs[provider.ID]; !ok {
					continue
				}
			}
			key := provider.ID.String()
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			providers = append(providers, provider)
		}
	}
	return providers, nil
}

func (s *Server) searchProviderMediaTypes(
	ctx context.Context,
	provider storage.MetadataProvider,
	request groupedMediaSearchRequest,
) []metadata.SearchResult {
	results := []metadata.SearchResult{}
	for _, mediaType := range request.mediaTypes {
		found, err := s.searchMetadataProvider(ctx, provider, metadata.SearchRequest{
			Query:     request.query,
			MediaType: mediaType,
			Year:      request.year,
		})
		if err != nil {
			continue
		}
		results = append(results, found...)
	}
	return results
}
