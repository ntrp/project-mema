import {
	clearIndexerSearchCache as clearIndexerSearchCacheRequest,
	clearIndexerSearchCacheByPattern as clearIndexerSearchCacheByPatternRequest,
	clearMetadataCache as clearMetadataCacheRequest,
	clearMetadataCacheByPattern as clearMetadataCacheByPatternRequest,
	getIndexerSearch as getIndexerSearchRequest,
	getMetadataCache as getMetadataCacheRequest,
	testDownloadClientConfig as testDownloadClientConfigRequest,
	testIndexer as testIndexerRequest,
	testMetadataProvider as testMetadataProviderRequest,
	updateIndexerSearchSettings as updateIndexerSearchSettingsRequest
} from '$lib/settings/api';
import type { DownloadClientForm as DownloadClientFormValue } from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

interface SettingsTestCacheDeps {
	clearNotice: () => void;
	loadSettings: () => Promise<void>;
}

export function createSettingsTestCacheActions(state: AppShellState, deps: SettingsTestCacheDeps) {
	const clearNotice = deps.clearNotice;
	const loadSettings = deps.loadSettings;
	async function testDownloadClientConfig(form: DownloadClientFormValue) {
		clearNotice();

		try {
			return await testDownloadClientConfigRequest(form);
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not test download client');
			throw error;
		}
	}

	async function testIndexer(id: string) {
		clearNotice();
		state.testingIndexerId = id;

		try {
			const result = await testIndexerRequest(id);
			state.indexerTests = { ...state.indexerTests, [id]: result };
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not test indexer');
		} finally {
			state.testingIndexerId = undefined;
		}
	}

	async function testMetadataProvider(id: string) {
		clearNotice();
		state.testingMetadataProviderId = id;

		try {
			const result = await testMetadataProviderRequest(id);
			state.metadataProviderTests = { ...state.metadataProviderTests, [id]: result };
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not test metadata provider');
		} finally {
			state.testingMetadataProviderId = undefined;
		}
	}

	async function refreshMetadataCache() {
		state.loadingMetadataCache = true;
		clearNotice();

		try {
			state.metadataCache = await getMetadataCacheRequest();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load metadata cache');
		} finally {
			state.loadingMetadataCache = false;
		}
	}

	async function refreshIndexerSearch() {
		state.loadingIndexerSearch = true;
		clearNotice();

		try {
			state.indexerSearch = await getIndexerSearchRequest();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load indexer search cache');
		} finally {
			state.loadingIndexerSearch = false;
		}
	}

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
		state.clearingIndexerSearchCache = true;
		clearNotice();

		try {
			const deletedCount = await clearIndexerSearchCacheRequest();
			state.indexerSearch = await getIndexerSearchRequest();
			state.message = `Indexer search cache reset: ${deletedCount} entries deleted`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not reset indexer search cache');
		} finally {
			state.clearingIndexerSearchCache = false;
		}
	}

	async function clearIndexerSearchCachePattern(pattern: string) {
		pattern = pattern.trim();
		if (!pattern) {
			return;
		}
		state.clearingIndexerSearchCache = true;
		clearNotice();

		try {
			const deletedCount = await clearIndexerSearchCacheByPatternRequest(pattern);
			state.indexerSearch = await getIndexerSearchRequest();
			state.message = `Indexer search cache reset: ${deletedCount} matching entries deleted`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(
				error,
				'Could not reset matching indexer search cache entries'
			);
		} finally {
			state.clearingIndexerSearchCache = false;
		}
	}

	async function clearMetadataCache() {
		state.clearingMetadataCache = true;
		clearNotice();

		try {
			const deletedCount = await clearMetadataCacheRequest();
			state.metadataCache = await getMetadataCacheRequest();
			state.message = `Metadata cache reset: ${deletedCount} entries deleted`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not reset metadata cache');
		} finally {
			state.clearingMetadataCache = false;
		}
	}

	async function clearMetadataCachePattern(pattern: string) {
		pattern = pattern.trim();
		if (!pattern) {
			return;
		}
		state.clearingMetadataCache = true;
		clearNotice();

		try {
			const deletedCount = await clearMetadataCacheByPatternRequest(pattern);
			state.metadataCache = await getMetadataCacheRequest();
			state.message = `Metadata cache reset: ${deletedCount} matching entries deleted`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(
				error,
				'Could not reset matching metadata cache entries'
			);
		} finally {
			state.clearingMetadataCache = false;
		}
	}

	return {
		testDownloadClientConfig,
		testIndexer,
		testMetadataProvider,
		refreshIndexerSearch,
		saveIndexerSearchSettings,
		clearIndexerSearchCache,
		clearIndexerSearchCachePattern,
		refreshMetadataCache,
		clearMetadataCache,
		clearMetadataCachePattern
	};
}
