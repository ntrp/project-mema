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

	"media-manager/internal/jobs"
	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func (s *Server) SearchMedia(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	var body MediaSearchRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	query := strings.TrimSpace(body.Query)
	if query == "" {
		writeError(w, http.StatusBadRequest, "invalid_query", "Search query is required")
		return
	}
	if !body.Type.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_type", "Media type is not supported")
		return
	}

	providers, err := s.settings.ListEnabledMetadataProviders(r.Context(), string(body.Type))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_provider_list_failed", "Could not list metadata providers")
		return
	}
	if len(providers) == 0 {
		writeJSON(w, http.StatusOK, MediaSearchResponse{
			Results: []MediaSearchResult{
				{
					Title: query,
					Type:  body.Type,
					Year:  body.Year,
				},
			},
		})
		return
	}

	response := MediaSearchResponse{Results: []MediaSearchResult{}}
	for _, provider := range providers {
		results, err := s.searchMetadataProvider(r.Context(), provider, metadata.SearchRequest{
			Query:     query,
			MediaType: string(body.Type),
			Year:      body.Year,
		})
		if err != nil {
			continue
		}
		for _, result := range results {
			response.Results = append(response.Results, metadataSearchResultResponse(result))
		}
		if len(response.Results) > 0 {
			break
		}
	}
	if len(response.Results) == 0 {
		response.Results = append(response.Results, MediaSearchResult{
			Title: query,
			Type:  body.Type,
			Year:  body.Year,
		})
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GetMediaDiscover(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	providers, err := s.settings.ListMetadataProviders(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_provider_list_failed", "Could not list metadata providers")
		return
	}
	blacklist, err := s.settings.ListDiscoverBlacklist(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "discover_blacklist_list_failed", "Could not list discover blacklist")
		return
	}

	response := MediaDiscoverResponse{Sections: make([]MediaDiscoverSection, 0, len(discoverSections))}
	for _, section := range discoverSections {
		response.Sections = append(response.Sections, s.discoverSectionResponse(r.Context(), providers, section, 20, 1, blacklist))
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GetMediaDiscoverSection(w http.ResponseWriter, r *http.Request, sectionId string, params GetMediaDiscoverSectionParams) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	section, ok := discoverSectionByID(sectionId)
	if !ok {
		writeError(w, http.StatusNotFound, "discover_section_not_found", "Discovery section was not found")
		return
	}
	providers, err := s.settings.ListMetadataProviders(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_provider_list_failed", "Could not list metadata providers")
		return
	}
	blacklist, err := s.settings.ListDiscoverBlacklist(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "discover_blacklist_list_failed", "Could not list discover blacklist")
		return
	}
	page := int32(1)
	if params.Page != nil {
		page = *params.Page
	}
	limit := int32(20)
	if params.Limit != nil {
		limit = *params.Limit
	}
	writeJSON(w, http.StatusOK, s.discoverSectionResponse(r.Context(), providers, section, int(limit), int(page), blacklist))
}

func (s *Server) AutocompleteMedia(w http.ResponseWriter, r *http.Request, params AutocompleteMediaParams) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	query := strings.TrimSpace(params.Query)
	if len(query) < 2 {
		writeError(w, http.StatusBadRequest, "invalid_query", "Search query must contain at least 2 characters")
		return
	}

	groups, err := s.groupedMediaSearch(r.Context(), groupedMediaSearchRequest{
		query:            query,
		mediaTypes:       []string{"movie", "series"},
		limit:            5,
		includeLibrary:   boolDefault(params.IncludeLibrary, true),
		includeProviders: boolDefault(params.IncludeProviders, true),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "media_autocomplete_failed", "Could not search media")
		return
	}
	writeJSON(w, http.StatusOK, MediaGroupedSearchResponse{Groups: groups})
}

