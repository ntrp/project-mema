import { activityForEpisode } from '../activity/activityQueue';
import { episodeKey, missingRow, seasonNumberFromName, type MediaFileRow } from './mediaFiles';
import type { SeasonEpisodeFile } from './mediaSeasonSummary';
import type {
	DownloadActivity,
	MediaMetadataEpisode,
	MediaMetadataSeason
} from '$lib/settings/types';

export type SeasonEpisodeRow = SeasonEpisodeFile & {
	episode: MediaMetadataEpisode;
};

export function episodeTitle(episode: MediaMetadataEpisode) {
	return `${episode.episodeNumber} - ${episode.name}`;
}

export function seasonMonitored(season: MediaMetadataSeason) {
	return (season.episodes ?? []).some((episode) => episode.monitored) || season.monitored === true;
}

export function seasonEpisodeRows(
	season: MediaMetadataSeason,
	seasonIndex: number,
	mediaRows: MediaFileRow[],
	activities: DownloadActivity[],
	itemId: string
): SeasonEpisodeRow[] {
	return (season.episodes ?? []).map((episode) => {
		const row = episodeFileRow(season, seasonIndex, episode, mediaRows);
		return {
			episode,
			row,
			activityStatus: activityForEpisode(activities, itemId, row.seasonNumber, row.episodeNumber)
		};
	});
}

function episodeFileRow(
	season: MediaMetadataSeason,
	seasonIndex: number,
	episode: MediaMetadataEpisode,
	mediaRows: MediaFileRow[]
) {
	const seasonNumber = seasonNumberFromName(season.name) ?? seasonIndex + 1;
	return (
		mediaRows.find(
			(row) =>
				episodeKey(row.seasonNumber, row.episodeNumber) ===
				episodeKey(seasonNumber, episode.episodeNumber)
		) ??
		missingRow(
			`s${seasonNumber}e${episode.episodeNumber}`,
			episode.name,
			seasonNumber,
			episode.episodeNumber
		)
	);
}
