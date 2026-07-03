import { describe, expect, it } from 'vitest';

import { appRouteKey, defaultRouteState, routeStateFromPath } from './routeState';

const noQuery = { get: () => null };
const query = (values: Record<string, string>) => ({
	get: (name: string) => values[name] ?? null
});

describe('route state parsing (SCN-MEDIA-004)', () => {
	it('maps top-level routes to app views and sections', () => {
		expect(routeStateFromPath('/', {}, noQuery)).toEqual(defaultRouteState());
		expect(routeStateFromPath('/settings/profiles', {}, noQuery)).toMatchObject({
			view: 'settings',
			settingsSection: 'profiles'
		});
		expect(routeStateFromPath('/system/events', {}, noQuery)).toMatchObject({
			view: 'system',
			systemSection: 'events'
		});
		expect(routeStateFromPath('/activity', {}, noQuery)).toMatchObject({
			homeSection: 'activity'
		});
	});

	it('maps dynamic media, collection, discover, request, and search routes', () => {
		expect(routeStateFromPath('/search/advanced', {}, { get: () => 'matrix' })).toMatchObject({
			view: 'advanced-search',
			advancedQuery: 'matrix'
		});
		expect(
			routeStateFromPath('/discover/trending', { sectionId: 'trending' }, noQuery)
		).toMatchObject({
			view: 'discover-section',
			discoverSectionId: 'trending'
		});
		expect(
			routeStateFromPath(
				'/media/tmdb/movie/123/recommendations',
				{ provider: 'tmdb', type: 'movie', externalId: '123' },
				noQuery
			)
		).toMatchObject({
			view: 'related-section',
			relatedSectionKind: 'recommendations',
			metadataExternalId: '123'
		});
		expect(routeStateFromPath('/movies/movie-1/cast', { id: 'movie-1' }, noQuery)).toMatchObject({
			view: 'media-people',
			homeSection: 'movies',
			selectedMediaItemId: 'movie-1',
			peopleSectionKind: 'cast'
		});
		expect(
			routeStateFromPath('/people/tmdb/42', { provider: 'tmdb', personId: '42' }, noQuery)
		).toMatchObject({
			view: 'person-detail',
			personProvider: 'tmdb',
			personId: '42'
		});
		expect(routeStateFromPath('/requests/request-1', { id: 'request-1' }, noQuery)).toMatchObject({
			homeSection: 'requests',
			selectedRequestId: 'request-1'
		});
	});

	it('maps discover preset query strings to submenu entries', () => {
		expect(
			routeStateFromPath(
				'/discover/movies',
				{},
				query({ genres: 'Animation', withoutKeywords: 'anime' })
			)
		).toMatchObject({
			view: 'discover-movies',
			discoverSubmenuSection: 'animated-movies'
		});
		expect(
			routeStateFromPath('/discover/series', {}, query({ genres: 'Animation', keywords: 'anime' }))
		).toMatchObject({
			view: 'discover-series',
			discoverSubmenuSection: 'anime-series'
		});
	});

	it('creates stable route keys', () => {
		const route = routeStateFromPath('/settings/users', {}, noQuery);

		expect(appRouteKey(route)).toContain('settings');
		expect(appRouteKey(route)).toContain('users');
	});
});
