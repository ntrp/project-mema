import {
	deleteCustomFormat as deleteCustomFormatRequest,
	deleteDownloadClient as deleteDownloadClientRequest,
	deleteIndexer as deleteIndexerRequest,
	deleteLanguage as deleteLanguageRequest,
	deleteLibraryFolder as deleteLibraryFolderRequest,
	deleteMediaProfile as deleteMediaProfileRequest,
	deletePathMapping as deletePathMappingRequest,
	deleteSubtitleProvider as deleteSubtitleProviderRequest,
	deleteTag as deleteTagRequest,
	deleteUser as deleteUserRequest,
	advancedSearchMedia,
	importLibraryScanItems as importLibraryScanItemsRequest,
	mediaTypeForLibraryKind,
	scanLibraryFolder as scanLibraryFolderRequest
} from '$lib/settings/api';
import {
	emptyCustomFormatForm,
	emptyDownloadClientForm,
	emptyIndexerForm,
	emptyLanguageForm,
	emptyMediaProfileForm,
	emptySubtitleProviderForm,
	emptyUserForm
} from '$lib/settings/forms';
import type { LibraryMediaKind, LibraryScan, LibraryScanImportRequest } from '$lib/settings/types';
import { emptyTagForm, errorMessageFrom, omitResult } from './helpers';
import type { AppShellState } from './state.svelte';

interface SettingsDeleteDeps {
	clearNotice: () => void;
	loadSettings: () => Promise<void>;
}

export function createSettingsDeleteActions(state: AppShellState, deps: SettingsDeleteDeps) {
	const clearNotice = deps.clearNotice;
	const loadSettings = deps.loadSettings;
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
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete library folder');
		}
	}

	async function scanLibraryFolder(id: string) {
		state.scanningLibraryFolderId = id;
		clearNotice();

		try {
			const scan = await scanLibraryFolderRequest(id);
			state.libraryScansByFolder = { ...state.libraryScansByFolder, [scan.folderId]: scan };
			state.openLibraryFolderId = scan.folderId;
			state.message = `Library scan completed: ${scan.manualCount} pending`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not scan library folder');
		} finally {
			state.scanningLibraryFolderId = undefined;
		}
	}

	async function deletePathMapping(id: string) {
		state.deletingPathMappingId = id;
		clearNotice();

		try {
			await deletePathMappingRequest(id);
			state.pathMappings = state.pathMappings.filter((mapping) => mapping.id !== id);
			state.message = 'Path mapping deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete path mapping');
		} finally {
			state.deletingPathMappingId = undefined;
		}
	}

	async function deleteUser(id: string) {
		clearNotice();

		try {
			await deleteUserRequest(id);
			if (state.userForm.id === id) {
				state.userForm = emptyUserForm();
			}
			state.users = state.users.filter((user) => user.id !== id);
			state.message = 'User deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete user');
		}
	}

	async function deleteTag(id: string) {
		state.deletingTagId = id;
		clearNotice();

		try {
			await deleteTagRequest(id);
			if (state.tagForm.id === id) {
				state.tagForm = emptyTagForm();
			}
			state.tags = state.tags.filter((tag) => tag.id !== id);
			state.message = 'Tag deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete tag');
		} finally {
			state.deletingTagId = undefined;
		}
	}

	async function deleteLanguage(code: string) {
		state.deletingLanguageCode = code;
		clearNotice();

		try {
			await deleteLanguageRequest(code);
			if (state.languageForm.originalCode === code) {
				state.languageForm = emptyLanguageForm();
			}
			state.languages = state.languages.filter((language) => language.code !== code);
			state.message = 'Language deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete language');
		} finally {
			state.deletingLanguageCode = undefined;
		}
	}

	async function deleteMediaProfile(id: string) {
		state.deletingMediaProfileId = id;
		clearNotice();

		try {
			await deleteMediaProfileRequest(id);
			if (state.mediaProfileForm.id === id) {
				state.mediaProfileForm = emptyMediaProfileForm();
			}
			state.mediaProfiles = state.mediaProfiles.filter((profile) => profile.id !== id);
			state.message = 'Profile deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete profile');
		} finally {
			state.deletingMediaProfileId = undefined;
		}
	}

	async function deleteCustomFormat(id: string) {
		state.deletingCustomFormatId = id;
		clearNotice();

		try {
			await deleteCustomFormatRequest(id);
			if (state.customFormatForm.id === id) {
				state.customFormatForm = emptyCustomFormatForm();
			}
			state.customFormats = state.customFormats.filter((format) => format.id !== id);
			state.message = 'Custom format deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete custom format');
		} finally {
			state.deletingCustomFormatId = undefined;
		}
	}

	async function searchLibraryMatch(kind: LibraryMediaKind, query: string, providerId?: string) {
		const groups = await advancedSearchMedia({
			type: mediaTypeForLibraryKind(kind),
			query: query.trim(),
			includeMedia: true,
			includePeople: false,
			providerIds: providerId ? [providerId] : undefined,
			limit: 8
		});
		return groups.flatMap((group) => group.results);
	}

	async function importLibraryScanRows(scan: LibraryScan, request: LibraryScanImportRequest) {
		clearNotice();

		try {
			const result = await importLibraryScanItemsRequest(scan.id, request);
			const importedMediaIds = result.mediaItems.map((item) => item.id);
			state.mediaItems = [
				...result.mediaItems,
				...state.mediaItems.filter((item) => !importedMediaIds.includes(item.id))
			];
			state.libraryScansByFolder = {
				...state.libraryScansByFolder,
				[scan.folderId]: result.scan
			};
			state.message = `Imported ${result.importedCount} media item${result.importedCount === 1 ? '' : 's'}`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not import library items');
		}
	}

	return {
		deleteDownloadClient,
		deleteIndexer,
		deleteSubtitleProvider,
		deleteLibraryFolder,
		scanLibraryFolder,
		deletePathMapping,
		deleteUser,
		deleteTag,
		deleteLanguage,
		deleteMediaProfile,
		deleteCustomFormat,
		searchLibraryMatch,
		importLibraryScanRows
	};
}
