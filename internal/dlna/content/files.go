package content

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"media-manager/internal/storage"
)

var episodeFilePattern = regexp.MustCompile(`(?i)s(\d{1,2})e(\d{1,3})`)

func (t *Tree) availableFiles(item storage.MediaItem) []File {
	files := make([]File, 0, len(item.FilePaths))
	for _, path := range item.FilePaths {
		path = strings.TrimSpace(path)
		if path == "" || !t.fileAvailable(path) {
			continue
		}
		file := File{Path: path, Hash: filePathHash(path)}
		if season, episode, ok := episodeNumbers(path); ok {
			file.Season = season
			file.Episode = episode
			file.HasNumber = true
		}
		files = append(files, file)
	}
	return files
}

func (t *Tree) fileAvailable(path string) bool {
	if t.stat == nil {
		return false
	}
	info, err := t.stat(path)
	return err == nil && !info.IsDir()
}

func episodeNumbers(path string) (int32, int32, bool) {
	match := episodeFilePattern.FindStringSubmatch(filepath.Base(path))
	if len(match) != 3 {
		return 0, 0, false
	}
	season, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, 0, false
	}
	episode, err := strconv.Atoi(match[2])
	if err != nil {
		return 0, 0, false
	}
	return int32(season), int32(episode), true
}

func filesForEpisode(files []File, seasonNumber int32, episodeNumber int32) []File {
	matched := make([]File, 0, len(files))
	for _, file := range files {
		if file.HasNumber && file.Season == seasonNumber && file.Episode == episodeNumber {
			matched = append(matched, file)
		}
	}
	return matched
}
