import type { QueryClient } from '@tanstack/svelte-query';
import { searchKeys } from './queries.svelte';

export function createSearchCache(client: QueryClient) {
	const clear = () => client.removeQueries({ queryKey: searchKeys.all });
	return { clear };
}
