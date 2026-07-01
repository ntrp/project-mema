import type { AppShellState } from './state.svelte';

export function createNoticeActions(state: AppShellState) {
	function clearNotice() {
		state.errorMessage = '';
		state.message = '';
	}

	function showProfile() {
		clearNotice();
		state.message = 'Profile settings are not implemented yet';
	}

	return { clearNotice, showProfile };
}
