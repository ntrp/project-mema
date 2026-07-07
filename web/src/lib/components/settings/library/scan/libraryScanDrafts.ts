import {
	defaultMonitorModeForMatch,
	isSeriesKind,
	scanMediaKind,
	type MatchDraft
} from '$lib/components/settings/library/scan/libraryScanImport';
import type {
	LibraryMediaKind,
	LibraryScanItem,
	MediaMonitorMode,
	MediaSearchResult,
	MetadataProvider,
	MinimumAvailability,
	SeriesType
} from '$lib/settings/types';

export function defaultMetadataProviderId(
	providers: MetadataProvider[],
	_kind: LibraryMediaKind
): string {
	const enabled = providers.filter((provider) => provider.enabled);
	const tmdb = enabled.find(
		(provider) => provider.type === 'tmdb' || provider.name.toLowerCase() === 'tmdb'
	);
	if (tmdb) return tmdb.id;
	return enabled[0]?.id ?? providers[0]?.id ?? '';
}

export function matchFromScanItem(item: LibraryScanItem): MediaSearchResult | undefined {
	if (!item.matchedTitle) return undefined;
	return {
		id: item.mediaItemId,
		title: item.matchedTitle,
		type: isSeriesKind(scanMediaKind(item)) ? 'serie' : 'movie',
		year: item.matchedYear,
		externalProvider: item.matchedExternalProvider,
		externalId: item.matchedExternalId
	};
}

export function initialMatchDraft(
	item: LibraryScanItem,
	metadataProviders: MetadataProvider[],
	bulk: {
		qualityProfileId: string;
		monitorMode: MediaMonitorMode;
		minimumAvailability: MinimumAvailability;
		seriesType: SeriesType;
	}
): MatchDraft {
	const mediaKind = scanMediaKind(item);
	const matched = matchFromScanItem(item);
	return {
		selected: Boolean(matched) && item.status === 'pending' && !item.imported,
		query: item.matchedTitle ?? item.detectedTitle,
		mediaKind,
		metadataProviderId:
			item.selectedMetadataProviderId ?? defaultMetadataProviderId(metadataProviders, mediaKind),
		matched,
		results: [],
		searching: false,
		searched: Boolean(matched),
		qualityProfileId: bulk.qualityProfileId,
		monitorMode: matched ? defaultMonitorModeForMatch(matched) : bulk.monitorMode,
		minimumAvailability: bulk.minimumAvailability,
		seriesType: bulk.seriesType,
		removeDuplicate: false
	};
}

export function ensureScanDrafts(
	items: LibraryScanItem[],
	drafts: Record<string, MatchDraft>,
	metadataProviders: MetadataProvider[],
	bulk: {
		qualityProfileId: string;
		monitorMode: MediaMonitorMode;
		minimumAvailability: MinimumAvailability;
		seriesType: SeriesType;
	}
) {
	for (const item of items) {
		if (drafts[item.id]) continue;
		drafts[item.id] = initialMatchDraft(item, metadataProviders, bulk);
	}
}
