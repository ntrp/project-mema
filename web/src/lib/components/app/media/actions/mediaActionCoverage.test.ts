import { describe, expect, it } from 'vitest';

import {
	matchingLibraryFolders,
	mediaPosterUrl,
	preselectLibraryFolderId,
	preselectQualityProfileId
} from '$lib/components/app/media/actions/mediaActionDefaults';
import { fileDetailRows } from '$lib/components/app/media/files/mediaFileDetails';
import {
	toggledEpisodeMonitor,
	toggledMediaMonitor,
	toggledSeasonMonitor
} from '$lib/components/app/media/series/mediaMonitoring';
import {
	detailsFromOverrideDraft,
	overrideDraftFromRelease
} from '$lib/components/app/media/release-override/releaseOverrideDetails';
import type {
	Language,
	LibraryFolder,
	MediaItem,
	MediaMetadataSeason,
	MediaSearchResult,
	QualityProfileOption,
	ReleaseCandidate
} from '$lib/settings/types';

describe('media action defaults (SCN-MEDIA-003)', () => {
	it('preselects quality profiles and folders from candidate metadata', () => {
		const anime = {
			title: 'Scenario Anime',
			type: 'movie',
			overview: 'A feature anime',
			genres: ['Animation'],
			originalLanguage: 'ja'
		} as MediaSearchResult;
		const profiles = [
			{ id: 'any', name: 'Any acceptable' },
			{ id: 'anime-1080p', name: 'Anime 1080p' },
			{ id: 'uhd', name: 'UHD 2160p' }
		] as QualityProfileOption[];
		const folders = [
			{ id: 'movies', path: '/media/movies', kind: 'movie' },
			{ id: 'anime-movies', path: '/media/anime/movies', kind: 'movie' },
			{ id: 'series', path: '/media/tv', kind: 'series' }
		] as LibraryFolder[];

		expect(preselectQualityProfileId(anime, profiles)).toBe('anime-1080p');
		expect(preselectLibraryFolderId(anime, folders)).toBe('anime-movies');
		expect(
			preselectLibraryFolderId({ title: 'Show', type: 'serie' } as MediaSearchResult, folders)
		).toBe('series');
		expect(matchingLibraryFolders('movie', folders).map((folder) => folder.id)).toEqual([
			'movies',
			'anime-movies'
		]);
		expect(matchingLibraryFolders('serie', folders).map((folder) => folder.id)).toEqual(['series']);
		expect(mediaPosterUrl('/poster.jpg')).toBe('https://image.tmdb.org/t/p/w780/poster.jpg');
		expect(mediaPosterUrl('https://image.test/poster.jpg')).toBe('https://image.test/poster.jpg');
		expect(mediaPosterUrl()).toBeUndefined();
	});
});

describe('media monitoring payloads (SCN-MEDIA-004)', () => {
	it('toggles title, season, and episode monitoring into API payloads', () => {
		const movie = {
			title: 'Scenario Movie',
			type: 'movie',
			monitored: false,
			monitorMode: 'none',
			minimumAvailability: 'released'
		} as unknown as MediaItem;
		const series = {
			title: 'Scenario Series',
			type: 'serie',
			monitored: false,
			monitorMode: 'none',
			minimumAvailability: 'released',
			seasons: []
		} as unknown as MediaItem;
		const seasons = [
			{
				name: 'Season 1',
				monitored: false,
				episodes: [
					{ episodeNumber: 1, monitored: false },
					{ episodeNumber: 2, monitored: true }
				]
			},
			{ name: 'Season 2', monitored: false, episodes: [{ episodeNumber: 1, monitored: false }] }
		] as MediaMetadataSeason[];

		expect(toggledMediaMonitor(movie)).toEqual({ monitored: true, monitorMode: 'only_media' });
		expect(toggledMediaMonitor(series)).toEqual({
			monitored: true,
			monitorMode: 'future_episodes'
		});
		expect(toggledEpisodeMonitor(series, seasons, seasons[1], seasons[1].episodes![0])).toEqual({
			monitored: true,
			monitorMode: 'all_episodes',
			monitorSeasonName: 'Season 2',
			monitorEpisodeNumber: 1,
			episodeMonitored: true
		});
		expect(toggledSeasonMonitor(series, seasons, seasons[0])).toEqual({
			monitored: false,
			monitorMode: 'none',
			monitorSeasonName: 'Season 1',
			seasonMonitored: false
		});
	});
});

describe('media file detail rows (SCN-MEDIA-004)', () => {
	it('describes tracks, chapters, unwanted streams, and missing expected languages', () => {
		const rows = fileDetailRows({
			quality: 'HD-1080p',
			expectedLanguages: ['English', 'German'],
			expectedRequiredLanguages: ['German'],
			expectedSubtitleLanguages: ['English'],
			removeNonEnabledLanguages: true,
			removeNonEnabledSubtitleLanguages: true,
			tracks: [
				{
					index: 0,
					type: 'video',
					codec: 'h264',
					width: 1920,
					height: 1080,
					profile: 'High',
					pixelFormat: 'yuv420p',
					frameRate: '24',
					bitRate: '4000000'
				},
				{
					index: 1,
					type: 'audio',
					language: 'eng',
					codec: 'aac',
					channels: 6,
					bitRate: '640000',
					title: 'Main'
				},
				{ index: 2, type: 'audio', language: 'spa', codec: 'aac' },
				{ index: 3, type: 'subtitle', language: 'spa', codec: 'srt' }
			],
			chapters: [{ index: 0, title: 'Opening', startTime: '00:00:00', endTime: '00:01:00' }]
		} as never);

		expect(rows.map((row) => [row.type, row.language, row.description])).toContainEqual([
			'video',
			'-',
			'h264 · 1920x1080 · High · yuv420p · 24 · 4000 kbps'
		]);
		expect(rows.find((row) => row.key === 'audio-2')?.unwanted).toBe(true);
		expect(rows.find((row) => row.key === 'subtitle-3')?.unwanted).toBe(true);
		expect(rows.find((row) => row.key === 'chapter-0')?.description).toBe(
			'Opening · 00:00:00 - 00:01:00'
		);
		expect(rows.find((row) => row.missing)?.language).toBe('German');
	});
});

describe('release override details (SCN-MEDIA-002)', () => {
	it('builds override drafts and trims API details', () => {
		const languages = [
			{ code: 'eng', displayName: 'English', aliases: ['English'] },
			{ code: 'deu', displayName: 'German', aliases: ['Deutsch'] }
		] as Language[];
		const draft = overrideDraftFromRelease(
			{ title: 'Fallback Series' } as MediaItem,
			{
				title: 'Scenario.Series.S02E03E04.1080p-WEB',
				match: {
					matchedMedia: '',
					quality: 'HD-1080p',
					languages: ['English', 'Deutsch']
				}
			} as ReleaseCandidate,
			languages
		);

		expect(draft).toMatchObject({
			seriesTitle: 'Scenario Series',
			seasonNumber: '2',
			episodeNumbers: '3, 4',
			releaseGroup: 'WEB',
			quality: 'HD-1080p',
			languages: ['eng', 'deu']
		});
		expect(detailsFromOverrideDraft({ ...draft, movieTitle: '  Movie  ' }, languages)).toEqual({
			movieTitle: 'Movie',
			seriesTitle: 'Scenario Series',
			seasonNumber: 2,
			episodeNumbers: [3, 4],
			releaseGroup: 'WEB',
			quality: 'HD-1080p',
			languages: ['English', 'German']
		});
	});
});
