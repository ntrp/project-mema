import { describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
	createQuery: vi.fn((options: () => unknown) => options()),
	sections: vi.fn(),
	section: vi.fn()
}));
vi.mock('@tanstack/svelte-query', () => ({ createQuery: mocks.createQuery }));
vi.mock('./api', () => ({
	getMediaDiscover: mocks.sections,
	getMediaDiscoverSection: mocks.section
}));

import {
	createDiscoverSectionQuery,
	createDiscoverSectionsQuery,
	discoverContentKeys
} from './query.svelte';

describe('discover content queries', () => {
	it('loads and selects the home sections', () => {
		const query = createDiscoverSectionsQuery(() => true) as unknown as QueryOptions;
		const signal = new AbortController().signal;
		query.queryFn({ signal });
		expect(query.queryKey).toEqual(discoverContentKeys.sections);
		expect(mocks.sections).toHaveBeenCalledWith({ signal });
		expect(query.select?.({ sections: ['popular'] })).toEqual(['popular']);
	});

	it('loads page one for the active section and owns pagination metadata', async () => {
		mocks.section.mockResolvedValue({ id: 'popular', results: [{ id: 'one' }] });
		const query = createDiscoverSectionQuery(
			() => 'popular',
			() => true
		) as unknown as QueryOptions;
		const signal = new AbortController().signal;
		await expect(query.queryFn({ signal })).resolves.toMatchObject({
			page: 1,
			hasMore: true,
			loadingMore: false
		});
		expect(query.queryKey).toEqual(discoverContentKeys.section('popular'));
		expect(mocks.section).toHaveBeenCalledWith('popular', { page: 1 }, { signal });
	});
});

type QueryOptions = {
	queryKey: readonly unknown[];
	queryFn: (_context: { signal: AbortSignal }) => unknown;
	select?: (_response: { sections?: unknown[] }) => unknown[];
};
