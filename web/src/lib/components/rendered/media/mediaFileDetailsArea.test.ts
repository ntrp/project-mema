import { describe, expect, it, vi } from 'vitest';

import MediaFileDetailsAccordion from '$lib/components/app/media/files/MediaFileDetailsAccordion.svelte';
import MediaFileOtherFilesPanel from '$lib/components/app/media/files/other-files/MediaFileOtherFilesPanel.svelte';
import MediaFileSummary from '$lib/components/app/media/files/MediaFileSummary.svelte';
import { provenanceFields } from '$lib/components/app/media/files/provenance/mediaFileTrackProvenance';
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
		expect(body).toContain('Other files');
		expect(body).toContain('No other files present.');
		expect(body).toContain('Missing');
		expect(body).toContain('Automatic search');
		expect(body).toContain('Manual search');
		expect(body).not.toContain('Delete file');
		expect(body).not.toContain('Formats');
	});

	it('renders media track details, collapsed chapters, and missing expected languages', () => {
		const { body } = renderWithTooltip(MediaFileDetailsAccordion, { row: detailedFileRow() });

		expect(body).toContain('Track Nr.');
		expect(body).toContain('Provenance');
		expect(body).toContain('Actions');
		expect(body).toContain('Video track');
		expect(body).toContain('h264');
		expect(body).toContain('1920x1080');
		expect(body).toContain('Audio track');
		expect(body).toContain('English');
		expect(body).toContain('DTS');
		expect(body).toContain('Track provenance');
		expect(body.match(/aria-label="Track provenance"/g) ?? []).toHaveLength(4);
		expect(body).toContain('Subtitle track');
		expect(body).toContain('Signs');
		expect(body).toContain('Chapter');
		expect(body).toContain('2-21');
		expect(body).toContain('20 chapters');
		expect(body.match(/Delete embedded track/g) ?? []).toHaveLength(4);
		expect(body).not.toContain('Scenario.Movie.2026.1080p.japanese.srt');
		expect(body).not.toContain('Opening');
		expect(body).toContain('German');
		expect(body).toContain('Missing expected audio track');
		expect(body.indexOf('Missing expected audio track')).toBeLessThan(
			body.indexOf('Subtitle track')
		);
		expect(body).toContain('border-t-4');
		expect(body).toContain('[&amp;>td]:border-t-4');
	});

	it('renders missing audio target rows when the profile target is not required', () => {
		const { body } = renderWithTooltip(MediaFileDetailsAccordion, {
			row: {
				...detailedFileRow(),
				missingTracks: [
					{
						key: 'missing-audio-italian',
						type: 'audio' as const,
						language: 'italian',
						description: 'Missing expected audio track',
						state: {
							visualState: 'missing_placeholder' as const,
							statusLabel: 'Missing',
							details: ['Missing expected audio: italian']
						}
					}
				]
			}
		});

		expect(body).toContain('Italian');
		expect(body).toContain('Missing expected audio track');
		expect(body.indexOf('Missing expected audio track')).toBeLessThan(
			body.indexOf('Subtitle track')
		);
		expect(body).toContain('bg-destructive/10 text-destructive');
	});

	it('renders other files with path, type, and subtitle state badges', () => {
		const { body } = renderWithTooltip(MediaFileOtherFilesPanel, {
			row: detailedFileRow(),
			canManage: true,
			onDelete: vi.fn()
		});

		expect(body).toContain('Other files');
		expect(body).toContain('Type');
		expect(body).toContain('Language');
		expect(body).toContain('Subtitle');
		expect(body).toContain('Metadata');
		expect(body).toContain('Unknown');
		expect(body).toContain('Japanese');
		expect(body).toContain('Scenario.Movie.2026.1080p.japanese.srt');
		expect(body).toContain('poster.jpg');
		expect(body).toContain('notes.bin');
		expect(body).toContain('Missing');
		expect(body).toContain('bg-destructive/10 text-destructive');
		expect(otherFilesOrder(body)).toEqual(['Other files', 'Type', 'Language', 'Actions']);
		expect(body.match(/>Type</g) ?? []).toHaveLength(1);
		expect(body.match(/>Language</g) ?? []).toHaveLength(1);
		expect(body.match(/>Actions</g) ?? []).toHaveLength(1);
		expect(body.match(/Delete other file/g) ?? []).toHaveLength(2);
	});

	it('renders an empty other files state', () => {
		const { body } = renderWithTooltip(MediaFileOtherFilesPanel, {
			row: { ...detailedFileRow(), otherFiles: [] },
			canManage: true,
			onDelete: vi.fn()
		});

		expect(body).toContain('Other files');
		expect(body).toContain('No other files present.');
	});

	it('marks detected external subtitles outside the target languages', () => {
		const { body } = renderWithTooltip(MediaFileOtherFilesPanel, {
			row: {
				...detailedFileRow(),
				otherFiles: [
					{
						type: 'subtitle' as const,
						path: '/library/Scenario Movie/Scenario.Movie.2026.1080p.spanish.srt',
						status: 'available' as const,
						language: 'spanish',
						state: {
							visualState: 'unwanted' as const,
							statusLabel: 'Unwanted',
							details: ['Subtitle language is outside enabled profile targets.']
						}
					}
				]
			},
			canManage: true,
			onDelete: vi.fn()
		});

		expect(body).toContain('Spanish');
		expect(body).toContain('bg-yellow-500/10 text-yellow-800 dark:text-yellow-300');
		expect(body).toContain('Unwanted');
	});

	it('formats all track provenance fields for the tooltip', () => {
		const fields = provenanceFields(detailedFileRow().tracks[1].provenance!);
		const values = fields.map((field) => `${field.label}: ${field.value}`);

		expect(values).toContain('ID: 11111111-1111-1111-1111-111111111111');
		expect(values).toContain('Media item: 22222222-2222-2222-2222-222222222222');
		expect(values).toContain('Component type: audio');
		expect(values).toContain('Component key: audio-source-1');
		expect(values).toContain('Release group: ARR');
		expect(values).toContain('Release name: Scenario.Release');
		expect(values).toContain('Source provider: Scenario Indexer');
		expect(values).toContain('Source file: /downloads/Scenario.Release.mkv');
		expect(values).toContain('Retained source: 33333333-3333-3333-3333-333333333333');
		expect(values).toContain('Source stream: 1');
		expect(values).toContain('Created: 2026-01-02T03:04:05Z');
		expect(values).toContain('Updated: 2026-01-02T04:05:06Z');
		expect(values.join('\n')).toContain('componentAssembly');
	});

	it('renders compact requirement states in file summaries', () => {
		const { body } = renderWithTooltip(MediaFileSummary, {
			mediaItemId: 'media-1',
			mediaTitle: 'Scenario Movie',
			row: detailedFileRow(),
			canManage: true,
			searching: false,
			fileLabel: 'Movie file',
			missingLabel: 'No matched file for this movie',
			onAutoSearch: vi.fn(),
			onManualSearch: vi.fn(),
			onDelete: vi.fn()
		});

		expect(body).toContain('Subtitles');
		expect(body).toContain('Audio');
		expect(body).toContain('Other files');
		expect(body).toContain('poster.jpg');
		expect(body).toContain('Missing');
		expect(body).toContain('text-emerald-600');
		expect(body).toContain('text-orange-500');
		expect(body).toContain('text-destructive');
		expect(body).toContain('Video requirements met');
		expect(body).toContain('Missing required audio: German');
		expect(body).toContain('Missing subtitles: Japanese');
		expect(body).not.toContain('>Video</strong>');
		expect(body).not.toContain('>Audio</strong>');
		expect(body).not.toContain('>Subtitles</strong>');
		expect(body).not.toContain('Other File');
		expect(body).not.toContain('Formats');
		expect(summaryOrder(body)).toEqual(['File', 'Size', 'Quality', 'Score', 'Status']);
	});
});

