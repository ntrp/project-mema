import { beforeEach, describe, expect, it, vi } from 'vitest';

const clientMock = vi.hoisted(() => ({
	GET: vi.fn(),
	POST: vi.fn(),
	PUT: vi.fn(),
	DELETE: vi.fn()
}));

vi.mock('$lib/api/client', () => ({ client: clientMock }));

import {
	clearMetadataCache,
	clearMetadataCacheByPattern,
	clearMetadataSearchHistory,
	createLibraryFolderOption,
	deleteIndexerSearchCacheEntry,
	deleteMetadataCacheEntry,
	getFileNamingSettings,
	getLibraryScan,
	getMetadataCache,
	listLibraryFolderOptions,
	matchLibraryScanItem,
	scanLibraryFolder,
	testCustomFormatParsing,
	testIndexer,
	testMetadataProvider,
	updateFileNamingSettings,
	updateIndexerSearchSettings
} from '../api';

describe('cache and library UI API helpers (SCN-SETTINGS-009)', () => {
	beforeEach(() => {
		clientMock.GET.mockReset();
		clientMock.POST.mockReset();
		clientMock.PUT.mockReset();
		clientMock.DELETE.mockReset();
	});

	it('maps parsing, integration tests, cache operations, and library helpers', async () => {
		clientMock.GET.mockResolvedValueOnce({ data: { movieFolderFormat: '{movie_title}' } })
			.mockResolvedValueOnce({ data: { entries: [] } })
			.mockResolvedValueOnce({ data: { options: [{ path: '/media' }] } })
			.mockResolvedValueOnce({ data: { id: 'scan-1' } });
		clientMock.PUT.mockResolvedValueOnce({
			data: { movieFolderFormat: '{title}' }
		}).mockResolvedValueOnce({ data: { id: 'result-1', deletedCount: 3 } });
		clientMock.POST.mockResolvedValue({ data: { id: 'result-1', deletedCount: 3 } });
		clientMock.DELETE.mockResolvedValue({ data: { deletedCount: 2 } });

		await expect(getFileNamingSettings()).resolves.toEqual({ movieFolderFormat: '{movie_title}' });
		await expect(
			updateFileNamingSettings({ movieFolderFormat: '{title}' } as never)
		).resolves.toEqual({
			movieFolderFormat: '{title}'
		});
		await expect(testCustomFormatParsing('Movie.2026.mkv')).resolves.toEqual({
			id: 'result-1',
			deletedCount: 3
		});
		await expect(testIndexer('indexer-1')).resolves.toEqual({ id: 'result-1', deletedCount: 3 });
		await expect(testMetadataProvider('provider-1')).resolves.toEqual({
			id: 'result-1',
			deletedCount: 3
		});
		await expect(
			updateIndexerSearchSettings({ cacheDurationMinutes: 1 } as never)
		).resolves.toEqual({
			id: 'result-1',
			deletedCount: 3
		});
		await expect(clearMetadataCache()).resolves.toBe(2);
		await expect(clearMetadataCacheByPattern('matrix')).resolves.toBe(3);
		await expect(clearMetadataSearchHistory()).resolves.toBe(2);
		await expect(
			deleteIndexerSearchCacheEntry({ indexerId: 'i', mediaType: 'movie', query: 'q' } as never)
		).resolves.toBe(2);
		await expect(
			deleteMetadataCacheEntry({
				providerId: 'p',
				mediaType: 'movie',
				query: 'q',
				year: 2026
			} as never)
		).resolves.toBe(2);
		await expect(getMetadataCache({ cacheLimit: 2 })).resolves.toEqual({ entries: [] });
		await expect(listLibraryFolderOptions('/media')).resolves.toEqual({
			options: [{ path: '/media' }]
		});
		await expect(createLibraryFolderOption('/media', 'Movies')).resolves.toEqual({
			id: 'result-1',
			deletedCount: 3
		});
		await expect(scanLibraryFolder('folder-1')).resolves.toEqual({
			id: 'result-1',
			deletedCount: 3
		});
		await expect(getLibraryScan('scan-1')).resolves.toEqual({ id: 'scan-1' });
		await expect(
			matchLibraryScanItem('scan-1', 'item-1', { mediaItemId: 'media-1' } as never)
		).resolves.toEqual({
			id: 'result-1',
			deletedCount: 3
		});
	});
});
