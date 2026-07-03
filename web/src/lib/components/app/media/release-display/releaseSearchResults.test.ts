import { describe, expect, it } from 'vitest';

import {
	activeFilterCount,
	defaultReleaseFilters,
	filteredSortedReleases,
	releaseQualityOptions
} from '$lib/components/app/media/release-display/releaseSearchResults';
import type { MediaItem, ReleaseCandidate } from '$lib/settings/types';

const item = { title: 'Scenario Movie', type: 'movie' } as MediaItem;

function match(overrides: Partial<ReleaseCandidate['match']> = {}): ReleaseCandidate['match'] {
	return {
		severity: 'info',
		details: [],
		qualityId: 'webdl-1080p',
		quality: 'WEBDL-1080p',
		score: 100,
		scoreContributors: [],
		languages: ['English'],
		matchedMedia: 'Scenario Movie',
		titleMatched: true,
		yearMatched: true,
		seasonEpisodeMatched: true,
		rejected: false,
		parsed: {},
		...overrides
	} as ReleaseCandidate['match'];
}

function release(
	overrides: Partial<Omit<ReleaseCandidate, 'match'>> & {
		match?: Partial<ReleaseCandidate['match']>;
	}
): ReleaseCandidate {
	return {
		title: 'Scenario.Movie.2026.1080p.WEB-DL',
		indexerName: 'Indexer',
		indexerType: 'torznab',
		sizeBytes: 8 * 1024 * 1024 * 1024,
		seeders: 10,
		peers: 12,
		publishedAt: '2026-01-02T00:00:00Z',
		...overrides,
		match: match(overrides.match)
	} as ReleaseCandidate;
}

describe('release result filtering and sorting (SCN-MEDIA-002)', () => {
	it('counts active filters against defaults', () => {
		expect(activeFilterCount(defaultReleaseFilters())).toBe(0);
		expect(
			activeFilterCount({ ...defaultReleaseFilters(), source: 'torrent', minScore: '50' })
		).toBe(2);
	});

	it('builds quality options from defaults plus release matches', () => {
		expect(
			releaseQualityOptions([
				release({ match: { severity: 'info', quality: 'Scenario-4K', score: 200, languages: [] } })
			])
		).toContain('Scenario-4K');
	});

	it('filters by source, quality, size, and score', () => {
		const visible = filteredSortedReleases(
			item,
			[
				release({ title: 'Torrent Match', indexerType: 'torznab', sizeBytes: 5 * 1024 ** 3 }),
				release({ title: 'NZB Match', indexerType: 'newznab', sizeBytes: 5 * 1024 ** 3 }),
				release({ title: 'Too Large', indexerType: 'torznab', sizeBytes: 15 * 1024 ** 3 })
			],
			{
				source: 'torrent',
				quality: 'WEBDL-1080p',
				minSize: '4',
				maxSize: '10',
				minScore: '90',
				maxScore: '150'
			},
			{ key: 'title', direction: 'asc' }
		);

		expect(visible.map((value) => value.title)).toEqual(['Torrent Match']);
	});

	it('keeps severe match problems ahead of normal sorted results', () => {
		const visible = filteredSortedReleases(
			item,
			[
				release({
					title: 'Normal A',
					match: { severity: 'info', quality: 'WEBDL-1080p', score: 10, languages: [] }
				}),
				release({
					title: 'Warning B',
					match: { severity: 'warning', quality: 'WEBDL-1080p', score: 20, languages: [] }
				}),
				release({
					title: 'Error C',
					match: { severity: 'error', quality: 'WEBDL-1080p', score: 30, languages: [] }
				})
			],
			defaultReleaseFilters(),
			{ key: 'score', direction: 'desc' }
		);

		expect(visible.map((value) => value.title)).toEqual(['Normal A', 'Warning B', 'Error C']);
	});

	it('sorts by release metadata fields and ignores invalid numeric filters', () => {
		const releases = [
			release({
				title: 'Second',
				indexerName: 'B Indexer',
				indexerType: 'newznab',
				sizeBytes: 2 * 1024 ** 3,
				peers: undefined,
				seeders: 3,
				publishedAt: undefined,
				match: { severity: 'info', quality: 'WEBDL-720p', score: 20, languages: [] }
			}),
			release({
				title: 'First',
				indexerName: 'A Indexer',
				indexerType: 'torznab',
				sizeBytes: 4 * 1024 ** 3,
				peers: 8,
				seeders: 1,
				publishedAt: '2026-07-03T00:00:00Z',
				match: { severity: 'info', quality: 'WEBDL-1080p', score: 40, languages: [] }
			})
		];
		const filters = {
			...defaultReleaseFilters(),
			minSize: 'not-a-number',
			maxScore: 'not-a-number'
		};

		expect(
			filteredSortedReleases(item, releases, filters, { key: 'source', direction: 'asc' }).map(
				(value) => value.title
			)
		).toEqual(['Second', 'First']);
		expect(
			filteredSortedReleases(item, releases, filters, { key: 'peers', direction: 'desc' }).map(
				(value) => value.title
			)
		).toEqual(['First', 'Second']);
		expect(
			filteredSortedReleases(item, releases, filters, { key: 'age', direction: 'desc' }).map(
				(value) => value.title
			)
		).toEqual(['First', 'Second']);
		expect(
			filteredSortedReleases(item, releases, filters, { key: 'match', direction: 'asc' }).map(
				(value) => value.title
			)
		).toEqual(['Second', 'First']);
	});
});
