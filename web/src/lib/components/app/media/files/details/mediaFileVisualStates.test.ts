import { describe, expect, it } from 'vitest';

import {
	externalSubtitleVisualState,
	trackVisualState
} from '$lib/components/app/media/files/details/mediaFileVisualStates';
import { fileTrackDetailRows } from '$lib/components/app/media/files/mediaFileDetails';
import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';

describe('media file detail visual states', () => {
	it('marks audio candidates as matching, partial, or unwanted from profile targets', () => {
		const base = row({
			expectedAudioTargets: [
				{
					languageId: 'english',
					targetCodec: 'aac',
					targetChannels: ['2.0'],
					minimumBitrateKbps: 160
				}
			],
			expectedLanguages: ['english']
		});

		expect(
			trackVisualState(base, {
				type: 'audio',
				language: 'english',
				codec: 'aac',
				channels: 2,
				bitRate: '192000'
			})
		).toMatchObject({ visualState: 'matching' });
		expect(
			trackVisualState(base, {
				type: 'audio',
				language: 'english',
				codec: 'ac3',
				channels: 6,
				bitRate: '128000'
			})
		).toMatchObject({ visualState: 'partial' });
		expect(
			trackVisualState(
				{ ...base, removeNonEnabledLanguages: true },
				{ type: 'audio', language: 'spanish', codec: 'aac' }
			)
		).toMatchObject({ visualState: 'unwanted' });
		expect(trackVisualState(base, { type: 'audio', language: 'spanish', codec: 'aac' })).toEqual(
			{}
		);
	});

	it('changes subtitle row state when subtitle mode changes', () => {
		const base = row({ expectedSubtitleLanguages: ['english'] });

		expect(
			trackVisualState(
				{ ...base, subtitleSatisfaction: subtitleStatus('external') },
				{ type: 'subtitle', language: 'english', codec: 'srt' }
			)
		).toMatchObject({ visualState: 'pending_operation', operationLabel: 'Extract subtitle' });
		expect(
			trackVisualState(
				{ ...base, subtitleSatisfaction: subtitleStatus('mixed') },
				{ type: 'subtitle', language: 'english', codec: 'srt' }
			)
		).toMatchObject({ visualState: 'matching' });
		expect(
			externalSubtitleVisualState(
				{ ...base, subtitleSatisfaction: subtitleStatus('embedded') },
				'english'
			)
		).toMatchObject({ visualState: 'pending_operation', operationLabel: 'Embed subtitle' });
		expect(
			externalSubtitleVisualState(
				{ ...base, subtitleSatisfaction: subtitleStatus('external') },
				'english'
			)
		).toMatchObject({ visualState: 'matching' });
	});

	it('renders external subtitle rows and missing placeholders as detail rows', () => {
		const rows = fileTrackDetailRows(
			row({
				path: '/library/movie/Movie.mkv',
				expectedSubtitleLanguages: ['english', 'italian'],
				subtitleSatisfaction: {
					...subtitleStatus('external'),
					missingLanguages: ['italian']
				},
				externalSubtitles: [
					{
						id: 'sub-1',
						filePath: '/library/movie/Movie.eng.srt',
						providerName: 'Manual',
						languageId: 'english',
						format: 'srt',
						retentionMode: 'external',
						downloadedAt: '2026-07-08T00:00:00Z',
						selected: true
					}
				]
			})
		);

		expect(rows.map((detail) => [detail.key, detail.visualState])).toContainEqual([
			'external-subtitle-sub-1',
			'matching'
		]);
		expect(rows.map((detail) => [detail.key, detail.visualState])).toContainEqual([
			'missing-subtitle-it',
			'missing_placeholder'
		]);
	});
});

function row(overrides: Partial<MediaFileRow>): MediaFileRow {
	return {
		key: 'row',
		path: '/library/movie/Movie.mkv',
		relativePath: 'Movie.mkv',
		exists: true,
		videoCodec: '-',
		audioInfo: '-',
		size: '-',
		languages: '-',
		quality: '-',
		formats: [],
		tracks: [],
		chapters: [],
		otherFiles: [],
		externalSubtitles: [],
		upgrade: { state: 'current', label: 'Current', reasons: [] },
		expectedAudioTargets: [],
		expectedLanguages: [],
		expectedRequiredLanguages: [],
		expectedSubtitleLanguages: [],
		removeNonEnabledLanguages: false,
		removeNonEnabledSubtitleLanguages: false,
		score: 0,
		...overrides
	};
}

function subtitleStatus(mode: 'embedded' | 'external' | 'mixed') {
	return {
		state: 'missing' as const,
		mode,
		wantedLanguages: ['english'],
		matchedLanguages: [],
		missingLanguages: []
	};
}
