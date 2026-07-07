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
		s.recordMetadataSearchHistory(ctx, provider, request.MediaType, cacheKey, request.Year, true, cached, nil)
		return cached, nil
	}

	results, err := s.metadata.Search(ctx, metadataProviderConfig(provider), request)
	if err != nil {
		s.recordMetadataSearchHistory(ctx, provider, request.MediaType, cacheKey, request.Year, false, nil, err)
		return nil, err
	}
	expiresAt := s.now().Add(24 * time.Hour)
	s.recordMetadataSearchHistory(ctx, provider, request.MediaType, cacheKey, request.Year, false, results, nil)
	if err := s.settings.SetMetadataSearchCache(ctx, provider.ID, request.MediaType, cacheKey, request.Year, results, expiresAt); err != nil {
		return nil, err
	}
	s.publishMetadataCacheUpdated(ctx, provider, request.MediaType, cacheKey, request.Year, results, expiresAt)
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
		s.recordMetadataSearchHistory(ctx, provider, mediaType, cacheKey, nil, true, cached, nil)
		return cached, nil
	}

	results, err := s.metadata.Discover(ctx, metadataProviderConfig(provider), metadata.DiscoverRequest{
		MediaType: mediaType,
		Section:   section,
		Limit:     limit,
		Page:      page,
	})
	if err != nil {
		s.recordMetadataSearchHistory(ctx, provider, mediaType, cacheKey, nil, false, nil, err)
		return nil, err
	}
	expiresAt := s.now().Add(24 * time.Hour)
	s.recordMetadataSearchHistory(ctx, provider, mediaType, cacheKey, nil, false, results, nil)
	if err := s.settings.SetMetadataSearchCache(ctx, provider.ID, cacheMediaType, cacheKey, nil, results, expiresAt); err != nil {
		return nil, err
	}
	s.publishMetadataCacheUpdated(ctx, provider, cacheMediaType, cacheKey, nil, results, expiresAt)
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
		s.recordMetadataSearchHistory(ctx, provider, request.MediaType, cacheKey, nil, true, cached, nil)
		return cached, nil
	}

	details, err := s.metadata.Details(ctx, metadataProviderConfig(provider), request)
	if err != nil {
		s.recordMetadataSearchHistory(ctx, provider, request.MediaType, cacheKey, nil, false, nil, err)
		return metadata.Details{}, err
	}
	expiresAt := s.now().Add(24 * time.Hour)
	s.recordMetadataSearchHistory(ctx, provider, request.MediaType, cacheKey, nil, false, details, nil)
	if err := s.settings.SetMetadataSearchCache(ctx, provider.ID, request.MediaType, cacheKey, nil, details, expiresAt); err != nil {
		return metadata.Details{}, err
	}
	s.publishMetadataCacheUpdated(ctx, provider, request.MediaType, cacheKey, nil, details, expiresAt)
	return details, nil
}

func (s *Server) freshMetadataProviderDetails(ctx context.Context, provider storage.MetadataProvider, request metadata.DetailsRequest) (metadata.Details, error) {
	details, err := s.metadata.Details(ctx, metadataProviderConfig(provider), request)
	if err != nil {
		cacheKey := metadataDetailsCacheKey(request.ExternalID)
		s.recordMetadataSearchHistory(ctx, provider, request.MediaType, cacheKey, nil, false, nil, err)
		return metadata.Details{}, err
	}
	cacheKey := metadataDetailsCacheKey(request.ExternalID)
	expiresAt := s.now().Add(24 * time.Hour)
	s.recordMetadataSearchHistory(ctx, provider, request.MediaType, cacheKey, nil, false, details, nil)
	if err := s.settings.SetMetadataSearchCache(ctx, provider.ID, request.MediaType, cacheKey, nil, details, expiresAt); err != nil {
		return metadata.Details{}, err
	}
	s.publishMetadataCacheUpdated(ctx, provider, request.MediaType, cacheKey, nil, details, expiresAt)
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
	includePeople       bool
	sortByPopularity    bool
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

	if len(request.mediaTypes) > 0 {
		providerGroups, err := s.providerMediaSearchGroups(ctx, request, limit)
		if err != nil {
			return nil, err
		}
		groups = append(groups, providerGroups...)
	}
	if request.includePeople {
		peopleGroups, err := s.providerPersonSearchGroups(ctx, request, limit)
		if err != nil {
			return nil, err
		}
		groups = append(groups, peopleGroups...)
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
		ExternalUrl:      result.ExternalURL,
		Overview:         result.Overview,
		PosterPath:       result.PosterPath,
		Popularity:       result.Popularity,
		ReleaseDate:      dateOnly(result.ReleaseDate),
		RuntimeMinutes:   result.RuntimeMinutes,
		VoteAverage:      result.VoteAverage,
		VoteCount:        result.VoteCount,
		OriginalLanguage: result.OriginalLanguage,
		ContentRating:    result.ContentRating,
		Genres:           optionalResponseStrings(result.Genres),
		Keywords:         optionalResponseStrings(result.Keywords),
		Studios:          optionalResponseStrings(result.Studios),
		BackdropPath:     result.BackdropPath,
	}
}

func dateOnly(value *string) *openapi_types.Date {
	if value == nil {
		return nil
	}
	date, err := time.Parse("2006-01-02", *value)
	if err != nil {
		return nil
	}
	return &openapi_types.Date{Time: date}
}

func optionalResponseStrings(values []string) *[]string {
	if len(values) == 0 {
		return nil
	}
	copy := append([]string(nil), values...)
	return &copy
}

func popularityValue(result metadata.SearchResult) float64 {
	if result.Popularity == nil {
		return 0
	}
	return *result.Popularity
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
		ExternalUrl:      mediaItemExternalURL(item),
		Overview:         item.Overview,
		PosterPath:       item.PosterPath,
	}
}
