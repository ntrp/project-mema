import {
	cleanMatchSearchTitle,
	defaultMonitorModeForMatch,
	scanItemGroupRows,
	type MatchDraft
} from './libraryScanImport';
import { normalizeDuplicateDrafts } from './libraryScanDuplicates';
import type { LibraryScanItem, MediaSearchResult } from '$lib/settings/types';

export function normalizedMatchTitle(value: string) {
	return cleanMatchSearchTitle(value)
		.normalize('NFD')
		.replace(/\p{Diacritic}/gu, '')
		.toLowerCase()
		.replace(/[^a-z0-9]+/g, '');
}

export function scanItemExternalReference(item: LibraryScanItem) {
	const source = `${item.fileName ?? ''} ${item.path ?? ''}`;
	const tagged = source.match(/\b(imdb(?:[ ._-]*id)?|tmdb|tvdb)[ ._:-]*((?:tt)?\d{3,12})\b/i);
	if (tagged) {
		return {
			provider: tagged[1].toLowerCase().startsWith('imdb') ? 'imdb' : tagged[1].toLowerCase(),
			id: tagged[2].toLowerCase()
		};
	}
	const imdb = source.match(/\b(tt\d{5,10})\b/i);
	return imdb ? { provider: 'imdb', id: imdb[1].toLowerCase() } : undefined;
}

export function autoMatchResult(item: LibraryScanItem, results: MediaSearchResult[]) {
	const reference = scanItemExternalReference(item);
	if (reference) {
		return results.find(
			(result) =>
				result.externalProvider?.toLowerCase() === reference.provider &&
				result.externalId?.toLowerCase() === reference.id
		);
	}
	if (results.length === 1) return results[0];
	const detectedTitle = normalizedMatchTitle(item.detectedTitle ?? '');
	if (!detectedTitle) return undefined;
	return results.find((result) => {
		if (normalizedMatchTitle(result.title) !== detectedTitle) return false;
		return (
			item.detectedYear === undefined ||
			result.year === undefined ||
			item.detectedYear === result.year
		);
	});
}

export function applyAutoMatch(
	item: LibraryScanItem,
	result: MediaSearchResult,
	rows: LibraryScanItem[],
	drafts: Record<string, MatchDraft>
) {
	const source = drafts[item.id];
	if (!source) return;
	for (const row of scanItemGroupRows(item, rows)) {
		const draft = drafts[row.id];
		if (!draft || row.imported || (row.status && row.status !== 'pending')) continue;
		draft.matched = result;
		draft.query = result.title;
		draft.results = [];
		draft.selected = true;
		draft.searched = true;
		draft.searching = false;
		draft.metadataProviderId = source.metadataProviderId;
		draft.monitorMode = defaultMonitorModeForMatch(result);
		normalizeDuplicateDrafts(rows, drafts, row.duplicateGroupId);
	}
}
