import { describe, expect, it } from 'vitest';

import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import {
	mediaFilePreviewInfoUrl,
	mediaFilePreviewUrl,
	mediaFileTextTracks
} from '$lib/components/app/media/files/preview/mediaFilePlayback';
import { addSourceTimeline } from '$lib/components/app/media/files/preview/mediaFileVideoSeeking';
import type Player from 'video.js/dist/types/player';

describe('media file playback helpers', () => {
	it('builds selected-audio preview URLs', () => {
		const url = mediaFilePreviewUrl('media 1', '/library/Movie File.mkv', 2);

		expect(url).toBe(
			'/api/media/items/media%201/files/preview?path=%2Flibrary%2FMovie+File.mkv&audioTrackIndex=2'
		);
	});

	it('builds selected-audio preview info URLs', () => {
		const url = mediaFilePreviewInfoUrl('media 1', '/library/Movie File.mkv', 2);

		expect(url).toBe(
			'/api/media/items/media%201/files/preview-info?path=%2Flibrary%2FMovie+File.mkv&audioTrackIndex=2'
		);
	});

	it('builds subtitle and chapter text tracks for Video.js', () => {
		const tracks = mediaFileTextTracks('media-1', playbackRow());

		expect(tracks[0]).toMatchObject({
			key: 'subtitle-4',
			kind: 'subtitles',
			label: 'Signs · English · srt · Track 4',
			src: '/api/media/items/media-1/files/subtitle?path=%2Flibrary%2FMovie.mkv&subtitleTrackIndex=4',
			srclang: 'en'
		});
		expect(tracks[1]).toMatchObject({
			key: 'chapters',
			kind: 'chapters',
			label: 'Chapters',
			default: true
		});
		expect(decodeURIComponent(tracks[1].src)).toContain(
			'WEBVTT\n\n00:00:00.000 --> 00:05:00.000\nOpening'
		);
	});

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

function playbackRow(): MediaFileRow {
	return {
		key: 'file-1',
		path: '/library/Movie.mkv',
		relativePath: 'Movie.mkv',
		exists: true,
		videoCodec: 'h264',
		audioInfo: 'AAC',
		size: '1.00 GiB',
		languages: 'English',
		quality: '1080p',
		formats: [],
		tracks: [
			{ type: 'audio', index: 2, codec: 'aac', language: 'eng' },
			{ type: 'subtitle', index: 4, codec: 'srt', language: 'eng', title: 'Signs' }
		],
		chapters: [
			{ index: 0, title: 'Opening', startTime: '0', endTime: '300' },
			{ index: 1, title: 'Middle', startTime: '00:05:00', endTime: '600' }
		],
		expectedLanguages: [],
		removeNonEnabledLanguages: false,
		score: 0
	};
}

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
