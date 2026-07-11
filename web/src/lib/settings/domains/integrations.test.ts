import { beforeEach, describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
	client: { GET: vi.fn(), POST: vi.fn(), PUT: vi.fn() },
	downloadBody: { normalized: 'download' },
	indexerBody: { normalized: 'indexer' }
}));
vi.mock('$lib/api/client', () => ({ client: mocks.client }));
vi.mock('../forms', () => ({
	normalizeDownloadClientForm: () => mocks.downloadBody,
	normalizeIndexerForm: () => mocks.indexerBody
}));

import {
	bulkUpdateIndexers,
	listIndexerAppProfiles,
	listIndexerCatalog,
	listIndexerProxies,
	saveDownloadClient,
	saveIndexer,
	testDownloadClient,
	testDownloadClientConfig,
	testIndexer,
	testIndexerConfig
} from './integrations';

describe('integration settings domain', () => {
	beforeEach(() => vi.clearAllMocks());

	it('creates and updates normalized integrations', async () => {
		mocks.client.POST.mockResolvedValue({});
		mocks.client.PUT.mockResolvedValue({});
		await saveDownloadClient({} as never);
		await saveDownloadClient({ id: 'client-1' } as never);
		await saveIndexer({} as never);
		await saveIndexer({ id: 'indexer-1' } as never);
		expect(mocks.client.POST).toHaveBeenCalledWith('/settings/download-clients', {
			body: mocks.downloadBody
		});
		expect(mocks.client.PUT).toHaveBeenCalledWith(
			'/settings/indexers/{id}',
			expect.objectContaining({ body: mocks.indexerBody })
		);
	});

	it('returns test, catalog, profile, proxy, and bulk results', async () => {
		mocks.client.POST.mockResolvedValue({ data: { success: true } });
		await expect(testDownloadClient('client')).resolves.toMatchObject({ success: true });
		await expect(testDownloadClientConfig({} as never)).resolves.toMatchObject({ success: true });
		await expect(testIndexer('indexer')).resolves.toMatchObject({ success: true });
		await expect(testIndexerConfig({} as never)).resolves.toMatchObject({ success: true });

		mocks.client.GET.mockResolvedValueOnce({ data: { definitions: ['catalog'] } })
			.mockResolvedValueOnce({ data: { profiles: ['profile'] } })
			.mockResolvedValueOnce({ data: { proxies: ['proxy'] } });
		await expect(listIndexerCatalog()).resolves.toMatchObject({ definitions: ['catalog'] });
		await expect(listIndexerAppProfiles()).resolves.toEqual(['profile']);
		await expect(listIndexerProxies()).resolves.toEqual(['proxy']);
		mocks.client.PUT.mockResolvedValueOnce({ data: { indexers: ['updated'] } });
		await expect(bulkUpdateIndexers({} as never)).resolves.toEqual(['updated']);
		mocks.client.PUT.mockResolvedValueOnce({});
		await expect(bulkUpdateIndexers({} as never)).resolves.toEqual([]);
	});

	it('surfaces API failures', async () => {
		mocks.client.POST.mockResolvedValue({ error: { message: 'post failed' } });
		for (const command of [
			() => saveDownloadClient({} as never),
			() => saveIndexer({} as never),
			() => testDownloadClient('id'),
			() => testDownloadClientConfig({} as never),
			() => testIndexer('id'),
			() => testIndexerConfig({} as never)
		])
			await expect(command()).rejects.toThrow('post failed');

		mocks.client.GET.mockResolvedValue({ error: { message: 'get failed' } });
		for (const read of [listIndexerCatalog, listIndexerAppProfiles, listIndexerProxies])
			await expect(read()).rejects.toThrow('get failed');
		mocks.client.PUT.mockResolvedValue({ error: { message: 'put failed' } });
		await expect(saveDownloadClient({ id: 'id' } as never)).rejects.toThrow('put failed');
		await expect(saveIndexer({ id: 'id' } as never)).rejects.toThrow('put failed');
		await expect(bulkUpdateIndexers({} as never)).rejects.toThrow('put failed');
	});

	it('rejects missing required command data', async () => {
		mocks.client.POST.mockResolvedValue({});
		for (const command of [
			() => testDownloadClient('id'),
			() => testDownloadClientConfig({} as never),
			() => testIndexer('id'),
			() => testIndexerConfig({} as never)
		])
			await expect(command()).rejects.toThrow('did not return');
		mocks.client.GET.mockResolvedValue({});
		for (const read of [listIndexerCatalog, listIndexerAppProfiles, listIndexerProxies])
			await expect(read()).rejects.toThrow('did not return');
	});
});
