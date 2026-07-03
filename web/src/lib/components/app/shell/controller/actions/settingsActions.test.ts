import { beforeEach, describe, expect, it, vi } from 'vitest';

const apiMock = vi.hoisted(() => ({
	saveLibraryFolder: vi.fn(),
	savePathMapping: vi.fn(),
	saveUser: vi.fn(),
	deleteLibraryFolder: vi.fn(),
	scanLibraryFolder: vi.fn(),
	deletePathMapping: vi.fn(),
	deleteUser: vi.fn(),
	searchMedia: vi.fn(),
	matchLibraryScanItem: vi.fn(),
	mediaTypeForLibraryKind: vi.fn((kind: string) => (kind.includes('series') ? 'series' : 'movie')),
	saveCustomFormat: vi.fn(),
	saveDownloadClient: vi.fn(),
	saveIndexer: vi.fn(),
	saveLanguage: vi.fn(),
	saveMediaProfile: vi.fn(),
	saveMetadataProvider: vi.fn(),
	saveTag: vi.fn(),
	deleteCustomFormat: vi.fn(),
	deleteDownloadClient: vi.fn(),
	deleteIndexer: vi.fn(),
	deleteLanguage: vi.fn(),
	deleteMediaProfile: vi.fn(),
	deleteTag: vi.fn()
}));

vi.mock('$lib/settings/api', () => apiMock);

import { createSettingsDeleteActions } from '../settingsDeleteActions';
import { createSettingsSaveActions } from '../settingsSaveActions';
import type { AppShellState } from '../state.svelte';
import type { LibraryScan } from '$lib/settings/types';

function submitEvent() {
	return { preventDefault: vi.fn() } as unknown as SubmitEvent;
}

function shellState(overrides: Record<string, unknown> = {}) {
	return {
		message: '',
		errorMessage: '',
		libraryFolderForm: { path: '/incoming', kind: 'movie' },
		pathMappingForm: { hostPath: '/downloads', mediaPath: '/media' },
		userForm: { id: 'user-1', username: 'editor', role: 'user', password: '' },
		libraryFolders: [{ id: 'old-folder', path: '/old' }],
		pathMappings: [{ id: 'old-map', hostPath: '/old', mediaPath: '/media-old' }],
		libraryScansByFolder: {},
		openLibraryFolderId: undefined,
		users: [{ id: 'user-1', username: 'old', role: 'user' }],
		currentUser: { id: 'user-1', username: 'old', role: 'user' },
		mediaItems: [{ id: 'existing-media' }],
		...overrides
	} as unknown as AppShellState;
}

describe('settings save actions (SCN-SETTINGS-009)', () => {
	beforeEach(() => {
		for (const value of Object.values(apiMock)) {
			if (typeof value === 'function' && 'mockReset' in value) value.mockReset();
		}
		apiMock.mediaTypeForLibraryKind.mockImplementation((kind: string) =>
			kind.includes('series') ? 'series' : 'movie'
		);
	});

	it('adds a library folder, stores its scan, and opens the scan result', async () => {
		const state = shellState();
		const clearNotice = vi.fn();
		const loadSettings = vi.fn();
		apiMock.saveLibraryFolder.mockResolvedValue({
			folder: { id: 'new-folder', path: '/incoming' },
			scan: { folderId: 'new-folder', manualCount: 2 }
		});

		await createSettingsSaveActions(state, { clearNotice, loadSettings }).saveLibraryFolder(
			submitEvent()
		);

		expect(clearNotice).toHaveBeenCalledOnce();
		expect(apiMock.saveLibraryFolder).toHaveBeenCalledWith({ path: '/incoming', kind: 'movie' });
		expect(state.libraryFolders.map((folder) => folder.id)).toEqual(['new-folder', 'old-folder']);
		expect(state.libraryScansByFolder).toMatchObject({ 'new-folder': { manualCount: 2 } });
		expect(state.openLibraryFolderId).toBe('new-folder');
		expect(state.message).toBe('Library scan completed: 2 pending');
		expect(state.savingLibraryFolder).toBe(false);
	});

	it('saves path mappings and reports user-facing errors without leaving spinners active', async () => {
		const state = shellState();
		const actions = createSettingsSaveActions(state, {
			clearNotice: vi.fn(),
			loadSettings: vi.fn()
		});
		apiMock.savePathMapping.mockResolvedValueOnce({
			id: 'map-1',
			hostPath: '/downloads',
			mediaPath: '/media'
		});

		await actions.savePathMapping(submitEvent());
		expect(state.pathMappings.map((mapping) => mapping.id)).toEqual(['map-1', 'old-map']);
		expect(state.message).toBe('Path mapping saved');

		apiMock.savePathMapping.mockRejectedValueOnce(new Error('Path already exists'));
		await actions.savePathMapping(submitEvent());
		expect(state.errorMessage).toBe('Path already exists');
		expect(state.savingPathMapping).toBe(false);
	});

	it('refreshes settings after saving the current user and updates the session label', async () => {
		const state = shellState();
		const loadSettings = vi.fn(async () => {
			state.users = [
				{
					id: 'user-1',
					username: 'editor-new',
					role: 'admin',
					createdAt: '2026-01-01T00:00:00Z',
					updatedAt: '2026-01-02T00:00:00Z'
				}
			];
		});
		apiMock.saveUser.mockResolvedValue(undefined);

		await createSettingsSaveActions(state, { clearNotice: vi.fn(), loadSettings }).saveUser(
			submitEvent()
		);

		expect(loadSettings).toHaveBeenCalledOnce();
		expect(state.currentUser).toEqual({ id: 'user-1', username: 'editor-new', role: 'admin' });
		expect(state.message).toBe('User saved');
	});
});

