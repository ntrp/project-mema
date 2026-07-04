import { describe, expect, it, vi } from 'vitest';

import { crewRolePreviews } from '$lib/components/app/media/people/mediaPeople';
import { mediaRootPreview } from '$lib/components/app/media/collection/mediaRootPreview';
import { seasonFileSummary } from '$lib/components/app/media/series/mediaSeasonSummary';
import {
	episodeTitle,
	seasonEpisodeRows,
	seasonMonitored
} from '$lib/components/app/media/series/mediaSeriesRows';
import {
	episodeLabel,
	episodeNumbers,
	episodeValueFromNumbers,
	seasonOptions,
	selectedSeason
} from '$lib/components/app/media/release-override/releaseOverrideSeriesOptions';
import {
	releaseSearchQuery,
	releaseSearchQueryVariants
} from '$lib/components/app/media/release-search/releaseSearchQuery';
import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import type {
	DownloadActivity,
	MediaItem,
	MediaMetadataDetails,
	MediaMetadataEpisode
} from '$lib/settings/types';

const activeStatus = {
	status: 'downloading',
	label: 'downloading',
	activity: { status: 'downloading' } as DownloadActivity
} as const;

describe('series helper projections (SCN-MEDIA-004)', () => {
	it('summarizes downloaded, queued, and missing season files', () => {
		const summary = seasonFileSummary([
			{ row: { exists: true, sizeBytes: 1024 } as MediaFileRow },
			{ row: { exists: false } as MediaFileRow, activityStatus: activeStatus },
			{ row: { exists: false } as MediaFileRow }
		]);

		expect(summary).toEqual({
			label: '1 + 1 / 3',
			size: '1.00 KiB',
			tone: 'active',
			hasActive: true
		});
	});

	it('creates episode rows with existing files, missing fallbacks, and active activity', () => {
		const rows = seasonEpisodeRows(
			{
				name: 'Season 2',
				episodes: [
					{ episodeNumber: 1, name: 'Start', monitored: false },
					{ episodeNumber: 2, name: 'Next', monitored: true }
				]
			},
			0,
			[
				{
					exists: true,
					seasonNumber: 2,
					episodeNumber: 1,
					relativePath: 'S02E01.mkv'
				} as MediaFileRow
			],
			[
				{
					id: 'activity-1',
					mediaItemId: 'media-1',
					mediaTitle: 'Scenario Series',
					mediaType: 'serie',
					releaseTitle: 'Scenario.Series.S02E02.1080p.WEB-DL',
					indexerName: 'Indexer',
					downloadClientName: 'Client',
					downloadUrl: 'https://example.test/download',
					status: 'queued'
				} as unknown as DownloadActivity
			],
			'media-1'
		);

		expect(
			rows.map((row) => [row.row.exists, row.row.episodeNumber, row.activityStatus?.status])
		).toEqual([
			[true, 1, undefined],
			[false, 2, 'queued']
		]);
		expect(episodeTitle({ episodeNumber: 2, name: 'Next' } as MediaMetadataEpisode)).toBe(
			'2 - Next'
		);
	});

	it('builds season and episode controls from metadata details', () => {
		const details = {
			seasons: [
				{ name: 'Season 10', episodes: [] },
				{ name: 'Specials', episodes: [{ episodeNumber: 1, name: 'Bonus' }] },
				{ name: 'Season 2', monitored: true, episodes: [] }
			]
		} as MediaMetadataDetails;
		const options = seasonOptions(details);

		expect(options.map((option) => option.value)).toEqual(['0', '2', '10']);
		expect(selectedSeason(options, '2')?.label).toBe('Season 2');
		expect(episodeLabel({ episodeNumber: 3, name: 'Third' })).toBe('E03 Third');
		expect(episodeNumbers('3, nope 1 3 0 -2')).toEqual([3, 1, 3]);
		expect(episodeValueFromNumbers([3, 1, 3, 0])).toBe('1, 3');
		expect(seasonMonitored(details.seasons![1])).toBe(false);
		expect(seasonMonitored(details.seasons![2])).toBe(true);
	});
});

describe('media metadata previews (SCN-MEDIA-004)', () => {
	it('limits crew previews to prioritized roles and first three names', () => {
		expect(
			crewRolePreviews([
				{ label: 'Director', value: 'One, Two, Three, Four' },
				{ label: 'Composer', value: 'Hidden' },
				{ label: 'Writer', value: 'Writer One' }
			])
		).toEqual([
			{ role: 'Director', names: ['One', 'Two', 'Three'] },
			{ role: 'Writer', names: ['Writer One'] }
		]);
	});

	it('renders sanitized library root previews from naming templates', () => {
		const item = { type: 'movie', title: 'Bad: /Title*', year: 2026 } as MediaItem;

		expect(
			mediaRootPreview(
				item,
				{ path: '/movies' } as never,
				{
					movieFolderFormat: '{Movie Title} ({Release Year})'
				} as never
			)
		).toBe('/movies/Bad -  Title (2026)');
		expect(mediaRootPreview(item, undefined, undefined)).toBe('-');
	});
});

describe('release search query helpers (SCN-MEDIA-002)', () => {
	it('builds movie, season, and episode query variants without duplicates', () => {
		const movie = { type: 'movie', title: ' Scenario Movie ', year: 2026 } as MediaItem;
		const series = { type: 'serie', title: 'Scenario Series', year: 2026 } as MediaItem;

		expect(releaseSearchQuery(movie)).toBe('Scenario Movie 2026');
		expect(releaseSearchQuery(series, { type: 'season', seasonNumber: 2 })).toBe(
			'Scenario Series s2'
		);
		expect(
			releaseSearchQueryVariants(series, {
				type: 'episode',
				seasonNumber: 2,
				episodeNumber: 3
			})
		).toEqual(['Scenario Series s2e3', 'Scenario Series S02E03']);
		expect(releaseSearchQueryVariants(series, { type: 'title' })).toEqual(['Scenario Series 2026']);
	});

	it('keeps release search logs deterministic when the clock is fixed', async () => {
		vi.useFakeTimers();
		vi.setSystemTime(new Date(2026, 0, 2, 3, 4, 5, 6));
		const { applyStatusToLog, createLogEntry, placeholderLogEntry } =
			await import('$lib/components/app/media/release-search/releaseSearchLog');

		expect(placeholderLogEntry()).toMatchObject({
			id: 'placeholder',
			message: 'Press search to start'
		});
		expect(createLogEntry('Started')).toMatchObject({
			timestamp: '03:04:05.006',
			message: 'Started'
		});
		expect(
			applyStatusToLog(
				[
					{
						id: 'indexer:Alpha:Scenario',
						timestamp: '03:04:05.006',
						message: 'Searching Alpha'
					}
				],
				{
					kind: 'indexer_finish',
					indexerName: 'Alpha',
					query: 'Scenario',
					message: '1 result',
					cacheHit: true,
					durationMs: 12
				}
			)
		).toMatchObject([{ resultMessage: '1 result', cacheHit: true, durationMs: 12 }]);
		vi.useRealTimers();
	});
});
