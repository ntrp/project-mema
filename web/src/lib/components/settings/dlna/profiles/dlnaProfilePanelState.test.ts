import { describe, expect, it, vi } from 'vitest';
import type { DLNARendererProfile } from '$lib/settings/types';

type MutationMock = {
	isPending: boolean;
	mutateAsync: ReturnType<typeof vi.fn>;
};

type MockResources = {
	profiles: { data: DLNARendererProfile[]; isFetching: boolean; refetch: ReturnType<typeof vi.fn> };
	overrides: { data: unknown[]; isFetching: boolean; refetch: ReturnType<typeof vi.fn> };
	devices: { data: unknown[]; isFetching: boolean; refetch: ReturnType<typeof vi.fn> };
	createProfile: MutationMock;
	updateProfile: MutationMock;
	cloneProfile: MutationMock;
	importProfile: MutationMock;
	resetProfile: MutationMock;
	restoreProfiles: MutationMock;
	deleteProfile: MutationMock;
	upsertOverride: MutationMock;
	deleteOverride: MutationMock;
};

const mocks = vi.hoisted(() => {
	const mutation = (): MutationMock => ({ isPending: false, mutateAsync: vi.fn() });
	const resources: MockResources = {
		profiles: { data: [] as DLNARendererProfile[], isFetching: false, refetch: vi.fn() },
		overrides: { data: [], isFetching: false, refetch: vi.fn() },
		devices: { data: [], isFetching: false, refetch: vi.fn() },
		createProfile: mutation(),
		updateProfile: mutation(),
		cloneProfile: mutation(),
		importProfile: mutation(),
		resetProfile: mutation(),
		restoreProfiles: mutation(),
		deleteProfile: mutation(),
		upsertOverride: mutation(),
		deleteOverride: mutation()
	};

	return {
		resources,
		exportProfile: vi.fn()
	};
});

vi.mock('../dlnaResources.svelte', () => ({
	createDLNAResources: () => mocks.resources
}));
vi.mock('$lib/settings/dlnaProfilesApi', () => ({
	exportDLNARendererProfile: mocks.exportProfile
}));

import { DLNAProfilePanelState } from './dlnaProfilePanelState.svelte';

const sampleProfile: DLNARendererProfile = {
	id: 'lg-webos',
	name: 'LG webOS',
	vendor: 'LG',
	deviceClass: 'MediaRenderer',
	enabled: true,
	priority: 120,
	iconKey: 'tv',
	notes: 'seeded',
	matchRules: { headers: ['LG'] },
	capabilityRules: { containers: ['mp4'] },
	deliverySettings: {},
	dlnaFlags: {},
	subtitleRules: {},
	artworkRules: {},
	metadataRules: {},
	quirks: {},
	source: 'mema_seed',
	sourceVersion: 1,
	customized: false,
	createdAt: '2026-07-08T08:00:00Z',
	updatedAt: '2026-07-08T08:00:00Z'
};

describe('DLNA profile panel state', () => {
	it('opens the profile editor in edit mode for a selected profile', () => {
		mocks.resources.profiles.data = [sampleProfile];
		const state = new DLNAProfilePanelState();

		state.openProfileEditor(sampleProfile);

		expect(state.editorOpen).toBe(true);
		expect(state.editorMode).toBe('edit');
		expect(state.selectedId).toBe(sampleProfile.id);
		expect(state.form?.id).toBe(sampleProfile.id);
	});

	it('opens a create modal from the current selection and closes cleanly', () => {
		mocks.resources.profiles.data = [sampleProfile];
		const state = new DLNAProfilePanelState();

		state.selectProfile(sampleProfile);
		state.newProfile();

		expect(state.editorOpen).toBe(true);
		expect(state.editorMode).toBe('create');
		expect(state.selectedId).toBe(sampleProfile.id);
		expect(state.form?.id).toBe('lg-webos-copy');

		state.closeEditor();

		expect(state.editorOpen).toBe(false);
		expect(state.editorMode).toBe('edit');
		expect(state.selectedId).toBe(sampleProfile.id);
		expect(state.form).toBeUndefined();
	});

	it('opens and closes the decision trace modal', () => {
		mocks.resources.devices.data = [{ ip: '192.168.1.20' }, { ip: '192.168.1.21' }];
		const state = new DLNAProfilePanelState();

		state.traceIp = '192.168.1.15';
		state.traceMediaPath = '/media/movie.mkv';
		state.openTrace();
		expect(state.traceOpen).toBe(true);
		expect(state.traceIp).toBe('192.168.1.20');
		state.closeTrace();
		expect(state.traceOpen).toBe(false);
		expect(state.traceIp).toBe('');
		expect(state.traceMediaPath).toBe('');
	});

	it('creates and updates profiles based on the active editor mode', async () => {
		mocks.resources.profiles.data = [sampleProfile];
		mocks.resources.createProfile.mutateAsync.mockResolvedValue({
			...sampleProfile,
			id: 'sony-bravia'
		});
		mocks.resources.updateProfile.mutateAsync.mockResolvedValue({
			...sampleProfile,
			name: 'LG webOS Plus'
		});
		const state = new DLNAProfilePanelState();

		state.selectProfile(sampleProfile);
		state.newProfile();
		state.form = {
			...state.form!,
			id: 'sony-bravia',
			name: 'Sony Bravia',
			vendor: 'Sony',
			deviceClass: 'MediaRenderer'
		};
		await state.saveProfile();

		expect(mocks.resources.createProfile.mutateAsync).toHaveBeenCalledOnce();
		expect(state.selectedId).toBe('sony-bravia');
		expect(state.editorOpen).toBe(false);

		state.openProfileEditor(sampleProfile);
		state.form = {
			...state.form!,
			name: 'LG webOS Plus'
		};
		await state.saveProfile();

		expect(mocks.resources.updateProfile.mutateAsync).toHaveBeenCalledOnce();
		expect(state.selectedId).toBe(sampleProfile.id);
		expect(state.editorOpen).toBe(false);
	});

	it('deletes profiles and restores seeded originals', async () => {
		mocks.resources.profiles.data = [
			sampleProfile,
			{ ...sampleProfile, id: 'user-profile', source: 'user', customized: true }
		];
		mocks.resources.deleteProfile.mutateAsync.mockResolvedValue(undefined);
		mocks.resources.restoreProfiles.mutateAsync.mockResolvedValue(undefined);
		const state = new DLNAProfilePanelState();

		state.selectProfile(sampleProfile);
		await state.deleteProfile({
			...sampleProfile,
			id: 'user-profile',
			source: 'user',
			customized: true
		});
		expect(mocks.resources.deleteProfile.mutateAsync).toHaveBeenCalledWith('user-profile');
		expect(state.message).toBe('Profile deleted');

		await state.restoreOriginalProfiles();
		expect(mocks.resources.restoreProfiles.mutateAsync).toHaveBeenCalledOnce();
		expect(state.message).toBe('Original profiles restored');
	});
});
