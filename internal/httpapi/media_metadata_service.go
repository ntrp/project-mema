package httpapi

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func writeMetadataDetailsError(w http.ResponseWriter, err error) {
	if errors.Is(err, metadata.ErrCredentialsRequired) {
		writeError(w, http.StatusBadRequest, "metadata_provider_credentials_required", "Metadata provider credentials are required")
		return
	}
	if errors.Is(err, metadata.ErrUnsupportedProvider) {
		writeError(w, http.StatusBadRequest, "metadata_provider_unsupported", "Metadata provider does not support details")
		return
	}
	if metadata.IsRateLimited(err) {
		writeError(w, http.StatusTooManyRequests, "metadata_provider_rate_limited", "Metadata provider rate limit reached")
		return
	}
	if statusCode, ok := metadata.ProviderStatusCode(err); ok {
		switch statusCode {
		case http.StatusNotFound:
			writeError(w, http.StatusNotFound, "metadata_details_not_found", "Could not find metadata details")
		case http.StatusUnauthorized, http.StatusForbidden:
			writeError(w, http.StatusBadGateway, "metadata_provider_auth_failed", "Metadata provider authentication failed")
		default:
			writeError(w, http.StatusBadGateway, "metadata_provider_failed", "Metadata provider could not load details")
		}
		return
	}
	writeError(w, http.StatusBadGateway, "metadata_provider_failed", "Metadata provider could not load details")
}

func (s *Server) searchMetadataProvider(ctx context.Context, provider storage.MetadataProvider, request metadata.SearchRequest) ([]metadata.SearchResult, error) {
	cacheKey := strings.ToLower(strings.Join(strings.Fields(request.Query), " "))
	cached := []metadata.SearchResult{}
	found, err := s.settings.GetMetadataSearchCache(ctx, provider.ID, request.MediaType, cacheKey, request.Year, &cached)
	if err != nil {
		return nil, err
	}
	if found {
		return cached, nil
	}

	results, err := s.metadata.Search(ctx, metadataProviderConfig(provider), request)
	if err != nil {
		return nil, err
	}
	if err := s.settings.SetMetadataSearchCache(ctx, provider.ID, request.MediaType, cacheKey, request.Year, results, s.now().Add(24*time.Hour)); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *Server) discoverMetadataProvider(
	ctx context.Context,
	provider storage.MetadataProvider,
	mediaType string,
	section string,
	limit int,
	page int,
) ([]metadata.SearchResult, error) {
	if !providerHasCredentials(provider) {
		return []metadata.SearchResult{}, nil
	}
	cacheMediaType := mediaType
	if mediaType == "mixed" {
		cacheMediaType = "movie"
	}
	cacheKey := "discover:" + strings.ToLower(strings.TrimSpace(section)) + ":" + strconv.Itoa(limit) + ":" + strconv.Itoa(page)
	cached := []metadata.SearchResult{}
	found, err := s.settings.GetMetadataSearchCache(ctx, provider.ID, cacheMediaType, cacheKey, nil, &cached)
	if err != nil {
		return nil, err
	}
	if found {
		return cached, nil
	}

	results, err := s.metadata.Discover(ctx, metadataProviderConfig(provider), metadata.DiscoverRequest{
		MediaType: mediaType,
		Section:   section,
		Limit:     limit,
		Page:      page,
	})
	if err != nil {
		return nil, err
	}
	if err := s.settings.SetMetadataSearchCache(ctx, provider.ID, cacheMediaType, cacheKey, nil, results, s.now().Add(24*time.Hour)); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *Server) metadataProviderDetails(ctx context.Context, provider storage.MetadataProvider, request metadata.DetailsRequest) (metadata.Details, error) {
	cacheKey := metadataDetailsCacheKey(request.ExternalID)
	var cached metadata.Details
	found, err := s.settings.GetMetadataSearchCache(ctx, provider.ID, request.MediaType, cacheKey, nil, &cached)
	if err != nil {
		return metadata.Details{}, err
	}
	if found {
		return cached, nil
	}

	details, err := s.metadata.Details(ctx, metadataProviderConfig(provider), request)
	if err != nil {
		return metadata.Details{}, err
	}
	if err := s.settings.SetMetadataSearchCache(ctx, provider.ID, request.MediaType, cacheKey, nil, details, s.now().Add(24*time.Hour)); err != nil {
		return metadata.Details{}, err
	}
	return details, nil
}

