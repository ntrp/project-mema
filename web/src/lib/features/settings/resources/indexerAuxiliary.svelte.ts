import { createQuery } from '@tanstack/svelte-query';
import { listIndexerAppProfiles, listIndexerCatalog, listIndexerProxies } from '$lib/settings/api';

export const indexerAuxiliaryKeys = {
	catalog: () => ['settings', 'indexers', 'catalog'] as const,
	appProfiles: () => ['settings', 'indexers', 'app-profiles'] as const,
	proxies: () => ['settings', 'indexers', 'proxies'] as const
};

export function createIndexerAuxiliaryQueries() {
	return {
		catalog: createQuery(() => ({
			queryKey: indexerAuxiliaryKeys.catalog(),
			queryFn: listIndexerCatalog
		})),
		appProfiles: createQuery(() => ({
			queryKey: indexerAuxiliaryKeys.appProfiles(),
			queryFn: listIndexerAppProfiles
		})),
		proxies: createQuery(() => ({
			queryKey: indexerAuxiliaryKeys.proxies(),
			queryFn: listIndexerProxies
		}))
	};
}
