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
	approveMediaRequest,
	cancelDownloadActivity,
	createMediaItem,
	createMediaRequest,
	deleteDiscoverBlacklistItem,
	deleteMediaItemFile,
	getSystemEventSettings,
	getSystemLogFileSettings,
	getSystemLogLevel,
	getSystemStatus,
	grabMediaRelease,
	listDiscoverBlacklist,
	listDownloadActivity,
	listMediaProfiles,
	listMediaRequests,
	listQualitySizeSettings,
	listSystemEvents,
	listSystemLogFiles,
	logout,
	manualImportDownloadActivity,
	refreshMediaItemMetadata,
	rescanMediaItemFiles,
	updateMediaItem,
	updateSystemEventSettings,
	updateSystemLogFileSettings,
	updateSystemLogLevel
} from '../api';
import {
	enqueueMediaComponentExtraction,
	getMediaComponentSource,
	listMediaComponentSources,
	releaseMediaComponentSource,
	retainMediaComponentSource
} from './mediaComponentSources';

describe('additional UI API command helpers (SCN-SETTINGS-009)', () => {
	beforeEach(() => {
		clientMock.GET.mockReset();
		clientMock.POST.mockReset();
		clientMock.PUT.mockReset();
		clientMock.DELETE.mockReset();
	});

	it('maps system reads, writes, list defaults, and command errors', async () => {
		clientMock.GET.mockResolvedValueOnce({ data: { level: 'debug' } })
			.mockResolvedValueOnce({ data: { status: 'ok' } })
			.mockResolvedValueOnce({ data: { enabled: true } })
			.mockResolvedValueOnce({ data: { files: [{ name: 'app.log' }] } })
			.mockResolvedValueOnce({ data: undefined })
			.mockResolvedValueOnce({ data: { enabled: true } });
		clientMock.PUT.mockResolvedValueOnce({ data: { level: 'info' } })
			.mockResolvedValueOnce({ data: { enabled: false } })
			.mockResolvedValueOnce({ data: { enabled: false } });
		clientMock.POST.mockResolvedValueOnce({ error: { message: 'logout failed' } });

		await expect(getSystemLogLevel()).resolves.toEqual({ level: 'debug' });
		await expect(updateSystemLogLevel('info' as never)).resolves.toEqual({ level: 'info' });
		await expect(getSystemStatus()).resolves.toEqual({ status: 'ok' });
		await expect(getSystemLogFileSettings()).resolves.toEqual({ enabled: true });
		await expect(updateSystemLogFileSettings({ enabled: false } as never)).resolves.toEqual({
			enabled: false
		});
		await expect(listSystemLogFiles()).resolves.toEqual([{ name: 'app.log' }]);
		await expect(listSystemEvents({ limit: 5 })).resolves.toEqual({ events: [], hasMore: false });
		await expect(getSystemEventSettings()).resolves.toEqual({ enabled: true });
		await expect(updateSystemEventSettings({ enabled: false } as never)).resolves.toEqual({
			enabled: false
		});
		await expect(logout()).rejects.toThrow('logout failed');
		expect(clientMock.GET).toHaveBeenCalledWith('/system/events', {
			params: { query: { limit: 5 } }
		});
	});

	it('maps media, request, release, activity, and discovery commands', async () => {
		clientMock.GET.mockResolvedValueOnce({ data: undefined })
			.mockResolvedValueOnce({ data: { profiles: [{ id: 'profile-1' }] } })
			.mockResolvedValueOnce({ data: { items: [{ id: 'blacklist-1' }] } })
			.mockResolvedValueOnce({ data: { requests: [{ id: 'request-1' }] } })
			.mockResolvedValueOnce({ data: { activities: [{ id: 'activity-1' }] } });
		clientMock.POST.mockResolvedValue({ data: { id: 'result-1' } });
		clientMock.PUT.mockResolvedValue({ data: { id: 'media-1' } });
		clientMock.DELETE.mockResolvedValue({ data: {} });

		await expect(listQualitySizeSettings()).rejects.toThrow(
			'Quality size settings were not returned'
		);
		await expect(listMediaProfiles()).resolves.toEqual([{ id: 'profile-1' }]);
		await expect(addDiscoverBlacklistItem({ title: 'Hidden' } as never)).resolves.toEqual({
			id: 'result-1'
		});
		await expect(listDiscoverBlacklist()).resolves.toEqual([{ id: 'blacklist-1' }]);
		await expect(deleteDiscoverBlacklistItem('blacklist-1')).resolves.toBeUndefined();
		await expect(advancedSearchMedia({ query: 'matrix' } as never)).resolves.toEqual([]);
		await expect(createMediaItem({ title: 'Movie' } as never)).resolves.toEqual({ id: 'result-1' });
		await expect(updateMediaItem('media-1', { monitored: true } as never)).resolves.toEqual({
			id: 'media-1'
		});
		await expect(refreshMediaItemMetadata('media-1')).resolves.toEqual({ id: 'result-1' });
		await expect(listMediaRequests()).resolves.toEqual([{ id: 'request-1' }]);
		await expect(createMediaRequest({ title: 'Movie' } as never)).resolves.toEqual({
			id: 'result-1'
		});
		await expect(approveMediaRequest('request-1', { monitored: true } as never)).resolves.toEqual({
			id: 'result-1'
		});
		await expect(rescanMediaItemFiles('media-1')).resolves.toEqual({ id: 'result-1' });
		await expect(deleteMediaItemFile('media-1', '/movie.mkv')).resolves.toEqual({ id: 'result-1' });
		await expect(grabMediaRelease('media-1', { id: 'release-1' } as never, true)).resolves.toEqual({
			id: 'result-1'
		});
		await expect(listDownloadActivity()).resolves.toEqual([{ id: 'activity-1' }]);
		await expect(cancelDownloadActivity('activity-1')).resolves.toEqual({ id: 'result-1' });
		await expect(
			manualImportDownloadActivity('activity-1', { sourcePath: '/x.mkv' } as never)
		).resolves.toEqual({
			id: 'result-1'
		});
	});

	it('maps media component source commands', async () => {
		clientMock.GET.mockResolvedValueOnce({
			data: { sources: [{ id: 'source-1' }] }
		}).mockResolvedValueOnce({ data: { id: 'source-1' } });
		clientMock.POST.mockResolvedValueOnce({ data: { id: 'source-1' } })
			.mockResolvedValueOnce({ data: { id: 'source-1' } })
			.mockResolvedValueOnce({
				data: { jobId: 42, message: 'queued', artifact: { id: 'artifact-1' } }
			});

		await expect(listMediaComponentSources('media-1')).resolves.toEqual({
			sources: [{ id: 'source-1' }]
		});
		await expect(
			retainMediaComponentSource('media-1', {
				sourceRole: 'baseVideo',
				sourceFilePath: '/library/Movie/Base.mkv'
			})
		).resolves.toEqual({ id: 'source-1' });
		await expect(getMediaComponentSource('media-1', 'source-1')).resolves.toEqual({
			id: 'source-1'
		});
		await expect(releaseMediaComponentSource('media-1', 'source-1')).resolves.toEqual({
			id: 'source-1'
		});
		await expect(
			enqueueMediaComponentExtraction('media-1', 'source-1', {
				streamId: 2,
				streamType: 'audio'
			})
		).resolves.toEqual({ jobId: 42, message: 'queued', artifact: { id: 'artifact-1' } });
		expect(clientMock.POST).toHaveBeenNthCalledWith(
			2,
			'/media/items/{id}/component-sources/{sourceId}/release',
			{
				params: { path: { id: 'media-1', sourceId: 'source-1' } }
			}
		);
		expect(clientMock.POST).toHaveBeenLastCalledWith(
			'/media/items/{id}/component-sources/{sourceId}/extractions',
			{
				params: { path: { id: 'media-1', sourceId: 'source-1' } },
				body: { streamId: 2, streamType: 'audio' }
			}
		);
	});
});
