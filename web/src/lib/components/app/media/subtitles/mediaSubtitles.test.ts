import { describe, expect, it } from 'vitest';

import { subtitleStateRows } from './mediaSubtitles';
import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';

describe('media subtitle state mapping', () => {
	it('maps embedded, external, missing, and in-progress states', () => {
		const row = fileRow();

		expect(subtitleStateRows(row).map((item) => [item.language, item.label])).toEqual([
			['English', 'Embedded'],
			['German', 'External'],
			['French', 'Missing']
		]);

		expect(subtitleStateRows(row, true).map((item) => [item.language, item.label])).toEqual([
			['English', 'Embedded'],
			['German', 'External'],
			['French', 'Downloading']
		]);
	});
});

function fileRow(): MediaFileRow {
	return {
		key: 'file-1',
		path: '/library/Scenario/Scenario.Movie.2026.mkv',
		relativePath: 'Scenario.Movie.2026.mkv',
		exists: true,
		videoCodec: 'h264',
		audioInfo: 'DTS',
		size: '5.00 GiB',
		sizeBytes: 5,
		languages: 'English',
		quality: '1080p',
		formats: [],
		tracks: [{ type: 'subtitle', index: 2, codec: 'SRT', language: 'eng' }],
		chapters: [],
		subtitleSatisfaction: {
			state: 'missing',
			wantedLanguages: ['english', 'german', 'french'],
			matchedLanguages: ['english', 'german'],
			missingLanguages: ['french']
		},
		externalSubtitles: [
			{
				id: 'sub-1',
				providerName: 'OpenSubtitles',
				languageId: 'german',
				format: 'srt',
				filePath: '/library/Scenario/Scenario.Movie.2026.german.srt',
				downloadedAt: '2026-07-06T10:00:00Z'
			}
		],
		upgrade: { state: 'current', label: 'Current', reasons: [] },
		expectedLanguages: [],
		expectedRequiredLanguages: [],
		expectedSubtitleLanguages: ['english', 'german', 'french'],
		removeNonEnabledLanguages: false,
		removeNonEnabledSubtitleLanguages: false,
		score: 0
	};
}
