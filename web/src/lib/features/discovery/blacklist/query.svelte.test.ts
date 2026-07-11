import { describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
	createQuery: vi.fn((options: () => unknown) => options()),
	list: vi.fn()
}));

vi.mock('@tanstack/svelte-query', () => ({ createQuery: mocks.createQuery }));
vi.mock('./api', () => ({ listDiscoverBlacklist: mocks.list }));

import { createDiscoverBlacklistQuery, discoverBlacklistKey } from './query.svelte';

describe('discover blacklist query', () => {
	it('uses a feature key, admin gate, cancellation, and item selection', () => {
		const query = createDiscoverBlacklistQuery(() => true) as unknown as {
			queryKey: readonly string[];
			enabled: boolean;
			queryFn: (_context: { signal: AbortSignal }) => unknown;
			select: (_response: { items?: unknown[] }) => unknown[];
		};
		const signal = new AbortController().signal;
		query.queryFn({ signal });

		expect(query.queryKey).toEqual(discoverBlacklistKey);
		expect(query.enabled).toBe(true);
		expect(mocks.list).toHaveBeenCalledWith({ signal });
		expect(query.select({ items: ['blocked'] })).toEqual(['blocked']);
		expect(query.select({})).toEqual([]);
	});
});
