import { useQueryClient } from '@tanstack/svelte-query';
import { subscribeToAppEvent } from '$lib/app/realtime/appEventSource';
import type { DownloadActivity } from './api';
import { createActivityCache } from './cache';

export function connectActivityQueryEvents() {
	const client = useQueryClient();
	const cache = createActivityCache(client);
	return subscribeToAppEvent<DownloadActivity>('activity.download.updated', ({ data }) => {
		if (!data) return;
		cache.upsert(data);
	});
}
