import {
	clearIndexerSearchCache as clearIndexerSearchCacheRequest,
	clearIndexerSearchCacheByPattern as clearIndexerSearchCacheByPatternRequest,
	clearIndexerSearchHistory as clearIndexerSearchHistoryRequest,
	clearMetadataCache as clearMetadataCacheRequest,
	clearMetadataCacheByPattern as clearMetadataCacheByPatternRequest,
	clearMetadataSearchHistory as clearMetadataSearchHistoryRequest,
	deleteIndexerSearchCacheEntry as deleteIndexerSearchCacheEntryRequest,
	deleteMetadataCacheEntry as deleteMetadataCacheEntryRequest,
	getIndexerSearch as getIndexerSearchRequest,
	getMetadataCache as getMetadataCacheRequest,
	updateIndexerSearchSettings as updateIndexerSearchSettingsRequest
} from '$lib/settings/api';
import type { IndexerSearchCacheEntry, MetadataCacheEntry } from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import { createSearchInspectionActions } from './searchInspectionActions';
import type { AppShellState } from './state.svelte';

export function createSearchCacheActions(state: AppShellState, clearNotice: () => void) {
	const inspectionActions = createSearchInspectionActions(state, clearNotice);

	async function saveIndexerSearchSettings(settings = state.indexerSearch.settings) {
		state.savingIndexerSearchSettings = true;
		clearNotice();
		try {
			state.indexerSearch = await updateIndexerSearchSettingsRequest(settings);
			state.message = 'Indexer search settings saved';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save indexer search settings');
		} finally {
			state.savingIndexerSearchSettings = false;
		}
	}

	async function clearIndexerSearchCache() {
		await updateIndexerCache(
			() => clearIndexerSearchCacheRequest(),
			'Indexer search cache reset',
			'Could not reset indexer search cache'
		);
	}

	async function clearIndexerSearchCachePattern(pattern: string) {
		pattern = pattern.trim();
		if (!pattern) return;
		await updateIndexerCache(
			() => clearIndexerSearchCacheByPatternRequest(pattern),
			'Indexer search cache reset',
			'Could not reset matching indexer search cache entries',
			' matching'
		);
	}

	async function deleteIndexerSearchCacheEntry(entry: IndexerSearchCacheEntry) {
		await updateIndexerCache(
			() => deleteIndexerSearchCacheEntryRequest(entry),
			'Indexer search cache entry deleted',
			'Could not delete indexer search cache entry'
		);
	}

	async function clearIndexerSearchHistory() {
		await updateIndexerCache(
			() => clearIndexerSearchHistoryRequest(),
			'Indexer query history cleared',
			'Could not clear indexer query history'
		);
	}

	async function clearMetadataCache() {
		await updateMetadataCache(
			() => clearMetadataCacheRequest(),
			'Metadata cache reset',
			'Could not reset metadata cache'
		);
	}

	async function clearMetadataCachePattern(pattern: string) {
		pattern = pattern.trim();
		if (!pattern) return;
		await updateMetadataCache(
			() => clearMetadataCacheByPatternRequest(pattern),
			'Metadata cache reset',
			'Could not reset matching metadata cache entries',
			' matching'
		);
	}

	async function deleteMetadataCacheEntry(entry: MetadataCacheEntry) {
		await updateMetadataCache(
			() => deleteMetadataCacheEntryRequest(entry),
			'Metadata cache entry deleted',
			'Could not delete metadata cache entry'
		);
	}

	async function clearMetadataSearchHistory() {
		await updateMetadataCache(
			() => clearMetadataSearchHistoryRequest(),
			'Metadata query history cleared',
			'Could not clear metadata query history'
		);
	}

	async function updateIndexerCache(
		deleteRequest: () => Promise<number>,
		successPrefix: string,
		errorText: string,
		countLabel = ''
	) {
		state.clearingIndexerSearchCache = true;
		clearNotice();
		try {
			const deletedCount = await deleteRequest();
			state.indexerSearch = await getIndexerSearchRequest();
			state.message = `${successPrefix}: ${deletedCount}${countLabel} entries deleted`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, errorText);
		} finally {
			state.clearingIndexerSearchCache = false;
		}
	}

	async function updateMetadataCache(
		deleteRequest: () => Promise<number>,
		successPrefix: string,
		errorText: string,
		countLabel = ''
	) {
		state.clearingMetadataCache = true;
		clearNotice();
		try {
			const deletedCount = await deleteRequest();
			state.metadataCache = await getMetadataCacheRequest();
			state.message = `${successPrefix}: ${deletedCount}${countLabel} entries deleted`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, errorText);
		} finally {
			state.clearingMetadataCache = false;
		}
	}

	return {
		...inspectionActions,
		saveIndexerSearchSettings,
		clearIndexerSearchCache,
		clearIndexerSearchCachePattern,
		deleteIndexerSearchCacheEntry,
		clearIndexerSearchHistory,
		clearMetadataCache,
		clearMetadataCachePattern,
		deleteMetadataCacheEntry,
		clearMetadataSearchHistory
	};
}
