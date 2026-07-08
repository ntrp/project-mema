package satisfaction

import (
	"fmt"
	"path/filepath"
	"strings"

	"media-manager/internal/storage"
	"media-manager/internal/targets"
)

type VideoResult struct {
	Target             targets.Target
	Candidates         []targets.Candidate
	FailedRequirements []string
}

func EvaluateVideoTarget(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	fact storage.MediaFileFact,
) VideoResult {
	targetID := "video:" + fact.ID.String()
	result := VideoResult{
		Target: targets.Target{
			ID:          targetID,
			Type:        targets.TypeVideo,
			State:       targets.StateMissing,
			MediaItemID: item.ID.String(),
			MediaFileID: fact.ID.String(),
		},
	}
	video := firstVideoTrack(fact)
	if fact.FilePath == "" || video == nil {
		result.Target.Reasons = []string{"No persisted video track exists."}
		return result
	}
	candidate := targets.Candidate{
		ID:          candidateID(fact, *video),
		Type:        targets.CandidateVideoTrack,
		VisualState: targets.VisualMatching,
		TargetIDs:   []string{targetID},
	}
	failed := videoFailures(profile, fact, *video)
	result.FailedRequirements = failed
	if len(failed) > 0 {
		candidate.VisualState = targets.VisualPartial
		result.Candidates = []targets.Candidate{candidate}
		if operation := pendingVideoOperation(profile, fact, failed); operation != nil {
			candidate.VisualState = targets.VisualPendingOperation
			candidate.Operation = operation
			result.Candidates[0] = candidate
			result.Target.State = targets.StatePending
			result.Target.RequiredOperation = operation
			result.Target.Reasons = []string{operation.Reason}
			return result
		}
		result.Target.State = targets.StatePartial
		result.Target.Reasons = failed
		return result
	}
	result.Candidates = []targets.Candidate{candidate}
	if qualityUpgradeWanted(profile, fact.QualityID) {
		result.Target.State = targets.StateUpgradeable
		result.Target.Reasons = []string{"Profile quality upgrade target is higher than current file quality."}
		return result
	}
	result.Target.State = targets.StateSatisfied
	result.Target.Reasons = []string{"Video target is satisfied by persisted file facts."}
	return result
}

func firstVideoTrack(fact storage.MediaFileFact) *storage.MediaFileTrackFact {
	for index := range fact.Tracks {
		if fact.Tracks[index].TrackType == "video" {
			return &fact.Tracks[index]
		}
	}
	return nil
}

func videoFailures(
	profile *storage.MediaProfile,
	fact storage.MediaFileFact,
	track storage.MediaFileTrackFact,
) []string {
	if profile == nil {
		return nil
	}
	failed := []string{}
	target := profile.VideoTarget
	if len(profile.QualityIDs) > 0 && !stringListHas(profile.QualityIDs, stringPtrValue(fact.QualityID)) {
		failed = append(failed, "quality is not enabled in the profile")
	}
	if len(target.Codecs) > 0 && !stringListHasNormalized(target.Codecs, normalizeVideoCodec(stringPtrValue(track.Codec))) {
		failed = append(failed, "video codec does not meet the profile target")
	}
	if len(target.HDRFormats) > 0 && !stringListHasNormalized(target.HDRFormats, stringPtrValue(track.HDRFormat)) {
		failed = append(failed, "HDR format does not meet the profile target")
	}
	if len(target.PixelFormats) > 0 && !stringListHasNormalized(target.PixelFormats, stringPtrValue(track.PixelFormat)) {
		failed = append(failed, "pixel format does not meet the profile target")
	}
	if profile.FinalContainer != "" && profile.FinalContainer != containerExtension(fact) {
		failed = append(failed, "container does not meet the profile target")
	}
	return failed
}

func pendingVideoOperation(
	profile *storage.MediaProfile,
	fact storage.MediaFileFact,
	failed []string,
) *targets.Operation {
	if profile == nil || len(failed) != 1 || failed[0] != "container does not meet the profile target" {
		return nil
	}
	return &targets.Operation{
		Type:      targets.OperationContainerRemux,
		Manual:    true,
		Automatic: true,
		Reason:    fmt.Sprintf("Remux %s to %s container.", containerExtension(fact), profile.FinalContainer),
	}
}

func qualityUpgradeWanted(profile *storage.MediaProfile, qualityID *string) bool {
	if profile == nil || profile.UpgradeUntilQualityID == nil || qualityID == nil {
		return false
	}
	return qualityRank(profile, *qualityID) > 0 && qualityRank(profile, *qualityID) < qualityRank(profile, *profile.UpgradeUntilQualityID)
}

func qualityRank(profile *storage.MediaProfile, qualityID string) int {
	for index, value := range profile.QualityIDs {
		if value == qualityID {
			return index + 1
		}
	}
	return 0
}

func candidateID(fact storage.MediaFileFact, track storage.MediaFileTrackFact) string {
	return fmt.Sprintf("%s:stream:%d", fact.ID.String(), track.StreamIndex)
}

func containerExtension(fact storage.MediaFileFact) string {
	if value := strings.TrimPrefix(strings.ToLower(stringPtrValue(fact.ContainerFormat)), "."); value != "" {
		return value
	}
	return strings.TrimPrefix(strings.ToLower(filepath.Ext(fact.FilePath)), ".")
}

func stringListHas(values []string, candidate string) bool {
	for _, value := range values {
		if value == candidate {
			return true
		}
	}
	return false
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

func stringPtrValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
