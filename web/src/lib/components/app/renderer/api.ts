import { client } from '$lib/api/client';
import type { Renderer, RendererStatus, RendererTransportRequest } from './types';

export async function listRenderers(): Promise<Renderer[]> {
	const { data, error } = await client.GET('/renderers');
	if (error) throw new Error(error.message);
	return data?.renderers ?? [];
}

export async function playOnRenderer(id: string, mediaItemId: string, path: string): Promise<void> {
	const { error } = await client.POST('/renderers/{id}/play', {
		params: { path: { id } },
		body: { mediaItemId, path }
	});
	if (error) throw new Error(error.message);
}

export async function rendererStatus(id: string): Promise<RendererStatus> {
	const { data, error } = await client.GET('/renderers/{id}/status', {
		params: { path: { id } }
	});
	if (error) throw new Error(error.message);
	if (!data) throw new Error('Renderer did not return its status');
	return data;
}

export async function rendererTransport(
	id: string,
	action: RendererTransportRequest['action']
): Promise<void> {
	const { error } = await client.POST('/renderers/{id}/transport', {
		params: { path: { id } },
		body: { action }
	});
	if (error) throw new Error(error.message);
}

export async function seekRenderer(id: string, positionSeconds: number): Promise<void> {
	const { error } = await client.POST('/renderers/{id}/seek', {
		params: { path: { id } },
		body: { positionSeconds }
	});
	if (error) throw new Error(error.message);
}

export async function setRendererVolume(id: string, volume: number): Promise<void> {
	const { error } = await client.POST('/renderers/{id}/volume', {
		params: { path: { id } },
		body: { volume }
	});
	if (error) throw new Error(error.message);
}

export async function setRendererMute(id: string, muted: boolean): Promise<void> {
	const { error } = await client.POST('/renderers/{id}/mute', {
		params: { path: { id } },
		body: { muted }
	});
	if (error) throw new Error(error.message);
}
