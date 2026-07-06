package decisions

import "strings"

func detectVideoCodec(title string) string {
	switch {
	case containsAnyNormalized(title, "x266", "h266", "vvc"):
		return "x266"
	case containsAnyNormalized(title, "x265", "h265", "hevc"):
		return "x265"
	case containsAnyNormalized(title, "x264", "h264", "avc"):
		return "x264"
	case containsAnyNormalized(title, "av1"):
		return "AV1"
	case containsAnyNormalized(title, "xvid"):
		return "Xvid"
	case containsAnyNormalized(title, "vp9"):
		return "VP9"
	default:
		return ""
	}
}

func detectHDRFormat(title string) string {
	lowerTitle := strings.ToLower(title)
	switch {
	case containsAnyNormalized(title, "dolby vision", "dovi", "dv"):
		return "dolby-vision"
	case containsAnyNormalized(title, "hdr10plus") || strings.Contains(lowerTitle, "hdr10+"):
		return "hdr10plus"
	case containsAnyNormalized(title, "hdr10", "hdr"):
		return "hdr10"
	case containsAnyNormalized(title, "hlg"):
		return "hlg"
	case containsAnyNormalized(title, "sdr"):
		return "sdr"
	default:
		return ""
	}
}

func detectPixelFormat(title string) string {
	switch {
	case containsAnyNormalized(title, "hi10p", "10bit", "10-bit"):
		return "yuv420p10le"
	case containsAnyNormalized(title, "8bit", "8-bit"):
		return "yuv420p"
	default:
		return ""
	}
}
