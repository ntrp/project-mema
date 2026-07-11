import { beforeEach, describe, expect, it, vi } from 'vitest';

const client = vi.hoisted(() => ({ GET: vi.fn(), POST: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client }));

import { listCustomFormats, listMediaProfiles, testCustomFormatParsing } from './catalogRead';
import { getIndexerSearch, getMetadataCache } from './cacheInspection';

describe('settings read-domain coverage', () => {
	beforeEach(() => vi.clearAllMocks());

	it('maps catalog reads and cache fallbacks', async () => {
		client.GET.mockResolvedValueOnce({ data: { profiles: ['profile'] } })
			.mockResolvedValueOnce({ data: { formats: ['format'] } })
			.mockResolvedValueOnce({})
			.mockResolvedValueOnce({});
		await expect(listMediaProfiles()).resolves.toEqual(['profile']);
		await expect(listCustomFormats()).resolves.toEqual(['format']);
		await expect(getMetadataCache()).resolves.toMatchObject({ entries: [] });
		await expect(getIndexerSearch()).resolves.toMatchObject({ cacheEntries: [] });
		client.POST.mockResolvedValueOnce({ data: { matches: [] } });
		await expect(testCustomFormatParsing('movie.mkv')).resolves.toEqual({ matches: [] });
	});

	it('surfaces catalog errors and missing parsing results', async () => {
		client.GET.mockResolvedValue({ error: { message: 'read failed' } });
		await expect(listMediaProfiles()).rejects.toThrow('read failed');
		await expect(listCustomFormats()).rejects.toThrow('read failed');
		client.POST.mockResolvedValueOnce({});
		await expect(testCustomFormatParsing('movie.mkv')).rejects.toThrow('not returned');
	});
});
