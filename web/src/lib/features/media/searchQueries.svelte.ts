import { createQuery } from '@tanstack/svelte-query';
import { getMediaMetadataDetails, searchMedia } from '$lib/settings/api';
import { searchMediaSubtitles } from '$lib/features/releases/api';
import type { ManualSubtitleSearchRequest, MediaType } from '$lib/settings/types';

export const mediaSearchKeys = {
	all: ['media-search'] as const,
	lookup: (type: MediaType, query: string) =>
		[...mediaSearchKeys.all, 'lookup', type, query] as const,
	details: (provider: string, type: MediaType, id: string) =>
		[...mediaSearchKeys.all, 'details', provider, type, id] as const,
	subtitles: (id: string, request?: ManualSubtitleSearchRequest) =>
		[...mediaSearchKeys.all, 'subtitles', id, request] as const
};

export function createMediaLookupQuery(type: MediaType, query: () => string, enabled = () => true) {
	return createQuery(() => {
		const value = query().trim();
		return {
			queryKey: mediaSearchKeys.lookup(type, value),
			enabled: enabled() && value.length >= 2,
			queryFn: () => searchMedia({ query: value, type }),
			select: (results) => results.slice(0, 6)
		};
	});
}

export function createTmdbSeriesDetailsQuery(id: () => string | undefined) {
	return createQuery(() => ({
		queryKey: mediaSearchKeys.details('tmdb', 'serie', id() ?? ''),
		enabled: Boolean(id()),
		queryFn: () => getMediaMetadataDetails('tmdb', 'serie', id()!)
	}));
}

export function createSubtitleSearchQuery(
	id: () => string,
	request: () => ManualSubtitleSearchRequest | undefined
) {
	return createQuery(() => ({
		queryKey: mediaSearchKeys.subtitles(id(), request()),
		enabled: false,
		queryFn: () => searchMediaSubtitles(id(), request()!)
	}));
}
