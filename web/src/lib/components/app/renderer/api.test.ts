import { beforeEach, describe, expect, it, vi } from 'vitest';

const client = vi.hoisted(() => ({ GET: vi.fn(), POST: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client }));

import * as api from './api';

describe('renderer API', () => {
	beforeEach(() => vi.clearAllMocks());

	it('lists devices and maps every remote command', async () => {
		client.GET.mockResolvedValueOnce({ data: { renderers: [{ id: 'uuid:tv' }] } });
		await expect(api.listRenderers()).resolves.toEqual([{ id: 'uuid:tv' }]);
		client.POST.mockResolvedValue({});
		client.GET.mockResolvedValueOnce({ data: { transportState: 'PLAYING' } });
		await expect(api.rendererStatus('uuid:tv')).resolves.toMatchObject({
			transportState: 'PLAYING'
		});
		await api.playOnRenderer('uuid:tv', 'media-1', '/movie.mkv');
		await api.rendererTransport('uuid:tv', 'pause');
		await api.seekRenderer('uuid:tv', 42);
		await api.setRendererVolume('uuid:tv', 35);
		await api.setRendererMute('uuid:tv', true);
		expect(client.POST).toHaveBeenCalledTimes(5);
		expect(client.POST).toHaveBeenCalledWith('/renderers/{id}/play', {
			params: { path: { id: 'uuid:tv' } },
			body: { mediaItemId: 'media-1', path: '/movie.mkv' }
		});
	});

	it('surfaces actionable API errors', async () => {
		client.GET.mockResolvedValue({ error: { message: 'Renderer is offline' } });
		await expect(api.listRenderers()).rejects.toThrow('Renderer is offline');
		client.POST.mockResolvedValue({ error: { message: 'Action unsupported' } });
		await expect(api.rendererTransport('uuid:tv', 'stop')).rejects.toThrow('Action unsupported');
	});
});
