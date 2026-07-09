package jobs

import "media-manager/internal/storage"

func outputSubtitleIndex(item storage.MediaItem, filePath string) int {
	count := 0
	for _, fact := range item.FileFacts {
		if fact.FilePath != filePath {
			continue
		}
		for _, track := range fact.Tracks {
			if track.TrackType == "subtitle" {
				count++
			}
		}
	}
	return count
}

func ffmpegSubtitleLanguageTag(languageID string) string {
	switch languageMatchKey(languageID) {
	case "english":
		return "eng"
	case "german":
		return "ger"
	case "french":
		return "fre"
	case "spanish":
		return "spa"
	case "japanese":
		return "jpn"
	default:
		normalized := languageMatchKey(languageID)
		if len(normalized) == 3 {
			return normalized
		}
		return ""
	}
}
