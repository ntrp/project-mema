import type { ActivityQueueStatus } from '$lib/components/app/activity/activityQueue';
import type { MediaFileRow } from '$lib/components/app/media/file-data/mediaFileRows';
import type {
	DownloadActivity,
	Language,
	LibraryFolder,
	MediaItem,
	MediaItemSubtitle,
	MediaItemUpdateRequest,
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
	subtitleSearching?: boolean;
	onAutoSearch: () => void;
	onManualSearch: () => void;
	onSearchSubtitle?: (_row: MediaFileRow, _languageId?: string) => void | Promise<void>;
	onDeleteSubtitle?: (_subtitle: MediaItemSubtitle) => void | Promise<void>;
	onDelete: (_row: MediaFileRow) => void;
}

export interface MediaFilesTableProps {
	item: MediaItem;
	activities: DownloadActivity[];
	searchingItemId?: string;
	grabbingKey?: string;
	canManage: boolean;
	libraryFolders: LibraryFolder[];
	languages: Language[];
	qualityProfiles: QualityProfileOption[];
	onSaveOptions: (_item: MediaItem, _request: MediaItemUpdateRequest) => void;
	onAutoSearch: (_item: MediaItem) => void;
	onSearchSubtitle: (_item: MediaItem, _request: SubtitleSearchRequest) => void | Promise<void>;
	onDeleteSubtitle: (_item: MediaItem, _subtitleId: string) => void | Promise<void>;
	onDeleteFile: (_item: MediaItem, _path: string) => void;
	onGrabRelease: (
		_item: MediaItem,
		_release: ReleaseCandidate,
		_overrideMatch?: boolean,
		_details?: ReleaseOverrideDetails
	) => void;
}
