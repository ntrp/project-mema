import { createQuery } from '@tanstack/svelte-query';
import { getMediaDiscover, getMediaDiscoverSection, type MediaDiscoverSection } from './api';

export interface DiscoverSectionCacheEntry {
	section: MediaDiscoverSection;
	page: number;
	hasMore: boolean;
	loadingMore: boolean;
}

export const discoverContentKeys = {
	sections: ['discovery', 'sections'] as const,
	section: (id: string) => ['discovery', 'section', id] as const,
	page: (id: string, page: number) => [...discoverContentKeys.section(id), 'page', page] as const
};

export function createDiscoverSectionsQuery(enabled: () => boolean) {
	return createQuery(() => ({
		queryKey: discoverContentKeys.sections,
		queryFn: ({ signal }) => getMediaDiscover({ signal }),
		select: (response) => response.sections ?? [],
		enabled: enabled()
	}));
}

export function createDiscoverSectionQuery(id: () => string | undefined, enabled: () => boolean) {
	return createQuery(() => ({
		queryKey: discoverContentKeys.section(id() ?? ''),
		queryFn: async ({ signal }) =>
			entry(await getMediaDiscoverSection(id() ?? '', { page: 1 }, { signal })),
		enabled: enabled() && Boolean(id())
	}));
}

function entry(section: MediaDiscoverSection): DiscoverSectionCacheEntry {
	return { section, page: 1, hasMore: (section.results ?? []).length > 0, loadingMore: false };
}
