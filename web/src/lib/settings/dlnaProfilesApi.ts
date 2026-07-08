import { client } from '$lib/api/client';
import type {
	DLNAClientDiagnostic,
	DLNARendererDeviceOverride,
	DLNARendererDeviceOverrideRequest,
	DLNARendererProfile,
	DLNARendererProfileCloneRequest,
	DLNARendererProfileCreateRequest,
	DLNARendererProfileRequest
} from '$lib/settings/types';

export async function listDLNARendererProfiles(): Promise<DLNARendererProfile[]> {
	const { data, error } = await client.GET('/settings/dlna/profiles');
	if (error) throw new Error(error.message);
	return data?.profiles ?? [];
}

export async function createDLNARendererProfile(
	request: DLNARendererProfileCreateRequest
): Promise<DLNARendererProfile> {
	const { data, error } = await client.POST('/settings/dlna/profiles', { body: request });
	if (error) throw new Error(error.message);
	if (!data) throw new Error('DLNA renderer profile create did not return a result');
	return data;
}

export async function importDLNARendererProfile(
	request: DLNARendererProfileCreateRequest
): Promise<DLNARendererProfile> {
	const { data, error } = await client.POST('/settings/dlna/profiles/import', { body: request });
	if (error) throw new Error(error.message);
	if (!data) throw new Error('DLNA renderer profile import did not return a result');
	return data;
}

export async function updateDLNARendererProfile(
	id: string,
	request: DLNARendererProfileRequest
): Promise<DLNARendererProfile> {
	const { data, error } = await client.PUT('/settings/dlna/profiles/{id}', {
		params: { path: { id } },
		body: request
	});
	if (error) throw new Error(error.message);
	if (!data) throw new Error('DLNA renderer profile update did not return a result');
	return data;
}

export async function cloneDLNARendererProfile(
	id: string,
	request: DLNARendererProfileCloneRequest
): Promise<DLNARendererProfile> {
	const { data, error } = await client.POST('/settings/dlna/profiles/{id}/clone', {
		params: { path: { id } },
		body: request
	});
	if (error) throw new Error(error.message);
	if (!data) throw new Error('DLNA renderer profile clone did not return a result');
	return data;
}

export async function resetDLNARendererProfile(id: string): Promise<DLNARendererProfile> {
	const { data, error } = await client.POST('/settings/dlna/profiles/{id}/reset', {
		params: { path: { id } }
	});
	if (error) throw new Error(error.message);
	if (!data) throw new Error('DLNA renderer profile reset did not return a result');
	return data;
}

export async function exportDLNARendererProfile(id: string): Promise<DLNARendererProfile> {
	const { data, error } = await client.GET('/settings/dlna/profiles/{id}/export', {
		params: { path: { id } }
	});
	if (error) throw new Error(error.message);
	if (!data) throw new Error('DLNA renderer profile export did not return a result');
	return data;
}

export async function deleteDLNARendererProfile(id: string): Promise<void> {
	const { error } = await client.DELETE('/settings/dlna/profiles/{id}', {
		params: { path: { id } }
	});
	if (error) throw new Error(error.message);
}

export async function listDLNARendererDeviceOverrides(): Promise<DLNARendererDeviceOverride[]> {
	const { data, error } = await client.GET('/settings/dlna/device-overrides');
	if (error) throw new Error(error.message);
	return data?.overrides ?? [];
}

export async function upsertDLNARendererDeviceOverride(
	request: DLNARendererDeviceOverrideRequest
): Promise<DLNARendererDeviceOverride> {
	const { data, error } = await client.POST('/settings/dlna/device-overrides', { body: request });
	if (error) throw new Error(error.message);
	if (!data) throw new Error('DLNA renderer override save did not return a result');
	return data;
}

export async function deleteDLNARendererDeviceOverride(id: string): Promise<void> {
	const { error } = await client.DELETE('/settings/dlna/device-overrides/{id}', {
		params: { path: { id } }
	});
	if (error) throw new Error(error.message);
}

export async function listDLNARecentDevices(): Promise<DLNAClientDiagnostic[]> {
	const { data, error } = await client.GET('/settings/dlna/recent-devices');
	if (error) throw new Error(error.message);
	return data?.devices ?? [];
}
