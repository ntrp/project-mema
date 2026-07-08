package httpapi

import (
	"os"
	"path/filepath"
	"strings"

	"media-manager/internal/delivery"
	"media-manager/internal/storage"
)

func mediaFileInfoResponses(
	paths []string,
	subtitleTargets []storage.MediaProfileSubtitleTarget,
	subtitleMode string,
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
			probe := delivery.Probe(mediaFileProbePath(path))
			tracks := mediaFileTracksFromDelivery(probe.Tracks)
			chapters := mediaFileChaptersFromDelivery(probe.Chapters)
			hydrateTrackProvenance(path, tracks, componentProvenance)
			if len(tracks) > 0 {
				file.Tracks = &tracks
			}
			if len(chapters) > 0 {
				file.Chapters = &chapters
			}
			file.SubtitleSatisfaction = mediaFileSubtitleSatisfaction(
				tracks,
				subtitleTargets,
				subtitleMode,
				externalSubtitleLanguagesForPath(externalSubtitles, sidecars, path),
			)
			otherFiles := mediaFileOtherFiles(path, paths, subtitleTargets, subtitleMode, externalSubtitles, sidecars, file.SubtitleSatisfaction)
			if len(otherFiles) > 0 {
				file.OtherFiles = &otherFiles
			}
		}
		rollup := mediaFileRollupSummary(file.Status)
		targetSatisfaction := targetSatisfactionSummaryResponse(nil, nil)
		file.Rollup = &rollup
		file.TargetSatisfaction = &targetSatisfaction
		files = append(files, file)
	}
	return &files
}

func mediaFileProbePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return absolute
}

func mediaFileSubtitleSatisfaction(
	tracks []MediaFileTrack,
	targets []storage.MediaProfileSubtitleTarget,
	subtitleMode string,
	externalLanguages []string,
) *MediaFileSubtitleSatisfaction {
	if len(targets) == 0 {
		return &MediaFileSubtitleSatisfaction{
			Mode:             mediaFileSubtitleMode(subtitleMode),
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
		if subtitleTargetSatisfied(target, subtitleMode, embedded, external) {
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
		Mode:             mediaFileSubtitleMode(subtitleMode),
		State:            state,
		WantedLanguages:  wanted,
		MatchedLanguages: matched,
		MissingLanguages: missing,
	}
}

func mediaFileSubtitleMode(value string) MediaProfileSubtitleMode {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(MediaProfileSubtitleModeEmbedded):
		return MediaProfileSubtitleModeEmbedded
	case string(MediaProfileSubtitleModeExternal):
		return MediaProfileSubtitleModeExternal
	default:
		return MediaProfileSubtitleModeMixed
	}
}

func subtitleTargetSatisfied(
	target storage.MediaProfileSubtitleTarget,
	subtitleMode string,
	embedded map[string]struct{},
	external map[string]struct{},
) bool {
	language := languageMatchKey(target.LanguageID)
	_, embeddedOK := embedded[language]
	_, externalOK := external[language]
	switch subtitleMode {
	case "embedded":
		return embeddedOK
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
