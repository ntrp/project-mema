import type {
	DownloadActivity,
	DiscoverBlacklistItem,
	HomeSection,
	Language,
	LibraryFolder,
	MediaDiscoverSection,
	MediaItem,
	MediaItemUpdateRequest,
	MediaRequest,
	MediaRequestApproveRequest,
	MediaSearchResult,
	QualityProfileOption,
	ReleaseCandidate,
	ReleaseSearchResults
} from '$lib/settings/types';

export interface HomeAreaProps {
	activeSection: HomeSection;
	selectedMediaItemId?: string;
	selectedRequestId?: string;
	mediaItems: MediaItem[];
	mediaRequests: MediaRequest[];
	discoverSections: MediaDiscoverSection[];
	discoverBlacklist: DiscoverBlacklistItem[];
	libraryFolders: LibraryFolder[];
	languages: Language[];
	qualityProfiles: QualityProfileOption[];
	releaseResults: ReleaseSearchResults;
	activities: DownloadActivity[];
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
	canManage: boolean;
	loadingActivity: boolean;
	onAddMedia: (_candidate: MediaSearchResult) => void;
	onBlacklistMedia: (_candidate: MediaSearchResult) => void;
	onRemoveBlacklistMedia: (_item: DiscoverBlacklistItem) => void;
	onApproveMediaRequest: (_request: MediaRequest, _approval: MediaRequestApproveRequest) => void;
	onFindReleases: (_item: MediaItem, _query?: string) => void;
	onAutoSearchMedia: (_item: MediaItem) => void;
	onRefreshMediaMetadata: (_item: MediaItem) => void;
	onSaveMediaItemOptions: (_item: MediaItem, _request: MediaItemUpdateRequest) => void;
	onDeleteMediaFile: (_item: MediaItem, _path: string) => void;
	onDeleteMedia: (_item: MediaItem) => void;
	onGrabRelease: (_item: MediaItem, _release: ReleaseCandidate) => void;
	onRefreshActivity: () => void;
	onCancelActivity: (_activity: DownloadActivity) => void;
	onDeleteActivity: (_activity: DownloadActivity) => void;
}
