import { beforeEach, describe, expect, it, vi } from 'vitest';

const apiMock = vi.hoisted(() => ({
	addDiscoverBlacklistItem: vi.fn(),
	deleteDiscoverBlacklistItem: vi.fn(),
	listDiscoverBlacklist: vi.fn(),
	loadMediaDiscoverSection: vi.fn(),
	loadMediaDiscoverSections: vi.fn()
}));

const navigationMock = vi.hoisted(() => ({
	goto: vi.fn(),
	resolve: vi.fn((path: string, params?: Record<string, string>) =>
		params?.sectionId ? `/discover/${params.sectionId}` : path
	)
}));

vi.mock('$lib/settings/api', () => apiMock);
vi.mock('$app/navigation', () => ({ goto: navigationMock.goto }));
vi.mock('$app/paths', () => ({ resolve: navigationMock.resolve }));

import { createDiscoveryActions } from '../discoveryActions';
import { createNavigationActions } from '../navigationActions';
import { createRouteActions } from '../routeActions';
import type { AppShellState } from '../state.svelte';

function state(overrides: Record<string, unknown> = {}) {
	return {
		isAdmin: true,
		authenticated: true,
		activeView: 'home',
		activeHomeSection: 'discover',
		activeSettingsSection: 'general',
		activeSystemSection: 'status',
		activePrimarySection: 'discover',
		activeDiscoverSectionId: 'popular',
		discoverSections: [],
		discoverSection: undefined,
		discoverSectionPage: 1,
		discoverSectionHasMore: true,
		discoverBlacklist: [],
		loadingDiscover: false,
		loadingDiscoverSection: false,
		loadingMoreDiscoverSection: false,
		loadingBlacklist: false,
		message: '',
		errorMessage: '',
		route: { view: 'home', homeSection: 'discover' },
		searchQuery: 'old',
		...overrides
	} as unknown as AppShellState;
}

function result(overrides: Record<string, unknown> = {}) {
	return {
		title: 'Scenario Pick',
		type: 'movie',
		year: 2026,
		externalProvider: 'tmdb',
		externalId: '123',
		...overrides
	};
}

describe('discovery actions (SCN-MEDIA-003)', () => {
	beforeEach(() => {
		for (const value of Object.values(apiMock)) value.mockReset();
		navigationMock.goto.mockReset();
	});

	it('loads sections and appends only new paginated results', async () => {
		const shell = state({
			discoverSection: { id: 'popular', results: [result({ externalId: '1' })] }
		});
		apiMock.loadMediaDiscoverSections.mockResolvedValue([{ id: 'popular', title: 'Popular' }]);
		apiMock.loadMediaDiscoverSection
			.mockResolvedValueOnce({ id: 'popular', results: [result({ externalId: '1' })] })
			.mockResolvedValueOnce({
				id: 'popular',
				results: [result({ externalId: '1' }), result({ externalId: '2' })]
			});
		const actions = createDiscoveryActions(shell);

		await actions.loadDiscoverSections();
		await actions.loadDiscoverSection();
		await actions.loadMoreDiscoverSection();

		expect(shell.discoverSections).toEqual([{ id: 'popular', title: 'Popular' }]);
		expect(shell.discoverSection?.results?.map((item) => item.externalId)).toEqual(['1', '2']);
		expect(shell.discoverSectionPage).toBe(2);
		expect(shell.loadingDiscoverSection).toBe(false);
	});

	it('blacklists and removes discover results for admins only', async () => {
		const blocked = { id: 'blacklist-1', ...result() };
		const shell = state({
			discoverSections: [{ id: 'popular', results: [result()] }],
			discoverSection: { id: 'popular', results: [result()] }
		});
		apiMock.listDiscoverBlacklist.mockResolvedValue([blocked]);
		apiMock.addDiscoverBlacklistItem.mockResolvedValue(blocked);
		apiMock.deleteDiscoverBlacklistItem.mockResolvedValue(undefined);
		const actions = createDiscoveryActions(shell);

		await actions.loadDiscoverBlacklist();
		await actions.blacklistDiscoverMedia(result() as never);
		await actions.removeDiscoverBlacklistItem(blocked as never);

		expect(apiMock.addDiscoverBlacklistItem).toHaveBeenCalledWith(
			expect.objectContaining({ title: 'Scenario Pick', externalId: '123' })
		);
		expect(shell.discoverSection?.results).toEqual([]);
		expect(shell.discoverBlacklist).toEqual([]);
		expect(shell.message).toBe('Scenario Pick removed from blacklist');
		expect(shell.removingBlacklistId).toBeUndefined();
	});

	it('does not expose blacklist controls to non-admin users', async () => {
		const shell = state({ isAdmin: false });
		const actions = createDiscoveryActions(shell);

		await actions.loadDiscoverBlacklist();
		await actions.blacklistDiscoverMedia(result() as never);
		await actions.removeDiscoverBlacklistItem({ id: 'blacklist-1' } as never);

		expect(apiMock.listDiscoverBlacklist).not.toHaveBeenCalled();
		expect(apiMock.addDiscoverBlacklistItem).not.toHaveBeenCalled();
		expect(apiMock.deleteDiscoverBlacklistItem).not.toHaveBeenCalled();
	});
});

describe('navigation and route actions (SCN-AUTH-003)', () => {
	it('routes primary and submenu selections through the correct app sections', async () => {
		const shell = state({ activePrimarySection: 'discover' });
		const loadDiscoverSection = vi.fn();
		const actions = createNavigationActions(shell, { loadDiscoverSection });

		actions.selectPrimarySection('settings');
		actions.selectSystemSection('events');
		actions.selectSubmenuSection('trending');

		expect(shell.activeSettingsSection).toBe('general');
		expect(shell.activeSystemSection).toBe('events');
		expect(shell.activeDiscoverSectionId).toBe('trending');
		expect(loadDiscoverSection).toHaveBeenCalledOnce();
		expect(navigationMock.goto).toHaveBeenLastCalledWith('/discover/trending');
	});

	it('applies routes, clears stale detail data, and redirects forbidden user views', async () => {
		const shell = state({
			isAdmin: false,
			metadataDetail: { title: 'Old detail' },
			mediaCollection: { title: 'Old collection' },
			route: {
				view: 'metadata-detail',
				homeSection: 'discover',
				metadataProvider: 'tmdb',
				metadataType: 'movie',
				metadataExternalId: '1'
			}
		});
		const deps = {
			loadDiscoverSection: vi.fn(),
			loadMediaCollection: vi.fn(),
			loadMetadataDetail: vi.fn(),
			loadPersonDetail: vi.fn(),
			loadProfile: vi.fn()
		};
		const actions = createRouteActions(shell, deps);

		await actions.applyRoute({
			view: 'advanced-search',
			homeSection: 'search',
			advancedQuery: 'scenario'
		} as never);
		await actions.applyRoute({ view: 'settings', homeSection: 'settings' } as never);

		expect(shell.metadataDetail).toBeUndefined();
		expect(shell.searchQuery).toBe('scenario');
		expect(shell.activeView).toBe('settings');
		expect(navigationMock.goto).toHaveBeenCalledWith('/discover');
		expect(deps.loadMetadataDetail).not.toHaveBeenCalled();
	});
});
