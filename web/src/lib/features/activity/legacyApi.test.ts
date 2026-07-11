import { beforeEach, describe, expect, it, vi } from 'vitest';
const client = vi.hoisted(() => ({ GET: vi.fn(), POST: vi.fn(), DELETE: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client }));
import * as api from './legacyApi';

describe('legacy activity API compatibility', () => {
	beforeEach(() => Object.values(client).forEach((mock) => mock.mockReset()));

	it('maps activity and blocklist commands', async () => {
		client.GET.mockResolvedValue({ data: {} });
		client.POST.mockResolvedValue({ data: { id: 'activity-1' } });
		client.DELETE.mockResolvedValue({});
		await expect(api.listDownloadActivity()).resolves.toEqual([]);
		await expect(api.listReleaseBlocklist()).resolves.toEqual([]);
		await api.deleteReleaseBlocklistItem('block-1');
		await api.clearReleaseBlocklist();
		await expect(api.cancelDownloadActivity('activity-1')).resolves.toEqual({ id: 'activity-1' });
		await api.deleteDownloadActivity('activity-1');
		await expect(api.manualImportDownloadActivity('activity-1', {} as never)).resolves.toEqual({
			id: 'activity-1'
		});
	});

	it('surfaces failures and missing activity results', async () => {
		client.GET.mockResolvedValueOnce({ error: { message: 'activity failed' } });
		await expect(api.listDownloadActivity()).rejects.toThrow('activity failed');
		client.POST.mockResolvedValueOnce({ data: undefined });
		await expect(api.cancelDownloadActivity('activity-1')).rejects.toThrow(
			'Download activity was not returned'
		);
	});
});
