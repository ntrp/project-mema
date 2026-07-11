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
import { createResourceEnablement } from './resourceEnablement';

export function createControllerResourceRuntime(state: AppShellState, client: QueryClient) {
	const enabled = createResourceEnablement(state);
	const catalogCache = createSettingsCatalogCache(client);
	const queries = {
		languages: createLanguagesQuery(enabled.languages),
		tags: createTagsQuery(enabled.tags),
		users: createUsersQuery(enabled.users),
		downloadClients: createDownloadClientsQuery(enabled.downloadClients),
		indexers: createIndexersQuery(enabled.indexers),
		metadataProviders: createMetadataProvidersQuery(enabled.metadataProviders),
		subtitleProviders: createSubtitleProvidersQuery(enabled.subtitleProviders),
		libraryFolders: createLibraryFoldersQuery(enabled.libraryFolders),
		pathMappings: createPathMappingsQuery(enabled.pathMappings),
		mediaProfiles: createMediaProfilesQuery(enabled.mediaProfiles),
		customFormats: createCustomFormatsQuery(enabled.customFormats)
	};
	const server = createServerResourceRuntime(client, enabled);
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
	return { client, catalogCache, queries, server, scans, properties };
}
