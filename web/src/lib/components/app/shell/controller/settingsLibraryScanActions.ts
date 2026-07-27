import {
	advancedSearchMedia,
	importLibraryScanItems as importLibraryScanItemsRequest,
	mediaTypeForLibraryKind,
	resetLibraryScanItemImport as resetLibraryScanItemImportRequest,
	scanLibraryFolder as scanLibraryFolderRequest,
	updateLibraryFolder as updateLibraryFolderRequest
} from '$lib/settings/api';
import type {
	LibraryFolder,
	LibraryMediaKind,
	LibraryScan,
	LibraryScanImportRequest
} from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';
import type { RunCommandMutation } from '$lib/app/query/commandMutation.svelte';

interface SettingsLibraryScanDeps {
	clearNotice: () => void;
	runMutation?: RunCommandMutation;
	refreshMediaItems: () => Promise<void>;
	upsertFolder: (_folder: LibraryFolder) => void;
	upsertScan: (_scan: LibraryScan) => void;
}

export function createSettingsLibraryScanActions(
	state: AppShellState,
	deps: SettingsLibraryScanDeps
) {
	const clearNotice = deps.clearNotice;
	const runMutation = deps.runMutation ?? ((command) => command());
	let scanController: AbortController | undefined;
	let activeScanFolderId: string | undefined;

	async function scanLibraryFolder(id: string) {
		scanController?.abort();
		const controller = new AbortController();
		scanController = controller;
		activeScanFolderId = id;
		state.scanningLibraryFolderId = id;
		clearNotice();

		try {
			const scan = await runMutation(() => scanLibraryFolderRequest(id, controller.signal));
			deps.upsertScan(scan);
			state.openLibraryFolderId = scan.folderId;
			state.message = `Library scan completed: ${scan.manualCount} pending`;
		} catch (error) {
			if (!(error instanceof Error && error.name === 'AbortError')) {
				state.errorMessage = errorMessageFrom(error, 'Could not scan library folder');
			}
		} finally {
			if (scanController === controller) {
				scanController = undefined;
				activeScanFolderId = undefined;
				state.scanningLibraryFolderId = undefined;
			}
		}
	}

	function stopLibraryScan(id: string) {
		if (activeScanFolderId !== id) return;
		scanController?.abort();
		state.scanningLibraryFolderId = undefined;
		state.message = 'Library scan stopped';
	}

	async function updateLibraryFolderManual(id: string, manual: boolean) {
		clearNotice();
		try {
			const folder = await runMutation(() => updateLibraryFolderRequest(id, manual));
			deps.upsertFolder(folder);
			state.message = manual
				? 'Library folder set to manual only'
				: 'Library folder automation enabled';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not update library folder');
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
			const result = await runMutation(() => importLibraryScanItemsRequest(scan.id, request));
			await deps.refreshMediaItems();
			deps.upsertScan(result.scan);
			state.message = `Imported ${result.importedCount} media item${result.importedCount === 1 ? '' : 's'}`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not import library items');
		}
	}

	async function resetLibraryScanImport(scan: LibraryScan, itemId: string) {
		clearNotice();

		try {
			const result = await runMutation(() => resetLibraryScanItemImportRequest(scan.id, itemId));
			if (result.removedMediaItemId) {
				await deps.refreshMediaItems();
			}
			deps.upsertScan(result.scan);
			state.message = 'Library import reset';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not reset library import');
			throw error;
		}
	}

	return {
		scanLibraryFolder,
		stopLibraryScan,
		updateLibraryFolderManual,
		searchLibraryMatch,
		importLibraryScanRows,
		resetLibraryScanImport
	};
}
