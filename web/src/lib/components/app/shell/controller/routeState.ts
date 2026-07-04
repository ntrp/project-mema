import type {
	ActivitySection,
	AppView,
	HomeSection,
	SettingsSection,
	SystemSection
} from '$lib/settings/types';
import type { PeopleSectionKind, RelatedSectionKind } from './types';
import { discoverSubmenu } from './routeStateHelpers';

export { appRouteKey } from './routeStateHelpers';

interface QueryParams {
	get: (_name: string) => string | null;
}

export interface AppRouteState {
	view: AppView;
	homeSection: HomeSection;
	activitySection: ActivitySection;
	settingsSection: SettingsSection;
	systemSection: SystemSection;
	selectedMediaItemId?: string;
	selectedRequestId?: string;
	advancedQuery: string;
	metadataProvider?: string;
	metadataType?: string;
	metadataExternalId?: string;
	personProvider?: string;
	personId?: string;
	collectionProvider?: string;
	collectionId?: string;
	discoverSectionId?: string;
	discoverSubmenuSection?: string;
	relatedSectionKind: RelatedSectionKind;
	peopleSectionKind: PeopleSectionKind;
}

const settingsSections = new Set<SettingsSection>([
	'general',
	'library',
	'download-clients',
	'indexers',
	'quality',
	'profiles',
	'custom-formats',
	'metadata',
	'languages',
	'tags',
	'users'
]);
const systemSections = new Set<SystemSection>([
	'status',
	'indexing',
	'metadata',
	'jobs',
	'logs',
	'events'
]);

export function defaultRouteState(): AppRouteState {
	return {
		view: 'home',
		homeSection: 'discover',
		activitySection: 'queue',
		settingsSection: 'general',
		systemSection: 'status',
		advancedQuery: '',
		relatedSectionKind: 'recommendations',
		peopleSectionKind: 'cast'
	};
}

export function routeStateFromPath(
	pathname: string,
	params: Record<string, string>,
	searchParams: QueryParams
): AppRouteState {
	const path = normalise(pathname);
	const segments = path.split('/').filter(Boolean);
	const route = defaultRouteState();

	if (path === '/') {
		return route;
	}
	if (path === '/search/advanced') {
		return { ...route, view: 'advanced-search', advancedQuery: searchParams.get('q') ?? '' };
	}
	if (path === '/profile') {
		return { ...route, view: 'profile' };
	}
	if (path === '/discover') {
		return route;
	}
	if (path === '/discover/movies') {
		return {
			...route,
			view: 'discover-movies',
			homeSection: 'discover',
			discoverSubmenuSection: discoverSubmenu(searchParams, 'movies')
		};
	}
	if (path === '/discover/series') {
		return {
			...route,
			view: 'discover-series',
			homeSection: 'discover',
			discoverSubmenuSection: discoverSubmenu(searchParams, 'series')
		};
	}
	if (segments[0] === 'discover' && params.sectionId) {
		return {
			...route,
			view: 'discover-section',
			homeSection: 'discover',
			discoverSectionId: params.sectionId
		};
	}
	if (segments[0] === 'settings') {
		const section = settingsSections.has(segments[1] as SettingsSection)
			? (segments[1] as SettingsSection)
			: 'general';
		return { ...route, view: 'settings', settingsSection: section };
	}
	if (segments[0] === 'system') {
		const section = systemSections.has(segments[1] as SystemSection)
			? (segments[1] as SystemSection)
			: 'status';
		return { ...route, view: 'system', systemSection: section };
	}
	if (segments[0] === 'media' && segments[1] === 'collections') {
		return {
			...route,
			view: 'media-collection',
			collectionProvider: params.provider,
			collectionId: params.collectionId
		};
	}
	if (segments[0] === 'media') {
		return mediaMetadataRoute(route, params, segments[4]);
	}
	if (segments[0] === 'people' && params.provider && params.personId) {
		return {
			...route,
			view: 'person-detail',
			personProvider: params.provider,
			personId: params.personId
		};
	}
	if (segments[0] === 'movies' || segments[0] === 'series') {
		return libraryRoute(route, segments[0], params.id, segments[2]);
	}
	if (segments[0] === 'requests') {
		return { ...route, homeSection: 'requests', selectedRequestId: params.id };
	}
	if (segments[0] === 'blacklist') {
		return { ...route, homeSection: 'blacklist' };
	}
	if (segments[0] === 'wanted') {
		return { ...route, homeSection: 'wanted' };
	}
	if (segments[0] === 'activity') {
		return {
			...route,
			homeSection: 'activity',
			activitySection:
				segments[1] === 'history' || segments[1] === 'blocklist' ? segments[1] : 'queue'
		};
	}
	return route;
}

function mediaMetadataRoute(
	route: AppRouteState,
	params: Record<string, string>,
	child?: string
): AppRouteState {
	const base = {
		...route,
		metadataProvider: params.provider,
		metadataType: params.type,
		metadataExternalId: params.externalId
	};
	if (child === 'cast' || child === 'crew') {
		return { ...base, view: 'media-people', peopleSectionKind: child };
	}
	if (child === 'recommendations' || child === 'similar') {
		return { ...base, view: 'related-section', relatedSectionKind: child };
	}
	return { ...base, view: 'metadata-detail' };
}

function libraryRoute(
	route: AppRouteState,
	section: 'movies' | 'series',
	id?: string,
	child?: string
): AppRouteState {
	const base = { ...route, homeSection: section, selectedMediaItemId: id };
	if (child === 'cast' || child === 'crew') {
		return { ...base, view: 'media-people', peopleSectionKind: child };
	}
	return base;
}

function normalise(pathname: string) {
	const path = pathname.replace(/\/+$/, '');
	return path.length > 0 ? path : '/';
}
