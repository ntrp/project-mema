import { describe, expect, it, vi } from 'vitest';

import HomeArea from '$lib/components/app/home/HomeArea.svelte';
import type { HomeAreaProps } from '$lib/components/app/home/homeAreaTypes';
import type {
	DiscoverBlacklistItem,
	DownloadActivity,
	HomeSection,
	MediaDiscoverSection,
	MediaItem,
	MediaRequest,
	MediaSearchResult
} from '$lib/settings/types';
import { renderWithTooltip } from '../renderHelpers';

describe('rendered home area sections (SCN-MEDIA-004)', () => {
	it('routes discover, blacklist, requests, library, wanted, and activity sections', () => {
		expect(renderHome('discover').body).toContain('Browse media from metadata providers');
		expect(renderHome('discover').body).toContain('Trending');

		const blacklist = renderHome('blacklist').body;
		expect(blacklist).toContain('Blacklist');
		expect(blacklist).toContain('Hidden Movie');

		const requests = renderHome('requests').body;
		expect(requests).toContain('Media requests');
		expect(requests).toContain('Requested Movie');

		const movies = renderHome('movies').body;
		expect(movies).toContain('Added movies');
		expect(movies).toContain('Library Movie');

		const series = renderHome('series').body;
		expect(series).toContain('Added series');
		expect(series).toContain('Library Series');

		const wanted = renderHome('wanted').body;
		expect(wanted).toContain('Wanted');
		expect(wanted).toContain('Library Movie');
		expect(wanted).not.toContain('Library Series');

		const activity = renderHome('activity').body;
		expect(activity).toContain('Activity queue');
		expect(activity).toContain('Scenario.Movie.2026.1080p.WEB-DL-GROUP');
	});
});

function renderHome(activeSection: HomeSection) {
	return renderWithTooltip(HomeArea, {
		...baseProps(),
		activeSection
	});
}

function baseProps(): HomeAreaProps {
	return {
		activeSection: 'discover',
		activitySection: 'queue',
		mediaItems: [
			mediaItem(),
			mediaItem({ id: 'series-1', title: 'Library Series', type: 'serie', status: 'downloaded' })
		],
		mediaRequests: [mediaRequest()],
		discoverSections: [discoverSection()],
		discoverBlacklist: [blacklistItem()],
		libraryFolders: [],
		languages: [],
		qualityProfiles: [],
		activities: [downloadActivity()],
		releaseBlocklist: [],
		clearingReleaseBlocklist: false,
		loadingDiscover: false,
		loadingBlacklist: false,
		loadingMediaItems: false,
		canManage: true,
		loadingActivity: false,
		onAddMedia: vi.fn(),
		onBlacklistMedia: vi.fn(),
		onRemoveBlacklistMedia: vi.fn(),
		onApproveMediaRequest: vi.fn(),
		onFindReleases: vi.fn(),
		onAutoSearchMedia: vi.fn(),
		onRefreshMediaMetadata: vi.fn(),
		onSaveMediaItemOptions: vi.fn(),
		onDeleteMediaFile: vi.fn(),
		onDeleteMedia: vi.fn(),
		onGrabRelease: vi.fn(),
		onRefreshActivity: vi.fn(),
		onRefreshReleaseBlocklist: vi.fn(),
		onCancelActivity: vi.fn(),
		onDeleteActivity: vi.fn(),
		onDeleteReleaseBlocklistItem: vi.fn(),
		onClearReleaseBlocklist: vi.fn()
	};
}

function discoverSection(): MediaDiscoverSection {
	return {
		id: 'trending',
		title: 'Trending',
		providerName: 'TMDb',
		mediaType: 'movie',
		results: [searchResult({ title: 'Discover Movie', externalId: 'discover-1' })]
	};
}

function mediaItem(overrides: Partial<MediaItem> = {}): MediaItem {
	return {
		id: 'movie-1',
		title: 'Library Movie',
		type: 'movie',
		year: 2026,
		status: 'missing',
		monitored: true,
		monitorMode: 'only_media',
		minimumAvailability: 'released',
		filePaths: [],
		metadataFilePaths: [],
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:01:00Z',
		...overrides
	} as MediaItem;
}

function mediaRequest(): MediaRequest {
	return {
		id: 'request-1',
		title: 'Requested Movie',
		type: 'movie',
		year: 2026,
		status: 'pending',
		requestedByUserId: 'user-1',
		requestedByUsername: 'user',
		monitorMode: 'only_media',
		minimumAvailability: 'released',
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:01:00Z'
	} as MediaRequest;
}

function blacklistItem(): DiscoverBlacklistItem {
	return {
		id: 'blacklist-1',
		title: 'Hidden Movie',
		type: 'movie',
		year: 2025,
		externalProvider: 'tmdb',
		externalId: 'hidden-1',
		createdAt: '2026-07-03T00:00:00Z'
	} as DiscoverBlacklistItem;
}

function searchResult(overrides: Partial<MediaSearchResult> = {}): MediaSearchResult {
	return {
		title: 'Scenario Result',
		type: 'movie',
		year: 2026,
		externalProvider: 'tmdb',
		externalId: 'result-1',
		overview: 'Scenario overview',
		...overrides
	} as MediaSearchResult;
}

function downloadActivity(): DownloadActivity {
	return {
		id: 'activity-1',
		mediaItemId: 'movie-1',
		mediaTitle: 'Library Movie',
		mediaType: 'movie',
		mediaYear: 2026,
		releaseTitle: 'Scenario.Movie.2026.1080p.WEB-DL-GROUP',
		indexerName: 'Scenario Indexer',
		downloadClientName: 'Scenario Client',
		downloadUrl: 'https://example.test/download',
		status: 'downloading',
		progressPercent: 50,
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:01:00Z'
	} as DownloadActivity;
}
