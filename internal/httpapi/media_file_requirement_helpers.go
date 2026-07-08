package httpapi

import (
	"fmt"
	"path/filepath"
	"strings"

	"media-manager/internal/storage"
)

func detailState(visual MediaFileDetailVisualState, label string, details ...string) *MediaFileDetailState {
	return &MediaFileDetailState{VisualState: visual, StatusLabel: label, Details: compactDetails(details)}
}

func operationDetailState(label string, detail string) *MediaFileDetailState {
	return &MediaFileDetailState{
		VisualState:    MediaFileDetailVisualStatePendingOperation,
		StatusLabel:    "Pending",
		OperationLabel: &label,
		Details:        []string{detail},
	}
}

func requirementStatus(state MediaFileRequirementState, label string, details ...string) MediaFileRequirementStatus {
	return MediaFileRequirementStatus{State: state, Label: label, Details: compactDetails(details)}
}

func requirementLabel(state MediaFileRequirementState) string {
	if state == MediaFileRequirementStateMissing {
		return "Missing"
	}
	if state == MediaFileRequirementStatePending {
		return "Pending"
	}
	return "Partial"
}

func compactDetails(values []string) []string {
	result := []string{}
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			result = append(result, value)
		}
	}
	if len(result) == 0 {
		return []string{"-"}
	}
	return result
}

func videoFailures(item storage.MediaItem, path string, track MediaFileTrack) []string {
	failures := []string{}
	target := item.VideoTarget
	if len(target.Codecs) > 0 && !stringListHasNormalized(target.Codecs, normalizeVideoCodec(optionalStringValue(track.Codec))) {
		failures = append(failures, "Video codec does not meet the profile target")
	}
	if failure := videoResolutionFailure(item, path, track); failure != "" {
		failures = append(failures, failure)
	}
	if len(target.PixelFormats) > 0 && !stringListHasNormalized(target.PixelFormats, optionalStringValue(track.PixelFormat)) {
		failures = append(failures, "Pixel format does not meet the profile target")
	}
	return failures
}

func videoResolutionFailure(item storage.MediaItem, path string, track MediaFileTrack) string {
	qualityID := mediaFileQualityID(item, path)
	if qualityID == "" {
		return ""
	}
	bounds, ok := storage.QualityResolutionForID(qualityID)
	if !ok {
		return ""
	}
	if track.Width == nil && track.Height == nil {
		return "Video resolution is unknown for the selected quality"
	}
	if trackResolutionMatchesQuality(track.Width, track.Height, bounds) {
		return ""
	}
	return fmt.Sprintf("Video resolution is below selected quality %dp", bounds.MinHeight)
}

func trackResolutionMatchesQuality(width *int32, height *int32, bounds storage.QualityResolutionBounds) bool {
	if width != nil && *width >= bounds.MinWidth {
		return true
	}
	return height != nil && *height >= bounds.MinHeight
}

func mediaFileQualityID(item storage.MediaItem, path string) string {
	for _, fact := range item.FileFacts {
		if sameMediaFilePath(fact.FilePath, path) && fact.QualityID != nil && strings.TrimSpace(*fact.QualityID) != "" {
			return strings.TrimSpace(*fact.QualityID)
		}
	}
	return storage.QualityIDFromPath(path)
}

func sameMediaFilePath(left string, right string) bool {
	return filepath.Clean(left) == filepath.Clean(right)
}

func audioTrackFailures(track MediaFileTrack, target storage.MediaProfileAudioTarget) []string {
	failures := []string{}
	if target.TargetCodec != nil && normalizeAudioCodec(optionalStringValue(track.Codec)) != normalizeAudioCodec(*target.TargetCodec) {
		failures = append(failures, "Audio codec does not meet the profile target")
	}
	if len(target.TargetChannels) > 0 && !stringListHasNormalized(target.TargetChannels, audioChannelValue(track)) {
		failures = append(failures, "Audio channels do not meet the profile target")
	}
	if target.MinimumBitrateKbps != nil && trackBitrateKbps(track) < *target.MinimumBitrateKbps {
		failures = append(failures, fmt.Sprintf("Audio bitrate is below %d kbps", *target.MinimumBitrateKbps))
	}
	return failures
}

func firstTrackOfType(tracks []MediaFileTrack, trackType MediaFileTrackType) *MediaFileTrack {
	for index := range tracks {
		if tracks[index].Type == trackType {
			return &tracks[index]
		}
	}
	return nil
}

func tracksOfType(tracks []MediaFileTrack, trackType MediaFileTrackType) []MediaFileTrack {
	result := []MediaFileTrack{}
	for _, track := range tracks {
		if track.Type == trackType {
			result = append(result, track)
		}
	}
	return result
}

func audioTargetsForTrack(targets []storage.MediaProfileAudioTarget, language string) []storage.MediaProfileAudioTarget {
	result := []storage.MediaProfileAudioTarget{}
	for _, target := range targets {
		if languageMatchKey(language) == languageMatchKey(target.LanguageID) {
			result = append(result, target)
		}
	}
	return result
}

func audioTargetsCandidates(tracks []MediaFileTrack, language string) []MediaFileTrack {
	result := []MediaFileTrack{}
	for _, track := range tracks {
		if languageMatchKey(optionalStringValue(track.Language)) == languageMatchKey(language) {
			result = append(result, track)
		}
	}
	return result
}

func anyAudioTrackMatches(tracks []MediaFileTrack, target storage.MediaProfileAudioTarget) bool {
	for _, track := range tracks {
		if len(audioTrackFailures(track, target)) == 0 {
			return true
		}
	}
	return false
}

func audioTargetLanguages(targets []storage.MediaProfileAudioTarget) []string {
	languages := make([]string, 0, len(targets))
	for _, target := range targets {
		languages = append(languages, target.LanguageID)
	}
	return languages
}

func subtitleTargetLanguages(targets []storage.MediaProfileSubtitleTarget) []string {
	languages := make([]string, 0, len(targets))
	for _, target := range targets {
		languages = append(languages, target.LanguageID)
	}
	return languages
}
