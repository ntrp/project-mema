import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import AdvancedSearchResults from '$lib/components/app/search/AdvancedSearchResults.svelte';
import type { MediaSearchGroup } from '$lib/settings/types';

describe('rendered advanced search results (SCN-MEDIA-008)', () => {
	it('renders provider results, library status, and add affordances', () => {
		const { body } = render(AdvancedSearchResults, {
			props: {
				groups: searchGroups(),
				addingKey: 'movie:Remote Movie:2026',
				actionLabel: 'Add media',
				onAdd: vi.fn()
			}
		});

		expect(body).toContain('TMDb');
		expect(body).toContain('metadata');
		expect(body).toContain('Remote Movie');
		expect(body).toContain('Working');
		expect(body).toContain('Open TMDB page in a new tab');
		expect(body).toContain('Local Series');
		expect(body).toContain('In library');
	});
});

function searchGroups(): MediaSearchGroup[] {
	return [
		{
			sourceType: 'metadata',
			sourceName: 'TMDb',
			results: [
				{
					title: 'Remote Movie',
					type: 'movie',
					year: 2026,
					overview: 'A remote provider result',
					externalProvider: 'tmdb',
					externalId: 'movie-1',
					posterPath: '/poster.jpg'
				}
			]
		},
		{
			sourceType: 'library',
			sourceName: 'Library',
			results: [
				{
					id: 'series-1',
					title: 'Local Series',
					type: 'series',
					year: 2025,
					overview: 'Already in the library'
				}
			]
		}
	] as MediaSearchGroup[];
}
