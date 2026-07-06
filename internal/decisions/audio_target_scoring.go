package decisions

import (
	"fmt"
	"strings"

	"media-manager/internal/storage"
)

func audioTargetScore(
	parsed ParsedRelease,
	profile *storage.MediaProfile,
) (int32, []ReleaseScoreContributor, string) {
	if profile == nil {
		return 0, nil, ""
	}
	var total int32
	contributors := []ReleaseScoreContributor{}
	for _, target := range profile.AudioTargets {
		score, matched := audioTargetDetailScore(parsed, target)
		total += score
		contributors = append(contributors, matched...)
		if reason := audioTargetBitrateReject(parsed, target); reason != "" {
			return total, contributors, reason
		}
	}
	return total, contributors, ""
}

func audioTargetDetailScore(
	parsed ParsedRelease,
	target storage.MediaProfileAudioTarget,
) (int32, []ReleaseScoreContributor) {
	var total int32
	contributors := []ReleaseScoreContributor{}
	add := func(label string) {
		if target.Score == 0 {
			return
		}
		total += target.Score
		contributors = append(contributors, ReleaseScoreContributor{Label: label, Score: target.Score})
	}
	if target.TargetCodec != nil && audioCodecMatches(parsed.AudioCodec, *target.TargetCodec) {
		add("Audio codec: " + normalizeAudioCodec(*target.TargetCodec))
	}
	if targetValueInList(parsed.AudioChannels, target.TargetChannels) {
		add("Audio channels: " + parsed.AudioChannels)
	}
	if bitrateMeetsTarget(parsed.AudioBitrateKbps, target) {
		add(fmt.Sprintf("Audio bitrate: %d kbps", parsed.AudioBitrateKbps))
	}
	return total, contributors
}

func audioTargetBitrateReject(
	parsed ParsedRelease,
	target storage.MediaProfileAudioTarget,
) string {
	if target.MinimumBitrateKbps == nil || parsed.AudioBitrateKbps == 0 {
		return ""
	}
	if parsed.AudioBitrateKbps >= *target.MinimumBitrateKbps {
		return ""
	}
	return fmt.Sprintf("Audio bitrate %d kbps is below the profile minimum.", parsed.AudioBitrateKbps)
}

func bitrateMeetsTarget(parsed int32, target storage.MediaProfileAudioTarget) bool {
	if parsed == 0 {
		return false
	}
	if target.PreferredBitrateKbps != nil {
		return parsed >= *target.PreferredBitrateKbps
	}
	if target.MinimumBitrateKbps != nil {
		return parsed >= *target.MinimumBitrateKbps
	}
	return false
}

func targetValueInList(parsed string, targets []string) bool {
	parsed = strings.ToLower(strings.TrimSpace(parsed))
	if parsed == "" || len(targets) == 0 {
		return false
	}
	for _, target := range targets {
		if parsed == strings.ToLower(strings.TrimSpace(target)) {
			return true
		}
	}
	return false
}

func audioCodecMatches(parsed, target string) bool {
	parsedCodec := normalizeAudioCodec(parsed)
	targetCodec := normalizeAudioCodec(target)
	if parsedCodec == "" || targetCodec == "" {
		return false
	}
	return parsedCodec == targetCodec
}

func normalizeAudioCodec(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), "DD+") {
		return "eac3"
	}
	switch normalizedToken(value) {
	case "ddp", "ddplus", "eac3":
		return "eac3"
	case "dd", "ac3", "dolbydigital":
		return "ac3"
	case "truehd", "truehdatmos":
		return "truehd"
	case "aac":
		return "aac"
	case "flac":
		return "flac"
	case "dts":
		return "dts"
	case "pcm":
		return "pcm"
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}
