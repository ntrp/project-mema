import { describe, expect, it } from 'vitest';

import {
	applyAutoMatch,
	autoMatchResult,
	normalizedMatchTitle
} from '$lib/components/settings/library/scan/libraryScanAutoMatch';
import {
	applyScanItemProvider,
	prepareProviderSearch,
	searchScanItem
} from '$lib/components/settings/library/scan/libraryScanTableActions';
import { applyMovieOptions, applySeriesOptions, matchedRowsByKind } from './libraryScanBulk';
import type { MatchDraft } from '$lib/components/settings/library/scan/libraryScanImport';
import type { LibraryScanItem, MediaSearchResult } from '$lib/settings/types';

describe('library scan auto matching', () => {
	it('matches titles after stripping accents and special characters', () => {
		expect(normalizedMatchTitle('Amélie')).toBe('amelie');
		expect(normalizedMatchTitle('Wall-e')).toBe('walle');
		expect(normalizedMatchTitle('WALL·E')).toBe('walle');
		expect(normalizedMatchTitle('The.Matrix.1999.[imdbid-tt0133093].1080p')).toBe('thematrix');
		expect(
			autoMatchResult(detectedScanItem('item-1', 'Amelie', 2001), [
				match('tmdb-other', 'movie', 'Amelia', 2001),
				match('tmdb-194', 'movie', 'Amélie', 2001)
			])
		).toMatchObject({ title: 'Amélie' });
		expect(
			autoMatchResult(detectedScanItem('item-2', 'Walle', 2008), [
				match('tmdb-other', 'movie', 'Wall Street', 1987),
				match('tmdb-10681', 'movie', 'WALL·E', 2008)
			])
		).toMatchObject({ title: 'WALL·E' });
	});

	it('matches provider ids embedded in filenames', () => {
		for (const [fileName, provider, id] of [
			['WALL-E.2008.tmdb-10681.1080p.WEB-DL.mkv', 'tmdb', '10681'],
			['WALL-E.2008.tt0910970.mkv', 'imdb', 'tt0910970'],
			['WALL-E.2008.tvdb-12345.mkv', 'tvdb', '12345']
		]) {
			const result = { ...match(id, 'movie', 'WALL·E', 2008), externalProvider: provider };
			expect(
				autoMatchResult({ ...detectedScanItem('item-3', 'Wall E', 2008), fileName }, [result])
			).toBe(result);
		}
		expect(
			autoMatchResult({ ...detectedScanItem('wrong', 'Wall E', 2008), fileName: 'tmdb-999.mkv' }, [
				match('10681')
			])
		).toBeUndefined();
	});

	it('falls back to title matching when the search provider cannot return the embedded id', () => {
		const result = match('1413', 'serie', 'The Apothecary Diaries', 2023);
		expect(
			autoMatchResult(
				{
					...detectedScanItem('item-imdb', 'The Apothecary Diaries', 2023),
					path: 'The Apothecary Diaries (2023) [imdbid-tt26743760]/Season 1/01.mkv'
				},
				[result]
			)
		).toBe(result);
	});

	it('leaves episodes unselected when a series match is applied', () => {
		const rows = [
			{
				...detectedScanItem('episode-1', 'The Apothecary Diaries', 2023),
				detectedMediaKind: 'series',
				seriesFolderPath: '/series/The Apothecary Diaries'
			},
			{
				...detectedScanItem('episode-2', 'The Apothecary Diaries', 2023),
				detectedMediaKind: 'series',
				seriesFolderPath: '/series/The Apothecary Diaries'
			}
		] as LibraryScanItem[];
		const drafts = Object.fromEntries(
			rows.map((row) => [row.id, movieDraft({ mediaKind: 'series', selected: true })])
		);

		applyAutoMatch(rows[0], match('1413', 'serie', 'The Apothecary Diaries', 2023), rows, drafts);

		expect(rows.map((row) => drafts[row.id].selected)).toEqual([false, false]);
	});

	it('stops the auto matching spinner when the provider title changes the query', async () => {
		const item = detectedScanItem('item-4', 'Amelie', 2001);
		const drafts = {
			[item.id]: movieDraft({
				matched: undefined,
				query: 'Amelie',
				searching: false,
				searched: false
			})
		};

		await searchScanItem({
			item,
			allRows: [item],
			drafts,
			searchCache: {},
			auto: true,
			onSearchMatch: async () => [match('tmdb-194', 'movie', 'Amélie', 2001)]
		});

		expect(drafts[item.id].matched).toMatchObject({ title: 'Amélie' });
		expect(drafts[item.id].query).toBe('Amélie');
		expect(drafts[item.id].searching).toBe(false);
		expect(drafts[item.id].searched).toBe(true);
	});

	it('clears stale matches before searching another provider', () => {
		const item = detectedScanItem('item-5', 'Scenario.Movie.2026.[imdbid-tt1234567].1080p', 2026);
		const draft = movieDraft({
			matched: match('tmdb-1'),
			query: 'Old Provider Title',
			results: [match('tmdb-1')],
			selected: true,
			searched: true
		});

		prepareProviderSearch(item, draft);

		expect(draft.query).toBe('Scenario Movie');
		expect(draft.matched).toBeUndefined();
		expect(draft.results).toEqual([]);
		expect(draft.selected).toBe(false);
		expect(draft.searched).toBe(false);
	});

	it('applies bulk options only to matched media of that type', () => {
		const rows = ['movie-1', 'series-1', 'unmatched-1'].map((id) => ({ id }) as LibraryScanItem);
		const drafts = {
			'movie-1': movieDraft({ matched: match('movie-1') }),
			'series-1': movieDraft({ matched: match('series-1', 'serie'), monitorMode: 'all_episodes' }),
			'unmatched-1': movieDraft({ matched: undefined, monitorMode: 'none' })
		};
		applyMovieOptions(matchedRowsByKind(rows, drafts, 'movie'), drafts, 'collection', 'announced');
		applySeriesOptions(
			matchedRowsByKind(rows, drafts, 'series'),
			drafts,
			'missing_episodes',
			'daily'
		);
		expect(drafts['movie-1']).toMatchObject({
			monitorMode: 'collection',
			minimumAvailability: 'announced'
		});
		expect(drafts['series-1']).toMatchObject({
			monitorMode: 'missing_episodes',
			seriesType: 'daily'
		});
		expect(drafts['unmatched-1'].monitorMode).toBe('none');
	});

	it('applies footer provider changes and starts a fresh auto search', () => {
		const item = detectedScanItem('item-6', 'Scenario Movie', 2026);
		const drafts = {
			[item.id]: movieDraft({
				metadataProviderId: 'metadata-1',
				matched: match('tmdb-1'),
				selected: true,
				searched: true
			})
		};
		const searched: string[] = [];

		applyScanItemProvider({
			rows: [item],
			drafts,
			providerId: 'metadata-2',
			search: async (row, auto) => {
				searched.push(`${row.id}:${auto}`);
			}
		});

		expect(drafts[item.id].metadataProviderId).toBe('metadata-2');
		expect(drafts[item.id].matched).toBeUndefined();
		expect(drafts[item.id].selected).toBe(false);
		expect(searched).toEqual(['item-6:true']);
	});
});

function detectedScanItem(id: string, title: string, year: number): LibraryScanItem {
	return {
		id,
		path: `/downloads/${id}.mkv`,
		detectedTitle: title,
		detectedYear: year,
		detectedMediaKind: 'movie'
	} as LibraryScanItem;
}

function match(
	externalId: string,
	type: MediaSearchResult['type'] = 'movie',
	title = 'Scenario Movie',
	year = 2026
): MediaSearchResult {
	return {
		title,
		type,
		year,
		externalProvider: 'tmdb',
		externalId
	};
}

function movieDraft(overrides: Partial<MatchDraft> = {}): MatchDraft {
	return {
		selected: true,
		query: 'Scenario Movie',
		mediaKind: 'movie',
		metadataProviderId: 'metadata-1',
		matched: { title: 'Scenario Movie', type: 'movie', year: 2026 } as MediaSearchResult,
		results: [],
		searching: false,
		searched: true,
		qualityProfileId: 'profile-1',
		monitorMode: 'only_media',
		minimumAvailability: 'released',
		seriesType: 'standard',
		removeDuplicate: false,
		...overrides
	};
}
