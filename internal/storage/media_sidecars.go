package storage

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"media-manager/internal/subtitleformats"
)

type MediaSidecarType string

const (
	MediaSidecarMetadata MediaSidecarType = "metadata"
	MediaSidecarSubtitle MediaSidecarType = "subtitle"
	MediaSidecarUnknown  MediaSidecarType = "unknown"
)

type MediaSidecar struct {
	Type       MediaSidecarType
	Path       string
	Subtype    string
	LanguageID string
	Format     string
}

var mediaSidecarSubtitleExtensions = map[string]struct{}{
	".ass": {},
	".idx": {},
	".srt": {},
	".ssa": {},
	".sub": {},
	".vtt": {},
}

var mediaSidecarMetadataExtensions = map[string]struct{}{
	".jpeg": {},
	".jpg":  {},
	".nfo":  {},
	".png":  {},
	".tbn":  {},
	".webp": {},
}

var mediaSidecarArtNames = map[string]struct{}{
	"backdrop": {}, "background": {}, "banner": {}, "clearart": {}, "clearlogo": {},
	"cover": {}, "disc": {}, "discart": {}, "fanart": {}, "folder": {},
	"keyart": {}, "landscape": {}, "logo": {}, "movie": {}, "poster": {},
	"thumb": {}, "thumbnail": {},
}

var mediaSidecarLanguageIDs = map[string]string{
	"ara": "arabic", "ar": "arabic", "arabic": "arabic",
	"chi": "chinese", "chinese": "chinese", "zh": "chinese", "zho": "chinese",
	"dan": "danish", "danish": "danish", "da": "danish",
	"de": "german", "deu": "german", "ger": "german", "german": "german",
	"dut": "dutch", "dutch": "dutch", "nl": "dutch", "nld": "dutch",
	"en": "english", "eng": "english", "english": "english",
	"es": "spanish", "spa": "spanish", "spanish": "spanish",
	"fi": "finnish", "fin": "finnish", "finnish": "finnish",
	"fr": "french", "fra": "french", "fre": "french", "french": "french",
	"hi": "hindi", "hin": "hindi", "hindi": "hindi",
	"it": "italian", "ita": "italian", "italian": "italian",
	"ja": "japanese", "japanese": "japanese", "jpn": "japanese",
	"ko": "korean", "kor": "korean", "korean": "korean",
	"no": "norwegian", "nor": "norwegian", "norwegian": "norwegian",
	"pl": "polish", "pol": "polish", "polish": "polish",
	"por": "portuguese", "portuguese": "portuguese", "pt": "portuguese",
	"ru": "russian", "rus": "russian", "russian": "russian",
	"sv": "swedish", "swe": "swedish", "swedish": "swedish",
}

var mediaSidecarSubtitleFlags = map[string]struct{}{
	"cc": {}, "default": {}, "forced": {}, "foreign": {}, "full": {},
	"hi": {}, "normal": {}, "sdh": {}, "signs": {}, "songs": {},
}

func MediaSidecarsForFile(mediaPath string) []MediaSidecar {
	entries, err := os.ReadDir(filepath.Dir(mediaPath))
	if err != nil {
		return nil
	}
	sidecars := []MediaSidecar{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(filepath.Dir(mediaPath), entry.Name())
		if filepath.Clean(path) == filepath.Clean(mediaPath) {
			continue
		}
		if _, ok := mediaFileExtensions[strings.ToLower(filepath.Ext(path))]; ok {
			continue
		}
		sidecar := ClassifyMediaSidecar(mediaPath, path)
		sidecars = append(sidecars, sidecar)
	}
	sort.Slice(sidecars, func(i, j int) bool { return sidecars[i].Path < sidecars[j].Path })
	return sidecars
}

