import { createCatalogSaveActions } from './settings-save/catalogSaveActions';
import { createIntegrationSaveActions } from './settings-save/integrationSaveActions';
import { createLibraryUserSaveActions } from './settings-save/libraryUserSaveActions';
import type { SettingsSaveDeps } from './settings-save/types';
import type { AppShellState } from './state.svelte';

export function createSettingsSaveActions(state: AppShellState, deps: SettingsSaveDeps) {
	const context = { state, ...deps };

	return {
		...createIntegrationSaveActions(context),
		...createLibraryUserSaveActions(context),
		...createCatalogSaveActions(context)
	};
}
