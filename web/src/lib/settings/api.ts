import { client } from '$lib/api/client';

import {
	normalizeDownloadClientForm,
	normalizeIndexerForm,
	normalizeLibraryFolderForm,
	normalizeMetadataProviderForm,
	normalizeUserCreateForm,
	normalizeUserUpdateForm
} from './forms';
import type {
	DownloadClientForm,
	IndexerForm,
	LibraryFolderForm,
	LibraryFolderOption,
	LibraryFolderOptionListResponse,
	LibraryMediaKind,
	LibraryScanItemMatchRequest,
	MediaAdvancedSearchRequest,
	MediaItemRequest,
	MediaRequestApproveRequest,
	MediaRequestCreateRequest,
	MediaSearchRequest,
	MediaType,
	MetadataCacheResponse,
	MetadataProviderForm,
	MetadataProviderType,
	ReleaseCandidate,
	SessionResponse,
	SettingsData,
	TagForm,
	UserForm
} from './types';

export async function currentSession(): Promise<SessionResponse | undefined> {
	const { data } = await client.GET('/auth/session');
	return data;
}

export async function currentSessionAuthenticated() {
	const data = await currentSession();
	return Boolean(data?.authenticated);
}

export async function login(username: string, password: string) {
	const { data, error } = await client.POST('/auth/login', {
		body: { username, password }
	});

	if (error || !data?.authenticated) {
		throw new Error(error?.message ?? 'Login failed');
	}
	return data;
}

export async function logout() {
	const { error } = await client.POST('/auth/logout');

	if (error) {
		throw new Error(error.message);
	}
}

export async function loadSettings(): Promise<SettingsData> {
	const [
		clientResult,
		indexerResult,
		metadataProviderResult,
		metadataCacheResult,
		libraryFolderResult,
		userResult,
		tagResult
	] = await Promise.all([
		client.GET('/settings/download-clients'),
		client.GET('/settings/indexers'),
		client.GET('/settings/metadata-providers'),
		client.GET('/settings/metadata-cache'),
		client.GET('/settings/library/folders'),
		client.GET('/settings/users'),
		client.GET('/settings/tags')
	]);

	if (clientResult.error) {
		throw new Error(clientResult.error.message);
	}
	if (indexerResult.error) {
		throw new Error(indexerResult.error.message);
	}
	if (metadataProviderResult.error) {
		throw new Error(metadataProviderResult.error.message);
	}
	if (metadataCacheResult.error) {
		throw new Error(metadataCacheResult.error.message);
	}
	if (libraryFolderResult.error) {
		throw new Error(libraryFolderResult.error.message);
	}
	if (userResult.error) {
		throw new Error(userResult.error.message);
	}
	if (tagResult.error) {
		throw new Error(tagResult.error.message);
	}

	return {
		downloadClients: clientResult.data?.clients ?? [],
		indexers: indexerResult.data?.indexers ?? [],
		metadataProviders: metadataProviderResult.data?.providers ?? [],
		metadataCache: metadataCacheResult.data ?? emptyMetadataCache(),
		libraryFolders: libraryFolderResult.data?.folders ?? [],
		users: userResult.data?.users ?? [],
		tags: tagResult.data?.tags ?? []
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
		entries: []
	};
}

export async function searchMedia(request: MediaSearchRequest) {
	const { data, error } = await client.POST('/media/search', { body: request });

	if (error) {
		throw new Error(error.message);
	}
	return data?.results ?? [];
}

export async function loadMediaDiscoverSections() {
	const { data, error } = await client.GET('/media/discover');

	if (error) {
		throw new Error(error.message);
	}
	return data?.sections ?? [];
}

export type AutocompleteSearchScope = 'library' | 'providers' | 'all';

export async function autocompleteMedia(query: string, scope: AutocompleteSearchScope = 'all') {
	const { data, error } = await client.GET('/media/autocomplete', {
		params: {
			query: {
				query,
				includeLibrary: scope !== 'providers',
				includeProviders: scope !== 'library'
			}
		}
	});

	if (error) {
		throw new Error(error.message);
	}
	return data?.groups ?? [];
}

export async function advancedSearchMedia(request: MediaAdvancedSearchRequest) {
	const { data, error } = await client.POST('/media/advanced-search', { body: request });

	if (error) {
		throw new Error(error.message);
	}
	return data?.groups ?? [];
}

export async function getMediaMetadataDetails(
	provider: MetadataProviderType,
	type: MediaType,
	externalId: string
) {
	const { data, error } = await client.GET('/media/metadata/{provider}/{type}/{externalId}', {
		params: { path: { provider, type, externalId } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media details were not returned');
	}
	return data;
}

export async function listMediaItems() {
	const { data, error } = await client.GET('/media/items');

	if (error) {
		throw new Error(error.message);
	}
	return data?.items ?? [];
}

export async function createMediaItem(request: MediaItemRequest) {
	const { data, error } = await client.POST('/media/items', { body: request });

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media item was not returned');
	}
	return data;
}

export async function listMediaRequests() {
	const { data, error } = await client.GET('/media/requests');

	if (error) {
		throw new Error(error.message);
	}
	return data?.requests ?? [];
}

export async function createMediaRequest(request: MediaRequestCreateRequest) {
	const { data, error } = await client.POST('/media/requests', { body: request });

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media request was not returned');
	}
	return data;
}

export async function getMediaRequest(id: string) {
	const { data, error } = await client.GET('/media/requests/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media request was not returned');
	}
	return data;
}

