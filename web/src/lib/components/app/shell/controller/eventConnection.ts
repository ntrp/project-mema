import type { DownloadActivity, SystemEvent, SystemJobExecution } from '$lib/settings/types';
import {
	startAppEventSource,
	stopAppEventSource,
	subscribeToAppEvent
} from '$lib/app/realtime/appEventSource';
import type { createEventActions } from './events';
import type { AppShellState } from './state.svelte';

type EventActions = ReturnType<typeof createEventActions>;

export interface EventConnectionDeps {
	loadMediaItems: () => Promise<void>;
	upsertActivity: (_activity: DownloadActivity) => void;
	updateMediaStatusFromActivity: (_activity: DownloadActivity) => void;
	appendIndexerSearchHistory: EventActions['appendIndexerSearchHistory'];
	upsertIndexerSearchCache: EventActions['upsertIndexerSearchCache'];
	upsertMetadataCache: EventActions['upsertMetadataCache'];
	appendMetadataSearchHistory: EventActions['appendMetadataSearchHistory'];
	updateFulfillmentJobExecution: (_execution: SystemJobExecution) => void;
	parseEventData: EventActions['parseEventData'];
}

export function connectAppEvents(state: AppShellState, deps: EventConnectionDeps) {
	if (!state.authenticated || connections.has(state)) return;
	const unsubscribers = [
		subscribeToAppEvent<DownloadActivity>('activity.download.updated', ({ data: activity }) => {
			if (!activity) return;
			deps.upsertActivity(activity);
			deps.updateMediaStatusFromActivity(activity);
			if (activity.status === 'completed') void deps.loadMediaItems();
		}),
		subscribeToAppEvent<SystemEvent>('system.event.created', ({ data: systemEvent }) => {
			if (systemEvent?.category === 'subtitles' || systemEvent?.category === 'media') {
				void deps.loadMediaItems();
			}
		}),
		subscribeToAppEvent<SystemJobExecution>(
			'system.job.execution.updated',
			({ data: execution }) => {
				if (!execution) return;
				deps.updateFulfillmentJobExecution(execution);
				if (
					execution.kind.startsWith('media.fulfillment.') &&
					(execution.status === 'completed' ||
						execution.status === 'cancelled' ||
						execution.status === 'discarded')
				) {
					void deps.loadMediaItems();
				}
			}
		),
		subscribeToAppEvent<Parameters<typeof deps.appendIndexerSearchHistory>[0]>(
			'indexer.search.history.created',
			({ data: entry }) => {
				if (entry) deps.appendIndexerSearchHistory(entry);
			}
		),
		subscribeToAppEvent<Parameters<typeof deps.upsertIndexerSearchCache>[0]>(
			'indexer.search.cache.updated',
			({ data: update }) => {
				if (update) deps.upsertIndexerSearchCache(update);
			}
		),
		subscribeToAppEvent<Parameters<typeof deps.upsertMetadataCache>[0]>(
			'metadata.cache.updated',
			({ data: update }) => {
				if (update) deps.upsertMetadataCache(update);
			}
		),
		subscribeToAppEvent<Parameters<typeof deps.appendMetadataSearchHistory>[0]>(
			'metadata.search.history.created',
			({ data: entry }) => {
				if (entry) deps.appendMetadataSearchHistory(entry);
			}
		)
	];
	connections.set(state, unsubscribers);
	startAppEventSource();
}

export function disconnectAppEvents(state: AppShellState) {
	for (const unsubscribe of connections.get(state) ?? []) unsubscribe();
	connections.delete(state);
	stopAppEventSource();
}

const connections = new WeakMap<AppShellState, (() => void)[]>();
