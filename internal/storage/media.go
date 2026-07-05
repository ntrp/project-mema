package storage

import (
	"context"
	"errors"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *SettingsStore) ListMediaItems(ctx context.Context) ([]MediaItem, error) {
	rows, err := storagegen.New(s.pool).ListMediaItems(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]MediaItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mediaItemFromListRow(row))
	}
	return hydrateMediaItems(ctx, s.pool, items)
}

func (s *SettingsStore) SearchMediaItems(ctx context.Context, query string, mediaType *string, limit int) ([]MediaItem, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	rows, err := storagegen.New(s.pool).SearchMediaItems(ctx, storagegen.SearchMediaItemsParams{
		Query:     query,
		MediaType: textValue(mediaType),
		RowLimit:  int32(limit),
	})
	if err != nil {
		return nil, err
	}

	items := make([]MediaItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mediaItemFromSearchRow(row))
	}
	return hydrateMediaItems(ctx, s.pool, items)
}

func (s *SettingsStore) GetMediaItem(ctx context.Context, id uuid.UUID) (MediaItem, error) {
	return getMediaItem(ctx, s.pool, id)
}

func getMediaItem(ctx context.Context, q mediaItemQuerier, id uuid.UUID) (MediaItem, error) {
	row, err := storagegen.New(q).GetMediaItem(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaItem{}, ErrNotFound
	}
	if err != nil {
		return MediaItem{}, err
	}
	item := mediaItemFromGetRow(row)
	return hydrateMediaItem(ctx, q, item)
}

func (s *SettingsStore) CreateMediaItem(ctx context.Context, input MediaItemInput) (MediaItem, error) {
	input = normalizeMediaItemOptions(input)
	metadataPayloads, err := marshalMediaMetadata(input.MediaMetadataSnapshot)
	if err != nil {
		return MediaItem{}, err
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return MediaItem{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	id := uuid.New()
	mediaFolderPath, err := ensureMediaMainFolder(ctx, tx, input)
	if err != nil {
		return MediaItem{}, err
	}
	itemID, err := storagegen.New(tx).CreateMediaItemRecord(ctx, mediaItemRecordParams(id, input, metadataPayloads, mediaFolderPath))
	if err != nil {
		return MediaItem{}, err
	}
	if err := assignMediaItemTags(ctx, tx, itemID, input.Tags); err != nil {
		return MediaItem{}, err
	}
	if err := materializeMediaSeriesSnapshot(ctx, tx, itemID, input); err != nil {
		return MediaItem{}, err
	}
	item, err := getMediaItem(ctx, tx, itemID)
	if err != nil {
		return MediaItem{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return MediaItem{}, err
	}
	return item, nil
}

func (s *SettingsStore) ListMissingMediaItems(ctx context.Context) ([]MediaItem, error) {
	rows, err := storagegen.New(s.pool).ListMissingMediaItems(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]MediaItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mediaItemFromMissingRow(row))
	}
	return hydrateMediaItems(ctx, s.pool, items)
}
