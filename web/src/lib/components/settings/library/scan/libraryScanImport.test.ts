import { describe, expect, it } from 'vitest';

import {
	canImportRows,
	importRequestForDraft,
	type MatchDraft
} from '$lib/components/settings/library/scan/libraryScanImport';
import type { LibraryScanItem, MediaSearchResult } from '$lib/settings/types';

describe('library scan import payloads', () => {
	it('uses the footer quality profile when row draft has not been updated yet', () => {
		const item = { id: 'item-1' } as LibraryScanItem;
		const draft = {
			selected: true,
			query: 'Scenario Movie',
			mediaKind: 'movie',
			matched: { title: 'Scenario Movie', type: 'movie', year: 2026 } as MediaSearchResult,
			results: [],
			searching: false,
			searched: true,
			qualityProfileId: '',
			monitorMode: 'only_media',
			minimumAvailability: 'released'
		} as MatchDraft;

		expect(canImportRows([item], { [item.id]: draft }, 'profile-1')).toBe(true);
		expect(
			importRequestForDraft(draft, draft.matched!, {
				qualityProfileId: 'profile-1',
				monitorMode: 'only_media',
				minimumAvailability: 'released'
			}).qualityProfileId
		).toBe('profile-1');
	});
});
