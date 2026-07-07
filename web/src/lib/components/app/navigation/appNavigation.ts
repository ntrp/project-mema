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
	| '/settings/subtitles'
	| '/settings/dlna'
	| '/settings/languages'
	| '/settings/tags'
	| '/settings/users';

export type SystemHref =
	| '/system/status'
	| '/system/indexing'
	| '/system/metadata'
	| '/system/jobs'
	| '/system/logs'
	| '/system/events';

export type HomeHref =
	| '/discover'
	| '/discover/movies'
	| '/discover/series'
	| `/discover/${string}`
	| '/blacklist'
	| '/requests'
	| '/movies'
	| '/series'
	| '/wanted'
	| '/activity'
	| '/activity/history'
	| '/activity/blocklist';

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
	{ value: 'subtitles', label: 'Subtitles', href: '/settings/subtitles' },
	{ value: 'dlna', label: 'DLNA', href: '/settings/dlna' },
	{ value: 'languages', label: 'Languages', href: '/settings/languages' },
	{ value: 'tags', label: 'Tags', href: '/settings/tags' },
	{ value: 'users', label: 'Users', href: '/settings/users' }
] satisfies PrimaryItem['children'];

export const systemItems = [
	{ value: 'status', label: 'Status', href: '/system/status' },
	{ value: 'indexing', label: 'Indexing', href: '/system/indexing' },
	{ value: 'metadata', label: 'Metadata', href: '/system/metadata' },
	{ value: 'jobs', label: 'Jobs', href: '/system/jobs' },
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
	{ value: 'movies', label: 'Movies', href: '/discover/movies' },
	{
		value: 'animated-movies',
		label: 'Animated Movies',
		href: '/discover/movies?genres=Animation&withoutKeywords=anime'
	},
	{
		value: 'anime-movies',
		label: 'Anime Movies',
		href: '/discover/movies?genres=Animation&keywords=anime'
	},
	{ value: 'series', label: 'Series', href: '/discover/series' },
	{
		value: 'animated-series',
		label: 'Animated Series',
		href: '/discover/series?genres=Animation&withoutKeywords=anime'
	},
	{
		value: 'anime-series',
		label: 'Anime Series',
		href: '/discover/series?genres=Animation&keywords=anime'
	}
] satisfies PrimaryItem['children'];

export const activityItems = [
	{ value: 'queue', label: 'Queue', href: '/activity' },
	{ value: 'history', label: 'History', href: '/activity/history' },
	{ value: 'blocklist', label: 'Blocklist', href: '/activity/blocklist' }
] satisfies PrimaryItem['children'];

export function discoverSectionHref(section: string): HomeHref {
	return discoverItems.find((item) => item.value === section)?.href ?? `/discover/${section}`;
}

export function activitySectionHref(section: string): HomeHref {
	return activityItems.find((item) => item.value === section)?.href ?? '/activity';
}

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
	{
		value: 'activity',
		label: 'Activity',
		icon: 'activity',
		href: '/activity',
		children: activityItems
	}
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
		case 'subtitles':
			return '/settings/subtitles';
		case 'dlna':
			return '/settings/dlna';
		case 'languages':
			return '/settings/languages';
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
		case 'indexing':
			return '/system/indexing';
		case 'metadata':
			return '/system/metadata';
		case 'jobs':
			return '/system/jobs';
		case 'logs':
			return '/system/logs';
		case 'events':
			return '/system/events';
	}
}
