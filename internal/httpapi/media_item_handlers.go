package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

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

func (s *Server) UpdateMediaItem(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body MediaItemUpdateRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	if body.QualityProfileId != nil && strings.TrimSpace(*body.QualityProfileId) == "" {
		writeError(w, http.StatusBadRequest, "invalid_media_item_settings", "Quality profile is required")
		return
	}
	if body.MinimumAvailability != nil && !body.MinimumAvailability.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_media_item_settings", "Minimum availability is not supported")
		return
	}
	if body.MonitorMode != nil && !body.MonitorMode.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_media_monitor_mode", "Monitor mode is not supported")
		return
	}
	if body.LibraryFolderId != nil {
		exists, err := s.settings.LibraryFolderExists(r.Context(), uuid.UUID(*body.LibraryFolderId))
		if err != nil {
			writeError(w, http.StatusInternalServerError, "media_target_validation_failed", "Could not validate media target")
			return
		}
		if !exists {
			writeError(w, http.StatusBadRequest, "library_folder_invalid", "Library folder is not supported")
			return
		}
	}
	if body.QualityProfileId != nil {
		exists, err := s.settings.MediaProfileExists(r.Context(), strings.TrimSpace(*body.QualityProfileId))
		if err != nil {
			writeError(w, http.StatusInternalServerError, "media_target_validation_failed", "Could not validate media target")
			return
		}
		if !exists {
			writeError(w, http.StatusBadRequest, "quality_profile_invalid", "Quality profile is not supported")
			return
		}
	}

	item, err := s.settings.UpdateMediaItemOptions(r.Context(), uuid.UUID(id), storage.MediaItemOptionsInput{
		QualityProfileID:     body.QualityProfileId,
		MinimumAvailability:  optionalMinimumAvailability(body.MinimumAvailability),
		LibraryFolderID:      optionalUUID(body.LibraryFolderId),
		Monitored:            body.Monitored,
		MonitorMode:          optionalMediaMonitorMode(body.MonitorMode),
		Seasons:              storageMediaSeasons(body.Seasons),
		MonitorSeasonName:    optionalTrimmedString(body.MonitorSeasonName),
		MonitorEpisodeNumber: body.MonitorEpisodeNumber,
		SeasonMonitored:      body.SeasonMonitored,
		EpisodeMonitored:     body.EpisodeMonitored,
	})
	if err != nil {
		writeSettingsError(w, err, "Could not update media item")
		return
	}
	writeJSON(w, http.StatusOK, mediaItemResponse(item))
}

func (s *Server) ListMediaRequests(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}

	requests, err := s.settings.ListMediaRequests(r.Context(), uuid.UUID(session.user.Id), session.user.Role == UserRoleAdmin)
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

	request, err := s.settings.GetMediaRequest(r.Context(), uuid.UUID(id), uuid.UUID(session.user.Id), session.user.Role == UserRoleAdmin)
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
	if !body.MonitorMode.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_media_monitor_mode", "Monitor mode is not supported")
		return
	}
	if !body.MinimumAvailability.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_media_item_settings", "Minimum availability is not supported")
		return
	}
	input := storage.MediaRequestApprovalInput{
		QualityProfileID:    qualityProfileID,
		LibraryFolderID:     uuid.UUID(body.LibraryFolderId),
		MonitorMode:         string(body.MonitorMode),
		SeriesType:          nil,
		MinimumAvailability: string(body.MinimumAvailability),
		Tags:                optionalStringSlice(body.Tags),
	}
	if body.SeriesType != nil {
		seriesType := string(*body.SeriesType)
		input.SeriesType = &seriesType
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
	input = storage.NormalizeMediaRequestApprovalOptions(existingRequest.Type, input)
	if err := s.validateMediaRequestLibraryFolderKind(r.Context(), input.LibraryFolderID, existingRequest.Type); err != nil {
		writeMediaTargetError(w, err)
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
	if input.MonitorMode != "none" && body.StartSearch {
		s.enqueueAutomaticSearch(r.Context(), items)
	}
	writeJSON(w, http.StatusOK, MediaRequestApproveResponse{
		Request:   mediaRequestResponse(request),
		MediaItem: mediaItemResponse(item),
	})
}
