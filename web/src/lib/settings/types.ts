import type { components } from '$lib/api/generated/schema';

export type DownloadClient = components['schemas']['DownloadClient'];
export type DownloadClientRequest = components['schemas']['DownloadClientRequest'];
export type DownloadClientType = components['schemas']['DownloadClientType'];
export type ManagedUser = components['schemas']['ManagedUser'];
export type SessionResponse = components['schemas']['SessionResponse'];
export type UserProfile = components['schemas']['UserProfile'];
export type UserProfileUpdateRequest = components['schemas']['UserProfileUpdateRequest'];
export type UserCreateRequest = components['schemas']['UserCreateRequest'];
export type UserRole = components['schemas']['UserRole'];
export type UserSummary = components['schemas']['UserSummary'];
export type UserUpdateRequest = components['schemas']['UserUpdateRequest'];
export type Indexer = components['schemas']['Indexer'];
export type IndexerHealthStatus = components['schemas']['IndexerHealthStatus'];
export type IndexerRequest = components['schemas']['IndexerRequest'];
export type IndexerSearchResponse = components['schemas']['IndexerSearchResponse'];
export type IndexerSearchSettings = components['schemas']['IndexerSearchSettings'];
export type IndexerSearchCacheEntry = components['schemas']['IndexerSearchCacheEntry'];
export type IndexerSearchCacheStats = components['schemas']['IndexerSearchCacheStats'];
export type IndexerSearchHistoryEntry = components['schemas']['IndexerSearchHistoryEntry'];
export type IndexerProtocol = components['schemas']['IndexerProtocol'];
export type IndexerPrivacy = components['schemas']['IndexerPrivacy'];
export type IndexerMediaType = components['schemas']['IndexerMediaType'];
export type IndexerCatalogResponse = components['schemas']['IndexerCatalogResponse'];
export type IndexerCatalogEntry = components['schemas']['IndexerCatalogEntry'];
export type IndexerAppProfile = components['schemas']['IndexerAppProfile'];
export type IndexerProxy = components['schemas']['IndexerProxy'];
export type IndexerProxyRequest = components['schemas']['IndexerProxyRequest'];
export type IndexerBulkUpdateRequest = components['schemas']['IndexerBulkUpdateRequest'];
export type MetadataProvider = components['schemas']['MetadataProvider'];
export type MetadataCacheClearResponse = components['schemas']['MetadataCacheClearResponse'];
export type MetadataCacheEntry = components['schemas']['MetadataCacheEntry'];
export type MetadataCacheResponse = components['schemas']['MetadataCacheResponse'];
export type MetadataCacheStats = components['schemas']['MetadataCacheStats'];
export type MetadataSearchHistoryEntry = components['schemas']['MetadataSearchHistoryEntry'];
export type MetadataProviderRequest = components['schemas']['MetadataProviderRequest'];
export type MetadataProviderType = components['schemas']['MetadataProviderType'];
export type SubtitleProvider = components['schemas']['SubtitleProvider'];
export type SubtitleProviderRequest = components['schemas']['SubtitleProviderRequest'];
export type SubtitleProviderType = components['schemas']['SubtitleProviderType'];
export type IntegrationTestResponse = components['schemas']['IntegrationTestResponse'];
export type SystemLogEntry = components['schemas']['SystemLogEntry'];
export type SystemStatusResponse = components['schemas']['SystemStatusResponse'];
export type SystemEvent = components['schemas']['SystemEvent'];
export type SystemEventListResponse = components['schemas']['SystemEventListResponse'];
export type SystemEventSeverity = components['schemas']['SystemEventSeverity'];
export type SystemEventSettings = components['schemas']['SystemEventSettings'];
export type SystemEventSettingsRequest = components['schemas']['SystemEventSettingsRequest'];
export type SystemLogFile = components['schemas']['SystemLogFile'];
export type SystemLogFileSettings = components['schemas']['SystemLogFileSettings'];
export type SystemLogFileSettingsRequest = components['schemas']['SystemLogFileSettingsRequest'];
export type SystemLogLevel = components['schemas']['SystemLogLevel'];
export type SystemLogLevelResponse = components['schemas']['SystemLogLevelResponse'];
export type SystemJob = components['schemas']['SystemJob'];
export type SystemJobListResponse = components['schemas']['SystemJobListResponse'];
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
export type DiscoverBlacklistItem = components['schemas']['DiscoverBlacklistItem'];
export type DiscoverBlacklistRequest = components['schemas']['DiscoverBlacklistRequest'];
export type DiscoverMovieFacetOption = components['schemas']['DiscoverMovieFacetOption'];
export type DiscoverMovieSearchResponse = components['schemas']['DiscoverMovieSearchResponse'];
export type MediaAdvancedSearchRequest = components['schemas']['MediaAdvancedSearchRequest'];
export type MediaDiscoverSection = components['schemas']['MediaDiscoverSection'];
export type MediaMetadataDetails = components['schemas']['MediaMetadataDetails'];
export type MediaMetadataFact = components['schemas']['MediaMetadataFact'];
export type MediaMetadataEpisode = components['schemas']['MediaMetadataEpisode'];
export type MediaMetadataSeason = components['schemas']['MediaMetadataSeason'];
export type PersonAppearance = components['schemas']['PersonAppearance'];
export type PersonDetails = components['schemas']['PersonDetails'];
export type PersonSearchResult = components['schemas']['PersonSearchResult'];
export type MediaCollection = components['schemas']['MediaCollection'];
export type MediaRequest = components['schemas']['MediaRequest'];
export type MediaRequestApproveRequest = components['schemas']['MediaRequestApproveRequest'];
export type MediaRequestCreateRequest = components['schemas']['MediaRequestCreateRequest'];
export type MediaRequestStatus = components['schemas']['MediaRequestStatus'];
export type MediaSearchGroup = components['schemas']['MediaSearchGroup'];
export type MediaItem = components['schemas']['MediaItem'];
export type MediaItemCreateRequest = components['schemas']['MediaItemCreateRequest'];
export type MediaItemUpdateRequest = components['schemas']['MediaItemUpdateRequest'];
export type MediaItemRequest = components['schemas']['MediaItemRequest'];
export type MediaItemStatus = components['schemas']['MediaItemStatus'];
export type MediaMonitorMode = components['schemas']['MediaMonitorMode'];
export type SeriesType = components['schemas']['SeriesType'];
export type MinimumAvailability = components['schemas']['MinimumAvailability'];
export type ReleaseCandidate = components['schemas']['ReleaseCandidate'];
export type ReleaseBlocklistItem = components['schemas']['ReleaseBlocklistItem'];
export type ReleaseOverrideDetails = components['schemas']['ReleaseOverrideDetails'];
export type DownloadActivity = components['schemas']['DownloadActivity'];
export type DownloadActivityStatus = components['schemas']['DownloadActivity']['status'];
export type ImportMode = components['schemas']['ImportMode'];
export type ManualImportRequest = components['schemas']['ManualImportRequest'];
export type JobEnqueueResponse = components['schemas']['JobEnqueueResponse'];
export type Tag = components['schemas']['Tag'];
export type TagRequest = components['schemas']['TagRequest'];
export type Language = components['schemas']['Language'];
export type LanguageRequest = components['schemas']['LanguageRequest'];
export type LanguageUpdateRequest = components['schemas']['LanguageUpdateRequest'];
export type QualitySizeSetting = components['schemas']['QualitySizeSetting'];
export type QualitySizeSettingRequest = components['schemas']['QualitySizeSettingRequest'];
export type QualitySizeSettingsResponse = components['schemas']['QualitySizeSettingsResponse'];
export type QualitySizeSettingsUpdateRequest =
	components['schemas']['QualitySizeSettingsUpdateRequest'];
