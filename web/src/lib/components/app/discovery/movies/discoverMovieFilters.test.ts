import { describe, expect, it } from 'vitest';

import {
	activeMovieFilterCount,
	defaultMovieFilters,
	filtersFromParams,
	movieFilterUrl,
	movieQuery,
	nextMovieSort
} from './discoverMovieFilters';

describe('discover movie filters', () => {
	it('uses full range defaults when range params are absent', () => {
		const filters = filtersFromParams(new URLSearchParams());

		expect(filters.runtime).toEqual([0, 400]);
		expect(filters.score).toEqual([0, 10]);
		expect(filters.minVoteCount).toBe(0);
		expect(movieFilterUrl(defaultMovieFilters())).toBe('/discover/movies');
	});

	it('parses explicit numeric range params', () => {
		const params = new URLSearchParams(
			'runtimeMin=80&runtimeMax=120&scoreMin=6.5&scoreMax=9&minVoteCount=250'
		);
		const filters = filtersFromParams(params);

		expect(filters.runtime).toEqual([80, 120]);
		expect(filters.score).toEqual([6.5, 9]);
		expect(filters.minVoteCount).toBe(250);
	});

	it('counts active filter groups without sorting', () => {
		const filters = {
			...defaultMovieFilters(),
			sort: 'title.asc',
			genres: ['Drama', 'Comedy'],
			withoutKeywords: ['zombie'],
			releaseDateFrom: '2026-01-01',
			runtime: [90, 120] as [number, number]
		};

		expect(activeMovieFilterCount(filters)).toBe(4);
	});

	it('serializes negated genre and keyword filters', () => {
		const filters = {
			...defaultMovieFilters(),
			withoutGenres: ['Horror'],
			withoutKeywords: ['ghost']
		};

		expect(movieFilterUrl(filters)).toBe(
			'/discover/movies?withoutGenres=Horror&withoutKeywords=ghost'
		);
		expect(movieQuery(filters)).toMatchObject({
			withoutGenres: ['Horror'],
			withoutKeywords: ['ghost']
		});
	});

	it('drops negated filters that are already included', () => {
		const filters = filtersFromParams(
			new URLSearchParams('genres=Drama&withoutGenres=Drama&keywords=hero&withoutKeywords=hero')
		);

		expect(filters.genres).toEqual(['Drama']);
		expect(filters.withoutGenres).toEqual([]);
		expect(filters.keywords).toEqual(['hero']);
		expect(filters.withoutKeywords).toEqual([]);
	});

	it('uses default sort direction for new rows and reverses the current row', () => {
		expect(nextMovieSort('popularity.desc', 'title')).toBe('title.asc');
		expect(nextMovieSort('title.asc', 'title')).toBe('title.desc');
		expect(nextMovieSort('title.desc', 'title')).toBe('title.asc');
	});
});
