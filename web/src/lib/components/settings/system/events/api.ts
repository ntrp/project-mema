import { client } from '$lib/api/client';
import type { SystemEventSettings, SystemEventSettingsRequest } from '$lib/settings/types';

export async function listSystemEvents(options: { before?: string; limit?: number } = {}) {
	const { data, error } = await client.GET('/system/events', { params: { query: options } });
	if (error) throw new Error(error.message);
	return data ?? { events: [], hasMore: false };
}

export async function deleteSystemEvent(id: string) {
	const { error } = await client.DELETE('/system/events/{id}', {
		params: { path: { id } }
	});
	if (error) throw new Error(error.message);
}

export async function clearSystemEvents() {
	const { error } = await client.DELETE('/system/events');
	if (error) throw new Error(error.message);
}

export async function getSystemEventSettings(): Promise<SystemEventSettings> {
	const { data, error } = await client.GET('/system/event-settings');
	if (error) throw new Error(error.message);
	if (!data) throw new Error('Event settings request did not return a result');
	return data;
}

export async function updateSystemEventSettings(
	request: SystemEventSettingsRequest
): Promise<SystemEventSettings> {
	const { data, error } = await client.PUT('/system/event-settings', { body: request });
	if (error) throw new Error(error.message);
	if (!data) throw new Error('Event settings update did not return a result');
	return data;
}
