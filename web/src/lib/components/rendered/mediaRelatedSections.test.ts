import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import MediaRelatedSections from '$lib/components/app/media/posters/MediaRelatedSections.svelte';
import type { MediaItem, MediaMetadataDetails, MediaSearchResult } from '$lib/settings/types';

describe('rendered media related sections (SCN-MEDIA-008)', () => {
	it('renders recommendations, similar media, detail links, and library state', () => {
		const { body } = render(MediaRelatedSections, {
			props: {
				detail: metadataDetail(),
				mediaItems: [
					mediaItem({
						title: 'Owned Recommendation',
						externalProvider: 'tmdb',
						externalId: 'owned-1'
					})
				],
				addingKey: 'movie:New Recommendation:2027',
				actionLabel: 'Add to library',
				onAdd: vi.fn()
			}
		});

		expect(body).toContain('Recommendations');
		expect(body).toContain('Similar Movies');
		expect(body).toContain('Owned Recommendation');
		expect(body).toContain('New Recommendation');
		expect(body).toContain('Similar Candidate');
		expect(body).toContain('/media/tmdb/movie/parent-1/recommendations');
		expect(body).toContain('/media/tmdb/movie/parent-1/similar');
		expect(body).toContain('/media/tmdb/movie/owned-1');
		expect(body).toContain('Open Recommendations');
		expect(body).toContain('Open Similar Movies');
		expect(body).toContain('View all');
		expect(body).toContain('In library');
		expect(body).toContain('Working');
	});

	it('uses series labels and omits empty sections', () => {
		const { body } = render(MediaRelatedSections, {
			props: {
				detail: metadataDetail({
					type: 'serie',
					externalId: 'series-1',
					recommendations: [],
					similar: [searchResult({ title: 'Similar Series Candidate', type: 'serie' })]
				}),
				actionLabel: 'Add series',
				onAdd: vi.fn()
			}
		});

		expect(body).not.toContain('Recommendations');
		expect(body).toContain('Similar Series');
		expect(body).toContain('Similar Series Candidate');
		expect(body).toContain('/media/tmdb/serie/series-1/similar');
	});
});

function metadataDetail(overrides: Partial<MediaMetadataDetails> = {}): MediaMetadataDetails {
	return {
		title: 'Scenario Parent',
		type: 'movie',
		year: 2026,
		externalProvider: 'tmdb',
		externalId: 'parent-1',
		recommendations: [
			searchResult({ title: 'Owned Recommendation', externalId: 'owned-1' }),
			searchResult({ title: 'New Recommendation', externalId: 'new-1', year: 2027 })
		],
		similar: [searchResult({ title: 'Similar Candidate', externalId: 'similar-1' })],
		...overrides
	} as MediaMetadataDetails;
}

function searchResult(overrides: Partial<MediaSearchResult> = {}): MediaSearchResult {
	return {
		title: 'Scenario Result',
		type: 'movie',
		year: 2026,
		externalProvider: 'tmdb',
		externalId: 'result-1',
		overview: 'A related scenario result.',
		...overrides
	} as MediaSearchResult;
}

function mediaItem(overrides: Partial<MediaItem> = {}): MediaItem {
	return {
		id: 'media-1',
		title: 'Scenario Movie',
		type: 'movie',
		year: 2026,
		status: 'missing',
		monitored: true,
		monitorMode: 'only_media',
		minimumAvailability: 'released',
		filePaths: [],
		metadataFilePaths: [],
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:01:00Z',
		...overrides
	} as MediaItem;
}
