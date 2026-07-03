import { describe, expect, it } from 'vitest';

import {
	discoverResultKey,
	filterDiscoverSection,
	filterDiscoverSections,
	relatedSectionFromDetail,
	sameDiscoverBlacklistItem
} from './discoverFilters';
import type {
	DiscoverBlacklistItem,
	MediaDiscoverSection,
	MediaMetadataDetails,
	MediaSearchResult
} from '$lib/settings/types';

const matrix = {
	type: 'movie',
	title: 'The Matrix',
	year: 1999,
	externalProvider: 'tmdb',
	externalId: '603'
} as MediaSearchResult;
const blacklist = [{ ...matrix, id: 'blacklist-1' }] as DiscoverBlacklistItem[];

describe('discovery filtering (SCN-MEDIA-004)', () => {
	it('uses external keys before title fallback', () => {
		expect(discoverResultKey(matrix)).toBe('movie:tmdb:603:The Matrix:1999');
		expect(sameDiscoverBlacklistItem(blacklist[0], matrix)).toBe(true);
		expect(
			sameDiscoverBlacklistItem(
				{ type: 'movie', title: ' the matrix ', year: 1999 } as DiscoverBlacklistItem,
				{ type: 'movie', title: 'The Matrix', year: 1999 } as MediaSearchResult
			)
		).toBe(true);
	});

	it('filters discover sections by blacklist entries', () => {
		const section = {
			id: 'trending',
			title: 'Trending',
			results: [matrix, { type: 'movie', title: 'Arrival', year: 2016 } as MediaSearchResult]
		} as MediaDiscoverSection;

		expect(filterDiscoverSection(section, blacklist).results.map((item) => item.title)).toEqual([
			'Arrival'
		]);
		expect(filterDiscoverSections([section], blacklist)[0].results).toHaveLength(1);
	});

	it('derives related sections from metadata details', () => {
		const detail = {
			type: 'movie',
			externalProvider: 'tmdb',
			recommendations: [
				matrix,
				{ type: 'movie', title: 'Dark City', year: 1998 } as MediaSearchResult
			],
			similar: [{ type: 'movie', title: 'Equilibrium', year: 2002 } as MediaSearchResult]
		} as MediaMetadataDetails;

		expect(relatedSectionFromDetail(undefined, 'similar', blacklist)).toBeUndefined();
		expect(relatedSectionFromDetail(detail, 'recommendations', blacklist)).toMatchObject({
			id: 'recommendations',
			title: 'Recommendations',
			providerName: 'TMDB',
			results: [{ title: 'Dark City' }]
		});
		expect(relatedSectionFromDetail(detail, 'similar', [])?.title).toBe('Similar Movies');
	});
});
