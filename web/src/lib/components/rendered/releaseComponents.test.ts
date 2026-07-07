import { render } from 'svelte/server';
import { describe, expect, it, vi } from 'vitest';

import ReleaseSearchFilters from '$lib/components/app/media/release-search/ReleaseSearchFilters.svelte';
import ReleaseSearchQueryInput from '$lib/components/app/media/release-search/ReleaseSearchQueryInput.svelte';
import ReleaseSearchResultsTable from '$lib/components/app/media/release-search/ReleaseSearchResultsTable.svelte';
import ReleaseSearchStatusLog from '$lib/components/app/media/release-search/ReleaseSearchStatusLog.svelte';
import SubtitleSearchResultsTable from '$lib/components/app/media/subtitle-search/SubtitleSearchResultsTable.svelte';
import { defaultReleaseFilters } from '$lib/components/app/media/release-display/releaseSearchResults';
import type { ReleaseSort } from '$lib/components/app/media/release-display/releaseSearchResults';
import type { MediaItem, ReleaseCandidate, SubtitleCandidate } from '$lib/settings/types';
import { renderWithTooltip } from './renderHelpers';

const mediaItem = {
	id: 'movie-1',
	type: 'movie',
	title: 'Scenario Movie'
} as MediaItem;
const scoreSort: ReleaseSort = { key: 'score', direction: 'desc' };
const titleSort: ReleaseSort = { key: 'title', direction: 'asc' };

function release(overrides: Partial<ReleaseCandidate> = {}): ReleaseCandidate {
	return {
		id: 'release-1',
		title: 'Scenario.Movie.2026.1080p.WEBDL.German-GRP',
		indexerName: 'Local Torznab',
		indexerProtocol: 'torrent',
		sizeBytes: 8 * 1024 ** 3,
		seeders: 42,
		peers: 7,
		publishedAt: '2026-07-03T04:00:00Z',
		match: {
			severity: 'info',
			details: ['Matches the requested resource.'],
			qualityId: 'webdl-1080p',
			quality: 'WEBDL-1080p',
			score: 1200,
			scoreContributors: [{ label: 'Quality', score: 1000 }],
			languages: ['German'],
			matchedMedia: 'Scenario Movie',
			customFormatScore: 0,
			customFormatContributors: [],
			languageContributors: [],
			rankContributors: [{ label: 'Seeders', score: 42 }],
			parsed: {
				release: { releaseTitle: 'Scenario.Movie.2026.1080p.WEBDL.German-GRP' },
				quality: { qualityId: 'webdl-1080p', quality: 'WEBDL-1080p' },
				languages: ['German'],
				details: { releaseType: 'movie', customFormatNames: [], matchedSpecCount: 0 }
			}
		},
		...overrides
	} as ReleaseCandidate;
}

function subtitleCandidate(overrides: Partial<SubtitleCandidate> = {}): SubtitleCandidate {
	return {
		id: 'subtitle-1',
		protocol: 'HTTP',
		providerId: 'provider-1',
		providerName: 'Mock Subtitles',
		title: 'Scenario.Movie.2026.english.srt',
		languageId: 'english',
		format: 'srt',
		match: {
			severity: 'success',
			label: 'Match',
			details: ['Language matches english']
		},
		...overrides
	} as SubtitleCandidate;
}

