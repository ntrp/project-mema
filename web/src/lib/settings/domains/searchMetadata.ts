import { client } from '$lib/api/client';
import type { MediaAdvancedSearchRequest, MediaType, MetadataProviderType } from '../types';

export type AutocompleteSearchScope = 'library' | 'providers' | 'all';

export async function autocompleteMedia(query: string, scope: AutocompleteSearchScope = 'all') {
	const { data, error } = await client.GET('/media/autocomplete', {
		params: {
			query: {
				query,
				includeLibrary: scope !== 'providers',
				includeProviders: scope !== 'library'
			}
		}
	});

	if (error) {
		throw new Error(error.message);
	}
	return data?.groups ?? [];
}

export async function advancedSearchMedia(request: MediaAdvancedSearchRequest) {
	const { data, error } = await client.POST('/media/advanced-search', { body: request });

	if (error) {
		throw new Error(error.message);
	}
	return data?.groups ?? [];
}

export async function getMediaMetadataDetails(
	provider: MetadataProviderType,
	type: MediaType,
	externalId: string
) {
	const { data, error } = await client.GET('/media/metadata/{provider}/{type}/{externalId}', {
		params: { path: { provider, type, externalId } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media details were not returned');
	}
	return data;
}

export async function getPersonDetails(provider: MetadataProviderType, personId: string) {
	const { data, error } = await client.GET('/people/{provider}/{personId}', {
		params: { path: { provider, personId } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Person details were not returned');
	}
	return data;
}

export async function getMediaCollection(provider: MetadataProviderType, collectionId: string) {
	const { data, error } = await client.GET('/media/collections/{provider}/{collectionId}', {
		params: { path: { provider, collectionId } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media collection was not returned');
	}
	return data;
}
