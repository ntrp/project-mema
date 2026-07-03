import type { DiscoverSeriesSearchQuery } from '$lib/settings/api';
import {
	contentRatingOptions,
	languageOptions,
	type DiscoverMovieFilters
} from '../movies/discoverMovieFilters';

export interface DiscoverSeriesFilters extends DiscoverMovieFilters {
	status: string[];
}

export const seriesSortOptions = [
	{ key: 'popularity', label: 'Popularity', defaultDirection: 'desc' },
	{ key: 'first_air_date', label: 'First air date', defaultDirection: 'desc' },
	{ key: 'vote_average', label: 'TMDB score', defaultDirection: 'desc' },
	{ key: 'name', label: 'Name', defaultDirection: 'asc' }
] as const;

export type SeriesSortKey = (typeof seriesSortOptions)[number]['key'];
export type SeriesSortDirection = 'asc' | 'desc';

const seriesSortDefaults = new Map<SeriesSortKey, SeriesSortDirection>(
	seriesSortOptions.map((option) => [option.key, option.defaultDirection])
);

export const seriesStatusOptions = [
	{ value: 'returning', label: 'Returning' },
	{ value: 'planned', label: 'Planned' },
	{ value: 'in production', label: 'In production' },
	{ value: 'ended', label: 'Ended' },
	{ value: 'canceled', label: 'Canceled' },
	{ value: 'pilot', label: 'Pilot' }
];

export { contentRatingOptions, languageOptions };

export function defaultSeriesFilters(): DiscoverSeriesFilters {
	return {
		sort: 'popularity.desc',
		releaseDateFrom: '',
		releaseDateTo: '',
		studios: [],
		genres: [],
		keywords: [],
		withoutGenres: [],
		withoutKeywords: [],
		originalLanguages: [],
		contentRatings: [],
		status: [],
		runtime: [0, 400],
		score: [0, 10],
		minVoteCount: 0
	};
}

export function filtersFromParams(params: URLSearchParams): DiscoverSeriesFilters {
	const genres = params.getAll('genres');
	const keywords = params.getAll('keywords');
	return {
		...defaultSeriesFilters(),
		sort: params.get('sort') || 'popularity.desc',
		releaseDateFrom: params.get('releaseDateFrom') ?? '',
		releaseDateTo: params.get('releaseDateTo') ?? '',
		studios: params.getAll('studios'),
		genres,
		keywords,
		withoutGenres: withoutIncluded(params.getAll('withoutGenres'), genres),
		withoutKeywords: withoutIncluded(params.getAll('withoutKeywords'), keywords),
		originalLanguages: params.getAll('originalLanguages'),
		contentRatings: params.getAll('contentRatings'),
		status: params.getAll('status'),
		runtime: [numberParam(params, 'runtimeMin', 0), numberParam(params, 'runtimeMax', 400)],
		score: [numberParam(params, 'scoreMin', 0), numberParam(params, 'scoreMax', 10)],
		minVoteCount: numberParam(params, 'minVoteCount', 0)
	};
}

export function seriesQuery(filters: DiscoverSeriesFilters, page = 1): DiscoverSeriesSearchQuery {
	return prune({
		sort: filters.sort,
		page,
		releaseDateFrom: filters.releaseDateFrom,
		releaseDateTo: filters.releaseDateTo,
		studios: filters.studios,
		genres: filters.genres,
		keywords: filters.keywords,
		withoutGenres: filters.withoutGenres,
		withoutKeywords: filters.withoutKeywords,
		originalLanguages: filters.originalLanguages,
		contentRatings: filters.contentRatings,
		status: filters.status,
		runtimeMin: filters.runtime[0],
		runtimeMax: filters.runtime[1],
		scoreMin: filters.score[0],
		scoreMax: filters.score[1],
		minVoteCount: filters.minVoteCount
	});
}

