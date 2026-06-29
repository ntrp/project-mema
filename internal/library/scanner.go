package library

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type MediaKind string

const (
	MediaKindMovie       MediaKind = "movie"
	MediaKindSeries      MediaKind = "series"
	MediaKindAnimeMovie  MediaKind = "anime_movie"
	MediaKindAnimeSeries MediaKind = "anime_series"
	MediaKindUnknown     MediaKind = "unknown"
)

type DiscoveredFile struct {
	Path          string
	FileName      string
	DetectedTitle string
	DetectedYear  *int32
	DetectedKind  MediaKind
	SafeMatch     bool
}

var (
	videoExtensions = map[string]struct{}{
		".avi":  {},
		".m4v":  {},
		".mkv":  {},
		".mov":  {},
		".mp4":  {},
		".mpeg": {},
		".mpg":  {},
		".ts":   {},
		".webm": {},
		".wmv":  {},
	}
	episodePattern = regexp.MustCompile(`(?i)\bS([0-9]{1,2})E([0-9]{1,3})\b`)
	yearPattern    = regexp.MustCompile(`\b(19[0-9]{2}|20[0-9]{2})\b`)
	releasePattern = regexp.MustCompile(`(?i)\b(480p|576p|720p|1080p|2160p|4k|bluray|brrip|webrip|web-dl|webdl|hdtv|hdrip|dvdrip|x264|x265|h264|h265|hevc|avc|aac|ac3|eac3|dts|atmos|truehd|proper|repack|remux|multi|dual)\b`)
)

func Discover(root string) ([]DiscoveredFile, error) {
	root = filepath.Clean(root)
	files := []DiscoveredFile{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if strings.HasPrefix(entry.Name(), ".") && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		if !isVideoFile(entry.Name()) {
			return nil
		}
		relativePath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, parseMediaFile(root, relativePath, entry.Name()))
		return nil
	})
	return files, err
}

func parseMediaFile(root string, relativePath string, fileName string) DiscoveredFile {
	withoutExtension := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	episodeMatch := episodePattern.FindStringIndex(withoutExtension)
	if episodeMatch != nil {
		title := seriesTitle(relativePath, withoutExtension, episodeMatch[0])
		return DiscoveredFile{
			Path:          filepath.ToSlash(relativePath),
			FileName:      fileName,
			DetectedTitle: title,
			DetectedKind:  MediaKindSeries,
			SafeMatch:     isCleanTitle(title) && !releasePattern.MatchString(withoutExtension),
		}
	}

	title, year := movieTitleAndYear(withoutExtension)
	return DiscoveredFile{
		Path:          filepath.ToSlash(relativePath),
		FileName:      fileName,
		DetectedTitle: title,
		DetectedYear:  year,
		DetectedKind:  MediaKindMovie,
		SafeMatch:     isCleanTitle(title) && !releasePattern.MatchString(withoutExtension),
	}
}

func seriesTitle(relativePath string, fileTitle string, episodeIndex int) string {
	parts := strings.Split(filepath.ToSlash(relativePath), "/")
	if len(parts) >= 3 && strings.HasPrefix(strings.ToLower(parts[len(parts)-2]), "season") {
		return cleanTitle(parts[len(parts)-3])
	}
	if len(parts) >= 2 {
		parent := cleanTitle(parts[len(parts)-2])
		if isCleanTitle(parent) {
			return parent
		}
	}
	return cleanTitle(fileTitle[:episodeIndex])
}

func movieTitleAndYear(value string) (string, *int32) {
	cleaned := cleanTitle(value)
	match := yearPattern.FindStringSubmatch(cleaned)
	if len(match) < 2 {
		return cleaned, nil
	}
	yearValue, err := strconv.ParseInt(match[1], 10, 32)
	if err != nil {
		return cleaned, nil
	}
	title := strings.TrimSpace(yearPattern.ReplaceAllString(cleaned, " "))
	title = cleanTitle(title)
	year := int32(yearValue)
	return title, &year
}

func cleanTitle(value string) string {
	value = strings.TrimSuffix(value, filepath.Ext(value))
	value = strings.NewReplacer(".", " ", "_", " ", "-", " ").Replace(value)
	value = strings.ReplaceAll(value, "(", " ")
	value = strings.ReplaceAll(value, ")", " ")
	value = strings.Join(strings.Fields(value), " ")
	return strings.TrimSpace(value)
}

func isCleanTitle(value string) bool {
	if strings.TrimSpace(value) == "" {
		return false
	}
	return !releasePattern.MatchString(value) && !episodePattern.MatchString(value)
}

func isVideoFile(name string) bool {
	_, ok := videoExtensions[strings.ToLower(filepath.Ext(name))]
	return ok
}
