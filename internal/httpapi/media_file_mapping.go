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
	subtitlePreferredMode string,
	externalSubtitles []storage.MediaItemSubtitle,
	componentProvenance []storage.MediaComponentProvenance,
	sidecars []storage.MediaItemSidecar,
) *[]MediaFileInfo {
	files := make([]MediaFileInfo, 0, len(paths))
	for _, path := range paths {
		file := MediaFileInfo{Path: path, Status: MediaFileInfoStatusMissing}
		if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
			file.Status = MediaFileInfoStatusAvailable
			size := stat.Size()
			file.SizeBytes = &size
			probe := mediaFileProbe(path)
			hydrateTrackProvenance(path, probe.tracks, componentProvenance)
			if len(probe.tracks) > 0 {
				file.Tracks = &probe.tracks
			}
			if len(probe.chapters) > 0 {
				file.Chapters = &probe.chapters
			}
			file.SubtitleSatisfaction = mediaFileSubtitleSatisfaction(
				probe.tracks,
				subtitleTargets,
				subtitlePreferredMode,
				externalSubtitleLanguagesForPath(externalSubtitles, sidecars, path),
			)
			otherFiles := mediaFileOtherFiles(path, paths, subtitleTargets, subtitlePreferredMode, externalSubtitles, sidecars, file.SubtitleSatisfaction)
			if len(otherFiles) > 0 {
				file.OtherFiles = &otherFiles
			}
		}
		files = append(files, file)
	}
	return &files
}

func mediaFileSubtitleSatisfaction(
	tracks []MediaFileTrack,
	targets []storage.MediaProfileSubtitleTarget,
	subtitlePreferredMode string,
	externalLanguages []string,
) *MediaFileSubtitleSatisfaction {
	if len(targets) == 0 {
		return &MediaFileSubtitleSatisfaction{
			PreferredMode:    mediaFileSubtitlePreferredMode(subtitlePreferredMode),
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
	seenTargets := map[string]struct{}{}
	for _, target := range targets {
		language := languageMatchKey(target.LanguageID)
		if language == "" {
			continue
		}
		if _, ok := seenTargets[language]; ok {
			continue
		}
		seenTargets[language] = struct{}{}
		wanted = append(wanted, target.LanguageID)
		if subtitleTargetSatisfied(target, subtitlePreferredMode, embedded, external) {
			matched = append(matched, target.LanguageID)
			continue
		}
		missing = append(missing, target.LanguageID)
	}
	state := MediaFileSubtitleSatisfactionStateSatisfied
	if len(missing) > 0 {
		state = MediaFileSubtitleSatisfactionStateMissing
	}
	return &MediaFileSubtitleSatisfaction{
		PreferredMode:    mediaFileSubtitlePreferredMode(subtitlePreferredMode),
		State:            state,
		WantedLanguages:  wanted,
		MatchedLanguages: matched,
		MissingLanguages: missing,
	}
}

func mediaFileSubtitlePreferredMode(value string) MediaProfileSubtitlePreferredMode {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(MediaProfileSubtitlePreferredModeEmbedded):
		return MediaProfileSubtitlePreferredModeEmbedded
	case string(MediaProfileSubtitlePreferredModeExternal):
		return MediaProfileSubtitlePreferredModeExternal
	default:
		return MediaProfileSubtitlePreferredModeMixed
	}
}

func subtitleTargetSatisfied(
	target storage.MediaProfileSubtitleTarget,
	subtitlePreferredMode string,
	embedded map[string]struct{},
	external map[string]struct{},
) bool {
	language := languageMatchKey(target.LanguageID)
	_, embeddedOK := embedded[language]
	_, externalOK := external[language]
	switch subtitlePreferredMode {
	case "embedded":
		return embeddedOK || externalOK
	case "external":
		return externalOK
	default:
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

func externalSubtitleLanguagesForPath(
	subtitles []storage.MediaItemSubtitle,
	sidecars []storage.MediaItemSidecar,
	path string,
) []string {
	languages := []string{}
	seen := map[string]struct{}{}
	for _, subtitle := range subtitles {
		if !sameSubtitleMediaBase(subtitle.FilePath, path) {
			continue
		}
		languages = appendLanguage(languages, seen, subtitle.LanguageID)
	}
	for _, sidecar := range sidecars {
		if sidecar.MediaFilePath != path ||
			sidecar.SidecarType != storage.MediaSidecarSubtitle ||
			sidecar.LanguageID == nil ||
			otherFileStatus(sidecar.FilePath) != MediaFileOtherFileStatusAvailable {
			continue
		}
		languages = appendLanguage(languages, seen, *sidecar.LanguageID)
	}
	for _, sidecar := range storage.MediaSidecarsForFile(path) {
		if sidecar.Type != storage.MediaSidecarSubtitle || sidecar.LanguageID == "" {
			continue
		}
		languages = appendLanguage(languages, seen, sidecar.LanguageID)
	}
	return languages
}

func appendLanguage(languages []string, seen map[string]struct{}, value string) []string {
	language := languageMatchKey(value)
	if language == "" {
		return languages
	}
	if _, ok := seen[language]; ok {
		return languages
	}
	seen[language] = struct{}{}
	return append(languages, value)
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