func (s *Server) AdvancedSearchMedia(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	var body MediaAdvancedSearchRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	query := strings.TrimSpace(valueOrEmpty(body.Query))
	if query == "" {
		writeError(w, http.StatusBadRequest, "invalid_query", "Search query is required")
		return
	}
	mediaTypes := []string{"movie", "series"}
	if body.Type != nil {
		if !body.Type.Valid() {
			writeError(w, http.StatusBadRequest, "invalid_type", "Media type is not supported")
			return
		}
		mediaTypes = []string{string(*body.Type)}
	}
	limit := int32(20)
	if body.Limit != nil {
		limit = *body.Limit
	}

	providerIDs := map[uuid.UUID]struct{}{}
	if body.ProviderIds != nil {
		for _, id := range *body.ProviderIds {
			providerIDs[uuid.UUID(id)] = struct{}{}
		}
	}

	groups, err := s.groupedMediaSearch(r.Context(), groupedMediaSearchRequest{
		query:               query,
		mediaTypes:          mediaTypes,
		year:                body.Year,
		providerIDs:         providerIDs,
		providerIDsProvided: body.ProviderIds != nil,
		limit:               int(limit),
		includeLibrary:      true,
		includeProviders:    true,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "media_advanced_search_failed", "Could not search media")
		return
	}
	writeJSON(w, http.StatusOK, MediaGroupedSearchResponse{Groups: groups})
}

func (s *Server) GetMediaMetadataDetails(w http.ResponseWriter, r *http.Request, providerType MetadataProviderType, mediaType MediaType, externalID string) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	if !providerType.Valid() || !mediaType.Valid() || strings.TrimSpace(externalID) == "" {
		writeError(w, http.StatusBadRequest, "invalid_metadata_request", "Metadata provider, media type, and external id are required")
		return
	}

	providers, err := s.settings.ListMetadataProviders(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "metadata_provider_list_failed", "Could not list metadata providers")
		return
	}
	provider, ok := metadataProviderByType(providers, string(providerType))
	if !ok {
		writeError(w, http.StatusNotFound, "metadata_provider_not_found", "Metadata provider is not configured")
		return
	}

	details, err := s.metadataProviderDetails(r.Context(), provider, metadata.DetailsRequest{
		MediaType:  string(mediaType),
		ExternalID: externalID,
	})
	if err != nil {
		writeMetadataDetailsError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, metadataDetailsResponse(details))
}

func (s *Server) ListMediaItems(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	items, err := s.settings.ListMediaItems(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "media_list_failed", "Could not list media items")
		return
	}
	response := MediaItemListResponse{Items: make([]MediaItem, 0, len(items))}
	for _, item := range items {
		response.Items = append(response.Items, mediaItemResponse(item))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreateMediaItem(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body MediaItemCreateRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := mediaItemInput(body)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_media_item", "Media item title and type are required")
		return
	}
	if err := s.validateMediaTarget(r.Context(), input.QualityProfileID, input.LibraryFolderID); err != nil {
		writeMediaTargetError(w, err)
		return
	}

	items, err := s.createMediaForAdd(r.Context(), input)
	if err != nil {
		if errors.Is(err, errMediaCollectionUnavailable) {
			writeError(w, http.StatusBadRequest, "collection_unavailable", "Selected media is not part of an available collection")
			return
		}
		writeError(w, http.StatusInternalServerError, "media_create_failed", "Could not add media item")
		return
	}
	if body.StartSearch {
		s.enqueueAutomaticSearch(r.Context(), items)
	}
	writeJSON(w, http.StatusCreated, mediaItemResponse(items[0]))
}

func (s *Server) ListMediaRequests(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}

	requests, err := s.settings.ListMediaRequests(r.Context(), uuid.UUID(session.user.Id), session.user.Role == Admin)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "media_request_list_failed", "Could not list media requests")
		return
	}
	response := MediaRequestListResponse{Requests: make([]MediaRequest, 0, len(requests))}
	for _, request := range requests {
		response.Requests = append(response.Requests, mediaRequestResponse(request))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreateMediaRequest(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}

	var body MediaRequestCreateRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := mediaRequestInput(body, uuid.UUID(session.user.Id))
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_media_request", "Media request title and type are required")
		return
	}

	request, err := s.settings.CreateMediaRequest(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "media_request_create_failed", "Could not create media request")
		return
	}
	writeJSON(w, http.StatusCreated, mediaRequestResponse(request))
}

