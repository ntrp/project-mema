package httpapi

import (
	"fmt"
	"strconv"
	"strings"

	"media-manager/internal/storage"
	"media-manager/internal/subtitleformats"
)

func unwantedLanguage(enabled bool, expected []string, language string) bool {
	if !enabled {
		return false
	}
	return outsideTargetLanguage(expected, language)
}

func outsideTargetLanguage(expected []string, language string) bool {
	if len(expected) == 0 {
		return false
	}
	key := languageMatchKey(language)
	if key == "" {
		return false
	}
	for _, value := range expected {
		if languageMatchKey(value) == key {
			return false
		}
	}
	return true
}

func unwantedAudioTracks(item storage.MediaItem, tracks []MediaFileTrack) []string {
	if !item.RemoveUnwantedAudio {
		return nil
	}
	values := []string{}
	for _, track := range tracks {
		language := optionalStringValue(track.Language)
		if unwantedLanguage(true, audioTargetLanguages(item.AudioTargets), language) {
			values = append(values, language)
		}
	}
	return values
}

func subtitleTargetMatches(targets []storage.MediaProfileSubtitleTarget, language string) bool {
	for _, target := range targets {
		if languageMatchKey(target.LanguageID) == languageMatchKey(language) {
			return true
		}
	}
	return false
}

func subtitleTargetFormatMatches(targets []storage.MediaProfileSubtitleTarget, language string, format string) bool {
	for _, target := range targets {
		if languageMatchKey(target.LanguageID) != languageMatchKey(language) {
			continue
		}
		return subtitleformats.AnyMatch(target.Formats, format)
	}
	return len(targets) == 0
}

func subtitleTargetTextConversionSupported(
	targets []storage.MediaProfileSubtitleTarget,
	language string,
	format string,
) bool {
	if !subtitleformats.Text(format) {
		return false
	}
	for _, target := range targets {
		if languageMatchKey(target.LanguageID) == languageMatchKey(language) {
			return subtitleformats.HasTextTarget(target.Formats)
		}
	}
	return false
}

func hasPendingSubtitleCandidate(tracks []MediaFileTrack, otherFiles []MediaFileOtherFile) bool {
	for _, track := range tracks {
		if track.State != nil && track.State.VisualState == MediaFileDetailVisualStatePendingOperation {
			return true
		}
	}
	for _, file := range otherFiles {
		if file.State != nil && file.State.VisualState == MediaFileDetailVisualStatePendingOperation {
			return true
		}
	}
	return false
}

func hasTrackLanguage(tracks []MediaFileTrack, trackType MediaFileTrackType, language string) bool {
	for _, track := range tracks {
		if track.Type == trackType && languageMatchKey(optionalStringValue(track.Language)) == languageMatchKey(language) {
			return true
		}
	}
	return false
}

func hasSubtitleCandidate(tracks []MediaFileTrack, otherFiles []MediaFileOtherFile, language string) bool {
	if hasTrackLanguage(tracks, MediaFileTrackTypeSubtitle, language) {
		return true
	}
	for _, file := range otherFiles {
		if file.Type == MediaFileOtherFileTypeSubtitle && languageMatchKey(optionalStringValue(file.Language)) == languageMatchKey(language) {
			return true
		}
	}
	return false
}

func audioChannelValue(track MediaFileTrack) string {
	if track.ChannelLayout != nil {
		return *track.ChannelLayout
	}
	if track.Channels != nil {
		return fmt.Sprintf("%g", float64(*track.Channels))
	}
	return ""
}

func trackBitrateKbps(track MediaFileTrack) int32 {
	if track.BitRate == nil {
		return 0
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(*track.BitRate), 64)
	if err != nil || value <= 0 {
		return 0
	}
	return int32(value / 1000)
}

func normalizeAudioCodec(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "dd+", "ddp", "ddplus", "eac3":
		return "eac3"
	case "dd", "ac3", "dolbydigital":
		return "ac3"
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}

func normalizeVideoCodec(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "x264", "h264", "avc":
		return "h264"
	case "x265", "h265", "hevc":
		return "hevc"
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}

func stringListHasNormalized(values []string, candidate string) bool {
	candidate = strings.ToLower(strings.TrimSpace(candidate))
	for _, value := range values {
		if strings.ToLower(strings.TrimSpace(value)) == candidate {
			return true
		}
	}
	return false
}
