package satisfaction

import (
	"fmt"
	"strings"

	"media-manager/internal/storage"
	"media-manager/internal/targets"
)

type AudioEvaluation struct {
	Results    []AudioResult
	Candidates []targets.Candidate
}

type AudioResult struct {
	Target             targets.Target
	FailedRequirements []string
}

func EvaluateAudioTargets(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	fact storage.MediaFileFact,
) AudioEvaluation {
	if profile == nil || len(profile.AudioTargets) == 0 {
		return AudioEvaluation{}
	}
	evaluation := AudioEvaluation{}
	audioTracks := tracksByType(fact, "audio")
	targetLanguages := map[string]struct{}{}
	for _, target := range profile.AudioTargets {
		targetLanguages[languageMatchKey(target.LanguageID)] = struct{}{}
		result, candidates := evaluateAudioTarget(item, profile, fact, target, audioTracks)
		evaluation.Results = append(evaluation.Results, result)
		evaluation.Candidates = append(evaluation.Candidates, candidates...)
	}
	if profile.RemoveUnwantedAudio {
		evaluation.Candidates = append(evaluation.Candidates, unwantedAudioCandidates(fact, audioTracks, targetLanguages)...)
	}
	return evaluation
}

func evaluateAudioTarget(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	fact storage.MediaFileFact,
	target storage.MediaProfileAudioTarget,
	audioTracks []storage.MediaFileTrackFact,
) (AudioResult, []targets.Candidate) {
	targetID := "audio:" + fact.ID.String() + ":" + target.LanguageID
	result := AudioResult{Target: targets.Target{
		ID:          targetID,
		Type:        targets.TypeAudio,
		State:       targets.StateMissing,
		MediaItemID: item.ID.String(),
		MediaFileID: fact.ID.String(),
		LanguageID:  target.LanguageID,
	}}
	related := relatedAudioTracks(audioTracks, target.LanguageID)
	if len(related) == 0 {
		result.Target.Reasons = []string{"No persisted audio track for language " + target.LanguageID + "."}
		return result, nil
	}
	candidates := []targets.Candidate{}
	bestFailures := []string{}
	bestIndex := -1
	for index, track := range related {
		failures := audioFailures(target, track)
		visual := targets.VisualMatching
		if len(failures) > 0 {
			visual = targets.VisualPartial
		}
		candidates = append(candidates, targets.Candidate{
			ID:          candidateID(fact, track),
			Type:        targets.CandidateAudioTrack,
			VisualState: visual,
			TargetIDs:   []string{targetID},
			LanguageID:  target.LanguageID,
		})
		if len(failures) == 0 {
			result.Target.State = targets.StateSatisfied
			result.Target.Reasons = []string{"Audio target is satisfied by persisted track facts."}
			return result, candidates
		}
		if bestIndex < 0 || len(failures) < len(bestFailures) {
			bestFailures = failures
			bestIndex = index
		}
	}
	result.FailedRequirements = bestFailures
	if operation := pendingAudioOperation(profile, bestFailures); operation != nil && bestIndex >= 0 {
		candidates[bestIndex].VisualState = targets.VisualPendingOperation
		candidates[bestIndex].Operation = operation
		result.Target.State = targets.StatePending
		result.Target.RequiredOperation = operation
		result.Target.Reasons = []string{operation.Reason}
		return result, candidates
	}
	result.Target.State = targets.StatePartial
	result.Target.Reasons = bestFailures
	return result, candidates
}

func tracksByType(fact storage.MediaFileFact, trackType string) []storage.MediaFileTrackFact {
	tracks := []storage.MediaFileTrackFact{}
	for _, track := range fact.Tracks {
		if track.TrackType == trackType {
			tracks = append(tracks, track)
		}
	}
	return tracks
}

func relatedAudioTracks(tracks []storage.MediaFileTrackFact, languageID string) []storage.MediaFileTrackFact {
	related := []storage.MediaFileTrackFact{}
	for _, track := range tracks {
		if LanguageMatches(stringPtrValue(track.LanguageID), languageID) {
			related = append(related, track)
		}
	}
	return related
}

func audioFailures(target storage.MediaProfileAudioTarget, track storage.MediaFileTrackFact) []string {
	failures := []string{}
	if target.TargetCodec != nil && normalizeAudioCodec(stringPtrValue(track.Codec)) != normalizeAudioCodec(*target.TargetCodec) {
		failures = append(failures, "audio codec does not meet the profile target")
	}
	if len(target.TargetChannels) > 0 && !stringListHasNormalized(target.TargetChannels, stringPtrValue(track.Channels)) {
		failures = append(failures, "audio channels do not meet the profile target")
	}
	bitrateTarget := target.MinimumBitrateKbps
	if target.PreferredBitrateKbps != nil {
		bitrateTarget = target.PreferredBitrateKbps
	}
	if bitrateTarget != nil && (track.BitrateKbps == nil || *track.BitrateKbps < *bitrateTarget) {
		failures = append(failures, "audio bitrate does not meet the profile target")
	}
	return failures
}

func pendingAudioOperation(profile *storage.MediaProfile, failed []string) *targets.Operation {
	if profile == nil || profile.AudioLossyTranscodePolicy == "disabled" || len(failed) == 0 {
		return nil
	}
	return &targets.Operation{
		Type:      targets.OperationAudioTranscode,
		Manual:    true,
		Automatic: true,
		Reason:    fmt.Sprintf("Transcode audio to satisfy %s.", strings.Join(failed, ", ")),
	}
}

func unwantedAudioCandidates(
	fact storage.MediaFileFact,
	tracks []storage.MediaFileTrackFact,
	targetLanguages map[string]struct{},
) []targets.Candidate {
	candidates := []targets.Candidate{}
	for _, track := range tracks {
		language := stringPtrValue(track.LanguageID)
		if _, ok := targetLanguages[languageMatchKey(language)]; ok {
			continue
		}
		candidates = append(candidates, targets.Candidate{
			ID:            candidateID(fact, track),
			Type:          targets.CandidateAudioTrack,
			VisualState:   targets.VisualUnwanted,
			LanguageID:    language,
			UnwantedRules: []string{"remove-unwanted-audio"},
		})
	}
	return candidates
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
