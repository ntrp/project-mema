import { describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
	createQuery: vi.fn((options: () => unknown) => options()),
	searchMedia: vi.fn(),
	getDetails: vi.fn(),
	searchSubtitles: vi.fn()
}));
vi.mock('@tanstack/svelte-query', () => ({ createQuery: mocks.createQuery }));
vi.mock('$lib/settings/api', () => ({
	searchMedia: mocks.searchMedia,
	getMediaMetadataDetails: mocks.getDetails
}));
vi.mock('$lib/features/releases/api', () => ({ searchMediaSubtitles: mocks.searchSubtitles }));

import {
	createMediaLookupQuery,
	createSubtitleSearchQuery,
	createTmdbSeriesDetailsQuery,
	mediaSearchKeys
} from './searchQueries.svelte';

describe('media search queries', () => {
	it('enables lookups only for active, meaningful input', async () => {
		createMediaLookupQuery(
			'movie',
			() => 'a',
			() => true
		);
		const disabled = captured(0);
		expect(disabled.enabled).toBe(false);
		createMediaLookupQuery(
			'movie',
			() => ' dune ',
			() => true
		);
		const active = captured(1);
		expect(active.queryKey).toEqual(mediaSearchKeys.lookup('movie', 'dune'));
		await active.queryFn();
		expect(mocks.searchMedia).toHaveBeenCalledWith({ query: 'dune', type: 'movie' });
	});

	it('keys details by provider, type and identity', async () => {
		createTmdbSeriesDetailsQuery(() => '42');
		const query = captured(2);
		expect(query.enabled).toBe(true);
		expect(query.queryKey).toEqual(mediaSearchKeys.details('tmdb', 'serie', '42'));
		await query.queryFn();
		expect(mocks.getDetails).toHaveBeenCalledWith('tmdb', 'serie', '42');
	});

	it('keeps manual subtitle searches disabled until explicitly submitted', async () => {
		const request = { query: 'Movie 2024', languageId: 'en', filePath: '/movie.mkv' };
		createSubtitleSearchQuery(
			() => 'media-1',
			() => request
		);
		const query = captured(3);
		expect(query.enabled).toBe(false);
		expect(query.queryKey).toEqual(mediaSearchKeys.subtitles('media-1', request));
		await query.queryFn();
		expect(mocks.searchSubtitles).toHaveBeenCalledWith('media-1', request);
	});
});

function captured(index: number) {
	return mocks.createQuery.mock.results[index].value as {
		queryKey: readonly unknown[];
		enabled: boolean;
		queryFn: () => Promise<unknown>;
	};
}
