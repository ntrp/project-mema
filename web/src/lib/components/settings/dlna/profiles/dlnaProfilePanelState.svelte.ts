import { exportDLNARendererProfile } from '$lib/settings/dlnaProfilesApi';
import { createDLNAResources } from '../dlnaResources.svelte';
import type {
	DLNAClientDiagnostic,
	DLNARendererDeviceOverride,
	DLNARendererDeviceOverrideRequest,
	DLNARendererProfile
} from '$lib/settings/types';
import {
	blankProfileForm,
	defaultOverrideRequest,
	downloadProfileJson,
	formToCreateRequest,
	formToRequest,
	importProfileText,
	profileToForm,
	type DLNAProfileForm
} from './dlnaProfileForms';
export class DLNAProfilePanelState {
	private resources = createDLNAResources();
	selectedId = $state('');
	form = $state<DLNAProfileForm>();
	overrideForm = $state<DLNARendererDeviceOverrideRequest>(defaultOverrideRequest());
	overrideJsonText = $state('{}');
	search = $state('');
	errorMessage = $state('');
	message = $state('');
	cloneSource = $state<DLNARendererProfile>();
	cloneId = $state('');
	cloneName = $state('');
	importOpen = $state(false);
	importText = $state('');
	traceIp = $state('');
	traceMediaPath = $state('');

	get profiles(): DLNARendererProfile[] {
		return this.resources.profiles.data ?? [];
	}

	get overrides(): DLNARendererDeviceOverride[] {
		return this.resources.overrides.data ?? [];
	}

	get devices(): DLNAClientDiagnostic[] {
		return this.resources.devices.data ?? [];
	}

	get loading() {
		return (
			this.resources.profiles.isFetching ||
			this.resources.overrides.isFetching ||
			this.resources.devices.isFetching
		);
	}

	get saving() {
		return [
			this.resources.createProfile,
			this.resources.updateProfile,
			this.resources.cloneProfile,
			this.resources.importProfile,
			this.resources.resetProfile,
			this.resources.upsertOverride,
			this.resources.deleteOverride
		].some((mutation) => mutation.isPending);
	}

	get selectedProfile() {
		return this.profiles.find((profile) => profile.id === this.selectedId);
	}

	get filteredProfiles() {
		const query = this.search.toLowerCase();
		return this.profiles.filter((profile) =>
			`${profile.name} ${profile.vendor} ${profile.deviceClass} ${profile.id}`
				.toLowerCase()
				.includes(query)
		);
	}

	load = async () => {
		this.errorMessage = '';
		try {
			await Promise.all([
				this.resources.profiles.refetch(),
				this.resources.overrides.refetch(),
				this.resources.devices.refetch()
			]);
			this.selectProfile(
				this.profiles.find((profile) => profile.id === this.selectedId) ?? this.profiles[0]
			);
			this.overrideForm = defaultOverrideRequest(this.profiles[0]?.id ?? '');
			this.traceIp ||= this.devices[0]?.ip ?? '';
		} catch (error) {
			this.errorMessage = error instanceof Error ? error.message : 'Could not load DLNA profiles';
		}
	};

	selectProfile = (profile?: DLNARendererProfile) => {
		this.selectedId = profile?.id ?? '';
		this.form = profile ? profileToForm(profile) : undefined;
	};

	newProfile = () => {
		const source = this.selectedProfile;
		this.selectedId = '';
		this.form = blankProfileForm(source);
	};

	saveProfile = async () => {
		if (!this.form) return;
		this.errorMessage = '';
		try {
			const current = this.selectedProfile;
			const saved = current
				? await this.resources.updateProfile.mutateAsync({
						id: current.id,
						request: formToRequest(this.form)
					})
				: await this.resources.createProfile.mutateAsync(formToCreateRequest(this.form));
			this.selectProfile(saved);
			this.message = current ? 'Profile saved' : 'Profile created';
		} catch (error) {
			this.errorMessage = error instanceof Error ? error.message : 'Could not save profile';
		}
	};

	resetProfile = async (profile: DLNARendererProfile) => {
		const saved = await this.resources.resetProfile.mutateAsync(profile.id);
		this.selectProfile(saved);
		this.message = 'Profile reset';
	};

	openClone = (profile: DLNARendererProfile) => {
		this.cloneSource = profile;
		this.cloneId = `${profile.id}-copy`;
		this.cloneName = `${profile.name} copy`;
	};

	cloneProfile = async () => {
		if (!this.cloneSource) return;
		const saved = await this.resources.cloneProfile.mutateAsync({
			id: this.cloneSource.id,
			request: { id: this.cloneId, name: this.cloneName }
		});
		this.selectProfile(saved);
		this.cloneSource = undefined;
		this.message = 'Profile cloned';
	};

	importProfile = async () => {
		const saved = await this.resources.importProfile.mutateAsync(
			importProfileText(this.importText)
		);
		this.selectProfile(saved);
		this.importOpen = false;
		this.importText = '';
		this.message = 'Profile imported';
	};

	exportProfile = async (profile: DLNARendererProfile) => {
		downloadProfileJson(await exportDLNARendererProfile(profile.id));
	};

	saveOverride = async () => {
		const parsed = JSON.parse(this.overrideJsonText || '{}') as Record<string, unknown>;
		await this.resources.upsertOverride.mutateAsync({
			...this.overrideForm,
			deliveryPolicyOverrides: parsed
		});
		this.message = 'Override saved';
	};

	quickAssign = async (device: DLNAClientDiagnostic, profileId: string) => {
		const existing = this.overrides.find((override) => override.ipAddress === device.ip);
		if (!profileId) {
			if (existing) await this.deleteOverride(existing.id);
			return;
		}
		await this.resources.upsertOverride.mutateAsync({
			id: existing?.id,
			ipAddress: device.ip,
			profileId,
			displayName: device.userAgent || device.ip,
			allowed: true,
			deliveryPolicyOverrides: {},
			notes: ''
		});
	};

	deleteOverride = async (id: string) => {
		await this.resources.deleteOverride.mutateAsync(id);
	};
}
