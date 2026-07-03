import type { AppRouteState } from './routeState';

interface QueryParams {
	get: (_name: string) => string | null;
}

export function appRouteKey(route: AppRouteState) {
	return [
		route.view,
		route.homeSection,
		route.activitySection,
		route.settingsSection,
		route.systemSection,
		route.selectedMediaItemId ?? '',
		route.selectedRequestId ?? '',
		route.advancedQuery,
		route.metadataProvider ?? '',
		route.metadataType ?? '',
		route.metadataExternalId ?? '',
		route.personProvider ?? '',
		route.personId ?? '',
		route.collectionProvider ?? '',
		route.collectionId ?? '',
		route.discoverSectionId ?? '',
		route.discoverSubmenuSection ?? '',
		route.relatedSectionKind,
		route.peopleSectionKind
	].join('|');
}

export function discoverSubmenu(searchParams: QueryParams, kind: 'movies' | 'series') {
	if (searchParams.get('genres') !== 'Animation') return kind;
	if (searchParams.get('keywords') === 'anime') return `anime-${kind}`;
	if (searchParams.get('withoutKeywords') === 'anime') return `animated-${kind}`;
	return kind;
}
