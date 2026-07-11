import { beforeEach, describe, expect, it, vi } from 'vitest';

const client = vi.hoisted(() => ({ GET: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client }));

import * as api from './api';

describe('settings catalog API', () => {
	beforeEach(() => client.GET.mockReset());

	it('maps every catalog collection and its empty fallback', async () => {
		client.GET.mockResolvedValue({ data: {} });
		const calls = [
			api.listLanguages(),
			api.listTags(),
			api.listUsers(),
			api.listDownloadClients(),
			api.listIndexers(),
			api.listMetadataProviders(),
			api.listSubtitleProviders(),
			api.listLibraryFolders(),
			api.listPathMappings(),
			api.listMediaProfiles(),
			api.listCustomFormats()
		];
		await expect(Promise.all(calls)).resolves.toEqual(Array.from({ length: 11 }, () => []));
		expect(client.GET).toHaveBeenCalledTimes(11);
	});

	it('surfaces catalog request errors', async () => {
		client.GET.mockResolvedValue({ error: { message: 'catalog failed' } });
		await expect(api.listLanguages()).rejects.toThrow('catalog failed');
		await expect(api.listMetadataProviders()).rejects.toThrow('catalog failed');
	});
});
