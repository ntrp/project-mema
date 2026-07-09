package httpapi

import (
	"path/filepath"
	"strings"

	"media-manager/internal/storage"
)

func applyMediaFileRequirementStates(
	file *MediaFileInfo,
	item storage.MediaItem,
	tracks []MediaFileTrack,
	otherFiles []MediaFileOtherFile,
) {
	if file.Status == MediaFileInfoStatusMissing {
		file.Requirements = &MediaFileRequirementSummary{
			Video:     requirementStatus(MediaFileRequirementStateMissing, "Missing", "File is missing"),
			Container: requirementStatus(MediaFileRequirementStateMissing, "Missing", "File is missing"),
			Audio:     requirementStatus(MediaFileRequirementStateMissing, "Missing", "File is missing"),
			Subtitles: subtitleRequirementStatus(file.SubtitleSatisfaction, nil, nil),
		}
		return
	}
	for index := range tracks {
		tracks[index].State = trackDetailState(item, file.Path, tracks[index])
	}
	for index := range otherFiles {
		otherFiles[index].State = otherFileDetailState(item, file.Path, otherFiles[index])
	}
	missing := missingRequirementRows(item, file.SubtitleSatisfaction, tracks, otherFiles)
	if len(missing) > 0 {
		file.MissingTracks = &missing
	}
	file.Requirements = &MediaFileRequirementSummary{
		Video:     videoRequirementStatus(item, file.Path, tracks),
		Container: containerRequirementStatus(item, file.Path),
		Audio:     audioRequirementStatus(item, tracks),
		Subtitles: subtitleRequirementStatus(file.SubtitleSatisfaction, tracks, otherFiles),
	}
}

func trackDetailState(item storage.MediaItem, path string, track MediaFileTrack) *MediaFileDetailState {
	switch track.Type {
	case MediaFileTrackTypeVideo:
		failures := videoTrackFailures(item, path, track)
		if len(failures) > 0 {
			return detailState(MediaFileDetailVisualStatePartial, "Partial", failures...)
		}
		return detailState(MediaFileDetailVisualStateMatching, "Matching", "Video track satisfies the profile target.")
	case MediaFileTrackTypeAudio:
		return audioTrackDetailState(item, track)
	case MediaFileTrackTypeSubtitle:
		return embeddedSubtitleDetailState(item, track)
	default:
		return nil
	}
}

func audioTrackDetailState(item storage.MediaItem, track MediaFileTrack) *MediaFileDetailState {
	if outsideTargetLanguage(audioTargetLanguages(item.AudioTargets), optionalStringValue(track.Language)) {
		return detailState(MediaFileDetailVisualStateUnwanted, "Unwanted", "Audio language is outside enabled profile targets.")
	}
	targets := audioTargetsForTrack(item.AudioTargets, optionalStringValue(track.Language))
	if len(item.AudioTargets) == 0 {
		return detailState(MediaFileDetailVisualStateMatching, "Matching", "Audio track is available.")
	}
	if len(targets) == 0 {
		return nil
	}
	failures := audioTrackFailures(track, targets[0])
	if len(failures) == 0 {
		return detailState(MediaFileDetailVisualStateMatching, "Matching", "Audio track satisfies a profile target.")
	}
	return detailState(MediaFileDetailVisualStatePartial, "Partial", failures...)
}

func embeddedSubtitleDetailState(item storage.MediaItem, track MediaFileTrack) *MediaFileDetailState {
	language := optionalStringValue(track.Language)
	if mediaFileSubtitleMode(item.SubtitleMode) == MediaProfileSubtitleModeExternal {
		if subtitleTargetMatches(item.SubtitleTargets, language) {
			return operationDetailState("Extract subtitle", "Embedded subtitle can satisfy the target after extraction.")
		}
		return detailState(MediaFileDetailVisualStateUnwanted, "Unwanted", "Embedded subtitles conflict with external subtitle mode.")
	}
	if unwantedLanguage(item.RemoveUnwantedSubtitles, subtitleTargetLanguages(item.SubtitleTargets), language) {
		return detailState(MediaFileDetailVisualStateUnwanted, "Unwanted", "Subtitle language is outside enabled profile targets.")
	}
	if outsideTargetLanguage(subtitleTargetLanguages(item.SubtitleTargets), language) {
		return detailState(MediaFileDetailVisualStateUnwanted, "Unwanted", "Subtitle language is outside enabled profile targets.")
	}
	if len(item.SubtitleTargets) == 0 || subtitleTargetMatches(item.SubtitleTargets, language) {
		if !subtitleTargetFormatMatches(item.SubtitleTargets, language, optionalStringValue(track.Codec)) {
			return subtitleFormatMismatchDetailState(item.SubtitleTargets, language, optionalStringValue(track.Codec))
		}
		return detailState(MediaFileDetailVisualStateMatching, "Matching", "Embedded subtitle satisfies the subtitle target.")
	}
	return nil
}