func (s *Server) GetMediaRequest(w http.ResponseWriter, r *http.Request, id ResourceId) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}

	request, err := s.settings.GetMediaRequest(r.Context(), uuid.UUID(id), uuid.UUID(session.user.Id), session.user.Role == Admin)
	if err != nil {
		writeSettingsError(w, err, "Could not find media request")
		return
	}
	writeJSON(w, http.StatusOK, mediaRequestResponse(request))
}

func (s *Server) ApproveMediaRequest(w http.ResponseWriter, r *http.Request, id ResourceId) {
	session, ok := s.requireAdmin(w, r)
	if !ok {
		return
	}

	var body MediaRequestApproveRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	qualityProfileID := strings.TrimSpace(body.QualityProfileId)
	input := storage.MediaRequestApprovalInput{
		QualityProfileID: qualityProfileID,
		LibraryFolderID:  uuid.UUID(body.LibraryFolderId),
	}
	if err := s.validateMediaTarget(r.Context(), &input.QualityProfileID, &input.LibraryFolderID); err != nil {
		writeMediaTargetError(w, err)
		return
	}

	existingRequest, err := s.settings.GetMediaRequest(r.Context(), uuid.UUID(id), uuid.UUID(session.user.Id), true)
	if err != nil {
		writeSettingsError(w, err, "Could not find media request")
		return
	}
	addInputs, err := s.mediaAddInputs(r.Context(), mediaInputFromRequest(existingRequest, input))
	if err != nil {
		if errors.Is(err, errMediaCollectionUnavailable) {
			writeError(w, http.StatusBadRequest, "collection_unavailable", "Selected media is not part of an available collection")
			return
		}
		writeError(w, http.StatusInternalServerError, "collection_lookup_failed", "Could not load media collection")
		return
	}
	for index := range addInputs {
		enriched, err := s.enrichMediaItemInput(r.Context(), addInputs[index])
		if err != nil {
			writeMetadataDetailsError(w, err)
			return
		}
		addInputs[index] = applySeriesMonitoring(enriched)
	}
	if len(addInputs) > 0 {
		input.MediaInput = &addInputs[0]
	}

	request, item, err := s.settings.ApproveMediaRequest(r.Context(), uuid.UUID(id), input)
	if err != nil {
		writeMediaRequestError(w, err)
		return
	}
	items := []storage.MediaItem{item}
	if len(addInputs) > 1 {
		items, err = s.createMediaInputs(r.Context(), addInputs)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "media_create_failed", "Could not add media collection")
			return
		}
	}
	if existingRequest.MonitorMode != "none" {
		s.enqueueAutomaticSearch(r.Context(), items)
	}
	writeJSON(w, http.StatusOK, MediaRequestApproveResponse{
		Request:   mediaRequestResponse(request),
		MediaItem: mediaItemResponse(item),
	})
}

func (s *Server) EnqueueMediaReleaseSearch(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	mediaItemID := uuid.UUID(id)
	if _, err := s.settings.GetMediaItem(r.Context(), mediaItemID); err != nil {
		writeSettingsError(w, err, "Could not find media item")
		return
	}

	jobID, err := s.jobs.EnqueueReleaseSearch(r.Context(), mediaItemID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "release_search_enqueue_failed", "Could not enqueue release search")
		s.recordEvent(r.Context(), eventSeverityError, "media", "Release search enqueue failed", map[string]any{"mediaItemId": mediaItemID.String(), "error": err.Error()})
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "media", "Release search queued", map[string]any{"mediaItemId": mediaItemID.String(), "jobId": jobID})
	writeJSON(w, http.StatusAccepted, JobEnqueueResponse{
		JobId:   jobID,
		Message: "Release search queued",
	})
}

