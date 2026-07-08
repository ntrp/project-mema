package storage

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

var mediaRenameResolutionPattern = regexp.MustCompile(`(?i)\b(2160|1080|720|576|480)[pi]\b`)

func mediaRenameQualityFull(path string) string {
	return mediaRenameQualityDefinition(path).Name
}

func mediaRenameQualityDefinition(path string) QualitySizeDefinition {
	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	source := mediaRenameSource(title)
	resolution := mediaRenameResolution(title)
	best := QualitySizeDefinition{}
	for _, definition := range QualitySizeDefinitions() {
		if mediaRenameQualityMatches(definition, source, resolution, title) &&
			definition.SortOrder > best.SortOrder {
			best = definition
		}
	}
	return best
}

func mediaRenameSource(title string) string {
	switch {
	case mediaRenameContainsAny(title, "webdl", "web-dl"):
		return "WEBDL"
	case mediaRenameContainsAny(title, "webrip", "web-rip"):
		return "WEBRip"
	case mediaRenameContainsAny(title, "remux"):
		return "Remux"
	case mediaRenameContainsAny(title, "bluray", "blu-ray", "bdrip", "brrip"):
		return "Bluray"
	case mediaRenameContainsAny(title, "hdtv"):
		return "HDTV"
	case mediaRenameContainsAny(title, "sdtv"):
		return "SDTV"
	case mediaRenameContainsAny(title, "dvd-r", "dvdr"):
		return "DVD-R"
	case mediaRenameContainsAny(title, "dvd"):
		return "DVD"
	case mediaRenameContainsAny(title, "cam"):
		return "CAM"
	case mediaRenameContainsAny(title, "telesync"):
		return "TELESYNC"
	case mediaRenameContainsAny(title, "telecine"):
		return "TELECINE"
	default:
		return ""
	}
}

func mediaRenameResolution(title string) string {
	value := mediaRenameResolutionPattern.FindString(title)
	if value == "" {
		return ""
	}
	if strings.HasSuffix(strings.ToLower(value), "p") || strings.HasSuffix(strings.ToLower(value), "i") {
		return strings.ToLower(value)
	}
	return value + "p"
}

func mediaRenameQualityMatches(definition QualitySizeDefinition, source string, resolution string, title string) bool {
	name := strings.ToLower(definition.Name)
	if source != "" && resolution != "" {
		return strings.Contains(name, strings.ToLower(source)) && strings.Contains(name, resolution)
	}
	normalizedTitle := mediaRenameNormalize(title)
	return strings.Contains(normalizedTitle, mediaRenameNormalize(definition.ID)) ||
		strings.Contains(normalizedTitle, mediaRenameNormalize(definition.Name))
}

func mediaRenameContainsAny(title string, values ...string) bool {
	normalized := mediaRenameNormalize(title)
	for _, value := range values {
		if strings.Contains(normalized, mediaRenameNormalize(value)) {
			return true
		}
	}
	return false
}

func mediaRenameNormalize(value string) string {
	var builder strings.Builder
	for _, r := range strings.ToLower(value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
