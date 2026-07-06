import type {
	DownloadActivity,
	DiscoverBlacklistItem,
	HomeSection,
	ActivitySection,
	Language,
	LibraryFolder,
	MediaDiscoverSection,
	MediaItem,
	MediaItemUpdateRequest,
	MediaRequest,
	MediaRequestApproveRequest,
	MediaSearchResult,
	QualityProfileOption,
	ReleaseBlocklistItem,
	ReleaseCandidate,
	ReleaseOverrideDetails
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
	refreshingMetadataItemId?: string;
	savingMediaItemOptionsId?: string;
	grabbingKey?: string;
	deletingMediaItemId?: string;
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
	onSearchMediaSubtitle?: (
		_item: MediaItem,
		_request: { languageId?: string; filePath?: string }
	) => void | Promise<void>;
	onDeleteMediaSubtitle?: (_item: MediaItem, _subtitleId: string) => void | Promise<void>;
	onRefreshMediaMetadata: (_item: MediaItem) => void;
	onSaveMediaItemOptions: (_item: MediaItem, _request: MediaItemUpdateRequest) => void;
	onDeleteMediaFile: (_item: MediaItem, _path: string) => void;
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
