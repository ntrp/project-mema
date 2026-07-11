import type { SettingsData, SettingsSection, SystemSection } from '$lib/settings/types';

type SettingsPatch = Partial<SettingsData>;

export async function loadSettingsSection(_section: SettingsSection): Promise<SettingsPatch> {
	return {};
}

export async function loadSystemSettings(_section: SystemSection): Promise<SettingsPatch> {
	return {};
}

export async function loadMediaActionSettings(): Promise<SettingsPatch> {
	return {};
}
