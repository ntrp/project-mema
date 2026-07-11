import { client } from '$lib/api/client';
import type {
	FileDeleteSettings,
	FileDeleteSettingsRequest,
	FileNamingSettings,
	FileNamingSettingsRequest
} from '$lib/settings/types';

function required<T>(data: T | undefined, message: string): T {
	if (!data) throw new Error(message);
	return data;
}

export async function getFileNamingSettings(): Promise<FileNamingSettings> {
	const { data, error } = await client.GET('/settings/file-naming');
	if (error) throw new Error(error.message);
	return required(data, 'File naming settings were not returned');
}

export async function updateFileNamingSettings(request: FileNamingSettingsRequest) {
	const { data, error } = await client.PUT('/settings/file-naming', { body: request });
	if (error) throw new Error(error.message);
	return required(data, 'File naming settings were not returned');
}

export async function getFileDeleteSettings(): Promise<FileDeleteSettings> {
	const { data, error } = await client.GET('/settings/file-delete');
	if (error) throw new Error(error.message);
	return required(data, 'File delete settings were not returned');
}

export async function updateFileDeleteSettings(request: FileDeleteSettingsRequest) {
	const { data, error } = await client.PUT('/settings/file-delete', { body: request });
	if (error) throw new Error(error.message);
	return required(data, 'File delete settings were not returned');
}