func (s *Server) SearchMediaReleases(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	mediaItemID := uuid.UUID(id)
	if _, err := s.settings.GetMediaItem(r.Context(), mediaItemID); err != nil {
		writeSettingsError(w, err, "Could not find media item")
		return
	}
	snapshot, err := s.settings.ListReleaseSearchResults(r.Context(), mediaItemID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "release_results_failed", "Could not list release search results")
		return
	}
	response := ReleaseSearchResponse{
		Releases: make([]ReleaseCandidate, 0, len(snapshot.Releases)),
		Errors:   snapshot.Errors,
	}
	for _, release := range snapshot.Releases {
		response.Releases = append(response.Releases, releaseCandidateResponse(release))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) GrabMediaRelease(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	item, err := s.settings.GetMediaItem(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find media item")
		return
	}
	var body GrabReleaseRequest
	if !decodeJSON(w, r, &body) {
		return
	}

	release, err := s.settings.GetReleaseCandidate(r.Context(), uuid.UUID(body.ReleaseId), item.ID)
	if err != nil {
		writeSettingsError(w, err, "Could not find release candidate")
		return
	}

	clients, err := s.settings.ListEnabledDownloadClients(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "download_client_list_failed", "Could not list enabled download clients")
		return
	}
	if len(clients) == 0 {
		writeError(w, http.StatusBadRequest, "no_download_client", "No enabled download client is configured")
		return
	}

	client := clients[0]
	activity, err := s.settings.CreateDownloadActivity(r.Context(), storage.DownloadActivityInput{
		MediaItemID:        item.ID,
		ReleaseTitle:       release.Title,
		IndexerName:        release.IndexerName,
		DownloadClientName: client.Name,
		DownloadURL:        release.DownloadURL,
		Status:             "queued",
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "activity_create_failed", "Could not record download activity")
		return
	}
	activity.MediaTitle = item.Title
	activity.MediaType = item.Type

	jobID, err := s.jobs.EnqueueGrabRelease(r.Context(), jobs.GrabReleaseArgs{
		ActivityID:  activity.ID.String(),
		MediaItemID: item.ID.String(),
		Title:       release.Title,
		DownloadURL: release.DownloadURL,
		IndexerName: release.IndexerName,
	})
	if err != nil {
		enqueueError := "Could not enqueue download job"
		_, _ = s.settings.FailDownloadActivity(r.Context(), activity.ID, &enqueueError, "download")
		s.recordEvent(r.Context(), eventSeverityError, "downloads", "Download enqueue failed", map[string]any{"mediaItemId": item.ID.String(), "activityId": activity.ID.String(), "releaseTitle": release.Title, "error": err.Error()})
		writeError(w, http.StatusInternalServerError, "download_enqueue_failed", enqueueError)
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "downloads", "Download queued", map[string]any{"mediaItemId": item.ID.String(), "activityId": activity.ID.String(), "releaseTitle": release.Title, "jobId": jobID})
	writeJSON(w, http.StatusAccepted, GrabReleaseResponse{
		JobId:    jobID,
		Message:  "Download queued",
		Activity: downloadActivityResponse(activity),
	})
}

func (s *Server) ListDownloadActivity(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	activities, err := s.settings.ListDownloadActivity(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "activity_list_failed", "Could not list download activity")
		return
	}
	response := DownloadActivityListResponse{Activities: make([]DownloadActivity, 0, len(activities))}
	for _, activity := range activities {
		response.Activities = append(response.Activities, downloadActivityResponse(activity))
	}
	writeJSON(w, http.StatusOK, response)
}

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
	cacheKey := "details:v3:" + strings.ToLower(strings.TrimSpace(request.ExternalID))
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

