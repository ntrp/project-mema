import type { QueryClient } from '@tanstack/svelte-query';
import type { ReleaseCandidate } from './api';
import { searchMediaReleases } from './api';

export interface ReleaseSearchResult {
	loaded: boolean;
	releases: ReleaseCandidate[];
	errors: string[];
}

export const releaseKeys = {
	all: ['releases'] as const,
	results: (mediaItemId: string) => [...releaseKeys.all, mediaItemId] as const
};

export function createReleaseCache(client: QueryClient) {
	const set = (id: string, result: ReleaseSearchResult) =>
		client.setQueryData(releaseKeys.results(id), result);
	const load = async (id: string) => {
		const result = await client.fetchQuery({
			queryKey: releaseKeys.results(id),
			queryFn: () => searchMediaReleases(id),
			staleTime: 0
		});
		return result;
	};
	const remove = (id: string) => client.removeQueries({ queryKey: releaseKeys.results(id) });
	const clear = () => client.removeQueries({ queryKey: releaseKeys.all });
	return { set, load, remove, clear };
}
