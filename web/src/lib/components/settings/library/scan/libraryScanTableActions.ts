import { applyAutoMatch, autoMatchResult } from './libraryScanAutoMatch';
import {
	cleanMatchSearchTitle,
	isScanItemGroupRepresentative,
	scanItemGroupKey,
	scanItemGroupRows,
	searchCacheKey,
	type MatchDraft
} from './libraryScanImport';
import type { LibraryScan, LibraryScanItem, MediaSearchResult } from '$lib/settings/types';

type SearchCache = Record<string, MediaSearchResult[] | undefined>;

export function prepareProviderSearch(item: LibraryScanItem, draft: MatchDraft) {
	draft.query = cleanMatchSearchTitle(item.detectedTitle || item.fileName || draft.query);
	draft.matched = undefined;
	draft.results = [];
	draft.selected = false;
	draft.searched = false;
}

export function scheduleScanItemSearch(input: {
	item: LibraryScanItem;
	drafts: Record<string, MatchDraft>;
	searchTimers: Record<string, ReturnType<typeof globalThis.setTimeout>>;
	search: (_item: LibraryScanItem, _auto: boolean) => Promise<void>;
}) {
	const { item, drafts, searchTimers, search } = input;
	const draft = drafts[item.id];
	if (!draft) return;
	globalThis.clearTimeout(searchTimers[item.id]);
	if (draft.query.trim().length < 2) {
		draft.results = [];
		draft.searched = false;
		draft.searching = false;
		return;
	}
	searchTimers[item.id] = globalThis.setTimeout(() => void search(item, false), 1000);
}

export function changeScanItemProvider(input: {
	item: LibraryScanItem;
	allRows: LibraryScanItem[];
	providerId: string;
	drafts: Record<string, MatchDraft>;
	search: (_item: LibraryScanItem, _auto: boolean) => Promise<void>;
}) {
	const { item, allRows, providerId, drafts, search } = input;
	const draft = drafts[item.id];
	if (!draft || draft.metadataProviderId === providerId) return;
	for (const row of scanItemGroupRows(item, allRows)) {
		const groupedDraft = drafts[row.id];
		if (!groupedDraft || row.imported || (row.status && row.status !== 'pending')) continue;
		groupedDraft.metadataProviderId = providerId;
		prepareProviderSearch(row, groupedDraft);
	}
	void search(item, true);
}

export function applyScanItemProvider(input: {
	rows: LibraryScanItem[];
	allRows?: LibraryScanItem[];
	drafts: Record<string, MatchDraft>;
	providerId: string;
	search: (_item: LibraryScanItem, _auto: boolean) => Promise<void>;
}) {
	const { rows, allRows = rows, drafts, providerId, search } = input;
	const searched = new Set<string>();
	for (const item of rows) {
		const key = scanItemGroupKey(item);
		if (searched.has(key)) continue;
		searched.add(key);
		for (const row of scanItemGroupRows(item, allRows)) {
			const draft = drafts[row.id];
			if (!draft || row.imported || (row.status && row.status !== 'pending')) continue;
			draft.metadataProviderId = providerId;
			prepareProviderSearch(row, draft);
		}
		void search(item, true);
	}
}

export function searchPendingScanItems(input: {
	rows: LibraryScanItem[];
	drafts: Record<string, MatchDraft>;
	autoSearchStarted: Record<string, boolean>;
	search: (_item: LibraryScanItem, _auto: boolean) => Promise<void>;
}) {
	const { rows, drafts, autoSearchStarted, search } = input;
	for (const item of rows) {
		const draft = drafts[item.id];
		if (!isScanItemGroupRepresentative(item, rows)) continue;
		if (!draft || draft.matched || item.imported || item.status !== 'pending') continue;
		if (autoSearchStarted[item.id] || draft.query.trim().length < 2) continue;
		autoSearchStarted[item.id] = true;
		void search(item, true);
	}
}

export async function searchScanItem(input: {
	item: LibraryScanItem;
	allRows: LibraryScanItem[];
	drafts: Record<string, MatchDraft>;
	searchCache: SearchCache;
	auto: boolean;
	onSearchMatch: (
		_kind: MatchDraft['mediaKind'],
		_query: string,
		_providerId?: string
	) => Promise<MediaSearchResult[]>;
}) {
	const { item, allRows, drafts, searchCache, auto, onSearchMatch } = input;
	const draft = drafts[item.id];
	if (!draft || draft.query.trim().length < 2) return;
	const query = draft.query;
	const key = searchCacheKey(draft.mediaKind, draft.metadataProviderId, query);
	let matchedCurrentSearch = false;
	draft.searching = true;
	try {
		const results =
			searchCache[key] ?? (await onSearchMatch(draft.mediaKind, query, draft.metadataProviderId));
		if (draft.query !== query) return;
		searchCache[key] = results;
		draft.results = results;
		if (auto) {
			const result = autoMatchResult(item, results);
			if (result) {
				applyAutoMatch(item, result, allRows, drafts);
				matchedCurrentSearch = true;
			}
		}
		draft.searched = true;
	} catch {
		if (draft.query !== query) return;
		draft.results = [];
		draft.searched = true;
	} finally {
		if (draft.query === query || matchedCurrentSearch) draft.searching = false;
	}
}

export async function resetScanItemImport(input: {
	item: LibraryScanItem;
	drafts: Record<string, MatchDraft>;
	resetting: { itemId: string };
	scan: LibraryScan;
	onResetImport: (_scan: LibraryScan, _itemId: string) => Promise<void>;
}) {
	const { item, drafts, resetting, scan, onResetImport } = input;
	if (!item.imported || resetting.itemId) return;
	const existing = { ...drafts[item.id] };
	resetting.itemId = item.id;
	try {
		await onResetImport(scan, item.id);
		drafts[item.id] = {
			...existing,
			selected: false,
			matched: undefined,
			results: [],
			searching: false,
			searched: false,
			removeDuplicate: false
		};
	} finally {
		resetting.itemId = '';
	}
}
