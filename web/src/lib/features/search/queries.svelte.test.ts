import { describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
	createQuery: vi.fn((options: () => unknown) => options()),
	autocomplete: vi.fn(),
	advanced: vi.fn()
}));
vi.mock('@tanstack/svelte-query', () => ({ createQuery: mocks.createQuery }));
vi.mock('./api', () => ({
	autocompleteMedia: mocks.autocomplete,
	advancedSearchMedia: mocks.advanced
}));

import { createAdvancedSearchQuery, createAutocompleteQuery, searchKeys } from './queries.svelte';

describe('search queries', () => {
	it('keys and executes library-only autocomplete with cancellation', () => {
		const query = createAutocompleteQuery(
			() => 'matrix',
			() => true
		) as unknown as QueryOptions;
		const signal = new AbortController().signal;
		query.queryFn({ signal });
		expect(query.queryKey).toEqual(searchKeys.autocomplete('matrix'));
		expect(query.enabled).toBe(true);
		expect(mocks.autocomplete).toHaveBeenCalledWith(
			{ query: 'matrix', includeLibrary: true, includeProviders: false },
			{ signal }
		);
		expect(query.select({ groups: ['library'] })).toEqual(['library']);
	});

	it('disables short autocomplete and missing advanced requests', () => {
		const autocomplete = createAutocompleteQuery(
			() => 'x',
			() => true
		) as unknown as QueryOptions;
		const advanced = createAdvancedSearchQuery(
			() => undefined,
			() => true
		) as unknown as QueryOptions;
		expect(autocomplete.enabled).toBe(false);
		expect(advanced.enabled).toBe(false);
	});

	it('executes POST-based advanced search as a query', () => {
		const request = { query: 'matrix', type: 'movie' } as const;
		const query = createAdvancedSearchQuery(
			() => request,
			() => true
		) as unknown as QueryOptions;
		const signal = new AbortController().signal;
		query.queryFn({ signal });
		expect(query.queryKey).toEqual(searchKeys.advanced(request));
		expect(mocks.advanced).toHaveBeenCalledWith(request, { signal });
		expect(query.select({})).toEqual([]);
	});
});

type QueryOptions = {
	queryKey: readonly unknown[];
	enabled: boolean;
	queryFn: (_context: { signal: AbortSignal }) => unknown;
	select: (_response: { groups?: unknown[] }) => unknown[];
};
