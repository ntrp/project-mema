import { beforeEach, describe, expect, it, vi } from 'vitest';

const apiMock = vi.hoisted(() => ({
	currentSession: vi.fn(),
	login: vi.fn(),
	logout: vi.fn()
}));

const navigationMock = vi.hoisted(() => ({
	goto: vi.fn(),
	resolve: vi.fn((path: string) => path)
}));

const eventsMock = vi.hoisted(() => ({
	connectAppEvents: vi.fn(),
	disconnectAppEvents: vi.fn()
}));

vi.mock('$lib/app/session/api', () => apiMock);
vi.mock('$app/navigation', () => ({ goto: navigationMock.goto }));
vi.mock('$app/paths', () => ({ resolve: navigationMock.resolve }));
vi.mock('../eventConnection', () => eventsMock);

import { createSessionActions } from '../sessionActions';
import type { AppShellState } from '../state.svelte';

function submitEvent() {
	return { preventDefault: vi.fn() } as unknown as SubmitEvent;
}

function shellState(overrides: Record<string, unknown> = {}) {
	return {
		loading: false,
		errorMessage: '',
		authenticated: false,
		currentUser: undefined,
		activeView: 'home',
		activeHomeSection: 'discover',
		route: {
			view: 'home',
			homeSection: 'discover',
			activitySection: 'queue',
			settingsSection: 'general',
			systemSection: 'status',
			advancedQuery: '',
			relatedSectionKind: 'recommendations',
			peopleSectionKind: 'cast'
		},
		downloadClients: [{ id: 'download-1' }],
		indexers: [{ id: 'indexer-1' }],
		metadataProviders: [{ id: 'metadata-1' }],
		mediaProfiles: [{ id: 'profile-1' }],
		customFormats: [{ id: 'format-1' }],
		users: [{ id: 'user-1' }],
		tags: [{ id: 'tag-1' }],
		languages: [{ code: 'EN' }],
		mediaItems: [{ id: 'media-1' }],
		mediaRequests: [{ id: 'request-1' }],
		activities: [{ id: 'activity-1' }],
		libraryFolders: [{ id: 'folder-1' }],
		pathMappings: [{ id: 'mapping-1' }],
		libraryScansByFolder: { 'folder-1': {} },
		openLibraryFolderId: 'folder-1',
		downloadForm: { name: 'Client' },
		indexerForm: { name: 'Indexer' },
		libraryFolderForm: { path: '/media', kind: 'movie' },
		pathMappingForm: { hostPath: '/downloads' },
		mediaProfileForm: { name: 'Profile' },
		customFormatForm: { name: 'Format' },
		tagForm: { name: 'Tag' },
		userForm: { username: 'User' },
		username: 'admin',
		password: 'password',
		...overrides
	} as unknown as AppShellState;
}

function deps() {
	const routeData = {
		loadSettingsSection: vi.fn(),
		loadSystemSettings: vi.fn(),
		loadMediaActionSettings: vi.fn(),
		loadMediaItems: vi.fn(),
		loadMediaRequests: vi.fn(),
		loadDownloadActivity: vi.fn(),
		loadReleaseBlocklist: vi.fn(),
		loadMetadataDetail: vi.fn(),
		loadPersonDetail: vi.fn(),
		loadMediaCollection: vi.fn(),
		loadProfile: vi.fn()
	};
	return {
		clearNotice: vi.fn(),
		clearActivityCache: vi.fn(),
		clearLibraryCache: vi.fn(),
		clearReleaseCache: vi.fn(),
		clearDiscoverBlacklistCache: vi.fn(),
		clearDiscoverContentCache: vi.fn(),
		clearSearchCache: vi.fn(),
		clearSettingsCatalogCache: vi.fn(),
		clearServerResourceCache: vi.fn(),
		clearLibraryScanCache: vi.fn(),
		routeData,
		events: {
			connect: vi.fn(),
			disconnect: vi.fn()
		}
	} as unknown as Parameters<typeof createSessionActions>[1];
}

describe('session actions (SCN-AUTH-002)', () => {
	beforeEach(() => {
		apiMock.currentSession.mockReset();
		apiMock.login.mockReset();
		apiMock.logout.mockReset();
		navigationMock.goto.mockReset();
		navigationMock.resolve.mockClear();
		eventsMock.connectAppEvents.mockReset();
		eventsMock.disconnectAppEvents.mockReset();
	});

	it('initialises an authenticated admin and loads route data only', async () => {
		const state = shellState({
			activeView: 'media-collection',
			route: {
				view: 'media-collection',
				homeSection: 'discover',
				activitySection: 'queue',
				settingsSection: 'general',
				systemSection: 'status',
				advancedQuery: '',
				collectionProvider: 'tmdb',
				collectionId: 'collection-1',
				relatedSectionKind: 'recommendations',
				peopleSectionKind: 'cast'
			}
		});
		const actionDeps = deps();
		apiMock.currentSession.mockResolvedValue({
			authenticated: true,
			user: { id: 'admin-1', username: 'admin', role: 'admin' }
		});

		await createSessionActions(state, actionDeps).initialise();

		expect(state.loading).toBe(false);
		expect(state.authenticated).toBe(true);
		expect(actionDeps.routeData.loadSystemSettings).not.toHaveBeenCalled();
		expect(eventsMock.connectAppEvents).toHaveBeenCalledWith(state, actionDeps.events);
	});

	it('redirects regular users away from admin-only views on login', async () => {
		const state = shellState({ activeView: 'settings', activeHomeSection: 'settings' });
		const actionDeps = deps();
		apiMock.login.mockResolvedValue({
			user: { id: 'user-1', username: 'viewer', role: 'user' }
		});

		await createSessionActions(state, actionDeps).login(submitEvent());

		expect(actionDeps.clearNotice).toHaveBeenCalledOnce();
		expect(state.authenticated).toBe(true);
		expect(state.activeView).toBe('home');
		expect(state.activeHomeSection).toBe('discover');
		expect(navigationMock.goto).toHaveBeenCalledWith('/discover');
		expect(actionDeps.routeData.loadSystemSettings).not.toHaveBeenCalled();
		expect(eventsMock.connectAppEvents).toHaveBeenCalledOnce();
	});

	it('disconnects events and clears authenticated state on logout even when the API fails', async () => {
		const state = shellState({
			authenticated: true,
			currentUser: { id: 'admin-1', role: 'admin' }
		});
		const actionDeps = deps();
		apiMock.logout.mockRejectedValue(new Error('network down'));

		await createSessionActions(state, actionDeps).logout();

		expect(eventsMock.disconnectAppEvents).toHaveBeenCalledWith(state);
		expect(state.authenticated).toBe(false);
		expect(state.currentUser).toBeUndefined();
		expect(state.activeView).toBe('home');
		expect(state.openLibraryFolderId).toBeUndefined();
		expect(actionDeps.clearReleaseCache).toHaveBeenCalledOnce();
		expect(actionDeps.clearDiscoverBlacklistCache).toHaveBeenCalledOnce();
		expect(actionDeps.clearDiscoverContentCache).toHaveBeenCalledOnce();
		expect(actionDeps.clearSearchCache).toHaveBeenCalledOnce();
		expect(state.errorMessage).toBe('network down');
	});
});
