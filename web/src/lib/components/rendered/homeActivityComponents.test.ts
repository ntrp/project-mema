import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import ActivityList from '$lib/components/app/activity/ActivityList.svelte';
import MediaItemList from '$lib/components/app/home/MediaItemList.svelte';
import WantedMediaTable from '$lib/components/app/home/WantedMediaTable.svelte';
import type { DownloadActivity, MediaItem, ReleaseBlocklistItem } from '$lib/settings/types';
import { renderWithTooltip } from './renderHelpers';

describe('rendered activity components (SCN-ACTIVITY-001)', () => {
	it('renders activity rows with parsed release details and management actions', () => {
		const { body } = renderWithTooltip(ActivityList, {
			activities: [
				downloadActivity({
					status: 'downloading',
					progressPercent: 42,
					releaseTitle: 'Scenario.Movie.2026.German.1080p.WEB-DL.Atmos-GROUP'
				}),
				downloadActivity({
					id: 'activity-2',
					status: 'failed',
					failureType: 'import',
					error: 'Import failed',
					clientFiles: [
						{ path: '/downloads/Broken.Movie.2026.mkv', sizeBytes: 1_073_741_824, complete: true }
					],
					releaseTitle: 'Broken.Movie.2026.720p.HDTV-GROUP'
				}),
				downloadActivity({
					id: 'activity-3',
					mediaTitle: 'Scenario Series',
					mediaType: 'serie',
					mediaYear: 2025,
					releaseTitle: 'Scenario.Series.S02E05.The.Return.1080p.WEB-DL-GROUP'
				})
			],
			loading: false,
			canManage: true,
			cancellingId: 'activity-1',
			deletingId: 'activity-2',
			onRefresh: vi.fn(),
			onCancel: vi.fn(),
			onRetry: vi.fn(),
			onDelete: vi.fn()
		});

		expect(body).toContain('Activity queue');
		expect(body).toContain('Scenario Movie (2026)');
		expect(body).toContain('aria-label="Movie media type"');
		expect(body).toContain('Scenario.Movie.2026.German.1080p.WEB-DL.Atmos-GROUP');
		expect(body).toContain('Scenario Series (2025) - 2x5 - The Return');
		expect(body).toContain('aria-label="Series media type"');
		expect(body).not.toContain('>Season<');
		expect(body).not.toContain('>Episode<');
		expect(body).not.toContain('>Languages<');
		expect(body).toContain('Protocol');
		expect(body).toContain('Client');
		expect(body).toContain('Status');
		expect(body).toContain('aria-label="downloading status"');
		expect(body).toContain('aria-label="failed status: Import failed: 1 client file(s)"');
		expect(body).toContain('torrent');
		expect(body).toContain('German');
		expect(body).toContain('1080p');
		expect(body).toContain('WEB-DL');
		expect(body).toContain('42%');
		expect(body).not.toContain('2026-07-03T00:00:00Z');
		expect(body).toContain('Retry Broken.Movie.2026.720p.HDTV-GROUP');
		expect(body).toContain('Manual import Broken.Movie.2026.720p.HDTV-GROUP');
		expect(body).toContain('Cancel Scenario.Movie.2026.German.1080p.WEB-DL.Atmos-GROUP');
		expect(body).toContain('Delete Broken.Movie.2026.720p.HDTV-GROUP');
	});

	it('renders empty activity state and refresh loading label', () => {
		const { body } = renderWithTooltip(ActivityList, {
			activities: [],
			loading: true,
			canManage: false,
			onRefresh: vi.fn(),
			onCancel: vi.fn(),
			onRetry: vi.fn(),
			onDelete: vi.fn()
		});

		expect(body).toContain('Refreshing');
		expect(body).toContain('No queued activity');
	});

	it('allows cancelled history downloads to be retried', () => {
		const { body } = renderWithTooltip(ActivityList, {
			section: 'history' as const,
			activities: [downloadActivity({ status: 'cancelled' })],
			loading: false,
			canManage: true,
			onRefresh: vi.fn(),
			onCancel: vi.fn(),
			onRetry: vi.fn(),
			onDelete: vi.fn()
		});

		expect(body).toContain('Activity history');
		expect(body).toContain('Retry Scenario.Movie.2026.German.1080p.WEB-DL.Atmos-GROUP');
	});

	it('renders blocklist rows with protocol, client, and management actions', () => {
		const { body } = renderWithTooltip(ActivityList, {
			section: 'blocklist' as const,
			activities: [],
			releaseBlocklist: [releaseBlocklistItem()],
			loading: false,
			canManage: true,
			deletingBlocklistId: 'block-1',
			clearingReleaseBlocklist: false,
			onRefresh: vi.fn(),
			onCancel: vi.fn(),
			onRetry: vi.fn(),
			onDelete: vi.fn(),
			onDeleteReleaseBlocklistItem: vi.fn(),
			onClearReleaseBlocklist: vi.fn()
		});

		expect(body).toContain('Release blocklist');
		expect(body).toContain('Protocol');
		expect(body).toContain('Usenet');
		expect(body).toContain('Scenario Indexer');
		expect(body).toContain('Scenario Client');
		expect(body).toContain('Clear all');
		expect(body).toContain('Remove Scenario.Movie.2026.1080p from blocklist');
		expect(body).not.toContain('download_client_rejected');
	});
});

