import { describe, expect, it, vi } from 'vitest';

import MediaFileDetailsAccordion from '$lib/components/app/media/files/MediaFileDetailsAccordion.svelte';
import MediaFileOtherFilesPanel from '$lib/components/app/media/files/other-files/MediaFileOtherFilesPanel.svelte';
import MediaFileSummary from '$lib/components/app/media/files/MediaFileSummary.svelte';
import { fileTrackDetailRows } from '$lib/components/app/media/files/mediaFileDetails';
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
		expect(body).toContain('Source audio');
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

	it('keeps raw language ids for fulfillment action payloads', () => {
		const rows = fileTrackDetailRows(detailedFileRow());

		const audio = rows.find((row) => row.type === 'audio' && row.language === 'English');
		const missing = rows.find((row) => row.type === 'audio' && row.missing);

		expect(audio?.languageId).toBe('eng');
		expect(audio?.trackId).toBe('aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa');
		expect(missing?.languageId).toBe('german');
	});

	it('shows transcode only when an existing track has a codec mismatch', () => {
		const codecMismatch = detailedFileRow();
		codecMismatch.tracks[1] = {
			...codecMismatch.tracks[1],
			state: {
				visualState: 'partial',
				statusLabel: 'Partial',
				details: ['Audio codec does not meet the profile target']
			}
		};
		const channelMismatch = detailedFileRow();
		channelMismatch.tracks[1] = {
			...channelMismatch.tracks[1],
			state: {
				visualState: 'partial',
				statusLabel: 'Partial',
				details: ['Audio channels do not meet the profile target']
			}
		};

		expect(
			renderWithTooltip(MediaFileDetailsAccordion, {
				row: codecMismatch,
				canManage: true
			}).body
		).toContain('Transcode audio');
		expect(
			renderWithTooltip(MediaFileDetailsAccordion, {
				row: channelMismatch,
				canManage: true
			}).body
		).not.toContain('Transcode audio');
	});

	it('shows video transcode for supported video codec or pixel mismatches', () => {
		const row = detailedFileRow();
		row.tracks[0] = {
			...row.tracks[0],
			id: 'bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb',
			state: {
				visualState: 'partial',
				statusLabel: 'Partial',
				details: ['Pixel format does not meet the profile target']
			}
		};

		expect(
			renderWithTooltip(MediaFileDetailsAccordion, {
				row,
				canManage: true
			}).body
		).toContain('Transcode video');
	});

	it('shows container remux on the file overview row', () => {
		const row = detailedFileRow();
		row.requirements = {
			...row.requirements!,
			video: { state: 'satisfied', label: 'Ok', details: ['Video requirements met'] },
			container: {
				state: 'pending',
				label: 'Pending',
				details: ['Container does not meet the profile target']
			}
		};

		const summary = renderWithTooltip(MediaFileSummary, {
			mediaItemId: 'media-1',
			mediaTitle: 'Scenario Movie',
			row,
			canManage: true,
			searching: false,
			onAutoSearch: vi.fn(),
			onManualSearch: vi.fn(),
			onDelete: vi.fn()
		}).body;
		const details = renderWithTooltip(MediaFileDetailsAccordion, { row, canManage: true }).body;

		expect(summary).toContain('Remux container');
		expect(summary).toContain('Pending');
		expect(details).not.toContain('Remux container');
	});

	it('hides the container badge when the container is ok', () => {
		const summary = renderWithTooltip(MediaFileSummary, {
			mediaItemId: 'media-1',
			mediaTitle: 'Scenario Movie',
			row: detailedFileRow(),
			canManage: true,
			searching: false,
			onAutoSearch: vi.fn(),
			onManualSearch: vi.fn(),
			onDelete: vi.fn()
		}).body;

		expect(summary).not.toContain('Container requirements met');
		expect(summary).not.toContain('Remux container');
	});

	it('shows pending fulfillment actions as busy buttons', () => {
		const row = detailedFileRow();
		row.tracks[1] = {
			...row.tracks[1],
			state: {
				visualState: 'partial',
				statusLabel: 'Partial',
				details: ['Audio codec does not meet the profile target']
			}
		};

		const { body } = renderWithTooltip(MediaFileDetailsAccordion, {
			row,
			canManage: true,
			pendingFulfillmentActionKeys: [
				'audio_transcode|aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa|audio|eng'
			]
		});

		expect(body).toContain('aria-busy="true"');
		expect(body).toContain('animate-spin');
	});

	it('renders other files with path, type, and subtitle state badges', () => {
		const { body } = renderWithTooltip(MediaFileOtherFilesPanel, {
			row: detailedFileRow(),
			canManage: true,
			onDelete: vi.fn()
		});

		expect(body).toContain('Other files');
		expect(body).toContain('Type');
		expect(body).toContain('Subtype');
		expect(body).toContain('Language');
		expect(body).toContain('Subtitle');
		expect(body).toContain('Metadata');
		expect(body).toContain('Unknown');
		expect(body).toContain('SubRip');
		expect(body).toContain('POSTER');
		expect(body).toContain('BIN');
		expect(body).toContain('Japanese');
		expect(body).toContain('Scenario.Movie.2026.1080p.japanese.srt');
		expect(body).toContain('poster.jpg');
		expect(body).toContain('notes.bin');
		expect(body).toContain('Missing');
		expect(body).toContain('Download subtitle');
		expect(body).toContain('bg-destructive/10 text-destructive');
		expect(otherFilesOrder(body)).toEqual([
			'Other files',
			'Type',
			'Subtype',
			'Language',
			'Actions'
		]);
		expect(body.match(/>Type</g) ?? []).toHaveLength(1);
		expect(body.match(/>Subtype</g) ?? []).toHaveLength(1);
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
		expect(body).toContain('bg-yellow-300/10 text-yellow-600 dark:text-yellow-200');
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
		expect(body).toContain('Video track satisfies the profile target.');
		expect(body).toContain('Missing expected audio track');
		expect(body).toContain('Missing expected external subtitle: japanese');
		expect(body).not.toContain('>Video</strong>');
		expect(body).not.toContain('>Audio</strong>');
		expect(body).not.toContain('>Subtitles</strong>');
		expect(body).not.toContain('Other File');
		expect(body).not.toContain('Formats');
		expect(summaryOrder(body)).toEqual(['File', 'Size', 'Quality', 'Score', 'Status']);
	});

	it('renders status icon tooltips with a list per track', () => {
		const base = detailedFileRow();
		const row = {
			...base,
			tracks: base.tracks.map((track) =>
				track.type === 'audio'
					? {
							...track,
							state: {
								visualState: 'partial' as const,
								statusLabel: 'Partial',
								details: [
									'Audio codec does not meet the profile target',
									'Audio channels do not meet the profile target'
								]
							}
						}
					: track
			),
			missingTracks: [],
			requirements: {
				...base.requirements!,
				audio: {
					state: 'partial' as const,
					label: 'Partial',
					details: ['Audio track does not meet the profile target']
				}
			}
		};

		const { body } = renderWithTooltip(MediaFileSummary, {
			mediaItemId: 'media-1',
			mediaTitle: 'Scenario Movie',
			row,
			canManage: true,
			searching: false,
			onAutoSearch: vi.fn(),
			onManualSearch: vi.fn(),
			onDelete: vi.fn()
		});

		expect(body).toContain('Audio track 1 - English');
		expect(body).toContain('Audio codec does not meet the profile target');
		expect(body).toContain('Audio channels do not meet the profile target');
	});
});

function summaryOrder(body: string) {
	return ['File', 'Size', 'Quality', 'Score', 'Status'].sort(
		(left, right) => body.indexOf(left) - body.indexOf(right)
	);
}

function otherFilesOrder(body: string) {
	return ['Other files', 'Type', 'Subtype', 'Language', 'Actions'].sort(
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
				id: 'aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa',
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
				subtype: 'subrip',
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
				status: 'available',
				subtype: 'poster'
			},
			{
				type: 'unknown',
				path: '/library/Scenario Movie/notes.bin',
				status: 'available',
				subtype: 'bin'
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
			container: { state: 'satisfied', label: 'Ok', details: ['Container requirements met'] },
			audio: { state: 'missing', label: 'Missing', details: ['Missing required audio: German'] },
			subtitles: { state: 'partial', label: 'Partial', details: ['Missing subtitles: Japanese'] }
		}
	};
}