describe('settings delete and import actions (SCN-SETTINGS-009)', () => {
	it('deletes library folders and path mappings from visible state', async () => {
		const state = shellState({
			openLibraryFolderId: 'old-folder',
			libraryScansByFolder: { 'old-folder': { folderId: 'old-folder' } }
		});
		const actions = createSettingsDeleteActions(state, {
			clearNotice: vi.fn(),
			loadSettings: vi.fn()
		});

		await actions.deleteLibraryFolder('old-folder');
		await actions.deletePathMapping('old-map');

		expect(state.libraryFolders).toEqual([]);
		expect(state.libraryScansByFolder).toEqual({});
		expect(state.openLibraryFolderId).toBeUndefined();
		expect(state.pathMappings).toEqual([]);
		expect(state.message).toBe('Path mapping deleted');
	});

	it('scans and imports selected library rows with updated scan counts', async () => {
		const scan = {
			id: 'scan-1',
			folderId: 'folder-1',
			manualCount: 2,
			items: [{ id: 'item-1' }, { id: 'item-2' }]
		} as LibraryScan;
		const state = shellState();
		const actions = createSettingsDeleteActions(state, {
			clearNotice: vi.fn(),
			loadSettings: vi.fn()
		});
		apiMock.scanLibraryFolder.mockResolvedValue({ ...scan, manualCount: 3 });
		apiMock.matchLibraryScanItem.mockResolvedValue({
			mediaItem: { id: 'media-1', title: 'Imported' },
			item: { id: 'item-1', matched: true }
		});

		await actions.scanLibraryFolder('folder-1');
		await actions.importLibraryScanRows(scan, [
			{ item: { id: 'item-1' }, request: { type: 'movie', monitored: true } } as never
		]);

		expect(state.openLibraryFolderId).toBe('folder-1');
		expect(state.mediaItems.map((item) => item.id)).toEqual(['media-1', 'existing-media']);
		expect(state.libraryScansByFolder['folder-1'].manualCount).toBe(1);
		expect(state.message).toBe('Imported 1 media item');
	});

	it('searches library matches using the media type for the selected library kind', async () => {
		const actions = createSettingsDeleteActions(shellState(), {
			clearNotice: vi.fn(),
			loadSettings: vi.fn()
		});
		apiMock.searchMedia.mockResolvedValue([{ title: 'Scenario Series' }]);

		await expect(actions.searchLibraryMatch('anime_series', '  Scenario  ')).resolves.toEqual([
			{ title: 'Scenario Series' }
		]);
		expect(apiMock.searchMedia).toHaveBeenCalledWith({ type: 'series', query: 'Scenario' });
	});
});
