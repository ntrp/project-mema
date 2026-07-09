import type {
	DownloadActivity,
	DiscoverBlacklistItem,
	GrabSubtitleRequest,
	HomeSection,
	ActivitySection,
	Language,
	LibraryFolder,
	MediaFileTrackDeleteRequest,
	MediaFulfillmentActionRequest,
	MediaDiscoverSection,
	MediaItem,
	MediaItemSubtitleSelectionRequest,
	MediaItemUpdateRequest,
	MediaRequest,
	MediaRequestApproveRequest,
	MediaSearchResult,
	QualityProfileOption,
	ReleaseBlocklistItem,
	ReleaseCandidate,
	ReleaseOverrideDetails,
	MediaComponentCompatibilityReviewState,
	MediaComponentSource
} from '$lib/settings/types';

export interface HomeAreaProps {
	activeSection: HomeSection;
	activitySection: ActivitySection;
	selectedMediaItemId?: string;
	selectedRequestId?: string;
	mediaItems: MediaItem[];
	mediaRequests: MediaRequest[];
	discoverSections: MediaDiscoverSection[];
	discoverBlacklist: DiscoverBlacklistItem[];
	libraryFolders: LibraryFolder[];
	languages: Language[];
	qualityProfiles: QualityProfileOption[];
	activities: DownloadActivity[];
	releaseBlocklist: ReleaseBlocklistItem[];
	loadingDiscover: boolean;
	loadingBlacklist: boolean;
	loadingMediaItems: boolean;
	addingKey?: string;
	blacklistingKey?: string;
	removingBlacklistId?: string;
	approvingRequestId?: string;
	searchingItemId?: string;
	scanningMediaItemId?: string;
	refreshingMetadataItemId?: string;
	savingMediaItemOptionsId?: string;
	grabbingKey?: string;
	deletingMediaItemId?: string;
	assemblingMediaItemId?: string;
	reviewingComponentDecisionId?: string;
	cancellingActivityId?: string;
	deletingActivityId?: string;
	deletingReleaseBlocklistId?: string;
	clearingReleaseBlocklist: boolean;
	canManage: boolean;
	loadingActivity: boolean;
	onAddMedia: (_candidate: MediaSearchResult) => void;
	onBlacklistMedia: (_candidate: MediaSearchResult) => void;
	onRemoveBlacklistMedia: (_item: DiscoverBlacklistItem) => void;
	onApproveMediaRequest: (_request: MediaRequest, _approval: MediaRequestApproveRequest) => void;
	onFindReleases: (_item: MediaItem, _query?: string) => void;
	onAutoSearchMedia: (_item: MediaItem) => void;
	onRescanMediaFiles: (_item: MediaItem) => void;
	onSearchMediaSubtitle?: (
		_item: MediaItem,
		_request: { languageId?: string; filePath?: string }
	) => void | Promise<void>;
	onGrabMediaSubtitle?: (_item: MediaItem, _request: GrabSubtitleRequest) => void | Promise<void>;
	onDeleteMediaSubtitle?: (_item: MediaItem, _subtitleId: string) => void | Promise<void>;
	onUpdateMediaSubtitle?: (
		_item: MediaItem,
		_subtitleId: string,
		_request: MediaItemSubtitleSelectionRequest
	) => void | Promise<void>;
	onRefreshMediaMetadata: (_item: MediaItem) => void;
	onSaveMediaItemOptions: (_item: MediaItem, _request: MediaItemUpdateRequest) => void;
	onDeleteMediaFile: (_item: MediaItem, _path: string) => void;
	onDeleteMediaFileTrack?: (
		_item: MediaItem,
		_request: MediaFileTrackDeleteRequest
	) => void | Promise<void>;
	onFulfillmentAction?: (
		_item: MediaItem,
		_request: MediaFulfillmentActionRequest
	) => void | Promise<void>;
	onAssembleMediaComponents?: (
		_item: MediaItem,
		_baseSourceId: string,
		_artifactIds: string[]
	) => void;
	onReviewComponentCompatibility?: (
		_item: MediaItem,
		_source: MediaComponentSource,
		_decisionId: string,
		_reviewState: MediaComponentCompatibilityReviewState
	) => void;
	onDeleteMedia: (_item: MediaItem) => void;
	onGrabRelease: (
		_item: MediaItem,
		_release: ReleaseCandidate,
		_overrideMatch?: boolean,
		_details?: ReleaseOverrideDetails
	) => void;
	onRefreshActivity: () => void;
	onRefreshReleaseBlocklist: () => void;
	onCancelActivity: (_activity: DownloadActivity) => void;
	onDeleteActivity: (_activity: DownloadActivity) => void;
	onDeleteReleaseBlocklistItem: (_item: ReleaseBlocklistItem) => void;
	onClearReleaseBlocklist: () => void;
}
