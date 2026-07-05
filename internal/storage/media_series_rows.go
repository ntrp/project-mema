package storage

import (
	"encoding/json"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

func mediaSeriesSeasonParams(
	mediaItemID uuid.UUID,
	input MediaSeriesSeasonInput,
) (storagegen.UpsertMediaSeasonRowParams, error) {
	source, err := mediaSeriesSourcePayload(input.Source)
	if err != nil {
		return storagegen.UpsertMediaSeasonRowParams{}, err
	}
	return storagegen.UpsertMediaSeasonRowParams{
		ID:               uuid.New(),
		MediaItemID:      mediaItemID,
		ExternalProvider: textValue(input.ExternalProvider),
		ExternalID:       textValue(input.ExternalID),
		SeasonNumber:     input.SeasonNumber,
		Name:             input.Name,
		Overview:         textValue(input.Overview),
		AirDate:          textValue(input.AirDate),
		PosterPath:       textValue(input.PosterPath),
		EpisodeCount:     int4Value(input.EpisodeCount),
		Monitored:        input.Monitored,
		Source:           source,
	}, nil
}

func mediaSeriesEpisodeParams(
	mediaItemID uuid.UUID,
	seasonID uuid.UUID,
	seasonNumber int32,
	input MediaSeriesEpisodeInput,
) (storagegen.UpsertMediaEpisodeRowParams, error) {
	source, err := mediaSeriesSourcePayload(input.Source)
	if err != nil {
		return storagegen.UpsertMediaEpisodeRowParams{}, err
	}
	return storagegen.UpsertMediaEpisodeRowParams{
		ID:               uuid.New(),
		SeasonID:         seasonID,
		MediaItemID:      mediaItemID,
		ExternalProvider: textValue(input.ExternalProvider),
		ExternalID:       textValue(input.ExternalID),
		SeasonNumber:     seasonNumber,
		EpisodeNumber:    input.EpisodeNumber,
		Name:             input.Name,
		Overview:         textValue(input.Overview),
		AirDate:          textValue(input.AirDate),
		StillPath:        textValue(input.StillPath),
		RuntimeMinutes:   int4Value(input.RuntimeMinutes),
		Monitored:        input.Monitored,
		Source:           source,
	}, nil
}

func mediaSeriesSeasonFromRow(row storagegen.AppMediaSeason) MediaSeriesSeason {
	return MediaSeriesSeason{
		ID:               row.ID,
		MediaItemID:      row.MediaItemID,
		ExternalProvider: textPtr(row.ExternalProvider),
		ExternalID:       textPtr(row.ExternalID),
		SeasonNumber:     row.SeasonNumber,
		Name:             row.Name,
		Overview:         textPtr(row.Overview),
		AirDate:          textPtr(row.AirDate),
		PosterPath:       textPtr(row.PosterPath),
		EpisodeCount:     int4Ptr(row.EpisodeCount),
		Monitored:        row.Monitored,
		Source:           mediaSeriesSourceFromPayload(row.Source),
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
}

func mediaSeriesEpisodeFromRow(row storagegen.AppMediaEpisode) MediaSeriesEpisode {
	return MediaSeriesEpisode{
		ID:               row.ID,
		SeasonID:         row.SeasonID,
		MediaItemID:      row.MediaItemID,
		ExternalProvider: textPtr(row.ExternalProvider),
		ExternalID:       textPtr(row.ExternalID),
		SeasonNumber:     row.SeasonNumber,
		EpisodeNumber:    row.EpisodeNumber,
		Name:             row.Name,
		Overview:         textPtr(row.Overview),
		AirDate:          textPtr(row.AirDate),
		StillPath:        textPtr(row.StillPath),
		RuntimeMinutes:   int4Ptr(row.RuntimeMinutes),
		Monitored:        row.Monitored,
		Source:           mediaSeriesSourceFromPayload(row.Source),
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
}

func mediaSeriesSourceFromPayload(payload []byte) map[string]any {
	source := map[string]any{}
	if len(payload) > 0 {
		_ = json.Unmarshal(payload, &source)
	}
	return source
}
