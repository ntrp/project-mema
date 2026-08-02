import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import LibraryScanImportTableBody from '$lib/components/settings/library/scan/LibraryScanImportTableBody.svelte';
import type { MatchDraft } from '$lib/components/settings/library/scan/libraryScanImport';
import type { LibraryScanItem } from '$lib/settings/types';

describe('rendered series scan rows', () => {
	it('renders one folder match row and an editable row for every episode file', () => {
		const rows = [episode('one', 1), episode('two', 2)];
		const drafts = Object.fromEntries(rows.map((row) => [row.id, draft(row)]));
		const { body } = render(LibraryScanImportTableBody, {
			props: {
				rows,
				folderPath: '/library',
				drafts,
				duplicateStates: {},
				qualityProfiles: [],
				metadataProviders: [],
				importingItemId: '',
				resettingItemId: '',
				onSearch: vi.fn(),
				onSelect: vi.fn(),
				onProviderChange: vi.fn(),
				onSeasonChange: vi.fn(),
				onEpisodeChange: vi.fn(),
				onResetImport: vi.fn()
			}
		});
		expect(body).toContain('Show');
		expect(body).toContain('Season 1/01.mkv');
		expect(body).toContain('Season 1/02.mkv');
		expect(body.match(/>Season</g)).toHaveLength(2);
		expect(body.match(/>Episode</g)).toHaveLength(2);
	});
});

function episode(id: string, episodeNumber: number): LibraryScanItem {
	return {
		id,
		scanId: 'scan',
		path: `/library/Show/Season 1/0${episodeNumber}.mkv`,
		fileName: `0${episodeNumber}.mkv`,
		seriesFolderPath: '/library/Show',
		detectedTitle: 'Show',
		detectedMediaKind: 'series',
		seasonNumber: 1,
		episodeNumber,
		status: 'pending',
		imported: false,
		duplicateRemovalAllowed: false,
		createdAt: '2026-07-27T00:00:00Z',
		updatedAt: '2026-07-27T00:00:00Z'
	};
}

function draft(item: LibraryScanItem): MatchDraft {
	return {
		selected: true,
		query: 'Show',
		mediaKind: 'series',
		metadataProviderId: '',
		matched: { title: 'Show', type: 'serie' },
		results: [],
		searching: false,
		searched: true,
		qualityProfileId: 'profile',
		monitorMode: 'all_episodes',
		minimumAvailability: 'released',
		seriesType: 'standard',
		seasonNumber: item.seasonNumber,
		episodeNumber: item.episodeNumber,
		removeDuplicate: false
	};
}
