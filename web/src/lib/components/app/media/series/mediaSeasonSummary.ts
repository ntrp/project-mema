import type { ActivityQueueStatus } from '$lib/components/app/activity/activityQueue';
import { formatBytes } from '$lib/components/app/media/file-data/mediaFileSize';
import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';

export interface SeasonEpisodeFile {
	row: MediaFileRow;
	activityStatus?: ActivityQueueStatus;
}

export interface SeasonFileSummary {
	label: string;
	size: string;
	tone: 'success' | 'active' | 'missing';
	hasActive: boolean;
}

export function seasonFileSummary(files: SeasonEpisodeFile[]): SeasonFileSummary {
	const total = files.length;
	const downloaded = files.filter((file) => file.row.exists).length;
	const queued = files.filter((file) => isActive(file.activityStatus)).length;
	const missing = files.some((file) => !file.row.exists && !isActive(file.activityStatus));
	return {
		label: queued > 0 ? `${downloaded} + ${queued} / ${total}` : `${downloaded} / ${total}`,
		size: seasonSize(files),
		tone: queued > 0 ? 'active' : missing ? 'missing' : 'success',
		hasActive: queued > 0
	};
}

function seasonSize(files: SeasonEpisodeFile[]) {
	const bytes = files.reduce((sum, file) => sum + (file.row.sizeBytes ?? 0), 0);
	return bytes > 0 ? formatBytes(bytes) : '-';
}

function isActive(status?: ActivityQueueStatus) {
	return (
		status?.status === 'queued' || status?.status === 'grabbed' || status?.status === 'downloading'
	);
}
