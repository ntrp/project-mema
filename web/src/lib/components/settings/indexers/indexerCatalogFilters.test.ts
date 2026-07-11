import { describe, expect, it } from 'vitest';

import {
	entryHasCategory,
	flattenCategories,
	fuzzyMatch,
	matches,
	matchesAny,
	privacyGroup,
	privacyLabel,
	textMatchRank,
	unique,
	uniqueCategories
} from './indexerCatalogFilters';

describe('indexer catalog filters', () => {
	it('normalizes scalar filters and unique values', () => {
		expect(unique(['torrent', '', 'torrent', 'usenet'])).toEqual(['torrent', 'usenet']);
		expect(matches('all', 'public')).toBe(true);
		expect(matches('private', 'public')).toBe(false);
		expect(matchesAny([], 'movie')).toBe(true);
		expect(matchesAny(['series'], 'movie')).toBe(false);
		expect(privacyGroup('semiPrivate')).toBe('private');
		expect(privacyLabel('semi-private')).toBe('semi-private');
	});

	it('flattens, deduplicates, and matches nested categories', () => {
		const categories = [
			{ id: 2, name: 'TV', children: [{ id: 3, name: 'Anime', children: [] }] },
			{ id: 1, name: 'Movies', children: [] }
		];
		const entries = [
			{ capabilities: { categories } },
			{ capabilities: { categories: [{ id: 2, name: 'Television', children: [] }] } }
		] as never;
		expect(flattenCategories(categories as never).map((item) => item.id)).toEqual([2, 3, 1]);
		expect(uniqueCategories(entries)).toEqual([
			{ id: 1, name: 'Movies' },
			{ id: 2, name: 'Television' },
			{ id: 3, name: 'Anime' }
		]);
		expect(entryHasCategory(entries[0], 3)).toBe(true);
		expect(entryHasCategory(entries[0], 99)).toBe(false);
	});

	it('ranks exact, fuzzy, empty, and missing text matches', () => {
		expect(fuzzyMatch('', ['Anything'])).toBe(true);
		expect(fuzzyMatch('nzb', ['Newznab'])).toBe(true);
		expect(fuzzyMatch('xyz', ['Newznab'])).toBe(false);
		expect(textMatchRank('', ['Value'])).toBe(0);
		expect(textMatchRank('newz', ['Newznab'])).toBe(0);
		expect(textMatchRank('nzb', ['Newznab'])).toBe(1);
		expect(textMatchRank('xyz', ['Newznab'])).toBe(-1);
	});
});
