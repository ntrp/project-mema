import { defaultMonitorModeForMatch, type MatchDraft } from './libraryScanImport';
import { normalizeDuplicateDrafts } from './libraryScanDuplicates';
import type { LibraryScanItem, MediaSearchResult } from '$lib/settings/types';

export function normalizedMatchTitle(value: string) {
	return value
		.normalize('NFD')
		.replace(/\p{Diacritic}/gu, '')
		.toLowerCase()
		.replace(/[^a-z0-9]+/g, '');
}

export function scanItemExternalReference(item: LibraryScanItem) {
	const source = `${item.fileName ?? ''} ${item.path ?? ''}`;
	const match = source.match(/\b(tmdb|tvdb|imdb)[-_. ]?([a-z0-9]+)\b/i);
	if (!match) return undefined;
	return {
		provider: match[1].toLowerCase(),
		id: match[2].toLowerCase()
	};
}

export function autoMatchResult(item: LibraryScanItem, results: MediaSearchResult[]) {
	const reference = scanItemExternalReference(item);
	if (reference) {
		const referenced = results.find(
			(result) =>
				result.externalProvider?.toLowerCase() === reference.provider &&
				result.externalId?.toLowerCase() === reference.id
		);
		if (referenced) return referenced;
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
	const draft = drafts[item.id];
	if (!draft) return;
	draft.matched = result;
	draft.query = result.title;
	draft.results = [];
	draft.selected = true;
	draft.monitorMode = defaultMonitorModeForMatch(result);
	normalizeDuplicateDrafts(rows, drafts, item.duplicateGroupId);
}
