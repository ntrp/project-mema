import { afterEach, describe, expect, it, vi } from 'vitest';

import { connectAppEvents, disconnectAppEvents, type EventConnectionDeps } from './eventConnection';
import type { AppShellState } from './state.svelte';

type EventSourceOptions = { withCredentials?: boolean };

describe('app shell event connection (SCN-SYSTEM-008)', () => {
	afterEach(() => {
		vi.unstubAllGlobals();
	});

	it('does not connect for anonymous sessions or duplicate existing streams', () => {
		const sources = installEventSource();
		const state = { authenticated: false } as AppShellState;

		connectAppEvents(state, deps());
		expect(sources).toHaveLength(0);

		state.authenticated = true;
		connectAppEvents(state, deps());
		connectAppEvents(state, deps());
		expect(sources).toHaveLength(1);
		disconnectAppEvents(state);
	});

	it('routes activity and cache events through event actions', () => {
		const sources = installEventSource();
		const state = { authenticated: true } as AppShellState;
		const dependencies = deps();

		connectAppEvents(state, dependencies);
		const source = sources[0];
		expect(source.url).toBe('/api/events');
		expect(source.init).toEqual({ withCredentials: true });

		source.emit('activity.download.updated', { data: { id: 'activity-1', status: 'completed' } });
		source.emit('indexer.search.history.created', { data: { id: 'history-1' } });
		source.emit('indexer.search.cache.updated', { data: { entry: { query: 'movie' }, stats: {} } });
		source.emit('metadata.cache.updated', { data: { entry: { query: 'movie' }, stats: {} } });
		source.emit('metadata.search.history.created', { data: { id: 'metadata-history-1' } });
		source.emit('system.event.created', { data: { category: 'media' } });
		source.emit('system.job.execution.updated', {
			data: { kind: 'media.fulfillment.audio_transcode', status: 'completed' }
		});

		expect(dependencies.upsertActivity).toHaveBeenCalledWith({
			id: 'activity-1',
			status: 'completed'
		});
		expect(dependencies.updateMediaStatusFromActivity).toHaveBeenCalled();
		expect(dependencies.loadMediaItems).toHaveBeenCalledTimes(3);
		expect(dependencies.updateFulfillmentJobExecution).toHaveBeenCalledWith({
			kind: 'media.fulfillment.audio_transcode',
			status: 'completed'
		});
		expect(dependencies.appendIndexerSearchHistory).toHaveBeenCalledWith({ id: 'history-1' });
		expect(dependencies.upsertIndexerSearchCache).toHaveBeenCalledWith({
			entry: { query: 'movie' },
			stats: {}
		});
		expect(dependencies.upsertMetadataCache).toHaveBeenCalled();
		expect(dependencies.appendMetadataSearchHistory).toHaveBeenCalled();
		disconnectAppEvents(state);
	});

	it('disconnects the shared stream and removes subscriptions', () => {
		const sources = installEventSource();
		const state = { authenticated: true } as AppShellState;
		const dependencies = deps();
		connectAppEvents(state, dependencies);

		disconnectAppEvents(state);
		sources[0].emit('activity.download.updated', { data: { id: 'activity-1' } });

		expect(sources[0].closed).toBe(true);
		expect(dependencies.upsertActivity).not.toHaveBeenCalled();
	});
});

function deps(): EventConnectionDeps {
	return {
		loadMediaItems: vi.fn(),
		upsertActivity: vi.fn(),
		updateMediaStatusFromActivity: vi.fn(),
		appendIndexerSearchHistory: vi.fn(),
		upsertIndexerSearchCache: vi.fn(),
		upsertMetadataCache: vi.fn(),
		appendMetadataSearchHistory: vi.fn(),
		updateFulfillmentJobExecution: vi.fn(),
		parseEventData: vi.fn((event: Event) => JSON.parse((event as MessageEvent<string>).data).data)
	};
}

function installEventSource() {
	const sources: FakeEventSource[] = [];
	vi.stubGlobal(
		'EventSource',
		class extends FakeEventSource {
			constructor(url: string, init: EventSourceOptions) {
				super(url, init);
				sources.push(this);
			}
		}
	);
	return sources;
}

class FakeEventSource {
	closed = false;
	onerror?: (event: Event) => void;
	private listeners = new Map<string, ((event: MessageEvent<string>) => void)[]>();

	constructor(
		public url: string,
		public init: EventSourceOptions
	) {}

	addEventListener(type: string, listener: (event: MessageEvent<string>) => void) {
		this.listeners.set(type, [...(this.listeners.get(type) ?? []), listener]);
	}

	close() {
		this.closed = true;
	}

	emit(type: string, envelope: unknown) {
		for (const listener of this.listeners.get(type) ?? []) {
			listener(new MessageEvent(type, { data: JSON.stringify(envelope) }));
		}
	}
}
