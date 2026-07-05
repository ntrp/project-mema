import { describe, expect, it, vi } from 'vitest';

import MediaFileDetailsAccordion from '$lib/components/app/media/files/MediaFileDetailsAccordion.svelte';
import MediaFileSummary from '$lib/components/app/media/files/MediaFileSummary.svelte';
import { missingRow, type MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import { renderWithTooltip } from '$lib/components/rendered/renderHelpers';

describe('rendered media file details (SCN-MEDIA-004)', () => {
	it('renders missing file search actions without delete affordance', () => {
		const { body } = renderWithTooltip(MediaFileSummary, {
			mediaItemId: 'media-1',
			mediaTitle: 'Scenario Movie',
			row: missingRow('movie-missing', 'Scenario Movie'),
			canManage: true,
			searching: false,
			fileLabel: 'Movie file',
			missingLabel: 'No matched file for this movie',
			onAutoSearch: vi.fn(),
			onManualSearch: vi.fn(),
			onDelete: vi.fn()
		});

		expect(body).toContain('Missing file');
		expect(body).toContain('No matched file for this movie');
		expect(body).toContain('Subtitles');
		expect(body).toContain('Missing');
		expect(body).toContain('Automatic search');
		expect(body).toContain('Manual search');
		expect(body).not.toContain('Delete file');
	});

	it('renders media track details, chapters, and missing expected languages', () => {
		const { body } = renderWithTooltip(MediaFileDetailsAccordion, { row: detailedFileRow() });

		expect(body).toContain('Track Nr.');
		expect(body).toContain('Video');
		expect(body).toContain('h264');
		expect(body).toContain('1920x1080');
		expect(body).toContain('Audio');
		expect(body).toContain('English');
		expect(body).toContain('DTS');
		expect(body).toContain('Subtitle');
		expect(body).toContain('Signs');
		expect(body).toContain('Chapter');
		expect(body).toContain('Opening');
		expect(body).toContain('German');
		expect(body).toContain('Missing expected audio track');
	});

	it('renders satisfied subtitle state in file summaries', () => {
		const { body } = renderWithTooltip(MediaFileSummary, {
			mediaItemId: 'media-1',
			mediaTitle: 'Scenario Movie',
			row: {
				...detailedFileRow(),
				subtitleSatisfaction: {
					state: 'satisfied' as const,
					wantedLanguages: ['english'],
					matchedLanguages: ['english'],
					missingLanguages: []
				}
			},
			canManage: true,
			searching: false,
			fileLabel: 'Movie file',
			missingLabel: 'No matched file for this movie',
			onAutoSearch: vi.fn(),
			onManualSearch: vi.fn(),
			onDelete: vi.fn()
		});

		expect(body).toContain('Subtitles');
		expect(body).toContain('Satisfied');
		expect(body).toContain('English');
	});
});

function detailedFileRow(): MediaFileRow {
	return {
		key: 'file-1',
		path: '/library/Scenario Movie/Scenario.Movie.2026.1080p.mkv',
		relativePath: 'Scenario.Movie.2026.1080p.mkv',
		exists: true,
		videoCodec: 'h264',
		audioInfo: 'DTS',
		size: '5.00 GiB',
		sizeBytes: 5 * 1024 * 1024 * 1024,
		languages: 'English',
		quality: '1080p',
		formats: ['WEB-DL'],
		upgrade: { state: 'current', label: 'Current', reasons: ['At or above upgrade target'] },
		expectedLanguages: ['english', 'german'],
		expectedRequiredLanguages: ['german'],
		expectedSubtitleLanguages: ['english'],
		removeNonEnabledLanguages: true,
		removeNonEnabledSubtitleLanguages: true,
		score: 120,
		tracks: [
			{ type: 'video', index: 0, codec: 'h264', width: 1920, height: 1080 },
			{ type: 'audio', index: 1, codec: 'DTS', language: 'eng', channels: 6 },
			{ type: 'subtitle', index: 2, codec: 'SRT', language: 'eng', title: 'Signs' },
			{ type: 'subtitle', index: 3, codec: 'SRT', language: 'spa', title: 'Spanish' }
		],
		chapters: [{ index: 0, title: 'Opening', startTime: '00:00:00', endTime: '00:05:00' }]
	};
}
