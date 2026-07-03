import { render } from 'svelte/server';
import { describe, expect, it } from 'vitest';

import MediaEpisodeRow from '$lib/components/app/media/series/MediaEpisodeRow.svelte';
import type { MediaMetadataEpisode } from '$lib/settings/types';

describe('rendered media episode components (SCN-MEDIA-003)', () => {
	it('renders episode title, date, overview, and still image', () => {
		const { body } = render(MediaEpisodeRow, {
			props: {
				title: 'S01E01 - Pilot',
				episode: episode({
					name: 'Pilot',
					overview: 'The first scenario episode.',
					airDate: '2026-07-03',
					stillPath: '/still.jpg'
				})
			}
		});

		expect(body).toContain('S01E01 - Pilot');
		expect(body).toContain('The first scenario episode.');
		expect(body).toContain('Jul 3, 2026');
		expect(body).toContain('/still.jpg');
	});

	it('renders the empty overview fallback when metadata is sparse', () => {
		const { body } = render(MediaEpisodeRow, {
			props: {
				title: 'S01E02 - Unknown',
				episode: episode({ overview: undefined, stillPath: undefined, airDate: undefined })
			}
		});

		expect(body).toContain('S01E02 - Unknown');
		expect(body).toContain('No episode overview available.');
		expect(body).not.toContain('<img');
	});
});

function episode(overrides: Partial<MediaMetadataEpisode> = {}): MediaMetadataEpisode {
	return {
		name: 'Pilot',
		episodeNumber: 1,
		overview: 'Episode overview',
		airDate: '2026-07-03',
		stillPath: '/still.jpg',
		...overrides
	} as MediaMetadataEpisode;
}
