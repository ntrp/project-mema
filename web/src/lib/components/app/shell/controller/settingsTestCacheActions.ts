import {
	testDownloadClientConfig as testDownloadClientConfigRequest,
	testIndexer as testIndexerRequest,
	testMetadataProvider as testMetadataProviderRequest
} from '$lib/settings/api';
import type { DownloadClientForm as DownloadClientFormValue } from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import { createSearchCacheActions } from './searchCacheActions';
import type { AppShellState } from './state.svelte';

interface SettingsTestCacheDeps {
	clearNotice: () => void;
	loadSettings: () => Promise<void>;
}

export function createSettingsTestCacheActions(state: AppShellState, deps: SettingsTestCacheDeps) {
	const clearNotice = deps.clearNotice;
	const loadSettings = deps.loadSettings;
	const cacheActions = createSearchCacheActions(state, clearNotice);
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

	return {
		testDownloadClientConfig,
		testIndexer,
		testMetadataProvider,
		...cacheActions
	};
}
