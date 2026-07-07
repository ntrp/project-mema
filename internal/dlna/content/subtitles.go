package content

import (
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"media-manager/internal/storage"
)

type Subtitle struct {
	URL      string
	FilePath string
	Language string
	Format   string
	Plan     SubtitlePlan
}

type SubtitlePlan string

const (
	SubtitleDirect  SubtitlePlan = "direct"
	SubtitleConvert SubtitlePlan = "convert"
	SubtitleOmit    SubtitlePlan = "omit"
)

func ApplySubtitleURLs(baseURL string, objects []Object) []Object {
	if strings.TrimSpace(baseURL) == "" {
		return objects
	}
	updated := make([]Object, len(objects))
	copy(updated, objects)
	for index := range updated {
		updated[index].Subtitles = append([]Subtitle{}, updated[index].Subtitles...)
		for subIndex := range updated[index].Subtitles {
			updated[index].Subtitles[subIndex].URL = SubtitleURL(baseURL, updated[index].ID, subIndex)
		}
	}
	return updated
}

func SubtitleURL(baseURL string, objectID string, index int) string {
	base, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		return "/dlna/subtitle/" + url.PathEscape(objectID) + "/" + strconv.Itoa(index)
	}
	base.Path = path.Join(base.Path, "/dlna/subtitle", objectID, strconv.Itoa(index))
	return base.String()
}

func SubtitlesForFile(item storage.MediaItem, mediaFilePath string) []Subtitle {
	var subtitles []Subtitle
	for _, sidecar := range item.Sidecars {
		if sidecar.MediaFilePath == mediaFilePath && sidecar.SidecarType == storage.MediaSidecarSubtitle {
			if subtitle, ok := subtitleFromPath(sidecar.FilePath, sidecar.LanguageID, sidecar.Format); ok {
				subtitles = append(subtitles, subtitle)
			}
		}
	}
	for _, external := range item.ExternalSubtitles {
		if external.RetentionMode == storage.SubtitleRetentionExternal {
			if subtitle, ok := subtitleFromPath(external.FilePath, &external.LanguageID, &external.Format); ok {
				subtitles = append(subtitles, subtitle)
			}
		}
	}
	return subtitles
}

func subtitleFromPath(filePath string, language *string, format *string) (Subtitle, bool) {
	value := strings.Trim(strings.ToLower(strings.TrimPrefix(optionalValue(format), ".")), " ")
	if value == "" {
		value = strings.TrimPrefix(strings.ToLower(filepath.Ext(filePath)), ".")
	}
	plan := PlanSubtitle(value)
	if plan == SubtitleOmit {
		return Subtitle{}, false
	}
	return Subtitle{FilePath: filePath, Language: optionalValue(language), Format: value, Plan: plan}, true
}

func PlanSubtitle(format string) SubtitlePlan {
	switch strings.ToLower(strings.TrimPrefix(format, ".")) {
	case "srt", "vtt":
		return SubtitleDirect
	case "ass", "ssa":
		return SubtitleConvert
	default:
		return SubtitleOmit
	}
}

func SubtitleProtocolInfo(format string) string {
	switch strings.ToLower(format) {
	case "vtt":
		return "http-get:*:text/vtt:*"
	default:
		return "http-get:*:application/x-subrip:*"
	}
}

func optionalValue(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}
