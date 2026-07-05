package httpapi

import (
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func (s *Server) ListLibraryFolderOptions(w http.ResponseWriter, r *http.Request, params ListLibraryFolderOptionsParams) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	response, ok := libraryFolderOptionList(w, s.cfg.MediaDataDir, params.Path)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreateLibraryFolderOption(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body LibraryFolderOptionCreateRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	option, ok := createLibraryFolderOption(w, s.cfg.MediaDataDir, body)
	if !ok {
		return
	}
	writeJSON(w, http.StatusCreated, option)
}

func (s *Server) ListLibraryFolders(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	folders, err := s.settings.ListLibraryFolders(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "library_folder_list_failed", "Could not list library folders")
		return
	}
	response := LibraryFolderListResponse{Folders: make([]LibraryFolder, 0, len(folders))}
	for _, folder := range folders {
		response.Folders = append(response.Folders, libraryFolderResponse(folder))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreateLibraryFolder(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body LibraryFolderRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	path, kind, ok := libraryFolderInput(w, s.cfg.MediaDataDir, body)
	if !ok {
		return
	}
	folder, err := s.settings.CreateLibraryFolder(r.Context(), path, kind)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "library_folder_create_failed", "Could not create library folder")
		return
	}
	scan, ok := s.createLibraryScan(w, r.Context(), folder)
	if !ok {
		return
	}
	writeJSON(w, http.StatusCreated, LibraryFolderCreateResponse{
		Folder: libraryFolderResponse(folder),
		Scan:   libraryScanResponse(scan),
	})
}

func (s *Server) DeleteLibraryFolder(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	if err := s.settings.DeleteLibraryFolder(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete library folder")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) ListPathMappings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	mappings, err := s.settings.ListPathMappings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "path_mapping_list_failed", "Could not list path mappings")
		return
	}
	response := PathMappingListResponse{Mappings: make([]PathMapping, 0, len(mappings))}
	for _, mapping := range mappings {
		response.Mappings = append(response.Mappings, pathMappingResponse(mapping))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) CreatePathMapping(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body PathMappingRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	input, ok := pathMappingInput(w, body)
	if !ok {
		return
	}
	mapping, err := s.settings.CreatePathMapping(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "path_mapping_save_failed", "Could not save path mapping")
		return
	}
	writeJSON(w, http.StatusCreated, pathMappingResponse(mapping))
}

func (s *Server) DeletePathMapping(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	if err := s.settings.DeletePathMapping(r.Context(), uuid.UUID(id)); err != nil {
		writeSettingsError(w, err, "Could not delete path mapping")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) GetLibraryScan(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	scan, err := s.settings.GetLibraryScan(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find library scan")
		return
	}
	writeJSON(w, http.StatusOK, libraryScanResponse(scan))
}

func (s *Server) MatchLibraryScanItem(w http.ResponseWriter, r *http.Request, id ResourceId, itemId openapi_types.UUID) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body LibraryScanItemMatchRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	scan, err := s.settings.GetLibraryScan(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find library scan")
		return
	}
	input, ok := libraryMatchInput(w, body)
	if !ok {
		return
	}
	if !libraryFolderKindAllowsMediaKind(scan.FolderKind, input.MediaKind) {
		writeError(w, http.StatusBadRequest, "invalid_media_kind", "Matched media type does not belong in this library folder")
		return
	}
	item, mediaItem, err := s.settings.MatchLibraryScanItem(r.Context(), uuid.UUID(id), uuid.UUID(itemId), input)
	if err != nil {
		writeSettingsError(w, err, "Could not match library scan item")
		return
	}
	writeJSON(w, http.StatusOK, LibraryScanItemMatchResponse{
		Item:      libraryScanItemResponse(item),
		MediaItem: mediaItemResponse(mediaItem),
	})
}

func (s *Server) ImportLibraryScanItems(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	var body LibraryScanImportRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	if len(body.Items) == 0 {
		writeError(w, http.StatusBadRequest, "invalid_import_rows", "At least one import row is required")
		return
	}
	scanID := uuid.UUID(id)
	scan, err := s.settings.GetLibraryScan(r.Context(), scanID)
	if err != nil {
		writeSettingsError(w, err, "Could not find library scan")
		return
	}
	mediaItems := make([]MediaItem, 0, len(body.Items))
	storageItems := make([]storage.MediaItem, 0, len(body.Items))
	for _, row := range body.Items {
		input, ok := libraryMatchInput(w, row.Match)
		if !ok {
			return
		}
		if !libraryFolderKindAllowsMediaKind(scan.FolderKind, input.MediaKind) {
			writeError(w, http.StatusBadRequest, "invalid_media_kind", "Matched media type does not belong in this library folder")
			return
		}
		input, err = s.enrichLibraryImportMatch(r.Context(), input)
		if err != nil {
			writeMetadataDetailsError(w, err)
			return
		}
		_, mediaItem, err := s.settings.ImportLibraryScanItem(r.Context(), scanID, uuid.UUID(row.ItemId), input)
		if err != nil {
			writeSettingsError(w, err, "Could not import library scan item")
			return
		}
		storageItems = append(storageItems, mediaItem)
		mediaItems = append(mediaItems, mediaItemResponse(mediaItem))
	}
	removed := int32(0)
	if body.RemoveDuplicatePaths != nil && len(*body.RemoveDuplicatePaths) > 0 && len(storageItems) > 0 {
		mediaItemID := storageItems[0].ID
		for _, path := range *body.RemoveDuplicatePaths {
			if err := s.settings.DeleteLibraryFolderFileForMedia(r.Context(), mediaItemID, scan.FolderID, path); err != nil {
				writeSettingsError(w, err, "Could not remove duplicate library file")
				return
			}
			removed++
		}
	}
	refreshed, err := s.settings.GetLibraryScan(r.Context(), scanID)
	if err != nil {
		writeSettingsError(w, err, "Could not reload library scan")
		return
	}
	writeJSON(w, http.StatusOK, LibraryScanImportResponse{
		Scan:                  libraryScanResponse(refreshed),
		MediaItems:            mediaItems,
		ImportedCount:         int32(len(body.Items)),
		RemovedDuplicateCount: removed,
	})
}
