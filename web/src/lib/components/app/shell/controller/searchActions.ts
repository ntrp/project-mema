import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import {
	advancedSearchMedia as advancedSearchMediaRequest,
	autocompleteMedia as autocompleteMediaRequest
} from '$lib/settings/api';
import type { MediaAdvancedSearchRequest, MediaSearchResult } from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

interface SearchDeps {
	clearNotice: () => void;
}

export function createSearchActions(state: AppShellState, deps: SearchDeps) {
	const clearNotice = deps.clearNotice;
	async function autocompleteMedia(query: string) {
		const trimmed = query.trim();
		if (trimmed.length < 2) {
			state.autocompleteGroups = [];
			return;
		}
		state.loadingAutocomplete = true;
		state.autocompleteGroups = [];
		try {
			const groups = await autocompleteMediaRequest(trimmed, 'library');
			if (state.searchQuery.trim() !== trimmed) {
				return;
			}
			state.autocompleteGroups = groups;
		} catch {
			state.autocompleteGroups = [];
		} finally {
			state.loadingAutocomplete = false;
		}
	}

	async function advancedSearch(request: MediaAdvancedSearchRequest) {
		state.searchingAdvanced = true;
		clearNotice();

		try {
			state.advancedSearchGroups = await advancedSearchMediaRequest(request);
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not search media');
		} finally {
			state.searchingAdvanced = false;
		}
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
