<script lang="ts">
	import { onMount } from 'svelte';

	import MediaProfileCustomFormatScores from '$lib/components/settings/MediaProfileCustomFormatScores.svelte';
	import MediaProfileQualitySelector from '$lib/components/settings/MediaProfileQualitySelector.svelte';
	import MediaProfileRules from '$lib/components/settings/MediaProfileRules.svelte';
	import MediaProfileTable from '$lib/components/settings/MediaProfileTable.svelte';
	import SettingsAddButton from '$lib/components/settings/shared/SettingsAddButton.svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
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

<Card class="p-5" aria-label="Profiles">
	<SectionHeading>
		{#snippet actions()}
			<SettingsAddButton label="Add profile" onclick={openModal} />
		{/snippet}
	</SectionHeading>

	<MediaProfileTable {profiles} {qualities} {deletingId} onEdit={editProfile} {onDelete} />

	{#if modalOpen}
		<SettingsFormModal
			title={form.id ? 'Edit profile' : 'Add profile'}
			modalClass="w-[min(2560px,calc(100vw-48px))] max-h-[min(880px,calc(100vh-48px))] max-[640px]:w-full max-[640px]:max-h-[calc(100vh-24px)]"
			onClose={closeModal}
		>
			<form class="grid gap-4" onsubmit={saveProfile}>
				<div
					class="grid min-w-0 items-start gap-4.5 min-[981px]:grid-cols-[minmax(340px,0.82fr)_minmax(440px,1fr)]"
				>
					<div class="grid min-w-0 gap-3.5">
						<div class="grid gap-1.5">
							<Label>Name</Label>
							<Input bind:value={form.name} type="text" maxlength={200} required />
						</div>

						<MediaProfileRules {form} {qualities} onChange={updateForm} />
						<MediaProfileCustomFormatScores {form} {customFormats} onChange={updateForm} />
					</div>

					<aside class="grid min-w-0 gap-3.5 min-[981px]:sticky min-[981px]:top-0">
						<MediaProfileQualitySelector
							{form}
							{qualities}
							loading={loadingQualities}
							error={qualityError}
							onChange={updateForm}
						/>
					</aside>
				</div>

				<div class="flex items-center justify-end gap-3">
					<Button type="button" variant="outline" onclick={closeModal}>Cancel</Button>
					<Button type="submit" disabled={saving || !canSave}>
						{saving ? 'Saving' : form.id ? 'Update profile' : 'Create profile'}
					</Button>
				</div>
			</form>
		</SettingsFormModal>
	{/if}
</Card>
