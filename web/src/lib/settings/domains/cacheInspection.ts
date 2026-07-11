import { client } from '$lib/api/client';
import type { IndexerSearchCacheEntry, IndexerSearchSettings, MetadataCacheEntry } from '../types';
import { emptyIndexerSearch, emptyMetadataCache } from './defaults';

export interface CacheInspectionLimits {
	cacheLimit?: number;
	historyLimit?: number;
}

export async function getMetadataCache(limits: CacheInspectionLimits = {}) {
	const { data, error } = await client.GET('/settings/metadata-cache', {
		params: { query: limits }
	});

	if (error) {
		throw new Error(error.message);
	}
	return data ?? emptyMetadataCache();
}

export async function getIndexerSearch(limits: CacheInspectionLimits = {}) {
	const { data, error } = await client.GET('/settings/indexer-search', {
		params: { query: limits }
	});

	if (error) {
		throw new Error(error.message);
	}
	return data ?? emptyIndexerSearch();
}

export async function updateIndexerSearchSettings(settings: IndexerSearchSettings) {
	const { data, error } = await client.PUT('/settings/indexer-search', { body: settings });

	if (error) {
		throw new Error(error.message);
	}
	return data ?? emptyIndexerSearch();
}

export async function clearIndexerSearchCache() {
	const { data, error } = await client.DELETE('/settings/indexer-search/cache');

	if (error) {
		throw new Error(error.message);
	}
	return data?.deletedCount ?? 0;
}

export async function clearIndexerSearchCacheByPattern(pattern: string) {
	const { data, error } = await client.POST('/settings/indexer-search/cache/reset', {
		body: { pattern }
	});

	if (error) {
		throw new Error(error.message);
	}
	return data?.deletedCount ?? 0;
}

export async function clearIndexerSearchHistory() {
	const { data, error } = await client.DELETE('/settings/indexer-search/history');

	if (error) {
		throw new Error(error.message);
	}
	return data?.deletedCount ?? 0;
}

export async function deleteIndexerSearchCacheEntry(entry: IndexerSearchCacheEntry) {
	const { data, error } = await client.DELETE('/settings/indexer-search/cache/entry', {
		params: {
			query: {
				indexerId: entry.indexerId,
				mediaType: entry.mediaType,
				query: entry.query
			}
		}
	});

	if (error) {
		throw new Error(error.message);
	}
	return data?.deletedCount ?? 0;
}

export async function clearMetadataCache() {
	const { data, error } = await client.DELETE('/settings/metadata-cache');

	if (error) {
		throw new Error(error.message);
	}
	return data?.deletedCount ?? 0;
}

export async function clearMetadataCacheByPattern(pattern: string) {
	const { data, error } = await client.POST('/settings/metadata-cache/reset', {
		body: { pattern }
	});

	if (error) {
		throw new Error(error.message);
	}
	return data?.deletedCount ?? 0;
}

export async function deleteMetadataCacheEntry(entry: MetadataCacheEntry) {
	const { data, error } = await client.DELETE('/settings/metadata-cache/entry', {
		params: {
			query: {
				providerId: entry.providerId,
				mediaType: entry.mediaType,
				query: entry.query,
				year: entry.year
			}
		}
	});

	if (error) {
		throw new Error(error.message);
	}
	return data?.deletedCount ?? 0;
}

export async function clearMetadataSearchHistory() {
	const { data, error } = await client.DELETE('/settings/metadata-cache/history');

	if (error) {
		throw new Error(error.message);
	}
	return data?.deletedCount ?? 0;
}
