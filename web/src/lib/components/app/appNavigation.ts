import type { HomeSection, SettingsSection, SystemSection } from '$lib/settings/types';

export type SettingsHref =
	| '/settings/library'
	| '/settings/download-clients'
	| '/settings/indexers'
	| '/settings/quality'
	| '/settings/file-naming'
	| '/settings/profiles'
	| '/settings/custom-formats'
	| '/settings/metadata'
	| '/settings/tags'
	| '/settings/users';

export type SystemHref = '/system/status' | '/system/logs' | '/system/log-files' | '/system/events';

export type HomeHref = '/discover' | '/requests' | '/movies' | '/series' | '/wanted' | '/activity';

export type PrimaryItem = {
	value: HomeSection | 'library' | 'settings' | 'system';
	label: string;
	icon: 'discover' | 'movies' | 'series' | 'activity' | 'settings' | 'computer';
	href: HomeHref | '/settings/library' | SystemHref;
	children?: readonly {
		value: HomeSection | SettingsSection | SystemSection;
		label: string;
		href: HomeHref | SettingsHref | SystemHref;
	}[];
};

export const settingsItems = [
	{ value: 'library', label: 'Library', href: '/settings/library' },
	{ value: 'download-clients', label: 'Download clients', href: '/settings/download-clients' },
	{ value: 'indexers', label: 'Indexers', href: '/settings/indexers' },
	{ value: 'quality', label: 'Quality', href: '/settings/quality' },
	{ value: 'file-naming', label: 'File naming', href: '/settings/file-naming' },
	{ value: 'profiles', label: 'Profiles', href: '/settings/profiles' },
	{ value: 'custom-formats', label: 'Custom formats', href: '/settings/custom-formats' },
	{ value: 'metadata', label: 'Metadata', href: '/settings/metadata' },
	{ value: 'tags', label: 'Tags', href: '/settings/tags' },
	{ value: 'users', label: 'Users', href: '/settings/users' }
] satisfies PrimaryItem['children'];

export const systemItems = [
	{ value: 'status', label: 'Status', href: '/system/status' },
	{ value: 'logs', label: 'Live logs', href: '/system/logs' },
	{ value: 'log-files', label: 'Log files', href: '/system/log-files' },
	{ value: 'events', label: 'Events', href: '/system/events' }
] satisfies PrimaryItem['children'];

export const libraryItems = [
	{ value: 'movies', label: 'Movies', href: '/movies' },
	{ value: 'series', label: 'Series', href: '/series' },
	{ value: 'wanted', label: 'Wanted', href: '/wanted' }
] satisfies PrimaryItem['children'];

export const basePrimaryItems = [
	{ value: 'discover', label: 'Discover', icon: 'discover', href: '/discover' },
	{ value: 'requests', label: 'Requests', icon: 'activity', href: '/requests' },
	{ value: 'library', label: 'Library', icon: 'movies', href: '/movies', children: libraryItems },
	{ value: 'activity', label: 'Activity', icon: 'activity', href: '/activity' }
] satisfies PrimaryItem[];

export const settingsPrimaryItem = {
	value: 'settings',
	label: 'Settings',
	icon: 'settings',
	href: '/settings/library',
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
		case 'download-clients':
			return '/settings/download-clients';
		case 'indexers':
			return '/settings/indexers';
		case 'quality':
			return '/settings/quality';
		case 'file-naming':
			return '/settings/file-naming';
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
		case 'log-files':
			return '/system/log-files';
		case 'events':
			return '/system/events';
	}
}
