import type { AppRouteState } from './routeState';
import type { SettingsSection, SystemSection } from '$lib/settings/types';

const settingsSections = new Set<SettingsSection>([
	'general',
	'library',
	'download-clients',
	'indexers',
	'quality',
	'profiles',
	'custom-formats',
	'metadata',
	'subtitles',
	'dlna',
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

export function settingsRouteSection(value?: string): SettingsSection {
	return settingsSections.has(value as SettingsSection) ? (value as SettingsSection) : 'general';
}

export function systemRouteSection(value?: string): SystemSection {
	return systemSections.has(value as SystemSection) ? (value as SystemSection) : 'status';
}
