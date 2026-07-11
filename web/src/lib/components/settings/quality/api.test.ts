import { beforeEach, describe, expect, it, vi } from 'vitest';

const clientMock = vi.hoisted(() => ({ GET: vi.fn(), PUT: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client: clientMock }));

import { listQualitySizeSettings, updateQualitySizeSettings } from './api';

describe('quality settings API', () => {
	beforeEach(() => {
		clientMock.GET.mockReset();
		clientMock.PUT.mockReset();
	});

	it('loads and updates quality sizes', async () => {
		clientMock.GET.mockResolvedValueOnce({ data: { qualities: [] } });
		await expect(listQualitySizeSettings()).resolves.toEqual({ qualities: [] });
		clientMock.PUT.mockResolvedValueOnce({ data: { qualities: [] } });
		await expect(updateQualitySizeSettings([])).resolves.toEqual({ qualities: [] });
		expect(clientMock.PUT).toHaveBeenCalledWith('/settings/quality-sizes', {
			body: { qualities: [] }
		});
	});

	it('rejects errors and empty responses', async () => {
		clientMock.GET.mockResolvedValueOnce({ error: { message: 'quality failed' } });
		await expect(listQualitySizeSettings()).rejects.toThrow('quality failed');
		clientMock.PUT.mockResolvedValueOnce({ data: undefined });
		await expect(updateQualitySizeSettings([])).rejects.toThrow(
			'Quality size settings were not returned'
		);
	});
});
