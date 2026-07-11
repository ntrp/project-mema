import { describe, expect, it } from 'vitest';

import { createAppQueryClient } from './queryClient';

describe('application query client', () => {
	it('uses consistent server-state defaults', () => {
		const defaults = createAppQueryClient().getDefaultOptions();

		expect(defaults.queries).toMatchObject({
			staleTime: 30_000,
			gcTime: 300_000,
			retry: 1,
			refetchOnWindowFocus: false,
			refetchOnReconnect: true
		});
		expect(defaults.mutations).toMatchObject({ retry: 0 });
	});

	it('creates isolated clients', () => {
		expect(createAppQueryClient()).not.toBe(createAppQueryClient());
	});
});
