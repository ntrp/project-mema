import { afterEach, describe, expect, it, vi } from 'vitest';

const apiMock = vi.hoisted(() => ({
	enqueueMediaReleaseSearch: vi.fn(),
	getMediaCollection: vi.fn(),
	getMediaMetadataDetails: vi.fn(),
	grabMediaRelease: vi.fn(),
	listMediaItems: vi.fn(),
	listMediaRequests: vi.fn(),
	searchMediaReleases: vi.fn()
}));

const navigationMock = vi.hoisted(() => ({
	goto: vi.fn(),
	resolve: vi.fn((path: string) => path)
}));

vi.mock('$lib/settings/api', () => apiMock);
vi.mock('$lib/features/releases/api', () => ({
	enqueueMediaReleaseSearch: apiMock.enqueueMediaReleaseSearch,
	grabMediaRelease: apiMock.grabMediaRelease
}));
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

import { createReleaseActions } from '../releaseActions';
import type { AppShellState } from '../state.svelte';

function state(overrides: Record<string, unknown> = {}) {
	return {
		message: '',
		errorMessage: '',
		route: {},
		mediaItems: [],
		mediaRequests: [],
		loadingMediaItems: false,
		...overrides
	} as unknown as AppShellState;
}

function mediaItem(overrides: Record<string, unknown> = {}) {
	return { id: 'media-1', title: 'Scenario Movie', type: 'movie', ...overrides };
}

describe('release actions (SCN-ACTIVITY-001)', () => {
	it('loads release results immediately when the queued search returns candidates', async () => {
		const shell = state();
		const setReleaseResult = vi.fn();
		stubImmediateWindowTimer();
		apiMock.enqueueMediaReleaseSearch.mockResolvedValue({ message: 'Search queued', jobId: 7 });
		apiMock.searchMediaReleases.mockResolvedValue({ releases: [{ id: 'release-1' }], errors: [] });
		const actions = createReleaseActions(shell, {
			clearNotice: vi.fn(),
			upsertActivity: vi.fn(),
			refreshActivity: vi.fn(),
			updateMediaStatusFromActivity: vi.fn(),
			setReleaseResult,
			loadReleaseResult: apiMock.searchMediaReleases
		});

		await actions.findReleases(mediaItem() as never, 'scenario');

		expect(apiMock.enqueueMediaReleaseSearch).toHaveBeenCalledWith('media-1', 'scenario');
		expect(setReleaseResult).toHaveBeenLastCalledWith('media-1', {
			loaded: true,
			releases: [{ id: 'release-1' }],
			errors: []
		});
		expect(shell.searchingItemId).toBeUndefined();
	});

	it('grabs a release, updates activity state, and schedules an activity refresh', async () => {
		const shell = state();
		const updateMediaStatusFromActivity = vi.fn();
		const upsertActivity = vi.fn();
		const refreshActivity = vi.fn();
		stubImmediateWindowTimer();
		apiMock.grabMediaRelease.mockResolvedValue({
			message: 'Download queued',
			jobId: 9,
			activity: { id: 'activity-1', mediaItemId: 'media-1' }
		});

		await createReleaseActions(shell, {
			clearNotice: vi.fn(),
			upsertActivity,
			refreshActivity,
			updateMediaStatusFromActivity,
			setReleaseResult: vi.fn(),
			loadReleaseResult: apiMock.searchMediaReleases
		}).grabRelease(mediaItem() as never, { id: 'release-1' } as never, true);

		expect(upsertActivity).toHaveBeenCalledWith({ id: 'activity-1', mediaItemId: 'media-1' });
		expect(updateMediaStatusFromActivity).toHaveBeenCalledWith({
			id: 'activity-1',
			mediaItemId: 'media-1'
		});
		expect(refreshActivity).toHaveBeenCalledOnce();
		expect(shell.message).toBe('Download queued (#9)');
		expect(navigationMock.goto).toHaveBeenCalledWith('/activity');
	});
});
