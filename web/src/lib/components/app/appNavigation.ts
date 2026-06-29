import type { HomeSection, SettingsSection } from '$lib/settings/types';

export type SettingsHref =
	| '/settings/library'
	| '/settings/download-clients'
	| '/settings/indexers'
	| '/settings/quality'
	| '/settings/file-naming'
	| '/settings/profiles'
	| '/settings/metadata'
	| '/settings/tags'
	| '/settings/users'
	| '/settings/system/logs';

export type PrimaryItem = {
	value: HomeSection | 'settings';
	label: string;
	icon: 'discover' | 'movies' | 'series' | 'activity' | 'settings';
	href: '/discover' | '/requests' | '/movies' | '/series' | '/activity' | '/settings/library';
	children?: readonly {
		value: SettingsSection;
		label: string;
		href: SettingsHref;
	}[];
};

export const settingsItems = [
	{ value: 'library', label: 'Library', href: '/settings/library' },
	{ value: 'download-clients', label: 'Download clients', href: '/settings/download-clients' },
	{ value: 'indexers', label: 'Indexers', href: '/settings/indexers' },
	{ value: 'quality', label: 'Quality', href: '/settings/quality' },
	{ value: 'file-naming', label: 'File naming', href: '/settings/file-naming' },
	{ value: 'profiles', label: 'Profiles', href: '/settings/profiles' },
	{ value: 'metadata', label: 'Metadata', href: '/settings/metadata' },
	{ value: 'tags', label: 'Tags', href: '/settings/tags' },
	{ value: 'users', label: 'Users', href: '/settings/users' },
	{ value: 'system-logs', label: 'System logs', href: '/settings/system/logs' }
] satisfies PrimaryItem['children'];

export const basePrimaryItems = [
	{ value: 'discover', label: 'Discover', icon: 'discover', href: '/discover' },
	{ value: 'requests', label: 'Requests', icon: 'activity', href: '/requests' },
	{ value: 'movies', label: 'Movies', icon: 'movies', href: '/movies' },
	{ value: 'series', label: 'Series', icon: 'series', href: '/series' },
	{ value: 'activity', label: 'Activity', icon: 'activity', href: '/activity' }
] satisfies PrimaryItem[];

export const settingsPrimaryItem = {
	value: 'settings',
	label: 'Settings',
	icon: 'settings',
	href: '/settings/library',
	children: settingsItems
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
		case 'metadata':
			return '/settings/metadata';
		case 'tags':
			return '/settings/tags';
		case 'users':
			return '/settings/users';
		case 'system-logs':
			return '/settings/system/logs';
		default:
			return '/settings/library';
	}
}
