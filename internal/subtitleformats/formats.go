package subtitleformats

import "strings"

func Normalize(value string) string {
	switch strings.TrimPrefix(strings.ToLower(strings.TrimSpace(value)), ".") {
	case "", "srt", "subrip":
		return "subrip"
	case "vtt", "webvtt":
		return "vtt"
	case "ass":
		return "ass"
	case "ssa":
		return "ssa"
	case "pgs", "sup":
		return "pgs"
	default:
		return ""
	}
}

func Match(target string, candidate string) bool {
	target = normalizedOrClean(target)
	candidate = normalizedOrClean(candidate)
	return target != "" && target == candidate
}

func AnyMatch(targets []string, candidate string) bool {
	if len(targets) == 0 {
		return true
	}
	for _, target := range targets {
		if Match(target, candidate) {
			return true
		}
	}
	return false
}

func Text(value string) bool {
	switch Normalize(value) {
	case "subrip", "vtt", "ass", "ssa":
		return true
	default:
		return false
	}
}

func HasTextTarget(targets []string) bool {
	for _, target := range targets {
		if Text(target) {
			return true
		}
	}
	return false
}

func normalizedOrClean(value string) string {
	if normalized := Normalize(value); normalized != "" {
		return normalized
	}
	return strings.TrimPrefix(strings.ToLower(strings.TrimSpace(value)), ".")
}
