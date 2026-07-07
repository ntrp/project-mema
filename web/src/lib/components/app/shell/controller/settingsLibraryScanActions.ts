import {
	advancedSearchMedia,
	importLibraryScanItems as importLibraryScanItemsRequest,
	mediaTypeForLibraryKind,
	resetLibraryScanItemImport as resetLibraryScanItemImportRequest,
	scanLibraryFolder as scanLibraryFolderRequest
} from '$lib/settings/api';
import type { LibraryMediaKind, LibraryScan, LibraryScanImportRequest } from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

interface SettingsLibraryScanDeps {
	clearNotice: () => void;
}

export function createSettingsLibraryScanActions(
	state: AppShellState,
	deps: SettingsLibraryScanDeps
) {
	const clearNotice = deps.clearNotice;

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

	async function searchLibraryMatch(kind: LibraryMediaKind, query: string, providerId?: string) {
		const groups = await advancedSearchMedia({
			type: mediaTypeForLibraryKind(kind),
			query: query.trim(),
			includeMedia: true,
			includePeople: false,
			providerIds: providerId ? [providerId] : undefined,
			limit: 8
		});
		return [
			...groups.filter((group) => group.sourceType === 'library').flatMap((group) => group.results),
			...groups.filter((group) => group.sourceType !== 'library').flatMap((group) => group.results)
		];
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

	async function resetLibraryScanImport(scan: LibraryScan, itemId: string) {
		clearNotice();

		try {
			const result = await resetLibraryScanItemImportRequest(scan.id, itemId);
			if (result.removedMediaItemId) {
				state.mediaItems = state.mediaItems.filter((item) => item.id !== result.removedMediaItemId);
			}
			state.libraryScansByFolder = {
				...state.libraryScansByFolder,
				[scan.folderId]: result.scan
			};
			state.message = 'Library import reset';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not reset library import');
			throw error;
		}
	}

	return {
		scanLibraryFolder,
		searchLibraryMatch,
		importLibraryScanRows,
		resetLibraryScanImport
	};
}
