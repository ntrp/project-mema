import type {
	LibraryMediaKind,
	LibraryScanItem,
	LibraryScanItemMatchRequest,
	MediaMonitorMode,
	MediaSearchResult,
	MinimumAvailability
} from '$lib/settings/types';

export interface MatchDraft {
	selected: boolean;
	query: string;
	mediaKind: LibraryMediaKind;
	matched?: MediaSearchResult;
	results: MediaSearchResult[];
	searching: boolean;
	searched: boolean;
	qualityProfileId: string;
	monitorMode: MediaMonitorMode;
	minimumAvailability: MinimumAvailability;
}

export interface LibraryScanImportRow {
	item: LibraryScanItem;
	request: LibraryScanItemMatchRequest;
}

export type LibraryScanSortMode = 'folders' | 'mixed';

export const monitorModeOptions: { value: MediaMonitorMode; label: string }[] = [
	{ value: 'only_media', label: 'Only this media' },
	{ value: 'collection', label: 'Entire collection' },
	{ value: 'none', label: 'None' }
];

export const minimumAvailabilityOptions: { value: MinimumAvailability; label: string }[] = [
	{ value: 'released', label: 'Released' },
	{ value: 'in_cinema', label: 'In cinema' },
	{ value: 'announced', label: 'Announced' }
];

export const scanSortModeOptions: { value: LibraryScanSortMode; label: string }[] = [
	{ value: 'folders', label: 'Folders' },
	{ value: 'mixed', label: 'Mixed' }
];

export function folderName(path: string) {
	const normalized = path.replaceAll('\\', '/');
	const parts = normalized.split('/').filter(Boolean);
	return parts.length > 1 ? parts.slice(0, -1).join('/') : '.';
}

export function sortedScanItems(items: LibraryScanItem[], sortMode: LibraryScanSortMode) {
	return [...items].sort((left, right) => {
		const leftKey =
			sortMode === 'folders' ? `${folderName(left.path)}/${left.fileName}` : left.fileName;
		const rightKey =
			sortMode === 'folders' ? `${folderName(right.path)}/${right.fileName}` : right.fileName;
		return leftKey.localeCompare(rightKey);
	});
}

export function scanMediaKind(item: LibraryScanItem): LibraryMediaKind {
	return item.detectedMediaKind === 'unknown' ? 'movie' : item.detectedMediaKind;
}

export function searchCacheKey(kind: LibraryMediaKind, query: string) {
	return `${kind}:${query.trim().toLowerCase()}`;
}

export function effectiveDraftValue<T extends string>(draftValue: T, bulkValue: T): T {
	return draftValue || bulkValue;
}

export function canImportRows(
	items: LibraryScanItem[],
	drafts: Record<string, MatchDraft>,
	bulkQualityProfileId: string
) {
	return (
		items.length > 0 &&
		items.every((item) => {
			const draft = drafts[item.id];
			return Boolean(
				draft?.matched && effectiveDraftValue(draft.qualityProfileId, bulkQualityProfileId)
			);
		})
	);
}

export function importRequestForDraft(
	draft: MatchDraft,
	match: MediaSearchResult,
	bulk: {
		qualityProfileId: string;
		monitorMode: MediaMonitorMode;
		minimumAvailability: MinimumAvailability;
	}
): LibraryScanItemMatchRequest {
	const monitorMode = effectiveDraftValue(draft.monitorMode, bulk.monitorMode);
	return {
		mediaKind: draft.mediaKind,
		title: match.title,
		year: match.year,
		monitored: monitorMode !== 'none',
		qualityProfileId: effectiveDraftValue(draft.qualityProfileId, bulk.qualityProfileId),
		monitorMode,
		minimumAvailability: effectiveDraftValue(draft.minimumAvailability, bulk.minimumAvailability),
		externalProvider: match.externalProvider,
		externalId: match.externalId,
		overview: match.overview,
		posterPath: match.posterPath
	};
}

export function wait(ms: number) {
	return new Promise((resolve) => window.setTimeout(resolve, ms));
}
