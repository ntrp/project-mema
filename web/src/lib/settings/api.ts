import { client } from '$lib/api/client';

import {
	normalizeCustomFormatForm,
	normalizeDownloadClientForm,
	normalizeIndexerForm,
	normalizeLibraryFolderForm,
	normalizeMediaProfileForm,
	normalizeMetadataProviderForm,
	normalizePathMappingForm,
	normalizeUserCreateForm,
	normalizeUserUpdateForm
} from './forms';
import type {
	CustomFormat,
	CustomFormatForm,
	DiscoverBlacklistItem,
	DiscoverBlacklistRequest,
	DownloadClientForm,
	FileNamingSettings,
	FileNamingSettingsRequest,
	IndexerForm,
	LibraryFolderForm,
	LibraryFolderOption,
	LibraryFolderOptionListResponse,
	LibraryMediaKind,
	LibraryScanItemMatchRequest,
	ManualImportRequest,
	MediaAdvancedSearchRequest,
	MediaItemCreateRequest,
	MediaItemUpdateRequest,
	MediaProfile,
	MediaProfileForm,
	MediaRequestApproveRequest,
	MediaRequestCreateRequest,
	MediaSearchRequest,
	MediaType,
	MetadataCacheResponse,
	MetadataProviderForm,
	MetadataProviderType,
	PathMappingForm,
	QualitySizeSettingRequest,
	QualitySizeSettingsResponse,
	ReleaseCandidate,
	SessionResponse,
	SettingsData,
	SystemEventSettings,
	SystemEventSettingsRequest,
	SystemLogFile,
	SystemLogFileSettings,
	SystemLogFileSettingsRequest,
	SystemLogLevel,
	SystemLogLevelResponse,
	SystemStatusResponse,
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

export async function getSystemLogLevel(): Promise<SystemLogLevelResponse> {
	const { data, error } = await client.GET('/system/log-level');

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Log level request did not return a result');
	}
	return data;
}

export async function updateSystemLogLevel(level: SystemLogLevel): Promise<SystemLogLevelResponse> {
	const { data, error } = await client.PUT('/system/log-level', {
		body: { level }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Log level update did not return a result');
	}
	return data;
}

export async function getSystemStatus(): Promise<SystemStatusResponse> {
	const { data, error } = await client.GET('/system/status');

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('System status request did not return a result');
	}
	return data;
}

export async function getSystemLogFileSettings(): Promise<SystemLogFileSettings> {
	const { data, error } = await client.GET('/system/log-file-settings');

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Log file settings request did not return a result');
	}
	return data;
}

export async function updateSystemLogFileSettings(
	request: SystemLogFileSettingsRequest
): Promise<SystemLogFileSettings> {
	const { data, error } = await client.PUT('/system/log-file-settings', { body: request });

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Log file settings update did not return a result');
	}
	return data;
}

export async function listSystemLogFiles(): Promise<SystemLogFile[]> {
	const { data, error } = await client.GET('/system/log-files');

	if (error) {
		throw new Error(error.message);
	}
	return data?.files ?? [];
}

export async function downloadSystemLogFile(name: string) {
	const response = await globalThis.fetch(
		`/api/system/log-files/${encodeURIComponent(name)}/download`,
		{
			credentials: 'include'
		}
	);
	if (!response.ok) {
		throw new Error('Could not download log file');
	}
	return response.blob();
}

export async function listSystemEvents(options: { before?: string; limit?: number } = {}) {
	const { data, error } = await client.GET('/system/events', {
		params: { query: options }
	});

	if (error) {
		throw new Error(error.message);
	}
	return data ?? { events: [], hasMore: false };
}

export async function deleteSystemEvent(id: string) {
	const { error } = await client.DELETE('/system/events/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function clearSystemEvents() {
	const { error } = await client.DELETE('/system/events');

	if (error) {
		throw new Error(error.message);
	}
}

export async function getSystemEventSettings(): Promise<SystemEventSettings> {
	const { data, error } = await client.GET('/system/event-settings');

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Event settings request did not return a result');
	}
	return data;
}

