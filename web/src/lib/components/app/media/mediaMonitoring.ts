import type {
	MediaItem,
	MediaItemUpdateRequest,
	MediaMetadataEpisode,
	MediaMetadataSeason,
	MediaMonitorMode
} from '$lib/settings/types';

export function monitorUpdate(
	updates: Pick<MediaItemUpdateRequest, 'monitored' | 'monitorMode' | 'seasons'>
): MediaItemUpdateRequest {
	return updates;
}

export function titleMonitorActive(item: MediaItem) {
	return item.type === 'series' ? item.monitorMode === 'future_episodes' : item.monitored;
}

export function titleMonitorStatus(item: MediaItem) {
	return titleMonitorActive(item) ? 'Monitored' : 'Not monitored';
}

export function titleMonitorHint(item: MediaItem) {
	if (item.type === 'series') {
		return titleMonitorActive(item)
			? 'Click to stop monitoring future episodes'
			: 'Click to monitor future episodes';
	}
	return item.monitored ? 'Click to stop monitoring this movie' : 'Click to monitor this movie';
}

export function toggledMediaMonitor(
	item: MediaItem
): Pick<MediaItemUpdateRequest, 'monitored' | 'monitorMode'> {
	const nextMonitored = !titleMonitorActive(item);
	const manualMonitored = item.type === 'series' && Boolean(item.seasons?.some(isSeasonMonitored));
	return {
		monitored: nextMonitored || manualMonitored,
		monitorMode: nextMonitorMode(item.type, nextMonitored)
	};
}

export function toggledEpisodeMonitor(
	item: MediaItem,
	seasons: MediaMetadataSeason[],
	targetSeason: MediaMetadataSeason,
	targetEpisode: MediaMetadataEpisode
): Pick<MediaItemUpdateRequest, 'monitored' | 'monitorMode' | 'seasons'> {
	const nextSeasons = seasons.map((season) => {
		if (season !== targetSeason) return season;
		const episodes = (season.episodes ?? []).map((episode) =>
			episode === targetEpisode ? { ...episode, monitored: !episode.monitored } : episode
		);
		return { ...season, episodes, monitored: episodes.some((episode) => episode.monitored) };
	});
	return seriesMonitorUpdate(item, nextSeasons);
}

export function toggledSeasonMonitor(
	item: MediaItem,
	seasons: MediaMetadataSeason[],
	targetSeason: MediaMetadataSeason
): Pick<MediaItemUpdateRequest, 'monitored' | 'monitorMode' | 'seasons'> {
	const nextSeasons = seasons.map((season) => {
		if (season !== targetSeason) return season;
		const nextMonitored = !isSeasonMonitored(season);
		const episodes = (season.episodes ?? []).map((episode) => ({
			...episode,
			monitored: nextMonitored
		}));
		return { ...season, episodes, monitored: nextMonitored };
	});
	return seriesMonitorUpdate(item, nextSeasons);
}

function seriesMonitorUpdate(
	item: MediaItem,
	seasons: MediaMetadataSeason[]
): Pick<MediaItemUpdateRequest, 'monitored' | 'monitorMode' | 'seasons'> {
	const manuallyMonitored = seasons.some(isSeasonMonitored);
	const futureMonitored = item.monitorMode === 'future_episodes';
	return {
		monitored: futureMonitored || manuallyMonitored,
		monitorMode: futureMonitored ? 'future_episodes' : manuallyMonitored ? 'all_episodes' : 'none',
		seasons
	};
}

function nextMonitorMode(type: MediaItem['type'], monitored: boolean): MediaMonitorMode {
	if (!monitored) return 'none';
	return type === 'series' ? 'future_episodes' : 'only_media';
}

function isSeasonMonitored(season: MediaMetadataSeason) {
	return (season.episodes ?? []).some((episode) => episode.monitored) || season.monitored === true;
}
