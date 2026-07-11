import { beforeEach, describe, expect, it, vi } from 'vitest';

const apiMock = vi.hoisted(() => ({
	addDiscoverBlacklistItem: vi.fn(),
	deleteDiscoverBlacklistItem: vi.fn(),
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
vi.mock('$lib/features/discovery/blacklist/api', () => apiMock);
vi.mock('$app/navigation', () => ({ goto: navigationMock.goto }));
vi.mock('$app/paths', () => ({ resolve: navigationMock.resolve }));

import { createDiscoveryActions, type DiscoveryBlacklistCache } from '../discoveryActions';
import type { DiscoverBlacklistItem } from '$lib/settings/types';
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
		loadingDiscover: false,
		loadingDiscoverSection: false,
		loadingMoreDiscoverSection: false,
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

	it('delegates pagination to the feature cache', async () => {
		const shell = state();
		const content = contentCache();
		const actions = createDiscoveryActions(shell, blacklistCache(), content);
		await actions.loadMoreDiscoverSection();

		expect(content.loadMore).toHaveBeenCalledWith('popular');
	});

	it('blacklists and removes discover results for admins only', async () => {
		const blocked = { id: 'blacklist-1', ...result() };
		const shell = state({
			discoverSections: [{ id: 'popular', results: [result()] }],
			discoverSection: { id: 'popular', results: [result()] }
		});
		apiMock.addDiscoverBlacklistItem.mockResolvedValue(blocked);
		apiMock.deleteDiscoverBlacklistItem.mockResolvedValue(undefined);
		const blacklist = blacklistCache();
		const content = contentCache();
		const actions = createDiscoveryActions(shell, blacklist, content);

		await actions.blacklistDiscoverMedia(result() as never);
		await actions.removeDiscoverBlacklistItem(blocked as never);

		expect(apiMock.addDiscoverBlacklistItem).toHaveBeenCalledWith(
			expect.objectContaining({ title: 'Scenario Pick', externalId: '123' })
		);
		expect(content.mapSections).toHaveBeenCalledOnce();
		expect(content.mapSection).toHaveBeenCalledWith('popular', expect.any(Function));
		expect(content.refresh).toHaveBeenCalledWith('popular');
		expect(blacklist.items()).toEqual([]);
		expect(shell.message).toBe('Scenario Pick removed from blacklist');
		expect(shell.removingBlacklistId).toBeUndefined();
	});

	it('does not expose blacklist controls to non-admin users', async () => {
		const shell = state({ isAdmin: false });
		const actions = createDiscoveryActions(shell, blacklistCache(), contentCache());

		await actions.blacklistDiscoverMedia(result() as never);
		await actions.removeDiscoverBlacklistItem({ id: 'blacklist-1' } as never);

		expect(apiMock.addDiscoverBlacklistItem).not.toHaveBeenCalled();
		expect(apiMock.deleteDiscoverBlacklistItem).not.toHaveBeenCalled();
	});
});

describe('navigation and route actions (SCN-AUTH-003)', () => {
	it('routes primary and submenu selections through the correct app sections', async () => {
		const shell = state({ activePrimarySection: 'discover' });
		const actions = createNavigationActions(shell);

		actions.selectPrimarySection('settings');
		actions.selectSystemSection('events');
		actions.selectSubmenuSection('trending');

		expect(shell.activeSettingsSection).toBe('general');
		expect(shell.activeSystemSection).toBe('events');
		expect(shell.activeDiscoverSectionId).toBe('trending');
		expect(navigationMock.goto).toHaveBeenLastCalledWith('/discover/trending');
	});

	it('applies routes, clears stale detail data, and redirects forbidden user views', async () => {
		const shell = state({
			isAdmin: false,
			route: {
				view: 'metadata-detail',
				homeSection: 'discover',
				metadataProvider: 'tmdb',
				metadataType: 'movie',
				metadataExternalId: '1'
			}
		});
		const deps = {
			routeData: {
				loadSettingsSection: vi.fn(),
				loadSystemSettings: vi.fn(),
				loadMediaActionSettings: vi.fn(),
				loadDiscoverSections: vi.fn(),
				loadDiscoverSection: vi.fn(),
				loadProfile: vi.fn()
			}
		};
		const actions = createRouteActions(shell, deps);

		await actions.applyRoute({
			view: 'advanced-search',
			homeSection: 'search',
			advancedQuery: 'scenario'
		} as never);
		await actions.applyRoute({ view: 'settings', homeSection: 'settings' } as never);

		expect(shell.searchQuery).toBe('scenario');
		expect(shell.activeView).toBe('settings');
		expect(navigationMock.goto).toHaveBeenCalledWith('/discover');
		expect(deps.routeData.loadMediaActionSettings).not.toHaveBeenCalled();
	});
});

function blacklistCache() {
	let entries: DiscoverBlacklistItem[] = [];
	return {
		items: () => entries,
		upsert: (item: DiscoverBlacklistItem) => {
			entries = [item, ...entries.filter((entry) => entry.id !== item.id)];
		},
		remove: (id: string) => {
			entries = entries.filter((entry) => entry.id !== id);
		}
	} satisfies DiscoveryBlacklistCache;
}

function contentCache() {
	return {
		mapSections: vi.fn(),
		mapSection: vi.fn(),
		refresh: vi.fn(),
		loadMore: vi.fn()
	};
}
