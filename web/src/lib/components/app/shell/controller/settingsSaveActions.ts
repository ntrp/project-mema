import {
	saveCustomFormat as saveCustomFormatRequest,
	saveDownloadClient as saveDownloadClientRequest,
	saveIndexer as saveIndexerRequest,
	saveLanguage as saveLanguageRequest,
	saveLibraryFolder as saveLibraryFolderRequest,
	saveMediaProfile as saveMediaProfileRequest,
	saveMetadataProvider as saveMetadataProviderRequest,
	savePathMapping as savePathMappingRequest,
	saveSubtitleProvider as saveSubtitleProviderRequest,
	saveTag as saveTagRequest,
	saveUser as saveUserRequest
} from '$lib/settings/api';
import {
	emptyCustomFormatForm,
	emptyDownloadClientForm,
	emptyIndexerForm,
	emptyLanguageForm,
	emptyLibraryFolderForm,
	emptyMediaProfileForm,
	emptyPathMappingForm,
	emptySubtitleProviderForm,
	emptyUserForm
} from '$lib/settings/forms';
import type {
	CustomFormatForm as CustomFormatFormValue,
	MetadataProviderForm as MetadataProviderFormValue,
	SubtitleProviderForm as SubtitleProviderFormValue
} from '$lib/settings/types';
import { emptyTagForm, errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

interface SettingsSaveDeps {
	clearNotice: () => void;
	loadSettings: () => Promise<void>;
}

export function createSettingsSaveActions(state: AppShellState, deps: SettingsSaveDeps) {
	const clearNotice = deps.clearNotice;
	const loadSettings = deps.loadSettings;
	async function saveDownloadClient(event: SubmitEvent) {
		event.preventDefault();
		state.savingDownloadClient = true;
		clearNotice();

		try {
			await saveDownloadClientRequest(state.downloadForm);
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
			await saveIndexerRequest(state.indexerForm);
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
			await saveMetadataProviderRequest(form);
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
			await saveSubtitleProviderRequest(form);
			state.subtitleProviderForm = emptySubtitleProviderForm();
			state.message = 'Subtitle provider saved';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save subtitle provider');
		} finally {
			state.savingSubtitleProviderId = undefined;
		}
	}

	async function saveLibraryFolder(event: SubmitEvent) {
		event.preventDefault();
		state.savingLibraryFolder = true;
		clearNotice();

		try {
			const result = await saveLibraryFolderRequest(state.libraryFolderForm);
			state.libraryFolderForm = emptyLibraryFolderForm();
			state.libraryFolders = [
				result.folder,
				...state.libraryFolders.filter((folder) => folder.id !== result.folder.id)
			];
			state.libraryScansByFolder = {
				...state.libraryScansByFolder,
				[result.folder.id]: result.scan
			};
			state.openLibraryFolderId = result.folder.id;
			state.message = `Library scan completed: ${result.scan.manualCount} pending`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not add library folder');
		} finally {
			state.savingLibraryFolder = false;
		}
	}

	async function savePathMapping(event: SubmitEvent) {
		event.preventDefault();
		state.savingPathMapping = true;
		clearNotice();

		try {
			const mapping = await savePathMappingRequest(state.pathMappingForm);
			state.pathMappingForm = emptyPathMappingForm();
			state.pathMappings = [
				mapping,
				...state.pathMappings.filter((item) => item.id !== mapping.id)
			];
			state.message = 'Path mapping saved';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save path mapping');
		} finally {
			state.savingPathMapping = false;
		}
	}

	async function saveUser(event: SubmitEvent) {
		event.preventDefault();
		state.savingUser = true;
		clearNotice();

		try {
			await saveUserRequest(state.userForm);
			state.userForm = emptyUserForm();
			state.message = 'User saved';
			await loadSettings();
			if (state.currentUser && state.users.some((user) => user.id === state.currentUser?.id)) {
				const updatedUser = state.users.find((user) => user.id === state.currentUser?.id);
				if (updatedUser) {
					state.currentUser = {
						id: updatedUser.id,
						username: updatedUser.username,
						role: updatedUser.role
					};
				}
			}
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save user');
		} finally {
			state.savingUser = false;
		}
	}

	async function saveTag(event: SubmitEvent) {
		event.preventDefault();
		state.savingTag = true;
		clearNotice();

		try {
			await saveTagRequest(state.tagForm);
			state.tagForm = emptyTagForm();
			state.message = 'Tag saved';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save tag');
		} finally {
			state.savingTag = false;
		}
	}

	async function saveLanguage(event: SubmitEvent) {
		event.preventDefault();
		state.savingLanguage = true;
		clearNotice();

		try {
			await saveLanguageRequest(state.languageForm);
			state.languageForm = emptyLanguageForm();
			state.message = 'Language saved';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save language');
		} finally {
			state.savingLanguage = false;
		}
	}

	async function saveMediaProfile(event: SubmitEvent) {
		event.preventDefault();
		state.savingMediaProfile = true;
		clearNotice();

		try {
			await saveMediaProfileRequest(state.mediaProfileForm);
			state.mediaProfileForm = emptyMediaProfileForm();
			state.message = 'Profile saved';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save profile');
		} finally {
			state.savingMediaProfile = false;
		}
	}

	async function saveCustomFormat(event: SubmitEvent) {
		event.preventDefault();
		state.savingCustomFormat = true;
		clearNotice();

		try {
			await saveCustomFormatRequest(state.customFormatForm);
			state.customFormatForm = emptyCustomFormatForm();
			state.message = 'Custom format saved';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save custom format');
		} finally {
			state.savingCustomFormat = false;
		}
	}

	async function importCustomFormat(format: CustomFormatFormValue) {
		state.savingCustomFormat = true;
		clearNotice();

		try {
			await saveCustomFormatRequest(format);
			state.message = 'Custom format imported';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not import custom format');
			throw error;
		} finally {
			state.savingCustomFormat = false;
		}
	}

	return {
		saveDownloadClient,
		saveIndexer,
		saveMetadataProvider,
		saveSubtitleProvider,
		saveLibraryFolder,
		savePathMapping,
		saveUser,
		saveTag,
		saveLanguage,
		saveMediaProfile,
		saveCustomFormat,
		importCustomFormat
	};
}
