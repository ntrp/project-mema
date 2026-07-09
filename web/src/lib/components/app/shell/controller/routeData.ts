import type { AppRouteState } from './routeState';

export interface RouteDataDeps {
	loadSettings: () => Promise<void>;
	loadDiscoverBlacklist: () => Promise<void>;
	loadDiscoverSections: () => Promise<void>;
	loadDiscoverSection: () => Promise<void>;
	loadMediaItems: () => Promise<void>;
	loadMediaRequests: () => Promise<void>;
	loadDownloadActivity: () => Promise<void>;
	loadReleaseBlocklist: () => Promise<void>;
	loadMetadataDetail: () => Promise<void>;
	loadPersonDetail: () => Promise<void>;
	loadMediaCollection: () => Promise<void>;
	loadProfile: () => Promise<void>;
}

export async function loadAppRouteData(
	route: AppRouteState,
	isAdmin: boolean,
	deps: RouteDataDeps
) {
	const tasks: Array<Promise<void>> = [];

	if (route.view === 'settings') {
		if (isAdmin) tasks.push(deps.loadSettings());
		return run(tasks);
	}
	if (route.view === 'system') {
		if (isAdmin && (route.systemSection === 'indexing' || route.systemSection === 'metadata')) {
			tasks.push(deps.loadSettings());
		}
		return run(tasks);
	}
	if (route.view === 'profile') {
		tasks.push(deps.loadProfile());
		return run(tasks);
	}
	if (
		route.view === 'metadata-detail' ||
		route.view === 'media-people' ||
		route.view === 'related-section'
	) {
		tasks.push(deps.loadMetadataDetail(), deps.loadMediaItems());
		if (isAdmin) tasks.push(deps.loadDiscoverBlacklist());
		return run(tasks);
	}
	if (route.view === 'media-collection') {
		tasks.push(deps.loadMediaCollection(), deps.loadMediaItems());
		return run(tasks);
	}
	if (route.view === 'person-detail') {
		tasks.push(deps.loadPersonDetail(), deps.loadMediaItems());
		return run(tasks);
	}
	if (route.view === 'discover-section') {
		tasks.push(deps.loadDiscoverSection(), deps.loadMediaItems());
		if (isAdmin) tasks.push(deps.loadDiscoverBlacklist());
		return run(tasks);
	}
	if (route.view === 'discover-movies' || route.view === 'discover-series') {
		tasks.push(deps.loadMediaItems());
		if (isAdmin) tasks.push(deps.loadDiscoverBlacklist());
		return run(tasks);
	}

	return loadPrimaryRouteData(route, isAdmin, deps);
}

function loadPrimaryRouteData(route: AppRouteState, isAdmin: boolean, deps: RouteDataDeps) {
	const tasks: Array<Promise<void>> = [];
	if (route.homeSection === 'discover') {
		tasks.push(deps.loadDiscoverSections(), deps.loadMediaItems());
		if (isAdmin) tasks.push(deps.loadDiscoverBlacklist());
	} else if (route.homeSection === 'movies' || route.homeSection === 'series') {
		tasks.push(deps.loadMediaItems());
		if (route.selectedMediaItemId && isAdmin)
			tasks.push(deps.loadSettings(), deps.loadDownloadActivity());
	} else if (route.homeSection === 'wanted') {
		tasks.push(deps.loadMediaItems());
	} else if (route.homeSection === 'requests') {
		tasks.push(deps.loadMediaRequests());
		if (isAdmin) tasks.push(deps.loadSettings());
	} else if (route.homeSection === 'blacklist') {
		if (isAdmin) tasks.push(deps.loadDiscoverBlacklist());
	} else if (route.homeSection === 'activity') {
		if (route.activitySection === 'blocklist') {
			tasks.push(deps.loadReleaseBlocklist());
		} else {
			tasks.push(deps.loadDownloadActivity());
		}
	}
	return run(tasks);
}

async function run(tasks: Array<Promise<void>>) {
	await Promise.all(tasks);
}
