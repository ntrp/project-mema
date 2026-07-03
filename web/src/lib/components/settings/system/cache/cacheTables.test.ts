import { describe, expect, it, vi } from 'vitest';

import IndexerSearchCacheTable from './IndexerSearchCacheTable.svelte';
import MetadataCacheTable from './MetadataCacheTable.svelte';
import { renderWithTooltip } from '$lib/components/rendered/renderHelpers';
import type { IndexerSearchResponse, MetadataCacheResponse } from '$lib/settings/types';

const now = '2026-07-03T00:00:00Z';
const future = '2026-07-04T00:00:00Z';

function metadataCacheResponse(entries: MetadataCacheResponse['entries']): MetadataCacheResponse {
	return {
		stats: {
			totalEntries: entries.length + 1,
			activeEntries: entries.filter((entry) => !entry.expired).length,
			expiredEntries: entries.filter((entry) => entry.expired).length,
			providerCount: 2
		},
		entries,
		historyEntries: [],
		historyTotalEntries: 0,
		historyStats: { totalEntries: 0, cacheHits: 0, cacheMisses: 0, failures: 0 }
	};
}

function indexerSearchResponse(
	cacheEntries: IndexerSearchResponse['cacheEntries']
): IndexerSearchResponse {
	return {
		settings: {
			cacheDurationMinutes: 60,
			historyRetentionDays: 7,
			automaticBlocklistExpiryDays: 7
		},
		stats: {
			totalEntries: cacheEntries.length + 1,
			activeEntries: cacheEntries.filter((entry) => !entry.expired).length,
			expiredEntries: cacheEntries.filter((entry) => entry.expired).length,
			indexerCount: 2
		},
		cacheEntries,
		historyEntries: [],
		historyTotalEntries: 0,
		historyStats: { totalEntries: 0, cacheHits: 0, cacheMisses: 0, failures: 0 }
	};
}

describe('cache tables (SCN-SETTINGS-015)', () => {
	it('renders metadata cache rows, expiry state, and load-more hint', () => {
		const { body } = renderWithTooltip(MetadataCacheTable, {
			cache: metadataCacheResponse([
				{
					providerId: 'provider-tmdb',
					providerName: 'TMDB Local',
					providerType: 'tmdb',
					mediaType: 'movie',
					query: 'edge of tomorrow',
					cacheKind: 'search',
					year: 2014,
					itemCount: 3,
					expiresAt: future,
					createdAt: now,
					updatedAt: now,
					expired: false
				},
				{
					providerId: 'provider-tvdb',
					providerName: 'TVDB Local',
					providerType: 'tvdb',
					mediaType: 'series',
					query: 'frieren',
					cacheKind: 'details',
					year: 2023,
					itemCount: 1,
					expiresAt: now,
					createdAt: now,
					updatedAt: now,
					expired: true
				}
			]),
			clearing: true,
			loading: false,
			onDeleteEntry: vi.fn(),
			onLoadMore: vi.fn()
		});

		expect(body).toContain('TMDB Local');
		expect(body).toContain('TVDB Local');
		expect(body).toContain('edge of tomorrow');
		expect(body).toContain('frieren');
		expect(body).toContain('movie · 2014');
		expect(body).toContain('Expired');
		expect(body).toContain('Scroll for more');
		expect(body).toContain('Delete cache entry');
	});

	it('renders indexer cache rows and empty states', () => {
		const populated = renderWithTooltip(IndexerSearchCacheTable, {
			search: indexerSearchResponse([
				{
					indexerId: 'indexer-1',
					indexerName: 'Torznab Local',
					indexerType: 'torznab',
					mediaType: 'movie',
					query: 'release snapshot 2160p',
					resultCount: 12,
					expiresAt: future,
					createdAt: now,
					updatedAt: now,
					expired: false
				}
			]),
			clearing: false,
			loading: true,
			onDeleteEntry: vi.fn(),
			onLoadMore: vi.fn()
		});

		expect(populated.body).toContain('Torznab Local');
		expect(populated.body).toContain('release snapshot 2160p');
		expect(populated.body).toContain('movie');
		expect(populated.body).toContain('12');
		expect(populated.body).toContain('Loading more...');

		const emptyMetadata = renderWithTooltip(MetadataCacheTable, {
			cache: metadataCacheResponse([]),
			clearing: false,
			loading: false,
			onDeleteEntry: vi.fn(),
			onLoadMore: vi.fn()
		});
		const emptyIndexer = renderWithTooltip(IndexerSearchCacheTable, {
			search: indexerSearchResponse([]),
			clearing: false,
			loading: false,
			onDeleteEntry: vi.fn(),
			onLoadMore: vi.fn()
		});

		expect(emptyMetadata.body).toContain('No metadata cache entries yet.');
		expect(emptyIndexer.body).toContain('No indexer cache entries yet.');
	});
});
