import { createQuery } from '@tanstack/svelte-query';
import { listSubtitleProviderCatalog } from '$lib/settings/api';

export const subtitleProviderAuxiliaryKeys = {
	catalog: () => ['settings', 'subtitle-providers', 'catalog'] as const
};

export function createSubtitleProviderAuxiliaryQueries() {
	return {
		catalog: createQuery(() => ({
			queryKey: subtitleProviderAuxiliaryKeys.catalog(),
			queryFn: listSubtitleProviderCatalog
		}))
	};
}