describe('rendered release components (SCN-MEDIA-002)', () => {
	it('renders release result headers, release metadata, and grab actions', () => {
		const { body } = renderWithTooltip(ReleaseSearchResultsTable, {
			item: mediaItem,
			releases: [release()],
			searching: false,
			sort: scoreSort,
			canManage: true,
			onSort: vi.fn(),
			onGrab: vi.fn()
		});

		expect(body).toContain('Sort by Score');
		expect(body).toContain('Sort by Peers');
		expect(body).not.toContain('Sort by Match');
		expect(body).toContain('Match');
		expect(body).toContain('Local Torznab');
		expect(body).toContain('Scenario.Movie.2026.1080p.WEBDL.German-GRP');
		expect(body).toContain('8.0 GiB');
		expect(body).toContain('WEBDL-1080p');
		expect(body).toContain('+1200');
		expect(body).toContain('Grab');
	});

	it('disables release grab actions when no compatible download client exists', () => {
		const { body } = renderWithTooltip(ReleaseSearchResultsTable, {
			item: mediaItem,
			releases: [
				release({ grabDisabledReason: 'No enabled torrent download client is configured' })
			],
			searching: false,
			sort: scoreSort,
			canManage: true,
			onSort: vi.fn(),
			onGrab: vi.fn()
		});

		expect(body).toContain('disabled');
		expect(body).toContain('Grab release');
	});

	it('renders release searching state without empty table rows', () => {
		const { body } = renderWithTooltip(ReleaseSearchResultsTable, {
			item: mediaItem,
			releases: [],
			searching: true,
			sort: titleSort,
			canManage: true,
			onSort: vi.fn(),
			onGrab: vi.fn()
		});

		expect(body).toContain('Searching releases');
		expect(body).not.toContain('Sort by Title');
	});

	it('renders subtitle search result headers with language after title', () => {
		const { body } = renderWithTooltip(SubtitleSearchResultsTable, {
			candidates: [subtitleCandidate()],
			searching: false,
			canManage: true,
			onGrab: vi.fn()
		});

		expect(body).toContain('Protocol');
		expect(body).toContain('Provider');
		expect(body).toContain('Title / filename');
		expect(body).toContain('Language');
		expect(body.indexOf('Title / filename')).toBeLessThan(body.indexOf('Language'));
		expect(body).toContain('HTTP');
		expect(body).toContain('Mock Subtitles');
		expect(body).toContain('Scenario.Movie.2026.english.srt');
		expect(body).toContain('English');
		expect(body).toContain('Grab subtitle');
	});

	it('renders release filter controls with selected values and quality options', () => {
		const { body } = render(ReleaseSearchFilters, {
			props: {
				filters: {
					...defaultReleaseFilters(),
					source: 'usenet',
					quality: 'WEBDL-1080p',
					minSize: '4',
					maxSize: '12',
					minScore: '100'
				},
				qualityOptions: ['WEBDL-1080p', 'Remux-2160p'],
				onChange: vi.fn(),
				onReset: vi.fn()
			}
		});

		expect(body).toContain('Protocol');
		expect(body).toContain('USENET');
		expect(body).toContain('Size GiB');
		expect(body).toContain('WEBDL-1080p');
		expect(body).toContain('Reset');
	});

	it('renders the release search query override control and default query affordance', () => {
		const { body } = renderWithTooltip(ReleaseSearchQueryInput, {
			overrideQuery: false,
			customQuery: 'Scenario Movie 2026',
			queryVariants: ['Scenario Movie 2026', 'Scenario Movie'],
			disabled: false
		});

		expect(body).toContain('Search query');
		expect(body).toContain('Scenario Movie 2026');
		expect(body).toContain('Show search query variants');
		expect(body).toContain('Override');
	});

	it('renders collapsed release search status messages and placeholders', () => {
		const empty = render(ReleaseSearchStatusLog, { props: { messages: [] } });
		expect(empty.body).toContain('Press search to start');

		const { body } = render(ReleaseSearchStatusLog, {
			props: {
				messages: [
					{
						id: 'log-1',
						timestamp: '05:00:00',
						message: 'Searching Local Torznab',
						resultMessage: '2 releases',
						durationMs: 42,
						cacheHit: true
					}
				]
			}
		});

		expect(body).toContain('[05:00:00]');
		expect(body).toContain('Searching Local Torznab');
		expect(body).toContain('2 releases');
		expect(body).toContain('cache hit');
	});
});
