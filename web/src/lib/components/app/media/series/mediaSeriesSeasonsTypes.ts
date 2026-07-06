import type {
	DownloadActivity,
	Language,
	LibraryFolder,
	MediaItem,
	MediaItemUpdateRequest,
	QualityProfileOption,
	ReleaseCandidate,
	ReleaseOverrideDetails
} from '$lib/settings/types';

type SubtitleSearchRequest = { languageId?: string; filePath?: string };

export interface MediaSeriesSeasonsProps {
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
