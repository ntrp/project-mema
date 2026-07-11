import { client } from '$lib/api/client';

export async function listLanguages() {
	const { data, error } = await client.GET('/settings/languages');
	if (error) throw new Error(error.message);
	return data?.languages ?? [];
}

export async function listTags() {
	const { data, error } = await client.GET('/settings/tags');
	if (error) throw new Error(error.message);
	return data?.tags ?? [];
}

export async function listUsers() {
	const { data, error } = await client.GET('/settings/users');
	if (error) throw new Error(error.message);
	return data?.users ?? [];
}

export async function listDownloadClients() {
	const { data, error } = await client.GET('/settings/download-clients');
	if (error) throw new Error(error.message);
	return data?.clients ?? [];
}

export async function listIndexers() {
	const { data, error } = await client.GET('/settings/indexers');
	if (error) throw new Error(error.message);
	return data?.indexers ?? [];
}

export async function listMetadataProviders() {
	const { data, error } = await client.GET('/settings/metadata-providers');
	if (error) throw new Error(error.message);
	return data?.providers ?? [];
}

export async function listSubtitleProviders() {
	const { data, error } = await client.GET('/settings/subtitle-providers');
	if (error) throw new Error(error.message);
	return data?.providers ?? [];
}

export async function listLibraryFolders() {
	const { data, error } = await client.GET('/settings/library/folders');
	if (error) throw new Error(error.message);
	return data?.folders ?? [];
}

export async function listPathMappings() {
	const { data, error } = await client.GET('/settings/library/path-mappings');
	if (error) throw new Error(error.message);
	return data?.mappings ?? [];
}

export async function listMediaProfiles() {
	const { data, error } = await client.GET('/settings/profiles');
	if (error) throw new Error(error.message);
	return data?.profiles ?? [];
}

export async function listCustomFormats() {
	const { data, error } = await client.GET('/settings/custom-formats');
	if (error) throw new Error(error.message);
	return data?.formats ?? [];
}
