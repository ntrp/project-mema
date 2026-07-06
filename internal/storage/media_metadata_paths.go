package storage

import "sort"

func collectMetadataFilePaths(mediaPaths []string) []string {
	paths := map[string]struct{}{}
	for _, mediaPath := range mediaPaths {
		for _, sidecar := range MediaSidecarsForFile(mediaPath) {
			if sidecar.Type == MediaSidecarMetadata {
				paths[sidecar.Path] = struct{}{}
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
