import {
	deleteDownloadClient as deleteDownloadClientRequest,
	deleteIndexer as deleteIndexerRequest,
	deleteLibraryFolder as deleteLibraryFolderRequest,
	deleteSubtitleProvider as deleteSubtitleProviderRequest
} from '$lib/settings/api';
import {
	emptyDownloadClientForm,
	emptyIndexerForm,
	emptySubtitleProviderForm
} from '$lib/settings/forms';
import { createSettingsEntityDeleteActions } from './settingsEntityDeleteActions';
import { errorMessageFrom, omitResult } from './helpers';
import { createSettingsLibraryScanActions } from './settingsLibraryScanActions';
import type { AppShellState } from './state.svelte';

interface SettingsDeleteDeps {
	clearNotice: () => void;
	loadSettings: () => Promise<void>;
	refreshMediaItems?: () => Promise<void>;
}

export function createSettingsDeleteActions(state: AppShellState, deps: SettingsDeleteDeps) {
	const clearNotice = deps.clearNotice;
	const loadSettings = deps.loadSettings;
	const entityDeleteActions = createSettingsEntityDeleteActions(state, { clearNotice });
	const libraryScanActions = createSettingsLibraryScanActions(state, {
		clearNotice,
		refreshMediaItems: deps.refreshMediaItems ?? (async () => {})
	});
	async function deleteDownloadClient(id: string) {
		clearNotice();

		try {
			await deleteDownloadClientRequest(id);
			if (state.downloadForm.id === id) {
				state.downloadForm = emptyDownloadClientForm();
			}
			state.message = 'Download client deleted';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete download client');
		}
	}

	async function deleteIndexer(id: string) {
		clearNotice();

		try {
			await deleteIndexerRequest(id);
			if (state.indexerForm.id === id) {
				state.indexerForm = emptyIndexerForm();
			}
			state.indexerTests = omitResult(state.indexerTests, id);
			state.message = 'Indexer deleted';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete indexer');
		}
	}

	async function deleteSubtitleProvider(id: string) {
		clearNotice();

		try {
			await deleteSubtitleProviderRequest(id);
			if (state.subtitleProviderForm.id === id) {
				state.subtitleProviderForm = emptySubtitleProviderForm();
			}
			state.subtitleProviderTests = omitResult(state.subtitleProviderTests, id);
			state.message = 'Subtitle provider deleted';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete subtitle provider');
		}
	}

	async function deleteLibraryFolder(id: string) {
		clearNotice();

		try {
			await deleteLibraryFolderRequest(id);
			state.libraryFolders = state.libraryFolders.filter((folder) => folder.id !== id);
			const remainingScans = { ...state.libraryScansByFolder };
			delete remainingScans[id];
			state.libraryScansByFolder = remainingScans;
			if (state.openLibraryFolderId === id) {
				state.openLibraryFolderId = undefined;
			}
			state.message = 'Library folder deleted';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete library folder');
		}
	}

	return {
		deleteDownloadClient,
		deleteIndexer,
		deleteSubtitleProvider,
		deleteLibraryFolder,
		...entityDeleteActions,
		...libraryScanActions
	};
}
