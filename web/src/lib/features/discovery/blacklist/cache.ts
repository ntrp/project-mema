import type { QueryClient } from '@tanstack/svelte-query';
import type { DiscoverBlacklistItem, DiscoverBlacklistResponse } from './api';
import { discoverBlacklistKey } from './query.svelte';

export function createDiscoverBlacklistCache(client: QueryClient) {
	const items = () =>
		client.getQueryData<DiscoverBlacklistResponse>(discoverBlacklistKey)?.items ?? [];
	const upsert = (item: DiscoverBlacklistItem) =>
		client.setQueryData<DiscoverBlacklistResponse>(discoverBlacklistKey, (current) => ({
			items: [item, ...(current?.items ?? []).filter((entry) => !sameEntry(entry, item))]
		}));
	const remove = (id: string) =>
		client.setQueryData<DiscoverBlacklistResponse>(discoverBlacklistKey, (current) => ({
			items: (current?.items ?? []).filter((item) => item.id !== id)
		}));
	const clear = () => client.removeQueries({ queryKey: discoverBlacklistKey });
	return { items, upsert, remove, clear };
}

function sameEntry(left: DiscoverBlacklistItem, right: DiscoverBlacklistItem) {
	return (
		left.externalProvider === right.externalProvider &&
		left.externalId === right.externalId &&
		left.type === right.type
	);
}
