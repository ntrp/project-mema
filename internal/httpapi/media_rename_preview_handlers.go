package httpapi

import (
	"net/http"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func (s *Server) PreviewMediaRename(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	preview, err := s.settings.PreviewMediaItemRename(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not preview file rename")
		return
	}
	writeJSON(w, http.StatusOK, mediaRenamePreviewResponse(preview))
}

func mediaRenamePreviewResponse(preview storage.MediaRenamePreview) MediaRenamePreviewResponse {
	response := MediaRenamePreviewResponse{Rows: make([]MediaRenamePreviewRow, 0, len(preview.Rows))}
	for _, row := range preview.Rows {
		response.Rows = append(response.Rows, MediaRenamePreviewRow{
			CurrentPath:  row.CurrentPath,
			ProposedPath: row.ProposedPath,
			Status:       MediaRenamePreviewRowStatus(row.Status),
			Messages:     row.Messages,
		})
	}
	return response
}
