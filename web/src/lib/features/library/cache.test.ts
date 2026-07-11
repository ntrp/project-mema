import type { QueryClient } from '@tanstack/svelte-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import type { MediaItem, MediaRequest } from './api';
import { createLibraryCache } from './cache';
import { libraryKeys } from './queries.svelte';

const client = {
	getQueryData: vi.fn(),
	setQueryData: vi.fn(),
	invalidateQueries: vi.fn(),
	removeQueries: vi.fn()
};

describe('library cache', () => {
	beforeEach(() => vi.clearAllMocks());

	it('reads items and defaults an empty cache', () => {
		const cache = createLibraryCache(client as unknown as QueryClient);
		client.getQueryData
			.mockReturnValueOnce({ items: [item('one')] })
			.mockReturnValueOnce(undefined);

		expect(cache.items()).toEqual([item('one')]);
		expect(cache.items()).toEqual([]);
		expect(client.getQueryData).toHaveBeenCalledWith(libraryKeys.items());
	});

	it('upserts, maps, and removes media items', () => {
		const cache = createLibraryCache(client as unknown as QueryClient);
		const current = { items: [item('one', 'old'), item('two')] };

		cache.upsertItem(item('one', 'new'));
		expect(runUpdater(current)).toEqual({ items: [item('one', 'new'), item('two')] });

		cache.mapItems((entry) => ({ ...entry, title: `${entry.title}!` }));
		expect(runUpdater(current)).toEqual({ items: [item('one', 'old!'), item('two', 'two!')] });

		cache.removeItem('one');
		expect(runUpdater(current)).toEqual({ items: [item('two')] });
	});

	it('initializes an empty item cache before applying an update', () => {
		const cache = createLibraryCache(client as unknown as QueryClient);
		cache.upsertItem(item('new'));

		expect(runUpdater(undefined)).toEqual({ items: [item('new')] });
	});

	it('upserts and maps media requests', () => {
		const cache = createLibraryCache(client as unknown as QueryClient);
		const current = { requests: [request('one', 'old'), request('two')] };

		cache.upsertRequest(request('one', 'new'));
		expect(runUpdater(current)).toEqual({
			requests: [request('one', 'new'), request('two')]
		});

		cache.mapRequests((entry) => ({ ...entry, title: `${entry.title}!` }));
		expect(runUpdater(current)).toEqual({
			requests: [request('one', 'old!'), request('two', 'two!')]
		});
	});

	it('defaults missing requests and exposes refresh and clear operations', () => {
		const cache = createLibraryCache(client as unknown as QueryClient);
		cache.mapRequests((entry) => entry);
		expect(runUpdater(undefined)).toEqual({ requests: [] });

		cache.refreshItems();
		cache.clear();
		expect(client.invalidateQueries).toHaveBeenCalledWith({ queryKey: libraryKeys.items() });
		expect(client.removeQueries).toHaveBeenCalledWith({ queryKey: libraryKeys.all });
	});
});

function runUpdater(value: unknown) {
	const updater = client.setQueryData.mock.calls.at(-1)?.[1] as (current: unknown) => unknown;
	return updater(value);
}

function item(id: string, title = id) {
	return { id, title } as MediaItem;
}

function request(id: string, title = id) {
	return { id, title } as MediaRequest;
}
