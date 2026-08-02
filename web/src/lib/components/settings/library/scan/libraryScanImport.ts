import type {
	LibraryMediaKind,
	LibraryScanItem,
	LibraryScanItemMatchRequest,
	MediaMonitorMode,
	MediaSearchResult,
	MinimumAvailability,
	SeriesType
} from '$lib/settings/types';
import type { QualityProfileOption } from '$lib/settings/types';
import { duplicateSelectionValid } from './libraryScanDuplicates';

export interface MatchDraft {
	selected: boolean;
	query: string;
	mediaKind: LibraryMediaKind;
	metadataProviderId: string;
	matched?: MediaSearchResult;
	results: MediaSearchResult[];
	searching: boolean;
	searched: boolean;
	qualityProfileId: string;
	monitorMode: MediaMonitorMode;
	minimumAvailability: MinimumAvailability;
	seriesType: SeriesType;
	seasonNumber?: number;
	episodeNumber?: number;
	removeDuplicate: boolean;
}

export interface LibraryScanImportRow {
	item: LibraryScanItem;
	request: LibraryScanItemMatchRequest;
}

export const movieMonitorModeOptions: { value: MediaMonitorMode; label: string }[] = [
	{ value: 'only_media', label: 'Only this media' },
	{ value: 'collection', label: 'Entire collection' },
	{ value: 'none', label: 'None' }
];

export const seriesMonitorModeOptions: { value: MediaMonitorMode; label: string }[] = [
	{ value: 'all_episodes', label: 'All episodes' },
	{ value: 'future_episodes', label: 'Future episodes' },
	{ value: 'missing_episodes', label: 'Missing episodes' },
	{ value: 'existing_episodes', label: 'Existing episodes' },
	{ value: 'no_specials', label: 'No specials' },
	{ value: 'none', label: 'None' }
];

export const minimumAvailabilityOptions: { value: MinimumAvailability; label: string }[] = [
	{ value: 'released', label: 'Released' },
	{ value: 'in_cinema', label: 'In cinema' },
	{ value: 'announced', label: 'Announced' }
];

export const seriesTypeOptions: { value: SeriesType; label: string }[] = [
	{ value: 'standard', label: 'Standard' },
	{ value: 'daily', label: 'Daily / Date' },
	{ value: 'absolute', label: 'Absolute' }
];

export function folderName(path: string) {
	const normalized = path.replaceAll('\\', '/');
	const parts = normalized.split('/').filter(Boolean);
	return parts.length > 1 ? parts.slice(0, -1).join('/') : '.';
}

export function sortedScanItems(items: LibraryScanItem[]) {
	return [...items].sort((left, right) => {
		const leftKey = `${left.seriesFolderPath ?? folderName(left.path)}/${left.fileName}`;
		const rightKey = `${right.seriesFolderPath ?? folderName(right.path)}/${right.fileName}`;
		return leftKey.localeCompare(rightKey);
	});
}

export function scanMediaKind(item: LibraryScanItem): LibraryMediaKind {
	return item.detectedMediaKind === 'unknown' ? 'movie' : item.detectedMediaKind;
}

export function isSeriesKind(kind: LibraryMediaKind) {
	return kind === 'series' || kind === 'anime_series';
}

export function isDraftSelectableForImport(item: LibraryScanItem, draft?: MatchDraft) {
	return Boolean(
		draft?.matched &&
		item.status === 'pending' &&
		!item.imported &&
		!draft.removeDuplicate &&
		(!isSeriesKind(draft.mediaKind) ||
			(draft.seasonNumber !== undefined && draft.episodeNumber !== undefined))
	);
}

export function scanItemGroupKey(item: LibraryScanItem) {
	if (isSeriesKind(scanMediaKind(item)) && item.seriesFolderPath) {
		return `series:${item.seriesFolderPath.replaceAll('\\', '/')}`;
	}
	return `item:${item.id}`;
}

export function scanItemGroupRows(item: LibraryScanItem, rows: LibraryScanItem[]) {
	const key = scanItemGroupKey(item);
	return rows.filter((row) => scanItemGroupKey(row) === key);
}

