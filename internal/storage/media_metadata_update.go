package storage

import (
	"context"

	"github.com/google/uuid"
)

func (s *SettingsStore) UpdateMediaItemMetadata(ctx context.Context, id uuid.UUID, input MediaItemInput) (MediaItem, error) {
	metadataPayloads, err := marshalMediaMetadata(input.MediaMetadataSnapshot)
	if err != nil {
		return MediaItem{}, err
	}
	_, err = s.pool.Exec(ctx, `
		update app.media_items
		set
			media_type = $2,
			title = $3,
			year = $4,
			external_provider = $5,
			external_id = $6,
			overview = $7,
			poster_path = $8,
			collection_id = $9,
			collection_name = $10,
			backdrop_path = $11,
			metadata_status = $12,
			original_language = $13,
			release_date = $14,
			first_air_date = $15,
			runtime_minutes = $16,
			season_count = $17,
			episode_count = $18,
			vote_average = $19,
			genres = $20::jsonb,
			keywords = $21::jsonb,
			facts = $22::jsonb,
			seasons = $23::jsonb,
			cast_members = $24::jsonb,
			crew_members = $25::jsonb,
			recommendations = $26::jsonb,
			similar_media = $27::jsonb,
			updated_at = now()
		where id = $1
	`, id, input.Type, input.Title, input.Year, input.ExternalProvider, input.ExternalID,
		input.Overview, input.PosterPath, input.CollectionID, input.CollectionName, input.BackdropPath,
		input.MetadataStatus, input.OriginalLanguage, input.ReleaseDate, input.FirstAirDate,
		input.RuntimeMinutes, input.SeasonCount, input.EpisodeCount, input.VoteAverage,
		metadataPayloads.genres, metadataPayloads.keywords, metadataPayloads.facts,
		metadataPayloads.seasons, metadataPayloads.cast, metadataPayloads.crew, metadataPayloads.recommendations,
		metadataPayloads.similar)
	if err != nil {
		return MediaItem{}, err
	}
	return s.GetMediaItem(ctx, id)
}
