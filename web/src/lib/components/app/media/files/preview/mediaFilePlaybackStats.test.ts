import { afterEach, describe, expect, it, vi } from 'vitest';

import { emptyPlaybackStats, watchPlaybackStats } from './mediaFilePlaybackStats';

describe('media file playback stats', () => {
	afterEach(() => {
		vi.useRealTimers();
		vi.unstubAllGlobals();
	});

	it('reports playback, bitrate variation, subtitles, and cleans up listeners', () => {
		vi.useFakeTimers();
		let stamp = 0;
		vi.stubGlobal('window', {
			location: { href: 'http://localhost/' },
			setInterval: globalThis.setInterval,
			clearInterval: globalThis.clearInterval
		});
		vi.stubGlobal('performance', {
			now: () => (stamp += 1000),
			getEntriesByName: () => []
		});
		const video = new FakeVideo();
		const trackListeners = new Set<() => void>();
		const tracks = Object.assign(
			[{ kind: 'subtitles', mode: 'showing', label: 'English', language: 'en' }],
			{
				on: (_type: string, handler: () => void) => trackListeners.add(handler),
				off: (_type: string, handler: () => void) => trackListeners.delete(handler)
			}
		);
		const player = { textTracks: () => tracks } as never;
		const changes: unknown[] = [];
		const stop = watchPlaybackStats(player, video as never, '/video.mkv', (stats) =>
			changes.push(stats)
		);

		for (const bytes of [1000, 2000, 4000, 8000]) {
			video.webkitVideoDecodedByteCount = bytes;
			vi.advanceTimersByTime(1000);
		}
		expect(changes[0]).toEqual(emptyPlaybackStats());
		expect(changes.at(-1)).toMatchObject({
			playing: true,
			variableBitRate: true,
			activeSubtitleLabel: 'English / en'
		});

		stop();
		expect(trackListeners.size).toBe(0);
		expect(video.listenerCount()).toBe(0);
	});

	it('uses resource timing bytes and handles paused media without subtitles', () => {
		vi.useFakeTimers();
		let stamp = 0;
		vi.stubGlobal('window', {
			location: { href: 'http://localhost/' },
			setInterval: globalThis.setInterval,
			clearInterval: globalThis.clearInterval
		});
		vi.stubGlobal('performance', {
			now: () => (stamp += 1000),
			getEntriesByName: () => [{ encodedBodySize: 0, transferSize: 2048, decodedBodySize: 0 }]
		});
		const video = new FakeVideo();
		video.paused = true;
		const changes: Array<{ playing: boolean; activeSubtitleLabel?: string }> = [];
		const stop = watchPlaybackStats(
			{ textTracks: () => [{ kind: 'captions', mode: 'showing' }] } as never,
			video as never,
			'/video.mkv',
			(stats) => changes.push(stats)
		);
		expect(changes.at(-1)).toMatchObject({ playing: false });
		expect(changes.at(-1)?.activeSubtitleLabel).toBeUndefined();
		stop();
	});
});

class FakeVideo extends EventTarget {
	paused = false;
	ended = false;
	webkitVideoDecodedByteCount = 0;
	webkitAudioDecodedByteCount = 0;
	private listeners = new Set<Listener>();

	addEventListener(type: string, callback: Listener | null) {
		if (callback) this.listeners.add(callback);
		super.addEventListener(type, callback);
	}

	removeEventListener(type: string, callback: Listener | null) {
		if (callback) this.listeners.delete(callback);
		super.removeEventListener(type, callback);
	}

	listenerCount() {
		return this.listeners.size;
	}
}

type Listener = ((event: Event) => void) | { handleEvent: (_event: Event) => void };