function summaryOrder(body: string) {
	return ['File', 'Size', 'Quality', 'Score', 'Status'].sort(
		(left, right) => body.indexOf(left) - body.indexOf(right)
	);
}

function otherFilesOrder(body: string) {
	return ['Other files', 'Type', 'Language', 'Actions'].sort(
		(left, right) => body.indexOf(left) - body.indexOf(right)
	);
}

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
		score: 120,
		tracks: [
			{
				type: 'video',
				index: 0,
				codec: 'h264',
				width: 1920,
				height: 1080,
				state: {
					visualState: 'matching',
					statusLabel: 'Matching',
					details: ['Video track satisfies the profile target.']
				}
			},
			{
				type: 'audio',
				index: 1,
				codec: 'DTS',
				language: 'eng',
				channels: 6,
				provenance: {
					id: '11111111-1111-1111-1111-111111111111',
					mediaItemId: '22222222-2222-2222-2222-222222222222',
					componentType: 'audio',
					componentKey: 'audio-source-1',
					releaseGroup: 'ARR',
					releaseName: 'Scenario.Release',
					sourceProvider: 'Scenario Indexer',
					sourceFilePath: '/downloads/Scenario.Release.mkv',
					retainedSourceId: '33333333-3333-3333-3333-333333333333',
					sourceStreamId: 1,
					transformationChain: [{ kind: 'componentAssembly', inputPath: '/downloads/audio.mka' }],
					createdAt: '2026-01-02T03:04:05Z',
					updatedAt: '2026-01-02T04:05:06Z'
				},
				state: {
					visualState: 'matching',
					statusLabel: 'Matching',
					details: ['Audio track satisfies a profile target.']
				}
			},
			{
				type: 'subtitle',
				index: 2,
				codec: 'SRT',
				language: 'eng',
				title: 'Signs',
				state: {
					visualState: 'matching',
					statusLabel: 'Matching',
					details: ['Embedded subtitle satisfies the subtitle target.']
				}
			},
			{
				type: 'subtitle',
				index: 3,
				codec: 'SRT',
				language: 'spa',
				title: 'Spanish',
				state: {
					visualState: 'unwanted',
					statusLabel: 'Unwanted',
					details: ['Subtitle language is outside enabled profile targets.']
				}
			}
		],
		chapters: Array.from({ length: 20 }, (_, index) => ({
			index: index + 1,
			title: index === 0 ? 'Opening' : `Chapter ${index + 1}`,
			startTime: `00:${String(index).padStart(2, '0')}:00`,
			endTime: `00:${String(index + 1).padStart(2, '0')}:00`
		})),
		otherFiles: [
			{
				type: 'subtitle',
				path: '/library/Scenario Movie/Scenario.Movie.2026.1080p.japanese.srt',
				status: 'missing',
				language: 'japanese',
				state: {
					visualState: 'missing_placeholder',
					statusLabel: 'Missing',
					details: ['Missing expected external subtitle: japanese']
				}
			},
			{
				type: 'metadata',
				path: '/library/Scenario Movie/poster.jpg',
				status: 'available'
			},
			{
				type: 'unknown',
				path: '/library/Scenario Movie/notes.bin',
				status: 'available'
			}
		],
		subtitleSatisfaction: {
			state: 'missing',
			mode: 'mixed',
			wantedLanguages: ['english', 'japanese'],
			matchedLanguages: ['english'],
			missingLanguages: ['japanese']
		},
		missingTracks: [
			{
				key: 'missing-audio-german',
				type: 'audio',
				language: 'german',
				description: 'Missing expected audio track',
				state: {
					visualState: 'missing_placeholder',
					statusLabel: 'Missing',
					details: ['Missing expected audio: german']
				}
			}
		],
		requirements: {
			video: { state: 'satisfied', label: 'Ok', details: ['Video requirements met'] },
			audio: { state: 'missing', label: 'Missing', details: ['Missing required audio: German'] },
			subtitles: { state: 'partial', label: 'Partial', details: ['Missing subtitles: Japanese'] }
		}
	};
}
