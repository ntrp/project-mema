import type {
	DownloadActivity,
	Language,
	LibraryFolder,
	MediaItem,
	MediaItemUpdateRequest,
	MediaMetadataDetails,
	MediaSearchResult,
	MediaType,
	QualityProfileOption,
	ReleaseCandidate,
	ReleaseOverrideDetails
} from '$lib/settings/types';

export interface MediaDetailProps {
	mediaType: MediaType;
	item?: MediaItem;
	loading?: boolean;
	mediaItems?: MediaItem[];
	libraryFolders: LibraryFolder[];
	languages: Language[];
	qualityProfiles: QualityProfileOption[];
	requestedItemId: string;
	activities: DownloadActivity[];
	searchingItemId?: string;
	refreshingMetadataItemId?: string;
	savingMediaItemOptionsId?: string;
	grabbingKey?: string;
	addingKey?: string;
	deletingMediaItemId?: string;
	canManage: boolean;
	actionLabel: string;
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
	onAddMedia: (_candidate: MediaSearchResult) => void;
}

export function imageUrl(path?: string, size = 'w780') {
	if (!path) return undefined;
	if (path.startsWith('http://') || path.startsWith('https://')) return path;
	return `https://image.tmdb.org/t/p/${size}${path}`;
}

export function mediaMetadataDetail(mediaItem: MediaItem): MediaMetadataDetails {
	return {
		title: mediaItem.title,
		type: mediaItem.type,
		year: mediaItem.year,
		monitored: mediaItem.monitored,
		seriesType: mediaItem.seriesType,
		externalProvider: mediaItem.externalProvider ?? 'local',
		externalId: mediaItem.externalId ?? mediaItem.id,
		overview: mediaItem.overview,
		posterPath: mediaItem.posterPath,
		collectionId: mediaItem.collectionId,
		collectionName: mediaItem.collectionName,
		backdropPath: mediaItem.backdropPath,
		status: mediaItem.metadataStatus,
		originalLanguage: mediaItem.originalLanguage,
		releaseDate: mediaItem.releaseDate,
		firstAirDate: mediaItem.firstAirDate,
		runtimeMinutes: mediaItem.runtimeMinutes,
		seasonCount: mediaItem.seasonCount,
		episodeCount: mediaItem.episodeCount,
		voteAverage: mediaItem.voteAverage,
		genres: mediaItem.genres,
		keywords: mediaItem.keywords,
		facts: mediaItem.facts,
		seasons: mediaItem.seasons,
		cast: mediaItem.cast,
		crew: mediaItem.crew,
		recommendations: mediaItem.recommendations,
		similar: mediaItem.similar
	};
}