export function scanItemGroupRepresentative(item: LibraryScanItem, rows: LibraryScanItem[]) {
	const group = scanItemGroupRows(item, rows);
	return group.find((row) => !row.imported && row.status === 'pending') ?? group[0];
}

export function isScanItemGroupRepresentative(item: LibraryScanItem, rows: LibraryScanItem[]) {
	return scanItemGroupRepresentative(item, rows)?.id === item.id;
}

export function relativeScanPath(path: string, root: string) {
	const normalizedPath = path.replaceAll('\\', '/').replace(/\/+$/g, '');
	const normalizedRoot = root.replaceAll('\\', '/').replace(/\/+$/g, '');
	const comparablePath = /^[a-z]:\//i.test(normalizedPath)
		? normalizedPath.toLowerCase()
		: normalizedPath;
	const comparableRoot = /^[a-z]:\//i.test(normalizedRoot)
		? normalizedRoot.toLowerCase()
		: normalizedRoot;
	if (!normalizedRoot || comparablePath === comparableRoot) return normalizedPath;
	if (!comparablePath.startsWith(`${comparableRoot}/`)) return normalizedPath;
	return normalizedPath.slice(normalizedRoot.length + 1) || normalizedPath;
}

export function defaultMonitorModeForMatch(match?: MediaSearchResult): MediaMonitorMode {
	return match?.type === 'serie' ? 'all_episodes' : 'only_media';
}

export function matchedDraftKind(draft?: MatchDraft) {
	if (!draft?.matched) return undefined;
	return draft.matched.type === 'serie' ? 'series' : 'movie';
}

export function searchCacheKey(kind: LibraryMediaKind, providerId: string, query: string) {
	return `${kind}:${providerId}:${query.trim().toLowerCase()}`;
}

export function cleanMatchSearchTitle(value: string) {
	const normalized = value
		.replace(/\.[a-z0-9]{2,5}$/i, '')
		.replace(/\b(?:imdb(?:[ ._-]*id)?|tmdb|tvdb)[ ._:-]*(?:tt)?\d{3,12}\b/gi, ' ')
		.replace(/\btt\d{5,10}\b/gi, ' ')
		.replace(/[()[\]{}]/g, ' ')
		.replace(/[._-]+/g, ' ')
		.replace(/\s+/g, ' ')
		.trim();
	const [beforeMediaToken] = normalized.split(
		/\s+\b(?:19|20)\d{2}\b|\b(?:s\d{1,2}e\d{1,3}|season\s+\d+|episode\s+\d+|2160p|1080p|720p|576p|480p|web-?dl|webrip|bluray|brrip|hdtv|dvdrip|remux|x26[45]|h\.?26[45]|hevc|av1)\b/i
	);
	return (beforeMediaToken || normalized).trim();
}

export function effectiveDraftValue<T extends string>(draftValue: T, bulkValue: T): T {
	return draftValue || bulkValue;
}

export function defaultQualityProfileId(profiles: QualityProfileOption[]) {
	return profiles.find((profile) => profile.isDefault)?.id ?? '';
}

export function canImportRows(
	items: LibraryScanItem[],
	drafts: Record<string, MatchDraft>,
	bulkQualityProfileId: string
) {
	return importDisabledReason(items, drafts, bulkQualityProfileId) === '';
}

export function importDisabledReason(
	items: LibraryScanItem[],
	drafts: Record<string, MatchDraft>,
	bulkQualityProfileId: string
) {
	if (items.length === 0) return 'Select at least one file to import';
	if (!duplicateSelectionValid(items, drafts))
		return 'Resolve duplicate selections before importing';
	for (const item of items) {
		const draft = drafts[item.id];
		if (!draft?.matched) return 'Match every selected item before importing';
		if (!effectiveDraftValue(draft.qualityProfileId, bulkQualityProfileId)) {
			return 'Choose a quality profile for every selected item';
		}
		if (
			isSeriesKind(draft.mediaKind) &&
			(draft.seasonNumber === undefined || draft.episodeNumber === undefined)
		) {
			return 'Set season and episode for every selected episode';
		}
	}
	return '';
}

export function wait(ms: number) {
	return new Promise((resolve) => window.setTimeout(resolve, ms));
}
