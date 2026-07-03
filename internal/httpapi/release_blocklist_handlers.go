package httpapi

import (
	"net/http"

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

func releaseBlocklistItemResponse(item storage.ReleaseBlocklistItem) ReleaseBlocklistItem {
	return ReleaseBlocklistItem{
		Id:           item.ID,
		MediaItemId:  item.MediaItemID,
		MediaTitle:   item.MediaTitle,
		MediaType:    MediaType(item.MediaType),
		ReleaseTitle: item.ReleaseTitle,
		IndexerName:  item.IndexerName,
		Reason:       item.Reason,
		Source:       item.Source,
		Temporary:    item.Temporary,
		ExpiresAt:    item.ExpiresAt,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	}
}
