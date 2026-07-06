import { beforeEach, describe, expect, it, vi } from 'vitest';

const clientMock = vi.hoisted(() => ({
	GET: vi.fn(),
	POST: vi.fn(),
	PUT: vi.fn(),
	DELETE: vi.fn()
}));

vi.mock('$lib/api/client', () => ({ client: clientMock }));

import {
	abortSystemJob,
	autocompleteMedia,
	clearIndexerSearchCache,
	clearIndexerSearchCacheByPattern,
	currentSessionAuthenticated,
	deleteDownloadClient,
	deleteMediaItem,
	downloadSystemLogFile,
	emptyIndexerSearch,
	emptyMetadataCache,
	getMediaCollection,
	getMediaMetadataDetails,
	listCustomFormats,
	listMediaItems,
	listSystemJobs,
	loadSettings,
	login,
	mediaTypeForLibraryKind,
	saveDownloadClient,
	saveIndexer,
	saveLanguage,
	saveTag,
	searchMediaReleases,
	testDownloadClientConfig,
	testIndexerConfig,
	updateQualitySizeSettings
} from './api';

describe('UI API helpers', () => {
	beforeEach(() => {
		clientMock.GET.mockReset();
		clientMock.POST.mockReset();
		clientMock.PUT.mockReset();
		clientMock.DELETE.mockReset();
	});

	it('SCN-SETTINGS-009 maps session, list, and default responses', async () => {
		clientMock.GET.mockResolvedValueOnce({ data: { authenticated: true } })
			.mockResolvedValueOnce({ data: { jobs: [{ id: 1 }] } })
			.mockResolvedValueOnce({ data: undefined })
			.mockResolvedValueOnce({ data: { releases: [{ id: 'r1' }], errors: ['bad'] } })
			.mockResolvedValueOnce({ data: { groups: [{ title: 'Matrix' }] } });

		await expect(currentSessionAuthenticated()).resolves.toBe(true);
		await expect(listSystemJobs({ status: ['running'], limit: 10 })).resolves.toEqual([{ id: 1 }]);
		await expect(listMediaItems()).resolves.toEqual([]);
		await expect(searchMediaReleases('media-1')).resolves.toEqual({
			releases: [{ id: 'r1' }],
			errors: ['bad']
		});
		await expect(autocompleteMedia('matrix', 'providers')).resolves.toEqual([{ title: 'Matrix' }]);
		expect(clientMock.GET).toHaveBeenLastCalledWith('/media/autocomplete', {
			params: { query: { query: 'matrix', includeLibrary: false, includeProviders: true } }
		});
		expect(emptyIndexerSearch().settings.cacheDurationMinutes).toBe(1440);
		expect(emptyIndexerSearch().settings.automaticBlocklistExpiryDays).toBe(7);
		expect(emptyMetadataCache().stats.providerCount).toBe(0);
		expect(mediaTypeForLibraryKind('anime_series')).toBe('serie');
		expect(mediaTypeForLibraryKind('movie')).toBe('movie');
	});

	it('SCN-SETTINGS-009 throws user-facing errors for failed or empty responses', async () => {
		clientMock.POST.mockResolvedValueOnce({
			error: { message: 'Invalid username or password' }
		}).mockResolvedValueOnce({ data: { authenticated: false } });
		clientMock.GET.mockResolvedValueOnce({
			error: { message: 'No collection' }
		}).mockResolvedValueOnce({ data: undefined });

		await expect(login('admin', 'bad')).rejects.toThrow('Invalid username or password');
		await expect(login('admin', 'admin')).rejects.toThrow('Login failed');
		await expect(getMediaCollection('tmdb', '1')).rejects.toThrow('No collection');
		await expect(getMediaMetadataDetails('tmdb', 'movie', '1')).rejects.toThrow(
			'Media details were not returned'
		);
	});

	it('SCN-SETTINGS-009 sends normalized settings writes', async () => {
		clientMock.POST.mockResolvedValue({ data: { success: true } });
		clientMock.PUT.mockResolvedValue({ data: { id: 'updated' } });

		await saveDownloadClient({
			name: ' SAB ',
			type: 'sabnzbd',
			protocol: 'usenet',
			baseUrl: ' http://sab.local ',
			username: '',
			password: ' secret ',
			apiKey: '',
			category: ' movies ',
			enabled: true,
			priority: 7
		});
		expect(clientMock.POST).toHaveBeenCalledWith('/settings/download-clients', {
			body: expect.objectContaining({ name: 'SAB', password: 'secret', category: 'movies' })
		});

		await saveIndexer({
			id: 'indexer-1',
			name: ' Torznab ',
			definitionId: 'generic-torznab',
			baseUrl: ' http://indexer.local ',
			apiKey: ' key ',
			categoriesText: '2000, bad, 2040',
			enabled: false,
			priority: 20
		});
		expect(clientMock.PUT).toHaveBeenCalledWith('/settings/indexers/{id}', {
			params: { path: { id: 'indexer-1' } },
			body: expect.objectContaining({ name: 'Torznab', categories: [2000, 2040] })
		});

		await saveTag({ id: 'tag-1', name: '  Action  ' });
		expect(clientMock.PUT).toHaveBeenLastCalledWith('/settings/tags/{id}', {
			params: { path: { id: 'tag-1' } },
			body: { name: 'Action' }
		});

		await saveLanguage({
			originalCode: 'DE',
			code: 'de',
			displayName: ' German ',
			aliasesText: 'deu, ger'
		});
		expect(clientMock.PUT).toHaveBeenLastCalledWith('/settings/languages/{code}', {
			params: { path: { code: 'DE' } },
			body: { displayName: 'German', aliases: ['deu', 'ger'] }
		});
	});

	it('SCN-SETTINGS-009 covers commands, deletes, fetch downloads, and fallback counts', async () => {
		clientMock.POST.mockResolvedValueOnce({ data: { id: 2 } })
			.mockResolvedValueOnce({ data: { success: false } })
			.mockResolvedValueOnce({ data: { success: true } })
			.mockResolvedValueOnce({ data: { deletedCount: 4 } });
		clientMock.PUT.mockResolvedValueOnce({ data: { qualities: [] } });
		clientMock.DELETE.mockResolvedValueOnce({})
			.mockResolvedValueOnce({})
			.mockResolvedValueOnce({ data: undefined })
			.mockResolvedValueOnce({ data: { deletedCount: 4 } });
		vi.stubGlobal(
			'fetch',
			vi.fn().mockResolvedValue({ ok: true, blob: () => Promise.resolve('blob-data') })
		);

		await expect(abortSystemJob(2)).resolves.toEqual({ id: 2 });
		await expect(
			testDownloadClientConfig({
				name: 'Client',
				type: 'transmission',
				protocol: 'torrent',
				baseUrl: 'http://client.local',
				username: '',
				password: '',
				apiKey: '',
				category: '',
				enabled: true,
				priority: 1
			})
		).resolves.toEqual({ success: false });
		await expect(
			testIndexerConfig({
				name: 'Indexer',
				definitionId: 'generic-torznab',
				baseUrl: 'http://indexer.local',
				apiKey: '',
				categoriesText: '2000',
				fields: [],
				redirect: true,
				enabled: true,
				priority: 1
			})
		).resolves.toEqual({ success: true });
		await expect(updateQualitySizeSettings([])).resolves.toEqual({ qualities: [] });
		await expect(deleteMediaItem('media-1', { keepFiles: true })).resolves.toBeUndefined();
		await expect(deleteDownloadClient('client-1')).resolves.toBeUndefined();
		await expect(clearIndexerSearchCache()).resolves.toBe(0);
		await expect(clearIndexerSearchCacheByPattern('matrix')).resolves.toBe(4);
		await expect(downloadSystemLogFile('app.log')).resolves.toBe('blob-data');
		expect(globalThis.fetch).toHaveBeenCalledWith('/api/system/log-files/app.log/download', {
			credentials: 'include'
		});
		vi.unstubAllGlobals();
	});

	it('SCN-SETTINGS-009 aggregates settings data from the individual endpoints', async () => {
		clientMock.GET.mockResolvedValue({ data: {} });
		const settings = await loadSettings();
		expect(settings.downloadClients).toEqual([]);
		expect(settings.indexers).toEqual([]);
		expect(settings.indexerSearch).toEqual({});
		expect(settings.metadataCache).toEqual({});
		expect(settings.languages).toEqual([]);
		expect(clientMock.GET).toHaveBeenCalledTimes(13);
		await expect(listCustomFormats()).resolves.toEqual([]);
	});
});
