import { client } from '$lib/api/client';
import type {
	MediaItemCreateRequest,
	MediaItemUpdateRequest,
	MediaRequestApproveRequest,
	MediaRequestCreateRequest
} from '$lib/settings/types';

export async function listMediaItems() {
	const { data, error } = await client.GET('/media/items');

	if (error) {
		throw new Error(error.message);
	}
	return data?.items ?? [];
}

export async function createMediaItem(request: MediaItemCreateRequest) {
	const { data, error } = await client.POST('/media/items', { body: request });

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media item was not returned');
	}
	return data;
}

export async function updateMediaItem(id: string, request: MediaItemUpdateRequest) {
	const { data, error } = await client.PUT('/media/items/{id}', {
		params: { path: { id } },
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media item was not returned');
	}
	return data;
}

export async function refreshMediaItemMetadata(id: string) {
	const { data, error } = await client.POST('/media/items/{id}/metadata/refresh', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media item was not returned');
	}
	return data;
}

export async function listMediaRequests() {
	const { data, error } = await client.GET('/media/requests');

	if (error) {
		throw new Error(error.message);
	}
	return data?.requests ?? [];
}

export async function createMediaRequest(request: MediaRequestCreateRequest) {
	const { data, error } = await client.POST('/media/requests', { body: request });

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media request was not returned');
	}
	return data;
}

export async function getMediaRequest(id: string) {
	const { data, error } = await client.GET('/media/requests/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media request was not returned');
	}
	return data;
}

export async function approveMediaRequest(id: string, request: MediaRequestApproveRequest) {
	const { data, error } = await client.POST('/media/requests/{id}/approve', {
		params: { path: { id } },
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media request approval was not returned');
	}
	return data;
}

export async function deleteMediaItem(id: string, options: { keepFiles?: boolean } = {}) {
	const { error } = await client.DELETE('/media/items/{id}', {
		params: { path: { id }, query: { keepFiles: options.keepFiles } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function rescanMediaItemFiles(id: string) {
	const { data, error } = await client.POST('/media/items/{id}/files/rescan', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Media item was not returned');
	}
	return data;
}
