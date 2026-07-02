package decisions

import (
	"path/filepath"
	"regexp"
	"strings"

	"media-manager/internal/storage"
)

type ParsedRelease struct {
	FileName      string
	ReleaseTitle  string
	MovieTitle    string
	SeriesTitle   string
	Year          string
	SeasonNumber  *int32
	EpisodeNumber *int32
	SeasonPack    bool
	Edition       string
	ReleaseGroup  string
	ReleaseHash   string
	QualityID     string
	Quality       string
	Source        string
	Resolution    string
	VideoCodec    string
	AudioCodec    string
	AudioChannels string
	Version       string
	Proper        bool
	Repack        bool
	Real          bool
	ReleaseType   string
	Languages     []string
}

var (
	yearPattern          = regexp.MustCompile(`\b(19\d{2}|20\d{2})\b`)
	resolutionPattern    = regexp.MustCompile(`(?i)\b(2160|1080|720|576|480)[pi]\b`)
	audioChannelsPattern = regexp.MustCompile(`(?i)(^|[^0-9])([257]\.1|[12345678]\.0)([^0-9]|$)`)
	versionPattern       = regexp.MustCompile(`(?i)\bv([0-9]+)\b`)
	hashPattern          = regexp.MustCompile(`(?i)\b[a-f0-9]{8,}\b`)
	extensionPattern     = regexp.MustCompile(`(?i)\.(mkv|mp4|avi|mov|wmv|ts|m2ts|iso)$`)
)

func ParseReleaseFileName(fileName string) ParsedRelease {
	base := filepath.Base(strings.TrimSpace(fileName))
	withoutExtension := extensionPattern.ReplaceAllString(base, "")
	title, group := splitReleaseGroup(withoutExtension)
	parsed := ParsedRelease{
		FileName:     base,
		ReleaseTitle: title,
		ReleaseGroup: group,
		ReleaseHash:  detectHash(title),
		Year:         detectYear(title),
		Source:       detectSource(title),
		Resolution:   detectResolution(title),
		VideoCodec:   detectVideoCodec(title),
		AudioCodec:   detectAudioCodec(title),
		Languages:    detectLanguages(title),
		Edition:      detectEdition(title),
		Proper:       containsAnyNormalized(title, "proper"),
		Repack:       containsAnyNormalized(title, "repack", "rerip"),
		Real:         containsAnyNormalized(title, "real"),
		Version:      detectVersion(title),
		ReleaseType:  detectReleaseType(title),
	}
	parsed.AudioChannels = detectAudioChannels(title)
	parsed.MovieTitle = releaseMovieTitle(title, parsed.Year)
	parsed.SeriesTitle = releaseSeriesTitle(title)
	parsed.SeasonNumber, parsed.EpisodeNumber = detectSeasonEpisode(title)
	parsed.SeasonPack = parsed.SeasonNumber != nil && parsed.EpisodeNumber == nil
	parsed.QualityID, parsed.Quality = detectReleaseQuality(parsed.Source, parsed.Resolution, title)
	return parsed
}

func splitReleaseGroup(title string) (string, string) {
	index := strings.LastIndex(title, "-")
	if index <= 0 || index >= len(title)-1 {
		return title, ""
	}
	group := strings.TrimSpace(title[index+1:])
	if strings.ContainsAny(group, " ._") || len(group) > 40 {
		return title, ""
	}
	return strings.TrimSpace(title[:index]), group
}

func detectYear(title string) string {
	return firstMatch(yearPattern, title)
}

func detectResolution(title string) string {
	value := firstMatch(resolutionPattern, title)
	if value == "" {
		return ""
	}
	if strings.HasSuffix(strings.ToLower(value), "p") || strings.HasSuffix(strings.ToLower(value), "i") {
		return strings.ToLower(value)
	}
	return value + "p"
}

func detectSource(title string) string {
	switch {
	case containsAnyNormalized(title, "webdl", "web-dl"):
		return "WEBDL"
	case containsAnyNormalized(title, "webrip", "web-rip"):
		return "WEBRip"
	case containsAnyNormalized(title, "remux"):
		return "Remux"
	case containsAnyNormalized(title, "bluray", "blu-ray", "bdrip", "brrip"):
		return "Bluray"
	case containsAnyNormalized(title, "hdtv"):
		return "HDTV"
	case containsAnyNormalized(title, "sdtv"):
		return "SDTV"
	case containsAnyNormalized(title, "dvd-r", "dvdr"):
		return "DVD-R"
	case containsAnyNormalized(title, "dvd"):
		return "DVD"
	case containsAnyNormalized(title, "cam"):
		return "CAM"
	case containsAnyNormalized(title, "telesync"):
		return "TELESYNC"
	case containsAnyNormalized(title, "telecine"):
		return "TELECINE"
	default:
		return ""
	}
}

