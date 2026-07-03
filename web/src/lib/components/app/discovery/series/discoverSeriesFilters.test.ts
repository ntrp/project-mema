import { describe, expect, it } from 'vitest';

import {
	activeSeriesFilterCount,
	defaultSeriesFilters,
	filtersFromParams,
	nextSeriesSort,
	seriesFilterUrl,
	seriesQuery
} from './discoverSeriesFilters';

describe('discover series filters', () => {
	it('uses full range defaults when range params are absent', () => {
		const filters = filtersFromParams(new URLSearchParams());

		expect(filters.runtime).toEqual([0, 400]);
		expect(filters.score).toEqual([0, 10]);
		expect(filters.minVoteCount).toBe(10);
		expect(seriesFilterUrl(defaultSeriesFilters())).toBe('/discover/series?minVoteCount=10');
		expect(seriesQuery(defaultSeriesFilters())).toMatchObject({ minVoteCount: 10 });
	});

	it('counts status as an active filter group without sorting', () => {
		const filters = {
			...defaultSeriesFilters(),
			sort: 'name.asc',
			status: ['ended'],
			genres: ['Animation'],
			withoutGenres: ['Reality']
		};

		expect(activeSeriesFilterCount(filters)).toBe(3);
	});

	it('serializes negated genre and keyword filters', () => {
		const filters = {
			...defaultSeriesFilters(),
			withoutGenres: ['Reality'],
			withoutKeywords: ['soap']
		};

		expect(seriesFilterUrl(filters)).toBe(
			'/discover/series?withoutGenres=Reality&withoutKeywords=soap&minVoteCount=10'
		);
		expect(seriesQuery(filters)).toMatchObject({
			withoutGenres: ['Reality'],
			withoutKeywords: ['soap']
		});
	});

	it('drops negated filters that are already included', () => {
		const filters = filtersFromParams(
			new URLSearchParams(
				'genres=Animation&withoutGenres=Animation&keywords=detective&withoutKeywords=detective'
			)
		);

		expect(filters.genres).toEqual(['Animation']);
		expect(filters.withoutGenres).toEqual([]);
		expect(filters.keywords).toEqual(['detective']);
		expect(filters.withoutKeywords).toEqual([]);
	});

	it('uses series sort defaults and reverses the current row', () => {
		expect(nextSeriesSort('popularity.desc', 'name')).toBe('name.asc');
		expect(nextSeriesSort('name.asc', 'name')).toBe('name.desc');
		expect(nextSeriesSort('popularity.desc', 'first_air_date')).toBe('first_air_date.desc');
	});
});
