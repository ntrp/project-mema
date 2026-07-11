import { client } from '$lib/api/client';
import type { QualitySizeSettingRequest, QualitySizeSettingsResponse } from '$lib/settings/types';

export async function listQualitySizeSettings(): Promise<QualitySizeSettingsResponse> {
	const { data, error } = await client.GET('/settings/quality-sizes');
	if (error) throw new Error(error.message);
	if (!data) throw new Error('Quality size settings were not returned');
	return data;
}

export async function updateQualitySizeSettings(qualities: QualitySizeSettingRequest[]) {
	const { data, error } = await client.PUT('/settings/quality-sizes', {
		body: { qualities }
	});
	if (error) throw new Error(error.message);
	if (!data) throw new Error('Quality size settings were not returned');
	return data;
}
