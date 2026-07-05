package httpapi

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/jobs"
)

func (s *Server) EnqueueMediaSubtitleSearch(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	item, err := s.settings.GetMediaItem(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find media item")
		return
	}
	body, ok := subtitleSearchRequestBody(w, r)
	if !ok {
		return
	}
	languageID := subtitleSearchValue(body.LanguageId)
	filePath := subtitleSearchValue(body.FilePath)
	if languageID == "" && len(item.SubtitleLanguages) == 0 {
		writeError(w, http.StatusBadRequest, "subtitle_language_missing", "No wanted subtitle language is configured")
		return
	}
	jobID, err := s.jobs.EnqueueSubtitleSearch(r.Context(), jobs.SubtitleSearchArgs{
		MediaItemID: item.ID.String(),
		LanguageID:  languageID,
		FilePath:    filePath,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "subtitle_search_enqueue_failed", "Could not enqueue subtitle search")
		s.recordEvent(r.Context(), eventSeverityError, "subtitles", "Subtitle search enqueue failed", map[string]any{"mediaItemId": item.ID.String(), "error": err.Error()})
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "subtitles", "Subtitle search queued", map[string]any{"mediaItemId": item.ID.String(), "jobId": jobID})
	writeJSON(w, http.StatusAccepted, JobEnqueueResponse{JobId: jobID, Message: "Subtitle search queued"})
}

func subtitleSearchRequestBody(w http.ResponseWriter, r *http.Request) (SubtitleSearchRequest, bool) {
	if r.Body == nil || r.ContentLength == 0 {
		return SubtitleSearchRequest{}, true
	}
	var body SubtitleSearchRequest
	if !decodeJSON(w, r, &body) {
		return SubtitleSearchRequest{}, false
	}
	return body, true
}

func subtitleSearchValue(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}
