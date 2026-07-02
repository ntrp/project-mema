package httpapi

import (
	"net/http"

	"github.com/google/uuid"

	"media-manager/internal/jobs"
)

func (s *Server) StreamMediaReleaseSearch(w http.ResponseWriter, r *http.Request, id ResourceId, params StreamMediaReleaseSearchParams) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming_unavailable", "Streaming is unavailable")
		return
	}

	mediaItemID := uuid.UUID(id)
	item, err := s.settings.GetMediaItem(r.Context(), mediaItemID)
	if err != nil {
		writeSSE(w, flusher, "media.release_search.error", jobs.ReleaseSearchProgressEvent{Kind: "error", Message: "Could not find media item"})
		return
	}
	if err := s.settings.ReplaceReleaseSearchResults(r.Context(), mediaItemID, nil, nil); err != nil {
		writeSSE(w, flusher, "media.release_search.error", jobs.ReleaseSearchProgressEvent{Kind: "error", Message: "Could not reset release search results"})
		s.recordEvent(r.Context(), eventSeverityError, "media", "Release search reset failed", map[string]any{"mediaItemId": mediaItemID.String(), "error": err.Error()})
		return
	}

	progress := func(event jobs.ReleaseSearchProgressEvent) {
		writeSSE(w, flusher, "media.release_search.status", event)
	}
	releases, searchErrors, err := jobs.SearchManualReleases(r.Context(), s.settings, s.indexers, item, stringValue(params.Query), s.events, progress)
	if err != nil {
		writeSSE(w, flusher, "media.release_search.error", jobs.ReleaseSearchProgressEvent{Kind: "error", Message: "Release search failed"})
		s.recordEvent(r.Context(), eventSeverityError, "media", "Release search failed", map[string]any{"mediaItemId": mediaItemID.String(), "error": err.Error()})
		return
	}
	if err := s.settings.ReplaceReleaseSearchResults(r.Context(), mediaItemID, releases, searchErrors); err != nil {
		writeSSE(w, flusher, "media.release_search.error", jobs.ReleaseSearchProgressEvent{Kind: "error", Message: "Could not store release search results"})
		s.recordEvent(r.Context(), eventSeverityError, "media", "Release search result store failed", map[string]any{"mediaItemId": mediaItemID.String(), "error": err.Error()})
		return
	}

	snapshot, err := s.settings.ListReleaseSearchResults(r.Context(), mediaItemID)
	if err != nil {
		writeSSE(w, flusher, "media.release_search.error", jobs.ReleaseSearchProgressEvent{Kind: "error", Message: "Could not load release search results"})
		return
	}
	response := ReleaseSearchResponse{Releases: make([]ReleaseCandidate, 0, len(snapshot.Releases)), Errors: snapshot.Errors}
	profile, formats, languages := s.releaseDecisionContext(r.Context(), item)
	for _, release := range snapshot.Releases {
		response.Releases = append(
			response.Releases,
			releaseCandidateResponse(item, release, profile, formats, languages),
		)
	}
	writeSSE(w, flusher, "media.release_search.result", response)
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
