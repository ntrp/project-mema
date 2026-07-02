import type {
	DownloadActivity,
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
	function upsertActivity(activity: DownloadActivity) {
		state.activities = [activity, ...state.activities.filter((item) => item.id !== activity.id)];
	}

	function appendIndexerSearchHistory(entry: IndexerSearchHistoryEntry) {
		state.indexerSearch = {
			...state.indexerSearch,
			historyEntries: [entry, ...state.indexerSearch.historyEntries].slice(0, 100)
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
			).slice(0, 100)
		};
	}

	function upsertMetadataCache(update: MetadataCacheUpdate) {
		state.metadataCache = {
			...state.metadataCache,
			stats: update.stats,
			entries: upsertByKey(state.metadataCache.entries, update.entry, metadataCacheKey).slice(0, 100)
		};
	}

	function appendMetadataSearchHistory(entry: MetadataSearchHistoryEntry) {
		state.metadataCache = {
			...state.metadataCache,
			historyEntries: [entry, ...state.metadataCache.historyEntries].slice(0, 100)
		};
	}

	function updateMediaStatusFromActivity(activity: DownloadActivity) {
		const status = mediaStatusFromActivity(activity.status);
		if (!status) {
			return;
		}
		state.mediaItems = state.mediaItems.map((item) =>
			item.id === activity.mediaItemId ? { ...item, status } : item
		);
	}

	function mediaStatusFromActivity(status: DownloadActivityStatus) {
		if (status === 'completed') {
			return 'downloaded';
		}
		if (status === 'queued' || status === 'grabbed' || status === 'downloading') {
			return 'downloading';
		}
		return undefined;
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
		upsertActivity,
		updateMediaStatusFromActivity,
		appendIndexerSearchHistory,
		upsertIndexerSearchCache,
		upsertMetadataCache,
		appendMetadataSearchHistory,
		parseEventData
	};
}

function upsertByKey<T>(items: T[], next: T, keyFor: (_item: T) => string) {
	const key = keyFor(next);
	return [next, ...items.filter((item) => keyFor(item) !== key)];
}

function indexerCacheKey(entry: IndexerSearchCacheEntry) {
	return `${entry.indexerName}:${entry.indexerType}:${entry.mediaType}:${entry.query}`;
}

function metadataCacheKey(entry: MetadataCacheEntry) {
	return `${entry.providerName}:${entry.providerType}:${entry.mediaType}:${entry.query}:${entry.year}`;
}
