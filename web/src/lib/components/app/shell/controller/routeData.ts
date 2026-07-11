import type { AppRouteState } from './routeState';
import type { SettingsSection, SystemSection } from '$lib/settings/types';

export interface RouteDataDeps {
	loadSettingsSection: (_section: SettingsSection) => Promise<void>;
	loadSystemSettings: (_section: SystemSection) => Promise<void>;
	loadMediaActionSettings: () => Promise<void>;
	loadProfile: () => Promise<void>;
}

export async function loadAppRouteData(
	route: AppRouteState,
	isAdmin: boolean,
	deps: RouteDataDeps
) {
	const tasks: Array<Promise<void>> = [];

	if (route.view === 'settings') {
		if (isAdmin) tasks.push(deps.loadSettingsSection(route.settingsSection));
		return run(tasks);
	}
	if (route.view === 'system') {
		if (isAdmin && (route.systemSection === 'indexing' || route.systemSection === 'metadata')) {
			tasks.push(deps.loadSystemSettings(route.systemSection));
		}
		return run(tasks);
	}
	if (route.view === 'profile') return run(tasks);
	if (
		route.view === 'metadata-detail' ||
		route.view === 'media-people' ||
		route.view === 'related-section'
	) {
		return run(tasks);
	}
	if (route.view === 'media-collection' || route.view === 'person-detail') {
		return run(tasks);
	}
	if (route.view === 'discover-section') {
		return run(tasks);
	}
	if (route.view === 'discover-movies' || route.view === 'discover-series') {
		return run(tasks);
	}

	return loadPrimaryRouteData(route, isAdmin, deps);
}

function loadPrimaryRouteData(route: AppRouteState, isAdmin: boolean, deps: RouteDataDeps) {
	const tasks: Array<Promise<void>> = [];
	if (route.homeSection === 'movies' || route.homeSection === 'series') {
		if (route.selectedMediaItemId && isAdmin) tasks.push(deps.loadMediaActionSettings());
	} else if (route.homeSection === 'requests') {
		if (isAdmin) tasks.push(deps.loadMediaActionSettings());
	}
	return run(tasks);
}

async function run(tasks: Array<Promise<void>>) {
	await Promise.all(tasks);
}