func detectVideoCodec(title string) string {
	switch {
	case containsAnyNormalized(title, "x266", "h266", "vvc"):
		return "x266"
	case containsAnyNormalized(title, "x265", "h265", "hevc"):
		return "x265"
	case containsAnyNormalized(title, "x264", "h264", "avc"):
		return "x264"
	case containsAnyNormalized(title, "av1"):
		return "AV1"
	case containsAnyNormalized(title, "xvid"):
		return "Xvid"
	case containsAnyNormalized(title, "vp9"):
		return "VP9"
	default:
		return ""
	}
}

func detectAudioCodec(title string) string {
	switch {
	case containsAnyNormalized(title, "truehd", "atmos"):
		return "TrueHD/Atmos"
	case containsAnyNormalized(title, "ddp", "ddplus", "eac3", "e-ac3"):
		return "DD+"
	case containsAnyNormalized(title, "ac3", "dolbydigital"):
		return "DD"
	case containsAnyNormalized(title, "dts"):
		return "DTS"
	case containsAnyNormalized(title, "aac"):
		return "AAC"
	case containsAnyNormalized(title, "flac"):
		return "FLAC"
	case containsAnyNormalized(title, "pcm"):
		return "PCM"
	default:
		return ""
	}
}

func detectAudioChannels(title string) string {
	match := audioChannelsPattern.FindStringSubmatch(title)
	if len(match) < 3 {
		return ""
	}
	return match[2]
}

func detectLanguages(title string) []string {
	languages := []string{}
	if containsAnyNormalized(title, "multi", "multilang") {
		languages = append(languages, "Multiple")
	}
	languageTokens := map[string]string{
		"english": "English", "eng": "English", "en": "English",
		"german": "German", "ger": "German", "de": "German",
		"french": "French", "fre": "French", "fr": "French",
		"spanish": "Spanish", "spa": "Spanish", "es": "Spanish",
		"japanese": "Japanese", "jpn": "Japanese", "ja": "Japanese",
		"korean": "Korean", "kor": "Korean", "ko": "Korean",
		"chinese": "Chinese", "chi": "Chinese", "zh": "Chinese",
	}
	normalized := tokenSet(title)
	for token, label := range languageTokens {
		if _, ok := normalized[token]; ok && !containsString(languages, label) {
			languages = append(languages, label)
		}
	}
	return languages
}

func detectEdition(title string) string {
	for _, label := range []string{"Directors Cut", "Extended", "Special Edition", "Unrated", "Uncut", "Remastered", "Theatrical"} {
		if containsAnyNormalized(title, label) {
			return label
		}
	}
	return ""
}

func detectVersion(title string) string {
	if match := firstMatch(versionPattern, title); match != "" {
		return match
	}
	return ""
}

func detectReleaseType(title string) string {
	if containsAnyNormalized(title, "episode") {
		return "Episode"
	}
	if containsAnyNormalized(title, "movie") {
		return "Movie"
	}
	return ""
}

func detectHash(title string) string {
	matches := hashPattern.FindAllString(title, -1)
	if len(matches) == 0 {
		return ""
	}
	return matches[len(matches)-1]
}

func releaseMovieTitle(title string, year string) string {
	end := len(title)
	if year != "" {
		if index := strings.Index(title, year); index > 0 {
			end = index
		}
	}
	return strings.Join(strings.Fields(releaseSeparator.ReplaceAllString(title[:end], " ")), " ")
}

func detectReleaseQuality(source string, resolution string, title string) (string, string) {
	best := storage.QualitySizeDefinition{}
	for _, definition := range storage.QualitySizeDefinitions() {
		if qualityMatches(definition, source, resolution, title) && definition.SortOrder > best.SortOrder {
			best = definition
		}
	}
	return best.ID, best.Name
}

func qualityMatches(definition storage.QualitySizeDefinition, source string, resolution string, title string) bool {
	name := strings.ToLower(definition.Name)
	if source != "" && resolution != "" {
		return strings.Contains(name, strings.ToLower(source)) && strings.Contains(name, resolution)
	}
	return containsAnyNormalized(title, definition.ID, definition.Name)
}

var releaseSeparator = regexp.MustCompile(`[._]+`)

func firstMatch(pattern *regexp.Regexp, title string) string {
	match := pattern.FindString(title)
	return strings.TrimSpace(match)
}
