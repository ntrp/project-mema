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
	addDiscoverBlacklistItem,
	approveMediaRequest,
	cancelDownloadActivity,
	createLibraryFolderOption,
	createMediaItem,
	createMediaRequest,
	deleteCustomFormat,
	deleteDownloadActivity,
	deleteIndexer,
	deleteLanguage,
	deleteLibraryFolder,
	deleteMediaProfile,
	deleteMetadataProvider,
	deletePathMapping,
	deleteTag,
	deleteUser,
	enqueueMediaAutomaticSearch,
	enqueueMediaReleaseSearch,
	getFileNamingSettings,
	getIndexerSearch,
	getLibraryScan,
	getMediaCollection,
	getMediaRequest,
	getSystemEventSettings,
	getSystemLogFileSettings,
	getSystemLogLevel,
	getSystemStatus,
	listLibraryFolderOptions,
	listQualitySizeSettings,
	manualImportDownloadActivity,
	matchLibraryScanItem,
	saveCustomFormat,
	saveLibraryFolder,
	saveMediaProfile,
	saveMetadataProvider,
	savePathMapping,
	saveUser,
	scanLibraryFolder,
	testDownloadClient,
	testIndexer,
	testMetadataProvider,
	updateFileNamingSettings,
	updateIndexerSearchSettings,
	updateMediaItem,
	updateQualitySizeSettings,
	updateSystemEventSettings,
	updateSystemLogFileSettings,
	updateSystemLogLevel
} from '../api';

