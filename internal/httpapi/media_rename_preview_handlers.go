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

func (s *Server) ApplyMediaRename(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	var body MediaRenameApplyRequest
	if r.Body != nil && r.ContentLength != 0 {
		if !decodeJSON(w, r, &body) {
			return
		}
	}
	result, err := s.settings.ApplySelectedMediaItemRename(
		r.Context(),
		uuid.UUID(id),
		body.CurrentPaths,
	)
	if err != nil {
		writeSettingsError(w, err, "Could not apply file rename")
		return
	}
	writeJSON(w, http.StatusOK, mediaRenameApplyResponse(result))
}

func mediaRenamePreviewResponse(preview storage.MediaRenamePreview) MediaRenamePreviewResponse {
	response := MediaRenamePreviewResponse{Rows: make([]MediaRenamePreviewRow, 0, len(preview.Rows))}
	for _, row := range preview.Rows {
		response.Rows = append(response.Rows, mediaRenamePreviewRow(row))
	}
	return response
}

func mediaRenameApplyResponse(result storage.MediaRenameApplyResult) MediaRenameApplyResponse {
	response := MediaRenameApplyResponse{
		Rows:         make([]MediaRenamePreviewRow, 0, len(result.Rows)),
		AppliedCount: result.AppliedCount,
		SkippedCount: result.SkippedCount,
		FailedCount:  result.FailedCount,
	}
	for _, row := range result.Rows {
		response.Rows = append(response.Rows, mediaRenamePreviewRow(row))
	}
	return response
}

func mediaRenamePreviewRow(row storage.MediaRenamePreviewRow) MediaRenamePreviewRow {
	return MediaRenamePreviewRow{
		CurrentPath:  row.CurrentPath,
		ProposedPath: row.ProposedPath,
		Status:       MediaRenamePreviewRowStatus(row.Status),
		Messages:     row.Messages,
	}
}
