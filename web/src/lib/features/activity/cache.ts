import type { QueryClient } from '@tanstack/svelte-query';
import type { DownloadActivity } from './api';
import { activityKeys } from './queries.svelte';

interface ActivityResponse {
	activities: DownloadActivity[];
}

export function createActivityCache(client: QueryClient) {
	function upsert(activity: DownloadActivity) {
		client.setQueryData<ActivityResponse>(activityKeys.downloads(), (current) => ({
			activities: [
				activity,
				...(current?.activities ?? []).filter((item) => item.id !== activity.id)
			]
		}));
	}

	function removeForMedia(mediaItemId: string) {
		client.setQueryData<ActivityResponse>(activityKeys.downloads(), (current) => ({
			activities: (current?.activities ?? []).filter((item) => item.mediaItemId !== mediaItemId)
		}));
	}

	function clear() {
		client.removeQueries({ queryKey: activityKeys.all });
	}

	function refresh() {
		return client.invalidateQueries({ queryKey: activityKeys.downloads() });
	}

	return { upsert, removeForMedia, clear, refresh };
}