export async function approveMediaRequest(id: string, request: MediaRequestApproveRequest) {
	const { data, error } = await client.POST('/media/requests/{id}/approve', {
		params: { path: { id } },
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media request approval was not returned');
	}
	return data;
}

export async function deleteMediaItem(id: string) {
	const { error } = await client.DELETE('/media/items/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function searchMediaReleases(id: string) {
	const { data, error } = await client.GET('/media/items/{id}/releases', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	return {
		releases: data?.releases ?? [],
		errors: data?.errors ?? []
	};
}

export async function enqueueMediaReleaseSearch(id: string) {
	const { data, error } = await client.POST('/media/items/{id}/release-searches', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Release search job was not returned');
	}
	return data;
}

export async function grabMediaRelease(id: string, release: ReleaseCandidate) {
	const { data, error } = await client.POST('/media/items/{id}/grab', {
		params: { path: { id } },
		body: {
			releaseId: release.id
		}
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Download grab did not return a result');
	}
	return data;
}

export async function listDownloadActivity() {
	const { data, error } = await client.GET('/activity/downloads');

	if (error) {
		throw new Error(error.message);
	}
	return data?.activities ?? [];
}

export async function saveDownloadClient(form: DownloadClientForm) {
	const body = normalizeDownloadClientForm(form);
	const result = form.id
		? await client.PUT('/settings/download-clients/{id}', {
				params: { path: { id: form.id } },
				body
			})
		: await client.POST('/settings/download-clients', { body });

	if (result.error) {
		throw new Error(result.error.message);
	}
}

export async function testDownloadClient(id: string) {
	const { data, error } = await client.POST('/settings/download-clients/{id}/test', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Download client test did not return a result');
	}
	return data;
}

export async function saveIndexer(form: IndexerForm) {
	const body = normalizeIndexerForm(form);
	const result = form.id
		? await client.PUT('/settings/indexers/{id}', {
				params: { path: { id: form.id } },
				body
			})
		: await client.POST('/settings/indexers', { body });

	if (result.error) {
		throw new Error(result.error.message);
	}
}

export async function testIndexer(id: string) {
	const { data, error } = await client.POST('/settings/indexers/{id}/test', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Indexer test did not return a result');
	}
	return data;
}

export async function saveMetadataProvider(form: MetadataProviderForm) {
	const body = normalizeMetadataProviderForm(form);
	const result = form.id
		? await client.PUT('/settings/metadata-providers/{id}', {
				params: { path: { id: form.id } },
				body
			})
		: await client.POST('/settings/metadata-providers', { body });

	if (result.error) {
		throw new Error(result.error.message);
	}
}

export async function saveUser(form: UserForm) {
	const result = form.id
		? await client.PUT('/settings/users/{id}', {
				params: { path: { id: form.id } },
				body: normalizeUserUpdateForm(form)
			})
		: await client.POST('/settings/users', {
				body: normalizeUserCreateForm(form)
			});

	if (result.error) {
		throw new Error(result.error.message);
	}
}

export async function saveTag(form: TagForm) {
	const body = { name: form.name.trim() };
	const result = form.id
		? await client.PUT('/settings/tags/{id}', {
				params: { path: { id: form.id } },
				body
			})
		: await client.POST('/settings/tags', { body });

	if (result.error) {
		throw new Error(result.error.message);
	}
}

export async function testMetadataProvider(id: string) {
	const { data, error } = await client.POST('/settings/metadata-providers/{id}/test', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Metadata provider test did not return a result');
	}
	return data;
}

export async function getMetadataCache() {
	const { data, error } = await client.GET('/settings/metadata-cache');

	if (error) {
		throw new Error(error.message);
	}
	return data ?? emptyMetadataCache();
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

export async function deleteDownloadClient(id: string) {
	const { error } = await client.DELETE('/settings/download-clients/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function deleteIndexer(id: string) {
	const { error } = await client.DELETE('/settings/indexers/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function deleteMetadataProvider(id: string) {
	const { error } = await client.DELETE('/settings/metadata-providers/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function deleteUser(id: string) {
	const { error } = await client.DELETE('/settings/users/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function deleteTag(id: string) {
	const { error } = await client.DELETE('/settings/tags/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function saveLibraryFolder(form: LibraryFolderForm) {
	const { data, error } = await client.POST('/settings/library/folders', {
		body: normalizeLibraryFolderForm(form)
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Library scan was not returned');
	}
	return data;
}

export async function listLibraryFolderOptions(
	path?: string
): Promise<LibraryFolderOptionListResponse> {
	const { data, error } = await client.GET('/settings/library/folder-options', {
		params: { query: { path } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Folder options were not returned');
	}
	return data;
}

export async function createLibraryFolderOption(
	parentPath: string,
	name: string
): Promise<LibraryFolderOption> {
	const { data, error } = await client.POST('/settings/library/folder-options', {
		body: { parentPath, name }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Folder was not returned');
	}
	return data;
}

export async function deleteLibraryFolder(id: string) {
	const { error } = await client.DELETE('/settings/library/folders/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function getLibraryScan(id: string) {
	const { data, error } = await client.GET('/settings/library/scans/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Library scan was not returned');
	}
	return data;
}

export async function matchLibraryScanItem(
	scanId: string,
	itemId: string,
	request: LibraryScanItemMatchRequest
) {
	const { data, error } = await client.POST('/settings/library/scans/{id}/items/{itemId}/match', {
		params: { path: { id: scanId, itemId } },
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Library match was not returned');
	}
	return data;
}

export function mediaTypeForLibraryKind(kind: LibraryMediaKind) {
	return kind === 'series' || kind === 'anime_series' ? 'series' : 'movie';
}
