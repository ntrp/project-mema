import { describe, expect, it } from 'vitest';

import { canImportRows, type MatchDraft } from './libraryScanImport';
import {
	duplicateRemovalPathsForRows,
	importPayloadForRows,
	importPayloadForSingleRow
} from './libraryScanImportPayloads';
import { duplicateDraftStatesForRows, normalizeDuplicateDrafts } from './libraryScanDuplicates';
import type { LibraryScanItem, MediaSearchResult } from '$lib/settings/types';

describe('library scan duplicate imports', () => {
	it('recalculates duplicate validation from changed matches', () => {
		const rows = [scanItem('item-1'), scanItem('item-2')];
		const drafts = {
			'item-1': movieDraft({ matched: match('tmdb-1') }),
			'item-2': movieDraft({ matched: match('tmdb-1') })
		};

		expect(canImportRows(rows, drafts, 'profile-1')).toBe(false);
		expect(duplicateDraftStatesForRows(rows, drafts)['item-1']?.duplicate).toBe(true);

		drafts['item-2'].matched = match('tmdb-2');

		expect(canImportRows(rows, drafts, 'profile-1')).toBe(true);
		expect(duplicateDraftStatesForRows(rows, drafts)['item-1']).toBeUndefined();
	});

	it('does not import stale duplicate removals after a match changes', () => {
		const rows = [scanItem('item-1'), scanItem('item-2'), scanItem('item-3')];
		const drafts = {
			'item-1': movieDraft({ matched: match('tmdb-1') }),
			'item-2': movieDraft({
				matched: match('tmdb-2'),
				selected: false,
				removeDuplicate: true
			}),
			'item-3': movieDraft({ matched: undefined, selected: false, removeDuplicate: true })
		};

		expect(
			importPayloadForRows([rows[0]], rows, drafts, {
				qualityProfileId: '',
				monitorMode: 'only_media',
				minimumAvailability: 'released',
				seriesType: 'standard'
			}).removeDuplicatePaths
		).toEqual(['/downloads/item-3.mkv']);
	});

	it('builds single-row import payloads for sequential imports', () => {
		const rows = [scanItem('item-1'), scanItem('item-2'), scanItem('item-3')];
		const drafts = {
			'item-1': movieDraft({ matched: match('tmdb-1') }),
			'item-2': movieDraft({ matched: match('tmdb-2') }),
			'item-3': movieDraft({ matched: undefined, selected: false, removeDuplicate: true })
		};
		const duplicatePaths = duplicateRemovalPathsForRows(rows, drafts);
		const payload = importPayloadForSingleRow(rows[0], drafts, bulkOptions(), duplicatePaths);

		expect(payload.items.map((row) => row.itemId)).toEqual(['item-1']);
		expect(payload.removeDuplicatePaths).toEqual(['/downloads/item-3.mkv']);
	});

	it('does not recalculate every duplicate group for non-duplicate rows', () => {
		const rows = [scanItem('item-1'), scanItem('item-2'), { id: 'item-3' } as LibraryScanItem];
		const drafts = {
			'item-1': movieDraft({ selected: false, matched: match('tmdb-1'), removeDuplicate: true }),
			'item-2': movieDraft({ matched: match('tmdb-1') }),
			'item-3': movieDraft({ matched: match('tmdb-3') })
		};

		normalizeDuplicateDrafts(rows, drafts, rows[2].duplicateGroupId);

		expect(drafts['item-1'].removeDuplicate).toBe(true);
	});
});

function scanItem(id: string): LibraryScanItem {
	return {
		id,
		path: `/downloads/${id}.mkv`,
		duplicateGroupId: 'dup:tmdb-1',
		duplicateRemovalAllowed: true
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

function bulkOptions() {
	return {
		qualityProfileId: '',
		monitorMode: 'only_media',
		minimumAvailability: 'released',
		seriesType: 'standard'
	} as const;
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
