package jobs

import (
	"regexp"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

var multiEpisodeReleasePattern = regexp.MustCompile(`(?i)\bs\d{1,2}e\d{1,3}(?:e\d{1,3}|\s*[-+&]\s*e?\d{1,3})`)

func markActiveEpisodeDownloads(item storage.MediaItem, activities []storage.DownloadActivity) (storage.MediaItem, bool) {
	for _, activity := range activities {
		if activity.MediaItemID != item.ID || !activeDownloadStatus(activity.Status) {
			continue
		}
		if multiEpisodeReleasePattern.MatchString(activity.ReleaseTitle) {
			return item, true
		}
		parsed := decisions.ParseReleaseFileName(activity.ReleaseTitle)
		if parsed.SeasonPack || decisions.IsSeriesPackTitle(activity.ReleaseTitle) {
			return item, true
		}
		if parsed.SeasonNumber == nil || parsed.EpisodeNumber == nil {
			return item, true
		}
		matched := false
		for _, season := range item.Seasons {
			seasonNumber, ok := seasonNumber(season)
			if !ok || seasonNumber != *parsed.SeasonNumber {
				continue
			}
			for _, episode := range season.Episodes {
				if episode.EpisodeNumber == *parsed.EpisodeNumber && episode.ID != nil {
					item.FileFacts = append(item.FileFacts, storage.MediaFileFact{
						SeasonID: season.ID, EpisodeID: episode.ID,
					})
					matched = true
					break
				}
			}
		}
		if !matched {
			return item, true
		}
	}
	return item, false
}

func activeDownloadStatus(status string) bool {
	return status == "queued" || status == "grabbed" || status == "downloading"
}
