import { describe, expect, it } from 'vitest';

import {
	autoMatchResult,
	normalizedMatchTitle
} from '$lib/components/settings/library/scan/libraryScanAutoMatch';
import {
	applyScanItemProvider,
	prepareProviderSearch,
	searchScanItem
} from '$lib/components/settings/library/scan/libraryScanTableActions';
import type { MatchDraft } from '$lib/components/settings/library/scan/libraryScanImport';
import type { LibraryScanItem, MediaSearchResult } from '$lib/settings/types';

describe('library scan auto matching', () => {
	it('matches titles after stripping accents and special characters', () => {
		expect(normalizedMatchTitle('Amélie')).toBe('amelie');
		expect(normalizedMatchTitle('Wall-e')).toBe('walle');
		expect(normalizedMatchTitle('WALL·E')).toBe('walle');
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
		expect(
			autoMatchResult(
				{
					...detectedScanItem('item-3', 'Wall E', 2008),
					fileName: 'WALL-E.2008.tmdb-10681.1080p.WEB-DL.mkv'
				},
				[match('tmdb-other', 'movie', 'Wall Street', 1987), match('10681', 'movie', 'WALL·E', 2008)]
			)
		).toMatchObject({ title: 'WALL·E' });
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
		const item = detectedScanItem('item-5', 'Scenario Movie', 2026);
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
