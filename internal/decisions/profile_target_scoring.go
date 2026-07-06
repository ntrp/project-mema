package decisions

import "media-manager/internal/storage"

func profileTargetScore(
	parsed ParsedRelease,
	profile *storage.MediaProfile,
) (int32, []ReleaseScoreContributor, string) {
	videoScore, videoContributors, videoReject := videoTargetScore(parsed, profile)
	if videoReject != "" {
		return videoScore, videoContributors, videoReject
	}
	audioScore, audioContributors, audioReject := audioTargetScore(parsed, profile)
	total := videoScore + audioScore
	contributors := append([]ReleaseScoreContributor{}, videoContributors...)
	contributors = append(contributors, audioContributors...)
	return total, contributors, audioReject
}
