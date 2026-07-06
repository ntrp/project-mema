package httpapi

import (
	"os"
	"path/filepath"
	"strings"

	"media-manager/internal/storage"
)

func mediaFileInfoResponses(
	paths []string,
	subtitleTargets []storage.MediaProfileSubtitleTarget,
	externalSubtitles []storage.MediaItemSubtitle,
) *[]MediaFileInfo {
	files := make([]MediaFileInfo, 0, len(paths))
	for _, path := range paths {
		file := MediaFileInfo{Path: path, Status: MediaFileInfoStatusMissing}
		if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
			file.Status = MediaFileInfoStatusAvailable
			size := stat.Size()
			file.SizeBytes = &size
			probe := mediaFileProbe(path)
			if len(probe.tracks) > 0 {
				file.Tracks = &probe.tracks
			}
			if len(probe.chapters) > 0 {
				file.Chapters = &probe.chapters
			}
			file.SubtitleSatisfaction = mediaFileSubtitleSatisfaction(
				probe.tracks,
				subtitleTargets,
				externalSubtitleLanguagesForPath(externalSubtitles, path),
			)
		}
		files = append(files, file)
	}
	return &files
}

func mediaFileSubtitleSatisfaction(
	tracks []MediaFileTrack,
	targets []storage.MediaProfileSubtitleTarget,
	externalLanguages []string,
) *MediaFileSubtitleSatisfaction {
	if len(targets) == 0 {
		return &MediaFileSubtitleSatisfaction{
			State:            MediaFileSubtitleSatisfactionStateIgnored,
			WantedLanguages:  []string{},
			MatchedLanguages: []string{},
			MissingLanguages: []string{},
		}
	}
	embedded := mediaFileSubtitleLanguageSet(tracks)
	external := languageSet(externalLanguages)
	wanted := make([]string, 0, len(targets))
	matched := []string{}
	missing := []string{}
	for _, target := range targets {
		wanted = append(wanted, target.LanguageID)
		if subtitleTargetSatisfied(target, embedded, external) {
			matched = append(matched, target.LanguageID)
			continue
		}
		if target.Required {
			missing = append(missing, target.LanguageID)
		}
	}
	state := MediaFileSubtitleSatisfactionStateSatisfied
	if len(missing) > 0 {
		state = MediaFileSubtitleSatisfactionStateMissing
	}
	return &MediaFileSubtitleSatisfaction{
		State:            state,
		WantedLanguages:  wanted,
		MatchedLanguages: matched,
		MissingLanguages: missing,
	}
}

func subtitleTargetSatisfied(
	target storage.MediaProfileSubtitleTarget,
	embedded map[string]struct{},
	external map[string]struct{},
) bool {
	language := languageMatchKey(target.LanguageID)
	switch target.Source {
	case "embedded":
		_, ok := embedded[language]
		return ok
	case "external":
		_, ok := external[language]
		return ok
	default:
		_, embeddedOK := embedded[language]
		_, externalOK := external[language]
		return embeddedOK || externalOK
	}
}

func mediaFileSubtitleLanguageSet(tracks []MediaFileTrack) map[string]struct{} {
	languages := map[string]struct{}{}
	for _, track := range tracks {
		if track.Type != Subtitle {
			continue
		}
		language := languageMatchKey(optionalStringValue(track.Language))
		if language != "" {
			languages[language] = struct{}{}
		}
	}
	return languages
}

func languageSet(values []string) map[string]struct{} {
	languages := map[string]struct{}{}
	for _, value := range values {
		language := languageMatchKey(value)
		if language != "" {
			languages[language] = struct{}{}
		}
	}
	return languages
}

func externalSubtitleLanguagesForPath(subtitles []storage.MediaItemSubtitle, path string) []string {
	languages := []string{}
	for _, subtitle := range subtitles {
		if !sameSubtitleMediaBase(subtitle.FilePath, path) {
			continue
		}
		languages = append(languages, subtitle.LanguageID)
	}
	return languages
}

func sameSubtitleMediaBase(subtitlePath string, mediaPath string) bool {
	subtitleBase := strings.TrimSuffix(filepath.Base(subtitlePath), filepath.Ext(subtitlePath))
	mediaBase := strings.TrimSuffix(filepath.Base(mediaPath), filepath.Ext(mediaPath))
	return strings.HasPrefix(strings.ToLower(subtitleBase), strings.ToLower(mediaBase)+".")
}

func languageMatchKey(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.TrimSuffix(normalized, " language")
	switch normalized {
	case "en", "eng", "english":
		return "english"
	case "de", "deu", "ger", "german":
		return "german"
	case "fr", "fra", "fre", "french":
		return "french"
	case "es", "spa", "spanish":
		return "spanish"
	case "ja", "jpn", "japanese":
		return "japanese"
	default:
		return strings.ReplaceAll(normalized, " ", "-")
	}
}

func optionalStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
