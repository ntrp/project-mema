package episodepath

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	seasonEpisodePattern = regexp.MustCompile(`(?i)s(\d{1,3})e(\d{1,3})`)
	seasonFolderPattern  = regexp.MustCompile(`(?i)^(?:season|s)[ ._-]*(\d{1,3})$`)
	episodeFilePattern   = regexp.MustCompile(`(?i)^(?:episode|ep|e)?[ ._-]*0*(\d{1,3})(?:\D|$)`)
)

// Numbers resolves conventional series paths such as Season 2/03.mkv and S02E03.mkv.
func Numbers(path string) (int32, int32, bool) {
	if match := seasonEpisodePattern.FindStringSubmatch(filepath.Base(path)); len(match) == 3 {
		return numbers(match[1], match[2])
	}

	season, ok := seasonNumber(filepath.Dir(path))
	if !ok {
		return 0, 0, false
	}
	match := episodeFilePattern.FindStringSubmatch(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
	if len(match) != 2 {
		return 0, 0, false
	}
	episode, err := strconv.ParseInt(match[1], 10, 32)
	if err != nil {
		return 0, 0, false
	}
	return season, int32(episode), true
}

func seasonNumber(path string) (int32, bool) {
	for path != "." && path != string(filepath.Separator) {
		name := filepath.Base(path)
		if strings.EqualFold(name, "specials") {
			return 0, true
		}
		if match := seasonFolderPattern.FindStringSubmatch(name); len(match) == 2 {
			value, err := strconv.ParseInt(match[1], 10, 32)
			return int32(value), err == nil
		}
		parent := filepath.Dir(path)
		if parent == path {
			break
		}
		path = parent
	}
	return 0, false
}

func numbers(seasonValue string, episodeValue string) (int32, int32, bool) {
	season, seasonErr := strconv.ParseInt(seasonValue, 10, 32)
	episode, episodeErr := strconv.ParseInt(episodeValue, 10, 32)
	return int32(season), int32(episode), seasonErr == nil && episodeErr == nil
}
