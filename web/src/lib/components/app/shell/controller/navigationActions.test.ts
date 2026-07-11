import { beforeEach, describe, expect, it, vi } from 'vitest';

const gotoMock = vi.hoisted(() => vi.fn());

vi.mock('$app/navigation', () => ({ goto: gotoMock }));
vi.mock('$app/paths', () => ({
	resolve: (path: string, params?: Record<string, string>) =>
		params ? path.replace('[sectionId]', params.sectionId) : path
}));

import { createNavigationActions } from './navigationActions';
import type { AppShellState } from './state.svelte';

describe('navigation actions (SCN-MEDIA-004)', () => {
	beforeEach(() => {
		gotoMock.mockReset();
	});

	it('guards admin-only sections and routes public library sections', () => {
		const state = testState({ isAdmin: false });
		const actions = createNavigationActions(state);

		actions.selectHomeSection('blacklist');
		expect(gotoMock).not.toHaveBeenCalled();
		expect(state.activeHomeSection).toBe('movies');

		actions.selectHomeSection('series');
		expect(state.activeView).toBe('home');
		expect(state.activeHomeSection).toBe('series');
		expect(gotoMock).toHaveBeenLastCalledWith('/series');

		actions.selectPrimarySection('settings');
		expect(state.activeView).toBe('home');
		expect(gotoMock).toHaveBeenCalledTimes(1);
	});

	it('routes admin settings and system sections', () => {
		const state = testState({ isAdmin: true });
		const actions = createNavigationActions(state);

		actions.selectPrimarySection('settings');
		expect(state.activeView).toBe('settings');
		expect(state.activeSettingsSection).toBe('general');
		expect(gotoMock).toHaveBeenLastCalledWith('/settings/general');

		actions.selectSettingsSection('download-clients');
		expect(state.activeSettingsSection).toBe('download-clients');
		expect(gotoMock).toHaveBeenLastCalledWith('/settings/download-clients');

		actions.selectPrimarySection('system');
		expect(state.activeView).toBe('system');
		expect(state.activeSystemSection).toBe('status');
		expect(gotoMock).toHaveBeenLastCalledWith('/system/status');

		actions.selectSystemSection('dlna');
		expect(state.activeSystemSection).toBe('dlna');
		expect(gotoMock).toHaveBeenLastCalledWith('/system/dlna');

		actions.selectSystemSection('jobs');
		expect(state.activeSystemSection).toBe('jobs');
		expect(gotoMock).toHaveBeenLastCalledWith('/system/jobs');
	});

	it('opens discover subsections through their query-owned route', () => {
		const state = testState({
			activeView: 'home',
			activePrimarySection: 'discover'
		});
		const actions = createNavigationActions(state);

		actions.selectSubmenuSection('trending');

		expect(state.activeView).toBe('discover-section');
		expect(state.activeHomeSection).toBe('discover');
		expect(state.activeDiscoverSectionId).toBe('trending');
		expect(gotoMock).toHaveBeenLastCalledWith('/discover/trending');
	});

	it('routes discover preset entries to filtered movie and series pages', () => {
		const state = testState({ activePrimarySection: 'discover' });
		const actions = createNavigationActions(state);

		actions.selectSubmenuSection('animated-movies');
		expect(state.activeView).toBe('discover-movies');
		expect(state.activeDiscoverSubmenuSection).toBe('animated-movies');
		expect(gotoMock).toHaveBeenLastCalledWith(
			'/discover/movies?genres=Animation&withoutKeywords=anime'
		);

		actions.selectSubmenuSection('anime-series');
		expect(state.activeView).toBe('discover-series');
		expect(state.activeDiscoverSubmenuSection).toBe('anime-series');
		expect(gotoMock).toHaveBeenLastCalledWith('/discover/series?genres=Animation&keywords=anime');
	});

	it('routes activity submenu sections', () => {
		const state = testState({ activePrimarySection: 'activity' });
		const actions = createNavigationActions(state);

		actions.selectSubmenuSection('history');
		expect(state.activeView).toBe('home');
		expect(state.activeHomeSection).toBe('activity');
		expect(state.activeActivitySection).toBe('history');
		expect(gotoMock).toHaveBeenLastCalledWith('/activity/history');
	});

	it('opens the current user profile route', () => {
		const state = testState();
		const actions = createNavigationActions(state);

		actions.showProfile();

		expect(state.activeView).toBe('profile');
		expect(gotoMock).toHaveBeenLastCalledWith('/profile');
	});
});

function testState(overrides: Partial<AppShellState> = {}): AppShellState {
	return {
		isAdmin: true,
		activeView: 'home',
		activePrimarySection: 'library',
		activeHomeSection: 'movies',
		activeActivitySection: 'queue',
		activeSettingsSection: 'general',
		activeSystemSection: 'status',
		...overrides
	} as AppShellState;
}
