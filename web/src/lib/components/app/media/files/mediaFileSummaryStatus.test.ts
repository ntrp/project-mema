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

	it('uses only audio in the fallback ok detail', () => {
		const status = audioSatisfaction(row({ audioLanguages: ['english'] }));

		expect(status.label).toBe('Ok');
		expect(status.details).toEqual(['Audio is available']);
	});

	it('marks audio ok when the target codec, channels, and bitrate match', () => {
		const status = audioSatisfaction(
			row({
				expectedAudioTargets: [
					{
						languageId: 'english',
						targetCodec: 'eac3',
						targetChannels: ['5.1'],
						minimumBitrateKbps: 640
					}
				],
				audioTracks: [
					{
						type: 'audio',
						index: 1,
						language: 'eng',
						codec: 'DD+',
						channels: 6,
						bitRate: '768000'
					}
				]
			})
		);

		expect(status.label).toBe('Ok');
		expect(status.details).toEqual(['Audio requirements met: English']);
	});

	it('marks audio partial when a present language misses target details', () => {
		const status = audioSatisfaction(
			row({
				expectedAudioTargets: [
					{
						languageId: 'english',
						targetCodec: 'aac',
						targetChannels: ['5.1'],
						minimumBitrateKbps: 384
					}
				],
				audioTracks: [
					{
						type: 'audio',
						index: 1,
						language: 'english',
						codec: 'dts',
						channels: 2,
						bitRate: '192000'
					}
				]
			})
		);

		expect(status.label).toBe('Partial');
		expect(status.details).toEqual([
			'English audio codec dts != aac',
			'English audio channels 2.0 not in 5.1',
			'English audio bitrate 192 kbps below 384 kbps'
		]);
	});

	it('does not mark minimum bitrate ok when bitrate is unknown', () => {
		const status = audioSatisfaction(
			row({
				expectedAudioTargets: [{ languageId: 'english', minimumBitrateKbps: 384 }],
				audioTracks: [{ type: 'audio', index: 1, language: 'english', codec: 'aac' }]
			})
		);

		expect(status.label).toBe('Partial');
		expect(status.details).toEqual(['English audio bitrate unknown below 384 kbps']);
	});

	it('uses partial subtitle status when only some subtitles match', () => {
		const status = subtitleSatisfaction(
			row({
				subtitleSatisfaction: {
					state: 'missing',
					mode: 'mixed',
					wantedLanguages: ['english', 'italian'],
					matchedLanguages: ['english'],
					missingLanguages: ['italian']
				}
			})
		);

		expect(status.label).toBe('Partial');
		expect(status.details).toEqual(['Mode: Mixed', 'Missing subtitles: Italian']);
	});

	it('marks embedded subtitles partial when they are available externally and need import', () => {
		const status = subtitleSatisfaction(
			row({
				subtitleSatisfaction: {
					state: 'missing',
					mode: 'embedded',
					wantedLanguages: ['english'],
					matchedLanguages: [],
					missingLanguages: ['english']
				},
				otherFiles: [
					{
						type: 'subtitle',
						path: '/library/movie/movie.english.srt',
						status: 'available',
						language: 'english'
					}
				]
			})
		);

		expect(status.label).toBe('Partial');
		expect(status.details).toEqual(['Mode: Embedded', 'Subtitles need to be imported: English']);
	});
});

function row(
	overrides: Partial<MediaFileRow> & {
		audioLanguages?: string[];
		audioTracks?: MediaFileRow['tracks'];
	} = {}
): MediaFileRow {
	const audioLanguages = overrides.audioLanguages ?? ['english'];
	const audioTracks =
		overrides.audioTracks ??
		audioLanguages.map((language, index) => ({
			type: 'audio' as const,
			index: index + 1,
			language
		}));
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
		tracks: [{ type: 'video', index: 0 }, ...audioTracks],
		chapters: [],
		otherFiles: [],
		externalSubtitles: [],
		subtitleSatisfaction: {
			state: 'satisfied',
			mode: 'mixed',
			wantedLanguages: ['english'],
			matchedLanguages: ['english'],
			missingLanguages: []
		},
		upgrade: { state: 'current', label: 'Current', reasons: [] },
		expectedAudioTargets: [],
		expectedLanguages: overrides.expectedLanguages ?? overrides.expectedRequiredLanguages ?? [],
		expectedRequiredLanguages: [],
		expectedSubtitleLanguages: [],
		removeNonEnabledLanguages: false,
		removeNonEnabledSubtitleLanguages: false,
		score: 0,
		...overrides
	};
}
