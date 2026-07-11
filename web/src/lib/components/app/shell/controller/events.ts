import type {
	DownloadActivityStatus,
	IndexerSearchCacheEntry,
	IndexerSearchCacheStats,
	IndexerSearchHistoryEntry,
	MetadataCacheEntry,
	MetadataCacheStats,
	MetadataSearchHistoryEntry
} from '$lib/settings/types';
import type { AppShellState } from './state.svelte';
import type { ServerEventEnvelope } from './types';

type IndexerCacheUpdate = {
	entry: IndexerSearchCacheEntry;
	stats: IndexerSearchCacheStats;
};

type MetadataCacheUpdate = {
	entry: MetadataCacheEntry;
	stats: MetadataCacheStats;
};

export function createEventActions(state: AppShellState) {
	function appendIndexerSearchHistory(entry: IndexerSearchHistoryEntry) {
		state.indexerSearch = {
			...state.indexerSearch,
			historyTotalEntries: state.indexerSearch.historyTotalEntries + 1,
			historyStats: incrementHistoryStats(state.indexerSearch.historyStats, entry),
			historyEntries: [entry, ...state.indexerSearch.historyEntries].slice(
				0,
				displayLimit(state.indexerSearch.historyEntries.length)
			)
		};
	}

	function upsertIndexerSearchCache(update: IndexerCacheUpdate) {
		state.indexerSearch = {
			...state.indexerSearch,
			stats: update.stats,
			cacheEntries: upsertByKey(
				state.indexerSearch.cacheEntries,
				update.entry,
				indexerCacheKey
			).slice(0, displayLimit(state.indexerSearch.cacheEntries.length))
		};
	}

	function upsertMetadataCache(update: MetadataCacheUpdate) {
		state.metadataCache = {
			...state.metadataCache,
			stats: update.stats,
			entries: upsertByKey(state.metadataCache.entries, update.entry, metadataCacheKey).slice(
				0,
				displayLimit(state.metadataCache.entries.length)
			)
		};
	}

	function appendMetadataSearchHistory(entry: MetadataSearchHistoryEntry) {
		state.metadataCache = {
			...state.metadataCache,
			historyTotalEntries: state.metadataCache.historyTotalEntries + 1,
			historyStats: incrementHistoryStats(state.metadataCache.historyStats, entry),
			historyEntries: [entry, ...state.metadataCache.historyEntries].slice(
				0,
				displayLimit(state.metadataCache.historyEntries.length)
			)
		};
	}

	function parseEventData<T>(event: Event) {
		const message = event as MessageEvent<string>;
		try {
			return (JSON.parse(message.data) as ServerEventEnvelope<T>).data;
		} catch {
			return undefined;
		}
	}

	return {
		appendIndexerSearchHistory,
		upsertIndexerSearchCache,
		upsertMetadataCache,
		appendMetadataSearchHistory,
		parseEventData
	};
}

export function mediaStatusFromActivity(status: DownloadActivityStatus) {
	if (status === 'completed') return 'downloaded';
	if (status === 'queued' || status === 'grabbed' || status === 'downloading') return 'downloading';
	return undefined;
}

function upsertByKey<T>(items: T[], next: T, keyFor: (_item: T) => string) {
	const key = keyFor(next);
	return [next, ...items.filter((item) => keyFor(item) !== key)];
}

function displayLimit(length: number) {
	return Math.max(length, 10);
}

function incrementHistoryStats<T extends { cacheHit: boolean; success: boolean }>(
	stats: { totalEntries: number; cacheHits: number; cacheMisses: number; failures: number },
	entry: T
) {
	return {
		totalEntries: stats.totalEntries + 1,
		cacheHits: stats.cacheHits + (entry.cacheHit ? 1 : 0),
		cacheMisses: stats.cacheMisses + (entry.cacheHit ? 0 : 1),
		failures: stats.failures + (entry.success ? 0 : 1)
	};
}

function indexerCacheKey(entry: IndexerSearchCacheEntry) {
	return `${entry.indexerId}:${entry.mediaType}:${entry.query}`;
}

function metadataCacheKey(entry: MetadataCacheEntry) {
	return `${entry.providerId}:${entry.mediaType}:${entry.query}:${entry.year}`;
}
