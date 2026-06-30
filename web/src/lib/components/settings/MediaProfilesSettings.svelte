<script lang="ts">
	import { onMount } from 'svelte';

	import MediaProfileCustomFormatScores from '$lib/components/settings/MediaProfileCustomFormatScores.svelte';
	import MediaProfileQualitySelector from '$lib/components/settings/MediaProfileQualitySelector.svelte';
	import MediaProfileRules from '$lib/components/settings/MediaProfileRules.svelte';
	import MediaProfileTable from '$lib/components/settings/MediaProfileTable.svelte';
	import SettingsFormModal from '$lib/components/settings/SettingsFormModal.svelte';
	import { listQualitySizeSettings } from '$lib/settings/api';
	import { emptyMediaProfileForm } from '$lib/settings/forms';
	import type {
		CustomFormat,
		MediaProfile,
		MediaProfileForm,
		QualitySizeSetting
	} from '$lib/settings/types';

	interface Props {
		profiles: MediaProfile[];
		customFormats: CustomFormat[];
		form: MediaProfileForm;
		saving: boolean;
		deletingId?: string;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
		onEdit: (_profile: MediaProfile) => void;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let {
		profiles,
		customFormats,
		form = $bindable(),
		saving,
		deletingId,
		onSave,
		onCancel,
		onEdit,
		onDelete
	}: Props = $props();

	let modalOpen = $state(false);
	let qualities = $state<QualitySizeSetting[]>([]);
	let loadingQualities = $state(false);
	let qualityError = $state('');

	let canSave = $derived(form.name.trim() !== '' && form.qualityIds.length > 0);

	onMount(() => {
		void loadQualities();
	});

	function openModal() {
		form = emptyMediaProfileForm();
		qualityError = '';
		modalOpen = true;
	}

	function editProfile(profile: MediaProfile) {
		onEdit(profile);
		qualityError = '';
		modalOpen = true;
	}

	function closeModal() {
		onCancel();
		modalOpen = false;
	}

	async function loadQualities() {
		loadingQualities = true;
		qualityError = '';
		try {
			const response = await listQualitySizeSettings();
			qualities = response.qualities;
		} catch (error) {
			qualityError = error instanceof Error ? error.message : 'Could not load qualities';
		} finally {
			loadingQualities = false;
		}
	}

	async function saveProfile(event: SubmitEvent) {
		await onSave(event);
		if (!form.id && form.name === '' && form.qualityIds.length === 0) {
			modalOpen = false;
		}
	}

	function updateForm(value: MediaProfileForm) {
		form = value;
	}
</script>

<div class="panel" aria-labelledby="profile-settings-title">
	<div class="section-heading">
		<div>
			<p class="section-kicker">Media</p>
			<h2 id="profile-settings-title">Profiles</h2>
		</div>
		<button type="button" onclick={openModal}>Add profile</button>
	</div>

	<MediaProfileTable {profiles} {qualities} {deletingId} onEdit={editProfile} {onDelete} />

	{#if modalOpen}
		<SettingsFormModal
			title={form.id ? 'Edit profile' : 'Add profile'}
			modalClass="profile-settings-modal"
			onClose={closeModal}
		>
			<form class="settings-form profile-form" onsubmit={saveProfile}>
				<div class="profile-editor-grid">
					<div class="profile-editor-main">
						<label>
							<span>Name</span>
							<input bind:value={form.name} type="text" maxlength="200" required />
						</label>

						<MediaProfileRules {form} {qualities} onChange={updateForm} />
						<MediaProfileCustomFormatScores {form} {customFormats} onChange={updateForm} />
					</div>

					<aside class="profile-editor-qualities">
						<MediaProfileQualitySelector
							{form}
							{qualities}
							loading={loadingQualities}
							error={qualityError}
							onChange={updateForm}
						/>
					</aside>
				</div>

				<div class="form-actions wide">
					<button type="button" class="secondary" onclick={closeModal}>Cancel</button>
					<button type="submit" disabled={saving || !canSave}>
						{saving ? 'Saving' : form.id ? 'Update profile' : 'Create profile'}
					</button>
				</div>
			</form>
		</SettingsFormModal>
	{/if}
</div>
