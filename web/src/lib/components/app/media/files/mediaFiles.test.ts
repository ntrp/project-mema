import { describe, expect, it } from 'vitest';

import { fileRow, mediaFileGroups } from '$lib/components/app/media/files/mediaFiles';
import {
	episodeKey,
	missingRow,
	seasonNumberFromName
} from '$lib/components/app/media/files/mediaFileMissing';
import { qualityInfo } from '$lib/components/app/media/files/mediaFileParsing';
import type { MediaItem } from '$lib/settings/types';

const movie = {
	type: 'movie',
	title: 'Scenario Movie',
	mediaFolderPath: '/library/Scenario Movie',
	filePaths: ['/library/Scenario Movie/Scenario.Movie.2026.1080p.WEB-DL.DDP5.1.EN.mkv'],
	files: [
		{
			path: '/library/Scenario Movie/Scenario.Movie.2026.1080p.WEB-DL.DDP5.1.EN.mkv',
			status: 'available',
			sizeBytes: 5 * 1024 * 1024 * 1024,
			tracks: [{ type: 'video', codec: 'h264' }],
			chapters: [{ title: 'Intro' }]
		}
	],
	qualityProfileId: 'profile-1'
} as MediaItem;

describe('media file display models (SCN-MEDIA-001)', () => {
	it('builds movie rows with relative paths, parsed quality, languages, and profile settings', () => {
		const row = fileRow(movie, movie.filePaths[0], [
			{
				id: 'profile-1',
				qualityIds: ['webdl-720p', 'webdl-1080p', 'bluray-1080p'],
				upgradesAllowed: true,
				upgradeUntilQualityId: 'bluray-1080p',
				audioTargets: [{ languageId: 'english', required: true }],
				subtitleTargets: [{ languageId: 'english' }],
				removeUnwantedAudio: true,
				removeUnwantedSubtitles: true
			}
		]);

		expect(row).toMatchObject({
			relativePath: 'Scenario.Movie.2026.1080p.WEB-DL.DDP5.1.EN.mkv',
			exists: true,
			videoCodec: '-',
			audioInfo: 'DDP DD+',
			size: '5.00 GiB',
			sizeBytes: 5368709120,
			languages: 'English',
			quality: 'WEBDL-1080p',
			formats: ['WEB-DL'],
			upgrade: {
				state: 'upgradeable',
				label: 'Upgradeable',
				reasons: ['Upgrade target is bluray-1080p']
			}
		});
		expect(row.tracks).toHaveLength(1);
		expect(row.chapters).toHaveLength(1);
	});

	it('creates missing rows for movies without files', () => {
		const groups = mediaFileGroups({ ...movie, filePaths: [], files: [] });

		expect(groups).toEqual([
			{ key: 'movie', title: 'Movie file', rows: [missingRow('movie-missing', 'Scenario Movie')] }
		]);
	});

	it('groups series files by season and adds missing episode rows', () => {
		const series = {
			...movie,
			type: 'serie',
			title: 'Scenario Series',
			filePaths: ['/tv/Scenario/S01E01.720p.Dual.mkv'],
			mediaFolderPath: '/tv/Scenario',
			seasons: [
				{
					name: 'Season 1',
					episodes: [
						{ episodeNumber: 1, name: 'Pilot' },
						{ episodeNumber: 2, name: 'Second' }
					]
				}
			]
		} as MediaItem;

		const groups = mediaFileGroups(series);

		expect(groups[0].key).toBe('season-1');
		expect(groups[0].rows.map((row) => [row.exists, row.episodeNumber, row.episodeTitle])).toEqual([
			[true, 1, undefined],
			[false, 2, 'Second']
		]);
		expect(groups[0].rows[0].languages).toBe('Dual');
	});

	it('parses season numbers and episode keys', () => {
		expect(seasonNumberFromName('Specials')).toBe(0);
		expect(seasonNumberFromName('Season 12')).toBe(12);
		expect(seasonNumberFromName('No number')).toBeUndefined();
		expect(episodeKey(2, 3)).toBe('2:3');
	});

	it('shows full quality labels when source and resolution are detected', () => {
		expect(qualityInfo('Scenario.Movie.2026.2160p.WEB-DL.mkv')).toBe('WEBDL-2160p');
		expect(qualityInfo('Madagascar.2005.2160p.BluRayRIP.TrueHD.7.1.Atmos.mkv')).toBe(
			'BluRay-2160p'
		);
		expect(qualityInfo('Scenario.Movie.2026.1080p.mkv')).toBe('1080p');
	});
});
