import { beforeEach, describe, expect, it, vi } from 'vitest';

const apiMock = vi.hoisted(() => ({
	saveLibraryFolder: vi.fn(),
	savePathMapping: vi.fn(),
	saveUser: vi.fn(),
	deleteLibraryFolder: vi.fn(),
	scanLibraryFolder: vi.fn(),
	deletePathMapping: vi.fn(),
	deleteUser: vi.fn(),
	advancedSearchMedia: vi.fn(),
	importLibraryScanItems: vi.fn(),
	mediaTypeForLibraryKind: vi.fn((kind: string) => (kind.includes('series') ? 'serie' : 'movie')),
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
vi.mock('$lib/components/settings/tags/api', () => ({
	saveTag: apiMock.saveTag,
	deleteTag: apiMock.deleteTag
}));
vi.mock('$lib/settings/domains/languages', () => ({
	saveLanguage: apiMock.saveLanguage,
	deleteLanguage: apiMock.deleteLanguage
}));
vi.mock('$lib/settings/domains/users', () => ({
	saveUser: apiMock.saveUser,
	deleteUser: apiMock.deleteUser
}));
vi.mock('$lib/settings/domains/customFormats', () => ({
	saveCustomFormat: apiMock.saveCustomFormat,
	deleteCustomFormat: apiMock.deleteCustomFormat
}));

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
		libraryFolders: [{ id: 'old-folder', path: '/old', kind: 'movie' }],
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
			kind.includes('series') ? 'serie' : 'movie'
		);
	});

	it('adds a library folder, stores its scan, and opens the scan result', async () => {
		const state = shellState();
		const clearNotice = vi.fn();
		const loadSettings = vi.fn();
		const upsertLibraryFolder = vi.fn();
		const upsertLibraryScan = vi.fn();
		apiMock.saveLibraryFolder.mockResolvedValue({
			folder: { id: 'new-folder', path: '/incoming', kind: 'movie' },
			scan: { folderId: 'new-folder', folderKind: 'movie', manualCount: 2 }
		});

		await createSettingsSaveActions(state, {
			clearNotice,
			loadSettings,
			upsertLibraryFolder,
			upsertLibraryScan
		}).saveLibraryFolder(submitEvent());

		expect(clearNotice).toHaveBeenCalledOnce();
		expect(apiMock.saveLibraryFolder).toHaveBeenCalledWith({ path: '/incoming', kind: 'movie' });
		expect(upsertLibraryFolder).toHaveBeenCalledWith(expect.objectContaining({ id: 'new-folder' }));
		expect(upsertLibraryScan).toHaveBeenCalledWith(expect.objectContaining({ manualCount: 2 }));
		expect(state.openLibraryFolderId).toBe('new-folder');
		expect(state.message).toBe('Library scan completed: 2 pending');
		expect(state.savingLibraryFolder).toBe(false);
	});

	it('saves path mappings and reports user-facing errors without leaving spinners active', async () => {
		const state = shellState();
		const actions = createSettingsSaveActions(state, {
			clearNotice: vi.fn(),
			loadSettings: vi.fn(),
			upsertPathMapping: vi.fn()
		});
		apiMock.savePathMapping.mockResolvedValueOnce({
			id: 'map-1',
			hostPath: '/downloads',
			mediaPath: '/media'
		});

		await actions.savePathMapping(submitEvent());
		expect(actions).toBeDefined();
		expect(state.message).toBe('Path mapping saved');

		apiMock.savePathMapping.mockRejectedValueOnce(new Error('Path already exists'));
		await actions.savePathMapping(submitEvent());
		expect(state.errorMessage).toBe('Path already exists');
		expect(state.savingPathMapping).toBe(false);
	});

	it('refreshes settings after saving the current user and updates the session label', async () => {
		const state = shellState();
		const updatedUsers = [
			{
				id: 'user-1',
				username: 'editor-new',
				role: 'admin' as const,
				createdAt: '2026-01-01T00:00:00Z',
				updatedAt: '2026-01-02T00:00:00Z'
			}
		];
		const loadSettings = vi.fn();
		apiMock.saveUser.mockResolvedValue(undefined);

		await createSettingsSaveActions(state, {
			clearNotice: vi.fn(),
			loadSettings,
			users: () => updatedUsers
		}).saveUser(submitEvent());

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
		const loadSettings = vi.fn();
		const removeLibraryFolder = vi.fn();
		const removePathMapping = vi.fn();
		const actions = createSettingsDeleteActions(state, {
			clearNotice: vi.fn(),
			loadSettings,
			removeLibraryFolder,
			removePathMapping
		});

		await actions.deleteLibraryFolder('old-folder');
		await actions.deletePathMapping('old-map');

		expect(removeLibraryFolder).toHaveBeenCalledWith('old-folder');
		expect(state.openLibraryFolderId).toBeUndefined();
		expect(removePathMapping).toHaveBeenCalledWith('old-map');
		expect(loadSettings).toHaveBeenCalledOnce();
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
		const upsertLibraryScan = vi.fn();
		const actions = createSettingsDeleteActions(state, {
			clearNotice: vi.fn(),
			loadSettings: vi.fn(),
			upsertLibraryScan
		});
		apiMock.scanLibraryFolder.mockResolvedValue({ ...scan, manualCount: 3 });
		apiMock.importLibraryScanItems.mockResolvedValue({
			importedCount: 1,
			removedDuplicateCount: 0,
			mediaItems: [{ id: 'media-1', title: 'Imported' }],
			scan: { ...scan, manualCount: 1 }
		});

		await actions.scanLibraryFolder('folder-1');
		await actions.importLibraryScanRows(scan, { items: [] } as never);

		expect(state.openLibraryFolderId).toBe('folder-1');
		expect(upsertLibraryScan).toHaveBeenLastCalledWith(expect.objectContaining({ manualCount: 1 }));
		expect(state.message).toBe('Imported 1 media item');
	});

	it('searches library matches using the media type for the selected library kind', async () => {
		const actions = createSettingsDeleteActions(shellState(), {
			clearNotice: vi.fn(),
			loadSettings: vi.fn()
		});
		apiMock.advancedSearchMedia.mockResolvedValue([
			{ sourceType: 'provider', sourceName: 'TMDB', results: [{ title: 'Provider Series' }] },
			{ sourceType: 'library', sourceName: 'Library', results: [{ title: 'Library Series' }] }
		]);

		await expect(actions.searchLibraryMatch('anime_series', '  Scenario  ')).resolves.toEqual([
			{ title: 'Library Series' },
			{ title: 'Provider Series' }
		]);
		expect(apiMock.advancedSearchMedia).toHaveBeenCalledWith({
			type: 'serie',
			query: 'Scenario',
			includeMedia: true,
			includePeople: false,
			providerIds: undefined,
			limit: 8
		});
	});
});
