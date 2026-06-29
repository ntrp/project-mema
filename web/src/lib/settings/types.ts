import type { components } from '$lib/api/generated/schema';

export type DownloadClient = components['schemas']['DownloadClient'];
export type DownloadClientRequest = components['schemas']['DownloadClientRequest'];
export type DownloadClientType = components['schemas']['DownloadClientType'];
export type ManagedUser = components['schemas']['ManagedUser'];
export type SessionResponse = components['schemas']['SessionResponse'];
export type UserCreateRequest = components['schemas']['UserCreateRequest'];
export type UserRole = components['schemas']['UserRole'];
export type UserSummary = components['schemas']['UserSummary'];
export type UserUpdateRequest = components['schemas']['UserUpdateRequest'];
export type Indexer = components['schemas']['Indexer'];
export type IndexerRequest = components['schemas']['IndexerRequest'];
export type IndexerType = components['schemas']['IndexerType'];
export type MetadataProvider = components['schemas']['MetadataProvider'];
export type MetadataCacheClearResponse = components['schemas']['MetadataCacheClearResponse'];
export type MetadataCacheEntry = components['schemas']['MetadataCacheEntry'];
export type MetadataCacheResponse = components['schemas']['MetadataCacheResponse'];
export type MetadataCacheStats = components['schemas']['MetadataCacheStats'];
export type MetadataProviderRequest = components['schemas']['MetadataProviderRequest'];
export type MetadataProviderType = components['schemas']['MetadataProviderType'];
export type IntegrationTestResponse = components['schemas']['IntegrationTestResponse'];
export type SystemLogEntry = components['schemas']['SystemLogEntry'];
export type SystemLogLevel = components['schemas']['SystemLogLevel'];
export type SystemLogLevelResponse = components['schemas']['SystemLogLevelResponse'];
export type LibraryFolder = components['schemas']['LibraryFolder'];
export type LibraryFolderOption = components['schemas']['LibraryFolderOption'];
export type LibraryFolderOptionCreateRequest =
	components['schemas']['LibraryFolderOptionCreateRequest'];
export type LibraryFolderOptionListResponse =
	components['schemas']['LibraryFolderOptionListResponse'];
export type LibraryFolderRequest = components['schemas']['LibraryFolderRequest'];
export type PathMapping = components['schemas']['PathMapping'];
export type PathMappingRequest = components['schemas']['PathMappingRequest'];
export type LibraryMediaKind = components['schemas']['LibraryMediaKind'];
export type LibraryScan = components['schemas']['LibraryScan'];
export type LibraryScanItem = components['schemas']['LibraryScanItem'];
export type LibraryScanItemMatchRequest = components['schemas']['LibraryScanItemMatchRequest'];
export type MediaType = components['schemas']['MediaType'];
export type MediaSearchRequest = components['schemas']['MediaSearchRequest'];
export type MediaSearchResult = components['schemas']['MediaSearchResult'];
export type MediaAdvancedSearchRequest = components['schemas']['MediaAdvancedSearchRequest'];
export type MediaDiscoverSection = components['schemas']['MediaDiscoverSection'];
export type MediaMetadataDetails = components['schemas']['MediaMetadataDetails'];
export type MediaRequest = components['schemas']['MediaRequest'];
export type MediaRequestApproveRequest = components['schemas']['MediaRequestApproveRequest'];
export type MediaRequestCreateRequest = components['schemas']['MediaRequestCreateRequest'];
export type MediaRequestStatus = components['schemas']['MediaRequestStatus'];
export type MediaSearchGroup = components['schemas']['MediaSearchGroup'];
export type MediaItem = components['schemas']['MediaItem'];
export type MediaItemRequest = components['schemas']['MediaItemRequest'];
export type MediaItemStatus = components['schemas']['MediaItemStatus'];
export type ReleaseCandidate = components['schemas']['ReleaseCandidate'];
export type DownloadActivity = components['schemas']['DownloadActivity'];
export type DownloadActivityStatus = components['schemas']['DownloadActivity']['status'];
export type JobEnqueueResponse = components['schemas']['JobEnqueueResponse'];
export type Tag = components['schemas']['Tag'];
export type TagRequest = components['schemas']['TagRequest'];
export type QualitySizeSetting = components['schemas']['QualitySizeSetting'];
export type QualitySizeSettingRequest = components['schemas']['QualitySizeSettingRequest'];
export type QualitySizeSettingsResponse = components['schemas']['QualitySizeSettingsResponse'];
export type QualitySizeSettingsUpdateRequest =
	components['schemas']['QualitySizeSettingsUpdateRequest'];
export type MediaProfile = components['schemas']['MediaProfile'];
export type MediaProfileRequest = components['schemas']['MediaProfileRequest'];
export type FileNamingSettings = components['schemas']['FileNamingSettings'];
export type FileNamingSettingsRequest = components['schemas']['FileNamingSettingsRequest'];

export type DownloadClientForm = DownloadClientRequest & { id?: string };
export type IndexerForm = Omit<IndexerRequest, 'categories'> & {
	id?: string;
	categoriesText: string;
};
export type MetadataProviderForm = MetadataProviderRequest & { id?: string };
export type LibraryFolderForm = LibraryFolderRequest;
export type PathMappingForm = PathMappingRequest;
export type MediaProfileForm = MediaProfileRequest & { id?: string };
export type UserForm = {
	id?: string;
	username: string;
	password: string;
	role: UserRole;
};

export interface QualityProfileOption {
	id: string;
	name: string;
}

export type AppView = 'home' | 'settings' | 'advanced-search' | 'metadata-detail';
export type HomeSection = 'discover' | 'requests' | 'movies' | 'series' | 'activity';
export type SettingsSection =
	| 'library'
	| 'download-clients'
	| 'indexers'
	| 'quality'
	| 'file-naming'
	| 'profiles'
	| 'metadata'
	| 'tags'
	| 'users'
	| 'system-logs';
export type TagForm = TagRequest & { id?: string };
export interface SettingsData {
	downloadClients: DownloadClient[];
	indexers: Indexer[];
	metadataProviders: MetadataProvider[];
	metadataCache: MetadataCacheResponse;
	libraryFolders: LibraryFolder[];
	pathMappings: PathMapping[];
	mediaProfiles: MediaProfile[];
	users: ManagedUser[];
	tags: Tag[];
}

export type IntegrationTestResults = Record<string, IntegrationTestResponse | undefined>;

export interface ReleaseSearchState {
	loaded: boolean;
	releases: ReleaseCandidate[];
	errors: string[];
}

export type ReleaseSearchResults = Record<string, ReleaseSearchState | undefined>;
