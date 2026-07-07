package subtitles

import (
	"fmt"
	"strconv"
	"strings"
)

const mockSubtitleCueCount = 20

func (s *Service) searchMock(config Config, request SearchRequest) []Candidate {
	candidates := []Candidate{}
	for index, row := range config.MockSubtitles {
		if !mockSubtitleMatches(row, request) {
			continue
		}
		format := mockSubtitleFormat(row.Format)
		candidates = append(candidates, Candidate{
			ProviderName:  config.Name,
			LanguageID:    row.LanguageID,
			FileID:        int64(index + 1),
			Format:        format,
			ReleaseName:   strings.TrimSpace(row.Title) + "." + format,
			DownloadCount: 1,
			SourceURL:     "mock://subtitle/" + strconv.Itoa(index+1),
			SourceRef:     "mock:" + strconv.Itoa(index+1),
		})
	}
	return candidates
}

func (s *Service) downloadMock(candidate Candidate) Download {
	format := mockSubtitleFormat(candidate.Format)
	return Download{
		Content: mockSubtitleContent(format),
		URL:     strings.TrimSpace(candidate.SourceURL),
	}
}

func mockSubtitleMatches(row MockSubtitle, request SearchRequest) bool {
	return normalizedMockText(row.Title) == normalizedMockText(request.Title) &&
		languageMatchKey(row.LanguageID) == languageMatchKey(request.LanguageID)
}

func mockSubtitleFormat(value string) string {
	switch strings.TrimPrefix(strings.ToLower(strings.TrimSpace(value)), ".") {
	case "vtt", "webvtt":
		return "vtt"
	case "ass":
		return "ass"
	case "ssa":
		return "ssa"
	default:
		return "srt"
	}
}

func mockSubtitleContent(format string) []byte {
	switch mockSubtitleFormat(format) {
	case "vtt":
		return []byte(mockSubtitleVTT())
	case "ass":
		return []byte(mockSubtitleASS("ass"))
	case "ssa":
		return []byte(mockSubtitleASS("ssa"))
	default:
		return []byte(mockSubtitleSRT())
	}
}

func mockSubtitleSRT() string {
	blocks := make([]string, 0, mockSubtitleCueCount)
	for index := 0; index < mockSubtitleCueCount; index++ {
		start := index * 3
		blocks = append(blocks, fmt.Sprintf("%d\n%s --> %s\nmock", index+1, srtTime(start), srtTime(start+1)))
	}
	return strings.Join(blocks, "\n\n") + "\n"
}

func mockSubtitleVTT() string {
	blocks := []string{"WEBVTT"}
	for index := 0; index < mockSubtitleCueCount; index++ {
		start := index * 3
		blocks = append(blocks, fmt.Sprintf("%s --> %s\nmock", vttTime(start), vttTime(start+1)))
	}
	return strings.Join(blocks, "\n\n") + "\n"
}

func mockSubtitleASS(format string) string {
	lines := []string{
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
	for index := 0; index < mockSubtitleCueCount; index++ {
		start := index * 3
		lines = append(lines, fmt.Sprintf("Dialogue: %s,%s,Default,mock", assTime(start), assTime(start+1)))
	}
	if format == "ssa" {
		lines[4] = "Format: Name, Fontname, Fontsize, PrimaryColour, Alignment"
	}
	return strings.Join(lines, "\n") + "\n"
}

func srtTime(seconds int) string {
	return fmt.Sprintf("00:%02d:%02d,000", seconds/60, seconds%60)
}

func vttTime(seconds int) string {
	return strings.ReplaceAll(srtTime(seconds), ",", ".")
}

func assTime(seconds int) string {
	return fmt.Sprintf("0:%02d:%02d.00", seconds/60, seconds%60)
}

func normalizedMockText(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(value), " "))
}
