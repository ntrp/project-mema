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
