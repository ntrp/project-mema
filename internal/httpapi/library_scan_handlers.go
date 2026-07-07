package httpapi

import (
	"context"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/text/unicode/norm"

	"media-manager/internal/library"
	"media-manager/internal/metadata"
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
	inputs, ok := s.libraryScanInputsForPath(w, ctx, folder)
	if !ok {
		return storage.LibraryScan{}, false
	}
	return s.storeLibraryScan(w, ctx, folder, inputs)
}

func (s *Server) libraryScanInputsForPath(w http.ResponseWriter, ctx context.Context, folder storage.LibraryFolder) ([]storage.LibraryScanItemInput, bool) {
	discovered, err := discoverLibraryFolder(folder)
	if err != nil {
		writeError(w, http.StatusBadRequest, "library_scan_failed", "Could not scan library folder")
		return nil, false
	}
	importedPaths, err := s.settings.ActiveImportedPathsForLibraryFolder(ctx, folder.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "library_scan_imported_lookup_failed", "Could not inspect imported files")
		return nil, false
	}
	inputs := s.libraryScanInputs(ctx, folder.Path, discovered, importedPaths)
	assignLibraryScanDuplicateGroups(inputs)
	return inputs, true
}

func discoverLibraryFolder(folder storage.LibraryFolder) ([]library.DiscoveredFile, error) {
	if folder.Kind == "movie" {
		return library.DiscoverMovies(folder.Path)
	}
	return library.Discover(folder.Path)
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

func (s *Server) libraryScanInputs(ctx context.Context, folderPath string, discovered []library.DiscoveredFile, importedPaths map[string]storage.ActiveImportedPath) []storage.LibraryScanItemInput {
	inputs := make([]storage.LibraryScanItemInput, 0, len(discovered))
	for _, item := range discovered {
		itemPath := libraryScanItemPath(folderPath, item.Path)
		input := storage.LibraryScanItemInput{
			Path:              itemPath,
			FileName:          item.FileName,
			SizeBytes:         item.SizeBytes,
			DetectedTitle:     item.DetectedTitle,
			DetectedYear:      item.DetectedYear,
			DetectedMediaKind: string(item.DetectedKind),
			SeasonNumber:      item.SeasonNumber,
			EpisodeNumber:     item.EpisodeNumber,
			SafeMatch:         item.SafeMatch,
		}
		if imported, ok := importedPaths[itemPath]; ok {
			applyImportedPathMatch(&input, imported)
		} else if imported, ok := importedPaths[item.Path]; ok {
			applyImportedPathMatch(&input, imported)
		}
		inputs = append(inputs, input)
	}
	return inputs
}

func applyImportedPathMatch(input *storage.LibraryScanItemInput, imported storage.ActiveImportedPath) {
	input.Imported = true
	input.MatchedTitle = &imported.MatchedTitle
	input.MatchedYear = imported.MatchedYear
	input.MatchedMediaKind = &input.DetectedMediaKind
	input.MatchSource = &imported.MatchedSource
	input.MediaItemID = &imported.MediaItemID
}

func libraryScanItemPath(folderPath string, discoveredPath string) string {
	if filepath.IsAbs(discoveredPath) {
		return filepath.Clean(discoveredPath)
	}
	return filepath.Clean(filepath.Join(folderPath, filepath.FromSlash(discoveredPath)))
}

func (s *Server) applyLibraryScanMatch(ctx context.Context, input *storage.LibraryScanItemInput) {
	mediaType, ok := scanInputMediaType(input.DetectedMediaKind)
	if !ok || strings.TrimSpace(input.DetectedTitle) == "" {
		return
	}
	if item, ok := s.localLibraryScanMatch(ctx, mediaType, input); ok {
		input.MatchedTitle = &item.Title
		input.MatchedYear = item.Year
		input.MatchedMediaKind = &input.DetectedMediaKind
		input.MediaItemID = &item.ID
		source := "library"
		input.MatchSource = &source
		return
	}
	provider, ok := s.defaultLibraryScanProvider(ctx, mediaType)
	if !ok {
		return
	}
	input.SelectedMetadataProviderID = &provider.ID
	results, err := s.searchMetadataProvider(ctx, provider, metadata.SearchRequest{
		Query:     input.DetectedTitle,
		MediaType: mediaType,
		Year:      input.DetectedYear,
	})
	if err != nil || len(results) == 0 {
		return
	}
	result, ok := confidentLibraryProviderMatch(input, results)
	if !ok {
		return
	}
	input.MatchedTitle = &result.Title
	input.MatchedYear = result.Year
	input.MatchedMediaKind = &input.DetectedMediaKind
	input.MatchedExternalProvider = optionalString(result.ExternalProvider)
	input.MatchedExternalID = optionalString(result.ExternalID)
	source := "provider"
	input.MatchSource = &source
}

func (s *Server) localLibraryScanMatch(ctx context.Context, mediaType string, input *storage.LibraryScanItemInput) (storage.MediaItem, bool) {
	items, err := s.settings.SearchMediaItems(ctx, input.DetectedTitle, &mediaType, 10)
	if err != nil {
		return storage.MediaItem{}, false
	}
	for _, item := range items {
		if !sameMediaTitle(item.Title, input.DetectedTitle) {
			continue
		}
		if mediaType == "movie" && input.DetectedYear != nil && item.Year != nil && *item.Year != *input.DetectedYear {
			continue
		}
		return item, true
	}
	return storage.MediaItem{}, false
}

func (s *Server) defaultLibraryScanProvider(ctx context.Context, mediaType string) (storage.MetadataProvider, bool) {
	providers, err := s.settings.ListEnabledMetadataProviders(ctx, mediaType)
	if err != nil || len(providers) == 0 {
		return storage.MetadataProvider{}, false
	}
	for _, provider := range providers {
		if strings.EqualFold(provider.Type, "tmdb") || strings.EqualFold(provider.Name, "tmdb") {
			return provider, true
		}
	}
	return providers[0], true
}

func confidentLibraryProviderMatch(input *storage.LibraryScanItemInput, results []metadata.SearchResult) (metadata.SearchResult, bool) {
	if len(results) == 1 {
		return results[0], true
	}
	for _, result := range results {
		if !sameMediaTitle(result.Title, input.DetectedTitle) {
			continue
		}
		if input.DetectedYear != nil && result.Year != nil && *result.Year != *input.DetectedYear {
			continue
		}
		return result, true
	}
	return metadata.SearchResult{}, false
}

func scanInputMediaType(kind string) (string, bool) {
	switch kind {
	case "movie", "anime_movie":
		return "movie", true
	case "series", "anime_series":
		return "serie", true
	default:
		return "", false
	}
}

func sameMediaTitle(left string, right string) bool {
	return normalizedMediaTitle(left) == normalizedMediaTitle(right)
}

func normalizedMediaTitle(value string) string {
	normalized := norm.NFD.String(strings.ToLower(value))
	var builder strings.Builder
	for _, r := range normalized {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func assignLibraryScanDuplicateGroups(inputs []storage.LibraryScanItemInput) {
	groups := map[string][]int{}
	for index, input := range inputs {
		key := libraryScanDuplicateKey(input)
		if key == "" {
			continue
		}
		groups[key] = append(groups[key], index)
	}
	for key, indexes := range groups {
		if len(indexes) < 2 {
			continue
		}
		groupID := "dup:" + key
		for _, index := range indexes {
			inputs[index].DuplicateGroupID = &groupID
			inputs[index].DuplicateRemovalAllowed = true
		}
	}
}

func libraryScanDuplicateKey(input storage.LibraryScanItemInput) string {
	if input.MediaItemID != nil {
		return "media:" + input.MediaItemID.String()
	}
	if input.MatchedExternalProvider != nil && input.MatchedExternalID != nil {
		key := strings.ToLower(*input.MatchedExternalProvider + ":" + *input.MatchedExternalID)
		if input.SeasonNumber != nil && input.EpisodeNumber != nil {
			key += ":s" + int32String(*input.SeasonNumber) + "e" + int32String(*input.EpisodeNumber)
		}
		return key
	}
	if input.DetectedTitle != "" {
		key := input.DetectedMediaKind + ":" + normalizedMediaTitle(input.DetectedTitle)
		if input.DetectedYear != nil {
			key += ":" + int32String(*input.DetectedYear)
		}
		if input.SeasonNumber != nil && input.EpisodeNumber != nil {
			key += ":s" + int32String(*input.SeasonNumber) + "e" + int32String(*input.EpisodeNumber)
		}
		return key
	}
	return ""
}

func int32String(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}
