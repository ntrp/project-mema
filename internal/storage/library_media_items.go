package storage

import (
	"context"
	"errors"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type mediaItemQuerier = storagegen.DBTX

func createMediaItemIfMissing(ctx context.Context, q mediaItemQuerier, input MediaItemInput) (MediaItem, error) {
	input = normalizeMediaItemOptions(input)
	metadataPayloads, err := marshalMediaMetadata(input.MediaMetadataSnapshot)
	if err != nil {
		return MediaItem{}, err
	}
	existingID, err := storagegen.New(q).FindExistingMediaItemID(ctx, storagegen.FindExistingMediaItemIDParams{
		MediaType: input.Type,
		Title:     input.Title,
		Year:      int4Value(input.Year),
	})
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
	if err := storagegen.New(q).UpdateExistingMediaItem(ctx, storagegen.UpdateExistingMediaItemParams{
		QualityProfileID:    textValue(input.QualityProfileID),
		LibraryFolderID:     input.LibraryFolderID,
		MediaFolderPath:     textValue(mediaFolderPath),
		MonitorMode:         input.MonitorMode,
		MinimumAvailability: input.MinimumAvailability,
		Monitored:           input.Monitored,
		SeriesType:          textValue(input.SeriesType),
		ID:                  id,
	}); err != nil {
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
	mediaFolderPath, err := ensureMediaMainFolder(ctx, q, input)
	if err != nil {
		return MediaItem{}, err
	}
	itemID, err := storagegen.New(q).CreateMediaItemRecord(ctx, mediaItemRecordParams(id, input, metadataPayloads, mediaFolderPath))
	if err != nil {
		return MediaItem{}, err
	}
	if err := assignMediaItemTags(ctx, q, itemID, input.Tags); err != nil {
		return MediaItem{}, err
	}
	if err := materializeMediaSeriesSnapshot(ctx, q, itemID, input); err != nil {
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
