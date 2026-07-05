import type {
	LibraryMediaKind,
	LibraryScanImportRequest,
	LibraryScanItem,
	LibraryScanItemMatchRequest,
	MediaMonitorMode,
	MediaSearchResult,
	MinimumAvailability,
	SeriesType
} from '$lib/settings/types';
import type { QualityProfileOption } from '$lib/settings/types';
import {
	duplicateDraftStatesForRows,
	duplicateSelectionValid
} from './libraryScanDuplicates';

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
		const leftKey = `${folderName(left.path)}/${left.fileName}`;
		const rightKey = `${folderName(right.path)}/${right.fileName}`;
		return leftKey.localeCompare(rightKey);
	});
}

export function scanMediaKind(item: LibraryScanItem): LibraryMediaKind {
	return item.detectedMediaKind === 'unknown' ? 'movie' : item.detectedMediaKind;
}

export function isSeriesKind(kind: LibraryMediaKind) {
	return kind === 'series' || kind === 'anime_series';
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
		.replace(/[._]+/g, ' ')
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
	if (items.length === 0) return false;
	if (!duplicateSelectionValid(items, drafts)) return false;
	return items.every((item) => {
		const draft = drafts[item.id];
		return Boolean(
			draft?.matched && effectiveDraftValue(draft.qualityProfileId, bulkQualityProfileId)
		);
	});
}

export function importRequestForDraft(
	draft: MatchDraft,
	match: MediaSearchResult,
	bulk: {
		qualityProfileId: string;
		monitorMode: MediaMonitorMode;
		minimumAvailability: MinimumAvailability;
		seriesType: SeriesType;
	}
): LibraryScanItemMatchRequest {
	const series = isSeriesKind(draft.mediaKind);
	const monitorMode = effectiveDraftValue(draft.monitorMode, bulk.monitorMode);
	return {
		mediaKind: draft.mediaKind,
		title: match.title,
		year: match.year,
		monitored: monitorMode !== 'none',
		qualityProfileId: effectiveDraftValue(draft.qualityProfileId, bulk.qualityProfileId),
		monitorMode,
		minimumAvailability: series
			? 'released'
			: effectiveDraftValue(draft.minimumAvailability, bulk.minimumAvailability),
		seriesType: series ? effectiveDraftValue(draft.seriesType, bulk.seriesType) : undefined,
		metadataProviderId: draft.metadataProviderId || undefined,
		mediaItemId: match.id,
		externalProvider: match.externalProvider,
		externalId: match.externalId,
		overview: match.overview,
		posterPath: match.posterPath
	};
}

export function importPayloadForRows(
	importRows: LibraryScanItem[],
	allRows: LibraryScanItem[],
	drafts: Record<string, MatchDraft>,
	bulk: {
		qualityProfileId: string;
		monitorMode: MediaMonitorMode;
		minimumAvailability: MinimumAvailability;
		seriesType: SeriesType;
	}
): LibraryScanImportRequest {
	const duplicateStates = duplicateDraftStatesForRows(allRows, drafts);
	return {
		items: importRows.map((item) => {
			const draft = drafts[item.id];
			return {
				itemId: item.id,
				match: importRequestForDraft(draft, draft.matched!, bulk)
			};
		}),
		removeDuplicatePaths: Object.entries(drafts)
			.filter(([id, draft]) => draft.removeDuplicate && (!draft.matched || duplicateStates[id]?.duplicate))
			.map(([id]) => allRows.find((row) => row.id === id)?.path)
			.filter((path): path is string => Boolean(path))
	};
}

export function wait(ms: number) {
	return new Promise((resolve) => window.setTimeout(resolve, ms));
}
