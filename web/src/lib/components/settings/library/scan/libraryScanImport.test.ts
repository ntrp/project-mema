import { describe, expect, it } from 'vitest';

import {
	canImportRows,
	defaultQualityProfileId,
	isDraftSelectableForImport,
	type MatchDraft
} from '$lib/components/settings/library/scan/libraryScanImport';
import { applyEpisodeSeason } from '$lib/components/settings/library/scan/libraryScanBulk';
import { importRequestForDraft } from '$lib/components/settings/library/scan/libraryScanImportPayloads';
import {
	defaultMetadataProviderId,
	ensureScanDrafts,
	initialMatchDraft,
	matchFromScanItem
} from '$lib/components/settings/library/scan/libraryScanDrafts';
import type { LibraryScanItem, MediaSearchResult, MetadataProvider } from '$lib/settings/types';

describe('library scan import payloads', () => {
	it('requires season and episode numbers before selecting a series episode', () => {
		const item = { id: 'episode', status: 'pending', imported: false } as LibraryScanItem;
		const draft = movieDraft({
			mediaKind: 'series',
			seasonNumber: undefined,
			episodeNumber: 1
		});
		expect(isDraftSelectableForImport(item, draft)).toBe(false);
		draft.seasonNumber = 1;
		expect(isDraftSelectableForImport(item, draft)).toBe(true);
		draft.episodeNumber = undefined;
		expect(isDraftSelectableForImport(item, draft)).toBe(false);
	});

	it('uses the footer quality profile when row draft has not been updated yet', () => {
		const item = { id: 'item-1' } as LibraryScanItem;
		const draft = movieDraft({ qualityProfileId: '' });
		expect(canImportRows([item], { [item.id]: draft }, 'profile-1')).toBe(true);
		expect(
			importRequestForDraft(draft, draft.matched!, {
				qualityProfileId: 'profile-1',
				monitorMode: 'only_media',
				minimumAvailability: 'released',
				seriesType: 'standard'
			}).qualityProfileId
		).toBe('profile-1');
	});

	it('keeps the row quality profile when the footer profile is empty', () => {
		const draft = movieDraft({ qualityProfileId: 'profile-row' });
		expect(
			importRequestForDraft(draft, draft.matched!, {
				qualityProfileId: '',
				monitorMode: 'only_media',
				minimumAvailability: 'released',
				seriesType: 'standard'
			}).qualityProfileId
		).toBe('profile-row');
	});

	it('defaults metadata provider selection to tmdb when it is enabled', () => {
		const providers = [
			{ id: 'tvdb-1', name: 'TVDB', type: 'tvdb', enabled: true },
			{ id: 'tmdb-1', name: 'TMDB', type: 'tmdb', enabled: true }
		] as MetadataProvider[];
		expect(defaultMetadataProviderId(providers, 'movie')).toBe('tmdb-1');
	});

	it('uses the configured default profile as the import profile default', () => {
		expect(
			defaultQualityProfileId([
				{ id: 'profile-1', name: 'HD' },
				{ id: 'profile-2', name: 'UHD', isDefault: true }
			])
		).toBe('profile-2');
	});

	it('applies a season to selected episodes and resets their episode numbers', () => {
		const rows = [{ id: 'one' }, { id: 'two' }] as LibraryScanItem[];
		const drafts = {
			one: movieDraft({ mediaKind: 'series', seasonNumber: 1, episodeNumber: 2 }),
			two: movieDraft({ mediaKind: 'series', seasonNumber: 1, episodeNumber: 3 })
		};
		applyEpisodeSeason(rows, drafts, 4);
		expect(drafts.one).toMatchObject({ seasonNumber: 4, episodeNumber: undefined });
		expect(drafts.two).toMatchObject({ seasonNumber: 4, episodeNumber: undefined });
		drafts.one.episodeNumber = 5;
		const request = importRequestForDraft(drafts.one, drafts.one.matched!, {
			qualityProfileId: '',
			monitorMode: 'all_episodes',
			minimumAvailability: 'released',
			seriesType: 'standard'
		});
		expect(request).toMatchObject({ seasonNumber: 4, episodeNumber: 5 });
	});

	it('prefills monitor mode from matched media type', () => {
		const bulk = {
			qualityProfileId: '',
			monitorMode: 'none',
			minimumAvailability: 'released',
			seriesType: 'standard'
		} as const;
		expect(initialMatchDraft(matchedItem('movie'), [], bulk)).toMatchObject({
			monitorMode: 'only_media',
			selected: true
		});
		expect(initialMatchDraft(matchedItem('series'), [], bulk)).toMatchObject({
			monitorMode: 'all_episodes',
			selected: false
		});
	});

	it('shows imported scan rows as matched even when rescans only carry media item links', () => {
		expect(
			matchFromScanItem({
				id: 'item-imported',
				imported: true,
				mediaItemId: 'media-1',
				detectedTitle: 'Imported Movie',
				detectedYear: 2026,
				detectedMediaKind: 'movie'
			} as LibraryScanItem)
		).toMatchObject({
			id: 'media-1',
			title: 'Imported Movie',
			type: 'movie',
			year: 2026
		});
	});

	it('rebuilds a row draft when reset import clears the scan item state', () => {
		const item = {
			id: 'item-imported',
			scanId: 'scan-1',
			path: '/library/Imported.Movie.2026.mkv',
			fileName: 'Imported.Movie.2026.mkv',
			imported: true,
			status: 'manually_added',
			mediaItemId: 'media-1',
			detectedTitle: 'Imported Movie',
			detectedYear: 2026,
			detectedMediaKind: 'movie',
			duplicateRemovalAllowed: false,
			createdAt: '2026-07-07T00:00:00Z',
			updatedAt: '2026-07-07T00:00:00Z'
		} as LibraryScanItem;
		const drafts: Record<string, MatchDraft> = {};
		const sources: Record<string, string> = {};
		const bulk = {
			qualityProfileId: 'profile-1',
			monitorMode: 'only_media',
			minimumAvailability: 'released',
			seriesType: 'standard'
		} as const;

		ensureScanDrafts([item], drafts, [], bulk, sources);
		expect(drafts[item.id].matched).toMatchObject({ id: 'media-1', title: 'Imported Movie' });

		ensureScanDrafts(
			[
				{
					...item,
					imported: false,
					status: 'pending',
					mediaItemId: undefined
				} as LibraryScanItem
			],
			drafts,
			[],
			bulk,
			sources
		);

		expect(drafts[item.id]).toMatchObject({
			selected: false,
			query: 'Imported Movie',
			searching: false,
			searched: false,
			removeDuplicate: false
		});
		expect(drafts[item.id].matched).toBeUndefined();
	});
});

function matchedItem(kind: 'movie' | 'series'): LibraryScanItem {
	return {
		id: `${kind}-item`,
		path: `/downloads/${kind}.mkv`,
		fileName: `${kind}.mkv`,
		detectedTitle: 'Scenario',
		detectedMediaKind: kind,
		status: 'pending',
		imported: false,
		matchedTitle: 'Scenario'
	} as LibraryScanItem;
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
