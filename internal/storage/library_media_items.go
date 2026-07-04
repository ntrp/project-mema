package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mediaItemQuerier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func createMediaItemIfMissing(ctx context.Context, q mediaItemQuerier, input MediaItemInput) (MediaItem, error) {
	input = normalizeMediaItemOptions(input)
	metadataPayloads, err := marshalMediaMetadata(input.MediaMetadataSnapshot)
	if err != nil {
		return MediaItem{}, err
	}
	var existingID uuid.UUID
	err = q.QueryRow(ctx, `
		select id
		from app.media_items
		where lower(media_type) = lower($1) and lower(title) = lower($2)
			and (($3::integer is null and year is null) or year = $3)
		order by created_at asc
		limit 1
	`, input.Type, input.Title, input.Year).Scan(&existingID)
	if err == nil {
		return updateExistingMediaItem(ctx, q, existingID, input)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return MediaItem{}, err
	}
	return insertMediaItem(ctx, q, input, metadataPayloads)
}

func updateExistingMediaItem(ctx context.Context, q mediaItemQuerier, id uuid.UUID, input MediaItemInput) (MediaItem, error) {
	mediaFolderPath, err := ensureMediaMainFolder(ctx, q, input)
	if err != nil {
		return MediaItem{}, err
	}
	if _, err := q.Exec(ctx, `
		update app.media_items
		set quality_profile_id = coalesce(quality_profile_id, $2::text),
			library_folder_id = coalesce(library_folder_id, $3::uuid),
			media_folder_path = coalesce(media_folder_path, $4::text),
			monitor_mode = $5,
			minimum_availability = $6,
			monitored = $7,
			series_type = coalesce($8::text, series_type),
			updated_at = case
				when (quality_profile_id is null and $2::text is not null)
					or (library_folder_id is null and $3::uuid is not null)
					or (media_folder_path is null and $4::text is not null)
					or monitor_mode <> $5
					or minimum_availability <> $6
					or monitored <> $7
					or ($8::text is not null and series_type is distinct from $8::text)
				then now()
				else updated_at
			end
		where id = $1
	`, id, input.QualityProfileID, input.LibraryFolderID, mediaFolderPath, input.MonitorMode,
		input.MinimumAvailability, input.Monitored, input.SeriesType); err != nil {
		return MediaItem{}, err
	}
	if len(input.Tags) > 0 {
		if err := assignMediaItemTags(ctx, q, id, input.Tags); err != nil {
			return MediaItem{}, err
		}
	}
	return getMediaItem(ctx, q, id)
}

func insertMediaItem(ctx context.Context, q mediaItemQuerier, input MediaItemInput, metadataPayloads mediaMetadataPayloads) (MediaItem, error) {
	id := uuid.New()
	var itemID uuid.UUID
	mediaFolderPath, err := ensureMediaMainFolder(ctx, q, input)
	if err != nil {
		return MediaItem{}, err
	}
	if err := q.QueryRow(ctx, `
		insert into app.media_items (
			id, media_type, title, year, monitored, external_provider, external_id, overview, poster_path,
			collection_id, collection_name, backdrop_path, metadata_status, original_language,
			series_type, release_date, first_air_date, runtime_minutes, season_count, episode_count, vote_average,
			genres, keywords, facts, seasons, cast_members, crew_members, recommendations, similar_media,
			monitor_mode, minimum_availability, quality_profile_id, library_folder_id, media_folder_path
		)
		values (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21,
			$22::jsonb, $23::jsonb, $24::jsonb, $25::jsonb, $26::jsonb, $27::jsonb, $28::jsonb, $29::jsonb,
			$30, $31, $32, $33, $34
		)
		returning id
	`, id, input.Type, input.Title, input.Year, input.Monitored, input.ExternalProvider, input.ExternalID,
		input.Overview, input.PosterPath, input.CollectionID, input.CollectionName, input.BackdropPath,
		input.MetadataStatus, input.OriginalLanguage, input.SeriesType, input.ReleaseDate, input.FirstAirDate,
		input.RuntimeMinutes, input.SeasonCount, input.EpisodeCount, input.VoteAverage,
		metadataPayloads.genres, metadataPayloads.keywords, metadataPayloads.facts, metadataPayloads.seasons,
		metadataPayloads.cast, metadataPayloads.crew, metadataPayloads.recommendations, metadataPayloads.similar,
		input.MonitorMode, input.MinimumAvailability, input.QualityProfileID, input.LibraryFolderID,
		mediaFolderPath).Scan(&itemID); err != nil {
		return MediaItem{}, err
	}
	if err := assignMediaItemTags(ctx, q, itemID, input.Tags); err != nil {
		return MediaItem{}, err
	}
	return getMediaItem(ctx, q, itemID)
}

func mediaKindToMediaType(kind string) (string, bool) {
	switch kind {
	case "movie", "anime_movie":
		return "movie", true
	case "series", "anime_series":
		return "serie", true
	default:
		return "", false
	}
}
