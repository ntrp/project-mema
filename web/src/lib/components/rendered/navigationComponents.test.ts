import { describe, expect, it, vi } from 'vitest';

import AppNavSearch from '$lib/components/app/navigation/AppNavSearch.svelte';
import SidebarMenu from '$lib/components/app/navigation/SidebarMenu.svelte';
import {
	basePrimaryItems,
	settingsPrimaryItem,
	systemPrimaryItem
} from '$lib/components/app/navigation/appNavigation';
import { renderWithSidebar } from './renderHelpers';

describe('rendered navigation components (SCN-MEDIA-003)', () => {
	it('renders primary navigation and the active settings submenu', () => {
		const { body } = renderWithSidebar(SidebarMenu, {
			title: 'Media Manager',
			items: [...basePrimaryItems, settingsPrimaryItem, systemPrimaryItem],
			active: 'settings',
			activeSubmenu: 'users',
			onSelect: vi.fn(),
			onSubmenuSelect: vi.fn()
		});

		expect(body).toContain('Media Manager');
		expect(body).toContain('Discover');
		expect(body).toContain('Library');
		expect(body).toContain('Settings');
		expect(body).toContain('Users');
		expect(body).toContain('href="/settings/users"');
		expect(body).toContain('aria-current="page"');
		expect(body).not.toContain('Status');
	});

	it('renders active system sections without expanding unrelated groups', () => {
		const { body } = renderWithSidebar(SidebarMenu, {
			title: 'Media Manager',
			items: [...basePrimaryItems, settingsPrimaryItem, systemPrimaryItem],
			active: 'system',
			activeSubmenu: 'jobs',
			onSelect: vi.fn(),
			onSubmenuSelect: vi.fn()
		});

		expect(body).toContain('System');
		expect(body).toContain('Jobs');
		expect(body).toContain('Events');
		expect(body).toContain('href="/system/jobs"');
		expect(body).not.toContain('Download clients');
	});

	it('renders the global search input with advanced search affordance hidden until active', () => {
		const { body } = renderWithSidebar(AppNavSearch, {
			searchQuery: 'scenario',
			groups: [],
			loading: false,
			onSearch: vi.fn(),
			onSelect: vi.fn(),
			onAdvancedSearch: vi.fn()
		});

		expect(body).toContain('Search Movies &amp; TV');
		expect(body).toContain('id="global-search"');
		expect(body).not.toContain('Search suggestions');
	});
});
