import { client } from '$lib/api/client';
import type {
	GrabSubtitleRequest,
	ManualSubtitleSearchRequest,
	MediaFulfillmentActionRequest,
	ReleaseCandidate,
	ReleaseOverrideDetails,
	SubtitleSearchRequest
} from '$lib/settings/types';
export type { ReleaseCandidate, ReleaseOverrideDetails } from '$lib/settings/types';

export async function searchMediaReleases(id: string) {
	const { data, error } = await client.GET('/media/items/{id}/releases', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	return {
		releases: data?.releases ?? [],
		errors: data?.errors ?? []
	};
}

export async function enqueueMediaReleaseSearch(id: string, query?: string) {
	const { data, error } = await client.POST('/media/items/{id}/release-searches', {
		params: { path: { id } },
		body: query ? { query } : undefined
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Release search job was not returned');
	}
	return data;
}

export async function enqueueMediaAutomaticSearch(id: string) {
	const { data, error } = await client.POST('/media/items/{id}/automatic-searches', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Automatic search job was not returned');
	}
	return data;
}

export async function enqueueMediaFulfillmentAction(
	id: string,
	request: MediaFulfillmentActionRequest
) {
	const { data, error } = await client.POST('/media/items/{id}/fulfillment-actions', {
		params: { path: { id } },
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Fulfillment job was not returned');
	}
	return data;
}

export async function enqueueMediaSubtitleSearch(id: string, request: SubtitleSearchRequest = {}) {
	const { data, error } = await client.POST('/media/items/{id}/subtitle-searches', {
		params: { path: { id } },
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Subtitle search job was not returned');
	}
	return data;
}

export async function searchMediaSubtitles(id: string, request: ManualSubtitleSearchRequest) {
	const { data, error } = await client.POST('/media/items/{id}/subtitle-search-results', {
		params: { path: { id } },
		body: request
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Subtitle search results were not returned');
	}
	return data;
}

export async function grabMediaSubtitle(id: string, request: GrabSubtitleRequest) {
	const { data, error } = await client.POST('/media/items/{id}/subtitle-grabs', {
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

export async function grabMediaRelease(
	id: string,
	release: ReleaseCandidate,
	overrideMatch = false,
	overrideDetails?: ReleaseOverrideDetails
) {
	const { data, error } = await client.POST('/media/items/{id}/grab', {
		params: { path: { id } },
		body: {
			releaseId: release.id,
			overrideMatch,
			...(overrideDetails ? { overrideDetails } : {})
		}
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Download grab did not return a result');
	}
	return data;
}
