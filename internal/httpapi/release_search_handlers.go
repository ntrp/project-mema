package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/storage"
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
	if err := s.settings.ReplaceReleaseSearchResults(r.Context(), mediaItemID, nil, nil); err != nil {
		writeError(w, http.StatusInternalServerError, "release_search_reset_failed", "Could not reset release search results")
		s.recordEvent(r.Context(), eventSeverityError, "media", "Release search reset failed", map[string]any{"mediaItemId": mediaItemID.String(), "error": err.Error()})
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
	profile, formats, languages := s.releaseDecisionContext(r.Context(), item)
	for _, release := range snapshot.Releases {
		block, blocked, err := s.settings.FindReleaseBlock(r.Context(), release)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "release_blocklist_failed", "Could not check release blocklist")
			return
		}
		var blockPtr *storage.ReleaseBlocklistItem
		if blocked {
			blockPtr = &block
		}
		response.Releases = append(
			response.Releases,
			releaseCandidateResponseWithBlock(item, release, profile, formats, languages, blockPtr),
		)
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) releaseDecisionContext(
	ctx context.Context,
	item storage.MediaItem,
) (*storage.MediaProfile, []storage.CustomFormat, []storage.Language) {
	var profile *storage.MediaProfile
	if item.QualityProfileID != nil {
		value, err := s.settings.GetMediaProfile(ctx, *item.QualityProfileID)
		if err == nil {
			profile = &value
		}
	}
	formats, _ := s.settings.ListCustomFormats(ctx)
	languages, _ := s.settings.ListLanguages(ctx)
	return profile, formats, languages
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
