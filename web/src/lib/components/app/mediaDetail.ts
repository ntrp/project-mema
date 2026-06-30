import type {
	LibraryFolder,
	MediaItem,
	MediaMetadataDetails,
	QualityProfileOption
} from '$lib/settings/types';

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
		facts: mediaItem.facts,
		seasons: mediaItem.seasons,
		cast: mediaItem.cast
	};
}

export function libraryFolderLabel(mediaItem: MediaItem | undefined, folders: LibraryFolder[]) {
	if (!mediaItem) return 'Not set';
	return (
		mediaItem.mediaFolderPath ??
		mediaItem.libraryFolderPath ??
		folders.find((folder) => folder.id === mediaItem.libraryFolderId)?.path ??
		'Not set'
	);
}

export function qualityProfileLabel(
	mediaItem: MediaItem | undefined,
	profiles: QualityProfileOption[]
) {
	if (!mediaItem) return 'Not set';
	return (
		mediaItem.qualityProfileName ??
		profiles.find((profile) => profile.id === mediaItem.qualityProfileId)?.name ??
		'Not set'
	);
}

export function monitorModeLabel(mediaItem: MediaItem | undefined) {
	if (!mediaItem?.monitored || mediaItem.monitorMode === 'none') return 'None';
	return mediaItem.monitorMode === 'collection' ? 'Entire collection' : 'This media only';
}
