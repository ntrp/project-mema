import { beforeEach, describe, expect, it, vi } from 'vitest';

const client = vi.hoisted(() => ({ GET: vi.fn(), POST: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client }));

import {
	advancedSearchMedia,
	autocompleteMedia,
	getMediaCollection,
	getMediaMetadataDetails,
	getPersonDetails
} from './searchMetadata';

describe('search metadata domain', () => {
	beforeEach(() => vi.clearAllMocks());

	it('maps search scopes and result groups', async () => {
		client.GET.mockResolvedValue({ data: { groups: ['result'] } });
		await expect(autocompleteMedia('matrix')).resolves.toEqual(['result']);
		await autocompleteMedia('matrix', 'library');
		await autocompleteMedia('matrix', 'providers');
		expect(client.GET).toHaveBeenNthCalledWith(
			2,
			'/media/autocomplete',
			expect.objectContaining({
				params: { query: expect.objectContaining({ includeProviders: false }) }
			})
		);
		expect(client.GET).toHaveBeenNthCalledWith(
			3,
			'/media/autocomplete',
			expect.objectContaining({
				params: { query: expect.objectContaining({ includeLibrary: false }) }
			})
		);
		client.POST.mockResolvedValueOnce({ data: { groups: ['advanced'] } });
		await expect(advancedSearchMedia({ query: 'matrix' })).resolves.toEqual(['advanced']);
	});

	it('loads required metadata detail resources', async () => {
		client.GET.mockResolvedValueOnce({ data: { title: 'Movie' } });
		await expect(getMediaMetadataDetails('tmdb', 'movie', '1')).resolves.toEqual({
			title: 'Movie'
		});
		client.GET.mockResolvedValueOnce({ data: { name: 'Person' } });
		await expect(getPersonDetails('tmdb', '2')).resolves.toEqual({ name: 'Person' });
		client.GET.mockResolvedValueOnce({ data: { name: 'Collection' } });
		await expect(getMediaCollection('tmdb', '3')).resolves.toEqual({ name: 'Collection' });
	});

	it('uses empty search fallbacks and rejects missing details', async () => {
		client.GET.mockResolvedValue({});
		client.POST.mockResolvedValue({});
		await expect(autocompleteMedia('none')).resolves.toEqual([]);
		await expect(advancedSearchMedia({})).resolves.toEqual([]);
		await expect(getMediaMetadataDetails('tmdb', 'movie', '1')).rejects.toThrow('not returned');
		await expect(getPersonDetails('tmdb', '1')).rejects.toThrow('not returned');
		await expect(getMediaCollection('tmdb', '1')).rejects.toThrow('not returned');
	});

	it('surfaces API errors from every operation', async () => {
		client.GET.mockResolvedValue({ error: { message: 'get failed' } });
		for (const read of [
			() => autocompleteMedia('query'),
			() => getMediaMetadataDetails('tmdb', 'movie', '1'),
			() => getPersonDetails('tmdb', '1'),
			() => getMediaCollection('tmdb', '1')
		])
			await expect(read()).rejects.toThrow('get failed');
		client.POST.mockResolvedValue({ error: { message: 'post failed' } });
		await expect(advancedSearchMedia({})).rejects.toThrow('post failed');
	});
});
