import type { QueryClient } from '@tanstack/svelte-query';
import { createSettingsCatalogCache } from '$lib/features/settings/catalog/cache';
import {
	createCustomFormatsQuery,
	createDownloadClientsQuery,
	createIndexersQuery,
	createLanguagesQuery,
	createLibraryFoldersQuery,
	createMediaProfilesQuery,
	createMetadataProvidersQuery,
	createPathMappingsQuery,
	createSubtitleProvidersQuery,
	createTagsQuery,
	createUsersQuery
} from '$lib/features/settings/catalog/queries.svelte';
import { createLibraryScansRuntime } from '$lib/features/settings/libraryScans.svelte';
import { createServerResourceRuntime } from '$lib/features/system/cacheResources.svelte';
import type { AppShellState } from './state.svelte';

export function createControllerResourceRuntime(state: AppShellState, client: QueryClient) {
	const admin = () => state.authenticated && state.isAdmin;
	const catalogCache = createSettingsCatalogCache(client);
	const queries = {
		languages: createLanguagesQuery(() => state.authenticated),
		tags: createTagsQuery(() => state.authenticated),
		users: createUsersQuery(admin),
		downloadClients: createDownloadClientsQuery(admin),
		indexers: createIndexersQuery(admin),
		metadataProviders: createMetadataProvidersQuery(admin),
		subtitleProviders: createSubtitleProvidersQuery(admin),
		libraryFolders: createLibraryFoldersQuery(admin),
		pathMappings: createPathMappingsQuery(admin),
		mediaProfiles: createMediaProfilesQuery(admin),
		customFormats: createCustomFormatsQuery(admin)
	};
	const server = createServerResourceRuntime(
		client,
		() => state.authenticated,
		() => state.isAdmin
	);
	const scans = createLibraryScansRuntime(client);
	Object.defineProperties(state, {
		indexerSearch: {
			get: () => server.indexerSearch.data ?? server.emptyIndexerSearch(),
			set: server.setIndexerSearch
		},
		metadataCache: {
			get: () => server.metadataCache.data ?? server.emptyMetadataCache(),
			set: server.setMetadataCache
		},
		profile: { get: () => server.profile.data, set: server.setProfile },
		loadingIndexerSearch: { get: () => server.indexerSearch.isFetching },
		loadingMetadataCache: { get: () => server.metadataCache.isFetching },
		loadingProfile: { get: () => server.profile.isFetching }
	});
	const properties: Record<string, { get: () => unknown }> = Object.fromEntries(
		Object.entries(queries).map(([name, query]) => [name, { get: () => query.data ?? [] }])
	);
	properties.libraryScansByFolder = { get: () => scans.scans.data ?? {} };
	return { catalogCache, queries, server, scans, properties };
}
