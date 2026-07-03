import { beforeEach, describe, expect, it, vi } from 'vitest';

const apiMock = vi.hoisted(() => ({
	cancelDownloadActivity: vi.fn(),
	clearIndexerSearchCache: vi.fn(),
	clearIndexerSearchCacheByPattern: vi.fn(),
	clearMetadataCacheByPattern: vi.fn(),
	deleteDownloadActivity: vi.fn(),
	getIndexerSearch: vi.fn(),
	getMetadataCache: vi.fn(),
	testDownloadClientConfig: vi.fn(),
	testIndexer: vi.fn(),
	testMetadataProvider: vi.fn(),
	updateIndexerSearchSettings: vi.fn()
}));

vi.mock('$lib/settings/api', () => apiMock);

import { createActivityActions } from '../activityActions';
import { createSearchCacheActions } from '../searchCacheActions';
import { createSettingsTestCacheActions } from '../settingsTestCacheActions';
import type { AppShellState } from '../state.svelte';

function state(overrides: Record<string, unknown> = {}) {
	return {
		message: '',
		errorMessage: '',
		indexerSearch: { settings: { retentionDays: 7 }, cacheEntries: [], historyEntries: [] },
		metadataCache: { entries: [], historyEntries: [] },
		indexerTests: {},
		metadataProviderTests: {},
		activities: [{ id: 'activity-1' }, { id: 'activity-2' }],
		...overrides
	} as unknown as AppShellState;
}

describe('cache and inspection actions (SCN-SYSTEM-004)', () => {
	beforeEach(() => {
		for (const value of Object.values(apiMock)) value.mockReset();
	});

	it('loads expanded cache inspection pages using stable page sizes', async () => {
		const shell = state({
			indexerSearch: {
				settings: {},
				cacheEntries: [{ key: 'cached' }],
				historyEntries: [{ query: 'old' }]
			},
			metadataCache: {
				entries: [{ key: 'movie:scenario' }, { key: 'movie:other' }],
				historyEntries: []
			}
		});
		apiMock.getIndexerSearch.mockResolvedValue({
			settings: {},
			cacheEntries: [],
			historyEntries: []
		});
		apiMock.getMetadataCache.mockResolvedValue({ entries: [], historyEntries: [] });
		const actions = createSearchCacheActions(shell, vi.fn());

		await actions.loadMoreIndexerSearchCache();
		await actions.loadMoreMetadataSearchHistory();

		expect(apiMock.getIndexerSearch).toHaveBeenCalledWith({ cacheLimit: 11, historyLimit: 1 });
		expect(apiMock.getMetadataCache).toHaveBeenCalledWith({ cacheLimit: 2, historyLimit: 20 });
		expect(shell.loadingIndexerSearch).toBe(false);
		expect(shell.loadingMetadataCache).toBe(false);
	});

	it('saves indexer settings and refreshes cache state after pattern clears', async () => {
		const shell = state();
		apiMock.updateIndexerSearchSettings.mockResolvedValue({ settings: { retentionDays: 14 } });
		apiMock.clearIndexerSearchCacheByPattern.mockResolvedValue(2);
		apiMock.clearMetadataCacheByPattern.mockResolvedValue(3);
		apiMock.getIndexerSearch.mockResolvedValue({ settings: {}, cacheEntries: [] });
		apiMock.getMetadataCache.mockResolvedValue({ entries: [] });
		const actions = createSearchCacheActions(shell, vi.fn());

		await actions.saveIndexerSearchSettings({ retentionDays: 14 } as never);
		await actions.clearIndexerSearchCachePattern(' scenario ');
		await actions.clearMetadataCachePattern(' movie ');
		await actions.clearIndexerSearchCachePattern('   ');

		expect(apiMock.updateIndexerSearchSettings).toHaveBeenCalledWith({ retentionDays: 14 });
		expect(apiMock.clearIndexerSearchCacheByPattern).toHaveBeenCalledWith('scenario');
		expect(apiMock.clearMetadataCacheByPattern).toHaveBeenCalledWith('movie');
		expect(apiMock.clearIndexerSearchCacheByPattern).toHaveBeenCalledTimes(1);
		expect(shell.message).toBe('Metadata cache reset: 3 matching entries deleted');
		expect(shell.clearingMetadataCache).toBe(false);
	});
});

describe('integration test and activity actions (SCN-ACTIVITY-002)', () => {
	it('stores integration test results and reloads settings after indexer tests', async () => {
		const shell = state();
		const loadSettings = vi.fn();
		apiMock.testDownloadClientConfig.mockResolvedValue({ success: true });
		apiMock.testIndexer.mockResolvedValue({ success: true, message: 'ok' });
		apiMock.testMetadataProvider.mockResolvedValue({ success: false, message: 'missing token' });
		const actions = createSettingsTestCacheActions(shell, { clearNotice: vi.fn(), loadSettings });

		await expect(actions.testDownloadClientConfig({ name: 'Client' } as never)).resolves.toEqual({
			success: true
		});
		await actions.testIndexer('indexer-1');
		await actions.testMetadataProvider('metadata-1');

		expect(shell.indexerTests).toEqual({ 'indexer-1': { success: true, message: 'ok' } });
		expect(shell.metadataProviderTests).toEqual({
			'metadata-1': { success: false, message: 'missing token' }
		});
		expect(loadSettings).toHaveBeenCalledOnce();
		expect(shell.testingIndexerId).toBeUndefined();
	});

	it('cancels and deletes activity while keeping visible state consistent', async () => {
		const shell = state();
		const upsertActivity = vi.fn();
		const loadMediaItems = vi.fn();
		apiMock.cancelDownloadActivity.mockResolvedValue({ id: 'activity-1', status: 'cancelled' });
		const actions = createActivityActions(shell, {
			clearNotice: vi.fn(),
			loadMediaItems,
			upsertActivity
		});

		await actions.cancelActivity({ id: 'activity-1' } as never);
		await actions.deleteActivity({ id: 'activity-2' } as never);

		expect(upsertActivity).toHaveBeenCalledWith({ id: 'activity-1', status: 'cancelled' });
		expect(loadMediaItems).toHaveBeenCalledOnce();
		expect(apiMock.deleteDownloadActivity).toHaveBeenCalledWith('activity-2');
		expect(shell.activities).toEqual([{ id: 'activity-1' }]);
		expect(shell.message).toBe('Download activity deleted');
	});
});
