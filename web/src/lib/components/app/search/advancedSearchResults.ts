import { resolve } from '$app/paths';
import { providerPageUrl } from '$lib/settings/providerLinks';
import type { MediaSearchResult, PersonSearchResult } from '$lib/settings/types';

export function imageUrl(path?: string) {
	if (!path) {
		return undefined;
	}
	if (path.startsWith('http://') || path.startsWith('https://')) {
		return path;
	}
	return `https://image.tmdb.org/t/p/w185${path}`;
}

export function mediaResultKey(result: MediaSearchResult) {
	return `${result.id ?? ''}:${result.type}:${result.externalProvider ?? ''}:${result.externalId ?? ''}:${result.title}:${result.year ?? ''}`;
}

export function mediaCandidateKey(candidate: MediaSearchResult) {
	return `${candidate.type}:${candidate.title}:${candidate.year ?? ''}`;
}

export function mediaHref(result: MediaSearchResult) {
	if (result.id) {
		return result.type === 'movie'
			? resolve('/movies/[id]', { id: result.id })
			: resolve('/series/[id]', { id: result.id });
	}
	if (result.externalProvider && result.externalId) {
		return resolve('/media/[provider]/[type]/[externalId]', {
			provider: result.externalProvider,
			type: result.type,
			externalId: result.externalId
		});
	}
	return undefined;
}

export function externalMediaUrl(result: MediaSearchResult) {
	return providerPageUrl(
		result.externalProvider,
		result.type,
		result.externalId,
		result.externalUrl
	);
}

export function personResultKey(result: PersonSearchResult) {
	return `${result.externalProvider}:${result.externalId}:${result.name}`;
}

export function personHref(result: PersonSearchResult) {
	return resolve('/people/[provider]/[personId]', {
		provider: result.externalProvider,
		personId: result.externalId
	});
}
