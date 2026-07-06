package decisions

import (
	"fmt"

	"media-manager/internal/storage"
)

type currentReleaseState struct {
	qualityID         string
	qualityScore      int32
	customFormatScore int32
}

func upgradeDecisionDetails(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
	qualityScore int32,
	customFormatScore int32,
) ([]string, string) {
	if len(item.FilePaths) == 0 {
		return []string{"No current file score is available."}, ""
	}
	current, ok := currentReleaseStateForItem(item, profile, formats)
	if !ok || current.qualityID == "" {
		return []string{"No current file score is available."}, ""
	}
	details := []string{fmt.Sprintf("Current file quality is %s.", current.qualityID)}
	if profile == nil {
		if qualityScore > current.qualityScore {
			return append(details, "Release quality is higher than the current file."), ""
		}
		return append(details, "Release quality is lower than or equal to the current file."), "quality not higher"
	}
	if !profile.UpgradesAllowed {
		return append(details, "Upgrades are disabled by the profile."), "upgrades disabled"
	}
	if reachedQualityTarget(current, profile) {
		return append(details, "Current file has reached the profile quality upgrade target."), "quality target reached"
	}
	if reachedCustomFormatTarget(current, profile) {
		return append(details, "Current file has reached the custom format upgrade target."), "custom format target reached"
	}
	if qualityScore > current.qualityScore {
		return append(details, "Release quality is higher than the current file."), ""
	}
	if qualityScore < current.qualityScore {
		return append(details, "Release quality is below the current file."), "quality downgrade"
	}
	increment := customFormatScore - current.customFormatScore
	minimumIncrement := profile.MinimumCustomFormatScoreIncrement
	if minimumIncrement < 0 {
		minimumIncrement = 0
	}
	if increment <= 0 {
		return append(details, "Custom format score is not higher than the current file."), "custom score not higher"
	}
	if increment < minimumIncrement {
		return append(
			details,
			fmt.Sprintf("Custom format score increment %d is below the profile minimum %d.", increment, minimumIncrement),
		), "custom score increment too low"
	}
	return append(details, fmt.Sprintf("Custom format score improves by %d.", increment)), ""
}

func currentReleaseStateForItem(
	item storage.MediaItem,
	profile *storage.MediaProfile,
	formats []storage.CustomFormat,
) (currentReleaseState, bool) {
	best := currentReleaseState{}
	found := false
	for _, path := range item.FilePaths {
		parsed := ParseReleaseFileName(path)
		qualityScore := profileQualityScore(parsed.QualityID, profile)
		customScore, _ := customFormatScore(parsed, profile, formats)
		if !found ||
			qualityScore > best.qualityScore ||
			(qualityScore == best.qualityScore && customScore > best.customFormatScore) {
			best = currentReleaseState{
				qualityID:         parsed.QualityID,
				qualityScore:      qualityScore,
				customFormatScore: customScore,
			}
			found = true
		}
	}
	return best, found
}

func reachedQualityTarget(current currentReleaseState, profile *storage.MediaProfile) bool {
	if profile == nil || profile.UpgradeUntilQualityID == nil {
		return false
	}
	targetScore := profileQualityScore(*profile.UpgradeUntilQualityID, profile)
	return targetScore > 0 && current.qualityScore >= targetScore
}

func reachedCustomFormatTarget(current currentReleaseState, profile *storage.MediaProfile) bool {
	if profile == nil || profile.UpgradeUntilCustomFormatScore <= 0 {
		return false
	}
	return current.customFormatScore >= profile.UpgradeUntilCustomFormatScore
}
