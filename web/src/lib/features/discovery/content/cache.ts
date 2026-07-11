import type { QueryClient } from '@tanstack/svelte-query';
import {
	getMediaDiscoverSection,
	type MediaDiscoverResponse,
	type MediaDiscoverSection
} from './api';
import { discoverContentKeys, type DiscoverSectionCacheEntry } from './query.svelte';

export function createDiscoverContentCache(client: QueryClient) {
	const sections = () =>
		client.getQueryData<MediaDiscoverResponse>(discoverContentKeys.sections)?.sections ?? [];
	const section = (id: string) =>
		client.getQueryData<DiscoverSectionCacheEntry>(discoverContentKeys.section(id));
	const mapSections = (map: (_sections: MediaDiscoverSection[]) => MediaDiscoverSection[]) =>
		client.setQueryData<MediaDiscoverResponse>(discoverContentKeys.sections, (current) => ({
			sections: map(current?.sections ?? [])
		}));
	const mapSection = (id: string, map: (_section: MediaDiscoverSection) => MediaDiscoverSection) =>
		client.setQueryData<DiscoverSectionCacheEntry>(discoverContentKeys.section(id), (current) =>
			current ? { ...current, section: map(current.section) } : current
		);
	const refresh = (id?: string) => {
		void client.invalidateQueries({ queryKey: discoverContentKeys.sections });
		if (id) void client.invalidateQueries({ queryKey: discoverContentKeys.section(id) });
	};
	const loadMore = (id: string) => loadMoreSection(client, id, section(id));
	const clear = () => {
		client.removeQueries({ queryKey: discoverContentKeys.sections });
		client.removeQueries({ queryKey: ['discovery', 'section'] });
	};
	return { sections, section, mapSections, mapSection, refresh, loadMore, clear };
}

async function loadMoreSection(
	client: QueryClient,
	id: string,
	current?: DiscoverSectionCacheEntry
) {
	if (!current || current.loadingMore || !current.hasMore) return;
	setLoading(client, id, current, true);
	try {
		const next = await getMediaDiscoverSection(id, { page: current.page + 1 });
		const known = new Set((current.section.results ?? []).map(resultKey));
		const additions = (next.results ?? []).filter((item) => !known.has(resultKey(item)));
		client.setQueryData<DiscoverSectionCacheEntry>(discoverContentKeys.section(id), {
			section: { ...next, results: [...(current.section.results ?? []), ...additions] },
			page: current.page + 1,
			hasMore: additions.length > 0,
			loadingMore: false
		});
	} catch (error) {
		setLoading(client, id, current, false);
		throw error;
	}
}

function setLoading(
	client: QueryClient,
	id: string,
	current: DiscoverSectionCacheEntry,
	loadingMore: boolean
) {
	client.setQueryData(discoverContentKeys.section(id), { ...current, loadingMore });
}

function resultKey(item: {
	externalProvider?: string;
	externalId?: string;
	type?: string;
	title?: string;
}) {
	return `${item.externalProvider ?? ''}:${item.externalId ?? ''}:${item.type ?? ''}:${item.title ?? ''}`;
}
