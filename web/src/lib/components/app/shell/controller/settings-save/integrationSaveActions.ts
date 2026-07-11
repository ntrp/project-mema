import {
	saveDownloadClient as saveDownloadClientRequest,
	saveIndexer as saveIndexerRequest,
	saveMetadataProvider as saveMetadataProviderRequest,
	saveSubtitleProvider as saveSubtitleProviderRequest
} from '$lib/settings/api';
import {
	emptyDownloadClientForm,
	emptyIndexerForm,
	emptySubtitleProviderForm
} from '$lib/settings/forms';
import type {
	MetadataProviderForm as MetadataProviderFormValue,
	SubtitleProviderForm as SubtitleProviderFormValue
} from '$lib/settings/types';
import { errorMessageFrom } from '../helpers';
import type { SettingsSaveContext } from './types';

export function createIntegrationSaveActions({
	state,
	clearNotice,
	loadSettings,
	runMutation = (command) => command()
}: SettingsSaveContext) {
	async function saveDownloadClient(event: SubmitEvent) {
		event.preventDefault();
		state.savingDownloadClient = true;
		clearNotice();

		try {
			await runMutation(() => saveDownloadClientRequest(state.downloadForm));
			state.downloadForm = emptyDownloadClientForm();
			state.message = 'Download client saved';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save download client');
		} finally {
			state.savingDownloadClient = false;
		}
	}

	async function saveIndexer(event: SubmitEvent) {
		event.preventDefault();
		state.savingIndexer = true;
		clearNotice();

		try {
			await runMutation(() => saveIndexerRequest(state.indexerForm));
			state.indexerForm = emptyIndexerForm();
			state.message = 'Indexer saved';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save indexer');
		} finally {
			state.savingIndexer = false;
		}
	}

	async function saveMetadataProvider(form: MetadataProviderFormValue) {
		state.savingMetadataProviderId = form.id;
		clearNotice();

		try {
			await runMutation(() => saveMetadataProviderRequest(form));
			state.message = 'Metadata provider saved';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save metadata provider');
		} finally {
			state.savingMetadataProviderId = undefined;
		}
	}

	async function saveSubtitleProvider(form: SubtitleProviderFormValue) {
		state.savingSubtitleProviderId = form.id;
		clearNotice();

		try {
			await runMutation(() => saveSubtitleProviderRequest(form));
			state.subtitleProviderForm = emptySubtitleProviderForm();
			state.message = 'Subtitle provider saved';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save subtitle provider');
		} finally {
			state.savingSubtitleProviderId = undefined;
		}
	}

	return {
		saveDownloadClient,
		saveIndexer,
		saveMetadataProvider,
		saveSubtitleProvider
	};
}
