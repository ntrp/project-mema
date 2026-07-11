import { createQuery } from '@tanstack/svelte-query';
import { listDiscoverBlacklist } from './api';

export const discoverBlacklistKey = ['discovery', 'blacklist'] as const;

export function createDiscoverBlacklistQuery(enabled: () => boolean) {
	return createQuery(() => ({
		queryKey: discoverBlacklistKey,
		queryFn: ({ signal }) => listDiscoverBlacklist({ signal }),
		select: (response) => response.items ?? [],
		enabled: enabled()
	}));
}
