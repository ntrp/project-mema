import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import ReleaseOverrideEpisodeSelect from '$lib/components/app/media/ReleaseOverrideEpisodeSelect.svelte';
import ReleaseOverrideMovieField from '$lib/components/app/media/ReleaseOverrideMovieField.svelte';
import ReleaseOverrideSeasonSelect from '$lib/components/app/media/ReleaseOverrideSeasonSelect.svelte';
import ReleaseOverrideSeriesFields from '$lib/components/app/media/ReleaseOverrideSeriesFields.svelte';
import type { SeasonOption } from '$lib/components/app/media/releaseOverrideSeriesOptions';
import type { ReleaseOverrideDraft } from '$lib/components/app/media/releaseOverrideDetails';
import type { MediaItem, MediaMetadataEpisode } from '$lib/settings/types';

const seriesItem = {
	id: 'series-1',
	type: 'series',
	title: 'Scenario Series',
	monitored: true,
	monitorMode: 'all_episodes',
	minimumAvailability: 'released',
	externalProvider: 'tmdb',
	externalId: '100',
	status: 'missing',
	filePaths: [],
	metadataFilePaths: [],
	createdAt: '2026-07-03T00:00:00Z',
	updatedAt: '2026-07-03T00:00:00Z'
} as MediaItem;

function draft(overrides: Partial<ReleaseOverrideDraft> = {}): ReleaseOverrideDraft {
	return {
		movieTitle: '',
		seriesTitle: 'Scenario Series',
		seasonNumber: '1',
		episodeNumbers: '1, 3',
		releaseGroup: '',
		quality: '',
		languages: [],
		...overrides
	};
}

function episodes(): MediaMetadataEpisode[] {
	return [
		{ episodeNumber: 1, name: 'Pilot' },
		{ episodeNumber: 2, name: 'Second' },
		{ episodeNumber: 3, name: 'Finale' }
	];
}

const seasons: SeasonOption[] = [
	{
		value: '1',
		label: 'Season 1',
		season: {
			name: 'Season 1',
			episodeCount: 3,
			episodes: episodes()
		}
	}
];

describe('rendered release override components (SCN-MEDIA-002)', () => {
	it('renders movie override title editing with the current value', () => {
		const { body } = render(ReleaseOverrideMovieField, {
			props: {
				value: 'Scenario Movie'
			}
		});

		expect(body).toContain('Movie');
		expect(body).toContain('value="Scenario Movie"');
		expect(body).not.toContain('Movie matches');
	});

	it('renders series override fields with fallback season and selected episodes', () => {
		const { body } = render(ReleaseOverrideSeriesFields, {
			props: {
				item: seriesItem,
				draft: draft()
			}
		});

		expect(body).toContain('Series');
		expect(body).toContain('Season');
		expect(body).toContain('Episodes');
		expect(body).toContain('value="Scenario Series"');
		expect(body).toContain('value="1"');
		expect(body).toContain('E01');
		expect(body).toContain('E03');
	});

	it('renders season options when metadata supplies known seasons', () => {
		const { body } = render(ReleaseOverrideSeasonSelect, {
			props: {
				value: '1',
				label: 'Season 1',
				seasons,
				onChange: vi.fn()
			}
		});

		expect(body).toContain('Season');
		expect(body).toContain('Season 1');
		expect(body).not.toContain('type="number"');
	});

	it('renders a numeric season fallback when provider metadata is unavailable', () => {
		const { body } = render(ReleaseOverrideSeasonSelect, {
			props: {
				value: '3',
				label: '',
				seasons: [],
				onChange: vi.fn()
			}
		});

		expect(body).toContain('Season');
		expect(body).toContain('type="number"');
		expect(body).toContain('value="3"');
	});

	it('renders episode labels and the empty placeholder state', () => {
		const selected = render(ReleaseOverrideEpisodeSelect, {
			props: {
				value: '1, 3',
				episodes: episodes(),
				onChange: vi.fn()
			}
		});
		expect(selected.body).toContain('E01 Pilot');
		expect(selected.body).toContain('E03 Finale');

		const empty = render(ReleaseOverrideEpisodeSelect, {
			props: {
				value: '',
				episodes: episodes(),
				onChange: vi.fn()
			}
		});
		expect(empty.body).toContain('Select episodes');
	});
});
