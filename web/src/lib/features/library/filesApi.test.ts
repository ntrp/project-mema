import { beforeEach, describe, expect, it, vi } from 'vitest';

const client = vi.hoisted(() => ({ GET: vi.fn(), POST: vi.fn(), PUT: vi.fn(), DELETE: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client }));
import * as api from './filesApi';

describe('media files API', () => {
	beforeEach(() => Object.values(client).forEach((mock) => mock.mockReset()));

	it('maps file, subtitle, history, and rename operations', async () => {
		client.GET.mockResolvedValue({ data: {} });
		client.POST.mockResolvedValue({ data: { id: 'media-1' } });
		client.PUT.mockResolvedValue({ data: { id: 'media-1' } });
		client.DELETE.mockResolvedValue({ data: { id: 'media-1' } });
		await expect(api.deleteMediaItemFile('media-1', '/movie.mkv')).resolves.toEqual({
			id: 'media-1'
		});
		await expect(api.deleteMediaItemFileTrack('media-1', {} as never)).resolves.toEqual({
			id: 'media-1'
		});
		await expect(api.listMediaItemSubtitles('media-1')).resolves.toEqual({});
		await expect(api.deleteMediaItemSubtitle('media-1', 'sub-1')).resolves.toEqual({
			id: 'media-1'
		});
		await expect(api.updateMediaItemSubtitle('media-1', 'sub-1', {} as never)).resolves.toEqual({
			id: 'media-1'
		});
		await expect(api.listMediaFileHistory('media-1')).resolves.toEqual({});
		await expect(api.previewMediaRename('media-1')).resolves.toEqual({});
		client.POST.mockResolvedValueOnce({ data: undefined });
		await expect(api.applyMediaRename('media-1')).resolves.toMatchObject({ appliedCount: 0 });
	});

	it('rejects server errors and required empty mutation responses', async () => {
		client.POST.mockResolvedValueOnce({ error: { message: 'file failed' } });
		await expect(api.deleteMediaItemFile('media-1', '/movie.mkv')).rejects.toThrow('file failed');
		client.DELETE.mockResolvedValueOnce({ data: undefined });
		await expect(api.deleteMediaItemSubtitle('media-1', 'sub-1')).rejects.toThrow(
			'Media item was not returned'
		);
	});
});
