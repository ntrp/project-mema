import { beforeEach, describe, expect, it, vi } from 'vitest';

const { createQuery, listDownloadActivity, listReleaseBlocklist } = vi.hoisted(() => ({
	createQuery: vi.fn((options: () => unknown) => options()),
	listDownloadActivity: vi.fn(),
	listReleaseBlocklist: vi.fn()
}));

vi.mock('@tanstack/svelte-query', () => ({ createQuery }));
vi.mock('./api', () => ({ listDownloadActivity, listReleaseBlocklist }));

import {
	activityKeys,
	createDownloadActivityQuery,
	createReleaseBlocklistQuery
} from './queries.svelte';

describe('activity queries', () => {
	beforeEach(() => vi.clearAllMocks());

	it('builds stable feature-owned keys', () => {
		expect(activityKeys.downloads()).toEqual(['activity', 'downloads']);
		expect(activityKeys.blocklist()).toEqual(['activity', 'blocklist']);
	});

	it('configures the download query and forwards cancellation', async () => {
		const query = createDownloadActivityQuery(() => false) as unknown as {
			queryKey: readonly string[];
			enabled: boolean;
			queryFn: (context: { signal: AbortSignal }) => unknown;
			select: (response: { activities: unknown[] }) => unknown[];
		};
		const signal = new AbortController().signal;

		query.queryFn({ signal });

		expect(query.queryKey).toEqual(activityKeys.downloads());
		expect(query.enabled).toBe(false);
		expect(listDownloadActivity).toHaveBeenCalledWith({ signal });
		expect(query.select({ activities: ['download'] })).toEqual(['download']);
	});

	it('configures the blocklist query with an enabled default', () => {
		const query = createReleaseBlocklistQuery() as unknown as {
			queryKey: readonly string[];
			enabled: boolean;
			queryFn: (context: { signal: AbortSignal }) => unknown;
			select: (response: { items: unknown[] }) => unknown[];
		};
		const signal = new AbortController().signal;

		query.queryFn({ signal });

		expect(query.queryKey).toEqual(activityKeys.blocklist());
		expect(query.enabled).toBe(true);
		expect(listReleaseBlocklist).toHaveBeenCalledWith({ signal });
		expect(query.select({ items: ['blocked'] })).toEqual(['blocked']);
	});
});
