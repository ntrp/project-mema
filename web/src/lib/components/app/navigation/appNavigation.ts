import type { HomeSection, SettingsSection, SystemSection } from '$lib/settings/types';

export type SettingsHref =
	| '/settings/general'
	| '/settings/library'
	| '/settings/download-clients'
	| '/settings/indexers'
	| '/settings/quality'
	| '/settings/profiles'
	| '/settings/custom-formats'
	| '/settings/metadata'
	| '/settings/tags'
	| '/settings/users';

export type SystemHref = '/system/status' | '/system/logs' | '/system/events';

export type HomeHref =
	| '/discover'
	| `/discover/${string}`
	| '/blacklist'
	| '/requests'
	| '/movies'
	| '/series'
	| '/wanted'
	| '/activity';

export type PrimaryItem = {
	value: HomeSection | 'library' | 'settings' | 'system';
	label: string;
	icon: 'discover' | 'movies' | 'series' | 'activity' | 'settings' | 'computer' | 'visibility_off';
	href: HomeHref | SettingsHref | SystemHref;
	children?: readonly {
		value: string;
		label: string;
		href: HomeHref | SettingsHref | SystemHref;
	}[];
};

export const settingsItems = [
	{ value: 'general', label: 'General', href: '/settings/general' },
	{ value: 'library', label: 'Library', href: '/settings/library' },
	{ value: 'download-clients', label: 'Download clients', href: '/settings/download-clients' },
	{ value: 'indexers', label: 'Indexers', href: '/settings/indexers' },
	{ value: 'quality', label: 'Quality', href: '/settings/quality' },
	{ value: 'profiles', label: 'Profiles', href: '/settings/profiles' },
	{ value: 'custom-formats', label: 'Custom formats', href: '/settings/custom-formats' },
	{ value: 'metadata', label: 'Metadata', href: '/settings/metadata' },
	{ value: 'tags', label: 'Tags', href: '/settings/tags' },
	{ value: 'users', label: 'Users', href: '/settings/users' }
] satisfies PrimaryItem['children'];

export const systemItems = [
	{ value: 'status', label: 'Status', href: '/system/status' },
	{ value: 'events', label: 'Events', href: '/system/events' },
	{ value: 'logs', label: 'Logs', href: '/system/logs' }
] satisfies PrimaryItem['children'];

export const libraryItems = [
	{ value: 'movies', label: 'Movies', href: '/movies' },
	{ value: 'series', label: 'Series', href: '/series' },
	{ value: 'wanted', label: 'Wanted', href: '/wanted' }
] satisfies PrimaryItem['children'];

export const discoverItems = [
	{ value: 'discover', label: 'Home', href: '/discover' },
	{ value: 'trending', label: 'Trending', href: '/discover/trending' },
	{ value: 'movie-popular', label: 'Popular movies', href: '/discover/movie-popular' },
	{ value: 'movie-upcoming', label: 'Upcoming movies', href: '/discover/movie-upcoming' },
	{ value: 'movie-top-rated', label: 'Top rated movies', href: '/discover/movie-top-rated' },
	{ value: 'series-popular', label: 'Popular series', href: '/discover/series-popular' },
	{ value: 'series-on-the-air', label: 'Airing series', href: '/discover/series-on-the-air' },
	{ value: 'series-top-rated', label: 'Top rated series', href: '/discover/series-top-rated' }
] satisfies PrimaryItem['children'];

export const basePrimaryItems = [
	{
		value: 'discover',
		label: 'Discover',
		icon: 'discover',
		href: '/discover',
		children: discoverItems
	},
	{ value: 'blacklist', label: 'Blacklist', icon: 'visibility_off', href: '/blacklist' },
	{ value: 'requests', label: 'Requests', icon: 'activity', href: '/requests' },
	{ value: 'library', label: 'Library', icon: 'movies', href: '/movies', children: libraryItems },
	{ value: 'activity', label: 'Activity', icon: 'activity', href: '/activity' }
] satisfies PrimaryItem[];

export const settingsPrimaryItem = {
	value: 'settings',
	label: 'Settings',
	icon: 'settings',
	href: '/settings/general',
	children: settingsItems
} satisfies PrimaryItem;

export const systemPrimaryItem = {
	value: 'system',
	label: 'System',
	icon: 'computer',
	href: '/system/status',
	children: systemItems
} satisfies PrimaryItem;

export function settingsSectionHref(section: SettingsSection): SettingsHref {
	switch (section) {
		case 'general':
			return '/settings/general';
		case 'download-clients':
			return '/settings/download-clients';
		case 'indexers':
			return '/settings/indexers';
		case 'quality':
			return '/settings/quality';
		case 'profiles':
			return '/settings/profiles';
		case 'custom-formats':
			return '/settings/custom-formats';
		case 'metadata':
			return '/settings/metadata';
		case 'tags':
			return '/settings/tags';
		case 'users':
			return '/settings/users';
		default:
			return '/settings/library';
	}
}

export function systemSectionHref(section: SystemSection): SystemHref {
	switch (section) {
		case 'status':
			return '/system/status';
		case 'logs':
			return '/system/logs';
		case 'events':
			return '/system/events';
	}
}
