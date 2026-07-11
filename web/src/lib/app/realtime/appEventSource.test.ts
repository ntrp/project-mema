import { afterEach, describe, expect, it, vi } from 'vitest';

import {
	hasActiveAppEventSource,
	startAppEventSource,
	stopAppEventSource,
	subscribeToAppEvent
} from './appEventSource';

describe('app event source', () => {
	afterEach(() => {
		stopAppEventSource();
		vi.unstubAllGlobals();
	});

	it('creates only one source and routes all subscribed event types', () => {
		const sources = installEventSource();
		const activity = vi.fn();
		const notification = vi.fn();
		const unsubscribeActivity = subscribeToAppEvent('activity.updated', activity);
		const unsubscribeNotification = subscribeToAppEvent('system.event.created', notification);

		startAppEventSource();
		startAppEventSource();
		sources[0].emit('activity.updated', { id: '1', data: { status: 'queued' } });
		sources[0].emit('system.event.created', { id: '2', data: { severity: 'error' } });

		expect(sources).toHaveLength(1);
		expect(hasActiveAppEventSource()).toBe(true);
		expect(activity).toHaveBeenCalledWith({ id: '1', data: { status: 'queued' } });
		expect(notification).toHaveBeenCalledWith({ id: '2', data: { severity: 'error' } });
		unsubscribeActivity();
		unsubscribeNotification();
	});

	it('ignores duplicate and malformed event envelopes', () => {
		const sources = installEventSource();
		const handler = vi.fn();
		const unsubscribe = subscribeToAppEvent('system.event.created', handler);
		startAppEventSource();

		sources[0].emit('system.event.created', { id: 'same', data: { message: 'once' } });
		sources[0].emit('system.event.created', { id: 'same', data: { message: 'twice' } });
		sources[0].emitRaw('system.event.created', '{broken');

		expect(handler).toHaveBeenCalledTimes(1);
		unsubscribe();
	});

	it('does not duplicate native listeners after resubscription', () => {
		const sources = installEventSource();
		const first = vi.fn();
		const second = vi.fn();
		const unsubscribe = subscribeToAppEvent('system.event.created', first);
		startAppEventSource();
		unsubscribe();
		const unsubscribeSecond = subscribeToAppEvent('system.event.created', second);

		sources[0].emit('system.event.created', { data: { message: 'once' } });

		expect(first).not.toHaveBeenCalled();
		expect(second).toHaveBeenCalledTimes(1);
		unsubscribeSecond();
	});
});

function installEventSource() {
	const sources: FakeEventSource[] = [];
	vi.stubGlobal(
		'EventSource',
		class extends FakeEventSource {
			constructor(url: string, init: EventSourceInit) {
				super(url, init);
				sources.push(this);
			}
		}
	);
	return sources;
}

class FakeEventSource {
	closed = false;
	private listeners = new Map<string, ((event: MessageEvent<string>) => void)[]>();

	constructor(
		public url: string,
		public init: EventSourceInit
	) {}

	addEventListener(type: string, listener: (event: MessageEvent<string>) => void) {
		this.listeners.set(type, [...(this.listeners.get(type) ?? []), listener]);
	}

	close() {
		this.closed = true;
	}

	emit(type: string, envelope: unknown) {
		this.emitRaw(type, JSON.stringify(envelope));
	}

	emitRaw(type: string, data: string) {
		for (const listener of this.listeners.get(type) ?? []) {
			listener(new MessageEvent(type, { data }));
		}
	}
}