export function seriesFilterUrl(filters: DiscoverSeriesFilters) {
	const params = new URLSearchParams();
	setParam(params, 'sort', filters.sort, 'popularity.desc');
	setParam(params, 'releaseDateFrom', filters.releaseDateFrom);
	setParam(params, 'releaseDateTo', filters.releaseDateTo);
	appendParams(params, 'studios', filters.studios);
	appendParams(params, 'genres', filters.genres);
	appendParams(params, 'keywords', filters.keywords);
	appendParams(params, 'withoutGenres', filters.withoutGenres);
	appendParams(params, 'withoutKeywords', filters.withoutKeywords);
	appendParams(params, 'originalLanguages', filters.originalLanguages);
	appendParams(params, 'contentRatings', filters.contentRatings);
	appendParams(params, 'status', filters.status);
	setParam(params, 'runtimeMin', String(filters.runtime[0]), '0');
	setParam(params, 'runtimeMax', String(filters.runtime[1]), '400');
	setParam(params, 'scoreMin', String(filters.score[0]), '0');
	setParam(params, 'scoreMax', String(filters.score[1]), '10');
	setParam(params, 'minVoteCount', String(filters.minVoteCount), '0');
	const query = params.toString();
	return query ? `/discover/series?${query}` : '/discover/series';
}

export function activeSeriesFilterCount(filters: DiscoverSeriesFilters) {
	const defaults = defaultSeriesFilters();
	return [
		filters.releaseDateFrom,
		filters.releaseDateTo,
		filters.studios.length > 0,
		filters.genres.length > 0,
		filters.keywords.length > 0,
		filters.withoutGenres.length > 0,
		filters.withoutKeywords.length > 0,
		filters.originalLanguages.length > 0,
		filters.contentRatings.length > 0,
		filters.status.length > 0,
		rangeChanged(filters.runtime, defaults.runtime),
		rangeChanged(filters.score, defaults.score),
		filters.minVoteCount !== defaults.minVoteCount
	].filter(Boolean).length;
}

export function seriesSortKey(sort: string): SeriesSortKey {
	const [key] = splitSeriesSort(sort);
	return seriesSortDefaults.has(key as SeriesSortKey) ? (key as SeriesSortKey) : 'popularity';
}

export function seriesSortDirection(sort: string): SeriesSortDirection {
	const [, direction] = splitSeriesSort(sort);
	return direction === 'asc' ? 'asc' : 'desc';
}

export function nextSeriesSort(currentSort: string, key: SeriesSortKey) {
	const currentKey = seriesSortKey(currentSort);
	if (currentKey === key) {
		return `${key}.${seriesSortDirection(currentSort) === 'asc' ? 'desc' : 'asc'}`;
	}
	return `${key}.${seriesSortDefaults.get(key) ?? 'desc'}`;
}

function splitSeriesSort(sort: string) {
	const [key, direction] = sort.split('.');
	return [key ?? 'popularity', direction ?? 'desc'] as const;
}

function numberParam(params: URLSearchParams, key: string, fallback: number) {
	if (!params.has(key)) return fallback;
	const value = Number(params.get(key));
	return Number.isFinite(value) ? value : fallback;
}

function appendParams(params: URLSearchParams, key: string, values: string[]) {
	for (const value of values.filter(Boolean)) params.append(key, value);
}

function withoutIncluded(values: string[], included: string[]) {
	const includedSet = new Set(included);
	return values.filter((value) => !includedSet.has(value));
}

function setParam(params: URLSearchParams, key: string, value: string, defaultValue = '') {
	if (value && value !== defaultValue) params.set(key, value);
}

function rangeChanged(left: [number, number], right: [number, number]) {
	return left[0] !== right[0] || left[1] !== right[1];
}

function prune(query: DiscoverSeriesSearchQuery): DiscoverSeriesSearchQuery {
	return Object.fromEntries(
		Object.entries(query).filter(
			([, value]) => value !== '' && (!Array.isArray(value) || value.length)
		)
	);
}
