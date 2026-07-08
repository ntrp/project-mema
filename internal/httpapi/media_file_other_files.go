package httpapi

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"media-manager/internal/storage"
)

var subtitleFileExtensions = map[string]struct{}{
	".ass": {},
	".idx": {},
	".srt": {},
	".ssa": {},
	".sub": {},
}

var metadataFileExtensions = map[string]struct{}{
	".jpeg": {},
	".jpg":  {},
	".nfo":  {},
	".png":  {},
	".tbn":  {},
	".webp": {},
}

func mediaFileOtherFiles(
	path string,
	mediaPaths []string,
	subtitleTargets []storage.MediaProfileSubtitleTarget,
	subtitleMode string,
	externalSubtitles []storage.MediaItemSubtitle,
	sidecars []storage.MediaItemSidecar,
	satisfaction *MediaFileSubtitleSatisfaction,
) []MediaFileOtherFile {
	seen := map[string]MediaFileOtherFile{}
	for _, file := range availableOtherFiles(path, mediaPaths) {
		seen[file.pathKey()] = file
	}
	for _, file := range storedOtherFiles(path, sidecars) {
		seen[file.pathKey()] = file
	}
	for _, subtitle := range externalSubtitles {
		if !sameSubtitleMediaBase(subtitle.FilePath, path) {
			continue
		}
		file := subtitleOtherFile(subtitle.FilePath, subtitle.LanguageID, otherFileStatus(subtitle.FilePath))
		seen[file.pathKey()] = file
	}
	if mediaFileSubtitleMode(subtitleMode) == MediaProfileSubtitleModeExternal {
		for _, file := range missingExternalSubtitleFiles(path, subtitleTargets, satisfaction) {
			seen[file.pathKey()] = file
		}
	}
	files := make([]MediaFileOtherFile, 0, len(seen))
	for _, file := range seen {
		files = append(files, file)
	}
	sort.Slice(files, func(i, j int) bool {
		if files[i].Status != files[j].Status {
			return files[i].Status == MediaFileOtherFileStatusMissing
		}
		if files[i].Type != files[j].Type {
			return files[i].Type < files[j].Type
		}
		return files[i].Path < files[j].Path
	})
	return files
}

func availableOtherFiles(path string, mediaPaths []string) []MediaFileOtherFile {
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		return nil
	}
	media := mediaPathSet(mediaPaths)
	files := []MediaFileOtherFile{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		candidate := filepath.Join(filepath.Dir(path), entry.Name())
		if _, ok := media[candidate]; ok {
			continue
		}
		files = append(files, availableOtherFile(path, candidate))
	}
	return files
}

func availableOtherFile(mediaPath string, path string) MediaFileOtherFile {
	sidecar := storage.ClassifyMediaSidecar(mediaPath, path)
	fileType := MediaFileOtherFileTypeUnknown
	if sidecar.Type == storage.MediaSidecarSubtitle {
		fileType = MediaFileOtherFileTypeSubtitle
	} else if sidecar.Type == storage.MediaSidecarMetadata || metadataFile(path) {
		fileType = MediaFileOtherFileTypeMetadata
	}
	return MediaFileOtherFile{
		Type:     fileType,
		Path:     path,
		Status:   MediaFileOtherFileStatusAvailable,
		Subtype:  optionalString(sidecar.Subtype),
		Language: optionalString(sidecar.LanguageID),
	}
}

func storedOtherFiles(path string, sidecars []storage.MediaItemSidecar) []MediaFileOtherFile {
	files := []MediaFileOtherFile{}
	for _, sidecar := range sidecars {
		if sidecar.MediaFilePath != path {
			continue
		}
		fileType := MediaFileOtherFileTypeUnknown
		switch sidecar.SidecarType {
		case storage.MediaSidecarSubtitle:
			fileType = MediaFileOtherFileTypeSubtitle
		case storage.MediaSidecarMetadata:
			fileType = MediaFileOtherFileTypeMetadata
		}
		files = append(files, MediaFileOtherFile{
			Type:     fileType,
			Path:     sidecar.FilePath,
			Status:   otherFileStatus(sidecar.FilePath),
			Subtype:  sidecar.Subtype,
			Language: sidecar.LanguageID,
		})
	}
	return files
}

func subtitleOtherFile(
	path string,
	language string,
	status MediaFileOtherFileStatus,
) MediaFileOtherFile {
	return MediaFileOtherFile{
		Type:     MediaFileOtherFileTypeSubtitle,
		Path:     path,
		Status:   status,
		Subtype:  optionalString(subtitleSubtype(path)),
		Language: optionalString(language),
	}
}

func subtitleSubtype(path string) string {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	if ext == "srt" || ext == "subrip" {
		return "subrip"
	}
	return ext
}

func missingExternalSubtitleFiles(
	path string,
	targets []storage.MediaProfileSubtitleTarget,
	satisfaction *MediaFileSubtitleSatisfaction,
) []MediaFileOtherFile {
	if satisfaction == nil || len(satisfaction.MissingLanguages) == 0 {
		return nil
	}
	missing := languageSet(satisfaction.MissingLanguages)
	files := []MediaFileOtherFile{}
	for _, target := range targets {
		language := languageMatchKey(target.LanguageID)
		if _, ok := missing[language]; !ok {
			continue
		}
		files = append(files, subtitleOtherFile(
			expectedSubtitlePath(path, target),
			target.LanguageID,
			MediaFileOtherFileStatusMissing,
		))
	}
	return files
}

func expectedSubtitlePath(path string, target storage.MediaProfileSubtitleTarget) string {
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	format := "srt"
	if len(target.Formats) > 0 && strings.TrimSpace(target.Formats[0]) != "" {
		format = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(target.Formats[0])), ".")
	}
	language := languageMatchKey(target.LanguageID)
	if language == "" {
		language = strings.ToLower(strings.TrimSpace(target.LanguageID))
	}
	return filepath.Join(filepath.Dir(path), base+"."+language+"."+format)
}

func metadataFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if _, ok := metadataFileExtensions[ext]; !ok {
		return false
	}
	base := strings.ToLower(strings.TrimSuffix(filepath.Base(path), ext))
	switch base {
	case "banner", "clearlogo", "cover", "fanart", "folder", "landscape", "movie", "poster":
		return true
	default:
		return ext == ".nfo"
	}
}

func mediaPathSet(paths []string) map[string]struct{} {
	result := map[string]struct{}{}
	for _, path := range paths {
		result[path] = struct{}{}
	}
	return result
}

func otherFileStatus(path string) MediaFileOtherFileStatus {
	if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
		return MediaFileOtherFileStatusAvailable
	}
	return MediaFileOtherFileStatusMissing
}

func (file MediaFileOtherFile) pathKey() string {
	return strings.ToLower(file.Path)
}
