import type { QueryClient } from '@tanstack/svelte-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const getSection = vi.hoisted(() => vi.fn());
vi.mock('./api', async (original) => ({
	...(await original()),
	getMediaDiscoverSection: getSection
}));

import { createDiscoverContentCache } from './cache';
import { discoverContentKeys } from './query.svelte';

const data = new Map<string, unknown>();
const client = {
	getQueryData: vi.fn((key: readonly unknown[]) => data.get(JSON.stringify(key))),
	setQueryData: vi.fn((key: readonly unknown[], value: unknown) => {
		const name = JSON.stringify(key);
		data.set(name, typeof value === 'function' ? value(data.get(name)) : value);
	}),
	invalidateQueries: vi.fn(),
	removeQueries: vi.fn()
};

describe('discover content cache', () => {
	beforeEach(() => {
		data.clear();
		vi.clearAllMocks();
	});

	it('maps home and active-section results', () => {
		data.set(JSON.stringify(discoverContentKeys.sections), {
			sections: [{ id: 'popular', results: [] }]
		});
		data.set(JSON.stringify(discoverContentKeys.section('popular')), entry([{ id: 'one' }]));
		const cache = createDiscoverContentCache(client as unknown as QueryClient);

		cache.mapSections((sections) => [...sections, { id: 'new' } as never]);
		cache.mapSection('popular', (section) => ({ ...section, title: 'Updated' }));
		expect(cache.sections()).toHaveLength(2);
		expect(cache.section('popular')?.section.title).toBe('Updated');
	});

	it('deduplicates pagination, advances metadata, and stops after an empty page', async () => {
		data.set(JSON.stringify(discoverContentKeys.section('popular')), entry([result('one')]));
		getSection.mockResolvedValueOnce({
			id: 'popular',
			results: [result('one'), result('two')]
		});
		const cache = createDiscoverContentCache(client as unknown as QueryClient);

		await cache.loadMore('popular');
		expect(cache.section('popular')).toMatchObject({ page: 2, hasMore: true, loadingMore: false });
		expect(cache.section('popular')?.section.results.map((item) => item.id)).toEqual([
			'one',
			'two'
		]);

		getSection.mockResolvedValueOnce({ id: 'popular', results: [result('one')] });
		await cache.loadMore('popular');
		expect(cache.section('popular')).toMatchObject({ page: 3, hasMore: false });
	});

	it('refreshes and clears only discover content keys', () => {
		const cache = createDiscoverContentCache(client as unknown as QueryClient);
		cache.refresh('popular');
		cache.clear();
		expect(client.invalidateQueries).toHaveBeenCalledTimes(2);
		expect(client.removeQueries).toHaveBeenCalledWith({ queryKey: discoverContentKeys.sections });
		expect(client.removeQueries).toHaveBeenCalledWith({ queryKey: ['discovery', 'section'] });
	});
});

function entry(results: Array<{ id: string }>) {
	return { section: { id: 'popular', results }, page: 1, hasMore: true, loadingMore: false };
}

function result(id: string) {
	return { id, externalProvider: 'tmdb', externalId: id, type: 'movie', title: id };
}
