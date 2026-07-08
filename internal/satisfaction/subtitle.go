package satisfaction

import (
	"fmt"

	"media-manager/internal/storage"
	"media-manager/internal/targets"
)

type SubtitleEvaluation struct {
	Results    []SubtitleResult
	Candidates []targets.Candidate
}

type SubtitleResult struct {
	Target             targets.Target
	FailedRequirements []string
}

type subtitleCandidateFact struct {
	Candidate targets.Candidate
	Format    string
}

func EvaluateSubtitleTargets(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	fact storage.MediaFileFact,
) SubtitleEvaluation {
	if profile == nil || len(profile.SubtitleTargets) == 0 {
		return SubtitleEvaluation{}
	}
	embedded := tracksByType(fact, "subtitle")
	evaluation := SubtitleEvaluation{}
	targetLanguages := map[string]struct{}{}
	for _, target := range profile.SubtitleTargets {
		targetLanguages[target.LanguageID] = struct{}{}
		result, candidates := evaluateSubtitleTarget(item, profile, fact, target, embedded)
		evaluation.Results = append(evaluation.Results, result)
		evaluation.Candidates = append(evaluation.Candidates, candidates...)
	}
	if profile.RemoveUnwantedSubtitles {
		evaluation.Candidates = append(evaluation.Candidates, unwantedSubtitleCandidates(item, fact, embedded, targetLanguages)...)
	}
	return evaluation
}

func evaluateSubtitleTarget(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	fact storage.MediaFileFact,
	target storage.MediaProfileSubtitleTarget,
	embedded []storage.MediaFileTrackFact,
) (SubtitleResult, []targets.Candidate) {
	targetID := "subtitle:" + fact.ID.String() + ":" + target.LanguageID
	result := SubtitleResult{Target: targets.Target{
		ID:          targetID,
		Type:        targets.TypeSubtitle,
		State:       targets.StateMissing,
		MediaItemID: item.ID.String(),
		MediaFileID: fact.ID.String(),
		LanguageID:  target.LanguageID,
	}}
	candidateFacts := subtitleCandidates(item, fact, targetID, target.LanguageID, embedded)
	candidates := subtitleTargetCandidates(candidateFacts)
	if len(candidates) == 0 {
		result.Target.Reasons = []string{"No persisted subtitle candidate for language " + target.LanguageID + "."}
		return result, nil
	}
	for index := range candidateFacts {
		failures, operation := subtitleCandidateFailures(profile, target, candidateFacts[index])
		if len(failures) == 0 {
			candidates[index].VisualState = targets.VisualMatching
			result.Target.State = targets.StateSatisfied
			result.Target.Reasons = []string{"Subtitle target is satisfied by persisted subtitle facts."}
			return result, candidates
		}
		candidates[index].VisualState = targets.VisualPartial
		if operation != nil {
			candidates[index].VisualState = targets.VisualPendingOperation
			candidates[index].Operation = operation
			result.Target.State = targets.StatePending
			result.Target.RequiredOperation = operation
			result.Target.Reasons = []string{operation.Reason}
			result.FailedRequirements = failures
			return result, candidates
		}
		if len(result.FailedRequirements) == 0 || len(failures) < len(result.FailedRequirements) {
			result.FailedRequirements = failures
		}
	}
	result.Target.State = targets.StatePartial
	result.Target.Reasons = result.FailedRequirements
	return result, candidates
}

func subtitleCandidates(
	item storage.MediaItem,
	fact storage.MediaFileFact,
	targetID string,
	languageID string,
	embedded []storage.MediaFileTrackFact,
) []subtitleCandidateFact {
	candidates := []subtitleCandidateFact{}
	for _, track := range embedded {
		if stringPtrValue(track.LanguageID) != languageID {
			continue
		}
		candidates = append(candidates, subtitleCandidateFact{
			Candidate: targets.Candidate{
				ID:          candidateID(fact, track),
				Type:        targets.CandidateEmbeddedSubtitle,
				VisualState: targets.VisualPartial,
				TargetIDs:   []string{targetID},
				LanguageID:  languageID,
			},
			Format: stringPtrValue(track.Format),
		})
	}
	for _, subtitle := range item.ExternalSubtitles {
		if !subtitle.Selected || subtitle.LanguageID != languageID {
			continue
		}
		candidates = append(candidates, subtitleCandidateFact{
			Candidate: targets.Candidate{
				ID:          subtitle.ID.String(),
				Type:        targets.CandidateExternalSubtitle,
				VisualState: targets.VisualPartial,
				TargetIDs:   []string{targetID},
				LanguageID:  languageID,
			},
			Format: subtitle.Format,
		})
	}
	return candidates
}

func subtitleCandidateFailures(
	profile *storage.MediaProfile,
	target storage.MediaProfileSubtitleTarget,
	candidate subtitleCandidateFact,
) ([]string, *targets.Operation) {
	mode := profile.SubtitleMode
	if mode == "" {
		mode = "mixed"
	}
	if candidate.Candidate.Type == targets.CandidateExternalSubtitle && mode == "embedded" {
		return []string{"subtitle must be embedded"}, &targets.Operation{
			Type:      targets.OperationSubtitleEmbed,
			Manual:    true,
			Automatic: true,
			Reason:    "Embed external subtitle for embedded subtitle mode.",
		}
	}
	if candidate.Candidate.Type == targets.CandidateEmbeddedSubtitle && mode == "external" {
		return []string{"subtitle must be external"}, &targets.Operation{
			Type:      targets.OperationSubtitleExtraction,
			Manual:    true,
			Automatic: true,
			Reason:    "Extract embedded subtitle for external subtitle mode.",
		}
	}
	if len(target.Formats) > 0 && !stringListHasNormalized(target.Formats, candidate.Format) {
		return []string{"subtitle format does not meet the profile target"}, &targets.Operation{
			Type:      targets.OperationSubtitleConversion,
			Manual:    true,
			Automatic: true,
			Reason:    fmt.Sprintf("Convert subtitle to one of %v.", target.Formats),
		}
	}
	return nil, nil
}

func subtitleTargetCandidates(facts []subtitleCandidateFact) []targets.Candidate {
	candidates := make([]targets.Candidate, 0, len(facts))
	for _, fact := range facts {
		candidates = append(candidates, fact.Candidate)
	}
	return candidates
}

func unwantedSubtitleCandidates(
	item storage.MediaItem,
	fact storage.MediaFileFact,
	embedded []storage.MediaFileTrackFact,
	targetLanguages map[string]struct{},
) []targets.Candidate {
	candidates := []targets.Candidate{}
	for _, track := range embedded {
		language := stringPtrValue(track.LanguageID)
		if _, ok := targetLanguages[language]; !ok {
			candidates = append(candidates, unwantedSubtitleCandidate(candidateID(fact, track), targets.CandidateEmbeddedSubtitle, language))
		}
	}
	for _, subtitle := range item.ExternalSubtitles {
		if _, ok := targetLanguages[subtitle.LanguageID]; !ok {
			candidates = append(candidates, unwantedSubtitleCandidate(subtitle.ID.String(), targets.CandidateExternalSubtitle, subtitle.LanguageID))
		}
	}
	return candidates
}

func unwantedSubtitleCandidate(id string, candidateType targets.CandidateType, language string) targets.Candidate {
	return targets.Candidate{
		ID:            id,
		Type:          candidateType,
		VisualState:   targets.VisualUnwanted,
		LanguageID:    language,
		UnwantedRules: []string{"remove-unwanted-subtitles"},
	}
}
