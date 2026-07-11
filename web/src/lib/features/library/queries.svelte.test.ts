import { beforeEach, describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
	createQuery: vi.fn((options: () => unknown) => options()),
	listMediaItems: vi.fn(),
	listMediaRequests: vi.fn()
}));

vi.mock('@tanstack/svelte-query', () => ({ createQuery: mocks.createQuery }));
vi.mock('./api', () => ({
	listMediaItems: mocks.listMediaItems,
	listMediaRequests: mocks.listMediaRequests
}));

import { createMediaItemsQuery, createMediaRequestsQuery, libraryKeys } from './queries.svelte';

describe('library queries', () => {
	beforeEach(() => vi.clearAllMocks());

	it('builds stable feature-owned keys', () => {
		expect(libraryKeys.items()).toEqual(['library', 'items']);
		expect(libraryKeys.requests()).toEqual(['library', 'requests']);
	});

	it('configures media item loading and selection', () => {
		const query = createMediaItemsQuery() as unknown as QueryOptions;
		const signal = new AbortController().signal;

		query.queryFn({ signal });

		expect(query.queryKey).toEqual(libraryKeys.items());
		expect(mocks.listMediaItems).toHaveBeenCalledWith({ signal });
		expect(query.select({ items: ['item'] })).toEqual(['item']);
	});

	it('configures media request loading and selection', () => {
		const query = createMediaRequestsQuery() as unknown as QueryOptions;
		const signal = new AbortController().signal;

		query.queryFn({ signal });

		expect(query.queryKey).toEqual(libraryKeys.requests());
		expect(mocks.listMediaRequests).toHaveBeenCalledWith({ signal });
		expect(query.select({ requests: ['request'] })).toEqual(['request']);
	});
});

type QueryOptions = {
	queryKey: readonly string[];
	queryFn: (context: { signal: AbortSignal }) => unknown;
	select: (response: { items?: unknown[]; requests?: unknown[] }) => unknown[];
};
