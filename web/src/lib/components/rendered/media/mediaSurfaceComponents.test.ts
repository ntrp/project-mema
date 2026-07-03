import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import MediaCollectionArea from '$lib/components/app/media/collection/MediaCollectionArea.svelte';
import MediaPeopleArea from '$lib/components/app/media/people/MediaPeopleArea.svelte';
import PersonDetailArea from '$lib/components/app/media/person-detail/PersonDetailArea.svelte';
import type {
	MediaCollection,
	MediaItem,
	MediaMetadataDetails,
	MediaSearchResult,
	PersonDetails
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
		expect(cast).toContain('/people/tmdb/person-1');

		const crew = render(MediaPeopleArea, {
			props: { detail: metadataDetail(), kind: 'crew', loading: false }
		}).body;
		expect(crew).toContain('Director');
		expect(crew).toContain('Ada Director');
		expect(crew).toContain('Writer');
		expect(crew).toContain('Grace Writer');
		expect(crew).toContain('/people/tmdb/director-1');
		expect(crew).toContain('/people/tmdb/writer-1');
	});
});

describe('rendered person detail area (SCN-MEDIA-008)', () => {
	it('renders profile details, filters, appearances, and add actions', () => {
		expect(
			render(PersonDetailArea, { props: personRenderProps({ loading: true }) }).body
		).toContain('Loading person details');
		expect(
			render(PersonDetailArea, { props: personRenderProps({ person: undefined }) }).body
		).toContain('Person not available');

		const { body } = render(PersonDetailArea, {
			props: personRenderProps({
				mediaItems: [mediaItem({ title: 'Scenario Movie', externalId: 'movie-1' })],
				addingKey: 'series:tmdb:series-1'
			})
		});

		expect(body).toContain('Scenario Person');
		expect(body).toContain('Born August 27, 1994 in Austin, Texas, USA');
		expect(body).toContain('Also Known As:');
		expect(body).toContain('Scenario biography.');
		expect(body).not.toContain('Backdrops');
		expect(body).toContain('Appearances');
		expect(body).toContain('Movies');
		expect(body).toContain('Series');
		expect(body).toContain('Scenario Movie');
		expect(body).toContain('as Mason');
		expect(body).toContain('/media/tmdb/movie/movie-1');
		expect(body).toContain('In library');
		expect(body).toContain('Working');
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
		crew: [
			{
				name: 'Ada Director',
				role: 'Director',
				profilePath: '/director.jpg',
				externalProvider: 'tmdb',
				externalId: 'director-1'
			},
			{
				name: 'Grace Writer',
				role: 'Writer',
				profilePath: '/writer.jpg',
				externalProvider: 'tmdb',
				externalId: 'writer-1'
			}
		],
		cast: [
			{
				name: 'Lead Actor',
				role: 'Hero',
				profilePath: '/profile.jpg',
				externalProvider: 'tmdb',
				externalId: 'person-1'
			}
		]
	} as MediaMetadataDetails;
}

type PersonRenderProps = {
	person?: PersonDetails;
	mediaItems: MediaItem[];
	loading: boolean;
	addingKey?: string;
	actionLabel: string;
	onAdd: (_candidate: MediaSearchResult) => void;
};

function personRenderProps(overrides: Partial<PersonRenderProps> = {}): PersonRenderProps {
	return {
		person: personDetail(),
		mediaItems: [],
		loading: false,
		actionLabel: 'Add to library',
		onAdd: vi.fn(),
		...overrides
	};
}

function personDetail(): PersonDetails {
	return {
		id: 'person-1',
		provider: 'tmdb',
		name: 'Scenario Person',
		birthday: '1994-08-27',
		placeOfBirth: 'Austin, Texas, USA',
		profilePath: '/profile.jpg',
		alsoKnownAs: ['Scenario Alias'],
		biography: 'Scenario biography.',
		appearances: [
			{
				title: 'Scenario Movie',
				type: 'movie',
				year: 2014,
				externalProvider: 'tmdb',
				externalId: 'movie-1',
				posterPath: '/poster.jpg',
				backdropPath: '/backdrop.jpg',
				role: 'Mason'
			},
			{
				title: 'Scenario Series',
				type: 'series',
				year: 2026,
				externalProvider: 'tmdb',
				externalId: 'series-1',
				posterPath: '/series.jpg',
				role: 'Self'
			}
		]
	};
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
		externalProvider: 'tmdb',
		externalId: 'scenario-1',
		monitorMode: 'only_media',
		minimumAvailability: 'released',
		filePaths: [],
		metadataFilePaths: [],
		createdAt: '2026-07-03T00:00:00Z',
		updatedAt: '2026-07-03T00:01:00Z',
		...overrides
	} as MediaItem;
}
