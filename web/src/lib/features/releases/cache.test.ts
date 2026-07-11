import type { QueryClient } from '@tanstack/svelte-query';
import { describe, expect, it, vi } from 'vitest';

const searchMediaReleases = vi.hoisted(() => vi.fn());
vi.mock('./api', () => ({ searchMediaReleases }));

import { createReleaseCache, releaseKeys } from './cache';

describe('release cache', () => {
	it('stores, loads, removes, and clears release results', async () => {
		const client = {
			setQueryData: vi.fn(),
			fetchQuery: vi.fn(async ({ queryFn }) => queryFn()),
			removeQueries: vi.fn()
		};
		const cache = createReleaseCache(client as unknown as QueryClient);
		const pending = { loaded: false, releases: [], errors: ['queued'] };

		cache.set('media-1', pending);
		expect(client.setQueryData).toHaveBeenCalledWith(releaseKeys.results('media-1'), pending);

		searchMediaReleases.mockResolvedValueOnce({ releases: [{ id: 'release-1' }], errors: [] });
		await expect(cache.load('media-1')).resolves.toEqual({
			releases: [{ id: 'release-1' }],
			errors: []
		});
		expect(searchMediaReleases).toHaveBeenCalledWith('media-1');

		cache.remove('media-1');
		expect(client.removeQueries).toHaveBeenCalledWith({
			queryKey: releaseKeys.results('media-1')
		});
		cache.clear();
		expect(client.removeQueries).toHaveBeenCalledWith({ queryKey: releaseKeys.all });
	});
});
