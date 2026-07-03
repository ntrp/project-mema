import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import MediaMetadataCore from '$lib/components/app/media/metadata/MediaMetadataCore.svelte';
import MediaMetadataHero from '$lib/components/app/media/metadata/MediaMetadataHero.svelte';
import type { MediaMetadataDetails } from '$lib/settings/types';
import { renderWithTooltip } from './renderHelpers';

function detail(overrides: Partial<MediaMetadataDetails> = {}): MediaMetadataDetails {
	return {
		title: 'Scenario Series',
		type: 'series',
		year: 2026,
		externalProvider: 'tmdb',
		externalId: 'series-1',
		overview: 'A carefully monitored scenario series.',
		posterPath: '/poster.jpg',
		collectionId: 'collection-1',
		collectionName: 'Scenario Collection',
		trailerUrl: 'https://video.test/trailer',
		runtimeMinutes: 95,
		status: 'returning',
		monitored: true,
		genres: ['Drama', 'Mystery'],
		keywords: ['scenario', 'coverage'],
		facts: [
			{ label: 'Certification', value: 'TV-14' },
			{ label: 'Director', value: 'Ada Example' },
			{ label: 'Writer', value: 'Grace Example' },
			{ label: 'Network', value: 'Local TV' }
		],
		seasons: [
			{
				name: 'Season 1',
				episodeCount: 2,
				episodes: [
					{
						name: 'Pilot',
						episodeNumber: 1,
						overview: 'The scenario starts.',
						airDate: '2026-01-01',
						monitored: true
					},
					{
						name: 'Second',
						episodeNumber: 2,
						overview: 'The scenario continues.',
						airDate: '2026-01-08',
						monitored: false
					}
				]
			}
		],
		cast: [
			{
				name: 'Actor One',
				role: 'Lead',
				profilePath: '/actor-one.jpg',
				externalProvider: 'tmdb',
				externalId: 'actor-1'
			},
			{ name: 'Actor Two', role: 'Support' }
		],
		crew: [
			{
				name: 'Ada Example',
				role: 'Director',
				externalProvider: 'tmdb',
				externalId: 'director-1'
			},
			{
				name: 'Grace Example',
				role: 'Writer',
				externalProvider: 'tmdb',
				externalId: 'writer-1'
			}
		],
		...overrides
	} as MediaMetadataDetails;
}

describe('rendered metadata components (SCN-MEDIA-008)', () => {
	it('renders metadata hero identity, status, collection, and actions', () => {
		const { body } = renderWithTooltip(MediaMetadataHero, {
			detail: detail(),
			titleId: 'metadata-title',
			mediaStatus: 'missing' as const,
			monitorStatusText: 'Monitored',
			monitorHintText: 'All episodes monitored',
			onToggleMonitor: vi.fn()
		});

		expect(body).toContain('Scenario Series');
		expect(body).toContain('TV-14');
		expect(body).toContain('2026');
		expect(body).toContain('1h 35m');
		expect(body).toContain('Missing');
		expect(body).toContain('Drama');
		expect(body).toContain('Mystery');
		expect(body).toContain('Scenario Collection');
		expect(body).toContain('Trailer');
	});

	it('renders metadata core overview, crew, keywords, seasons, and cast', () => {
		const { body } = render(MediaMetadataCore, {
			props: {
				detail: detail()
			}
		});

		expect(body).toContain('Overview');
		expect(body).toContain('A carefully monitored scenario series.');
		expect(body).toContain('Crew');
		expect(body).toContain('Ada Example');
		expect(body).toContain('Grace Example');
		expect(body).toContain('/people/tmdb/director-1');
		expect(body).toContain('/people/tmdb/writer-1');
		expect(body).toContain('scenario');
		expect(body).toContain('coverage');
		expect(body).toContain('Season 1');
		expect(body).toContain('2 episodes');
		expect(body).not.toContain('>-</span>');
		expect(body).toContain('aria-expanded="false"');
		expect(body).toContain('Actor One');
		expect(body).toContain('Lead');
	});

	it('renders metadata fallbacks when optional provider fields are absent', () => {
		const { body } = render(MediaMetadataCore, {
			props: {
				detail: detail({
					overview: undefined,
					facts: [],
					keywords: [],
					seasons: [],
					cast: []
				})
			}
		});

		expect(body).toContain('No overview available.');
		expect(body).not.toContain('Seasons');
		expect(body).not.toContain('Cast');
	});
});
