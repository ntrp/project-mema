import { describe, expect, it } from 'vitest';
import type { AppShellState } from './state.svelte';
import { createEventActions } from './events';
import { createNoticeActions } from './noticeActions';
import type {
	DownloadActivity,
	IndexerSearchCacheEntry,
	IndexerSearchHistoryEntry,
	MediaItem,
	MetadataCacheEntry,
	MetadataSearchHistoryEntry
} from '$lib/settings/types';
describe('app shell event actions (SCN-SYSTEM-008)', () => {
	it('upserts activity events and maps active downloads to media item status', () => {
		const state = testState();
		state.activities = [activity({ id: 'activity-1', status: 'queued' })];
		state.mediaItems = [{ id: 'media-1', title: 'Scenario Movie', status: 'missing' } as MediaItem];
		const actions = createEventActions(state);
		const update = activity({ id: 'activity-1', status: 'downloading' });
		actions.upsertActivity(update);
		actions.updateMediaStatusFromActivity(update);

		expect(state.activities).toHaveLength(1);
		expect(state.activities[0].status).toBe('downloading');
		expect(state.mediaItems[0].status).toBe('downloading');
		actions.updateMediaStatusFromActivity(activity({ status: 'completed' }));
		expect(state.mediaItems[0].status).toBe('downloaded');
	});

	it('appends search history and cache updates with bounded dedupe', () => {
		const state = testState();
		const actions = createEventActions(state);
		actions.appendIndexerSearchHistory(indexerHistory({ cacheHit: true, success: false }));
		actions.upsertIndexerSearchCache({
			entry: indexerCache({ query: 'scenario' }),
			stats: { totalEntries: 1, activeEntries: 1, expiredEntries: 0, indexerCount: 1 }
		});
		actions.upsertIndexerSearchCache({
			entry: indexerCache({ query: 'scenario', resultCount: 2 }),
			stats: { totalEntries: 1, activeEntries: 1, expiredEntries: 0, indexerCount: 1 }
		});

		expect(state.indexerSearch.historyTotalEntries).toBe(1);
		expect(state.indexerSearch.historyStats).toMatchObject({
			totalEntries: 1,
			cacheHits: 1,
			failures: 1
		});
		expect(state.indexerSearch.cacheEntries).toHaveLength(1);
		expect(state.indexerSearch.cacheEntries[0].resultCount).toBe(2);
	});

	it('appends metadata cache events and ignores malformed event payloads', () => {
		const state = testState();
		const actions = createEventActions(state);

		actions.appendMetadataSearchHistory(metadataHistory({ cacheHit: false, success: true }));
		actions.upsertMetadataCache({
			entry: metadataCache({ query: 'scenario' }),
			stats: { totalEntries: 1, activeEntries: 1, expiredEntries: 0, providerCount: 1 }
		});
		expect(state.metadataCache.historyTotalEntries).toBe(1);
		expect(state.metadataCache.historyStats).toMatchObject({
			totalEntries: 1,
			cacheMisses: 1,
			failures: 0
		});
		expect(state.metadataCache.entries[0].query).toBe('scenario');
		expect(actions.parseEventData(new MessageEvent('message', { data: '{' }))).toBeUndefined();
		expect(
			actions.parseEventData<{ ok: boolean }>(
				new MessageEvent('message', { data: '{"data":{"ok":true}}' })
			)
		).toEqual({ ok: true });
	});
});

describe('app shell notice actions (SCN-MEDIA-004)', () => {
	it('clears current messages before showing profile guidance', () => {
		const state = testState();
		const actions = createNoticeActions(state);

		state.errorMessage = 'Previous error';
		state.message = 'Previous message';
		actions.clearNotice();
		expect(state.errorMessage).toBe('');
		expect(state.message).toBe('');

		state.errorMessage = 'Previous error';
		actions.showProfile();
		expect(state.errorMessage).toBe('');
		expect(state.message).toBe('Profile settings are not implemented yet');
	});
});

function testState(): AppShellState {
	return {
		message: '',
		errorMessage: '',
		activities: [],
		mediaItems: [],
		indexerSearch: {
			stats: { totalEntries: 0, activeEntries: 0, expiredEntries: 0, indexerCount: 0 },
			cacheEntries: [],
			historyStats: { totalEntries: 0, cacheHits: 0, cacheMisses: 0, failures: 0 },
			historyEntries: [],
			historyTotalEntries: 0
		},
		metadataCache: {
			stats: { totalEntries: 0, activeEntries: 0, expiredEntries: 0, providerCount: 0 },
			entries: [],
			historyStats: { totalEntries: 0, cacheHits: 0, cacheMisses: 0, failures: 0 },
			historyEntries: [],
			historyTotalEntries: 0
		}
	} as unknown as AppShellState;
}

function activity(overrides: Partial<DownloadActivity> = {}): DownloadActivity {
	return {
		id: 'activity-1',
		mediaItemId: 'media-1',
		status: 'queued',
		title: 'Scenario Movie',
		...overrides
	} as DownloadActivity;
}

function indexerHistory(
	overrides: Partial<IndexerSearchHistoryEntry> = {}
): IndexerSearchHistoryEntry {
	return {
		id: 'history-1',
		indexerId: 'indexer-1',
		indexerName: 'Scenario Indexer',
		indexerProtocol: 'torrent',
		mediaType: 'movie',
		query: 'scenario',
		cacheHit: false,
		success: true,
		resultCount: 1,
		createdAt: '2026-07-03T00:00:00Z',
		...overrides
	} as IndexerSearchHistoryEntry;
}

function indexerCache(overrides: Partial<IndexerSearchCacheEntry> = {}): IndexerSearchCacheEntry {
	return {
		indexerId: 'indexer-1',
		indexerName: 'Scenario Indexer',
		indexerProtocol: 'torrent',
		mediaType: 'movie',
		query: 'scenario',
		resultCount: 1,
		expired: false,
		createdAt: '2026-07-03T00:00:00Z',
		expiresAt: '2026-07-04T00:00:00Z',
		...overrides
	} as IndexerSearchCacheEntry;
}

function metadataHistory(
	overrides: Partial<MetadataSearchHistoryEntry> = {}
): MetadataSearchHistoryEntry {
	return {
		id: 'metadata-history-1',
		providerId: 'provider-1',
		providerName: 'Scenario Provider',
		providerType: 'tmdb',
		mediaType: 'movie',
		query: 'scenario',
		cacheKind: 'search',
		cacheHit: false,
		success: true,
		itemCount: 1,
		createdAt: '2026-07-03T00:00:00Z',
		...overrides
	} as MetadataSearchHistoryEntry;
}

function metadataCache(overrides: Partial<MetadataCacheEntry> = {}): MetadataCacheEntry {
	return {
		providerId: 'provider-1',
		providerName: 'Scenario Provider',
		providerType: 'tmdb',
		mediaType: 'movie',
		query: 'scenario',
		year: 2026,
		cacheKind: 'search',
		itemCount: 1,
		expired: false,
		createdAt: '2026-07-03T00:00:00Z',
		expiresAt: '2026-07-04T00:00:00Z',
		...overrides
	} as MetadataCacheEntry;
}
