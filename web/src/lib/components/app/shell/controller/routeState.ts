import type { AppView, HomeSection, SettingsSection, SystemSection } from '$lib/settings/types';
import type { PeopleSectionKind, RelatedSectionKind } from './types';

interface QueryParams {
	get: (_name: string) => string | null;
}

export interface AppRouteState {
	view: AppView;
	homeSection: HomeSection;
	settingsSection: SettingsSection;
	systemSection: SystemSection;
	selectedMediaItemId?: string;
	selectedRequestId?: string;
	advancedQuery: string;
	metadataProvider?: string;
	metadataType?: string;
	metadataExternalId?: string;
	collectionProvider?: string;
	collectionId?: string;
	discoverSectionId?: string;
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
	'tags',
	'users'
]);
const systemSections = new Set<SystemSection>([
	'status',
	'indexing',
	'metadata',
	'logs',
	'events'
]);

export function defaultRouteState(): AppRouteState {
	return {
		view: 'home',
		homeSection: 'discover',
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
	if (path === '/discover') {
		return route;
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
		return { ...route, homeSection: 'activity' };
	}
	return route;
}

export function appRouteKey(route: AppRouteState) {
	return [
		route.view,
		route.homeSection,
		route.settingsSection,
		route.systemSection,
		route.selectedMediaItemId ?? '',
		route.selectedRequestId ?? '',
		route.advancedQuery,
		route.metadataProvider ?? '',
		route.metadataType ?? '',
		route.metadataExternalId ?? '',
		route.collectionProvider ?? '',
		route.collectionId ?? '',
		route.discoverSectionId ?? '',
		route.relatedSectionKind,
		route.peopleSectionKind
	].join('|');
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
