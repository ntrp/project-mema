import { describe, expect, it } from 'vitest';

import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import {
	mediaFilePreviewClientProfile,
	mediaFilePreviewInfoUrl,
	mediaFilePreviewSourceType,
	mediaFilePreviewUrl,
	mediaFileTextTracks
} from '$lib/components/app/media/files/preview/mediaFilePlayback';

describe('media file playback helpers', () => {
	it('builds selected-audio preview URLs', () => {
		const url = mediaFilePreviewUrl('media 1', '/library/Movie File.mkv', 2);

		expect(url).toBe(
			'/api/media/items/media%201/files/preview?path=%2Flibrary%2FMovie+File.mkv&audioTrackIndex=2'
		);
	});

	it('adds client profile to preview URLs when needed', () => {
		const url = mediaFilePreviewUrl('media 1', '/library/Movie File.mkv', 2, 'webkit');

		expect(url).toBe(
			'/api/media/items/media%201/files/preview?path=%2Flibrary%2FMovie+File.mkv&audioTrackIndex=2&clientProfile=webkit'
		);
	});

	it('builds selected-audio preview info URLs', () => {
		const url = mediaFilePreviewInfoUrl('media 1', '/library/Movie File.mkv', 2);

		expect(url).toBe(
			'/api/media/items/media%201/files/preview-info?path=%2Flibrary%2FMovie+File.mkv&audioTrackIndex=2'
		);
	});

	it('adds client profile to preview info URLs when needed', () => {
		const url = mediaFilePreviewInfoUrl('media 1', '/library/Movie File.mkv', 2, 'webkit');

		expect(url).toBe(
			'/api/media/items/media%201/files/preview-info?path=%2Flibrary%2FMovie+File.mkv&audioTrackIndex=2&clientProfile=webkit'
		);
	});

	it('selects the video.js source type from preview delivery', () => {
		expect(mediaFilePreviewSourceType('hls')).toBe('application/x-mpegURL');
		expect(mediaFilePreviewSourceType('file')).toBe('video/mp4');
		expect(mediaFilePreviewSourceType()).toBe('video/mp4');
	});

	it('detects WebKit preview clients', () => {
		expect(
			mediaFilePreviewClientProfile({
				userAgent:
					'Mozilla/5.0 (iPad; CPU OS 17_5 like Mac OS X) AppleWebKit/605.1.15 Mobile/15E148 Safari/604.1'
			})
		).toBe('webkit');
		expect(
			mediaFilePreviewClientProfile({
				userAgent:
					'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 Version/17.5 Safari/605.1.15',
				platform: 'MacIntel',
				maxTouchPoints: 5
			})
		).toBe('webkit');
		expect(
			mediaFilePreviewClientProfile({
				userAgent:
					'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Chrome/126.0.0.0 Safari/537.36'
			})
		).toBeUndefined();
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
		otherFiles: [],
		missingTracks: [],
		upgrade: { state: 'current', label: 'Current', reasons: ['At or above upgrade target'] },
		score: 0
	};
}
