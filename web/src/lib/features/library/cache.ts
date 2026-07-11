import type { QueryClient } from '@tanstack/svelte-query';
import type { MediaItem, MediaRequest } from './api';
import { libraryKeys } from './queries.svelte';

type ItemsResponse = { items: MediaItem[] };
type RequestsResponse = { requests: MediaRequest[] };

export function createLibraryCache(client: QueryClient) {
	const items = () => client.getQueryData<ItemsResponse>(libraryKeys.items())?.items ?? [];
	const updateItems = (fn: (_items: MediaItem[]) => MediaItem[]) =>
		client.setQueryData<ItemsResponse>(libraryKeys.items(), (value) => ({
			items: fn(value?.items ?? [])
		}));
	const upsertItem = (item: MediaItem) =>
		updateItems((current) => [item, ...current.filter((entry) => entry.id !== item.id)]);
	const mapItems = (fn: (_item: MediaItem) => MediaItem) =>
		updateItems((current) => current.map(fn));
	const removeItem = (id: string) =>
		updateItems((current) => current.filter((item) => item.id !== id));

	const upsertRequest = (request: MediaRequest) =>
		client.setQueryData<RequestsResponse>(libraryKeys.requests(), (value) => ({
			requests: [request, ...(value?.requests ?? []).filter((entry) => entry.id !== request.id)]
		}));
	const mapRequests = (fn: (_request: MediaRequest) => MediaRequest) =>
		client.setQueryData<RequestsResponse>(libraryKeys.requests(), (value) => ({
			requests: (value?.requests ?? []).map(fn)
		}));

	const refreshItems = () => client.invalidateQueries({ queryKey: libraryKeys.items() });
	const clear = () => client.removeQueries({ queryKey: libraryKeys.all });
	return {
		items,
		upsertItem,
		mapItems,
		removeItem,
		upsertRequest,
		mapRequests,
		refreshItems,
		clear
	};
}
