package httpapi

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/storage"
)

func libraryFolderInput(w http.ResponseWriter, mediaDataDir string, request LibraryFolderRequest) (string, string, bool) {
	path := strings.TrimSpace(request.Path)
	if path == "" {
		writeError(w, http.StatusBadRequest, "invalid_path", "Library folder path is required")
		return "", "", false
	}
	kind := string(request.Kind)
	if kind != "movie" && kind != "series" {
		writeError(w, http.StatusBadRequest, "invalid_folder_kind", "Library folder type must be movie or series")
		return "", "", false
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(mediaDataDir, path)
	}
	cleaned, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_path", "Library folder path is invalid")
		return "", "", false
	}
	info, err := os.Stat(cleaned)
	if err != nil {
		writeError(w, http.StatusBadRequest, "path_unavailable", "Library folder path is not available")
		return "", "", false
	}
	if !info.IsDir() {
		writeError(w, http.StatusBadRequest, "path_not_directory", "Library folder path must be a directory")
		return "", "", false
	}
	return cleaned, kind, true
}

func libraryFolderOptionList(w http.ResponseWriter, mediaDataDir string, requestedPath *string) (LibraryFolderOptionListResponse, bool) {
	path := strings.TrimSpace(mediaDataDir)
	relativeToMediaDir := false
	if requestedPath != nil && strings.TrimSpace(*requestedPath) != "" {
		path = strings.TrimSpace(*requestedPath)
		relativeToMediaDir = true
	}
	if path == "" {
		writeError(w, http.StatusBadRequest, "invalid_path", "Folder path is required")
		return LibraryFolderOptionListResponse{}, false
	}
	if relativeToMediaDir && !filepath.IsAbs(path) {
		path = filepath.Join(mediaDataDir, path)
	}
	cleaned, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_path", "Folder path is invalid")
		return LibraryFolderOptionListResponse{}, false
	}
	info, err := os.Stat(cleaned)
	if err != nil {
		writeError(w, http.StatusBadRequest, "path_unavailable", "Folder path is not available")
		return LibraryFolderOptionListResponse{}, false
	}
	if !info.IsDir() {
		writeError(w, http.StatusBadRequest, "path_not_directory", "Folder path must be a directory")
		return LibraryFolderOptionListResponse{}, false
	}
	entries, err := os.ReadDir(cleaned)
	if err != nil {
		writeError(w, http.StatusBadRequest, "path_unreadable", "Folder path is not readable")
		return LibraryFolderOptionListResponse{}, false
	}

	response := LibraryFolderOptionListResponse{
		CurrentPath: cleaned,
		Entries:     []LibraryFolderOption{},
	}
	if parent := filepath.Dir(cleaned); parent != cleaned {
		response.ParentPath = &parent
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		response.Entries = append(response.Entries, LibraryFolderOption{
			Name: entry.Name(),
			Path: filepath.Join(cleaned, entry.Name()),
		})
	}
	sort.Slice(response.Entries, func(i, j int) bool {
		return strings.ToLower(response.Entries[i].Name) < strings.ToLower(response.Entries[j].Name)
	})
	return response, true
}

