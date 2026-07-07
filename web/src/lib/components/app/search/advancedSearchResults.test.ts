import { describe, expect, it } from 'vitest';

import { externalMediaUrl } from './advancedSearchResults';
import type { MediaSearchResult } from '$lib/settings/types';

describe('advanced search results', () => {
	it('uses provider supplied external URLs before deriving provider links', () => {
		expect(
			externalMediaUrl({
				title: 'WALL-E',
				type: 'movie',
				externalProvider: 'tvdb',
				externalId: '516',
				externalUrl: 'https://thetvdb.com/movies/walle'
			} as MediaSearchResult)
		).toBe('https://thetvdb.com/movies/walle');

		expect(
			externalMediaUrl({
				title: 'WALL-E',
				type: 'movie',
				externalProvider: 'tvdb',
				externalId: '516'
			} as MediaSearchResult)
		).toBe('https://thetvdb.com/dereferrer/movie/516');
	});
});
