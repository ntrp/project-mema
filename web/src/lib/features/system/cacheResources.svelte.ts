import { createQuery, type QueryClient } from '@tanstack/svelte-query';
import { getIndexerSearch, getMetadataCache } from '$lib/settings/domains/cacheInspection';
import { emptyIndexerSearch, emptyMetadataCache } from '$lib/settings/domains/defaults';
import { getProfile } from '$lib/profile/profileApi';
import type { IndexerSearchResponse, MetadataCacheResponse } from '$lib/settings/types';
import type { UserProfile } from '$lib/profile/types';

export const serverResourceKeys = {
	all: ['server-resources'] as const,
	indexerSearch: () => [...serverResourceKeys.all, 'indexer-search'] as const,
	metadataCache: () => [...serverResourceKeys.all, 'metadata-cache'] as const,
	profile: () => [...serverResourceKeys.all, 'profile'] as const
};

export function createServerResourceRuntime(
	client: QueryClient,
	enabled: { indexerSearch: () => boolean; metadataCache: () => boolean; profile: () => boolean }
) {
	const indexerSearch = createQuery(() => ({
		queryKey: serverResourceKeys.indexerSearch(),
		queryFn: () => getIndexerSearch(),
		enabled: enabled.indexerSearch()
	}));
	const metadataCache = createQuery(() => ({
		queryKey: serverResourceKeys.metadataCache(),
		queryFn: () => getMetadataCache(),
		enabled: enabled.metadataCache()
	}));
	const profile = createQuery(() => ({
		queryKey: serverResourceKeys.profile(),
		queryFn: getProfile,
		enabled: enabled.profile()
	}));
	async function refetchProfile() {
		const result = await profile.refetch();
		if (result.error) throw result.error;
		if (!result.data) throw new Error('Profile request did not return a result');
		return result.data;
	}
	return {
		indexerSearch,
		metadataCache,
		profile,
		refetchProfile,
		setIndexerSearch: (value: IndexerSearchResponse) =>
			client.setQueryData(serverResourceKeys.indexerSearch(), value),
		setMetadataCache: (value: MetadataCacheResponse) =>
			client.setQueryData(serverResourceKeys.metadataCache(), value),
		setProfile: (value: UserProfile | undefined) =>
			client.setQueryData(serverResourceKeys.profile(), value),
		clear: () => client.removeQueries({ queryKey: serverResourceKeys.all }),
		emptyIndexerSearch,
		emptyMetadataCache
	};
}
