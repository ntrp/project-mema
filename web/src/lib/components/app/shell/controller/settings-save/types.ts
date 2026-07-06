import type { AppShellState } from '../state.svelte';

export interface SettingsSaveDeps {
	clearNotice: () => void;
	loadSettings: () => Promise<void>;
	loadMediaItems?: () => Promise<void>;
}

export interface SettingsSaveContext extends SettingsSaveDeps {
	state: AppShellState;
}
