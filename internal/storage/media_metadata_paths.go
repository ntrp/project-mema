package storage

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var metadataFileExtensions = map[string]struct{}{
	".ass":  {},
	".idx":  {},
	".jpeg": {},
	".jpg":  {},
	".nfo":  {},
	".png":  {},
	".srt":  {},
	".ssa":  {},
	".sub":  {},
	".tbn":  {},
	".txt":  {},
	".webp": {},
}

func collectMetadataFilePaths(mediaPaths []string) []string {
	paths := map[string]struct{}{}
	for _, mediaPath := range mediaPaths {
		dir := filepath.Dir(mediaPath)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		mediaBase := strings.TrimSuffix(filepath.Base(mediaPath), filepath.Ext(mediaPath))
		mediaBase = strings.ToLower(mediaBase)
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			ext := strings.ToLower(filepath.Ext(name))
			if _, ok := metadataFileExtensions[ext]; !ok {
				continue
			}
			base := strings.ToLower(strings.TrimSuffix(name, filepath.Ext(name)))
			if !isRelatedMetadataBase(base, mediaBase) {
				continue
			}
			fullPath := filepath.Join(dir, name)
			if fullPath != mediaPath {
				paths[fullPath] = struct{}{}
			}
		}
	}

	result := make([]string, 0, len(paths))
	for path := range paths {
		result = append(result, path)
	}
	sort.Strings(result)
	return result
}

func isRelatedMetadataBase(base string, mediaBase string) bool {
	if base == mediaBase || strings.HasPrefix(base, mediaBase+".") || strings.HasPrefix(base, mediaBase+"-") {
		return true
	}
	switch base {
	case "banner", "clearlogo", "cover", "fanart", "folder", "landscape", "movie", "poster":
		return true
	default:
		return false
	}
}
