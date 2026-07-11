import { describe, expect, it, vi } from 'vitest';

vi.mock('@tanstack/svelte-query', () => ({
	createQuery: vi.fn((options: () => unknown) => options())
}));
vi.mock('$lib/settings/domains/cacheInspection', () => ({
	getIndexerSearch: vi.fn(),
	getMetadataCache: vi.fn()
}));
vi.mock('$lib/profile/profileApi', () => ({ getProfile: vi.fn() }));

import { createServerResourceRuntime, serverResourceKeys } from './cacheResources.svelte';

describe('server resource runtime', () => {
	it('guards queries and writes or clears their shared cache namespace', () => {
		const client = { setQueryData: vi.fn(), removeQueries: vi.fn() };
		const runtime = createServerResourceRuntime(
			client as never,
			() => true,
			() => false
		);
		expect(runtime.indexerSearch).toMatchObject({ enabled: false });
		expect(runtime.profile).toMatchObject({ enabled: true });
		runtime.setProfile({ id: 'user-1' } as never);
		expect(client.setQueryData).toHaveBeenCalledWith(serverResourceKeys.profile(), {
			id: 'user-1'
		});
		runtime.clear();
		expect(client.removeQueries).toHaveBeenCalledWith({ queryKey: serverResourceKeys.all });
	});
});
