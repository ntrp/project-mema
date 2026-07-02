import type {
	MediaItem,
	MediaItemUpdateRequest,
	MediaMetadataEpisode,
	MediaMetadataSeason
} from '$lib/settings/types';

export function optimisticMediaItem(item: MediaItem, request: MediaItemUpdateRequest): MediaItem {
	const seasons = optimisticSeasons(item.seasons, request);
	return {
		...item,
		qualityProfileId: request.qualityProfileId ?? item.qualityProfileId,
		minimumAvailability: request.minimumAvailability ?? item.minimumAvailability,
		libraryFolderId: request.libraryFolderId ?? item.libraryFolderId,
		monitored: request.monitored ?? item.monitored,
		monitorMode: request.monitorMode ?? item.monitorMode,
		seasons
	};
}

export function mediaUpdateMessage(
	item: MediaItem,
	nextItem: MediaItem,
	request: MediaItemUpdateRequest
) {
	if (isSeasonPatch(request) && item.type === 'series') {
		return seriesSeasonMessage(item.seasons ?? [], nextItem.seasons ?? []);
	}
	if (request.libraryFolderId) {
		return 'Media root updated';
	}
	if (request.monitored !== undefined || request.monitorMode !== undefined) {
		return titleMonitorMessage(item, nextItem);
	}
	return 'Media settings saved';
}

function optimisticSeasons(
	seasons: MediaMetadataSeason[] | undefined,
	request: MediaItemUpdateRequest
) {
	if (!seasons) return request.seasons;
	if (request.seasons) return request.seasons;
	if (!isSeasonPatch(request)) return seasons;
	return seasons.map((season) => optimisticSeason(season, request));
}

function optimisticSeason(season: MediaMetadataSeason, request: MediaItemUpdateRequest) {
	if (season.name !== request.monitorSeasonName) return season;
	if (request.seasonMonitored !== undefined) {
		const episodes = (season.episodes ?? []).map((episode) => ({
			...episode,
			monitored: request.seasonMonitored
		}));
		return { ...season, episodes, monitored: request.seasonMonitored };
	}
	const episodes = (season.episodes ?? []).map((episode) =>
		episode.episodeNumber === request.monitorEpisodeNumber
			? { ...episode, monitored: request.episodeMonitored }
			: episode
	);
	return { ...season, episodes, monitored: episodes.some((episode) => episode.monitored) };
}

function isSeasonPatch(request: MediaItemUpdateRequest) {
	return Boolean(
		request.monitorSeasonName &&
		(request.seasonMonitored !== undefined || request.episodeMonitored !== undefined)
	);
}

function titleMonitorMessage(item: MediaItem, nextItem: MediaItem) {
	const monitored = item.type === 'series' ? nextItem.monitorMode !== 'none' : nextItem.monitored;
	const label = item.type === 'series' ? 'Series' : 'Movie';
	return `${label} is now ${monitored ? 'monitored' : 'not monitored'}`;
}

function seriesSeasonMessage(
	currentSeasons: MediaMetadataSeason[],
	nextSeasons: MediaMetadataSeason[]
) {
	const changed = changedSeasons(currentSeasons, nextSeasons);
	if (changed.length !== 1) return 'Series monitoring updated';

	const [change] = changed;
	const episodes = changedEpisodes(change.current?.episodes ?? [], change.next.episodes ?? []);
	if (episodes.length === 1) {
		const episode = episodes[0].next;
		return `Episode "${episode.name}" is now ${episode.monitored ? 'monitored' : 'not monitored'}`;
	}
	return `Season "${change.next.name}" is now ${isSeasonMonitored(change.next) ? 'monitored' : 'not monitored'}`;
}

function changedSeasons(currentSeasons: MediaMetadataSeason[], nextSeasons: MediaMetadataSeason[]) {
	return nextSeasons
		.map((next, index) => ({ current: currentSeasons[index], next }))
		.filter(({ current, next }) => seasonMonitorKey(current) !== seasonMonitorKey(next));
}

function changedEpisodes(
	currentEpisodes: MediaMetadataEpisode[],
	nextEpisodes: MediaMetadataEpisode[]
) {
	return nextEpisodes
		.map((next, index) => ({ current: currentEpisodes[index], next }))
		.filter(({ current, next }) => current?.monitored !== next.monitored);
}

function seasonMonitorKey(season?: MediaMetadataSeason) {
	if (!season) return '';
	return [season.monitored, ...(season.episodes ?? []).map((episode) => episode.monitored)].join(
		':'
	);
}

function isSeasonMonitored(season: MediaMetadataSeason) {
	return (season.episodes ?? []).some((episode) => episode.monitored) || season.monitored === true;
}
