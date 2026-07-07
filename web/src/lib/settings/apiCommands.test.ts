import { beforeEach, describe, expect, it, vi } from 'vitest';

const clientMock = vi.hoisted(() => ({
	GET: vi.fn(),
	POST: vi.fn(),
	PUT: vi.fn(),
	DELETE: vi.fn()
}));

vi.mock('$lib/api/client', () => ({ client: clientMock }));

import {
	addDiscoverBlacklistItem,
	advancedSearchMedia,
	applyMediaRename,
	approveMediaRequest,
	clearSystemEvents,
	createMediaItem,
	createMediaRequest,
	deleteDiscoverBlacklistItem,
	deleteMediaItemFile,
	deleteSystemEvent,
	enqueueMediaAutomaticSearch,
	enqueueMediaReleaseSearch,
	getFileNamingSettings,
	getSystemEventSettings,
	getSystemLogFileSettings,
	getSystemLogLevel,
	getSystemStatus,
	grabMediaRelease,
	listDiscoverBlacklist,
	listMediaProfiles,
	listSystemEvents,
	listSystemLogFiles,
	loadMediaDiscoverSection,
	loadMediaDiscoverSections,
	logout,
	refreshMediaItemMetadata,
	rescanMediaItemFiles,
	searchMedia,
	testCustomFormatParsing,
	updateFileNamingSettings,
	updateMediaItem,
	updateSystemEventSettings,
	updateSystemLogFileSettings,
	updateSystemLogLevel
} from './api';
import type { ReleaseCandidate } from './types';

describe('UI API command helpers (SCN-SETTINGS-009)', () => {
	beforeEach(() => {
		clientMock.GET.mockReset().mockResolvedValue({ data: {} });
		clientMock.POST.mockReset().mockResolvedValue({ data: {} });
		clientMock.PUT.mockReset().mockResolvedValue({ data: {} });
		clientMock.DELETE.mockReset().mockResolvedValue({ data: {} });
	});

	it('maps system, discovery, metadata, and media commands to API calls', async () => {
		await expect(logout()).resolves.toBeUndefined();
		await expect(getSystemLogLevel()).resolves.toEqual({});
		await expect(updateSystemLogLevel('debug')).resolves.toEqual({});
		await expect(getSystemStatus()).resolves.toEqual({});
		await expect(getSystemLogFileSettings()).resolves.toEqual({});
		await expect(
			updateSystemLogFileSettings({ enabled: true, directory: '/logs', retentionDays: 7 })
		).resolves.toEqual({});
		await expect(listSystemLogFiles()).resolves.toEqual([]);
		await expect(listSystemEvents({ before: 'now', limit: 5 })).resolves.toEqual({});
		await expect(deleteSystemEvent('event-1')).resolves.toBeUndefined();
		await expect(clearSystemEvents()).resolves.toBeUndefined();
		await expect(getSystemEventSettings()).resolves.toEqual({});
		await expect(updateSystemEventSettings({ retentionDays: 7 })).resolves.toEqual({});
		await expect(getFileNamingSettings()).resolves.toEqual({});
		await expect(updateFileNamingSettings({ movieTemplate: '{Title}' } as never)).resolves.toEqual(
			{}
		);
		await expect(listMediaProfiles()).resolves.toEqual([]);
		await expect(testCustomFormatParsing('Movie.2026.mkv')).resolves.toEqual({});
		await expect(searchMedia({ query: 'Scenario', type: 'movie' })).resolves.toEqual([]);
		await expect(loadMediaDiscoverSections()).resolves.toEqual([]);
		await expect(loadMediaDiscoverSection('popular', 2, 10)).resolves.toEqual({});
		await expect(listDiscoverBlacklist()).resolves.toEqual([]);
		await expect(addDiscoverBlacklistItem({ externalProvider: 'tmdb' } as never)).resolves.toEqual(
			{}
		);
		await expect(deleteDiscoverBlacklistItem('blacklist-1')).resolves.toBeUndefined();
		await expect(advancedSearchMedia({ query: 'Scenario' } as never)).resolves.toEqual([]);
		await expect(createMediaItem({ title: 'Scenario Movie' } as never)).resolves.toEqual({});
		await expect(updateMediaItem('media-1', { title: 'Updated' } as never)).resolves.toEqual({});
		await expect(refreshMediaItemMetadata('media-1')).resolves.toEqual({});
		await expect(createMediaRequest({ title: 'Scenario Movie' } as never)).resolves.toEqual({});
		await expect(approveMediaRequest('request-1', {} as never)).resolves.toEqual({});
		await expect(rescanMediaItemFiles('media-1')).resolves.toEqual({});
		await expect(applyMediaRename('media-1', ['/library/old.mkv'])).resolves.toEqual({});
		expect(clientMock.POST).toHaveBeenLastCalledWith('/media/items/{id}/rename-apply', {
			params: { path: { id: 'media-1' } },
			body: { currentPaths: ['/library/old.mkv'] }
		});
		await expect(deleteMediaItemFile('media-1', 'Movie.mkv')).resolves.toEqual({});
		await expect(enqueueMediaReleaseSearch('media-1', 'custom')).resolves.toEqual({});
		await expect(enqueueMediaAutomaticSearch('media-1')).resolves.toEqual({});
		await expect(
			grabMediaRelease('media-1', { id: 'release-1' } as ReleaseCandidate)
		).resolves.toEqual({});
	});
});