describe('rendered home media components (SCN-MEDIA-003)', () => {
	it('renders library cards with poster URLs, routes, and status labels', () => {
		const { body } = render(MediaItemList, {
			props: {
				mediaType: 'serie',
				items: [
					mediaItem({
						id: 'series-1',
						title: 'Continuing Series',
						type: 'serie',
						status: 'downloaded',
						metadataStatus: 'continuing',
						posterPath: '/poster.jpg'
					}),
					mediaItem({ id: 'series-2', title: 'Missing Series', type: 'serie', status: 'missing' })
				]
			}
		});

		expect(body).toContain('Added series');
		expect(body).toContain('Continuing Series');
		expect(body).toContain('Missing Series');
		expect(body).toContain('https://image.tmdb.org/t/p/w500/poster.jpg');
		expect(body).toContain('Open Continuing Series details');
		expect(body).toContain('Downloaded available episodes');
		expect(body).toContain('Missing');
	});

	it('renders wanted table with monitor, profile, and search state', () => {
		const { body } = renderWithTooltip(WantedMediaTable, {
			items: [
				mediaItem({
					monitorMode: 'collection',
					qualityProfileName: 'Scenario Profile',
					minimumAvailability: 'released'
				})
			],
			languages: [],
			searchingItemId: 'media-1',
			canManage: true,
			onFindReleases: vi.fn(),
			onGrabRelease: vi.fn()
		});

		expect(body).toContain('Wanted');
		expect(body).toContain('Scenario Movie');
		expect(body).toContain('Entire collection');
		expect(body).toContain('Scenario Profile');
		expect(body).toContain('released');
		expect(body).toContain('Automatic search Scenario Movie');
		expect(body).toContain('Manual search Scenario Movie');
	});

	it('renders library skeletons while loading', () => {
		const { body } = render(MediaItemList, {
			props: { mediaType: 'movie', items: [], loading: true }
		});
		expect(body).toContain('aria-label="Loading movies"');
		expect(body).not.toContain('No movies added yet.');
	});

	it('renders empty library and wanted states', () => {
		const library = render(MediaItemList, {
			props: { mediaType: 'movie', items: [] }
		});
		expect(library.body).toContain('No movies added yet.');
		expect(library.body).toContain('Discover');

		const wanted = render(WantedMediaTable, {
			props: {
				items: [],
				languages: [],
				canManage: false,
				onFindReleases: vi.fn(),
				onGrabRelease: vi.fn()
			}
		});
		expect(wanted.body).toContain('No missing media.');
	});
});

function downloadActivity(overrides: Partial<DownloadActivity> = {}): DownloadActivity {
	return {
		id: 'activity-1',
		mediaItemId: 'media-1',
		mediaTitle: 'Scenario Movie',
		mediaType: 'movie',
		mediaYear: 2026,
		releaseTitle: 'Scenario.Movie.2026.German.1080p.WEB-DL.Atmos-GROUP',
		indexerName: 'Scenario Indexer',
		indexerProtocol: 'torrent',
		downloadClientName: 'Scenario Client',
		downloadUrl: 'https://example.test/download',
		status: 'downloading',
		progressPercent: 42,
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:01:00Z',
		...overrides
	} as DownloadActivity;
}

function releaseBlocklistItem(overrides: Partial<ReleaseBlocklistItem> = {}): ReleaseBlocklistItem {
	return {
		id: 'block-1',
		mediaItemId: 'media-1',
		mediaTitle: 'Scenario Movie',
		mediaType: 'movie',
		releaseTitle: 'Scenario.Movie.2026.1080p',
		indexerName: 'Scenario Indexer',
		indexerProtocol: 'usenet',
		downloadClientName: 'Scenario Client',
		reason: 'Download client rejected release',
		source: 'download_client_rejected',
		temporary: true,
		expiresAt: '2026-07-04T12:00:00Z',
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:01:00Z',
		...overrides
	} as ReleaseBlocklistItem;
}

function mediaItem(overrides: Partial<MediaItem> = {}): MediaItem {
	return {
		id: 'media-1',
		title: 'Scenario Movie',
		type: 'movie',
		year: 2026,
		status: 'missing',
		monitorMode: 'onlyMedia',
		monitored: true,
		minimumAvailability: 'released',
		filePaths: [],
		metadataFilePaths: [],
		overview: 'A rendered media item.',
		...overrides
	} as MediaItem;
}
