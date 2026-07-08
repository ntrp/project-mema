import { activityForEpisode } from '$lib/components/app/activity/activityQueue';
import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import {
	episodeKey,
	missingRow,
	seasonNumberFromName
} from '$lib/components/app/media/files/mediaFileMissing';
import type { ReleaseSearchContext } from '$lib/components/app/media/release-search/releaseSearchQuery';
import type { SeasonEpisodeFile } from '$lib/components/app/media/series/mediaSeasonSummary';
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

export function seasonReleaseSearchContext(
	season: MediaMetadataSeason,
	index: number
): ReleaseSearchContext {
	return { type: 'season', seasonNumber: seasonNumberFromName(season.name) ?? index + 1 };
}

export function episodeReleaseSearchContext(row: MediaFileRow): ReleaseSearchContext {
	return { type: 'episode', seasonNumber: row.seasonNumber, episodeNumber: row.episodeNumber };
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
