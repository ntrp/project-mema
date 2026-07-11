export const movieSortOptions = [
	{ key: 'popularity', label: 'Popularity', defaultDirection: 'desc' },
	{ key: 'release_date', label: 'Release date', defaultDirection: 'desc' },
	{ key: 'vote_average', label: 'TMDB score', defaultDirection: 'desc' },
	{ key: 'title', label: 'Name', defaultDirection: 'asc' }
] as const;

export type MovieSortKey = (typeof movieSortOptions)[number]['key'];
export type MovieSortDirection = 'asc' | 'desc';

const defaults = new Map<MovieSortKey, MovieSortDirection>(
	movieSortOptions.map((option) => [option.key, option.defaultDirection])
);

export function movieSortKey(sort: string): MovieSortKey {
	const [key] = split(sort);
	return defaults.has(key as MovieSortKey) ? (key as MovieSortKey) : 'popularity';
}

export function movieSortDirection(sort: string): MovieSortDirection {
	const [, direction] = split(sort);
	return direction === 'asc' ? 'asc' : 'desc';
}

export function nextMovieSort(currentSort: string, key: MovieSortKey) {
	if (movieSortKey(currentSort) === key) {
		return `${key}.${movieSortDirection(currentSort) === 'asc' ? 'desc' : 'asc'}`;
	}
	return `${key}.${defaults.get(key) ?? 'desc'}`;
}

function split(sort: string) {
	const [key, direction] = sort.split('.');
	return [key ?? 'popularity', direction ?? 'desc'] as const;
}
