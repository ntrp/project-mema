import type { ActivitySection, DownloadActivity } from '$lib/settings/types';

export interface ActivitySectionHeading {
	title: string;
	empty: string;
}

export function visibleInActivitySection(activity: DownloadActivity, section: ActivitySection) {
	if (section === 'history') {
		return activity.status === 'completed' || activity.status === 'cancelled';
	}
	if (section === 'blocklist') {
		return false;
	}
	return ['queued', 'grabbed', 'downloading', 'failed'].includes(activity.status);
}

export function activitySectionHeading(section: ActivitySection): ActivitySectionHeading {
	if (section === 'history') {
		return {
			title: 'Activity history',
			empty: 'No completed background activity yet'
		};
	}
	if (section === 'blocklist') {
		return {
			title: 'Release blocklist',
			empty: 'No blocked releases yet'
		};
	}
	return {
		title: 'Activity queue',
		empty: 'No queued activity'
	};
}
