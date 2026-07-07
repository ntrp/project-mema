package httpapi

import (
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *Server) ResetLibraryScanItemImport(w http.ResponseWriter, r *http.Request, id ResourceId, itemId openapi_types.UUID) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	scanID := uuid.UUID(id)
	result, err := s.settings.ResetLibraryScanItemImport(r.Context(), scanID, uuid.UUID(itemId))
	if err != nil {
		writeSettingsError(w, err, "Could not reset library scan item")
		return
	}
	scan, err := s.settings.GetLibraryScan(r.Context(), scanID)
	if err != nil {
		writeSettingsError(w, err, "Could not reload library scan")
		return
	}
	writeJSON(w, http.StatusOK, LibraryScanItemResetResponse{
		Scan:               libraryScanResponse(scan),
		Item:               libraryScanItemResponse(result.Item),
		RemovedMediaItemId: result.RemovedMediaItemID,
	})
}
