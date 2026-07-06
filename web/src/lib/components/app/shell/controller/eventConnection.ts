import type { DownloadActivity, SystemEvent } from '$lib/settings/types';
import type { createEventActions } from './events';
import type { AppShellState } from './state.svelte';

type EventActions = ReturnType<typeof createEventActions>;

export interface EventConnectionDeps {
	loadMediaItems: () => Promise<void>;
	upsertActivity: EventActions['upsertActivity'];
	updateMediaStatusFromActivity: EventActions['updateMediaStatusFromActivity'];
	appendIndexerSearchHistory: EventActions['appendIndexerSearchHistory'];
	upsertIndexerSearchCache: EventActions['upsertIndexerSearchCache'];
	upsertMetadataCache: EventActions['upsertMetadataCache'];
	appendMetadataSearchHistory: EventActions['appendMetadataSearchHistory'];
	parseEventData: EventActions['parseEventData'];
}

export function connectAppEvents(state: AppShellState, deps: EventConnectionDeps) {
	if (!state.authenticated || state.eventSource) return;
	const source = new EventSource('/api/events', { withCredentials: true });
	state.eventSource = source;
	source.addEventListener('activity.download.updated', (event) => {
		const activity = deps.parseEventData<DownloadActivity>(event);
		if (!activity) return;
		deps.upsertActivity(activity);
		deps.updateMediaStatusFromActivity(activity);
		if (activity.status === 'completed') {
			void deps.loadMediaItems();
		}
	});
	source.addEventListener('system.event.created', (event) => {
		const systemEvent = deps.parseEventData<SystemEvent>(event);
		if (systemEvent?.category === 'subtitles' || systemEvent?.category === 'media') {
			void deps.loadMediaItems();
		}
	});
	source.addEventListener('indexer.search.history.created', (event) => {
		const entry = deps.parseEventData<Parameters<typeof deps.appendIndexerSearchHistory>[0]>(event);
		if (entry) deps.appendIndexerSearchHistory(entry);
	});
	source.addEventListener('indexer.search.cache.updated', (event) => {
		const update = deps.parseEventData<Parameters<typeof deps.upsertIndexerSearchCache>[0]>(event);
		if (update) deps.upsertIndexerSearchCache(update);
	});
	source.addEventListener('metadata.cache.updated', (event) => {
		const update = deps.parseEventData<Parameters<typeof deps.upsertMetadataCache>[0]>(event);
		if (update) deps.upsertMetadataCache(update);
	});
	source.addEventListener('metadata.search.history.created', (event) => {
		const entry =
			deps.parseEventData<Parameters<typeof deps.appendMetadataSearchHistory>[0]>(event);
		if (entry) deps.appendMetadataSearchHistory(entry);
	});
	source.onerror = () => {
		if (!state.authenticated) {
			disconnectAppEvents(state);
		}
	};
}

export function disconnectAppEvents(state: AppShellState) {
	state.eventSource?.close();
	state.eventSource = undefined;
}
