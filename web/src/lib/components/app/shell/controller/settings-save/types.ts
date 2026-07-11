import type { AppShellState } from '../state.svelte';
import type { RunCommandMutation } from '$lib/app/query/commandMutation.svelte';

export interface SettingsSaveDeps {
	clearNotice: () => void;
	runMutation?: RunCommandMutation;
	loadSettings: () => Promise<void>;
	loadMediaItems?: () => Promise<void>;
	mediaItems?: () => import('$lib/settings/types').MediaItem[];
	users?: () => import('$lib/settings/types').ManagedUser[];
	upsertLibraryFolder?: (_item: import('$lib/settings/types').LibraryFolder) => void;
	upsertPathMapping?: (_item: import('$lib/settings/types').PathMapping) => void;
	upsertLibraryScan?: (_scan: import('$lib/settings/types').LibraryScan) => void;
}

export interface SettingsSaveContext extends SettingsSaveDeps {
	state: AppShellState;
}
