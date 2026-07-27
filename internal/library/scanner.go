package library

import (
	"context"
	"io/fs"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"media-manager/internal/decisions"
	"media-manager/internal/episodepath"
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
	Path             string
	FileName         string
	SeriesFolderPath string
	SizeBytes        int64
	DetectedTitle    string
	DetectedYear     *int32
	DetectedKind     MediaKind
	SeasonNumber     *int32
	EpisodeNumber    *int32
	SafeMatch        bool
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
	episodePattern    = regexp.MustCompile(`(?i)\bS([0-9]{1,2})E([0-9]{1,3})\b`)
	yearPattern       = regexp.MustCompile(`\b(19[0-9]{2}|20[0-9]{2})\b`)
	externalIDPattern = regexp.MustCompile(`(?i)(?:\b(?:imdb(?:[ ._\-]*id)?|tmdb|tvdb)[ ._:\-]*(?:tt)?[0-9]{3,12}\b|\btt[0-9]{5,10}\b)`)
	releasePattern    = regexp.MustCompile(`(?i)\b(480p|576p|720p|1080p|2160p|4k|bluray|brrip|webrip|web-dl|webdl|hdtv|hdrip|dvdrip|x264|x265|h264|h265|hevc|avc|aac|ac3|eac3|dts|atmos|truehd|proper|repack|remux|multi|dual)\b`)
)

func Discover(root string) ([]DiscoveredFile, error) {
	return discover(context.Background(), root, parseMediaFile)
}

func DiscoverMovies(root string) ([]DiscoveredFile, error) {
	return DiscoverMoviesContext(context.Background(), root)
}

func DiscoverMoviesContext(ctx context.Context, root string) ([]DiscoveredFile, error) {
	return discover(ctx, root, parseMovieFile)
}

func DiscoverSeries(root string) ([]DiscoveredFile, error) {
	return DiscoverSeriesContext(context.Background(), root)
}

func DiscoverSeriesContext(ctx context.Context, root string) ([]DiscoveredFile, error) {
	return discover(ctx, root, parseSeriesFile)
}

func discover(
	ctx context.Context,
	root string,
	parse func(relativePath string, fileName string, info fs.FileInfo) DiscoveredFile,
) ([]DiscoveredFile, error) {
	root = filepath.Clean(root)
	files := []DiscoveredFile{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if contextErr := ctx.Err(); contextErr != nil {
			return contextErr
		}
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
		info, err := entry.Info()
		if err != nil {
			return err
		}
		files = append(files, parse(relativePath, entry.Name(), info))
		return nil
	})
	return files, err
}

func parseSeriesFile(relativePath string, fileName string, info fs.FileInfo) DiscoveredFile {
	folderPath, folderName := seriesFolder(relativePath)
	title, year := movieTitleAndYear(folderName)
	parsed := decisions.ParseReleaseFileName(fileName)
	if folderPath == "." {
		title = firstNonEmpty(cleanTitle(parsed.SeriesTitle), title)
	}
	season, episode, identified := episodepath.Numbers(relativePath)
	return DiscoveredFile{
		Path:             filepath.ToSlash(relativePath),
		FileName:         fileName,
		SeriesFolderPath: folderPath,
		SizeBytes:        info.Size(),
		DetectedTitle:    title,
		DetectedYear:     year,
		DetectedKind:     MediaKindSeries,
		SeasonNumber:     optionalEpisodeNumber(season, identified),
		EpisodeNumber:    optionalEpisodeNumber(episode, identified),
		SafeMatch:        isCleanTitle(title) && identified,
	}
}

func optionalEpisodeNumber(value int32, identified bool) *int32 {
	if !identified {
		return nil
	}
	return &value
}

func seriesFolder(relativePath string) (string, string) {
	parts := strings.Split(filepath.ToSlash(relativePath), "/")
	if len(parts) < 2 {
		return ".", strings.TrimSuffix(parts[0], filepath.Ext(parts[0]))
	}
	return parts[0], parts[0]
}

func parseMediaFile(relativePath string, fileName string, info fs.FileInfo) DiscoveredFile {
	withoutExtension := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	parsed := decisions.ParseReleaseFileName(fileName)
	if parsed.SeasonNumber != nil {
		title := firstNonEmpty(seriesTitleFromParsed(relativePath, parsed), seriesTitle(relativePath, withoutExtension, 0))
		return DiscoveredFile{
			Path:          filepath.ToSlash(relativePath),
			FileName:      fileName,
			SizeBytes:     info.Size(),
			DetectedTitle: title,
			DetectedKind:  MediaKindSeries,
			SeasonNumber:  parsed.SeasonNumber,
			EpisodeNumber: parsed.EpisodeNumber,
			SafeMatch:     isCleanTitle(title) && !releasePattern.MatchString(withoutExtension),
		}
	}

	return parseMovieFile(relativePath, fileName, info)
}

func parseMovieFile(relativePath string, fileName string, info fs.FileInfo) DiscoveredFile {
	withoutExtension := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	parsed := decisions.ParseReleaseFileName(fileName)
	title, year := movieTitleAndYear(firstNonEmpty(parsed.MovieTitle, withoutExtension))
	if folderTitle, folderYear := movieFolderTitle(relativePath, title); folderTitle != "" {
		title = folderTitle
		if folderYear != nil {
			year = folderYear
		}
	}
	return DiscoveredFile{
		Path:          filepath.ToSlash(relativePath),
		FileName:      fileName,
		SizeBytes:     info.Size(),
		DetectedTitle: title,
		DetectedYear:  firstYear(year, parsed.Year),
		DetectedKind:  MediaKindMovie,
		SafeMatch:     isCleanTitle(title) && !releasePattern.MatchString(withoutExtension),
	}
}

func seriesTitleFromParsed(relativePath string, parsed decisions.ParsedRelease) string {
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
	return cleanTitle(parsed.SeriesTitle)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func firstYear(existing *int32, parsed string) *int32 {
	if existing != nil || parsed == "" {
		return existing
	}
	value, err := strconv.ParseInt(parsed, 10, 32)
	if err != nil {
		return nil
	}
	year := int32(value)
	return &year
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

func movieFolderTitle(relativePath string, fileTitle string) (string, *int32) {
	directory := filepath.Dir(relativePath)
	if directory == "." {
		return "", nil
	}
	raw := filepath.Base(directory)
	title, year := movieTitleAndYear(raw)
	fileTitle = strings.ToLower(cleanTitle(fileTitle))
	folderTitle := strings.ToLower(cleanTitle(title))
	if year != nil || externalIDPattern.MatchString(raw) || fileTitle == folderTitle || fileTitle == "" || fileTitle == "video" || fileTitle == "movie" {
		return title, year
	}
	return "", nil
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
	value = externalIDPattern.ReplaceAllString(value, " ")
	if extension := strings.ToLower(filepath.Ext(value)); extension != "" {
		if _, ok := videoExtensions[extension]; ok {
			value = strings.TrimSuffix(value, filepath.Ext(value))
		}
	}
	value = strings.NewReplacer(".", " ", "_", " ", "-", " ").Replace(value)
	value = strings.NewReplacer("(", " ", ")", " ", "[", " ", "]", " ", "{", " ", "}", " ").Replace(value)
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
