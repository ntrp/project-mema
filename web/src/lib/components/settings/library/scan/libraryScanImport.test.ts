import { describe, expect, it } from 'vitest';

import {
	canImportRows,
	defaultQualityProfileId,
	type MatchDraft
} from '$lib/components/settings/library/scan/libraryScanImport';
import { importRequestForDraft } from '$lib/components/settings/library/scan/libraryScanImportPayloads';
import {
	defaultMetadataProviderId,
	initialMatchDraft
} from '$lib/components/settings/library/scan/libraryScanDrafts';
import {
	applyMovieOptions,
	applySeriesOptions,
	matchedRowsByKind
} from '$lib/components/settings/library/scan/libraryScanBulk';
import type { LibraryScanItem, MediaSearchResult, MetadataProvider } from '$lib/settings/types';

describe('library scan import payloads', () => {
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

	it('prefills monitor mode from matched media type', () => {
		const bulk = {
			qualityProfileId: '',
			monitorMode: 'none',
			minimumAvailability: 'released',
			seriesType: 'standard'
		} as const;
		expect(initialMatchDraft(matchedItem('movie'), [], bulk).monitorMode).toBe('only_media');
		expect(initialMatchDraft(matchedItem('series'), [], bulk).monitorMode).toBe('all_episodes');
	});

	it('applies footer movie and series sections only to matched media of that type', () => {
		const rows = [
			{ id: 'movie-1' } as LibraryScanItem,
			{ id: 'series-1' } as LibraryScanItem,
			{ id: 'unmatched-1' } as LibraryScanItem
		];
		const drafts = {
			'movie-1': movieDraft({ matched: match('movie-1', 'movie') }),
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
		expect(drafts['movie-1'].monitorMode).toBe('collection');
		expect(drafts['movie-1'].minimumAvailability).toBe('announced');
		expect(drafts['series-1'].monitorMode).toBe('missing_episodes');
		expect(drafts['series-1'].seriesType).toBe('daily');
		expect(drafts['unmatched-1'].monitorMode).toBe('none');
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
