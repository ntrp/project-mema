import { beforeEach, describe, expect, it, vi } from 'vitest';

const apiMock = vi.hoisted(() => ({
	saveCustomFormat: vi.fn(),
	saveDownloadClient: vi.fn(),
	saveIndexer: vi.fn(),
	saveLanguage: vi.fn(),
	saveLibraryFolder: vi.fn(),
	saveMediaProfile: vi.fn(),
	saveMetadataProvider: vi.fn(),
	savePathMapping: vi.fn(),
	saveTag: vi.fn(),
	saveUser: vi.fn(),
	deleteCustomFormat: vi.fn(),
	deleteDownloadClient: vi.fn(),
	deleteIndexer: vi.fn(),
	deleteLanguage: vi.fn(),
	deleteLibraryFolder: vi.fn(),
	deleteMediaProfile: vi.fn(),
	deletePathMapping: vi.fn(),
	deleteTag: vi.fn(),
	deleteUser: vi.fn(),
	matchLibraryScanItem: vi.fn(),
	mediaTypeForLibraryKind: vi.fn(),
	scanLibraryFolder: vi.fn(),
	searchMedia: vi.fn()
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

function submitEvent() {
	return { preventDefault: vi.fn() } as unknown as SubmitEvent;
}

function shellState(overrides: Record<string, unknown> = {}) {
	return {
		message: '',
		errorMessage: '',
		downloadForm: { id: 'download-1', name: 'Transmission' },
		indexerForm: { id: 'indexer-1', name: 'Torznab' },
		languageForm: { code: 'de', originalCode: 'de', displayName: 'German' },
		mediaProfileForm: { id: 'profile-1', name: 'HD' },
		metadataProviderForm: { id: 'metadata-1', name: 'TMDB' },
		tagForm: { id: 'tag-1', name: 'Kids' },
		customFormatForm: { id: 'format-1', name: 'HDR' },
		downloadClients: [{ id: 'download-1' }],
		indexers: [{ id: 'indexer-1' }],
		indexerTests: { 'indexer-1': { status: 'success' } },
		languages: [{ code: 'de' }, { code: 'en' }],
		mediaProfiles: [{ id: 'profile-1' }],
		mediaItems: [],
		tags: [{ id: 'tag-1' }],
		customFormats: [{ id: 'format-1' }],
		users: [],
		...overrides
	} as unknown as AppShellState;
}

describe('settings entity save actions (SCN-SETTINGS-009)', () => {
	beforeEach(() => {
		for (const value of Object.values(apiMock)) {
			if (typeof value === 'function' && 'mockReset' in value) value.mockReset();
		}
	});

	it('saves configured integrations and reloads settings', async () => {
		const state = shellState();
		const loadSettings = vi.fn();
		const actions = createSettingsSaveActions(state, {
			clearNotice: vi.fn(),
			loadSettings
		});
		apiMock.saveDownloadClient.mockResolvedValue(undefined);
		apiMock.saveIndexer.mockResolvedValue(undefined);
		apiMock.saveMetadataProvider.mockResolvedValue(undefined);

		await actions.saveDownloadClient(submitEvent());
		await actions.saveIndexer(submitEvent());
		await actions.saveMetadataProvider({
			id: 'metadata-1',
			name: 'TMDB',
			type: 'tmdb',
			baseUrl: 'https://metadata.test',
			enabled: true,
			priority: 100
		});

		expect(apiMock.saveDownloadClient).toHaveBeenCalledWith({
			id: 'download-1',
			name: 'Transmission'
		});
		expect(apiMock.saveIndexer).toHaveBeenCalledWith({ id: 'indexer-1', name: 'Torznab' });
		expect(apiMock.saveMetadataProvider).toHaveBeenCalledWith({
			id: 'metadata-1',
			name: 'TMDB',
			type: 'tmdb',
			baseUrl: 'https://metadata.test',
			enabled: true,
			priority: 100
		});
		expect(loadSettings).toHaveBeenCalledTimes(3);
		expect(state.message).toBe('Metadata provider saved');
		expect(state.savingDownloadClient).toBe(false);
		expect(state.savingIndexer).toBe(false);
		expect(state.savingMetadataProviderId).toBeUndefined();
	});

	it('saves catalog entities and keeps user-visible errors', async () => {
		const state = shellState();
		const actions = createSettingsSaveActions(state, {
			clearNotice: vi.fn(),
			loadSettings: vi.fn()
		});
		apiMock.saveTag.mockResolvedValue(undefined);
		apiMock.saveLanguage.mockResolvedValue(undefined);
		apiMock.saveMediaProfile.mockResolvedValue(undefined);
		apiMock.saveCustomFormat.mockRejectedValueOnce(new Error('Format name exists'));

		await actions.saveTag(submitEvent());
		await actions.saveLanguage(submitEvent());
		await actions.saveMediaProfile(submitEvent());
		await actions.saveCustomFormat(submitEvent());

		expect(apiMock.saveTag).toHaveBeenCalledWith({ id: 'tag-1', name: 'Kids' });
		expect(apiMock.saveLanguage).toHaveBeenCalledWith({
			code: 'de',
			originalCode: 'de',
			displayName: 'German'
		});
		expect(apiMock.saveMediaProfile).toHaveBeenCalledWith({ id: 'profile-1', name: 'HD' });
		expect(state.message).toBe('Profile saved');
		expect(state.errorMessage).toBe('Format name exists');
		expect(state.savingTag).toBe(false);
		expect(state.savingLanguage).toBe(false);
		expect(state.savingMediaProfile).toBe(false);
		expect(state.savingCustomFormat).toBe(false);
	});

	it('reloads media items using the saved profile so track status is recalculated', async () => {
		const state = shellState({
			mediaItems: [
				{ id: 'media-1', qualityProfileId: 'profile-1' },
				{ id: 'media-2', qualityProfileId: 'profile-2' }
			]
		});
		const loadSettings = vi.fn();
		const loadMediaItems = vi.fn();
		const actions = createSettingsSaveActions(state, {
			clearNotice: vi.fn(),
			loadSettings,
			loadMediaItems,
			mediaItems: () =>
				(state as unknown as { mediaItems: import('$lib/settings/types').MediaItem[] }).mediaItems
		});
		apiMock.saveMediaProfile.mockResolvedValue(undefined);

		await actions.saveMediaProfile(submitEvent());

		expect(apiMock.saveMediaProfile).toHaveBeenCalledWith({ id: 'profile-1', name: 'HD' });
		expect(loadSettings).toHaveBeenCalledOnce();
		expect(loadMediaItems).toHaveBeenCalledOnce();
	});

	it('rethrows failed custom format imports after recording the error', async () => {
		const state = shellState();
		const actions = createSettingsSaveActions(state, {
			clearNotice: vi.fn(),
			loadSettings: vi.fn()
		});
		apiMock.saveCustomFormat.mockRejectedValue(new Error('Bad import'));

		await expect(
			actions.importCustomFormat({
				name: 'Imported',
				includeInRenameTemplate: false,
				includeSpecs: [],
				excludeSpecs: []
			})
		).rejects.toThrow('Bad import');

		expect(state.errorMessage).toBe('Bad import');
		expect(state.savingCustomFormat).toBe(false);
	});
});

describe('settings entity delete actions (SCN-SETTINGS-009)', () => {
	it('deletes configured integrations and clears matching edit state', async () => {
		const state = shellState();
		const loadSettings = vi.fn();
		const actions = createSettingsDeleteActions(state, {
			clearNotice: vi.fn(),
			loadSettings
		});
		apiMock.deleteDownloadClient.mockResolvedValue(undefined);
		apiMock.deleteIndexer.mockResolvedValue(undefined);

		await actions.deleteDownloadClient('download-1');
		await actions.deleteIndexer('indexer-1');

		expect(loadSettings).toHaveBeenCalledTimes(2);
		expect(state.indexerTests).toEqual({});
		expect(state.message).toBe('Indexer deleted');
	});

	it('deletes catalog entities from visible state and resets spinners', async () => {
		const state = shellState();
		const removeTag = vi.fn();
		const removeLanguage = vi.fn();
		const removeMediaProfile = vi.fn();
		const removeCustomFormat = vi.fn();
		const actions = createSettingsDeleteActions(state, {
			clearNotice: vi.fn(),
			loadSettings: vi.fn(),
			removeTag,
			removeLanguage,
			removeMediaProfile,
			removeCustomFormat
		});

		await actions.deleteTag('tag-1');
		await actions.deleteLanguage('de');
		await actions.deleteMediaProfile('profile-1');
		await actions.deleteCustomFormat('format-1');

		expect(removeTag).toHaveBeenCalledWith('tag-1');
		expect(removeLanguage).toHaveBeenCalledWith('de');
		expect(removeMediaProfile).toHaveBeenCalledWith('profile-1');
		expect(removeCustomFormat).toHaveBeenCalledWith('format-1');
		expect(state.message).toBe('Custom format deleted');
		expect(state.deletingTagId).toBeUndefined();
		expect(state.deletingLanguageCode).toBeUndefined();
		expect(state.deletingMediaProfileId).toBeUndefined();
		expect(state.deletingCustomFormatId).toBeUndefined();
	});
});
