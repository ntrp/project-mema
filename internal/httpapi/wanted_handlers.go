package httpapi

import (
	"net/http"

	"media-manager/internal/satisfaction"
	"media-manager/internal/storage"
)

func (s *Server) ListWantedRows(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	items, err := s.settings.ListMissingMediaItems(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "wanted_rows_failed", "Could not list wanted rows")
		return
	}
	writeJSON(w, http.StatusOK, wantedRowsResponse(items))
}

func wantedRowsResponse(items []storage.MediaItem) WantedRowListResponse {
	response := WantedRowListResponse{Rows: []WantedRow{}}
	for _, item := range items {
		rows := satisfaction.BuildWantedRows(satisfaction.WantedRowsInput{
			Item:           item,
			HasUsableMedia: false,
		})
		for _, row := range rows {
			response.Rows = append(response.Rows, wantedRowResponse(row))
		}
	}
	return response
}
