package jobs

import (
	"fmt"
	"strings"

	"media-manager/internal/storage"
	"media-manager/internal/subtitleformats"
)

func subtitleTargetFormat(item storage.MediaItem, languageID string) string {
	for _, target := range item.SubtitleTargets {
		if languageMatchKey(target.LanguageID) != languageMatchKey(languageID) {
			continue
		}
		for _, format := range target.Formats {
			if normalized := subtitleformats.Normalize(format); normalized != "" {
				return normalized
			}
		}
	}
	return ""
}

func convertSubtitleContent(content []byte, targetFormat string) ([]byte, string, error) {
	format := subtitleformats.Normalize(targetFormat)
	switch format {
	case "subrip":
		return content, format, nil
	case "vtt":
		return srtToVTT(content), format, nil
	case "ass", "ssa":
		if assSubtitleContent(content) {
			return content, format, nil
		}
		return srtToASS(content, format), format, nil
	case "pgs":
		return nil, "", fmt.Errorf("subtitle format pgs requires bitmap subtitle extraction and cannot be converted from provider text")
	default:
		return nil, "", fmt.Errorf("subtitle format %q is unsupported", targetFormat)
	}
}

func assSubtitleContent(content []byte) bool {
	body := strings.TrimSpace(string(content))
	return strings.HasPrefix(body, "[Script Info]") && strings.Contains(body, "[Events]")
}

func srtToVTT(content []byte) []byte {
	body := strings.ReplaceAll(string(content), ",", ".")
	if strings.HasPrefix(strings.TrimSpace(body), "WEBVTT") {
		return []byte(body)
	}
	return []byte("WEBVTT\n\n" + body)
}

func srtToASS(content []byte, format string) []byte {
	events := []string{
		"[Script Info]",
		"ScriptType: v4.00+",
		"",
		"[V4+ Styles]",
		"Format: Name, Fontname, Fontsize, PrimaryColour, Alignment",
		"Style: Default,Arial,20,&H00FFFFFF,2",
		"",
		"[Events]",
		"Format: Start, End, Style, Text",
	}
	for _, block := range strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n\n") {
		lines := nonEmptyStrings(strings.Split(block, "\n"))
		if len(lines) < 3 || !strings.Contains(lines[1], " --> ") {
			continue
		}
		times := strings.SplitN(lines[1], " --> ", 2)
		text := strings.Join(lines[2:], `\N`)
		events = append(events, fmt.Sprintf(
			"Dialogue: %s,%s,Default,%s",
			assTimestamp(times[0]),
			assTimestamp(times[1]),
			text,
		))
	}
	if format == "ssa" {
		events[0] = "[Script Info]"
	}
	return []byte(strings.Join(events, "\n") + "\n")
}

func assTimestamp(value string) string {
	value = strings.TrimSpace(strings.ReplaceAll(value, ",", "."))
	parts := strings.Split(value, ".")
	if len(parts) == 1 {
		return parts[0] + ".00"
	}
	fraction := parts[1]
	if len(fraction) > 2 {
		fraction = fraction[:2]
	}
	for len(fraction) < 2 {
		fraction += "0"
	}
	return parts[0] + "." + fraction
}