func createLibraryFolderOption(w http.ResponseWriter, mediaDataDir string, request LibraryFolderOptionCreateRequest) (LibraryFolderOption, bool) {
	parentPath, ok := cleanLibraryFolderPath(w, mediaDataDir, request.ParentPath, true)
	if !ok {
		return LibraryFolderOption{}, false
	}
	info, err := os.Stat(parentPath)
	if err != nil {
		writeError(w, http.StatusBadRequest, "parent_path_unavailable", "Parent folder is not available")
		return LibraryFolderOption{}, false
	}
	if !info.IsDir() {
		writeError(w, http.StatusBadRequest, "parent_path_not_directory", "Parent path must be a directory")
		return LibraryFolderOption{}, false
	}

	name := strings.TrimSpace(request.Name)
	if name == "" {
		writeError(w, http.StatusBadRequest, "invalid_folder_name", "Folder name is required")
		return LibraryFolderOption{}, false
	}
	if len(name) > 255 {
		writeError(w, http.StatusBadRequest, "invalid_folder_name", "Folder name must be 255 characters or fewer")
		return LibraryFolderOption{}, false
	}
	if name == "." || name == ".." || name != filepath.Base(name) || strings.ContainsAny(name, `/\`) {
		writeError(w, http.StatusBadRequest, "invalid_folder_name", "Folder name must not contain path separators")
		return LibraryFolderOption{}, false
	}

	target := filepath.Join(parentPath, name)
	if err := os.Mkdir(target, 0o755); err != nil {
		if os.IsExist(err) {
			existing, statErr := os.Stat(target)
			if statErr == nil && existing.IsDir() {
				return LibraryFolderOption{Name: name, Path: target}, true
			}
			writeError(w, http.StatusBadRequest, "folder_exists_as_file", "A file already exists with that name")
			return LibraryFolderOption{}, false
		}
		writeError(w, http.StatusBadRequest, "folder_create_failed", "Could not create folder")
		return LibraryFolderOption{}, false
	}
	return LibraryFolderOption{Name: name, Path: target}, true
}

func cleanLibraryFolderPath(w http.ResponseWriter, mediaDataDir string, path string, relativeToMediaDir bool) (string, bool) {
	path = strings.TrimSpace(path)
	if path == "" {
		writeError(w, http.StatusBadRequest, "invalid_path", "Folder path is required")
		return "", false
	}
	if relativeToMediaDir && !filepath.IsAbs(path) {
		path = filepath.Join(mediaDataDir, path)
	}
	cleaned, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_path", "Folder path is invalid")
		return "", false
	}
	return cleaned, true
}

func libraryMatchInput(w http.ResponseWriter, request LibraryScanItemMatchRequest) (storage.LibraryMatchInput, bool) {
	title := strings.TrimSpace(request.Title)
	if title == "" {
		writeError(w, http.StatusBadRequest, "invalid_title", "Matched title is required")
		return storage.LibraryMatchInput{}, false
	}
	if request.MediaKind == LibraryMediaKindUnknown || !request.MediaKind.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_media_kind", "A movie, series, or anime type is required")
		return storage.LibraryMatchInput{}, false
	}
	qualityProfileID := strings.TrimSpace(request.QualityProfileId)
	if qualityProfileID == "" {
		writeError(w, http.StatusBadRequest, "invalid_quality_profile", "Quality profile is required")
		return storage.LibraryMatchInput{}, false
	}
	if !request.MonitorMode.Valid() || !request.MinimumAvailability.Valid() {
		writeError(w, http.StatusBadRequest, "invalid_monitoring_options", "Monitor and availability options are required")
		return storage.LibraryMatchInput{}, false
	}
	return storage.LibraryMatchInput{
		MediaKind:           string(request.MediaKind),
		Title:               title,
		Year:                request.Year,
		Monitored:           request.Monitored,
		QualityProfileID:    qualityProfileID,
		MonitorMode:         string(request.MonitorMode),
		MinimumAvailability: string(request.MinimumAvailability),
		SeriesType:          optionalLibrarySeriesType(request.SeriesType),
		MetadataProviderID:  optionalUUID(request.MetadataProviderId),
		MediaItemID:         optionalUUID(request.MediaItemId),
		ExternalProvider:    optionalTrimmedString(request.ExternalProvider),
		ExternalID:          optionalTrimmedString(request.ExternalId),
		Overview:            optionalTrimmedString(request.Overview),
		PosterPath:          optionalTrimmedString(request.PosterPath),
	}, true
}

func libraryFolderResponse(folder storage.LibraryFolder) LibraryFolder {
	return LibraryFolder{
		Id:        openapi_types.UUID(folder.ID),
		Path:      folder.Path,
		Kind:      LibraryFolderKind(folder.Kind),
		CreatedAt: folder.CreatedAt,
		UpdatedAt: folder.UpdatedAt,
	}
}

func pathMappingInput(w http.ResponseWriter, request PathMappingRequest) (storage.PathMappingInput, bool) {
	clientPath := strings.TrimSpace(request.ClientPath)
	appPath := strings.TrimSpace(request.AppPath)
	if clientPath == "" || appPath == "" {
		writeError(w, http.StatusBadRequest, "invalid_path_mapping", "Client and app paths are required")
		return storage.PathMappingInput{}, false
	}
	return storage.PathMappingInput{ClientPath: clientPath, AppPath: appPath}, true
}

func pathMappingResponse(mapping storage.PathMapping) PathMapping {
	return PathMapping{
		Id:         openapi_types.UUID(mapping.ID),
		ClientPath: mapping.ClientPath,
		AppPath:    mapping.AppPath,
		CreatedAt:  mapping.CreatedAt,
		UpdatedAt:  mapping.UpdatedAt,
	}
}

func libraryScanResponse(scan storage.LibraryScan) LibraryScan {
	items := make([]LibraryScanItem, 0, len(scan.Items))
	for _, item := range scan.Items {
		items = append(items, libraryScanItemResponse(item))
	}
	return LibraryScan{
		Id:               openapi_types.UUID(scan.ID),
		FolderId:         openapi_types.UUID(scan.FolderID),
		FolderPath:       scan.FolderPath,
		FolderKind:       LibraryFolderKind(scan.FolderKind),
		Status:           LibraryScanStatus(scan.Status),
		TotalFiles:       scan.TotalFiles,
		AutoMatchedCount: scan.AutoMatchedCount,
		ManualCount:      scan.ManualCount,
		Items:            items,
		CreatedAt:        scan.CreatedAt,
		CompletedAt:      scan.CompletedAt,
	}
}

func libraryScanItemResponse(item storage.LibraryScanItem) LibraryScanItem {
	var mediaItemID *openapi_types.UUID
	if item.MediaItemID != nil {
		value := openapi_types.UUID(*item.MediaItemID)
		mediaItemID = &value
	}
	var matchedKind *LibraryMediaKind
	if item.MatchedMediaKind != nil {
		value := LibraryMediaKind(*item.MatchedMediaKind)
		matchedKind = &value
	}
	var selectedProviderID *openapi_types.UUID
	if item.SelectedMetadataProviderID != nil {
		value := openapi_types.UUID(*item.SelectedMetadataProviderID)
		selectedProviderID = &value
	}
	return LibraryScanItem{
		Id:                         openapi_types.UUID(item.ID),
		ScanId:                     openapi_types.UUID(item.ScanID),
		Path:                       item.Path,
		FileName:                   item.FileName,
		SizeBytes:                  optionalInt64(item.SizeBytes),
		DetectedTitle:              item.DetectedTitle,
		DetectedYear:               item.DetectedYear,
		DetectedMediaKind:          LibraryMediaKind(item.DetectedMediaKind),
		SeasonNumber:               item.SeasonNumber,
		EpisodeNumber:              item.EpisodeNumber,
		Status:                     LibraryScanItemStatus(item.Status),
		Imported:                   item.Imported,
		MatchedTitle:               item.MatchedTitle,
		MatchedYear:                item.MatchedYear,
		MatchedMediaKind:           matchedKind,
		MatchedExternalProvider:    item.MatchedExternalProvider,
		MatchedExternalId:          item.MatchedExternalID,
		MatchSource:                libraryMatchSourceResponse(item.MatchSource),
		SelectedMetadataProviderId: selectedProviderID,
		DuplicateGroupId:           item.DuplicateGroupID,
		DuplicateRemovalAllowed:    item.DuplicateRemovalAllowed,
		MediaItemId:                mediaItemID,
		CreatedAt:                  item.CreatedAt,
		UpdatedAt:                  item.UpdatedAt,
	}
}
