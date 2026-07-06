import { client } from '$lib/api/client';

import type { MediaComponentSourceListResponse, MediaComponentSourceRetainRequest } from '../types';

export async function listMediaComponentSources(
	id: string
): Promise<MediaComponentSourceListResponse> {
	const { data, error } = await client.GET('/media/items/{id}/component-sources', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	return data ?? { sources: [] };
}

export async function retainMediaComponentSource(
	id: string,
	request: MediaComponentSourceRetainRequest
) {
	const { data, error } = await client.POST('/media/items/{id}/component-sources', {
		params: { path: { id } },
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Component source was not returned');
	}
	return data;
}

export async function getMediaComponentSource(id: string, sourceId: string) {
	const { data, error } = await client.GET('/media/items/{id}/component-sources/{sourceId}', {
		params: { path: { id, sourceId } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Component source was not returned');
	}
	return data;
}

export async function releaseMediaComponentSource(id: string, sourceId: string) {
	const { data, error } = await client.POST(
		'/media/items/{id}/component-sources/{sourceId}/release',
		{
			params: { path: { id, sourceId } }
		}
	);

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Component source was not returned');
	}
	return data;
}
