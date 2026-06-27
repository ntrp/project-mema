import { client } from '$lib/api/client';

import { normalizeDownloadClientForm, normalizeIndexerForm } from './forms';
import type { DownloadClientForm, IndexerForm, SettingsData } from './types';

export async function currentSessionAuthenticated() {
	const { data } = await client.GET('/auth/session');
	return Boolean(data?.authenticated);
}

export async function login(username: string, password: string) {
	const { data, error } = await client.POST('/auth/login', {
		body: { username, password }
	});

	if (error || !data?.authenticated) {
		throw new Error(error?.message ?? 'Login failed');
	}
}

export async function logout() {
	const { error } = await client.POST('/auth/logout');

	if (error) {
		throw new Error(error.message);
	}
}

export async function loadSettings(): Promise<SettingsData> {
	const [clientResult, indexerResult] = await Promise.all([
		client.GET('/settings/download-clients'),
		client.GET('/settings/indexers')
	]);

	if (clientResult.error) {
		throw new Error(clientResult.error.message);
	}
	if (indexerResult.error) {
		throw new Error(indexerResult.error.message);
	}

	return {
		downloadClients: clientResult.data?.clients ?? [],
		indexers: indexerResult.data?.indexers ?? []
	};
}

export async function saveDownloadClient(form: DownloadClientForm) {
	const body = normalizeDownloadClientForm(form);
	const result = form.id
		? await client.PUT('/settings/download-clients/{id}', {
				params: { path: { id: form.id } },
				body
			})
		: await client.POST('/settings/download-clients', { body });

	if (result.error) {
		throw new Error(result.error.message);
	}
}

export async function testDownloadClient(id: string) {
	const { data, error } = await client.POST('/settings/download-clients/{id}/test', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Download client test did not return a result');
	}
	return data;
}

export async function saveIndexer(form: IndexerForm) {
	const body = normalizeIndexerForm(form);
	const result = form.id
		? await client.PUT('/settings/indexers/{id}', {
				params: { path: { id: form.id } },
				body
			})
		: await client.POST('/settings/indexers', { body });

	if (result.error) {
		throw new Error(result.error.message);
	}
}

export async function testIndexer(id: string) {
	const { data, error } = await client.POST('/settings/indexers/{id}/test', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Indexer test did not return a result');
	}
	return data;
}

export async function deleteDownloadClient(id: string) {
	const { error } = await client.DELETE('/settings/download-clients/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function deleteIndexer(id: string) {
	const { error } = await client.DELETE('/settings/indexers/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}