func metadataDetailsResponse(details metadata.Details) MediaMetadataDetails {
	genres := append([]string(nil), details.Genres...)
	keywords := append([]string(nil), details.Keywords...)
	facts := make([]MediaMetadataFact, 0, len(details.Facts))
	for _, fact := range details.Facts {
		facts = append(facts, MediaMetadataFact{
			Label: fact.Label,
			Value: fact.Value,
		})
	}
	seasons := make([]MediaMetadataSeason, 0, len(details.Seasons))
	for _, season := range details.Seasons {
		episodes := make([]MediaMetadataEpisode, 0, len(season.Episodes))
		for _, episode := range season.Episodes {
			episodes = append(episodes, MediaMetadataEpisode{
				Name:          episode.Name,
				EpisodeNumber: episode.EpisodeNumber,
				Overview:      episode.Overview,
				AirDate:       episode.AirDate,
				StillPath:     episode.StillPath,
			})
		}
		seasons = append(seasons, MediaMetadataSeason{
			Name:         season.Name,
			EpisodeCount: season.EpisodeCount,
			AirDate:      season.AirDate,
			PosterPath:   season.PosterPath,
			Episodes:     &episodes,
		})
	}
	cast := make([]MediaMetadataPerson, 0, len(details.Cast))
	for _, person := range details.Cast {
		cast = append(cast, MediaMetadataPerson{
			Name:        person.Name,
			Role:        person.Role,
			ProfilePath: person.ProfilePath,
		})
	}
	recommendations := make([]MediaSearchResult, 0, len(details.Recommendations))
	for _, result := range details.Recommendations {
		recommendations = append(recommendations, metadataSearchResultResponse(result))
	}
	similar := make([]MediaSearchResult, 0, len(details.Similar))
	for _, result := range details.Similar {
		similar = append(similar, metadataSearchResultResponse(result))
	}
	return MediaMetadataDetails{
		Title:            details.Title,
		Type:             MediaType(details.Type),
		Year:             details.Year,
		ExternalProvider: details.ExternalProvider,
		ExternalId:       details.ExternalID,
		Overview:         details.Overview,
		PosterPath:       details.PosterPath,
		CollectionId:     details.CollectionID,
		CollectionName:   details.CollectionName,
		BackdropPath:     details.BackdropPath,
		Status:           details.Status,
		OriginalLanguage: details.OriginalLanguage,
		ReleaseDate:      details.ReleaseDate,
		FirstAirDate:     details.FirstAirDate,
		RuntimeMinutes:   details.RuntimeMinutes,
		SeasonCount:      details.SeasonCount,
		EpisodeCount:     details.EpisodeCount,
		VoteAverage:      details.VoteAverage,
		Genres:           &genres,
		Keywords:         &keywords,
		Facts:            &facts,
		Seasons:          &seasons,
		Cast:             &cast,
		Recommendations:  &recommendations,
		Similar:          &similar,
	}
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func boolDefault(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

func (s *Server) validateMediaTarget(ctx context.Context, qualityProfileID *string, libraryFolderID *uuid.UUID) error {
	if qualityProfileID == nil || strings.TrimSpace(*qualityProfileID) == "" {
		return errMissingQualityProfile
	}
	exists, err := s.settings.MediaProfileExists(ctx, strings.TrimSpace(*qualityProfileID))
	if err != nil {
		return err
	}
	if !exists {
		return errUnsupportedQualityProfile
	}
	if libraryFolderID == nil {
		return errMissingLibraryFolder
	}
	exists, err = s.settings.LibraryFolderExists(ctx, *libraryFolderID)
	if err != nil {
		return err
	}
	if !exists {
		return storage.ErrNotFound
	}
	return nil
}

func writeMediaTargetError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errMissingQualityProfile):
		writeError(w, http.StatusBadRequest, "quality_profile_required", "Quality profile is required")
	case errors.Is(err, errUnsupportedQualityProfile):
		writeError(w, http.StatusBadRequest, "quality_profile_invalid", "Quality profile is not supported")
	case errors.Is(err, errMissingLibraryFolder):
		writeError(w, http.StatusBadRequest, "library_folder_required", "Library folder is required")
	case errors.Is(err, storage.ErrNotFound):
		writeError(w, http.StatusNotFound, "library_folder_not_found", "Library folder was not found")
	default:
		writeError(w, http.StatusInternalServerError, "media_target_validation_failed", "Could not validate media target")
	}
}

func writeMediaRequestError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, storage.ErrRequestClosed):
		writeError(w, http.StatusBadRequest, "media_request_closed", "Media request is no longer pending")
	case errors.Is(err, storage.ErrNotFound):
		writeError(w, http.StatusNotFound, "media_request_not_found", "Could not find media request")
	default:
		writeError(w, http.StatusInternalServerError, "media_request_update_failed", "Could not update media request")
	}
}

