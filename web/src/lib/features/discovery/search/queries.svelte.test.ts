import { beforeEach, describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
	createInfiniteQuery: vi.fn((options: () => unknown) => options()),
	createQuery: vi.fn((options: () => unknown) => options()),
	searchMovies: vi.fn(),
	searchSeries: vi.fn(),
	movieFacet: vi.fn(),
	seriesFacet: vi.fn()
}));

vi.mock('@tanstack/svelte-query', () => ({
	createInfiniteQuery: mocks.createInfiniteQuery,
	createQuery: mocks.createQuery
}));
vi.mock('$lib/settings/api', () => ({
	searchDiscoverMovies: mocks.searchMovies,
	searchDiscoverSeries: mocks.searchSeries,
	autocompleteDiscoverMovieFacet: mocks.movieFacet,
	autocompleteDiscoverSeriesFacet: mocks.seriesFacet
}));

import {
	createDiscoverFacetQuery,
	createMovieSearchQuery,
	createSeriesSearchQuery,
	discoverSearchKeys
} from './queries.svelte';

describe('discovery search queries', () => {
	beforeEach(() => vi.clearAllMocks());

	it('uses stable kind-specific keys and query-owned pagination', async () => {
		createMovieSearchQuery(() => ({ sort: 'popularity' }));
		createSeriesSearchQuery(() => ({ sort: 'rating' }));
		const [movies, series] = mocks.createInfiniteQuery.mock.results.map(
			(result) => result.value as CapturedInfiniteOptions
		);
		expect(movies.queryKey).toEqual(discoverSearchKeys.results('movies', { sort: 'popularity' }));
		await movies.queryFn({ pageParam: 3 });
		expect(mocks.searchMovies).toHaveBeenCalledWith({ sort: 'popularity', page: 3 });
		expect(series.queryKey[1]).toBe('series');
		expect(movies.getNextPageParam({ hasMore: true }, [{}, {}])).toBe(3);
		expect(movies.getNextPageParam({ hasMore: false }, [{}])).toBeUndefined();
	});

	it('does not fetch facets until two trimmed characters are present', async () => {
		let input = ' a ';
		createDiscoverFacetQuery('movies', 'genres', () => input);
		const facet = mocks.createQuery.mock.results[0].value as CapturedQueryOptions;
		expect(facet.enabled).toBe(false);
		input = ' sci ';
		createDiscoverFacetQuery('movies', 'genres', () => input);
		const active = mocks.createQuery.mock.results[1].value as CapturedQueryOptions;
		expect(active.enabled).toBe(true);
		await active.queryFn();
		expect(mocks.movieFacet).toHaveBeenCalledWith('genres', 'sci');
	});
});

interface CapturedQueryOptions {
	queryKey: readonly unknown[];
	enabled: boolean;
	queryFn: () => Promise<unknown>;
}

interface CapturedInfiniteOptions extends Omit<CapturedQueryOptions, 'queryFn'> {
	queryFn: (context: { pageParam: number }) => Promise<unknown>;
	getNextPageParam: (last: { hasMore: boolean }, pages: unknown[]) => number | undefined;
}
