import type { SettingsSection, SystemSection } from '$lib/settings/types';
import {
	loadMediaActionSettings as loadMediaActionSettingsRequest,
	loadSettingsSection as loadSettingsSectionRequest,
	loadSystemSettings as loadSystemSettingsRequest
} from '$lib/features/settings/sectionData';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

export function createLoadActions(state: AppShellState) {
	async function loadSettingsSection(section: SettingsSection) {
		await apply(() => loadSettingsSectionRequest(section), `Could not load ${section} settings`);
	}
	async function loadSystemSettings(section: SystemSection) {
		await apply(() => loadSystemSettingsRequest(section), `Could not load ${section} system data`);
	}
	async function loadMediaActionSettings() {
		await apply(loadMediaActionSettingsRequest, 'Could not load media action settings');
	}
	async function apply(load: () => Promise<object>, message: string) {
		try {
			Object.assign(state, await load());
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, message);
		}
	}

	return { loadSettingsSection, loadSystemSettings, loadMediaActionSettings };
}
