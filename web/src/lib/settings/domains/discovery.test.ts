import { beforeEach, describe, expect, it, vi } from 'vitest';

const client = vi.hoisted(() => ({ GET: vi.fn(), POST: vi.fn(), DELETE: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client }));

import {
	addDiscoverBlacklistItem,
	autocompleteDiscoverMovieFacet,
	autocompleteDiscoverSeriesFacet,
	deleteDiscoverBlacklistItem,
	listDiscoverBlacklist,
	loadMediaDiscoverSection,
	loadMediaDiscoverSections,
	searchDiscoverMovies,
	searchDiscoverSeries,
	searchMedia
} from './discovery';

describe('discovery settings domain', () => {
	beforeEach(() => vi.clearAllMocks());

	it('maps successful discovery reads and commands', async () => {
		client.POST.mockResolvedValueOnce({ data: { results: ['result'] } });
		await expect(searchMedia({ query: 'matrix' } as never)).resolves.toEqual(['result']);
		client.GET.mockResolvedValueOnce({ data: { sections: ['popular'] } });
		await expect(loadMediaDiscoverSections()).resolves.toEqual(['popular']);
		client.GET.mockResolvedValueOnce({ data: { id: 'popular' } });
		await expect(loadMediaDiscoverSection('popular', 2, 10)).resolves.toEqual({ id: 'popular' });

		for (const search of [searchDiscoverMovies, searchDiscoverSeries]) {
			client.GET.mockResolvedValueOnce({ data: { results: ['title'], hasMore: true } });
			await expect(search({ page: 2 })).resolves.toMatchObject({ hasMore: true });
		}
		for (const autocomplete of [autocompleteDiscoverMovieFacet, autocompleteDiscoverSeriesFacet]) {
			client.GET.mockResolvedValueOnce({ data: { options: ['Drama'] } });
			await expect(autocomplete('genres', 'dra')).resolves.toEqual(['Drama']);
		}
		client.GET.mockResolvedValueOnce({ data: { items: ['blocked'] } });
		await expect(listDiscoverBlacklist()).resolves.toEqual(['blocked']);
		client.POST.mockResolvedValueOnce({ data: { id: 'blocked-1' } });
		await expect(addDiscoverBlacklistItem({ title: 'Hidden' } as never)).resolves.toEqual({
			id: 'blocked-1'
		});
		client.DELETE.mockResolvedValueOnce({});
		await expect(deleteDiscoverBlacklistItem('blocked-1')).resolves.toBeUndefined();
	});

	it('returns empty collection fallbacks', async () => {
		client.POST.mockResolvedValueOnce({});
		await expect(searchMedia({} as never)).resolves.toEqual([]);
		client.GET.mockResolvedValue({});
		await expect(loadMediaDiscoverSections()).resolves.toEqual([]);
		await expect(searchDiscoverMovies({})).resolves.toEqual({ results: [], hasMore: false });
		await expect(searchDiscoverSeries({})).resolves.toEqual({ results: [], hasMore: false });
		await expect(autocompleteDiscoverMovieFacet('genres', '')).resolves.toEqual([]);
		await expect(autocompleteDiscoverSeriesFacet('genres', '')).resolves.toEqual([]);
		await expect(listDiscoverBlacklist()).resolves.toEqual([]);
	});

	it('surfaces API errors and missing required responses', async () => {
		client.GET.mockResolvedValue({ error: { message: 'read failed' } });
		for (const read of [
			() => loadMediaDiscoverSections(),
			() => loadMediaDiscoverSection('popular'),
			() => searchDiscoverMovies({}),
			() => searchDiscoverSeries({}),
			() => autocompleteDiscoverMovieFacet('genres', ''),
			() => autocompleteDiscoverSeriesFacet('genres', ''),
			() => listDiscoverBlacklist()
		])
			await expect(read()).rejects.toThrow('read failed');

		client.POST.mockResolvedValue({ error: { message: 'write failed' } });
		await expect(searchMedia({} as never)).rejects.toThrow('write failed');
		await expect(addDiscoverBlacklistItem({} as never)).rejects.toThrow('write failed');
		client.DELETE.mockResolvedValue({ error: { message: 'delete failed' } });
		await expect(deleteDiscoverBlacklistItem('id')).rejects.toThrow('delete failed');

		client.GET.mockResolvedValueOnce({});
		await expect(loadMediaDiscoverSection('popular')).rejects.toThrow('did not return');
		client.POST.mockResolvedValueOnce({});
		await expect(addDiscoverBlacklistItem({} as never)).rejects.toThrow('did not return');
	});
});
