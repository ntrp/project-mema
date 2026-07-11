import { describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
	createQuery: vi.fn((options: () => unknown) => options()),
	createMutation: vi.fn((options: () => unknown) => options()),
	setQueryData: vi.fn(),
	invalidateQueries: vi.fn()
}));
vi.mock('@tanstack/svelte-query', () => ({
	createQuery: mocks.createQuery,
	createMutation: mocks.createMutation,
	useQueryClient: () => ({
		setQueryData: mocks.setQueryData,
		invalidateQueries: mocks.invalidateQueries
	})
}));
vi.mock('./api', () => ({
	getDLNASettings: vi.fn(),
	restartDLNA: vi.fn(),
	updateDLNASettings: vi.fn()
}));
vi.mock('$lib/settings/dlnaProfilesApi', () => ({
	cloneDLNARendererProfile: vi.fn(),
	createDLNARendererProfile: vi.fn(),
	deleteDLNARendererDeviceOverride: vi.fn(),
	importDLNARendererProfile: vi.fn(),
	listDLNARecentDevices: vi.fn(),
	listDLNARendererDeviceOverrides: vi.fn(),
	listDLNARendererProfiles: vi.fn(),
	resetDLNARendererProfile: vi.fn(),
	updateDLNARendererProfile: vi.fn(),
	upsertDLNARendererDeviceOverride: vi.fn()
}));

import { createDLNAResources, dlnaKeys } from './dlnaResources.svelte';

describe('DLNA resources', () => {
	it('uses stable keys for every server collection', () => {
		createDLNAResources();
		const options = mocks.createQuery.mock.results.map(
			(result) => result.value as { queryKey: readonly string[] }
		);
		expect(options.map((option) => option.queryKey)).toEqual([
			dlnaKeys.settings(),
			dlnaKeys.profiles(),
			dlnaKeys.overrides(),
			dlnaKeys.devices()
		]);
	});

	it('reconciles settings and invalidates changed collections', () => {
		mocks.createMutation.mockClear();
		createDLNAResources();
		const mutations = mocks.createMutation.mock.results.map(
			(result) => result.value as { onSuccess?: (data?: unknown) => void }
		);
		const settings = { enabled: true };
		mutations[0].onSuccess?.(settings);
		expect(mocks.setQueryData).toHaveBeenCalledWith(dlnaKeys.settings(), settings);
		mutations[1].onSuccess?.();
		mutations[2].onSuccess?.();
		mutations[7].onSuccess?.();
		expect(mocks.invalidateQueries).toHaveBeenCalledWith({ queryKey: dlnaKeys.settings() });
		expect(mocks.invalidateQueries).toHaveBeenCalledWith({ queryKey: dlnaKeys.profiles() });
		expect(mocks.invalidateQueries).toHaveBeenCalledWith({ queryKey: dlnaKeys.overrides() });
	});
});
