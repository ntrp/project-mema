import type { MediaType } from './types';

export function providerPageUrl(provider?: string, type?: MediaType, externalId?: string) {
	if (!provider || !type || !externalId) {
		return undefined;
	}
	if (provider === 'tmdb') {
		return `https://www.themoviedb.org/${type === 'movie' ? 'movie' : 'tv'}/${externalId}`;
	}
	if (provider === 'tvdb') {
		return `https://thetvdb.com/${type === 'movie' ? 'movies' : 'series'}/${externalId}`;
	}
	return undefined;
}

export function providerDisplayName(provider?: string) {
	if (provider === 'tmdb') {
		return 'TMDB';
	}
	if (provider === 'tvdb') {
		return 'TVDB';
	}
	return provider?.toUpperCase() ?? 'Provider';
}
