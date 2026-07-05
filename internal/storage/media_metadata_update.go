package storage

import (
	"context"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

func (s *SettingsStore) UpdateMediaItemMetadata(ctx context.Context, id uuid.UUID, input MediaItemInput) (MediaItem, error) {
	metadataPayloads, err := marshalMediaMetadata(input.MediaMetadataSnapshot)
	if err != nil {
		return MediaItem{}, err
	}
	if err := storagegen.New(s.pool).UpdateMediaItemMetadataRecord(ctx, mediaItemMetadataParams(id, input, metadataPayloads)); err != nil {
		return MediaItem{}, err
	}
	return s.GetMediaItem(ctx, id)
}

func mediaItemMetadataParams(
	id uuid.UUID,
	input MediaItemInput,
	payloads mediaMetadataPayloads,
) storagegen.UpdateMediaItemMetadataRecordParams {
	return storagegen.UpdateMediaItemMetadataRecordParams{
		MediaType:        input.Type,
		Title:            input.Title,
		Year:             int4Value(input.Year),
		ExternalProvider: textValue(input.ExternalProvider),
		ExternalID:       textValue(input.ExternalID),
		Overview:         textValue(input.Overview),
		PosterPath:       textValue(input.PosterPath),
		CollectionID:     textValue(input.CollectionID),
		CollectionName:   textValue(input.CollectionName),
		BackdropPath:     textValue(input.BackdropPath),
		MetadataStatus:   textValue(input.MetadataStatus),
		OriginalLanguage: textValue(input.OriginalLanguage),
		ReleaseDate:      textValue(input.ReleaseDate),
		FirstAirDate:     textValue(input.FirstAirDate),
		RuntimeMinutes:   int4Value(input.RuntimeMinutes),
		SeasonCount:      int4Value(input.SeasonCount),
		EpisodeCount:     int4Value(input.EpisodeCount),
		VoteAverage:      float8Value(input.VoteAverage),
		Genres:           payloads.genres,
		Keywords:         payloads.keywords,
		Facts:            payloads.facts,
		Seasons:          payloads.seasons,
		CastMembers:      payloads.cast,
		CrewMembers:      payloads.crew,
		Recommendations:  payloads.recommendations,
		SimilarMedia:     payloads.similar,
		ID:               id,
	}
}
