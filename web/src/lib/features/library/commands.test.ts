import { beforeEach, describe, expect, it, vi } from 'vitest';
const client = vi.hoisted(() => ({ GET: vi.fn(), POST: vi.fn(), PUT: vi.fn(), DELETE: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client }));
import * as api from './commands';

describe('library commands API', () => {
	beforeEach(() => Object.values(client).forEach((mock) => mock.mockReset()));

	it('maps item and request CRUD commands', async () => {
		client.GET.mockResolvedValue({ data: {} });
		client.POST.mockResolvedValue({ data: { id: 'created' } });
		client.PUT.mockResolvedValue({ data: { id: 'updated' } });
		client.DELETE.mockResolvedValue({});
		await expect(api.listMediaItems()).resolves.toEqual([]);
		await expect(api.createMediaItem({} as never)).resolves.toEqual({ id: 'created' });
		await expect(api.updateMediaItem('media-1', {} as never)).resolves.toEqual({ id: 'updated' });
		await expect(api.refreshMediaItemMetadata('media-1')).resolves.toEqual({ id: 'created' });
		await expect(api.listMediaRequests()).resolves.toEqual([]);
		await expect(api.createMediaRequest({} as never)).resolves.toEqual({ id: 'created' });
		await expect(api.getMediaRequest('request-1')).resolves.toEqual({});
		await expect(api.approveMediaRequest('request-1', {} as never)).resolves.toEqual({
			id: 'created'
		});
		await api.deleteMediaItem('media-1', { keepFiles: true });
		await expect(api.rescanMediaItemFiles('media-1')).resolves.toEqual({ id: 'created' });
	});

	it('rejects request errors and missing required responses', async () => {
		client.POST.mockResolvedValueOnce({ error: { message: 'create failed' } });
		await expect(api.createMediaItem({} as never)).rejects.toThrow('create failed');
		client.PUT.mockResolvedValueOnce({ data: undefined });
		await expect(api.updateMediaItem('media-1', {} as never)).rejects.toThrow(
			'Media item was not returned'
		);
	});
});
