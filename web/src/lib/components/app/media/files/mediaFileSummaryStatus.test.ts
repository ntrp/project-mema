import { describe, expect, it } from 'vitest';

import { audioSatisfaction, subtitleSatisfaction } from './mediaFileSummaryStatus';
import type { MediaFileRow } from './mediaFiles';

describe('media file summary requirement status', () => {
	it('marks audio missing when no required language is present', () => {
		const status = audioSatisfaction(
			row({ expectedRequiredLanguages: ['italian'], audioLanguages: ['english'] })
		);

		expect(status.label).toBe('Missing');
		expect(status.details).toEqual(['Missing required audio: Italian']);
	});

	it('marks audio partial when some required languages are present', () => {
		const status = audioSatisfaction(
			row({
				expectedRequiredLanguages: ['english', 'italian'],
				audioLanguages: ['english']
			})
		);

		expect(status.label).toBe('Partial');
		expect(status.details).toEqual(['Missing required audio: Italian']);
	});

	it('marks audio partial when unwanted tracks remain', () => {
		const status = audioSatisfaction(
			row({
				expectedLanguages: ['italian'],
				expectedRequiredLanguages: ['italian'],
				audioLanguages: ['italian', 'english'],
				removeNonEnabledLanguages: true
			})
		);

		expect(status.label).toBe('Partial');
		expect(status.details).toEqual(['Unwanted audio tracks: English']);
	});

	it('marks audio ok when requirements are satisfied', () => {
		const status = audioSatisfaction(
			row({ expectedRequiredLanguages: ['italian'], audioLanguages: ['italian'] })
		);

		expect(status.label).toBe('Ok');
		expect(status.details).toEqual(['Required audio present: Italian']);
	});

	it('uses partial subtitle status when only some subtitles match', () => {
		const status = subtitleSatisfaction(
			row({
				subtitleSatisfaction: {
					state: 'missing',
					preferredMode: 'mixed',
					wantedLanguages: ['english', 'italian'],
					matchedLanguages: ['english'],
					missingLanguages: ['italian']
				}
			})
		);

		expect(status.label).toBe('Partial');
		expect(status.details).toEqual(['Missing subtitles: Italian']);
	});
});

function row(overrides: Partial<MediaFileRow> & { audioLanguages?: string[] } = {}): MediaFileRow {
	const audioLanguages = overrides.audioLanguages ?? ['english'];
	return {
		key: 'file-1',
		path: '/library/movie.mkv',
		relativePath: 'movie.mkv',
		exists: true,
		videoCodec: 'h264',
		audioInfo: 'aac',
		size: '1 GiB',
		languages: 'English',
		quality: '1080p',
		formats: [],
		tracks: [
			{ type: 'video', index: 0 },
			...audioLanguages.map((language, index) => ({
				type: 'audio' as const,
				index: index + 1,
				language
			}))
		],
		chapters: [],
		otherFiles: [],
		externalSubtitles: [],
		subtitleSatisfaction: {
			state: 'satisfied',
			preferredMode: 'mixed',
			wantedLanguages: ['english'],
			matchedLanguages: ['english'],
			missingLanguages: []
		},
		upgrade: { state: 'current', label: 'Current', reasons: [] },
		expectedLanguages: overrides.expectedLanguages ?? overrides.expectedRequiredLanguages ?? [],
		expectedRequiredLanguages: [],
		expectedSubtitleLanguages: [],
		removeNonEnabledLanguages: false,
		removeNonEnabledSubtitleLanguages: false,
		score: 0,
		...overrides
	};
}
