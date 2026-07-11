import type { DiscoverMovieSearchQuery } from '$lib/settings/api';

export interface DiscoverMovieFilters {
	sort: string;
	releaseDateFrom: string;
	releaseDateTo: string;
	studios: string[];
	genres: string[];
	keywords: string[];
	withoutGenres: string[];
	withoutKeywords: string[];
	originalLanguages: string[];
	contentRatings: string[];
	runtime: [number, number];
	score: [number, number];
	minVoteCount: number;
}

export {
	movieSortDirection,
	movieSortKey,
	movieSortOptions,
	nextMovieSort,
	type MovieSortDirection,
	type MovieSortKey
} from './discoverMovieSort';

export function activeMovieFilterCount(filters: DiscoverMovieFilters) {
	const defaults = defaultMovieFilters();
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
		rangeChanged(filters.runtime, defaults.runtime),
		rangeChanged(filters.score, defaults.score),
		filters.minVoteCount !== defaults.minVoteCount
	].filter(Boolean).length;
}

function rangeChanged(left: [number, number], right: [number, number]) {
	return left[0] !== right[0] || left[1] !== right[1];
}
export const languageOptions = [
	{ value: 'en', label: 'English' },
	{ value: 'de', label: 'German' },
	{ value: 'fr', label: 'French' },
	{ value: 'es', label: 'Spanish' },
	{ value: 'it', label: 'Italian' },
	{ value: 'ja', label: 'Japanese' },
	{ value: 'ko', label: 'Korean' },
	{ value: 'zh', label: 'Chinese' },
	{ value: 'hi', label: 'Hindi' }
];

export const contentRatingOptions = ['G', 'PG', 'PG-13', 'R', 'NC-17', 'NR'].map((value) => ({
	value,
	label: value
}));

export const defaultMinVoteCount = 10;

export function defaultMovieFilters(): DiscoverMovieFilters {
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
		runtime: [0, 400],
		score: [0, 10],
		minVoteCount: defaultMinVoteCount
	};
}

export function filtersFromParams(params: URLSearchParams): DiscoverMovieFilters {
	const genres = params.getAll('genres');
	const keywords = params.getAll('keywords');
	return {
		...defaultMovieFilters(),
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
		runtime: [numberParam(params, 'runtimeMin', 0), numberParam(params, 'runtimeMax', 400)],
		score: [numberParam(params, 'scoreMin', 0), numberParam(params, 'scoreMax', 10)],
		minVoteCount: numberParam(params, 'minVoteCount', defaultMinVoteCount)
	};
}

export function movieQuery(filters: DiscoverMovieFilters, page = 1): DiscoverMovieSearchQuery {
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
		runtimeMin: filters.runtime[0],
		runtimeMax: filters.runtime[1],
		scoreMin: filters.score[0],
		scoreMax: filters.score[1],
		minVoteCount: filters.minVoteCount
	});
}

export function movieFilterUrl(filters: DiscoverMovieFilters) {
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
	setParam(params, 'runtimeMin', String(filters.runtime[0]), '0');
	setParam(params, 'runtimeMax', String(filters.runtime[1]), '400');
	setParam(params, 'scoreMin', String(filters.score[0]), '0');
	setParam(params, 'scoreMax', String(filters.score[1]), '10');
	params.set('minVoteCount', String(filters.minVoteCount));
	const query = params.toString();
	return query ? `/discover/movies?${query}` : '/discover/movies';
}

function numberParam(params: URLSearchParams, key: string, fallback: number) {
	if (!params.has(key)) {
		return fallback;
	}
	const value = Number(params.get(key));
	return Number.isFinite(value) ? value : fallback;
}

function appendParams(params: URLSearchParams, key: string, values: string[]) {
	for (const value of values.filter(Boolean)) params.append(key, value);
}

function withoutIncluded(values: string[], included: string[]) {
	return values.filter((value) => !included.includes(value));
}

function setParam(params: URLSearchParams, key: string, value: string, defaultValue = '') {
	if (value && value !== defaultValue) params.set(key, value);
}

function prune(query: DiscoverMovieSearchQuery): DiscoverMovieSearchQuery {
	return Object.fromEntries(
		Object.entries(query).filter(
			([, value]) => value !== '' && (!Array.isArray(value) || value.length)
		)
	);
}
