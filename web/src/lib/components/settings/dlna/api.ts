import { client } from '$lib/api/client';
import type { DLNASettings, DLNASettingsRequest } from '$lib/settings/types';

function result(data: DLNASettings | undefined, message: string) {
	if (!data) throw new Error(message);
	return data;
}

export async function getDLNASettings() {
	const { data, error } = await client.GET('/settings/dlna');
	if (error) throw new Error(error.message);
	return result(data, 'DLNA settings request did not return a result');
}

export async function updateDLNASettings(request: DLNASettingsRequest) {
	const { data, error } = await client.PUT('/settings/dlna', { body: request });
	if (error) throw new Error(error.message);
	return result(data, 'DLNA settings update did not return a result');
}

export async function restartDLNA() {
	const { data, error } = await client.POST('/settings/dlna/restart');
	if (error) throw new Error(error.message);
	return result(data, 'DLNA restart did not return a result');
}
