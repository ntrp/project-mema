import type { ActivityQueueStatus } from '$lib/components/app/activity/activityQueue';
import type { MediaFileRow } from '$lib/components/app/media/file-data/mediaFileRows';
import type {
	DownloadActivity,
	Language,
	LibraryFolder,
	MediaFileTrackDeleteRequest,
	MediaFulfillmentActionRequest,
	MediaItem,
	MediaItemSubtitle,
	MediaItemSubtitleSelectionRequest,
	MediaItemUpdateRequest,
	GrabSubtitleRequest,
	QualityProfileOption,
	ReleaseCandidate,
	ReleaseOverrideDetails
} from '$lib/settings/types';

type SubtitleSearchRequest = { languageId?: string; filePath?: string };

export interface MediaFileSummaryProps {
	mediaItemId: string;
	mediaTitle: string;
	row: MediaFileRow;
	activityStatus?: ActivityQueueStatus;
	canManage: boolean;
	searching: boolean;
	fileLabel?: string;
	missingLabel?: string;
	showSearchActions?: boolean;
	onAutoSearch: () => void;
	onManualSearch: () => void;
	onSearchSubtitle?: (_row: MediaFileRow, _languageId?: string) => void | Promise<void>;
	onManualSubtitleSearch?: (_row: MediaFileRow, _languageId?: string) => void;
	onDeleteSubtitle?: (_subtitle: MediaItemSubtitle) => void | Promise<void>;
	onUpdateSubtitle?: (
		_subtitle: MediaItemSubtitle,
		_request: MediaItemSubtitleSelectionRequest
	) => void | Promise<void>;
	onDeleteTrack?: (
		_row: MediaFileRow,
		_request: MediaFileTrackDeleteRequest
	) => void | Promise<void>;
	onFulfillmentAction?: (
		_row: MediaFileRow,
		_request: MediaFulfillmentActionRequest
	) => void | Promise<void>;
	onDelete: (_row: MediaFileRow) => void;
}

export interface MediaFilesTableProps {
	item: MediaItem;
	activities: DownloadActivity[];
	searchingItemId?: string;
	scanningMediaItemId?: string;
	grabbingKey?: string;
	canManage: boolean;
	libraryFolders: LibraryFolder[];
	languages: Language[];
	qualityProfiles: QualityProfileOption[];
	onSaveOptions: (_item: MediaItem, _request: MediaItemUpdateRequest) => void;
	onAutoSearch: (_item: MediaItem) => void;
	onRescanMediaFiles: (_item: MediaItem) => void;
	onSearchSubtitle: (_item: MediaItem, _request: SubtitleSearchRequest) => void | Promise<void>;
	onGrabSubtitle: (_item: MediaItem, _request: GrabSubtitleRequest) => void | Promise<void>;
	onDeleteSubtitle: (_item: MediaItem, _subtitleId: string) => void | Promise<void>;
	onUpdateSubtitle: (
		_item: MediaItem,
		_subtitleId: string,
		_request: MediaItemSubtitleSelectionRequest
	) => void | Promise<void>;
	onDeleteFile: (_item: MediaItem, _path: string) => void;
	onDeleteFileTrack: (
		_item: MediaItem,
		_request: MediaFileTrackDeleteRequest
	) => void | Promise<void>;
	onFulfillmentAction?: (
		_item: MediaItem,
		_request: MediaFulfillmentActionRequest
	) => void | Promise<void>;
	onGrabRelease: (
		_item: MediaItem,
		_release: ReleaseCandidate,
		_overrideMatch?: boolean,
		_details?: ReleaseOverrideDetails
	) => void;
}
