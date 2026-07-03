import { describe, expect, it } from 'vitest';

import {
	basePrimaryItems,
	discoverItems,
	libraryItems,
	settingsItems,
	settingsSectionHref,
	systemItems,
	systemSectionHref
} from './appNavigation';

describe('app navigation helpers (SCN-MEDIA-003)', () => {
	it('exposes stable primary navigation groups', () => {
		expect(basePrimaryItems.map((item) => item.value)).toEqual([
			'discover',
			'blacklist',
			'requests',
			'library',
			'activity'
		]);
		expect(discoverItems[0]).toMatchObject({ value: 'discover', href: '/discover' });
		expect(libraryItems.map((item) => item.href)).toEqual(['/movies', '/series', '/wanted']);
		expect(settingsItems.at(-1)).toMatchObject({ value: 'users', href: '/settings/users' });
		expect(systemItems.map((item) => item.value)).toContain('events');
	});

	it('maps settings and system sections to routed pages', () => {
		expect(settingsSectionHref('general')).toBe('/settings/general');
		expect(settingsSectionHref('download-clients')).toBe('/settings/download-clients');
		expect(settingsSectionHref('custom-formats')).toBe('/settings/custom-formats');
		expect(settingsSectionHref('library')).toBe('/settings/library');
		expect(systemSectionHref('status')).toBe('/system/status');
		expect(systemSectionHref('events')).toBe('/system/events');
		expect(systemSectionHref('logs')).toBe('/system/logs');
	});
});
