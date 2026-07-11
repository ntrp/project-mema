import { client } from '$lib/api/client';
import type { ManualImportRequest, ReleaseBlocklistItem } from '$lib/settings/types';

export async function listDownloadActivity() {
	const { data, error } = await client.GET('/activity/downloads');

	if (error) {
		throw new Error(error.message);
	}
	return data?.activities ?? [];
}

export async function listReleaseBlocklist(): Promise<ReleaseBlocklistItem[]> {
	const { data, error } = await client.GET('/activity/blocklist');

	if (error) {
		throw new Error(error.message);
	}
	return data?.items ?? [];
}

export async function deleteReleaseBlocklistItem(id: string) {
	const { error } = await client.DELETE('/activity/blocklist/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function clearReleaseBlocklist() {
	const { error } = await client.DELETE('/activity/blocklist');

	if (error) {
		throw new Error(error.message);
	}
}

export async function cancelDownloadActivity(id: string) {
	const { data, error } = await client.POST('/activity/downloads/{id}/cancel', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Download activity was not returned');
	}
	return data;
}

export async function deleteDownloadActivity(id: string) {
	const { error } = await client.DELETE('/activity/downloads/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function manualImportDownloadActivity(id: string, body: ManualImportRequest) {
	const { data, error } = await client.POST('/activity/downloads/{id}/manual-import', {
		params: { path: { id } },
		body
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Download activity was not returned');
	}
	return data;
}
