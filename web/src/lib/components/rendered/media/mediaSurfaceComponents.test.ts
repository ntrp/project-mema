import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import MediaCollectionArea from '$lib/components/app/media/MediaCollectionArea.svelte';
import MediaPeopleArea from '$lib/components/app/media/MediaPeopleArea.svelte';
import type {
	MediaCollection,
	MediaItem,
	MediaMetadataDetails,
	MediaSearchResult
} from '$lib/settings/types';

describe('rendered media collection area (SCN-MEDIA-008)', () => {
	it('renders collection states, provider routes, library status, and add actions', () => {
		expect(collectionRender({ loading: true }).body).toContain('Loading collection');
		expect(collectionRender({ collection: undefined }).body).toContain('Collection not available');

		const { body } = collectionRender({
			collection: mediaCollection(),
			mediaItems: [mediaItem({ title: 'Owned Entry', externalId: 'owned-1' })],
			addingKey: 'movie:tmdb:new-1:New Entry:2027'
		});

		expect(body).toContain('Scenario Collection');
		expect(body).toContain('Collection overview');
		expect(body).toContain('Collection media');
		expect(body).toContain('2 titles');
		expect(body).toContain('Owned Entry');
		expect(body).toContain('New Entry');
		expect(body).toContain('/media/tmdb/movie/owned-1');
		expect(body).toContain('https://image.tmdb.org/t/p/w342/new-poster.jpg');
		expect(body).toContain('In library');
		expect(body).toContain('Working');
	});
});

describe('rendered media people area (SCN-MEDIA-008)', () => {
	it('renders loading, missing, cast, and crew groups from metadata', () => {
		expect(render(MediaPeopleArea, { props: { loading: true } }).body).toContain('Loading cast');
		expect(render(MediaPeopleArea, { props: { loading: false } }).body).toContain(
			'Cast not available'
		);

		const cast = render(MediaPeopleArea, {
			props: { detail: metadataDetail(), kind: 'cast', loading: false }
		}).body;
		expect(cast).toContain('Scenario Movie');
		expect(cast).toContain('Cast');
		expect(cast).toContain('Lead Actor');
		expect(cast).toContain('Hero');
		expect(cast).toContain('https://image.tmdb.org/t/p/w185/profile.jpg');

		const crew = render(MediaPeopleArea, {
			props: { detail: metadataDetail(), kind: 'crew', loading: false }
		}).body;
		expect(crew).toContain('Director');
		expect(crew).toContain('Ada Director');
		expect(crew).toContain('Writer');
		expect(crew).toContain('Grace Writer');
	});
});

type CollectionRenderProps = {
	collection?: MediaCollection;
	mediaItems: MediaItem[];
	loading: boolean;
	addingKey?: string;
	actionLabel: string;
	onAdd: (_candidate: MediaSearchResult) => void;
};

function collectionRender(overrides: Partial<CollectionRenderProps> = {}) {
	return render(MediaCollectionArea, {
		props: {
			collection: mediaCollection(),
			mediaItems: [],
			loading: false,
			actionLabel: 'Add to library',
			onAdd: vi.fn(),
			...overrides
		}
	});
}

function mediaCollection(): MediaCollection {
	return {
		id: 'collection-1',
		name: 'Scenario Collection',
		provider: 'tmdb',
		overview: 'Collection overview',
		results: [
			searchResult({ title: 'Owned Entry', externalId: 'owned-1' }),
			searchResult({
				title: 'New Entry',
				externalId: 'new-1',
				year: 2027,
				posterPath: '/new-poster.jpg'
			})
		]
	};
}

function metadataDetail(): MediaMetadataDetails {
	return {
		title: 'Scenario Movie',
		type: 'movie',
		year: 2026,
		externalProvider: 'tmdb',
		externalId: 'scenario-1',
		posterPath: '/poster.jpg',
		runtimeMinutes: 120,
		genres: ['Drama'],
		facts: [
			{ label: 'Director', value: 'Ada Director' },
			{ label: 'Writer', value: 'Grace Writer' }
		],
		cast: [{ name: 'Lead Actor', role: 'Hero', profilePath: '/profile.jpg' }]
	} as MediaMetadataDetails;
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

function mediaItem(overrides: Partial<MediaItem> = {}): MediaItem {
	return {
		id: 'media-1',
		title: 'Scenario Movie',
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
