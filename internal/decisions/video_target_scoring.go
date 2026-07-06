package decisions

import (
	"fmt"
	"strings"

	"media-manager/internal/storage"
)

func videoTargetScore(
	parsed ParsedRelease,
	profile *storage.MediaProfile,
) (int32, []ReleaseScoreContributor, string) {
	if profile == nil {
		return 0, nil, ""
	}
	target := profile.VideoTarget
	var total int32
	contributors := []ReleaseScoreContributor{}
	addScore := func(label string, score int32) {
		if score == 0 {
			return
		}
		total += score
		contributors = append(contributors, ReleaseScoreContributor{Label: label, Score: score})
	}
	codec := normalizeVideoCodec(parsed.VideoCodec)
	if matched, reason := targetValueMatches("Video codec", codec, target.Codecs, target.CodecRequired); reason != "" {
		return total, contributors, reason
	} else if matched {
		addScore("Video codec: "+codec, target.CodecScore)
	}
	if matched, reason := targetValueMatches("HDR format", parsed.HDRFormat, target.HDRFormats, target.HDRRequired); reason != "" {
		return total, contributors, reason
	} else if matched {
		addScore("HDR format: "+parsed.HDRFormat, target.HDRScore)
	}
	if matched, reason := targetValueMatches("Pixel format", parsed.PixelFormat, target.PixelFormats, target.PixelFormatRequired); reason != "" {
		return total, contributors, reason
	} else if matched {
		addScore("Pixel format: "+parsed.PixelFormat, target.PixelFormatScore)
	}
	return total, contributors, ""
}

func targetValueMatches(label, parsed string, targets []string, required bool) (bool, string) {
	parsed = strings.ToLower(strings.TrimSpace(parsed))
	if parsed == "" || len(targets) == 0 {
		return false, ""
	}
	for _, target := range targets {
		if parsed == strings.ToLower(strings.TrimSpace(target)) {
			return true, ""
		}
	}
	if required {
		return false, fmt.Sprintf("%s %s does not meet the profile target.", label, parsed)
	}
	return false, ""
}

func normalizeVideoCodec(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "x264", "h264", "avc":
		return "h264"
	case "x265", "h265", "hevc":
		return "hevc"
	case "x266", "h266", "vvc":
		return "vvc"
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}
