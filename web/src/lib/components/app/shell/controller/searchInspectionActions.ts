import {
	getIndexerSearch as getIndexerSearchRequest,
	getMetadataCache as getMetadataCacheRequest,
	type CacheInspectionLimits
} from '$lib/settings/api';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

const INSPECTION_PAGE_SIZE = 10;

export function createSearchInspectionActions(state: AppShellState, clearNotice: () => void) {
	async function loadMetadataCache(limits: CacheInspectionLimits = {}) {
		state.loadingMetadataCache = true;
		clearNotice();
		try {
			state.metadataCache = await getMetadataCacheRequest(limits);
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load metadata cache');
		} finally {
			state.loadingMetadataCache = false;
		}
	}

	async function loadIndexerSearch(limits: CacheInspectionLimits = {}) {
		state.loadingIndexerSearch = true;
		clearNotice();
		try {
			state.indexerSearch = await getIndexerSearchRequest(limits);
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load indexer search cache');
		} finally {
			state.loadingIndexerSearch = false;
		}
	}

	async function refreshMetadataCache() {
		await loadMetadataCache();
	}

	async function refreshIndexerSearch() {
		await loadIndexerSearch();
	}

	async function loadMoreIndexerSearchCache() {
		await loadIndexerSearch({
			cacheLimit: nextLimit(state.indexerSearch.cacheEntries.length),
			historyLimit: currentLimit(state.indexerSearch.historyEntries.length)
		});
	}

	async function loadMoreIndexerSearchHistory() {
		await loadIndexerSearch({
			cacheLimit: currentLimit(state.indexerSearch.cacheEntries.length),
			historyLimit: nextLimit(state.indexerSearch.historyEntries.length)
		});
	}

	async function loadMoreMetadataCache() {
		await loadMetadataCache({
			cacheLimit: nextLimit(state.metadataCache.entries.length),
			historyLimit: currentLimit(state.metadataCache.historyEntries.length)
		});
	}

	async function loadMoreMetadataSearchHistory() {
		await loadMetadataCache({
			cacheLimit: currentLimit(state.metadataCache.entries.length),
			historyLimit: nextLimit(state.metadataCache.historyEntries.length)
		});
	}

	return {
		refreshIndexerSearch,
		loadMoreIndexerSearchCache,
		loadMoreIndexerSearchHistory,
		refreshMetadataCache,
		loadMoreMetadataCache,
		loadMoreMetadataSearchHistory
	};
}

function currentLimit(count: number) {
	return count > 0 ? count : INSPECTION_PAGE_SIZE;
}

function nextLimit(count: number) {
	return currentLimit(count) + INSPECTION_PAGE_SIZE;
}
