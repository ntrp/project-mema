import { beforeEach, describe, expect, it, vi } from 'vitest';

const apiMock = vi.hoisted(() => ({
	approveMediaRequest: vi.fn(),
	createMediaItem: vi.fn(),
	createMediaRequest: vi.fn(),
	deleteMediaItem: vi.fn(),
	deleteMediaItemFile: vi.fn(),
	enqueueMediaAutomaticSearch: vi.fn(),
	rescanMediaItemFiles: vi.fn(),
	updateMediaItem: vi.fn()
}));

const navigationMock = vi.hoisted(() => ({
	goto: vi.fn(),
	resolve: vi.fn((path: string) => path)
}));

vi.mock('$lib/settings/api', () => apiMock);
vi.mock('$app/navigation', () => ({ goto: navigationMock.goto }));
vi.mock('$app/paths', () => ({ resolve: navigationMock.resolve }));

import { createMediaActions } from '../mediaActions';
import type { AppShellState } from '../state.svelte';

function state(overrides: Record<string, unknown> = {}) {
	return {
		isAdmin: true,
		message: '',
		errorMessage: '',
		activeHomeSection: 'discover',
		activeMediaCandidate: undefined,
		addingKey: undefined,
		savingMediaAction: false,
		mediaItems: [mediaItem()],
		mediaRequests: [],
		releaseResults: { 'media-1': { loaded: true, releases: [], errors: [] } },
		activities: [{ id: 'activity-1', mediaItemId: 'media-1' }],
		selectedMediaItemId: undefined,
		...overrides
	} as unknown as AppShellState;
}

function deps() {
	return {
		clearNotice: vi.fn(),
		loadMediaItems: vi.fn(),
		loadSettings: vi.fn()
	};
}

function candidate(overrides: Record<string, unknown> = {}) {
	return {
		title: 'Scenario Movie',
		type: 'movie',
		year: 2026,
		externalProvider: 'tmdb',
		externalId: '123',
		overview: 'A scenario candidate',
		posterPath: '/poster.jpg',
		...overrides
	};
}

function mediaItem(overrides: Record<string, unknown> = {}) {
	return {
		id: 'media-1',
		title: 'Scenario Movie',
		type: 'movie',
		monitored: true,
		monitorMode: 'all',
		filePaths: ['/media/movie.mkv'],
		metadataFilePaths: [],
		...overrides
	};
}

describe('media actions (SCN-MEDIA-001)', () => {
	beforeEach(() => {
		for (const value of Object.values(apiMock)) value.mockReset();
		navigationMock.goto.mockReset();
		navigationMock.resolve.mockClear();
	});

	it('adds monitored media for admins and routes to the matching library section', async () => {
		const shell = state();
		const actionDeps = deps();
		const selected = candidate({ type: 'series' });
		const actions = createMediaActions(shell, actionDeps);

		actions.addMedia(selected as never);
		await actions.confirmMediaAction({
			monitorMode: 'all_episodes',
			seriesType: 'standard',
			minimumAvailability: 'released',
			startSearch: true,
			qualityProfileId: 'profile-1',
			libraryFolderId: 'folder-1',
			tags: ['tracked']
		});

		expect(apiMock.createMediaItem).toHaveBeenCalledWith(
			expect.objectContaining({ type: 'series', seriesType: 'standard', startSearch: true })
		);
		expect(actionDeps.loadMediaItems).toHaveBeenCalledOnce();
		expect(actionDeps.loadSettings).toHaveBeenCalledOnce();
		expect(shell.message).toBe('Media item added to monitored');
		expect(shell.activeHomeSection).toBe('series');
		expect(shell.activeMediaCandidate).toBeUndefined();
		expect(navigationMock.goto).toHaveBeenCalledWith('/series');
	});

	it('creates a user request instead of adding directly for non-admin users', async () => {
		const shell = state({ isAdmin: false, mediaRequests: [{ id: 'old-request' }] });
		const request = { id: 'request-1', title: 'Scenario Movie' };
		apiMock.createMediaRequest.mockResolvedValue(request);
		const actions = createMediaActions(shell, deps());

		actions.addMedia(candidate() as never);
		await actions.confirmMediaAction({ monitorMode: 'none', tags: ['later'] } as never);

		expect(apiMock.createMediaRequest).toHaveBeenCalledWith(
			expect.objectContaining({ title: 'Scenario Movie', monitorMode: 'none' })
		);
		expect(shell.mediaRequests).toEqual([request, { id: 'old-request' }]);
		expect(shell.message).toBe('Media request created');
		expect(shell.activeHomeSection).toBe('requests');
		expect(navigationMock.goto).toHaveBeenCalledWith('/requests');
	});

	it('updates visible media after automatic search, file rescan, and file deletion', async () => {
		const shell = state();
		const actions = createMediaActions(shell, deps());
		apiMock.enqueueMediaAutomaticSearch.mockResolvedValue({ message: 'Queued search', jobId: 42 });
		apiMock.rescanMediaItemFiles.mockResolvedValue(
			mediaItem({ filePaths: ['/media/movie.mkv', '/media/extra.mkv'] })
		);
		apiMock.deleteMediaItemFile.mockResolvedValue(mediaItem({ filePaths: [] }));

		await actions.autoSearchMedia(mediaItem() as never);
		await actions.rescanMediaFiles(mediaItem() as never);
		await actions.deleteMediaFile(mediaItem() as never, '/media/movie.mkv');

		expect(shell.message).toBe('Media file deleted');
		expect(apiMock.deleteMediaItemFile).toHaveBeenCalledWith('media-1', '/media/movie.mkv');
		expect(shell.mediaItems[0].filePaths).toEqual([]);
		expect(shell.searchingItemId).toBeUndefined();
		expect(shell.scanningMediaItemId).toBeUndefined();
	});

	it('rolls back optimistic option changes when saving fails', async () => {
		const original = mediaItem({ monitored: false, monitorMode: 'none' });
		const shell = state({ mediaItems: [original] });
		apiMock.updateMediaItem.mockRejectedValue(new Error('write failed'));

		await createMediaActions(shell, deps()).saveMediaItemOptions(original as never, {
			monitored: true,
			monitorMode: 'only_media'
		});

		expect(shell.mediaItems).toEqual([original]);
		expect(shell.message).toBe('');
		expect(shell.errorMessage).toBe('write failed');
		expect(shell.savingMediaItemOptionsId).toBeUndefined();
	});

	it('removes media and request state through destructive actions', async () => {
		const shell = state({
			selectedMediaItemId: 'media-1',
			mediaRequests: [{ id: 'request-1', status: 'pending' }]
		});
		const actions = createMediaActions(shell, deps());
		apiMock.approveMediaRequest.mockResolvedValue({
			request: { id: 'request-1', status: 'approved' },
			mediaItem: mediaItem({ id: 'media-2' })
		});

		actions.deleteMediaItem(mediaItem() as never);
		await actions.confirmMediaDelete(false);
		await actions.approveMediaRequest({ id: 'request-1' } as never, {
			qualityProfileId: 'p1',
			libraryFolderId: 'folder-1'
		});

		expect(shell.mediaItems.map((item) => item.id)).toEqual(['media-2']);
		expect(shell.releaseResults).toEqual({});
		expect(shell.activities).toEqual([]);
		expect(shell.selectedMediaItemId).toBeUndefined();
		expect(shell.mediaRequests[0]).toMatchObject({ id: 'request-1', status: 'approved' });
		expect(navigationMock.goto).toHaveBeenCalledWith('/movies');
	});
});
