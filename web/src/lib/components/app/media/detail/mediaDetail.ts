import type { MediaItem, MediaMetadataDetails } from '$lib/settings/types';

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
		recommendations: mediaItem.recommendations,
		similar: mediaItem.similar
	};
}
