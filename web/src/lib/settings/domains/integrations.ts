import { client } from '$lib/api/client';
import { normalizeDownloadClientForm, normalizeIndexerForm } from '../forms';
import type { DownloadClientForm, IndexerBulkUpdateRequest, IndexerForm } from '../types';

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

export async function testDownloadClientConfig(form: DownloadClientForm) {
	const { data, error } = await client.POST('/settings/download-clients/test', {
		body: normalizeDownloadClientForm(form)
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

export async function listIndexerCatalog() {
	const { data, error } = await client.GET('/settings/indexer-catalog');
	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Indexer catalog did not return a result');
	}
	return data;
}

export async function listIndexerAppProfiles() {
	const { data, error } = await client.GET('/settings/indexer-app-profiles');
	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Indexer app profiles did not return a result');
	}
	return data.profiles;
}

export async function listIndexerProxies() {
	const { data, error } = await client.GET('/settings/indexer-proxies');
	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Indexer proxies did not return a result');
	}
	return data.proxies;
}

export async function bulkUpdateIndexers(body: IndexerBulkUpdateRequest) {
	const { data, error } = await client.PUT('/settings/indexers/bulk', { body });
	if (error) {
		throw new Error(error.message);
	}
	return data?.indexers ?? [];
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

export async function testIndexerConfig(form: IndexerForm) {
	const { data, error } = await client.POST('/settings/indexers/test', {
		body: normalizeIndexerForm(form)
	});

	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Indexer test did not return a result');
	}
	return data;
}
