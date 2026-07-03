import { describe, expect, it, vi } from 'vitest';

import DiscoverBlacklistArea from '$lib/components/app/discovery/DiscoverBlacklistArea.svelte';
import DiscoverSectionArea from '$lib/components/app/discovery/DiscoverSectionArea.svelte';
import MediaSearchPanel from '$lib/components/app/discovery/MediaSearchPanel.svelte';
import type {
	DiscoverBlacklistItem,
	MediaDiscoverSection,
	MediaItem,
	MediaSearchResult
} from '$lib/settings/types';
import { renderWithTooltip } from './renderHelpers';

describe('rendered discovery components (SCN-MEDIA-004)', () => {
	it('renders filtered discovery sections and hides existing or blacklisted media', () => {
		const { body } = renderWithTooltip(MediaSearchPanel, {
			sections: [discoverSection()],
			mediaItems: [mediaItem({ title: 'Already In Library', externalId: 'lib-1' })],
			blacklist: [blacklistItem({ title: 'Hidden Movie', externalId: 'hidden-1' })],
			loading: false,
			actionLabel: 'Add',
			canManage: true,
			onAdd: vi.fn(),
			onBlacklist: vi.fn()
		});

		expect(body).toContain('Browse media from metadata providers');
		expect(body).toContain('Trending');
		expect(body).toContain('Visible Movie');
		expect(body).toContain('View all');
		expect(body).not.toContain('Already In Library');
		expect(body).not.toContain('Hidden Movie');
	});

	it('renders discovery loading and empty states', () => {
		const loading = renderWithTooltip(MediaSearchPanel, {
			sections: [],
			mediaItems: [],
			blacklist: [],
			loading: true,
			actionLabel: 'Add',
			canManage: false,
			onAdd: vi.fn(),
			onBlacklist: vi.fn()
		});
		expect(loading.body).toContain('Loading discovery');

		const empty = renderWithTooltip(MediaSearchPanel, {
			sections: [],
			mediaItems: [],
			blacklist: [],
			loading: false,
			actionLabel: 'Add',
			canManage: false,
			onAdd: vi.fn(),
			onBlacklist: vi.fn()
		});
		expect(empty.body).toContain('No discovery sections available');
	});

	it('renders a discover section page with library and blacklist filtering', () => {
		const { body } = renderWithTooltip(DiscoverSectionArea, {
			section: discoverSection(),
			mediaItems: [mediaItem({ title: 'Already In Library', externalId: 'lib-1' })],
			blacklist: [blacklistItem({ title: 'Hidden Movie', externalId: 'hidden-1' })],
			loading: false,
			loadingMore: true,
			actionLabel: 'Add',
			canManage: true,
			onAdd: vi.fn(),
			onBlacklist: vi.fn(),
			onLoadMore: vi.fn()
		});

		expect(body).toContain('Trending');
		expect(body).toContain('Visible Movie');
		expect(body).toContain('Loading more');
		expect(body).toContain('Already In Library');
		expect(body).not.toContain('Hidden Movie');
	});

	it('renders blacklist cards, placeholders, and empty state', () => {
		const { body } = renderWithTooltip(DiscoverBlacklistArea, {
			items: [
				blacklistItem({ posterPath: '/poster.jpg' }),
				blacklistItem({ id: 'blacklist-2', title: 'No Poster', posterPath: undefined })
			],
			loading: false,
			removingId: 'blacklist-1',
			onRemove: vi.fn()
		});

		expect(body).toContain('Blacklist');
		expect(body).toContain('2 hidden titles');
		expect(body).toContain('Visible Movie');
		expect(body).toContain('No Poster');
		expect(body).toContain('https://image.tmdb.org/t/p/w342/poster.jpg');
		expect(body).toContain('Remove Visible Movie from blacklist');
		expect(body).toContain('Remove No Poster from blacklist');

		const empty = renderWithTooltip(DiscoverBlacklistArea, {
			items: [],
			loading: false,
			onRemove: vi.fn()
		});
		expect(empty.body).toContain('No blacklisted media');
	});
});

function discoverSection(): MediaDiscoverSection {
	return {
		id: 'trending',
		title: 'Trending',
		mediaType: 'mixed',
		providerName: 'TMDb',
		results: [
			searchResult({ title: 'Visible Movie', externalId: 'visible-1' }),
			searchResult({ title: 'Already In Library', externalId: 'lib-1' }),
			searchResult({ title: 'Hidden Movie', externalId: 'hidden-1' })
		]
	} as MediaDiscoverSection;
}

function searchResult(overrides: Partial<MediaSearchResult> = {}): MediaSearchResult {
	return {
		title: 'Visible Movie',
		type: 'movie',
		year: 2026,
		externalProvider: 'tmdb',
		externalId: 'visible-1',
		overview: 'A visible discovery result.',
		posterPath: '/poster.jpg',
		popularity: 42,
		...overrides
	} as MediaSearchResult;
}

function mediaItem(overrides: Partial<MediaItem> = {}): MediaItem {
	return {
		id: 'media-1',
		title: 'Already In Library',
		type: 'movie',
		year: 2026,
		externalProvider: 'tmdb',
		externalId: 'lib-1',
		filePaths: [],
		metadataFilePaths: [],
		...overrides
	} as MediaItem;
}

function blacklistItem(overrides: Partial<DiscoverBlacklistItem> = {}): DiscoverBlacklistItem {
	return {
		id: 'blacklist-1',
		title: 'Visible Movie',
		type: 'movie',
		year: 2026,
		externalProvider: 'tmdb',
		externalId: 'visible-1',
		overview: 'A hidden discovery result.',
		...overrides
	} as DiscoverBlacklistItem;
}
