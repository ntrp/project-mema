import type { MediaType } from './types';

export function providerPageUrl(
	provider?: string,
	type?: MediaType,
	externalId?: string,
	externalUrl?: string
) {
	if (externalUrl) {
		return externalUrl;
	}
	if (!provider || !type || !externalId) {
		return undefined;
	}
	if (provider === 'tmdb') {
		return `https://www.themoviedb.org/${type === 'movie' ? 'movie' : 'tv'}/${externalId}`;
	}
	if (provider === 'tvdb') {
		return `https://thetvdb.com/dereferrer/${type === 'movie' ? 'movie' : 'series'}/${externalId}`;
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
