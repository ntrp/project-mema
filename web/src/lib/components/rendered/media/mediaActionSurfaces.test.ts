import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import MediaActionOptions from '$lib/components/app/media/actions/MediaActionOptions.svelte';
import MediaHero from '$lib/components/app/media/detail/MediaHero.svelte';
import type { LibraryFolder, MediaItem, QualityProfileOption } from '$lib/settings/types';

describe('rendered media action surfaces (SCN-MEDIA-004)', () => {
	it('renders admin and requester option controls for media actions', () => {
		const admin = render(MediaActionOptions, {
			props: {
				mediaType: 'series',
				isAdmin: true,
				libraryFolders: [libraryFolder()],
				qualityProfiles: [qualityProfile()],
				qualityProfileId: 'profile-1',
				libraryFolderId: 'folder-1',
				monitorMode: 'all_episodes',
				seriesType: 'standard',
				minimumAvailability: 'released',
				onMonitorModeChange: vi.fn(),
				onSeriesTypeChange: vi.fn()
			}
		}).body;

		expect(admin).toContain('Library folder');
		expect(admin).toContain('Quality profile');
		expect(admin).toContain('/library');
		expect(admin).toContain('HD Profile');
		expect(admin).toContain('Monitor');
		expect(admin).toContain('All episodes');
		expect(admin).toContain('Series type');
		expect(admin).toContain('Standard');
		expect(admin).toContain('Minimum availability');
		expect(admin).toContain('Released');

		const request = render(MediaActionOptions, {
			props: {
				mediaType: 'movie',
				isAdmin: false,
				libraryFolders: [libraryFolder()],
				qualityProfiles: [qualityProfile()],
				qualityProfileId: '',
				libraryFolderId: '',
				monitorMode: 'only_media',
				seriesType: 'standard',
				minimumAvailability: 'announced',
				onMonitorModeChange: vi.fn(),
				onSeriesTypeChange: vi.fn()
			}
		}).body;

		expect(request).not.toContain('Library folder');
		expect(request).not.toContain('Quality profile');
		expect(request).toContain('Monitor');
		expect(request).toContain('Only this media');
		expect(request).toContain('Minimum availability');
		expect(request).toContain('Announced');
	});

	it('renders library hero metadata and management actions', () => {
		const { body } = render(MediaHero, {
			props: {
				mediaType: 'movie',
				item: mediaItem({ tags: ['favorite'], posterPath: '/poster.jpg' }),
				qualityProfileLabel: 'HD Profile',
				canManage: true,
				searchingItemId: 'media-1',
				scanningMediaItemId: 'media-1',
				deletingMediaItemId: 'other',
				onFindReleases: vi.fn(),
				onRescanMediaFiles: vi.fn(),
				onDeleteMedia: vi.fn()
			}
		});

		expect(body).toContain('Scenario Movie');
		expect(body).toContain('Movie');
		expect(body).toContain('Year');
		expect(body).toContain('2026');
		expect(body).toContain('Downloaded');
		expect(body).toContain('HD Profile');
		expect(body).toContain('Monitored');
		expect(body).toContain('favorite');
		expect(body).toContain('https://image.tmdb.org/t/p/w500/poster.jpg');
		expect(body).toContain('Queued');
		expect(body).toContain('Scanning');
		expect(body).toContain('Remove');
	});
});

function mediaItem(overrides: Partial<MediaItem> = {}): MediaItem {
	return {
		id: 'media-1',
		title: 'Scenario Movie',
		type: 'movie',
		year: 2026,
		status: 'downloaded',
		monitored: true,
		monitorMode: 'only_media',
		minimumAvailability: 'released',
		qualityProfileId: 'profile-1',
		mediaFolderPath: '/library/Scenario Movie',
		filePaths: ['/library/Scenario Movie/Scenario.Movie.2026.1080p.mkv'],
		metadataFilePaths: ['/library/Scenario Movie/movie.nfo'],
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:01:00Z',
		...overrides
	} as MediaItem;
}

function libraryFolder(): LibraryFolder {
	return { id: 'folder-1', path: '/library', createdAt: '', updatedAt: '' };
}

function qualityProfile(): QualityProfileOption {
	return { id: 'profile-1', name: 'HD Profile' };
}