func otherFileDetailState(item storage.MediaItem, mediaPath string, file MediaFileOtherFile) *MediaFileDetailState {
	if file.Type != MediaFileOtherFileTypeSubtitle {
		return nil
	}
	language := optionalStringValue(file.Language)
	if file.Status == MediaFileOtherFileStatusMissing {
		return detailState(MediaFileDetailVisualStateMissingPlaceholder, "Missing", "Missing expected external subtitle: "+language)
	}
	if unwantedLanguage(item.RemoveUnwantedSubtitles, subtitleTargetLanguages(item.SubtitleTargets), language) {
		return detailState(MediaFileDetailVisualStateUnwanted, "Unwanted", "Subtitle language is outside enabled profile targets.")
	}
	if outsideTargetLanguage(subtitleTargetLanguages(item.SubtitleTargets), language) {
		return detailState(MediaFileDetailVisualStateUnwanted, "Unwanted", "Subtitle language is outside enabled profile targets.")
	}
	if !subtitleTargetMatches(item.SubtitleTargets, language) {
		return nil
	}
	if mediaFileSubtitleMode(item.SubtitleMode) == MediaProfileSubtitleModeEmbedded {
		return operationDetailState("Embed subtitle", "External subtitle can satisfy the target after embedding.")
	}
	format := strings.TrimPrefix(filepath.Ext(file.Path), ".")
	if !subtitleTargetFormatMatches(item.SubtitleTargets, language, format) {
		return subtitleFormatMismatchDetailState(item.SubtitleTargets, language, format)
	}
	return detailState(MediaFileDetailVisualStateMatching, "Matching", "External subtitle satisfies the subtitle target.")
}

func subtitleFormatMismatchDetailState(
	targets []storage.MediaProfileSubtitleTarget,
	language string,
	format string,
) *MediaFileDetailState {
	if subtitleTargetTextConversionSupported(targets, language, format) {
		return operationDetailState("Convert subtitle", "Subtitle format does not meet the profile target.")
	}
	return detailState(
		MediaFileDetailVisualStatePartial,
		"Partial",
		"Subtitle format requires non-text conversion support.",
	)
}

func videoRequirementStatus(item storage.MediaItem, path string, tracks []MediaFileTrack) MediaFileRequirementStatus {
	video := firstTrackOfType(tracks, MediaFileTrackTypeVideo)
	if video == nil {
		return requirementStatus(MediaFileRequirementStateMissing, "Missing", "Video track is missing")
	}
	failures := videoTrackFailures(item, path, *video)
	if len(failures) > 0 {
		return requirementStatus(MediaFileRequirementStatePartial, "Partial", failures...)
	}
	return requirementStatus(MediaFileRequirementStateSatisfied, "Ok", "Video requirements met")
}

func containerRequirementStatus(item storage.MediaItem, path string) MediaFileRequirementStatus {
	if item.FinalContainer == "" {
		return requirementStatus(MediaFileRequirementStateIgnored, "Ignored", "No container target")
	}
	container := mediaFileContainer(item, path)
	if container == "" {
		return requirementStatus(MediaFileRequirementStateMissing, "Missing", "Container format is unknown")
	}
	if videoContainerMismatch(item, path) {
		return requirementStatus(MediaFileRequirementStatePending, "Pending", "Container does not meet the profile target")
	}
	return requirementStatus(MediaFileRequirementStateSatisfied, "Ok", "Container requirements met")
}

