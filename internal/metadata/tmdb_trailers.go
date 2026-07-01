package metadata

import "strings"

func tmdbTrailerURL(videos tmdbVideos) *string {
	key := tmdbTrailerKey(videos)
	if key == "" {
		return nil
	}
	value := "https://www.youtube.com/watch?v=" + key
	return &value
}

func tmdbTrailerKey(videos tmdbVideos) string {
	firstTrailer := ""
	for _, video := range videos.Results {
		if !tmdbIsYouTubeTrailer(video) {
			continue
		}
		key := strings.TrimSpace(video.Key)
		if key == "" {
			continue
		}
		if video.Official {
			return key
		}
		if firstTrailer == "" {
			firstTrailer = key
		}
	}
	return firstTrailer
}

func tmdbIsYouTubeTrailer(video tmdbVideo) bool {
	return strings.EqualFold(strings.TrimSpace(video.Site), "YouTube") &&
		strings.EqualFold(strings.TrimSpace(video.Type), "Trailer")
}
