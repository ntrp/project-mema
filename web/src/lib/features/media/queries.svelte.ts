import { createQuery } from '@tanstack/svelte-query';
import type { MediaType, MetadataProviderType } from '$lib/settings/types';
import { getMediaCollection, getMediaMetadataDetails, getPersonDetails } from './api';

export const mediaDetailKeys = {
	all: ['media-detail'] as const,
	metadata: (provider?: string, type?: string, id?: string) =>
		[...mediaDetailKeys.all, 'metadata', provider, type, id] as const,
	person: (provider?: string, id?: string) =>
		[...mediaDetailKeys.all, 'person', provider, id] as const,
	collection: (provider?: string, id?: string) =>
		[...mediaDetailKeys.all, 'collection', provider, id] as const
};

export function createMetadataDetailQuery(params: {
	provider: () => string | undefined;
	type: () => string | undefined;
	id: () => string | undefined;
}) {
	return createQuery(() => {
		const provider = params.provider();
		const type = params.type();
		const id = params.id();
		return {
			queryKey: mediaDetailKeys.metadata(provider, type, id),
			enabled: Boolean(provider && type && id),
			queryFn: ({ signal }) =>
				getMediaMetadataDetails(provider as MetadataProviderType, type as MediaType, id!, signal)
		};
	});
}

export function createPersonDetailQuery(
	provider: () => string | undefined,
	id: () => string | undefined
) {
	return createQuery(() => ({
		queryKey: mediaDetailKeys.person(provider(), id()),
		enabled: Boolean(provider() && id()),
		queryFn: ({ signal }) => getPersonDetails(provider() as MetadataProviderType, id()!, signal)
	}));
}

export function createMediaCollectionQuery(
	provider: () => string | undefined,
	id: () => string | undefined
) {
	return createQuery(() => ({
		queryKey: mediaDetailKeys.collection(provider(), id()),
		enabled: Boolean(provider() && id()),
		queryFn: ({ signal }) => getMediaCollection(provider() as MetadataProviderType, id()!, signal)
	}));
}