export type MediaProfile = components['schemas']['MediaProfile'];
export type MediaProfileRequest = components['schemas']['MediaProfileRequest'];
export type MediaProfileLanguageScore = components['schemas']['MediaProfileLanguageScore'];
export type MediaProfileSubtitleLanguage = components['schemas']['MediaProfileSubtitleLanguage'];
export type MediaProfileCustomFormatScore = components['schemas']['MediaProfileCustomFormatScore'];
export type FileNamingSettings = components['schemas']['FileNamingSettings'];
export type FileNamingSettingsRequest = components['schemas']['FileNamingSettingsRequest'];
export type CustomFormat = components['schemas']['CustomFormat'];
export type CustomFormatRequest = components['schemas']['CustomFormatRequest'];
export type CustomFormatParsingResponse = components['schemas']['CustomFormatParsingResponse'];
export type CustomFormatSpec = components['schemas']['CustomFormatSpec'];
export type CustomFormatSpecType = components['schemas']['CustomFormatSpecType'];

export type DownloadClientForm = DownloadClientRequest & { id?: string };
export type IndexerForm = Omit<IndexerRequest, 'categories'> & {
	id?: string;
	categoriesText: string;
};
export type IndexerProxyForm = IndexerProxyRequest & { id?: string };
export type MetadataProviderForm = MetadataProviderRequest & { id?: string };
export type SubtitleProviderForm = SubtitleProviderRequest & { id?: string };
export type LibraryFolderForm = LibraryFolderRequest;
export type PathMappingForm = PathMappingRequest;
export type MediaProfileForm = MediaProfileRequest & { id?: string };
export type CustomFormatForm = CustomFormatRequest & { id?: string };
export type LanguageForm = {
	code: string;
	originalCode?: string;
	displayName: string;
	aliasesText: string;
};
export type UserForm = {
	id?: string;
	username: string;
	password: string;
	role: UserRole;
};