func providerHasCredentials(provider storage.MetadataProvider) bool {
	return optionalTrimmedString(provider.APIKey) != nil || optionalTrimmedString(provider.AccessToken) != nil
}

func discoverProvider(providers []storage.MetadataProvider) (storage.MetadataProvider, bool) {
	for _, provider := range providers {
		if provider.Enabled && provider.Type == "tmdb" && providerHasCredentials(provider) {
			return provider, true
		}
	}
	return storage.MetadataProvider{}, false
}

func metadataProviderByType(providers []storage.MetadataProvider, providerType string) (storage.MetadataProvider, bool) {
	for _, provider := range providers {
		if provider.Enabled && provider.Type == providerType && providerHasCredentials(provider) {
			return provider, true
		}
	}
	return storage.MetadataProvider{}, false
}

func (s *Server) discoverSectionResponse(
	ctx context.Context,
	providers []storage.MetadataProvider,
	section discoverSection,
	limit int,
	page int,
	blacklist []storage.DiscoverBlacklistItem,
) MediaDiscoverSection {
	providerName := "TMDB"
	results := []MediaSearchResult{}
	if provider, ok := discoverProvider(providers); ok {
		providerName = provider.Name
		for _, request := range section.requests {
			providerResults, err := s.discoverMetadataProvider(ctx, provider, request.mediaType, request.id, limit, page)
			if err != nil {
				continue
			}
			for _, result := range providerResults {
				results = append(results, metadataSearchResultResponse(result))
			}
		}
	}
	return MediaDiscoverSection{
		Id:           section.responseID,
		Title:        section.title,
		ProviderName: providerName,
		MediaType:    MediaDiscoverMediaType(section.mediaType),
		Results:      filterDiscoverBlacklist(dedupeMediaSearchResults(results), blacklist),
	}
}

func dedupeMediaSearchResults(results []MediaSearchResult) []MediaSearchResult {
	seen := map[string]struct{}{}
	deduped := make([]MediaSearchResult, 0, len(results))
	for _, result := range results {
		key := string(result.Type) + ":" + valueOrEmpty(result.ExternalProvider) + ":" + valueOrEmpty(result.ExternalId)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		deduped = append(deduped, result)
	}
	return deduped
}

type discoverSection struct {
	responseID string
	title      string
	mediaType  string
	requests   []discoverSectionRequest
}

type discoverSectionRequest struct {
	mediaType string
	id        string
}

var discoverSections = []discoverSection{
	{responseID: "trending", title: "Trending", mediaType: "mixed", requests: []discoverSectionRequest{
		{mediaType: "mixed", id: "trending"},
	}},
	{responseID: "movie-popular", title: "Popular Movies", mediaType: "movie", requests: []discoverSectionRequest{{mediaType: "movie", id: "popular"}}},
	{responseID: "movie-upcoming", title: "Upcoming Movies", mediaType: "movie", requests: []discoverSectionRequest{{mediaType: "movie", id: "upcoming"}}},
	{responseID: "movie-top-rated", title: "Top Rated Movies", mediaType: "movie", requests: []discoverSectionRequest{{mediaType: "movie", id: "top_rated"}}},
	{responseID: "series-popular", title: "Popular Series", mediaType: "series", requests: []discoverSectionRequest{{mediaType: "series", id: "popular"}}},
	{responseID: "series-on-the-air", title: "Airing Series", mediaType: "series", requests: []discoverSectionRequest{{mediaType: "series", id: "on_the_air"}}},
	{responseID: "series-top-rated", title: "Top Rated Series", mediaType: "series", requests: []discoverSectionRequest{{mediaType: "series", id: "top_rated"}}},
}

func discoverSectionByID(id string) (discoverSection, bool) {
	for _, section := range discoverSections {
		if section.responseID == id {
			return section, true
		}
	}
	return discoverSection{}, false
}

var (
	errMissingQualityProfile     = errors.New("quality profile is required")
	errUnsupportedQualityProfile = errors.New("quality profile is not supported")
	errMissingLibraryFolder      = errors.New("library folder is required")
)