func ClassifyMediaSidecar(mediaPath string, path string) MediaSidecar {
	ext := strings.ToLower(filepath.Ext(path))
	base := sidecarBase(path)
	mediaBase := sidecarBase(mediaPath)
	switch {
	case isMediaSidecarSubtitleExt(ext):
		format := strings.TrimPrefix(ext, ".")
		language := sidecarSubtitleLanguage(mediaBase, base)
		if language == "" {
			language = sidecarSubtitleStandaloneLanguage(base)
		}
		if language != "" || sidecarRelatedBase(base, mediaBase) {
			return MediaSidecar{
				Type:       MediaSidecarSubtitle,
				Path:       path,
				Subtype:    sidecarSubtitleSubtype(format),
				LanguageID: language,
				Format:     format,
			}
		}
	case isMediaSidecarMetadataExt(ext) && sidecarMetadataBase(base, mediaBase, ext):
		return MediaSidecar{Type: MediaSidecarMetadata, Path: path, Subtype: sidecarMetadataSubtype(base, mediaBase, ext)}
	}
	return MediaSidecar{Type: MediaSidecarUnknown, Path: path, Subtype: strings.TrimPrefix(ext, ".")}
}

func sidecarSubtitleSubtype(format string) string {
	if normalized := subtitleformats.Normalize(format); normalized != "" {
		return normalized
	}
	return strings.TrimPrefix(strings.ToLower(strings.TrimSpace(format)), ".")
}

func sidecarMetadataSubtype(base string, mediaBase string, ext string) string {
	if ext == ".nfo" {
		return "nfo"
	}
	if subtype := sidecarArtSubtype(base); subtype != "" {
		return subtype
	}
	if suffix, ok := sidecarSuffix(mediaBase, base); ok {
		if subtype := sidecarArtSubtype(strings.Trim(suffix, ".-_ ")); subtype != "" {
			return subtype
		}
	}
	return strings.TrimPrefix(ext, ".")
}

func sidecarArtSubtype(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "background":
		return "backdrop"
	case "thumbnail":
		return "thumb"
	default:
		if _, ok := mediaSidecarArtNames[value]; ok {
			return value
		}
	}
	return ""
}

func isMediaSidecarSubtitleExt(ext string) bool {
	_, ok := mediaSidecarSubtitleExtensions[ext]
	return ok
}

func isMediaSidecarMetadataExt(ext string) bool {
	_, ok := mediaSidecarMetadataExtensions[ext]
	return ok
}

func sidecarMetadataBase(base string, mediaBase string, ext string) bool {
	if base == mediaBase || ext == ".nfo" && sidecarRelatedBase(base, mediaBase) {
		return true
	}
	if _, ok := mediaSidecarArtNames[base]; ok {
		return true
	}
	suffix, ok := sidecarSuffix(mediaBase, base)
	if !ok {
		return false
	}
	_, ok = mediaSidecarArtNames[strings.Trim(suffix, ".-_ ")]
	return ok
}

func sidecarSubtitleLanguage(mediaBase string, base string) string {
	suffix, ok := sidecarSuffix(mediaBase, base)
	if !ok {
		return ""
	}
	for _, token := range subtitleLanguageTokens(suffix) {
		if language, ok := mediaSidecarLanguageIDs[token]; ok {
			return language
		}
	}
	return ""
}

func sidecarSubtitleStandaloneLanguage(base string) string {
	for _, token := range subtitleLanguageTokens(base) {
		if language, ok := mediaSidecarLanguageIDs[token]; ok {
			return language
		}
	}
	return ""
}

func subtitleLanguageTokens(value string) []string {
	tokens := []string{}
	for _, token := range strings.FieldsFunc(value, func(r rune) bool {
		return r == '.' || r == '-' || r == '_' || r == ' '
	}) {
		token = strings.ToLower(strings.TrimSpace(token))
		if token == "" {
			continue
		}
		if _, ok := mediaSidecarSubtitleFlags[token]; ok {
			continue
		}
		tokens = append(tokens, token)
	}
	return tokens
}

func sidecarRelatedBase(base string, mediaBase string) bool {
	_, ok := sidecarSuffix(mediaBase, base)
	return ok
}

func sidecarSuffix(mediaBase string, base string) (string, bool) {
	if base == mediaBase {
		return "", true
	}
	for _, separator := range []string{".", "-", "_", " "} {
		prefix := mediaBase + separator
		if strings.HasPrefix(base, prefix) {
			return strings.TrimPrefix(base, prefix), true
		}
	}
	return "", false
}

func sidecarBase(path string) string {
	return strings.ToLower(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
}
