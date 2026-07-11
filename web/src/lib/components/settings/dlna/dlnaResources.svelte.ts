import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
import { restartDLNA, updateDLNASettings } from './api';
import { createDLNAStatusQuery } from './dlnaStatus.svelte';
import {
	cloneDLNARendererProfile,
	createDLNARendererProfile,
	deleteDLNARendererDeviceOverride,
	deleteDLNARendererProfile,
	importDLNARendererProfile,
	listDLNARecentDevices,
	listDLNARendererDeviceOverrides,
	listDLNARendererProfiles,
	resetDLNARendererProfile,
	restoreDLNARendererProfiles,
	updateDLNARendererProfile,
	upsertDLNARendererDeviceOverride
} from '$lib/settings/dlnaProfilesApi';
import type {
	DLNARendererDeviceOverrideRequest,
	DLNARendererProfileCloneRequest,
	DLNARendererProfileCreateRequest,
	DLNARendererProfileRequest,
	DLNASettingsRequest
} from '$lib/settings/types';

export const dlnaKeys = {
	all: ['settings', 'dlna'] as const,
	settings: () => [...dlnaKeys.all, 'settings'] as const,
	profiles: () => [...dlnaKeys.all, 'profiles'] as const,
	overrides: () => [...dlnaKeys.all, 'overrides'] as const,
	devices: () => [...dlnaKeys.all, 'devices'] as const
};

export function createDLNAResources() {
	const client = useQueryClient();
	const invalidateSettings = () => client.invalidateQueries({ queryKey: dlnaKeys.settings() });
	const invalidateProfiles = () => client.invalidateQueries({ queryKey: dlnaKeys.profiles() });
	const invalidateOverrides = () => client.invalidateQueries({ queryKey: dlnaKeys.overrides() });
	return {
		settings: createDLNAStatusQuery(),
		profiles: createQuery(() => ({
			queryKey: dlnaKeys.profiles(),
			queryFn: listDLNARendererProfiles
		})),
		overrides: createQuery(() => ({
			queryKey: dlnaKeys.overrides(),
			queryFn: listDLNARendererDeviceOverrides
		})),
		devices: createQuery(() => ({ queryKey: dlnaKeys.devices(), queryFn: listDLNARecentDevices })),
		updateSettings: createMutation(() => ({
			mutationFn: (request: DLNASettingsRequest) => updateDLNASettings(request),
			onSuccess: (data) => client.setQueryData(dlnaKeys.settings(), data)
		})),
		restart: createMutation(() => ({ mutationFn: restartDLNA, onSuccess: invalidateSettings })),
		createProfile: createMutation(() => ({
			mutationFn: (request: DLNARendererProfileCreateRequest) => createDLNARendererProfile(request),
			onSuccess: invalidateProfiles
		})),
		importProfile: createMutation(() => ({
			mutationFn: (request: DLNARendererProfileCreateRequest) => importDLNARendererProfile(request),
			onSuccess: invalidateProfiles
		})),
		updateProfile: createMutation(() => ({
			mutationFn: ({ id, request }: { id: string; request: DLNARendererProfileRequest }) =>
				updateDLNARendererProfile(id, request),
			onSuccess: invalidateProfiles
		})),
		cloneProfile: createMutation(() => ({
			mutationFn: ({ id, request }: { id: string; request: DLNARendererProfileCloneRequest }) =>
				cloneDLNARendererProfile(id, request),
			onSuccess: invalidateProfiles
		})),
		resetProfile: createMutation(() => ({
			mutationFn: (id: string) => resetDLNARendererProfile(id),
			onSuccess: invalidateProfiles
		})),
		restoreProfiles: createMutation(() => ({
			mutationFn: restoreDLNARendererProfiles,
			onSuccess: invalidateProfiles
		})),
		deleteProfile: createMutation(() => ({
			mutationFn: (id: string) => deleteDLNARendererProfile(id),
			onSuccess: invalidateProfiles
		})),
		upsertOverride: createMutation(() => ({
			mutationFn: (request: DLNARendererDeviceOverrideRequest) =>
				upsertDLNARendererDeviceOverride(request),
			onSuccess: invalidateOverrides
		})),
		deleteOverride: createMutation(() => ({
			mutationFn: (id: string) => deleteDLNARendererDeviceOverride(id),
			onSuccess: invalidateOverrides
		}))
	};
}
