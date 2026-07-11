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
	removeLanguage?: (_code: string) => void;
	removeTag?: (_id: string) => void;
	removeUser?: (_id: string) => void;
	removeDownloadClient?: (_id: string) => void;
	removeIndexer?: (_id: string) => void;
	removeSubtitleProvider?: (_id: string) => void;
	removeLibraryFolder?: (_id: string) => void;
	removePathMapping?: (_id: string) => void;
	removeMediaProfile?: (_id: string) => void;
	removeCustomFormat?: (_id: string) => void;
	upsertLibraryScan?: (_scan: import('$lib/settings/types').LibraryScan) => void;
	removeLibraryScan?: (_folderId: string) => void;
}

export function createSettingsDeleteActions(state: AppShellState, deps: SettingsDeleteDeps) {
	const clearNotice = deps.clearNotice;
	const loadSettings = deps.loadSettings;
	const entityDeleteActions = createSettingsEntityDeleteActions(state, {
		clearNotice,
		removeLanguage: deps.removeLanguage,
		removeTag: deps.removeTag,
		removeUser: deps.removeUser,
		removePathMapping: deps.removePathMapping,
		removeMediaProfile: deps.removeMediaProfile,
		removeCustomFormat: deps.removeCustomFormat
	});
	const libraryScanActions = createSettingsLibraryScanActions(state, {
		clearNotice,
		refreshMediaItems: deps.refreshMediaItems ?? (async () => {}),
		upsertScan: deps.upsertLibraryScan ?? (() => {})
	});
	async function deleteDownloadClient(id: string) {
		clearNotice();

		try {
			await deleteDownloadClientRequest(id);
			if (state.downloadForm.id === id) {
				state.downloadForm = emptyDownloadClientForm();
			}
			state.message = 'Download client deleted';
			deps.removeDownloadClient?.(id);
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
			deps.removeIndexer?.(id);
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
			deps.removeSubtitleProvider?.(id);
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete subtitle provider');
		}
	}

	async function deleteLibraryFolder(id: string) {
		clearNotice();

		try {
			await deleteLibraryFolderRequest(id);
			deps.removeLibraryFolder?.(id);
			deps.removeLibraryScan?.(id);
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
