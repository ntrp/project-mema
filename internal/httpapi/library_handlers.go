package httpapi

import (
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
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
	path, ok := libraryFolderInput(w, s.cfg.MediaDataDir, body)
	if !ok {
		return
	}
	inputs, ok := libraryScanInputsForPath(w, path)
	if !ok {
		return
	}
	folder, err := s.settings.CreateLibraryFolder(r.Context(), path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "library_folder_create_failed", "Could not create library folder")
		return
	}
	scan, ok := s.storeLibraryScan(w, r.Context(), folder, inputs)
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
	input, ok := libraryMatchInput(w, body)
	if !ok {
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
