package httpapi

import (
	"net/http"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func (s *Server) ListReleaseBlocklist(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	items, err := s.settings.ListReleaseBlocklist(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "release_blocklist_failed", "Could not list blocked releases")
		return
	}
	response := ReleaseBlocklistListResponse{Items: make([]ReleaseBlocklistItem, 0, len(items))}
	for _, item := range items {
		response.Items = append(response.Items, releaseBlocklistItemResponse(item))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) DeleteReleaseBlocklistItem(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	if err := s.settings.DeleteReleaseBlocklistItem(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete release blocklist entry")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "activity", "Release blocklist entry deleted", map[string]any{"blocklistId": uuid.UUID(id).String()})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) ClearReleaseBlocklist(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	count, err := s.settings.ClearReleaseBlocklist(r.Context())
	if err != nil {
		writeSettingsError(w, err, "Could not clear release blocklist")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "activity", "Release blocklist cleared", map[string]any{"count": count})
	w.WriteHeader(http.StatusNoContent)
}

func releaseBlocklistItemResponse(item storage.ReleaseBlocklistItem) ReleaseBlocklistItem {
	return ReleaseBlocklistItem{
		Id:                 item.ID,
		MediaItemId:        item.MediaItemID,
		MediaTitle:         item.MediaTitle,
		MediaType:          MediaType(item.MediaType),
		ReleaseTitle:       item.ReleaseTitle,
		IndexerName:        item.IndexerName,
		IndexerProtocol:    IndexerProtocol(item.IndexerProtocol),
		DownloadClientName: item.DownloadClientName,
		Reason:             item.Reason,
		Source:             item.Source,
		Temporary:          item.Temporary,
		ExpiresAt:          item.ExpiresAt,
		CreatedAt:          item.CreatedAt,
		UpdatedAt:          item.UpdatedAt,
	}
}
