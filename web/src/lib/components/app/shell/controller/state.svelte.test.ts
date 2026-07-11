import { describe, expect, it } from 'vitest';

import { AppShellState } from './state.svelte';
import { routeStateFromPath } from './routeState';

const noQuery = { get: () => null };

describe('app shell state (SCN-MEDIA-004)', () => {
	it('initializes active sections from route state', () => {
		const state = new AppShellState(routeStateFromPath('/settings/quality', {}, noQuery));

		expect(state.activeView).toBe('settings');
		expect(state.activeSettingsSection).toBe('quality');
		expect(state.activePrimarySection).toBe('settings');
		expect(state.activeSubmenuSection).toBe('quality');
	});

	it('hides admin-only primary items for regular users', () => {
		const state = new AppShellState();

		state.currentUser = { id: 'user-1', username: 'viewer', role: 'user' };
		expect(state.isAdmin).toBe(false);
		expect(state.primaryItems.map((item) => item.value)).not.toContain('settings');
		expect(state.primaryItems.map((item) => item.value)).not.toContain('blacklist');

		state.currentUser = { id: 'admin-1', username: 'admin', role: 'admin' };
		expect(state.isAdmin).toBe(true);
		expect(state.primaryItems.map((item) => item.value)).toContain('settings');
		expect(state.primaryItems.map((item) => item.value)).toContain('system');
	});
});
