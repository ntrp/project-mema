import {
	clearMetadataCache as clearMetadataCacheRequest,
	clearMetadataCacheByPattern as clearMetadataCacheByPatternRequest,
	getMetadataCache as getMetadataCacheRequest,
	testDownloadClientConfig as testDownloadClientConfigRequest,
	testIndexer as testIndexerRequest,
	testMetadataProvider as testMetadataProviderRequest
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

	async function clearMetadataCachePattern(event: SubmitEvent) {
		event.preventDefault();
		const pattern = state.metadataCachePattern.trim();
		if (!pattern) {
			return;
		}
		state.clearingMetadataCache = true;
		clearNotice();

		try {
			const deletedCount = await clearMetadataCacheByPatternRequest(pattern);
			state.metadataCachePattern = '';
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
		refreshMetadataCache,
		clearMetadataCache,
		clearMetadataCachePattern
	};
}
