import type { QueryClient } from '@tanstack/svelte-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { DiscoverBlacklistItem } from './api';
import { createDiscoverBlacklistCache } from './cache';
import { discoverBlacklistKey } from './query.svelte';

const client = {
	getQueryData: vi.fn(),
	setQueryData: vi.fn(),
	removeQueries: vi.fn()
};

describe('discover blacklist cache', () => {
	beforeEach(() => vi.clearAllMocks());

	it('reads, upserts by provider identity, and removes entries', () => {
		const cache = createDiscoverBlacklistCache(client as unknown as QueryClient);
		const old = item('old', '123');
		const other = item('other', '456');
		client.getQueryData.mockReturnValueOnce({ items: [old] }).mockReturnValueOnce(undefined);
		expect(cache.items()).toEqual([old]);
		expect(cache.items()).toEqual([]);

		const updated = item('new', '123');
		cache.upsert(updated);
		expect(runUpdater({ items: [old, other] })).toEqual({ items: [updated, other] });

		cache.remove('other');
		expect(runUpdater({ items: [updated, other] })).toEqual({ items: [updated] });
	});

	it('initializes missing cache data and clears the feature key', () => {
		const cache = createDiscoverBlacklistCache(client as unknown as QueryClient);
		cache.upsert(item('new', '123'));
		expect(runUpdater(undefined)).toEqual({ items: [item('new', '123')] });
		cache.remove('missing');
		expect(runUpdater(undefined)).toEqual({ items: [] });

		cache.clear();
		expect(client.removeQueries).toHaveBeenCalledWith({ queryKey: discoverBlacklistKey });
	});
});

function runUpdater(value: unknown) {
	const update = client.setQueryData.mock.calls.at(-1)?.[1] as (_value: unknown) => unknown;
	return update(value);
}

function item(id: string, externalId: string) {
	return {
		id,
		externalId,
		externalProvider: 'tmdb',
		type: 'movie',
		title: id
	} as DiscoverBlacklistItem;
}
