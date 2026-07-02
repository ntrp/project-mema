package httpapi

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (s *Server) EnqueueMediaReleaseSearch(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	query, ok := releaseSearchQuery(w, r)
	if !ok {
		return
	}

	mediaItemID := uuid.UUID(id)
	if _, err := s.settings.GetMediaItem(r.Context(), mediaItemID); err != nil {
		writeSettingsError(w, err, "Could not find media item")
		return
	}

	jobID, err := s.jobs.EnqueueReleaseSearch(r.Context(), mediaItemID, query)
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
	item, err := s.settings.GetMediaItem(r.Context(), mediaItemID)
	if err != nil {
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
		response.Releases = append(response.Releases, releaseCandidateResponse(item, release))
	}
	writeJSON(w, http.StatusOK, response)
}

func releaseSearchQuery(w http.ResponseWriter, r *http.Request) (string, bool) {
	var body ReleaseSearchRequest
	if r.Body == nil || r.ContentLength == 0 {
		return "", true
	}
	if !decodeJSON(w, r, &body) {
		return "", false
	}
	if body.Query == nil {
		return "", true
	}
	return strings.TrimSpace(*body.Query), true
}
