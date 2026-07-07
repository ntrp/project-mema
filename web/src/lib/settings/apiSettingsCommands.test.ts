import { beforeEach, describe, expect, it, vi } from 'vitest';

const clientMock = vi.hoisted(() => ({
	GET: vi.fn(),
	POST: vi.fn(),
	PUT: vi.fn(),
	DELETE: vi.fn()
}));

vi.mock('$lib/api/client', () => ({ client: clientMock }));

import {
	cancelDownloadActivity,
	clearMetadataCache,
	clearMetadataCacheByPattern,
	clearMetadataSearchHistory,
	createLibraryFolderOption,
	deleteCustomFormat,
	deleteDownloadActivity,
	deleteIndexer,
	deleteIndexerSearchCacheEntry,
	deleteLanguage,
	deleteLibraryFolder,
	deleteMediaProfile,
	deleteMetadataCacheEntry,
	deleteMetadataProvider,
	deletePathMapping,
	deleteTag,
	deleteUser,
	getIndexerSearch,
	getLibraryScan,
	getMetadataCache,
	listLibraryFolderOptions,
	manualImportDownloadActivity,
	matchLibraryScanItem,
	saveCustomFormat,
	saveLibraryFolder,
	saveMediaProfile,
	saveMetadataProvider,
	savePathMapping,
	saveUser,
	scanLibraryFolder,
	testIndexerConfig,
	testIndexer,
	testMetadataProvider,
	updateIndexerSearchSettings
} from './api';
import type { IndexerSearchCacheEntry, MetadataCacheEntry } from './types';

describe('UI API settings command helpers (SCN-SETTINGS-009)', () => {
	beforeEach(() => {
		clientMock.GET.mockReset().mockResolvedValue({ data: {} });
		clientMock.POST.mockReset().mockResolvedValue({ data: {} });
		clientMock.PUT.mockReset().mockResolvedValue({ data: {} });
		clientMock.DELETE.mockReset().mockResolvedValue({ data: {} });
	});

	it('maps activity, integration, cache, and library commands to API calls', async () => {
		await expect(cancelDownloadActivity('activity-1')).resolves.toEqual({});
		await expect(deleteDownloadActivity('activity-1')).resolves.toBeUndefined();
		await expect(manualImportDownloadActivity('activity-1', {} as never)).resolves.toEqual({});
		await expect(testIndexerConfig(indexerForm())).resolves.toEqual({});
		await expect(testIndexer('indexer-1')).resolves.toEqual({});
		await expect(saveMetadataProvider(metadataProviderForm())).resolves.toBeUndefined();
		await expect(saveUser(userForm())).resolves.toBeUndefined();
		await expect(saveMediaProfile(mediaProfileForm() as never)).resolves.toBeUndefined();
		await expect(saveCustomFormat(customFormatForm())).resolves.toBeUndefined();
		await expect(testMetadataProvider('provider-1')).resolves.toEqual({});
		await expect(getMetadataCache({ cacheLimit: 1 })).resolves.toEqual({});
		await expect(getIndexerSearch({ historyLimit: 1 })).resolves.toEqual({});
		await expect(
			updateIndexerSearchSettings({
				cacheDurationMinutes: 1,
				historyRetentionDays: 7,
				automaticBlocklistExpiryDays: 7
			})
		).resolves.toEqual({});
		await expect(clearMetadataCache()).resolves.toBe(0);
		await expect(clearMetadataCacheByPattern('matrix')).resolves.toBe(0);
		await expect(clearMetadataSearchHistory()).resolves.toBe(0);
		await expect(deleteIndexerSearchCacheEntry(cacheEntry())).resolves.toBe(0);
		await expect(deleteMetadataCacheEntry(metadataEntry())).resolves.toBe(0);
		await expect(deleteIndexer('indexer-1')).resolves.toBeUndefined();
		await expect(deleteMetadataProvider('provider-1')).resolves.toBeUndefined();
		await expect(deleteUser('user-1')).resolves.toBeUndefined();
		await expect(deleteTag('tag-1')).resolves.toBeUndefined();
		await expect(deleteLanguage('en')).resolves.toBeUndefined();
		await expect(deleteMediaProfile('profile-1')).resolves.toBeUndefined();
		await expect(deleteCustomFormat('format-1')).resolves.toBeUndefined();
		await expect(saveLibraryFolder({ path: '/media', kind: 'movie' })).resolves.toEqual({});
		await expect(listLibraryFolderOptions('/media')).resolves.toEqual({});
		await expect(createLibraryFolderOption('/media', 'Movies')).resolves.toEqual({});
		await expect(deleteLibraryFolder('folder-1')).resolves.toBeUndefined();
		await expect(scanLibraryFolder('folder-1')).resolves.toEqual({});
		await expect(savePathMapping({ clientPath: '/downloads', appPath: '/media' })).resolves.toEqual(
			{}
		);
		await expect(deletePathMapping('mapping-1')).resolves.toBeUndefined();
		await expect(getLibraryScan('scan-1')).resolves.toEqual({});
		await expect(matchLibraryScanItem('scan-1', 'item-1', {} as never)).resolves.toEqual({});
	});
});

function cacheEntry(): IndexerSearchCacheEntry {
	return {
		indexerId: 'indexer-1',
		mediaType: 'movie',
		query: 'scenario'
	} as IndexerSearchCacheEntry;
}

function indexerForm() {
	return {
		name: 'Indexer',
		definitionId: 'generic-torznab',
		baseUrl: 'http://indexer.local',
		apiKey: '',
		categoriesText: '2000',
		fields: [],
		redirect: true,
		enabled: true,
		priority: 1
	};
}

function metadataEntry(): MetadataCacheEntry {
	return {
		providerId: 'provider-1',
		mediaType: 'movie',
		query: 'scenario',
		year: 2026
	} as MetadataCacheEntry;
}

function metadataProviderForm() {
	return {
		name: 'TMDB',
		type: 'tmdb',
		baseUrl: 'https://metadata.local',
		apiKey: '',
		pin: '',
		accessToken: '',
		enabled: true,
		priority: 1
	} as const;
}

function userForm() {
	return { username: 'viewer', password: 'secret', role: 'user' } as const;
}

function customFormatForm() {
	return { name: 'HQ', includeInRenameTemplate: true, includeSpecs: [], excludeSpecs: [] };
}

function mediaProfileForm() {
	return {
		name: 'HD',
		isDefault: false,
		qualityIds: [],
		upgradesAllowed: true,
		minimumCustomFormatScore: 0,
		upgradeUntilCustomFormatScore: 0,
		minimumCustomFormatScoreIncrement: 1,
		removeUnwantedAudio: false,
		audioLossyTranscodePolicy: 'disabled',
		removeUnwantedSubtitles: false,
		subtitleMode: 'mixed',
		allowSubtitleReleaseFallback: false,
		preferredProtocol: 'any',
		seriesPackPreference: 'auto',
		audioTargets: [
			{
				languageId: 'english',
				score: 0
			}
		],
		subtitleTargets: [{ languageId: 'english', score: 0 }],
		customFormatScores: []
	} as const;
}
