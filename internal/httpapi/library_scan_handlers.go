package httpapi

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"media-manager/internal/library"
	"media-manager/internal/storage"
)

func (s *Server) ScanLibraryFolder(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	folder, err := s.settings.GetLibraryFolder(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find library folder")
		return
	}
	scan, ok := s.createLibraryScan(w, r.Context(), folder)
	if !ok {
		return
	}
	writeJSON(w, http.StatusCreated, libraryScanResponse(scan))
}

func (s *Server) createLibraryScan(w http.ResponseWriter, ctx context.Context, folder storage.LibraryFolder) (storage.LibraryScan, bool) {
	inputs, ok := libraryScanInputsForPath(w, folder.Path)
	if !ok {
		return storage.LibraryScan{}, false
	}
	return s.storeLibraryScan(w, ctx, folder, inputs)
}

func libraryScanInputsForPath(w http.ResponseWriter, path string) ([]storage.LibraryScanItemInput, bool) {
	discovered, err := library.Discover(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "library_scan_failed", "Could not scan library folder")
		return nil, false
	}
	return libraryScanInputs(discovered), true
}

func (s *Server) storeLibraryScan(
	w http.ResponseWriter,
	ctx context.Context,
	folder storage.LibraryFolder,
	inputs []storage.LibraryScanItemInput,
) (storage.LibraryScan, bool) {
	scan, err := s.settings.CreateLibraryScan(ctx, folder, inputs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "library_scan_store_failed", "Could not store library scan")
		return storage.LibraryScan{}, false
	}
	return scan, true
}

func libraryScanInputs(discovered []library.DiscoveredFile) []storage.LibraryScanItemInput {
	inputs := make([]storage.LibraryScanItemInput, 0, len(discovered))
	for _, item := range discovered {
		inputs = append(inputs, storage.LibraryScanItemInput{
			Path:              item.Path,
			FileName:          item.FileName,
			DetectedTitle:     item.DetectedTitle,
			DetectedYear:      item.DetectedYear,
			DetectedMediaKind: string(item.DetectedKind),
			SafeMatch:         item.SafeMatch,
		})
	}
	return inputs
}
