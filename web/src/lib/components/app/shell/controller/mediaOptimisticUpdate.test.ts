import { describe, expect, it } from 'vitest';

import { mediaUpdateMessage, optimisticMediaItem } from './mediaOptimisticUpdate';
import type { MediaItem, MediaItemUpdateRequest } from '$lib/settings/types';

const series = {
	type: 'series',
	title: 'Scenario Series',
	monitored: true,
	monitorMode: 'all_episodes',
	qualityProfileId: 'profile-old',
	minimumAvailability: 'released',
	libraryFolderId: 'folder-old',
	seasons: [
		{
			name: 'Season 1',
			monitored: true,
			episodes: [
				{ name: 'Pilot', episodeNumber: 1, monitored: true },
				{ name: 'Second', episodeNumber: 2, monitored: true }
			]
		},
		{
			name: 'Season 2',
			monitored: false,
			episodes: [{ name: 'Future', episodeNumber: 1, monitored: false }]
		}
	]
} as MediaItem;

describe('optimistic media update state (SCN-MEDIA-007)', () => {
	it('applies top-level media settings while preserving untouched fields', () => {
		const request: MediaItemUpdateRequest = {
			qualityProfileId: 'profile-new',
			minimumAvailability: 'announced',
			libraryFolderId: 'folder-new',
			monitored: false
		};

		const next = optimisticMediaItem(series, request);

		expect(next.qualityProfileId).toBe('profile-new');
		expect(next.minimumAvailability).toBe('announced');
		expect(next.libraryFolderId).toBe('folder-new');
		expect(next.monitored).toBe(false);
		expect(next.monitorMode).toBe('all_episodes');
		expect(next.seasons).toBe(series.seasons);
		expect(mediaUpdateMessage(series, next, request)).toBe('Media root updated');
	});

	it('updates every episode in a patched season and reports season state', () => {
		const request: MediaItemUpdateRequest = {
			monitorSeasonName: 'Season 1',
			seasonMonitored: false
		};

		const next = optimisticMediaItem(series, request);

		expect(next.seasons?.[0].monitored).toBe(false);
		expect(next.seasons?.[0].episodes?.map((episode) => episode.monitored)).toEqual([false, false]);
		expect(next.seasons?.[1]).toBe(series.seasons?.[1]);
		expect(mediaUpdateMessage(series, next, request)).toBe(
			'Season "Season 1" is now not monitored'
		);
	});

	it('updates a single episode and derives the parent season monitored state', () => {
		const request: MediaItemUpdateRequest = {
			monitorSeasonName: 'Season 2',
			monitorEpisodeNumber: 1,
			episodeMonitored: true
		};

		const next = optimisticMediaItem(series, request);

		expect(next.seasons?.[1].monitored).toBe(true);
		expect(next.seasons?.[1].episodes?.[0].monitored).toBe(true);
		expect(mediaUpdateMessage(series, next, request)).toBe('Episode "Future" is now monitored');
	});

	it('falls back to stable messages for broad monitoring changes', () => {
		const movie = { type: 'movie', title: 'Scenario Movie', monitored: true } as MediaItem;
		const movieRequest: MediaItemUpdateRequest = { monitored: false };
		const nextMovie = optimisticMediaItem(movie, movieRequest);

		expect(mediaUpdateMessage(movie, nextMovie, movieRequest)).toBe('Movie is now not monitored');

		const seriesRequest: MediaItemUpdateRequest = { monitorMode: 'none' };
		const nextSeries = optimisticMediaItem(series, seriesRequest);
		expect(mediaUpdateMessage(series, nextSeries, seriesRequest)).toBe(
			'Series is now not monitored'
		);
		expect(mediaUpdateMessage(series, nextSeries, {})).toBe('Media settings saved');
	});
});
