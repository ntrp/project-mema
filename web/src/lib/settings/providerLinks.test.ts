import { describe, expect, it } from 'vitest';

import { providerDisplayName, providerPageUrl } from './providerLinks';

describe('provider links (SCN-SETTINGS-012)', () => {
	it('builds provider page URLs for movies and series', () => {
		expect(providerPageUrl('tmdb', 'movie', '123')).toBe('https://www.themoviedb.org/movie/123');
		expect(providerPageUrl('tmdb', 'serie', '456')).toBe('https://www.themoviedb.org/tv/456');
		expect(providerPageUrl('tvdb', 'movie', '789')).toBe(
			'https://thetvdb.com/dereferrer/movie/789'
		);
		expect(providerPageUrl('tvdb', 'serie', '101')).toBe(
			'https://thetvdb.com/dereferrer/series/101'
		);
		expect(providerPageUrl('tvdb', 'movie', '789', 'https://thetvdb.com/movies/walle')).toBe(
			'https://thetvdb.com/movies/walle'
		);
	});

	it('falls back for missing or unknown providers', () => {
		expect(providerPageUrl(undefined, 'movie', '123')).toBeUndefined();
		expect(providerPageUrl('unknown', 'movie', '123')).toBeUndefined();
		expect(providerDisplayName('tmdb')).toBe('TMDB');
		expect(providerDisplayName('tvdb')).toBe('TVDB');
		expect(providerDisplayName('imdb')).toBe('IMDB');
		expect(providerDisplayName()).toBe('Provider');
	});
});
