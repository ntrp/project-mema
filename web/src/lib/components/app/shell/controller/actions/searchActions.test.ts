import { beforeEach, describe, expect, it, vi } from 'vitest';

const navigationMock = vi.hoisted(() => ({
	goto: vi.fn(),
	resolve: vi.fn((path: string) => `resolved:${path}`)
}));

vi.mock('$app/navigation', () => ({ goto: navigationMock.goto }));
vi.mock('$app/paths', () => ({ resolve: navigationMock.resolve }));

import { createSearchActions } from '../searchActions';
import type { AppShellState } from '../state.svelte';

describe('search actions (SCN-SETTINGS-009)', () => {
	beforeEach(() => {
		navigationMock.goto.mockReset();
		navigationMock.resolve.mockClear();
	});

	it('normalizes autocomplete input before handing it to the query owner', () => {
		const dependencies = deps();
		const actions = createSearchActions(shellState(), dependencies);

		actions.autocompleteMedia(' e ');
		expect(dependencies.setAutocompleteQuery).toHaveBeenLastCalledWith('');
		actions.autocompleteMedia(' Example ');
		expect(dependencies.setAutocompleteQuery).toHaveBeenLastCalledWith('Example');
	});

	it('hands advanced requests to the query owner and clears notices', () => {
		const dependencies = deps();
		const request = { query: 'Example', type: 'movie' } as const;
		createSearchActions(shellState(), dependencies).advancedSearch(request);
		expect(dependencies.clearNotice).toHaveBeenCalledOnce();
		expect(dependencies.setAdvancedRequest).toHaveBeenCalledWith(request);
	});

	it('routes selected results by library id, provider id, or fallback query', () => {
		const state = shellState();
		const actions = createSearchActions(state, deps());

		actions.selectAutocompleteResult({ id: 'movie-1', type: 'movie', title: 'Library Movie' });
		expect(navigationMock.resolve).toHaveBeenCalledWith('/movies/[id]', { id: 'movie-1' });

		actions.selectAutocompleteResult({
			type: 'serie',
			title: 'Provider Series',
			externalProvider: 'tmdb',
			externalId: '123'
		});
		expect(navigationMock.resolve).toHaveBeenCalledWith('/media/[provider]/[type]/[externalId]', {
			provider: 'tmdb',
			type: 'serie',
			externalId: '123'
		});

		actions.selectAutocompleteResult({ type: 'movie', title: 'Loose Query' });
		expect(state.searchQuery).toBe('Loose Query');
		expect(navigationMock.goto).toHaveBeenLastCalledWith(
			'resolved:/search/advanced?q=Loose%20Query'
		);

		actions.openAdvancedSearch('Manual Search');
		expect(state.searchQuery).toBe('Manual Search');
	});
});

function shellState() {
	return { searchQuery: '' } as AppShellState;
}

function deps() {
	return {
		clearNotice: vi.fn(),
		setAutocompleteQuery: vi.fn(),
		setAdvancedRequest: vi.fn()
	};
}
