import {
	addDiscoverBlacklistItem as addDiscoverBlacklistItemRequest,
	deleteDiscoverBlacklistItem as deleteDiscoverBlacklistItemRequest
} from '$lib/features/discovery/blacklist/api';
import type { DiscoverBlacklistItem, MediaSearchResult } from '$lib/settings/types';
import type { MediaDiscoverSection } from '$lib/features/discovery/content/api';
import {
	discoverResultKey,
	filterDiscoverSection,
	filterDiscoverSections
} from './discoverFilters';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

export interface DiscoveryBlacklistCache {
	items: () => DiscoverBlacklistItem[];
	upsert: (_item: DiscoverBlacklistItem) => void;
	remove: (_id: string) => void;
}

export interface DiscoveryContentCache {
	mapSections: (_map: (_sections: MediaDiscoverSection[]) => MediaDiscoverSection[]) => void;
	mapSection: (_id: string, _map: (_section: MediaDiscoverSection) => MediaDiscoverSection) => void;
	refresh: (_id?: string) => void;
	loadMore: (_id: string) => Promise<void>;
}

export function createDiscoveryActions(
	state: AppShellState,
	blacklist: DiscoveryBlacklistCache,
	content: DiscoveryContentCache
) {
	async function loadMoreDiscoverSection() {
		if (!state.activeDiscoverSectionId) return;
		try {
			await content.loadMore(state.activeDiscoverSectionId);
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load more discover results');
		}
	}

	async function blacklistDiscoverMedia(candidate: MediaSearchResult) {
		if (!state.isAdmin) {
			return;
		}
		state.blacklistingKey = discoverResultKey(candidate);
		try {
			const item = await addDiscoverBlacklistItemRequest({
				title: candidate.title,
				type: candidate.type,
				year: candidate.year,
				externalProvider: candidate.externalProvider,
				externalId: candidate.externalId,
				overview: candidate.overview,
				posterPath: candidate.posterPath
			});
			blacklist.upsert(item);
			const blacklistItems = blacklist.items();
			content.mapSections((sections) => filterDiscoverSections(sections, blacklistItems));
			if (state.activeDiscoverSectionId)
				content.mapSection(state.activeDiscoverSectionId, (section) =>
					filterDiscoverSection(section, blacklistItems)
				);
			state.message = `${candidate.title} hidden from discover`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not add media to discover blacklist');
		} finally {
			state.blacklistingKey = undefined;
		}
	}

	async function removeDiscoverBlacklistItem(item: DiscoverBlacklistItem) {
		if (!state.isAdmin) {
			return;
		}
		state.removingBlacklistId = item.id;
		try {
			await deleteDiscoverBlacklistItemRequest(item.id);
			blacklist.remove(item.id);
			state.message = `${item.title} removed from blacklist`;
			content.refresh(state.activeDiscoverSectionId);
		} catch (error) {
			state.errorMessage = errorMessageFrom(
				error,
				'Could not remove media from discover blacklist'
			);
		} finally {
			state.removingBlacklistId = undefined;
		}
	}

	return {
		loadMoreDiscoverSection,
		blacklistDiscoverMedia,
		removeDiscoverBlacklistItem
	};
}
