import { beforeEach, describe, expect, it, vi } from 'vitest';

const apiMock = vi.hoisted(() => ({
	advancedSearchMedia: vi.fn(),
	autocompleteMedia: vi.fn()
}));

const navigationMock = vi.hoisted(() => ({
	goto: vi.fn(),
	resolve: vi.fn((path: string) => `resolved:${path}`)
}));

vi.mock('$lib/settings/api', () => apiMock);
vi.mock('$app/navigation', () => ({ goto: navigationMock.goto }));
vi.mock('$app/paths', () => ({ resolve: navigationMock.resolve }));

import { createSearchActions } from '../searchActions';
import type { AppShellState } from '../state.svelte';

function shellState(overrides: Record<string, unknown> = {}) {
	return {
		searchQuery: '',
		loadingAutocomplete: false,
		autocompleteGroups: [{ sourceType: 'library', sourceName: 'Library', results: [] }],
		searchingAdvanced: false,
		advancedSearchGroups: [],
		errorMessage: '',
		...overrides
	} as unknown as AppShellState;
}

describe('search actions (SCN-SETTINGS-009)', () => {
	beforeEach(() => {
		apiMock.advancedSearchMedia.mockReset();
		apiMock.autocompleteMedia.mockReset();
		navigationMock.goto.mockReset();
		navigationMock.resolve.mockClear();
	});

	it('ignores short autocomplete queries and applies matching provider results', async () => {
		const state = shellState({ searchQuery: 'Example' });
		const groups = [
			{ sourceType: 'provider', sourceName: 'TMDb', results: [{ title: 'Example' }] }
		];
		apiMock.autocompleteMedia.mockResolvedValue(groups);
		const actions = createSearchActions(state, { clearNotice: vi.fn() });

		await actions.autocompleteMedia(' e ');
		expect(apiMock.autocompleteMedia).not.toHaveBeenCalled();
		expect(state.autocompleteGroups).toEqual([]);

		await actions.autocompleteMedia(' Example ');
		expect(apiMock.autocompleteMedia).toHaveBeenCalledWith('Example', 'library');
		expect(state.autocompleteGroups).toEqual(groups);
		expect(state.loadingAutocomplete).toBe(false);
	});

	it('drops stale autocomplete responses and clears failures', async () => {
		const state = shellState({ searchQuery: 'Different' });
		apiMock.autocompleteMedia.mockResolvedValueOnce([
			{ sourceType: 'provider', sourceName: 'TMDb', results: [{ title: 'Example' }] }
		]);
		const actions = createSearchActions(state, { clearNotice: vi.fn() });

		await actions.autocompleteMedia('Example');
		expect(state.autocompleteGroups).toEqual([]);

		state.searchQuery = 'Example';
		apiMock.autocompleteMedia.mockRejectedValueOnce(new Error('network failed'));
		await actions.autocompleteMedia('Example');
		expect(state.autocompleteGroups).toEqual([]);
		expect(state.loadingAutocomplete).toBe(false);
	});

	it('runs advanced search and reports user-facing failures', async () => {
		const state = shellState();
		const clearNotice = vi.fn();
		const groups = [
			{ sourceType: 'provider', sourceName: 'TMDb', results: [{ title: 'Example' }] }
		];
		apiMock.advancedSearchMedia.mockResolvedValueOnce(groups);
		const actions = createSearchActions(state, { clearNotice });

		await actions.advancedSearch({ query: 'Example', type: 'movie' });
		expect(clearNotice).toHaveBeenCalledOnce();
		expect(state.advancedSearchGroups).toEqual(groups);
		expect(state.searchingAdvanced).toBe(false);

		apiMock.advancedSearchMedia.mockRejectedValueOnce(new Error('Search failed'));
		await actions.advancedSearch({ query: 'Broken' });
		expect(state.errorMessage).toBe('Search failed');
		expect(state.searchingAdvanced).toBe(false);
	});

	it('routes selected autocomplete results by library id, provider id, or fallback query', () => {
		const state = shellState();
		const actions = createSearchActions(state, { clearNotice: vi.fn() });

		actions.selectAutocompleteResult({ id: 'movie-1', type: 'movie', title: 'Library Movie' });
		expect(navigationMock.resolve).toHaveBeenCalledWith('/movies/[id]', { id: 'movie-1' });
		expect(navigationMock.goto).toHaveBeenLastCalledWith('resolved:/movies/[id]');

		actions.selectAutocompleteResult({
			type: 'series',
			title: 'Provider Series',
			externalProvider: 'tmdb',
			externalId: '123'
		});
		expect(navigationMock.resolve).toHaveBeenCalledWith('/media/[provider]/[type]/[externalId]', {
			provider: 'tmdb',
			type: 'series',
			externalId: '123'
		});

		actions.selectAutocompleteResult({ type: 'movie', title: 'Loose Query' });
		expect(state.searchQuery).toBe('Loose Query');
		expect(navigationMock.goto).toHaveBeenLastCalledWith(
			'resolved:/search/advanced?q=Loose%20Query'
		);

		actions.openAdvancedSearch('Manual Search');
		expect(state.searchQuery).toBe('Manual Search');
		expect(navigationMock.goto).toHaveBeenLastCalledWith(
			'resolved:/search/advanced?q=Manual%20Search'
		);
	});
});
