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

export function wait(ms: number) {
	return new Promise((resolve) => window.setTimeout(resolve, ms));
}