export async function updateSystemEventSettings(
	request: SystemEventSettingsRequest
): Promise<SystemEventSettings> {
	const { data, error } = await client.PUT('/system/event-settings', { body: request });

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Event settings update did not return a result');
	}
	return data;
}

export async function loadSettings(): Promise<SettingsData> {
	const [
		clientResult,
		indexerResult,
		metadataProviderResult,
		metadataCacheResult,
		libraryFolderResult,
		pathMappingResult,
		mediaProfileResult,
		customFormatResult,
		userResult,
		tagResult
	] = await Promise.all([
		client.GET('/settings/download-clients'),
		client.GET('/settings/indexers'),
		client.GET('/settings/metadata-providers'),
		client.GET('/settings/metadata-cache'),
		client.GET('/settings/library/folders'),
		client.GET('/settings/library/path-mappings'),
		client.GET('/settings/profiles'),
		client.GET('/settings/custom-formats'),
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
	if (pathMappingResult.error) {
		throw new Error(pathMappingResult.error.message);
	}
	if (mediaProfileResult.error) {
		throw new Error(mediaProfileResult.error.message);
	}
	if (customFormatResult.error) {
		throw new Error(customFormatResult.error.message);
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
		pathMappings: pathMappingResult.data?.mappings ?? [],
		mediaProfiles: mediaProfileResult.data?.profiles ?? [],
		customFormats: customFormatResult.data?.formats ?? [],
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

export async function listQualitySizeSettings(): Promise<QualitySizeSettingsResponse> {
	const { data, error } = await client.GET('/settings/quality-sizes');

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Quality size settings were not returned');
	}
	return data;
}

export async function updateQualitySizeSettings(qualities: QualitySizeSettingRequest[]) {
	const { data, error } = await client.PUT('/settings/quality-sizes', {
		body: { qualities }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Quality size settings were not returned');
	}
	return data;
}

export async function getFileNamingSettings(): Promise<FileNamingSettings> {
	const { data, error } = await client.GET('/settings/file-naming');

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('File naming settings were not returned');
	}
	return data;
}

export async function updateFileNamingSettings(request: FileNamingSettingsRequest) {
	const { data, error } = await client.PUT('/settings/file-naming', {
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('File naming settings were not returned');
	}
	return data;
}

export async function listMediaProfiles(): Promise<MediaProfile[]> {
	const { data, error } = await client.GET('/settings/profiles');

	if (error) {
		throw new Error(error.message);
	}
	return data?.profiles ?? [];
}

export async function listCustomFormats(): Promise<CustomFormat[]> {
	const { data, error } = await client.GET('/settings/custom-formats');

	if (error) {
		throw new Error(error.message);
	}
	return data?.formats ?? [];
}

export async function testCustomFormatParsing(fileName: string) {
	const { data, error } = await client.POST('/settings/custom-formats/test-parsing', {
		body: { fileName }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Parsing result was not returned');
	}
	return data;
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

export async function loadMediaDiscoverSection(sectionId: string, page = 1, limit = 20) {
	const { data, error } = await client.GET('/media/discover/{sectionId}', {
		params: { path: { sectionId }, query: { page, limit } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Discovery section request did not return a result');
	}
	return data;
}

export async function listDiscoverBlacklist(): Promise<DiscoverBlacklistItem[]> {
	const { data, error } = await client.GET('/media/discover/blacklist');

	if (error) {
		throw new Error(error.message);
	}
	return data?.items ?? [];
}

export async function addDiscoverBlacklistItem(
	request: DiscoverBlacklistRequest
): Promise<DiscoverBlacklistItem> {
	const { data, error } = await client.POST('/media/discover/blacklist', { body: request });

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Blacklist request did not return a result');
	}
	return data;
}

export async function deleteDiscoverBlacklistItem(id: string) {
	const { error } = await client.DELETE('/media/discover/blacklist/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
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

export async function getMediaCollection(provider: MetadataProviderType, collectionId: string) {
	const { data, error } = await client.GET('/media/collections/{provider}/{collectionId}', {
		params: { path: { provider, collectionId } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media collection was not returned');
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

export async function createMediaItem(request: MediaItemCreateRequest) {
	const { data, error } = await client.POST('/media/items', { body: request });

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media item was not returned');
	}
	return data;
}

export async function updateMediaItem(id: string, request: MediaItemUpdateRequest) {
	const { data, error } = await client.PUT('/media/items/{id}', {
		params: { path: { id } },
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media item was not returned');
	}
	return data;
}

export async function refreshMediaItemMetadata(id: string) {
	const { data, error } = await client.POST('/media/items/{id}/metadata/refresh', {
		params: { path: { id } }
	});

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

export async function deleteMediaItem(id: string, options: { keepFiles?: boolean } = {}) {
	const { error } = await client.DELETE('/media/items/{id}', {
		params: { path: { id }, query: { keepFiles: options.keepFiles } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function rescanMediaItemFiles(id: string) {
	const { data, error } = await client.POST('/media/items/{id}/files/rescan', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media item was not returned');
	}
	return data;
}

export async function deleteMediaItemFile(id: string, path: string) {
	const { data, error } = await client.POST('/media/items/{id}/files/delete', {
		params: { path: { id } },
		body: { path }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media item was not returned');
	}
	return data;
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

export async function enqueueMediaAutomaticSearch(id: string) {
	const { data, error } = await client.POST('/media/items/{id}/automatic-searches', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Automatic search job was not returned');
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

export async function cancelDownloadActivity(id: string) {
	const { data, error } = await client.POST('/activity/downloads/{id}/cancel', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Download activity was not returned');
	}
	return data;
}

export async function deleteDownloadActivity(id: string) {
	const { error } = await client.DELETE('/activity/downloads/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function manualImportDownloadActivity(id: string, body: ManualImportRequest) {
	const { data, error } = await client.POST('/activity/downloads/{id}/manual-import', {
		params: { path: { id } },
		body
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Download activity was not returned');
	}
	return data;
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

export async function testDownloadClientConfig(form: DownloadClientForm) {
	const { data, error } = await client.POST('/settings/download-clients/test', {
		body: normalizeDownloadClientForm(form)
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

export async function saveMediaProfile(form: MediaProfileForm) {
	const body = normalizeMediaProfileForm(form);
	const result = form.id
		? await client.PUT('/settings/profiles/{id}', {
				params: { path: { id: form.id } },
				body
			})
		: await client.POST('/settings/profiles', { body });

	if (result.error) {
		throw new Error(result.error.message);
	}
}

export async function saveCustomFormat(form: CustomFormatForm) {
	const body = normalizeCustomFormatForm(form);
	const result = form.id
		? await client.PUT('/settings/custom-formats/{id}', {
				params: { path: { id: form.id } },
				body
			})
		: await client.POST('/settings/custom-formats', { body });

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

export async function deleteMediaProfile(id: string) {
	const { error } = await client.DELETE('/settings/profiles/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function deleteCustomFormat(id: string) {
	const { error } = await client.DELETE('/settings/custom-formats/{id}', {
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

export async function scanLibraryFolder(id: string) {
	const { data, error } = await client.POST('/settings/library/folders/{id}/scan', {
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

export async function savePathMapping(form: PathMappingForm) {
	const { data, error } = await client.POST('/settings/library/path-mappings', {
		body: normalizePathMappingForm(form)
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Path mapping was not returned');
	}
	return data;
}

export async function deletePathMapping(id: string) {
	const { error } = await client.DELETE('/settings/library/path-mappings/{id}', {
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
