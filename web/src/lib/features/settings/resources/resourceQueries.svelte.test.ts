import { beforeEach, describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
	createQuery: vi.fn((options: () => unknown) => options()),
	createMutation: vi.fn((options: () => unknown) => options()),
	client: { setQueryData: vi.fn(), fetchQuery: vi.fn() },
	api: vi.fn()
}));

vi.mock('@tanstack/svelte-query', () => ({
	createQuery: mocks.createQuery,
	createMutation: mocks.createMutation,
	useQueryClient: () => mocks.client
}));
vi.mock('$lib/components/settings/library/filePoliciesApi', () => ({
	getFileNamingSettings: mocks.api,
	updateFileNamingSettings: mocks.api,
	getFileDeleteSettings: mocks.api,
	updateFileDeleteSettings: mocks.api
}));
vi.mock('$lib/settings/api', () => ({
	listIndexerCatalog: mocks.api,
	listIndexerAppProfiles: mocks.api,
	listIndexerProxies: mocks.api
}));

import {
	createFileDeleteResource,
	createFileNamingResource,
	filePolicyKeys
} from './filePolicies.svelte';
import { createIndexerAuxiliaryQueries, indexerAuxiliaryKeys } from './indexerAuxiliary.svelte';

describe('settings resources', () => {
	beforeEach(() => vi.clearAllMocks());

	it('uses stable, separate file-policy keys and reconciles saves', () => {
		const naming = createFileNamingResource() as unknown as {
			query: { queryKey: readonly string[] };
			save: { onSuccess: (data: unknown) => void };
		};
		const deletion = createFileDeleteResource() as unknown as {
			query: { queryKey: readonly string[] };
		};
		expect(naming.query.queryKey).toEqual(filePolicyKeys.naming());
		expect(deletion.query.queryKey).toEqual(filePolicyKeys.deletion());
		naming.save.onSuccess({ movieFileFormat: 'movie' });
		expect(mocks.client.setQueryData).toHaveBeenCalledWith(filePolicyKeys.naming(), {
			movieFileFormat: 'movie'
		});
	});

	it('creates independently cached indexer auxiliary reads', () => {
		const queries = createIndexerAuxiliaryQueries() as unknown as Record<
			'catalog' | 'appProfiles' | 'proxies',
			{ queryKey: readonly string[] }
		>;
		expect(queries.catalog.queryKey).toEqual(indexerAuxiliaryKeys.catalog());
		expect(queries.appProfiles.queryKey).toEqual(indexerAuxiliaryKeys.appProfiles());
		expect(queries.proxies.queryKey).toEqual(indexerAuxiliaryKeys.proxies());
	});
});
