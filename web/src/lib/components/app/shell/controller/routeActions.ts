import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import type { AppShellState } from './state.svelte';
import { appRouteKey, type AppRouteState } from './routeState';

interface RouteDeps {
	loadDiscoverSection: () => Promise<void>;
	loadMediaCollection: () => Promise<void>;
	loadMetadataDetail: () => Promise<void>;
	loadPersonDetail: () => Promise<void>;
}

export function createRouteActions(state: AppShellState, deps: RouteDeps) {
	let currentRouteKey = appRouteKey(state.route);

	async function applyRoute(route: AppRouteState) {
		const nextRouteKey = appRouteKey(route);
		if (nextRouteKey === currentRouteKey) {
			return;
		}

		const previousMetadataKey = metadataKey(state.route);
		const previousPersonKey = personKey(state.route);
		const previousCollectionKey = collectionKey(state.route);
		const previousDiscoverId = state.route.discoverSectionId;
		state.route = route;
		currentRouteKey = nextRouteKey;

		state.activeView = route.view;
		state.activeHomeSection = route.homeSection;
		state.activeActivitySection = route.activitySection;
		state.activeSettingsSection = route.settingsSection;
		state.activeSystemSection = route.systemSection;
		state.activeDiscoverSectionId = route.discoverSectionId;
		state.activeDiscoverSubmenuSection = route.discoverSubmenuSection;
		state.activeRelatedSectionKind = route.relatedSectionKind;
		state.activePeopleSectionKind = route.peopleSectionKind;
		state.selectedMediaItemId = route.selectedMediaItemId;
		state.selectedRequestId = route.selectedRequestId;
		state.searchQuery = route.view === 'advanced-search' ? route.advancedQuery : state.searchQuery;

		if (metadataKey(route) !== previousMetadataKey) {
			state.metadataDetail = undefined;
		}
		if (personKey(route) !== previousPersonKey) {
			state.personDetail = undefined;
		}
		if (collectionKey(route) !== previousCollectionKey) {
			state.mediaCollection = undefined;
		}
		if (route.discoverSectionId !== previousDiscoverId) {
			state.discoverSection = undefined;
			state.discoverSectionPage = 1;
			state.discoverSectionHasMore = true;
		}
		if (!state.authenticated) {
			return;
		}
		if (!state.isAdmin && forbiddenForUser(route)) {
			void goto(resolve('/discover'));
			return;
		}
		await loadRouteData(route);
	}

	async function loadRouteData(route: AppRouteState) {
		if (
			route.view === 'metadata-detail' ||
			route.view === 'media-people' ||
			route.view === 'related-section'
		) {
			await deps.loadMetadataDetail();
		} else if (route.view === 'media-collection') {
			await deps.loadMediaCollection();
		} else if (route.view === 'person-detail') {
			await deps.loadPersonDetail();
		} else if (route.view === 'discover-section') {
			await deps.loadDiscoverSection();
		}
	}

	return { applyRoute };
}

function forbiddenForUser(route: AppRouteState) {
	return route.view === 'settings' || route.view === 'system' || route.homeSection === 'blacklist';
}

function metadataKey(route: AppRouteState) {
	return `${route.metadataProvider ?? ''}:${route.metadataType ?? ''}:${route.metadataExternalId ?? ''}`;
}

function collectionKey(route: AppRouteState) {
	return `${route.collectionProvider ?? ''}:${route.collectionId ?? ''}`;
}

function personKey(route: AppRouteState) {
	return `${route.personProvider ?? ''}:${route.personId ?? ''}`;
}
