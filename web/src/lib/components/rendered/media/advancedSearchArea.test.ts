import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import AdvancedSearchArea from '$lib/components/app/search/AdvancedSearchArea.svelte';
import type { MediaSearchGroup, MediaSearchResult, MetadataProvider } from '$lib/settings/types';

describe('rendered advanced search area (SCN-MEDIA-008)', () => {
	it('renders search filters, enabled providers, result groups, and external links', () => {
		const { body } = render(AdvancedSearchArea, {
			props: {
				initialQuery: 'scenario',
				metadataProviders: [
					metadataProvider({ name: 'TMDb' }),
					metadataProvider({ id: 'disabled-1', name: 'Disabled Provider', enabled: false })
				],
				groups: searchGroups(),
				searching: true,
				addingKey: 'movie:Provider Result:2026',
				actionLabel: 'Request',
				onSearch: vi.fn(),
				onAdd: vi.fn()
			}
		});

		expect(body).toContain('Advanced search');
		expect(body).toContain('Title');
		expect(body).toContain('Search');
		expect(body).toContain('Year');
		expect(body).toContain('Metadata providers');
		expect(body).toContain('TMDb');
		expect(body).not.toContain('Disabled Provider');
		expect(body).toContain('2 results');
		expect(body).toContain('Existing Library Result');
		expect(body).toContain('/movies/media-1');
		expect(body).toContain('In library');
		expect(body).toContain('Provider Result');
		expect(body).toContain('/media/tmdb/movie/provider-1');
		expect(body).toContain('Open TMDB page in a new tab');
		expect(body).toContain('Working');
	});
});

function searchGroups(): MediaSearchGroup[] {
	return [
		{
			sourceType: 'library',
			sourceName: 'Library',
			results: [searchResult({ id: 'media-1', title: 'Existing Library Result' })]
		},
		{
			sourceType: 'provider',
			sourceName: 'TMDb',
			results: [searchResult({ title: 'Provider Result', externalId: 'provider-1' })]
		}
	];
}

function searchResult(overrides: Partial<MediaSearchResult> = {}): MediaSearchResult {
	return {
		title: 'Scenario Result',
		type: 'movie',
		year: 2026,
		externalProvider: 'tmdb',
		externalId: 'result-1',
		overview: 'Scenario overview',
		...overrides
	} as MediaSearchResult;
}

function metadataProvider(overrides: Partial<MetadataProvider> = {}): MetadataProvider {
	return {
		id: 'provider-1',
		name: 'Scenario Provider',
		type: 'tmdb',
		enabled: true,
		baseUrl: 'https://metadata.test',
		apiKey: 'secret',
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:01:00Z',
		...overrides
	} as MetadataProvider;
}
