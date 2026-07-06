package jobs

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
)

var subtitleEpisodePattern = regexp.MustCompile(`(?i)s(\d{1,2})e(\d{1,3})`)

func subtitleSearchRequestsForItem(item storage.MediaItem) []SubtitleSearchArgs {
	if !item.Monitored || len(item.FilePaths) == 0 || len(item.SubtitleTargets) == 0 {
		return nil
	}
	items := []SubtitleSearchArgs{}
	for _, target := range item.SubtitleTargets {
		if target.Source == "embedded" || !target.Required {
			continue
		}
		for _, path := range item.FilePaths {
			if externalSubtitleExists(item, target.LanguageID, path) {
				continue
			}
			items = append(items, SubtitleSearchArgs{
				MediaItemID: item.ID.String(),
				LanguageID:  target.LanguageID,
				FilePath:    path,
			})
		}
	}
	return items
}

func subtitleSearchRequest(
	item storage.MediaItem,
	args SubtitleSearchArgs,
) (subtitles.SearchRequest, bool) {
	language := strings.TrimSpace(args.LanguageID)
	filePath := strings.TrimSpace(args.FilePath)
	if filePath == "" && len(item.FilePaths) > 0 {
		filePath = item.FilePaths[0]
	}
	if language == "" {
		language = firstMissingSubtitleLanguage(item, filePath)
	}
	if language == "" || filePath == "" || externalSubtitleExists(item, language, filePath) {
		return subtitles.SearchRequest{}, false
	}
	request := subtitles.SearchRequest{
		MediaType:  item.Type,
		Title:      item.Title,
		LanguageID: language,
		Year:       item.Year,
		FilePath:   filePath,
	}
	season, episode, ok := subtitleEpisodeNumbers(filePath)
	if ok {
		request.SeasonNumber = &season
		request.EpisodeNumber = &episode
	}
	return request, true
}

func firstMissingSubtitleLanguage(item storage.MediaItem, filePath string) string {
	targets := append([]storage.MediaProfileSubtitleTarget(nil), item.SubtitleTargets...)
	sort.SliceStable(targets, func(i, j int) bool {
		return targets[i].LanguageID < targets[j].LanguageID
	})
	for _, target := range targets {
		if target.Source == "embedded" || !target.Required {
			continue
		}
		if !externalSubtitleExists(item, target.LanguageID, filePath) {
			return target.LanguageID
		}
	}
	return ""
}

func externalSubtitleExists(item storage.MediaItem, languageID string, filePath string) bool {
	for _, subtitle := range item.ExternalSubtitles {
		if languageMatchKey(subtitle.LanguageID) != languageMatchKey(languageID) {
			continue
		}
		if sameMediaBase(subtitle.FilePath, filePath) {
			return true
		}
	}
	return false
}

func sameMediaBase(subtitlePath string, filePath string) bool {
	subtitleBase := strings.TrimSuffix(filepath.Base(subtitlePath), filepath.Ext(subtitlePath))
	mediaBase := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	return strings.HasPrefix(strings.ToLower(subtitleBase), strings.ToLower(mediaBase)+".")
}

func subtitleEpisodeNumbers(filePath string) (int32, int32, bool) {
	matches := subtitleEpisodePattern.FindStringSubmatch(filepath.Base(filePath))
	if len(matches) != 3 {
		return 0, 0, false
	}
	season := parseSmallInt(matches[1])
	episode := parseSmallInt(matches[2])
	return season, episode, season > 0 && episode > 0
}

func subtitleSidecarPath(filePath string, languageID string) string {
	base := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	return base + "." + languageMatchKey(languageID) + ".srt"
}
