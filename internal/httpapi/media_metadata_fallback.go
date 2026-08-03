package httpapi

import (
	"context"
	"strings"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

type metadataResourceLoader func(
	context.Context,
	storage.MetadataProvider,
	metadata.DetailsRequest,
) (metadata.MediaResource, error)

type metadataResourceSearcher func(
	context.Context,
	storage.MetadataProvider,
	metadata.SearchRequest,
) ([]metadata.SearchResult, error)

func (s *Server) resourceWithMetadataFallback(
	ctx context.Context,
	providers []storage.MetadataProvider,
	primary metadata.MediaResource,
	load metadataResourceLoader,
	search metadataResourceSearcher,
) metadata.MediaResource {
	if primary.ExternalProvider == "tmdb" {
		return primary
	}
	provider, ok := metadataFallbackProvider(providers)
	if !ok {
		return primary
	}
	externalID, ok := providerMappingID(primary.ProviderMappings, "tmdb")
	if !ok {
		results, err := search(ctx, provider, metadata.SearchRequest{
			Query: primary.Title, MediaType: primary.Type, Year: primary.Year,
		})
		if err != nil {
			return primary
		}
		externalID, ok = fallbackSearchResultID(results, primary)
		if !ok {
			return primary
		}
	}
	fallback, err := load(ctx, provider, metadata.DetailsRequest{
		MediaType: primary.Type, ExternalID: externalID,
	})
	if err != nil {
		return primary
	}
	return metadata.MergeMediaResources(primary, fallback)
}

func metadataFallbackProvider(providers []storage.MetadataProvider) (storage.MetadataProvider, bool) {
	for _, provider := range providers {
		if provider.Type == "tmdb" && provider.Enabled && provider.FallbackEnabled && providerHasCredentials(provider) {
			return provider, true
		}
	}
	return storage.MetadataProvider{}, false
}

func fallbackSearchResultID(results []metadata.SearchResult, primary metadata.MediaResource) (string, bool) {
	for _, result := range results {
		if !strings.EqualFold(strings.TrimSpace(result.Title), strings.TrimSpace(primary.Title)) ||
			(primary.Year != nil && result.Year != nil && *primary.Year != *result.Year) {
			continue
		}
		if id := strings.TrimSpace(result.ExternalID); id != "" {
			return id, true
		}
	}
	return "", false
}

func providerMappingID(mappings []metadata.ProviderMapping, provider string) (string, bool) {
	for _, mapping := range mappings {
		if mapping.ProviderName == provider && strings.TrimSpace(mapping.ExternalID) != "" {
			return strings.TrimSpace(mapping.ExternalID), true
		}
	}
	return "", false
}
