package httpapi

import "media-manager/internal/storage"

func storageMediaSeasons(values *[]MediaMetadataSeason) *[]storage.MediaSeason {
	if values == nil {
		return nil
	}
	items := make([]storage.MediaSeason, 0, len(*values))
	for _, value := range *values {
		items = append(items, storage.MediaSeason{
			Name:         value.Name,
			EpisodeCount: value.EpisodeCount,
			AirDate:      value.AirDate,
			PosterPath:   value.PosterPath,
			Monitored:    optionalBoolValue(value.Monitored),
			Episodes:     storageMediaEpisodes(value.Episodes),
		})
	}
	return &items
}

func storageMediaEpisodes(values *[]MediaMetadataEpisode) []storage.MediaEpisode {
	if values == nil {
		return []storage.MediaEpisode{}
	}
	items := make([]storage.MediaEpisode, 0, len(*values))
	for _, value := range *values {
		items = append(items, storage.MediaEpisode{
			Name:          value.Name,
			EpisodeNumber: value.EpisodeNumber,
			Overview:      value.Overview,
			AirDate:       value.AirDate,
			StillPath:     value.StillPath,
			Monitored:     optionalBoolValue(value.Monitored),
		})
	}
	return items
}

func optionalBoolValue(value *bool) bool {
	return value != nil && *value
}

func optionalMediaMonitorMode(value *MediaMonitorMode) *string {
	if value == nil {
		return nil
	}
	mode := string(*value)
	return &mode
}

func optionalMinimumAvailability(value *MinimumAvailability) *string {
	if value == nil {
		return nil
	}
	availability := string(*value)
	return &availability
}
