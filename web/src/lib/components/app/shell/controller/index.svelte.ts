import { useQueryClient } from '@tanstack/svelte-query';
import { createActivityCache } from '$lib/features/activity/cache';
import { createLibraryCache } from '$lib/features/library/cache';
import { createReleaseCache } from '$lib/features/releases/cache';
import { createDiscoverBlacklistCache } from '$lib/features/discovery/blacklist/cache';
import { createDiscoverBlacklistQuery } from '$lib/features/discovery/blacklist/query.svelte';
import { createDiscoverContentCache } from '$lib/features/discovery/content/cache';
import { createSearchCache } from '$lib/features/search/cache';
import {
	createAdvancedSearchQuery,
	createAutocompleteQuery
} from '$lib/features/search/queries.svelte';
import type { MediaAdvancedSearchRequest } from '$lib/features/search/api';
import { createControllerResourceRuntime } from './controllerResourceRuntime.svelte';
import { createControllerActions } from './controllerActions.svelte';
import {
	createDiscoverSectionQuery,
	createDiscoverSectionsQuery
} from '$lib/features/discovery/content/query.svelte';
import { AppShellState } from './state.svelte';
import { defaultRouteState, type AppRouteState } from './routeState';

export type { PeopleSectionKind, RelatedSectionKind } from './types';
export type { AppRouteState } from './routeState';

export function createAppShellController(route: AppRouteState = defaultRouteState()) {
	const state = new AppShellState(route);
	const resources = createControllerResourceRuntime(state, useQueryClient());
	const activityCache = createActivityCache(useQueryClient());
	const libraryCache = createLibraryCache(useQueryClient());
	const releaseCache = createReleaseCache(useQueryClient());
	const discoverBlacklistCache = createDiscoverBlacklistCache(useQueryClient());
	const discoverBlacklist = createDiscoverBlacklistQuery(() => state.isAdmin);
	const discoverContentCache = createDiscoverContentCache(useQueryClient());
	const discoverSections = createDiscoverSectionsQuery(
		() =>
			state.authenticated && state.activeView === 'home' && state.activeHomeSection === 'discover'
	);
	const discoverSection = createDiscoverSectionQuery(
		() => state.activeDiscoverSectionId,
		() => state.authenticated && state.activeView === 'discover-section'
	);
	let autocompleteRequest = $state('');
	let advancedRequest = $state<MediaAdvancedSearchRequest | undefined>();
	const searchCache = createSearchCache(useQueryClient());
	const autocomplete = createAutocompleteQuery(
		() => autocompleteRequest,
		() => state.authenticated
	);
	const advancedSearch = createAdvancedSearchQuery(
		() => advancedRequest,
		() => state.authenticated && state.activeView === 'advanced-search'
	);
	const actions = createControllerActions(state, {
		activity: activityCache,
		library: libraryCache,
		releases: releaseCache,
		blacklist: discoverBlacklistCache,
		discovery: discoverContentCache,
		search: searchCache,
		resources,
		setAutocomplete: (query) => (autocompleteRequest = query),
		setAdvanced: (request) => (advancedRequest = request)
	});
	const controller = Object.assign(state, ...actions);
	return Object.defineProperties(controller, {
		...resources.properties,
		discoverBlacklist: { get: () => discoverBlacklist.data ?? [] },
		loadingBlacklist: { get: () => discoverBlacklist.isFetching },
		discoverSections: { get: () => discoverSections.data ?? [] },
		loadingDiscover: { get: () => discoverSections.isFetching },
		discoverSection: { get: () => discoverSection.data?.section },
		loadingDiscoverSection: { get: () => discoverSection.isFetching },
		loadingMoreDiscoverSection: { get: () => discoverSection.data?.loadingMore ?? false },
		discoverSectionHasMore: { get: () => discoverSection.data?.hasMore ?? false },
		autocompleteGroups: { get: () => autocomplete.data ?? [] },
		loadingAutocomplete: { get: () => autocomplete.isFetching },
		advancedSearchGroups: { get: () => advancedSearch.data ?? [] },
		searchingAdvanced: { get: () => advancedSearch.isFetching }
	});
}
