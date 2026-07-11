import {
	testDownloadClientConfig as testDownloadClientConfigRequest,
	testIndexerConfig as testIndexerConfigRequest,
	testIndexer as testIndexerRequest,
	testMetadataProvider as testMetadataProviderRequest,
	testSubtitleProvider as testSubtitleProviderRequest
} from '$lib/settings/api';
import type {
	DownloadClientForm as DownloadClientFormValue,
	IndexerForm as IndexerFormValue
} from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import { createSearchCacheActions } from './searchCacheActions';
import type { AppShellState } from './state.svelte';
import type { RunCommandMutation } from '$lib/app/query/commandMutation.svelte';
import type { QueryClient } from '@tanstack/svelte-query';

interface SettingsTestCacheDeps {
	clearNotice: () => void;
	loadSettings: () => Promise<void>;
	runMutation?: RunCommandMutation;
	queryClient?: QueryClient;
}

export function createSettingsTestCacheActions(state: AppShellState, deps: SettingsTestCacheDeps) {
	const clearNotice = deps.clearNotice;
	const loadSettings = deps.loadSettings;
	const runMutation = deps.runMutation ?? ((command) => command());
	const cacheActions = createSearchCacheActions(state, clearNotice, runMutation, deps.queryClient);
	async function testDownloadClientConfig(form: DownloadClientFormValue) {
		clearNotice();

		try {
			return await runMutation(() => testDownloadClientConfigRequest(form));
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not test download client');
			throw error;
		}
	}

	async function testIndexer(id: string) {
		clearNotice();
		state.testingIndexerId = id;

		try {
			const result = await runMutation(() => testIndexerRequest(id));
			state.indexerTests = { ...state.indexerTests, [id]: result };
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not test indexer');
		} finally {
			state.testingIndexerId = undefined;
		}
	}

	async function testIndexerConfig(form: IndexerFormValue) {
		clearNotice();

		try {
			return await runMutation(() => testIndexerConfigRequest(form));
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not test indexer');
			throw error;
		}
	}

	async function testMetadataProvider(id: string) {
		clearNotice();
		state.testingMetadataProviderId = id;

		try {
			const result = await runMutation(() => testMetadataProviderRequest(id));
			state.metadataProviderTests = { ...state.metadataProviderTests, [id]: result };
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not test metadata provider');
		} finally {
			state.testingMetadataProviderId = undefined;
		}
	}

	async function testSubtitleProvider(id: string) {
		clearNotice();
		state.testingSubtitleProviderId = id;

		try {
			const result = await runMutation(() => testSubtitleProviderRequest(id));
			state.subtitleProviderTests = { ...state.subtitleProviderTests, [id]: result };
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not test subtitle provider');
		} finally {
			state.testingSubtitleProviderId = undefined;
		}
	}

	return {
		testDownloadClientConfig,
		testIndexer,
		testIndexerConfig,
		testMetadataProvider,
		testSubtitleProvider,
		...cacheActions
	};
}
