import { client } from '$lib/api/client';
import type {
	DiscoverBlacklistItem,
	DiscoverBlacklistRequest,
	DiscoverMovieFacetOption,
	DiscoverMovieSearchResponse,
	MediaSearchRequest
} from '../types';

export async function searchMedia(request: MediaSearchRequest) {
	const { data, error } = await client.POST('/media/search', { body: request });

	if (error) {
		throw new Error(error.message);
	}
	return data?.results ?? [];
}

export async function loadMediaDiscoverSections() {
	const { data, error } = await client.GET('/media/discover');

	if (error) {
		throw new Error(error.message);
	}
	return data?.sections ?? [];
}

export async function loadMediaDiscoverSection(sectionId: string, page = 1, limit = 20) {
	const { data, error } = await client.GET('/media/discover/{sectionId}', {
		params: { path: { sectionId }, query: { page, limit } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Discovery section request did not return a result');
	}
	return data;
}

export interface DiscoverMovieSearchQuery {
	sort?: string;
	page?: number;
	releaseDateFrom?: string;
	releaseDateTo?: string;
	studios?: string[];
	genres?: string[];
	keywords?: string[];
	withoutGenres?: string[];
	withoutKeywords?: string[];
	originalLanguages?: string[];
	contentRatings?: string[];
	runtimeMin?: number;
	runtimeMax?: number;
	scoreMin?: number;
	scoreMax?: number;
	minVoteCount?: number;
}

export interface DiscoverSeriesSearchQuery extends DiscoverMovieSearchQuery {
	status?: string[];
}

export async function searchDiscoverMovies(
	query: DiscoverMovieSearchQuery
): Promise<DiscoverMovieSearchResponse> {
	const { data, error } = await client.GET('/media/discover/movies/search', {
		params: { query }
	});

	if (error) {
		throw new Error(error.message);
	}
	return data ?? { results: [], hasMore: false };
}

export async function searchDiscoverSeries(
	query: DiscoverSeriesSearchQuery
): Promise<DiscoverMovieSearchResponse> {
	const { data, error } = await client.GET('/media/discover/series/search', {
		params: { query }
	});

	if (error) {
		throw new Error(error.message);
	}
	return data ?? { results: [], hasMore: false };
}

export async function autocompleteDiscoverMovieFacet(
	facet: 'genres' | 'studios' | 'keywords',
	query: string
): Promise<DiscoverMovieFacetOption[]> {
	const { data, error } = await client.GET('/media/discover/movies/facets/{facet}', {
		params: { path: { facet }, query: { query } }
	});

	if (error) {
		throw new Error(error.message);
	}
	return data?.options ?? [];
}

export async function autocompleteDiscoverSeriesFacet(
	facet: 'genres' | 'studios' | 'keywords',
	query: string
): Promise<DiscoverMovieFacetOption[]> {
	const { data, error } = await client.GET('/media/discover/series/facets/{facet}', {
		params: { path: { facet }, query: { query } }
	});

	if (error) {
		throw new Error(error.message);
	}
	return data?.options ?? [];
}

export async function listDiscoverBlacklist(): Promise<DiscoverBlacklistItem[]> {
	const { data, error } = await client.GET('/media/discover/blacklist');

	if (error) {
		throw new Error(error.message);
	}
	return data?.items ?? [];
}

export async function addDiscoverBlacklistItem(
	request: DiscoverBlacklistRequest
): Promise<DiscoverBlacklistItem> {
	const { data, error } = await client.POST('/media/discover/blacklist', { body: request });

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Blacklist request did not return a result');
	}
	return data;
}

export async function deleteDiscoverBlacklistItem(id: string) {
	const { error } = await client.DELETE('/media/discover/blacklist/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}
