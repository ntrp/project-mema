import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

const apiMock = vi.hoisted(() => ({
	enqueueMediaReleaseSearch: vi.fn(),
	getMediaCollection: vi.fn(),
	getMediaMetadataDetails: vi.fn(),
	grabMediaRelease: vi.fn(),
	listDownloadActivity: vi.fn(),
	listMediaItems: vi.fn(),
	listMediaRequests: vi.fn(),
	loadSettings: vi.fn(),
	searchMediaReleases: vi.fn()
}));

const navigationMock = vi.hoisted(() => ({
	goto: vi.fn(),
	resolve: vi.fn((path: string) => path)
}));

vi.mock('$lib/settings/api', () => apiMock);
vi.mock('$app/navigation', () => ({ goto: navigationMock.goto }));
vi.mock('$app/paths', () => ({ resolve: navigationMock.resolve }));

afterEach(() => {
	vi.restoreAllMocks();
	vi.unstubAllGlobals();
});

function stubImmediateWindowTimer() {
	vi.stubGlobal('window', {
		setTimeout(callback: () => void) {
			callback();
			return 1;
		}
	});
}

import { createLoadActions } from '../loaders';
import { createReleaseActions } from '../releaseActions';
import type { AppShellState } from '../state.svelte';

function state(overrides: Record<string, unknown> = {}) {
	return {
		message: '',
		errorMessage: '',
		route: {},
		mediaItems: [],
		mediaRequests: [],
		activities: [],
		releaseResults: {},
		loadingMediaItems: false,
		loadingActivity: false,
		loadingMetadataDetail: false,
		loadingMediaCollection: false,
		...overrides
	} as unknown as AppShellState;
}

function mediaItem(overrides: Record<string, unknown> = {}) {
	return { id: 'media-1', title: 'Scenario Movie', type: 'movie', ...overrides };
}

describe('load actions (SCN-SYSTEM-001)', () => {
	beforeEach(() => {
		for (const value of Object.values(apiMock)) value.mockReset();
		navigationMock.goto.mockReset();
	});

	it('loads settings, library lists, metadata detail, and collection state', async () => {
		const shell = state({
			route: {
				metadataProvider: 'tmdb',
				metadataType: 'movie',
				metadataExternalId: '123',
				collectionProvider: 'tmdb',
				collectionId: 'collection-1'
			}
		});
		apiMock.loadSettings.mockResolvedValue({
			downloadClients: [{ id: 'client-1' }],
			indexers: [{ id: 'indexer-1' }],
			indexerSearch: { stats: {} },
			metadataProviders: [{ id: 'metadata-1' }],
			metadataCache: { stats: {} },
			libraryFolders: [{ id: 'folder-1' }],
			pathMappings: [{ id: 'mapping-1' }],
			mediaProfiles: [{ id: 'profile-1' }],
			customFormats: [{ id: 'format-1' }],
			users: [{ id: 'user-1' }],
			tags: [{ id: 'tag-1' }],
			languages: [{ code: 'EN' }]
		});
		apiMock.listMediaItems.mockResolvedValue([mediaItem()]);
		apiMock.listMediaRequests.mockResolvedValue([{ id: 'request-1' }]);
		apiMock.listDownloadActivity.mockResolvedValue([{ id: 'activity-1' }]);
		apiMock.getMediaMetadataDetails.mockResolvedValue({ title: 'Scenario Movie' });
		apiMock.getMediaCollection.mockResolvedValue({ title: 'Scenario Collection' });

		const actions = createLoadActions(shell);
		await actions.loadSettings();
		await actions.loadLibrary();
		await actions.loadMetadataDetail();
		await actions.loadMediaCollection();

		expect(shell.downloadClients).toEqual([{ id: 'client-1' }]);
		expect(shell.mediaItems).toEqual([mediaItem()]);
		expect(shell.mediaRequests).toEqual([{ id: 'request-1' }]);
		expect(shell.activities).toEqual([{ id: 'activity-1' }]);
		expect(shell.metadataDetail).toEqual({ title: 'Scenario Movie' });
		expect(shell.mediaCollection).toEqual({ title: 'Scenario Collection' });
		expect(shell.loadingMediaCollection).toBe(false);
	});

	it('leaves detail loaders idle when the route has no target identifiers', async () => {
		const shell = state();
		const actions = createLoadActions(shell);

		await actions.loadMetadataDetail();
		await actions.loadMediaCollection();

		expect(apiMock.getMediaMetadataDetails).not.toHaveBeenCalled();
		expect(apiMock.getMediaCollection).not.toHaveBeenCalled();
		expect(shell.loadingMetadataDetail).toBe(false);
	});
});

describe('release actions (SCN-ACTIVITY-001)', () => {
	it('loads release results immediately when the queued search returns candidates', async () => {
		const shell = state();
		stubImmediateWindowTimer();
		apiMock.enqueueMediaReleaseSearch.mockResolvedValue({ message: 'Search queued', jobId: 7 });
		apiMock.searchMediaReleases.mockResolvedValue({ releases: [{ id: 'release-1' }], errors: [] });
		const actions = createReleaseActions(shell, {
			clearNotice: vi.fn(),
			loadDownloadActivity: vi.fn(),
			updateMediaStatusFromActivity: vi.fn()
		});

		await actions.findReleases(mediaItem() as never, 'scenario');

		expect(apiMock.enqueueMediaReleaseSearch).toHaveBeenCalledWith('media-1', 'scenario');
		expect(shell.releaseResults['media-1']).toEqual({
			loaded: true,
			releases: [{ id: 'release-1' }],
			errors: []
		});
		expect(shell.searchingItemId).toBeUndefined();
	});

	it('grabs a release, updates activity state, and schedules an activity refresh', async () => {
		const shell = state({ activities: [{ id: 'old-activity' }] });
		const updateMediaStatusFromActivity = vi.fn();
		const loadDownloadActivity = vi.fn();
		stubImmediateWindowTimer();
		apiMock.grabMediaRelease.mockResolvedValue({
			message: 'Download queued',
			jobId: 9,
			activity: { id: 'activity-1', mediaItemId: 'media-1' }
		});

		await createReleaseActions(shell, {
			clearNotice: vi.fn(),
			loadDownloadActivity,
			updateMediaStatusFromActivity
		}).grabRelease(mediaItem() as never, { id: 'release-1' } as never, true);

		expect(shell.activities.map((activity) => activity.id)).toEqual(['activity-1', 'old-activity']);
		expect(updateMediaStatusFromActivity).toHaveBeenCalledWith({
			id: 'activity-1',
			mediaItemId: 'media-1'
		});
		expect(loadDownloadActivity).toHaveBeenCalledOnce();
		expect(shell.message).toBe('Download queued (#9)');
		expect(navigationMock.goto).toHaveBeenCalledWith('/activity');
	});
});
