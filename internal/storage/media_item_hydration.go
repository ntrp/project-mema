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
	item, err = hydrateMediaItemSidecars(ctx, q, item)
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
	item, err = hydrateMediaItemAssemblyRuns(ctx, q, item)
	if err != nil {
		return item, err
	}
	return hydrateMediaItemComponentProvenance(ctx, q, item)
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
	profile, err := storagegen.New(q).GetMediaProfile(ctx, *item.QualityProfileID)
	if err != nil {
		return item, err
	}
	targets, err := loadMediaProfileSubtitleTargets(ctx, q, *item.QualityProfileID)
	if err != nil {
		return item, err
	}
	item.SubtitleTargets = targets
	item.SubtitleMode = profile.SubtitleMode
	item.AllowSubtitleReleaseFallback = profile.AllowSubtitleReleaseFallback
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

func hydrateMediaItemSidecars(
	ctx context.Context,
	q storagegen.DBTX,
	item MediaItem,
) (MediaItem, error) {
	sidecars, err := listMediaItemSidecars(ctx, q, item.ID)
	if err != nil {
		return item, err
	}
	item.Sidecars = sidecars
	item.MetadataFilePaths = mergedStringSet(item.MetadataFilePaths, metadataSidecarFilePaths(sidecars))
	return item, nil
}

func metadataSidecarFilePaths(sidecars []MediaItemSidecar) []string {
	paths := make([]string, 0, len(sidecars))
	for _, sidecar := range sidecars {
		if sidecar.SidecarType == MediaSidecarMetadata {
			paths = append(paths, sidecar.FilePath)
		}
	}
	return paths
}

func mergedStringSet(left []string, right []string) []string {
	seen := map[string]struct{}{}
	merged := make([]string, 0, len(left)+len(right))
	for _, value := range append(left, right...) {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		merged = append(merged, value)
	}
	return merged
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

func hydrateMediaItemComponentProvenance(
	ctx context.Context,
	q storagegen.DBTX,
	item MediaItem,
) (MediaItem, error) {
	rows, err := storagegen.New(q).ListMediaComponentProvenance(ctx, item.ID)
	if err != nil {
		return item, err
	}
	item.ComponentProvenance = make([]MediaComponentProvenance, 0, len(rows))
	for _, row := range rows {
		item.ComponentProvenance = append(item.ComponentProvenance, mediaComponentProvenanceFromRow(row))
	}
	return item, nil
}
