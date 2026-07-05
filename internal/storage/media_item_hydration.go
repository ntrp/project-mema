package storage

import (
	"context"

	storagegen "media-manager/internal/storage/generated"
)

func hydrateMediaItem(ctx context.Context, q storagegen.DBTX, item MediaItem) (MediaItem, error) {
	item, err := hydrateMediaItemSeries(ctx, q, item)
	if err != nil {
		return item, err
	}
	return hydrateMediaItemProfile(ctx, q, item)
}

func hydrateMediaItems(ctx context.Context, q storagegen.DBTX, items []MediaItem) ([]MediaItem, error) {
	for index := range items {
		item, err := hydrateMediaItem(ctx, q, items[index])
		if err != nil {
			return nil, err
		}
		items[index] = item
	}
	return items, nil
}

func hydrateMediaItemProfile(
	ctx context.Context,
	q storagegen.DBTX,
	item MediaItem,
) (MediaItem, error) {
	if item.QualityProfileID == nil {
		return item, nil
	}
	languages, err := loadMediaProfileSubtitleLanguages(ctx, q, *item.QualityProfileID)
	if err != nil {
		return item, err
	}
	item.SubtitleLanguages = languages
	return item, nil
}
