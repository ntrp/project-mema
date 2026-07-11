import type { SettingsSection, SystemSection } from '$lib/settings/types';
import type { SettingsHref, SystemHref } from './appNavigation';

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
		case 'dlna':
			return '/system/dlna';
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
