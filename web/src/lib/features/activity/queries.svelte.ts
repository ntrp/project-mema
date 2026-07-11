import { createQuery } from '@tanstack/svelte-query';
import { listDownloadActivity, listReleaseBlocklist } from './api';

export const activityKeys = {
	all: ['activity'] as const,
	downloads: () => [...activityKeys.all, 'downloads'] as const,
	blocklist: () => [...activityKeys.all, 'blocklist'] as const
};

export function createDownloadActivityQuery(enabled: () => boolean = () => true) {
	return createQuery(() => ({
		queryKey: activityKeys.downloads(),
		queryFn: ({ signal }) => listDownloadActivity({ signal }),
		select: (response) => response.activities,
		enabled: enabled()
	}));
}

export function createReleaseBlocklistQuery(enabled: () => boolean = () => true) {
	return createQuery(() => ({
		queryKey: activityKeys.blocklist(),
		queryFn: ({ signal }) => listReleaseBlocklist({ signal }),
		select: (response) => response.items,
		enabled: enabled()
	}));
}
