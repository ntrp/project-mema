import { describe, expect, it } from 'vitest';

import { duplicateDraftStatesForRows } from '$lib/components/settings/library/scan/libraryScanDuplicates';
import type { MatchDraft } from '$lib/components/settings/library/scan/libraryScanImport';
import type { LibraryScanItem, MediaSearchResult } from '$lib/settings/types';

describe('library scan episode duplicates', () => {
	it('resolves series duplicates using edited season and episode numbers', () => {
		const rows = [11, 21, 33].map(
			(episodeNumber) =>
				({
					id: `episode-${episodeNumber}`,
					detectedMediaKind: 'series',
					duplicateGroupId: 'dup:show'
				}) as LibraryScanItem
		);
		const matched = {
			id: 'media-show',
			title: 'The Apothecary Diaries',
			type: 'serie'
		} as MediaSearchResult;
		const drafts = Object.fromEntries(
			rows.map((row, index) => [row.id, seriesDraft(matched, [11, 21, 33][index])])
		);

		expect(duplicateDraftStatesForRows(rows, drafts)).toEqual({});
		drafts['episode-21'].episodeNumber = 11;
		expect(duplicateDraftStatesForRows(rows, drafts)['episode-11']?.duplicate).toBe(true);
		drafts['episode-21'].episodeNumber = undefined;
		expect(duplicateDraftStatesForRows(rows, drafts)).toEqual({});
	});
});

function seriesDraft(matched: MediaSearchResult, episodeNumber: number): MatchDraft {
	return {
		selected: false,
		query: matched.title,
		mediaKind: 'series',
		metadataProviderId: '',
		matched,
		results: [],
		searching: false,
		searched: true,
		qualityProfileId: 'profile',
		monitorMode: 'all_episodes',
		minimumAvailability: 'released',
		seriesType: 'standard',
		seasonNumber: 2,
		episodeNumber,
		removeDuplicate: false
	};
}
