import { beforeEach, describe, expect, it, vi } from 'vitest';
const client = vi.hoisted(() => ({ GET: vi.fn(), POST: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client }));
import * as api from './api';

describe('release and subtitle commands API', () => {
	beforeEach(() => {
		client.GET.mockReset();
		client.POST.mockReset();
	});

	it('maps release searches and every enqueue command', async () => {
		client.GET.mockResolvedValue({ data: {} });
		client.POST.mockResolvedValue({ data: { jobId: 1, message: 'queued' } });
		await expect(api.searchMediaReleases('media-1')).resolves.toEqual({ releases: [], errors: [] });
		await api.enqueueMediaReleaseSearch('media-1', 'query');
		await api.enqueueMediaAutomaticSearch('media-1');
		await api.enqueueMediaFulfillmentAction('media-1', {} as never);
		await api.enqueueMediaSubtitleSearch('media-1');
		await api.searchMediaSubtitles('media-1', {} as never);
		await api.grabMediaSubtitle('media-1', {} as never);
		await api.grabMediaRelease('media-1', { id: 'release-1' } as never, true, {} as never);
		expect(client.POST).toHaveBeenCalledTimes(7);
	});

	it('surfaces failures and missing command responses', async () => {
		client.GET.mockResolvedValueOnce({ error: { message: 'search failed' } });
		await expect(api.searchMediaReleases('media-1')).rejects.toThrow('search failed');
		client.POST.mockResolvedValueOnce({ data: undefined });
		await expect(api.enqueueMediaAutomaticSearch('media-1')).rejects.toThrow(
			'Automatic search job was not returned'
		);
	});
});
