import { client } from '$lib/api/client';
import type { MediaType, MetadataProviderType } from '$lib/settings/types';

export async function getMediaMetadataDetails(
	provider: MetadataProviderType,
	type: MediaType,
	externalId: string,
	signal?: AbortSignal
) {
	const { data, error } = await client.GET('/media/metadata/{provider}/{type}/{externalId}', {
		params: { path: { provider, type, externalId } },
		signal
	});
	if (error) throw new Error(error.message);
	if (!data) throw new Error('Media details were not returned');
	return data;
}

export async function getPersonDetails(
	provider: MetadataProviderType,
	personId: string,
	signal?: AbortSignal
) {
	const { data, error } = await client.GET('/people/{provider}/{personId}', {
		params: { path: { provider, personId } },
		signal
	});
	if (error) throw new Error(error.message);
	if (!data) throw new Error('Person details were not returned');
	return data;
}

export async function getMediaCollection(
	provider: MetadataProviderType,
	collectionId: string,
	signal?: AbortSignal
) {
	const { data, error } = await client.GET('/media/collections/{provider}/{collectionId}', {
		params: { path: { provider, collectionId } },
		signal
	});
	if (error) throw new Error(error.message);
	if (!data) throw new Error('Media collection was not returned');
	return data;
}
