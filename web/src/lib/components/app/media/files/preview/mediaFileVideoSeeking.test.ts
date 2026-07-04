import { describe, expect, it } from 'vitest';

import { addSourceTimeline } from '$lib/components/app/media/files/preview/mediaFileVideoSeeking';
import type Player from 'video.js/dist/types/player';

describe('media file video seeking', () => {
	it('maps restarted remux previews back onto the source timeline', () => {
		const requestedSeeks: number[] = [];
		const player = timelinePlayer(5, [[0, 20]]);

		const removeTimeline = addSourceTimeline(player, true, 120, 600, (timeSeconds) => {
			requestedSeeks.push(timeSeconds);
		});

		expect(player.duration()).toBe(600);
		expect(player.currentTime()).toBe(125);
		expect(player.buffered().start(0)).toBe(120);
		expect(player.buffered().end(0)).toBe(140);
		expect(player.seekable().start(0)).toBe(0);
		expect(player.seekable().end(0)).toBe(600);

		player.currentTime(130);
		expect(player.currentTime()).toBe(130);
		expect(requestedSeeks).toEqual([]);

		player.currentTime(300);
		expect(requestedSeeks).toEqual([300]);

		removeTimeline?.();
	});

	it('keeps source time visible when a restarted preview reports local zero', () => {
		const requestedSeeks: number[] = [];
		const player = timelinePlayer(0, [[0, 10]]);

		const removeTimeline = addSourceTimeline(player, true, 300, 600, (timeSeconds) => {
			requestedSeeks.push(timeSeconds);
		});

		player.currentTime(0);

		expect(player.currentTime()).toBe(300);
		expect(requestedSeeks).toEqual([]);

		removeTimeline?.();
	});

	it('restarts from the requested source time when seeking outside the active segment', () => {
		const requestedSeeks: number[] = [];
		const player = timelinePlayer(5, [[0, 10]]);

		const removeTimeline = addSourceTimeline(player, true, 300, 600, (timeSeconds) => {
			requestedSeeks.push(timeSeconds);
		});

		player.currentTime(120);
		player.currentTime(0);

		expect(requestedSeeks).toEqual([120, 0]);

		removeTimeline?.();
	});
});

function timelinePlayer(currentTime: number, bufferedRanges: [number, number][]) {
	let localTime = currentTime;
	let localDuration: number | undefined = 10;
	const listeners = new Map<string, () => void>();
	const player = {
		currentTime: (seconds?: number | string) => {
			if (seconds !== undefined) localTime = Number(seconds);
			return localTime;
		},
		duration: (seconds?: number) => {
			if (seconds !== undefined) localDuration = seconds;
			return localDuration;
		},
		buffered: () => timeRanges(bufferedRanges),
		seekable: () => timeRanges(bufferedRanges),
		one: (event: string, callback: () => void) => listeners.set(event, callback),
		off: (event: string) => listeners.delete(event),
		trigger: (event: string) => listeners.get(event)?.()
	};
	return player as unknown as Player;
}

function timeRanges(ranges: [number, number][]) {
	return {
		length: ranges.length,
		start: (index: number) => ranges[index][0],
		end: (index: number) => ranges[index][1]
	};
}
