package httpapi

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/jobs"
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

	writeJSON(w, http.StatusOK, MediaSearchResponse{
		Results: []MediaSearchResult{
			{
				Title: query,
				Type:  body.Type,
				Year:  body.Year,
			},
		},
	})
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
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	var body MediaItemRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := mediaItemInput(body)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_media_item", "Media item title and type are required")
		return
	}

	item, err := s.settings.CreateMediaItem(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "media_create_failed", "Could not add media item")
		return
	}
	writeJSON(w, http.StatusCreated, mediaItemResponse(item))
}

func (s *Server) DeleteMediaItem(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	if err := s.settings.DeleteMediaItem(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not find media item")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) EnqueueMediaReleaseSearch(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireSession(w, r); !ok {
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
		return
	}
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
	if _, ok := s.requireSession(w, r); !ok {
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
		_, _ = s.settings.UpdateDownloadActivityStatus(r.Context(), activity.ID, "failed", &enqueueError)
		writeError(w, http.StatusInternalServerError, "download_enqueue_failed", enqueueError)
		return
	}
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
