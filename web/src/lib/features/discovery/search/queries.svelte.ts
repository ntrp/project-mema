import { createInfiniteQuery, createQuery } from '@tanstack/svelte-query';
import {
	autocompleteDiscoverMovieFacet,
	autocompleteDiscoverSeriesFacet,
	searchDiscoverMovies,
	searchDiscoverSeries,
	type DiscoverMovieSearchQuery,
	type DiscoverSeriesSearchQuery
} from '$lib/settings/api';

export type DiscoverFacet = 'genres' | 'studios' | 'keywords';
export type DiscoverKind = 'movies' | 'series';

export const discoverSearchKeys = {
	all: ['discover-search'] as const,
	results: (kind: DiscoverKind, query: object) =>
		[...discoverSearchKeys.all, kind, 'results', query] as const,
	facet: (kind: DiscoverKind, facet: DiscoverFacet, query: string) =>
		[...discoverSearchKeys.all, kind, 'facet', facet, query] as const
};

export function createMovieSearchQuery(query: () => DiscoverMovieSearchQuery) {
	return createInfiniteQuery(() => ({
		queryKey: discoverSearchKeys.results('movies', query()),
		queryFn: ({ pageParam }) => searchDiscoverMovies({ ...query(), page: pageParam }),
		initialPageParam: 1,
		getNextPageParam: (last, pages) => (last.hasMore ? pages.length + 1 : undefined)
	}));
}

export function createSeriesSearchQuery(query: () => DiscoverSeriesSearchQuery) {
	return createInfiniteQuery(() => ({
		queryKey: discoverSearchKeys.results('series', query()),
		queryFn: ({ pageParam }) => searchDiscoverSeries({ ...query(), page: pageParam }),
		initialPageParam: 1,
		getNextPageParam: (last, pages) => (last.hasMore ? pages.length + 1 : undefined)
	}));
}

export function createDiscoverFacetQuery(
	kind: DiscoverKind,
	facet: DiscoverFacet,
	query: () => string
) {
	return createQuery(() => {
		const value = query().trim();
		return {
			queryKey: discoverSearchKeys.facet(kind, facet, value),
			enabled: value.length >= 2,
			queryFn: () =>
				kind === 'movies'
					? autocompleteDiscoverMovieFacet(facet, value)
					: autocompleteDiscoverSeriesFacet(facet, value)
		};
	});
}