func audioRequirementStatus(item storage.MediaItem, tracks []MediaFileTrack) MediaFileRequirementStatus {
	audio := tracksOfType(tracks, MediaFileTrackTypeAudio)
	if len(audio) == 0 {
		return requirementStatus(MediaFileRequirementStateMissing, "Missing", "Audio track is missing")
	}
	issues := []string{}
	missingTargets := 0
	for _, target := range item.AudioTargets {
		candidates := audioTargetsCandidates(audio, target.LanguageID)
		if len(candidates) == 0 {
			missingTargets++
			issues = append(issues, "Missing required audio: "+target.LanguageID)
			continue
		}
		if !anyAudioTrackMatches(candidates, target) {
			issues = append(issues, audioTrackFailures(candidates[0], target)...)
		}
	}
	if len(issues) > 0 {
		state := MediaFileRequirementStatePartial
		if len(item.AudioTargets) > 0 && missingTargets == len(item.AudioTargets) {
			state = MediaFileRequirementStateMissing
		}
		return requirementStatus(state, requirementLabel(state), issues...)
	}
	if unwanted := unwantedAudioTracks(item, audio); len(unwanted) > 0 {
		return requirementStatus(MediaFileRequirementStatePartial, "Partial", "Unwanted audio tracks: "+strings.Join(unwanted, ", "))
	}
	return requirementStatus(MediaFileRequirementStateSatisfied, "Ok", "Audio requirements met")
}

func subtitleRequirementStatus(
	satisfaction *MediaFileSubtitleSatisfaction,
	tracks []MediaFileTrack,
	otherFiles []MediaFileOtherFile,
) MediaFileRequirementStatus {
	if satisfaction == nil || satisfaction.State == MediaFileSubtitleSatisfactionStateIgnored {
		return requirementStatus(MediaFileRequirementStateIgnored, "Ignored", "Subtitle requirements are ignored")
	}
	if satisfaction.State == MediaFileSubtitleSatisfactionStateSatisfied {
		return requirementStatus(MediaFileRequirementStateSatisfied, "Ok", "Subtitle requirements met")
	}
	if hasPendingSubtitleCandidate(tracks, otherFiles) {
		return requirementStatus(MediaFileRequirementStatePending, "Pending", "Subtitle candidate needs a follow-up operation")
	}
	if len(satisfaction.MatchedLanguages) > 0 {
		return requirementStatus(MediaFileRequirementStatePartial, "Partial", "Missing subtitles: "+strings.Join(satisfaction.MissingLanguages, ", "))
	}
	return requirementStatus(MediaFileRequirementStateMissing, "Missing", "Missing subtitles: "+strings.Join(satisfaction.MissingLanguages, ", "))
}

func missingRequirementRows(item storage.MediaItem, satisfaction *MediaFileSubtitleSatisfaction, tracks []MediaFileTrack, otherFiles []MediaFileOtherFile) []MediaFileMissingTrack {
	rows := []MediaFileMissingTrack{}
	for _, target := range item.AudioTargets {
		if hasTrackLanguage(tracks, MediaFileTrackTypeAudio, target.LanguageID) {
			continue
		}
		rows = append(rows, MediaFileMissingTrack{
			Key:         "missing-audio-" + languageMatchKey(target.LanguageID),
			Type:        MediaFileMissingTrackTypeAudio,
			Language:    target.LanguageID,
			Description: "Missing expected audio track",
			State:       *detailState(MediaFileDetailVisualStateMissingPlaceholder, "Missing", "Missing expected audio: "+target.LanguageID),
		})
	}
	if satisfaction == nil {
		return rows
	}
	for _, language := range satisfaction.MissingLanguages {
		if hasSubtitleCandidate(tracks, otherFiles, language) {
			continue
		}
		rows = append(rows, MediaFileMissingTrack{
			Key:         "missing-subtitle-" + languageMatchKey(language),
			Type:        MediaFileMissingTrackTypeSubtitle,
			Language:    language,
			Description: "Missing expected subtitle track",
			State:       *detailState(MediaFileDetailVisualStateMissingPlaceholder, "Missing", "Missing expected subtitle: "+language),
		})
	}
	return rows
}
