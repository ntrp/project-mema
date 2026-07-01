/* global MessageEvent */
import type { DownloadActivity, DownloadActivityStatus } from '$lib/settings/types';
import type { AppShellState } from './state.svelte';
import type { ServerEventEnvelope } from './types';

export function createEventActions(state: AppShellState) {
	function upsertActivity(activity: DownloadActivity) {
		state.activities = [activity, ...state.activities.filter((item) => item.id !== activity.id)];
	}

	function updateMediaStatusFromActivity(activity: DownloadActivity) {
		const status = mediaStatusFromActivity(activity.status);
		if (!status) {
			return;
		}
		state.mediaItems = state.mediaItems.map((item) =>
			item.id === activity.mediaItemId ? { ...item, status } : item
		);
	}

	function mediaStatusFromActivity(status: DownloadActivityStatus) {
		if (status === 'completed') {
			return 'downloaded';
		}
		if (status === 'queued' || status === 'grabbed' || status === 'downloading') {
			return 'downloading';
		}
		return undefined;
	}

	function parseEventData<T>(event: Event) {
		const message = event as MessageEvent<string>;
		try {
			return (JSON.parse(message.data) as ServerEventEnvelope<T>).data;
		} catch {
			return undefined;
		}
	}

	return { upsertActivity, updateMediaStatusFromActivity, parseEventData };
}
