import { describe, expect, it } from 'vitest';
import { createResourceEnablement } from './resourceEnablement';
import type { AppShellState } from './state.svelte';

describe('root resource enablement', () => {
	it('does not enable settings, system, profile, or library support data on discovery home', () => {
		const enabled = createResourceEnablement(state());
		expect(enabled.languages()).toBe(false);
		expect(enabled.tags()).toBe(false);
		expect(enabled.users()).toBe(false);
		expect(enabled.downloadClients()).toBe(false);
		expect(enabled.indexerSearch()).toBe(false);
		expect(enabled.metadataCache()).toBe(false);
		expect(enabled.profile()).toBe(false);
		expect(enabled.discovery()).toBe(true);
	});

	it('enables only resources consumed by the active settings section', () => {
		const value = state({ activeView: 'settings', activeSettingsSection: 'library' });
		const enabled = createResourceEnablement(value);
		expect(enabled.libraryFolders()).toBe(true);
		expect(enabled.pathMappings()).toBe(true);
		expect(enabled.mediaProfiles()).toBe(true);
		expect(enabled.metadataProviders()).toBe(true);
		expect(enabled.downloadClients()).toBe(false);
		expect(enabled.indexers()).toBe(false);
	});

	it('enables profile and system caches only on their owning routes', () => {
		expect(createResourceEnablement(state({ activeView: 'profile' })).profile()).toBe(true);
		const indexing = createResourceEnablement(
			state({ activeView: 'system', activeSystemSection: 'indexing' })
		);
		expect(indexing.indexerSearch()).toBe(true);
		expect(indexing.metadataCache()).toBe(false);
	});
});

function state(overrides: Partial<AppShellState> = {}) {
	return {
		authenticated: true,
		isAdmin: true,
		activeView: 'home',
		activeHomeSection: 'discover',
		activeSettingsSection: 'general',
		activeSystemSection: 'status',
		...overrides
	} as AppShellState;
}
