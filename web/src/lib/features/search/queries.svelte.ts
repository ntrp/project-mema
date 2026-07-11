import { createQuery } from '@tanstack/svelte-query';
import { advancedSearchMedia, autocompleteMedia, type MediaAdvancedSearchRequest } from './api';

export const searchKeys = {
	all: ['search'] as const,
	autocomplete: (query: string) => [...searchKeys.all, 'autocomplete', query] as const,
	advanced: (request?: MediaAdvancedSearchRequest) =>
		[...searchKeys.all, 'advanced', request] as const
};

export function createAutocompleteQuery(query: () => string, enabled: () => boolean) {
	return createQuery(() => ({
		queryKey: searchKeys.autocomplete(query()),
		queryFn: ({ signal }) =>
			autocompleteMedia(
				{ query: query(), includeLibrary: true, includeProviders: false },
				{ signal }
			),
		select: (response) => response.groups ?? [],
		enabled: enabled() && query().length >= 2
	}));
}

export function createAdvancedSearchQuery(
	request: () => MediaAdvancedSearchRequest | undefined,
	enabled: () => boolean
) {
	return createQuery(() => ({
		queryKey: searchKeys.advanced(request()),
		queryFn: ({ signal }) => advancedSearchMedia(request() ?? {}, { signal }),
		select: (response) => response.groups ?? [],
		enabled: enabled() && Boolean(request())
	}));
}