export interface QualityProfileOption {
	id: string;
	name: string;
	targetLanguages?: string[];
	removeNonEnabledLanguages?: boolean;
}

export type AppView =
	| 'home'
	| 'settings'
	| 'system'
	| 'profile'
	| 'advanced-search'
	| 'metadata-detail'
	| 'media-people'
	| 'person-detail'
	| 'media-collection'
	| 'related-section'
	| 'discover-section'
	| 'discover-movies'
	| 'discover-series';
export type HomeSection =
	| 'discover'
	| 'blacklist'
	| 'requests'
	| 'movies'
	| 'series'
	| 'wanted'
	| 'activity';
export type ActivitySection = 'queue' | 'history' | 'blocklist';
export type SystemSection = 'status' | 'indexing' | 'metadata' | 'jobs' | 'logs' | 'events';
export type SettingsSection =
	| 'general'
	| 'library'
	| 'download-clients'
	| 'indexers'
	| 'quality'
	| 'profiles'
	| 'custom-formats'
	| 'metadata'
	| 'subtitles'
	| 'languages'
	| 'tags'
	| 'users';
export type TagForm = TagRequest & { id?: string };
export interface SettingsData {
	downloadClients: DownloadClient[];
	indexers: Indexer[];
	indexerSearch: IndexerSearchResponse;
	metadataProviders: MetadataProvider[];
	subtitleProviders: SubtitleProvider[];
	metadataCache: MetadataCacheResponse;
	libraryFolders: LibraryFolder[];
	pathMappings: PathMapping[];
	mediaProfiles: MediaProfile[];
	customFormats: CustomFormat[];
	users: ManagedUser[];
	tags: Tag[];
	languages: Language[];
}

export type IntegrationTestResults = Record<string, IntegrationTestResponse | undefined>;

export interface ReleaseSearchState {
	loaded: boolean;
	releases: ReleaseCandidate[];
	errors: string[];
}

export type ReleaseSearchResults = Record<string, ReleaseSearchState | undefined>;
