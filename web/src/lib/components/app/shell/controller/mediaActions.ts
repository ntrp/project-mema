import { createMediaCreateActions } from './media-actions/createActions';
import { createMediaLibraryActions } from './media-actions/libraryActions';
import type { MediaDeps } from './media-actions/types';
import type { AppShellState } from './state.svelte';

export function createMediaActions(state: AppShellState, deps: MediaDeps) {
	return {
		...createMediaCreateActions(state, deps),
		...createMediaLibraryActions(state, deps)
	};
}
