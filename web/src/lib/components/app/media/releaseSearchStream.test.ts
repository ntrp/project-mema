import { afterEach, describe, expect, it, vi } from 'vitest';

import { subscribeReleaseSearchStream } from './releaseSearchStream';

type EventSourceOptions = { withCredentials?: boolean };

describe('release search stream (SCN-MEDIA-002)', () => {
	afterEach(() => {
		vi.unstubAllGlobals();
	});

	it('subscribes with encoded query and emits status then result', () => {
		const handlers = {
			onStatus: vi.fn(),
			onResult: vi.fn(),
			onError: vi.fn()
		};
		const sources = installEventSource();

		const unsubscribe = subscribeReleaseSearchStream('media-1', 'Movie 2026', handlers);
		const source = sources[0];
		expect(source.url).toBe('/api/media/items/media-1/release-searches/stream?query=Movie%202026');

		source.emit('media.release_search.status', { data: { message: 'Searching indexer' } });
		source.emit('media.release_search.result', {
			data: { releases: [{ id: 'rel-1' }], errors: [] }
		});

		expect(handlers.onStatus).toHaveBeenCalledWith({ message: 'Searching indexer' });
		expect(handlers.onResult).toHaveBeenCalledWith({ releases: [{ id: 'rel-1' }], errors: [] });
		expect(handlers.onError).not.toHaveBeenCalled();
		expect(source.closed).toBe(true);

		unsubscribe();
		expect(source.close).toHaveBeenCalledTimes(2);
	});

	it('reports stream errors and ignores malformed status payloads', () => {
		const handlers = {
			onStatus: vi.fn(),
			onResult: vi.fn(),
			onError: vi.fn()
		};
		const sources = installEventSource();

		subscribeReleaseSearchStream('media-1', '', handlers);
		const source = sources[0];
		expect(source.url).toBe('/api/media/items/media-1/release-searches/stream');

		source.emitRaw('media.release_search.status', '{');
		source.emit('media.release_search.error', {});
		expect(handlers.onStatus).not.toHaveBeenCalled();
		expect(handlers.onError).toHaveBeenCalledWith('Release search failed');

		source.closed = false;
		source.emitRaw('error', '');
		expect(handlers.onError).toHaveBeenLastCalledWith('Release search stream disconnected');
		expect(source.closed).toBe(true);
	});
});

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
	close = vi.fn(() => {
		this.closed = true;
	});
	private listeners = new Map<string, ((event: MessageEvent<string>) => void)[]>();

	constructor(
		public url: string,
		public init: EventSourceOptions
	) {}

	addEventListener(type: string, listener: (event: MessageEvent<string>) => void) {
		this.listeners.set(type, [...(this.listeners.get(type) ?? []), listener]);
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
