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
	item, err = hydrateMediaItemProfile(ctx, q, item)
	if err != nil {
		return item, err
	}
	item, err = hydrateMediaItemSubtitles(ctx, q, item)
	if err != nil {
		return item, err
	}
	item, err = hydrateMediaItemComponentSources(ctx, q, item)
	if err != nil {
		return item, err
	}
	item, err = hydrateMediaItemAnime(ctx, q, item)
	if err != nil {
		return item, err
	}
	return hydrateMediaItemAssemblyRuns(ctx, q, item)
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
	targets, err := loadMediaProfileSubtitleTargets(ctx, q, *item.QualityProfileID)
	if err != nil {
		return item, err
	}
	item.SubtitleTargets = targets
	return item, nil
}

func hydrateMediaItemSubtitles(
	ctx context.Context,
	q storagegen.DBTX,
	item MediaItem,
) (MediaItem, error) {
	subtitles, err := listMediaItemSubtitles(ctx, q, item.ID)
	if err != nil {
		return item, err
	}
	item.ExternalSubtitles = subtitles
	item.MetadataFilePaths = append(item.MetadataFilePaths, subtitleFilePaths(subtitles)...)
	return item, nil
}

func hydrateMediaItemComponentSources(
	ctx context.Context,
	q storagegen.DBTX,
	item MediaItem,
) (MediaItem, error) {
	sources, err := listMediaComponentSources(ctx, q, item.ID)
	if err != nil {
		return item, err
	}
	for index := range sources {
		artifacts, err := listMediaComponentArtifactsForSource(ctx, q, sources[index].ID)
		if err != nil {
			return item, err
		}
		sources[index].Artifacts = artifacts
		compatibility, err := listMediaComponentCompatibilityForSource(ctx, q, sources[index].ID)
		if err != nil {
			return item, err
		}
		sources[index].Compatibility = compatibility
	}
	item.ComponentSources = sources
	return item, nil
}

func hydrateMediaItemAssemblyRuns(
	ctx context.Context,
	q storagegen.DBTX,
	item MediaItem,
) (MediaItem, error) {
	runs, err := listMediaComponentAssemblyRuns(ctx, q, item.ID)
	if err != nil {
		return item, err
	}
	item.AssemblyRuns = runs
	return item, nil
}
