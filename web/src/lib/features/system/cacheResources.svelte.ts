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
	authenticated: () => boolean,
	admin: () => boolean
) {
	const indexerSearch = createQuery(() => ({
		queryKey: serverResourceKeys.indexerSearch(),
		queryFn: () => getIndexerSearch(),
		enabled: authenticated() && admin()
	}));
	const metadataCache = createQuery(() => ({
		queryKey: serverResourceKeys.metadataCache(),
		queryFn: () => getMetadataCache(),
		enabled: authenticated() && admin()
	}));
	const profile = createQuery(() => ({
		queryKey: serverResourceKeys.profile(),
		queryFn: getProfile,
		enabled: authenticated()
	}));
	return {
		indexerSearch,
		metadataCache,
		profile,
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
