import { client } from '$lib/api/client';

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

export async function deleteMetadataProvider(id: string) {
	const { error } = await client.DELETE('/settings/metadata-providers/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function deleteSubtitleProvider(id: string) {
	const { error } = await client.DELETE('/settings/subtitle-providers/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}

export async function deleteMediaProfile(id: string) {
	const { error } = await client.DELETE('/settings/profiles/{id}', {
		params: { path: { id } }
	});

	if (error) {
		throw new Error(error.message);
	}
}
