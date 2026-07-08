import type { components } from '$lib/api/generated/schema';

type S = components['schemas'];

export type * from './schema/platformTypes';

export type DownloadClient = S['DownloadClient'];
export type DownloadClientRequest = S['DownloadClientRequest'];
export type DownloadClientType = S['DownloadClientType'];
export type ManagedUser = S['ManagedUser'];
export type SessionResponse = S['SessionResponse'];
export type UserProfile = S['UserProfile'];
export type UserProfileUpdateRequest = S['UserProfileUpdateRequest'];
export type UserCreateRequest = S['UserCreateRequest'];
export type UserRole = S['UserRole'];
export type UserSummary = S['UserSummary'];
export type UserUpdateRequest = S['UserUpdateRequest'];
export type Indexer = S['Indexer'];
export type IndexerHealthStatus = S['IndexerHealthStatus'];
export type IndexerRequest = S['IndexerRequest'];
export type IndexerSearchResponse = S['IndexerSearchResponse'];
export type IndexerSearchSettings = S['IndexerSearchSettings'];
export type IndexerSearchCacheEntry = S['IndexerSearchCacheEntry'];
export type IndexerSearchCacheStats = S['IndexerSearchCacheStats'];
export type IndexerSearchHistoryEntry = S['IndexerSearchHistoryEntry'];
export type IndexerProtocol = S['IndexerProtocol'];
export type IndexerPrivacy = S['IndexerPrivacy'];
export type IndexerMediaType = S['IndexerMediaType'];
export type IndexerCatalogResponse = S['IndexerCatalogResponse'];
export type IndexerCatalogEntry = S['IndexerCatalogEntry'];
export type IndexerAppProfile = S['IndexerAppProfile'];
export type IndexerProxy = S['IndexerProxy'];
export type IndexerProxyRequest = S['IndexerProxyRequest'];
export type IndexerBulkUpdateRequest = S['IndexerBulkUpdateRequest'];
export type MetadataProvider = S['MetadataProvider'];
export type MetadataCacheClearResponse = S['MetadataCacheClearResponse'];
export type MetadataCacheEntry = S['MetadataCacheEntry'];
export type MetadataCacheResponse = S['MetadataCacheResponse'];
export type MetadataCacheStats = S['MetadataCacheStats'];
export type MetadataSearchHistoryEntry = S['MetadataSearchHistoryEntry'];
export type MetadataProviderRequest = S['MetadataProviderRequest'];
export type MetadataProviderType = S['MetadataProviderType'];
export type SubtitleProvider = S['SubtitleProvider'];
export type SubtitleProviderRequest = S['SubtitleProviderRequest'];
export type SubtitleProviderType = S['SubtitleProviderType'];
export type IntegrationTestResponse = S['IntegrationTestResponse'];
export type SystemLogEntry = S['SystemLogEntry'];
export type SystemStatusResponse = S['SystemStatusResponse'];
export type SystemEvent = S['SystemEvent'];
export type SystemEventListResponse = S['SystemEventListResponse'];
export type SystemEventSeverity = S['SystemEventSeverity'];
export type SystemEventSettings = S['SystemEventSettings'];
export type SystemEventSettingsRequest = S['SystemEventSettingsRequest'];
export type SystemLogFile = S['SystemLogFile'];
export type SystemLogFileSettings = S['SystemLogFileSettings'];
export type SystemLogFileSettingsRequest = S['SystemLogFileSettingsRequest'];
export type SystemLogLevel = S['SystemLogLevel'];
export type SystemLogLevelResponse = S['SystemLogLevelResponse'];
export type LibraryFolder = S['LibraryFolder'];
export type LibraryFolderKind = S['LibraryFolderKind'];
export type LibraryFolderOption = S['LibraryFolderOption'];
export type LibraryFolderOptionCreateRequest = S['LibraryFolderOptionCreateRequest'];
export type LibraryFolderOptionListResponse = S['LibraryFolderOptionListResponse'];
export type LibraryFolderRequest = S['LibraryFolderRequest'];
export type PathMapping = S['PathMapping'];
export type PathMappingRequest = S['PathMappingRequest'];
export type LibraryMediaKind = S['LibraryMediaKind'];
export type LibraryScan = S['LibraryScan'];
export type LibraryScanItem = S['LibraryScanItem'];
export type LibraryScanItemMatchRequest = S['LibraryScanItemMatchRequest'];
export type LibraryScanImportRequest = S['LibraryScanImportRequest'];
export type LibraryScanImportResponse = S['LibraryScanImportResponse'];
export type LibraryScanItemResetResponse = S['LibraryScanItemResetResponse'];
export type MediaType = S['MediaType'];
export type MediaSearchRequest = S['MediaSearchRequest'];
export type MediaSearchResult = S['MediaSearchResult'];
export type DiscoverBlacklistItem = S['DiscoverBlacklistItem'];
export type DiscoverBlacklistRequest = S['DiscoverBlacklistRequest'];
export type DiscoverMovieFacetOption = S['DiscoverMovieFacetOption'];
export type DiscoverMovieSearchResponse = S['DiscoverMovieSearchResponse'];
export type MediaAdvancedSearchRequest = S['MediaAdvancedSearchRequest'];
export type MediaDiscoverSection = S['MediaDiscoverSection'];
export type MediaMetadataDetails = S['MediaMetadataDetails'];
export type MediaMetadataFact = S['MediaMetadataFact'];
export type MediaMetadataEpisode = S['MediaMetadataEpisode'];
export type MediaMetadataSeason = S['MediaMetadataSeason'];
export type PersonAppearance = S['PersonAppearance'];
export type PersonDetails = S['PersonDetails'];
export type PersonSearchResult = S['PersonSearchResult'];
export type MediaCollection = S['MediaCollection'];
export type MediaRequest = S['MediaRequest'];
export type MediaRequestApproveRequest = S['MediaRequestApproveRequest'];
export type MediaRequestCreateRequest = S['MediaRequestCreateRequest'];
export type MediaRequestStatus = S['MediaRequestStatus'];
export type MediaSearchGroup = S['MediaSearchGroup'];
export type MediaItem = S['MediaItem'];
export type MediaItemSubtitle = S['MediaItemSubtitle'];
export type MediaItemSubtitleListResponse = S['MediaItemSubtitleListResponse'];
export type MediaItemSubtitleRetentionMode = S['MediaItemSubtitleRetentionMode'];
export type MediaItemSubtitleSelectionRequest = S['MediaItemSubtitleSelectionRequest'];
export type ManualSubtitleSearchRequest = S['ManualSubtitleSearchRequest'];
export type ManualSubtitleSearchResponse = S['ManualSubtitleSearchResponse'];
export type SubtitleCandidate = S['SubtitleCandidate'];
export type GrabSubtitleRequest = S['GrabSubtitleRequest'];
export type MediaComponentSource = S['MediaComponentSource'];
export type MediaComponentSourceListResponse = S['MediaComponentSourceListResponse'];
export type MediaComponentSourceRetainRequest = S['MediaComponentSourceRetainRequest'];
export type MediaComponentSourceRetentionState = S['MediaComponentSourceRetentionState'];
export type MediaComponentSourceRole = S['MediaComponentSourceRole'];
export type MediaComponentArtifact = S['MediaComponentArtifact'];
export type MediaComponentAssemblyRun = S['MediaComponentAssemblyRun'];
export type MediaComponentAssemblyInput = S['MediaComponentAssemblyInput'];
export type MediaComponentAssemblyRequest = S['MediaComponentAssemblyRequest'];
export type MediaComponentCompatibilityDecision = S['MediaComponentCompatibilityDecision'];
export type MediaComponentCompatibilityReviewState = S['MediaComponentCompatibilityReviewState'];
export type MediaItemCreateRequest = S['MediaItemCreateRequest'];
export type MediaItemUpdateRequest = S['MediaItemUpdateRequest'];
export type MediaItemRequest = S['MediaItemRequest'];
export type MediaItemStatus = S['MediaItemStatus'];
export type MediaFileHistoryEntry = S['MediaFileHistoryEntry'];
export type MediaFileHistoryResponse = S['MediaFileHistoryResponse'];
export type MediaFileTrackDeleteRequest = S['MediaFileTrackDeleteRequest'];
export type MediaRenameApplyResponse = S['MediaRenameApplyResponse'];
export type MediaRenamePreviewResponse = S['MediaRenamePreviewResponse'];
export type MediaRenamePreviewRow = S['MediaRenamePreviewRow'];
export type SubtitleSearchRequest = S['SubtitleSearchRequest'];
export type MediaMonitorMode = S['MediaMonitorMode'];
export type SeriesType = S['SeriesType'];
export type MinimumAvailability = S['MinimumAvailability'];
export type ReleaseCandidate = S['ReleaseCandidate'];
export type ReleaseBlocklistItem = S['ReleaseBlocklistItem'];
export type ReleaseOverrideDetails = S['ReleaseOverrideDetails'];
export type DownloadActivity = S['DownloadActivity'];
export type DownloadActivityStatus = S['DownloadActivity']['status'];
export type ImportMode = S['ImportMode'];
export type ManualImportRequest = S['ManualImportRequest'];
export type JobEnqueueResponse = S['JobEnqueueResponse'];
export type Tag = S['Tag'];
export type TagRequest = S['TagRequest'];
export type Language = S['Language'];
export type LanguageRequest = S['LanguageRequest'];
export type LanguageUpdateRequest = S['LanguageUpdateRequest'];
export type QualitySizeSetting = S['QualitySizeSetting'];
export type QualitySizeSettingRequest = S['QualitySizeSettingRequest'];
export type QualitySizeSettingsResponse = S['QualitySizeSettingsResponse'];
export type QualitySizeSettingsUpdateRequest = S['QualitySizeSettingsUpdateRequest'];
export type MediaProfile = S['MediaProfile'];
export type MediaProfileRequest = S['MediaProfileRequest'];
export type MediaProfileVideoTarget = S['MediaProfileVideoTarget'];
export type MediaProfileAudioTarget = S['MediaProfileAudioTarget'];
export type MediaProfileSubtitleTarget = S['MediaProfileSubtitleTarget'];
export type MediaProfileLossyTranscodePolicy = S['MediaProfileLossyTranscodePolicy'];
export type MediaProfileSubtitleMode = S['MediaProfileSubtitleMode'];
export type MediaProfileCustomFormatScore = S['MediaProfileCustomFormatScore'];
export type FileNamingSettings = S['FileNamingSettings'];
export type FileNamingSettingsRequest = S['FileNamingSettingsRequest'];
export type FileDeleteMode = S['FileDeleteMode'];
export type FileDeleteSettings = S['FileDeleteSettings'];
export type FileDeleteSettingsRequest = S['FileDeleteSettingsRequest'];
export type CustomFormat = S['CustomFormat'];
export type CustomFormatRequest = S['CustomFormatRequest'];
export type CustomFormatParsingResponse = S['CustomFormatParsingResponse'];
export type CustomFormatSpec = S['CustomFormatSpec'];
export type CustomFormatSpecType = S['CustomFormatSpecType'];

export type {
	ActivitySection,
	AppView,
	CustomFormatForm,
	DownloadClientForm,
	HomeSection,
	IndexerForm,
	IndexerProxyForm,
	IntegrationTestResults,
	LanguageForm,
	LibraryFolderForm,
	MediaProfileForm,
	MetadataProviderForm,
	PathMappingForm,
	QualityProfileOption,
	ReleaseSearchResults,
	ReleaseSearchState,
	SettingsData,
	SettingsSection,
	SubtitleProviderForm,
	SystemSection,
	TagForm,
	UserForm
} from './formTypes';
