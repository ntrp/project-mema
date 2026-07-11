import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import type { AppShellState } from './state.svelte';
import { appRouteKey, type AppRouteState } from './routeState';
import { loadAppRouteData, type RouteDataDeps } from './routeData';

interface RouteDeps {
	routeData: RouteDataDeps;
}

export function createRouteActions(state: AppShellState, deps: RouteDeps) {
	let currentRouteKey = appRouteKey(state.route);

	async function applyRoute(route: AppRouteState) {
		const nextRouteKey = appRouteKey(route);
		if (nextRouteKey === currentRouteKey) {
			return;
		}

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

		if (!state.authenticated) {
			return;
		}
		if (!state.isAdmin && forbiddenForUser(route)) {
			void goto(resolve('/discover'));
			return;
		}
		await loadAppRouteData(route, state.isAdmin, deps.routeData);
	}

	return { applyRoute };
}

function forbiddenForUser(route: AppRouteState) {
	return route.view === 'settings' || route.view === 'system' || route.homeSection === 'blacklist';
}
