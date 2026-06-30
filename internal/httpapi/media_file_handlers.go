package httpapi

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) RescanMediaItemFiles(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	item, err := s.settings.RescanMediaItemFiles(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not rescan media folder")
		return
	}
	writeJSON(w, http.StatusOK, mediaItemResponse(item))
}

func (s *Server) DeleteMediaItemFile(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body MediaFileDeleteRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	item, err := s.settings.DeleteMediaItemFile(r.Context(), uuid.UUID(id), body.Path)
	if err != nil {
		writeSettingsError(w, err, "Could not delete media file")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "media", "Media file deleted", map[string]any{"mediaItemId": uuid.UUID(id).String(), "path": body.Path})
	writeJSON(w, http.StatusOK, mediaItemResponse(item))
}

func (s *Server) EnqueueMediaAutomaticSearch(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	mediaItemID := uuid.UUID(id)
	if _, err := s.settings.GetMediaItem(r.Context(), mediaItemID); err != nil {
		writeSettingsError(w, err, "Could not find media item")
		return
	}
	jobID, err := s.jobs.EnqueueAutoSearchDownload(r.Context(), mediaItemID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "automatic_search_enqueue_failed", "Could not enqueue automatic search")
		s.recordEvent(r.Context(), eventSeverityError, "media", "Automatic search enqueue failed", map[string]any{"mediaItemId": mediaItemID.String(), "error": err.Error()})
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "media", "Automatic search queued", map[string]any{"mediaItemId": mediaItemID.String(), "jobId": jobID})
	writeJSON(w, http.StatusAccepted, JobEnqueueResponse{JobId: jobID, Message: "Automatic search queued"})
}
