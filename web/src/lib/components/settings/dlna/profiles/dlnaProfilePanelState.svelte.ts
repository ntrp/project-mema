import {
	cloneDLNARendererProfile,
	createDLNARendererProfile,
	deleteDLNARendererDeviceOverride,
	exportDLNARendererProfile,
	importDLNARendererProfile,
	listDLNARecentDevices,
	listDLNARendererDeviceOverrides,
	listDLNARendererProfiles,
	resetDLNARendererProfile,
	updateDLNARendererProfile,
	upsertDLNARendererDeviceOverride
} from '$lib/settings/dlnaProfilesApi';
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
	profiles = $state<DLNARendererProfile[]>([]);
	overrides = $state<DLNARendererDeviceOverride[]>([]);
	devices = $state<DLNAClientDiagnostic[]>([]);
	selectedId = $state('');
	form = $state<DLNAProfileForm>();
	overrideForm = $state<DLNARendererDeviceOverrideRequest>(defaultOverrideRequest());
	overrideJsonText = $state('{}');
	search = $state('');
	loading = $state(true);
	saving = $state(false);
	errorMessage = $state('');
	message = $state('');
	cloneSource = $state<DLNARendererProfile>();
	cloneId = $state('');
	cloneName = $state('');
	importOpen = $state(false);
	importText = $state('');
	traceIp = $state('');
	traceMediaPath = $state('');

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
		this.loading = true;
		this.errorMessage = '';
		try {
			[this.profiles, this.overrides, this.devices] = await Promise.all([
				listDLNARendererProfiles(),
				listDLNARendererDeviceOverrides(),
				listDLNARecentDevices()
			]);
			this.selectProfile(
				this.profiles.find((profile) => profile.id === this.selectedId) ?? this.profiles[0]
			);
			this.overrideForm = defaultOverrideRequest(this.profiles[0]?.id ?? '');
			this.traceIp ||= this.devices[0]?.ip ?? '';
		} catch (error) {
			this.errorMessage = error instanceof Error ? error.message : 'Could not load DLNA profiles';
		} finally {
			this.loading = false;
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
		this.saving = true;
		this.errorMessage = '';
		try {
			const current = this.selectedProfile;
			const saved = current
				? await updateDLNARendererProfile(current.id, formToRequest(this.form))
				: await createDLNARendererProfile(formToCreateRequest(this.form));
			this.upsertProfile(saved);
			this.selectProfile(saved);
			this.message = current ? 'Profile saved' : 'Profile created';
		} catch (error) {
			this.errorMessage = error instanceof Error ? error.message : 'Could not save profile';
		} finally {
			this.saving = false;
		}
	};

	resetProfile = async (profile: DLNARendererProfile) => {
		const saved = await resetDLNARendererProfile(profile.id);
		this.upsertProfile(saved);
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
		this.saving = true;
		try {
			const saved = await cloneDLNARendererProfile(this.cloneSource.id, {
				id: this.cloneId,
				name: this.cloneName
			});
			this.upsertProfile(saved);
			this.selectProfile(saved);
			this.cloneSource = undefined;
			this.message = 'Profile cloned';
		} finally {
			this.saving = false;
		}
	};

	importProfile = async () => {
		this.saving = true;
		try {
			const saved = await importDLNARendererProfile(importProfileText(this.importText));
			this.upsertProfile(saved);
			this.selectProfile(saved);
			this.importOpen = false;
			this.importText = '';
			this.message = 'Profile imported';
		} finally {
			this.saving = false;
		}
	};

	exportProfile = async (profile: DLNARendererProfile) => {
		downloadProfileJson(await exportDLNARendererProfile(profile.id));
	};

	saveOverride = async () => {
		const parsed = JSON.parse(this.overrideJsonText || '{}') as Record<string, unknown>;
		const saved = await upsertDLNARendererDeviceOverride({
			...this.overrideForm,
			deliveryPolicyOverrides: parsed
		});
		this.overrides = [saved, ...this.overrides.filter((override) => override.id !== saved.id)];
		this.message = 'Override saved';
	};

	quickAssign = async (device: DLNAClientDiagnostic, profileId: string) => {
		const existing = this.overrides.find((override) => override.ipAddress === device.ip);
		if (!profileId) {
			if (existing) await this.deleteOverride(existing.id);
			return;
		}
		const saved = await upsertDLNARendererDeviceOverride({
			id: existing?.id,
			ipAddress: device.ip,
			profileId,
			displayName: device.userAgent || device.ip,
			allowed: true,
			deliveryPolicyOverrides: {},
			notes: ''
		});
		this.overrides = [saved, ...this.overrides.filter((override) => override.id !== saved.id)];
	};

	deleteOverride = async (id: string) => {
		await deleteDLNARendererDeviceOverride(id);
		this.overrides = this.overrides.filter((override) => override.id !== id);
	};

	upsertProfile(profile: DLNARendererProfile) {
		this.profiles = [profile, ...this.profiles.filter((item) => item.id !== profile.id)];
	}
}
