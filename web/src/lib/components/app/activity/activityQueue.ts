import type { DownloadActivity } from '$lib/settings/types';

export interface ActivityQueueStatus {
	status: DownloadActivity['status'];
	label: string;
	progress?: number;
	activity: DownloadActivity;
}

const activeStatuses = ['queued', 'grabbed', 'downloading'] as const;

export function activityForMovie(
	activities: DownloadActivity[],
	mediaItemId: string
): ActivityQueueStatus | undefined {
	return queueStatus(activities.find((activity) => activity.mediaItemId === mediaItemId));
}

export function activityForEpisode(
	activities: DownloadActivity[],
	mediaItemId: string,
	seasonNumber?: number,
	episodeNumber?: number
): ActivityQueueStatus | undefined {
	if (!seasonNumber || !episodeNumber) return undefined;
	return queueStatus(
		activities.find(
			(activity) =>
				activity.mediaItemId === mediaItemId &&
				episodeKeys(activity.releaseTitle).has(episodeKey(seasonNumber, episodeNumber))
		)
	);
}

export function hasActiveActivity(activity: DownloadActivity) {
	return activeStatuses.includes(activity.status as (typeof activeStatuses)[number]);
}

function queueStatus(activity: DownloadActivity | undefined): ActivityQueueStatus | undefined {
	if (!activity) return undefined;
	const progress = activity.status === 'completed' ? 100 : (activity.progressPercent ?? undefined);
	return {
		status: activity.status,
		label: statusLabel(activity.status, progress),
		progress,
		activity
	};
}

function statusLabel(status: DownloadActivity['status'], progress?: number) {
	if (typeof progress === 'number' && (status === 'downloading' || status === 'completed')) {
		return `${status} ${progress}%`;
	}
	return status;
}

function episodeKeys(title: string) {
	const keys = new Set<string>();
	for (const match of title.matchAll(/s(\d{1,2})e(\d{1,3})/gi)) {
		keys.add(episodeKey(Number(match[1]), Number(match[2])));
	}
	return keys;
}

function episodeKey(season: number, episode: number) {
	return `${season}:${episode}`;
}
