import {
	addDiscoverBlacklistItem as addDiscoverBlacklistItemRequest,
	deleteDiscoverBlacklistItem as deleteDiscoverBlacklistItemRequest,
	listDiscoverBlacklist as listDiscoverBlacklistRequest,
	loadMediaDiscoverSection as loadMediaDiscoverSectionRequest,
	loadMediaDiscoverSections as loadMediaDiscoverSectionsRequest
} from '$lib/settings/api';
import type { DiscoverBlacklistItem, MediaSearchResult } from '$lib/settings/types';
import {
	discoverResultKey,
	filterDiscoverSection,
	filterDiscoverSections,
	sameDiscoverBlacklistItem
} from './discoverFilters';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

export function createDiscoveryActions(state: AppShellState) {
	async function loadDiscoverSections() {
		state.loadingDiscover = true;
		try {
			state.discoverSections = await loadMediaDiscoverSectionsRequest();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load discover sections');
		} finally {
			state.loadingDiscover = false;
		}
	}

	async function loadDiscoverSection() {
		if (!state.activeDiscoverSectionId) {
			return;
		}
		state.loadingDiscoverSection = true;
		state.discoverSectionPage = 1;
		state.discoverSectionHasMore = true;
		try {
			state.discoverSection = await loadMediaDiscoverSectionRequest(
				state.activeDiscoverSectionId,
				1
			);
			state.discoverSectionHasMore = (state.discoverSection.results ?? []).length > 0;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load discover section');
		} finally {
			state.loadingDiscoverSection = false;
		}
	}

	async function loadMoreDiscoverSection() {
		if (
			!state.activeDiscoverSectionId ||
			state.loadingDiscoverSection ||
			state.loadingMoreDiscoverSection ||
			!state.discoverSectionHasMore
		) {
			return;
		}
		state.loadingMoreDiscoverSection = true;
		const nextPage = state.discoverSectionPage + 1;
		try {
			const nextSection = await loadMediaDiscoverSectionRequest(
				state.activeDiscoverSectionId,
				nextPage
			);
			const existingKeys = (state.discoverSection?.results ?? []).map(discoverResultKey);
			const nextResults = (nextSection.results ?? []).filter((result) => {
				const key = discoverResultKey(result);
				if (existingKeys.includes(key)) {
					return false;
				}
				existingKeys.push(key);
				return true;
			});
			state.discoverSectionPage = nextPage;
			state.discoverSectionHasMore = nextResults.length > 0;
			state.discoverSection = {
				...nextSection,
				results: [...(state.discoverSection?.results ?? []), ...nextResults]
			};
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load more discover results');
		} finally {
			state.loadingMoreDiscoverSection = false;
		}
	}

	async function loadDiscoverBlacklist() {
		if (!state.isAdmin) {
			return;
		}
		state.loadingBlacklist = true;
		try {
			state.discoverBlacklist = await listDiscoverBlacklistRequest();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load discover blacklist');
		} finally {
			state.loadingBlacklist = false;
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
			state.discoverBlacklist = [
				item,
				...state.discoverBlacklist.filter((entry) => !sameDiscoverBlacklistItem(entry, item))
			];
			state.discoverSections = filterDiscoverSections(
				state.discoverSections,
				state.discoverBlacklist
			);
			if (state.discoverSection) {
				state.discoverSection = filterDiscoverSection(
					state.discoverSection,
					state.discoverBlacklist
				);
			}
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
			state.discoverBlacklist = state.discoverBlacklist.filter((entry) => entry.id !== item.id);
			state.message = `${item.title} removed from blacklist`;
			await loadDiscoverSections();
			if (state.activeView === 'discover-section') {
				await loadDiscoverSection();
			}
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
		loadDiscoverSections,
		loadDiscoverSection,
		loadMoreDiscoverSection,
		loadDiscoverBlacklist,
		blacklistDiscoverMedia,
		removeDiscoverBlacklistItem
	};
}
