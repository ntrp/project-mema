import type { IndexerSearchResponse, MetadataCacheResponse } from '../types';

export function emptyIndexerSearch(): IndexerSearchResponse {
	return {
		settings: {
			cacheDurationMinutes: 1440,
			historyRetentionDays: 7,
			automaticBlocklistExpiryDays: 7
		},
		stats: {
			totalEntries: 0,
			activeEntries: 0,
			expiredEntries: 0,
			indexerCount: 0
		},
		cacheEntries: [],
		historyEntries: [],
		historyTotalEntries: 0,
		historyStats: {
			totalEntries: 0,
			cacheHits: 0,
			cacheMisses: 0,
			failures: 0
		}
	};
}

export function emptyMetadataCache(): MetadataCacheResponse {
	return {
		stats: {
			totalEntries: 0,
			activeEntries: 0,
			expiredEntries: 0,
			providerCount: 0
		},
		entries: [],
		historyEntries: [],
		historyTotalEntries: 0,
		historyStats: {
			totalEntries: 0,
			cacheHits: 0,
			cacheMisses: 0,
			failures: 0
		}
	};
}
