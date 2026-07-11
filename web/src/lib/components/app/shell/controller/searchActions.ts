import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import type { MediaAdvancedSearchRequest, MediaSearchResult } from '$lib/settings/types';
import type { AppShellState } from './state.svelte';

interface SearchDeps {
	clearNotice: () => void;
	setAutocompleteQuery: (_query: string) => void;
	setAdvancedRequest: (_request: MediaAdvancedSearchRequest) => void;
}

export function createSearchActions(state: AppShellState, deps: SearchDeps) {
	const clearNotice = deps.clearNotice;
	function autocompleteMedia(query: string) {
		const trimmed = query.trim();
		deps.setAutocompleteQuery(trimmed.length >= 2 ? trimmed : '');
	}

	function advancedSearch(request: MediaAdvancedSearchRequest) {
		clearNotice();
		deps.setAdvancedRequest(request);
	}

	function selectAutocompleteResult(result: MediaSearchResult) {
		if (result.id) {
			void goto(
				resolve(result.type === 'movie' ? '/movies/[id]' : '/series/[id]', { id: result.id })
			);
			return;
		}
		if (result.externalProvider && result.externalId) {
			void goto(
				resolve('/media/[provider]/[type]/[externalId]', {
					provider: result.externalProvider,
					type: result.type,
					externalId: result.externalId
				})
			);
			return;
		}
		state.searchQuery = result.title;
		void goto(resolve(`/search/advanced?q=${encodeURIComponent(result.title)}`));
	}

	function openAdvancedSearch(query: string) {
		state.searchQuery = query;
		void goto(resolve(`/search/advanced?q=${encodeURIComponent(query)}`));
	}

	return { autocompleteMedia, advancedSearch, selectAutocompleteResult, openAdvancedSearch };
}