describe('UI API edge cases (SCN-SETTINGS-009)', () => {
	beforeEach(() => {
		clientMock.GET.mockReset().mockResolvedValue({ data: undefined });
		clientMock.POST.mockReset().mockResolvedValue({ data: undefined });
		clientMock.PUT.mockReset().mockResolvedValue({ data: undefined });
		clientMock.DELETE.mockReset().mockResolvedValue({});
	});

	it('throws clear messages when command endpoints omit required data', async () => {
		const cases: [() => Promise<unknown>, string][] = [
			[() => getSystemLogLevel(), 'Log level request did not return a result'],
			[() => updateSystemLogLevel('info' as never), 'Log level update did not return a result'],
			[() => getSystemStatus(), 'System status request did not return a result'],
			[() => getSystemLogFileSettings(), 'Log file settings request did not return a result'],
			[
				() => updateSystemLogFileSettings({ enabled: true } as never),
				'Log file settings update did not return a result'
			],
			[() => getSystemEventSettings(), 'Event settings request did not return a result'],
			[
				() => updateSystemEventSettings({ enabled: true } as never),
				'Event settings update did not return a result'
			],
			[() => abortSystemJob(1), 'Job abort did not return a job'],
			[() => listQualitySizeSettings(), 'Quality size settings were not returned'],
			[() => updateQualitySizeSettings([]), 'Quality size settings were not returned'],
			[() => getFileNamingSettings(), 'File naming settings were not returned'],
			[
				() => updateFileNamingSettings({ movieFileFormat: '{title}' } as never),
				'File naming settings were not returned'
			],
			[
				() => addDiscoverBlacklistItem({ title: 'Hidden' } as never),
				'Blacklist request did not return a result'
			],
			[() => getMediaCollection('tmdb', '1'), 'Media collection was not returned'],
			[() => createMediaItem({ title: 'Movie' } as never), 'Media item was not returned'],
			[
				() => updateMediaItem('media-1', { monitored: true } as never),
				'Media item was not returned'
			],
			[() => createMediaRequest({ title: 'Movie' } as never), 'Media request was not returned'],
			[() => getMediaRequest('request-1'), 'Media request was not returned'],
			[
				() => approveMediaRequest('request-1', { monitored: true } as never),
				'Media request approval was not returned'
			],
			[() => enqueueMediaReleaseSearch('media-1'), 'Release search job was not returned'],
			[() => enqueueMediaAutomaticSearch('media-1'), 'Automatic search job was not returned'],
			[() => cancelDownloadActivity('activity-1'), 'Download activity was not returned'],
			[
				() => manualImportDownloadActivity('activity-1', { sourcePath: '/movie.mkv' } as never),
				'Download activity was not returned'
			],
			[() => testDownloadClient('client-1'), 'Download client test did not return a result'],
			[() => testIndexer('indexer-1'), 'Indexer test did not return a result'],
			[() => testMetadataProvider('provider-1'), 'Metadata provider test did not return a result'],
			[() => listLibraryFolderOptions('/media'), 'Folder options were not returned'],
			[() => createLibraryFolderOption('/media', 'Movies'), 'Folder was not returned'],
			[() => scanLibraryFolder('folder-1'), 'Library scan was not returned'],
			[
				() => savePathMapping({ clientPath: '/a', appPath: '/b' } as never),
				'Path mapping was not returned'
			],
			[() => getLibraryScan('scan-1'), 'Library scan was not returned'],
			[
				() => matchLibraryScanItem('scan-1', 'item-1', {} as never),
				'Library match was not returned'
			]
		];

		for (const [call, message] of cases) {
			await expect(call()).rejects.toThrow(message);
		}
	});

	it('returns default cache shapes and sends less common save/delete commands', async () => {
		clientMock.POST.mockResolvedValue({ data: { id: 'created' } });
		clientMock.PUT.mockResolvedValue({ data: { cacheEntries: [] } });

		await expect(getIndexerSearch()).resolves.toMatchObject({
			settings: { cacheDurationMinutes: 1440 }
		});
		await expect(
			updateIndexerSearchSettings({ cacheDurationMinutes: 5 } as never)
		).resolves.toEqual({
			cacheEntries: []
		});
		await expect(
			saveLibraryFolder({ path: ' /media ', mediaKind: 'movie' } as never)
		).resolves.toEqual({
			id: 'created'
		});
		await expect(
			savePathMapping({ clientPath: ' /downloads ', appPath: ' /media ' } as never)
		).resolves.toEqual({
			id: 'created'
		});

		await saveMetadataProvider(providerForm());
		await saveMetadataProvider({
			...providerForm(),
			id: 'provider-1'
		} as never);
		await saveUser({ username: 'admin', password: 'secret', roles: ['admin'] } as never);
		await saveUser({ id: 'user-1', username: 'admin', roles: ['admin'] } as never);
		await saveMediaProfile(mediaProfileForm() as never);
		await saveMediaProfile({ ...mediaProfileForm(), id: 'profile-1' } as never);
		await saveCustomFormat({ name: 'Format', includeSpecs: [], excludeSpecs: [] } as never);
		await saveCustomFormat({
			id: 'format-1',
			name: 'Format',
			includeSpecs: [],
			excludeSpecs: []
		} as never);

		const deletes = [
			deleteDownloadActivity,
			deleteIndexer,
			deleteMetadataProvider,
			deleteUser,
			deleteTag
		];
		deletes.push(
			deleteLanguage,
			deleteMediaProfile,
			deleteCustomFormat,
			deleteLibraryFolder,
			deletePathMapping
		);
		for (const [index, remove] of deletes.entries()) {
			await remove(`id-${index}`);
		}
		expect(clientMock.DELETE).toHaveBeenCalledTimes(10);
	});
});

function providerForm() {
	return {
		name: 'TMDb',
		type: 'tmdb',
		baseUrl: ' https://api.themoviedb.org/3 ',
		apiKey: ' key ',
		pin: '',
		accessToken: '',
		enabled: true,
		priority: 10
	} as const;
}

function mediaProfileForm() {
	return {
		name: 'Profile',
		qualityIds: [],
		upgradesAllowed: true,
		minimumCustomFormatScore: 0,
		upgradeUntilCustomFormatScore: 0,
		minimumCustomFormatScoreIncrement: 1,
		removeNonEnabledLanguages: false,
		preferredProtocol: 'any',
		seriesPackPreference: 'auto',
		targetLanguages: ['english'],
		targetLanguageScores: [{ languageId: 'english', score: 0, required: false }],
		customFormatScores: []
	} as const;
}
