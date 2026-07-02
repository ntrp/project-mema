package httpapi

import "os"

func mediaFileInfoResponses(paths []string) *[]MediaFileInfo {
	files := make([]MediaFileInfo, 0, len(paths))
	for _, path := range paths {
		file := MediaFileInfo{Path: path}
		if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
			size := stat.Size()
			file.SizeBytes = &size
			probe := mediaFileProbe(path)
			if len(probe.tracks) > 0 {
				file.Tracks = &probe.tracks
			}
			if len(probe.chapters) > 0 {
				file.Chapters = &probe.chapters
			}
		}
		files = append(files, file)
	}
	return &files
}
