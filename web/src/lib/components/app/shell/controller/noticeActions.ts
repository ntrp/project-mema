import type { AppShellState } from './state.svelte';

export function createNoticeActions(state: AppShellState) {
	function clearNotice() {
		state.errorMessage = '';
		state.message = '';
	}

	return { clearNotice };
}