func (s *Server) freshMetadataProviderDetails(ctx context.Context, provider storage.MetadataProvider, request metadata.DetailsRequest) (metadata.Details, error) {
	details, err := s.metadata.Details(ctx, metadataProviderConfig(provider), request)
	if err != nil {
		return metadata.Details{}, err
	}
	cacheKey := metadataDetailsCacheKey(request.ExternalID)
	if err := s.settings.SetMetadataSearchCache(ctx, provider.ID, request.MediaType, cacheKey, nil, details, s.now().Add(24*time.Hour)); err != nil {
		return metadata.Details{}, err
	}
	return details, nil
}

func metadataDetailsCacheKey(externalID string) string {
	return "details:v3:" + strings.ToLower(strings.TrimSpace(externalID))
}

type groupedMediaSearchRequest struct {
	query               string
	mediaTypes          []string
	year                *int32
	providerIDs         map[uuid.UUID]struct{}
	providerIDsProvided bool
	limit               int
	includeLibrary      bool
	includeProviders    bool
}

func (s *Server) groupedMediaSearch(ctx context.Context, request groupedMediaSearchRequest) ([]MediaSearchGroup, error) {
	limit := request.limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	groups := []MediaSearchGroup{}
	if request.includeLibrary {
		libraryItems, err := s.settings.SearchMediaItems(ctx, request.query, nil, limit)
		if err != nil {
			return nil, err
		}
		if len(libraryItems) > 0 {
			results := make([]MediaSearchResult, 0, len(libraryItems))
			for _, item := range libraryItems {
				results = append(results, mediaItemSearchResultResponse(item))
			}
			groups = append(groups, MediaSearchGroup{
				SourceType: "library",
				SourceName: "Library",
				Results:    results,
			})
		}
	}

	if !request.includeProviders {
		return groups, nil
	}

	providerGroups := map[uuid.UUID]int{}
	if request.providerIDsProvided && len(request.providerIDs) == 0 {
		return groups, nil
	}
	for _, mediaType := range request.mediaTypes {
		providers, err := s.settings.ListEnabledMetadataProviders(ctx, mediaType)
		if err != nil {
			return nil, err
		}
		for _, provider := range providers {
			if len(request.providerIDs) > 0 {
				if _, ok := request.providerIDs[provider.ID]; !ok {
					continue
				}
			}
			results, err := s.searchMetadataProvider(ctx, provider, metadata.SearchRequest{
				Query:     request.query,
				MediaType: mediaType,
				Year:      request.year,
			})
			if err != nil {
				continue
			}
			if len(results) == 0 {
				continue
			}
			index, ok := providerGroups[provider.ID]
			if !ok {
				providerID := openapi_types.UUID(provider.ID)
				groups = append(groups, MediaSearchGroup{
					SourceType: "provider",
					SourceName: provider.Name,
					ProviderId: &providerID,
					Results:    []MediaSearchResult{},
				})
				index = len(groups) - 1
				providerGroups[provider.ID] = index
			}
			for _, result := range results {
				if len(groups[index].Results) >= limit {
					break
				}
				groups[index].Results = append(groups[index].Results, metadataSearchResultResponse(result))
			}
		}
	}
	return groups, nil
}

func metadataSearchResultResponse(result metadata.SearchResult) MediaSearchResult {
	return MediaSearchResult{
		Title:            result.Title,
		Type:             MediaType(result.Type),
		Year:             result.Year,
		ExternalProvider: optionalString(result.ExternalProvider),
		ExternalId:       optionalString(result.ExternalID),
		Overview:         result.Overview,
		PosterPath:       result.PosterPath,
	}
}

func mediaItemSearchResultResponse(item storage.MediaItem) MediaSearchResult {
	id := openapi_types.UUID(item.ID)
	return MediaSearchResult{
		Id:               &id,
		Title:            item.Title,
		Type:             MediaType(item.Type),
		Year:             item.Year,
		ExternalProvider: item.ExternalProvider,
		ExternalId:       item.ExternalID,
		Overview:         item.Overview,
		PosterPath:       item.PosterPath,
	}
}
