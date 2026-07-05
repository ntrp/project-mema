package httpapi

import (
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func (s *Server) ListMediaFileHistory(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	mediaItemID := uuid.UUID(id)
	if _, err := s.settings.GetMediaItem(r.Context(), mediaItemID); err != nil {
		writeSettingsError(w, err, "Could not find media item")
		return
	}
	entries, err := s.settings.ListMediaFileHistory(r.Context(), mediaItemID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "file_history_list_failed", "Could not list file history")
		return
	}
	response := MediaFileHistoryResponse{Entries: make([]MediaFileHistoryEntry, 0, len(entries))}
	for _, entry := range entries {
		response.Entries = append(response.Entries, mediaFileHistoryResponse(entry))
	}
	writeJSON(w, http.StatusOK, response)
}

func mediaFileHistoryResponse(entry storage.MediaFileHistoryEntry) MediaFileHistoryEntry {
	response := MediaFileHistoryEntry{
		Id:             openapi_types.UUID(entry.ID),
		FilePath:       entry.FilePath,
		Operation:      MediaFileHistoryEntryOperation(entry.Operation),
		Status:         MediaFileHistoryEntryStatus(entry.Status),
		ActorType:      MediaFileHistoryEntryActorType(entry.ActorType),
		ActorId:        entry.ActorID,
		JobId:          entry.JobID,
		Details:        entry.Details,
		FailureDetails: entry.FailureDetails,
		CreatedAt:      entry.CreatedAt,
	}
	if entry.MediaItemID != nil {
		value := openapi_types.UUID(*entry.MediaItemID)
		response.MediaItemId = &value
	}
	response.SourcePath = entry.SourcePath
	response.DestinationPath = entry.DestinationPath
	return response
}
