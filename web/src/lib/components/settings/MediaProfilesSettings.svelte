<script lang="ts">
	import { onMount } from 'svelte';

	import SettingsFormModal from '$lib/components/settings/SettingsFormModal.svelte';
	import { listQualitySizeSettings } from '$lib/settings/api';
	import { emptyMediaProfileForm } from '$lib/settings/forms';
	import {
		groupQualitiesByResolution,
		type QualityResolutionGroup
	} from '$lib/settings/qualityGroups';
	import type { MediaProfile, MediaProfileForm, QualitySizeSetting } from '$lib/settings/types';

	interface Props {
		profiles: MediaProfile[];
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

	let selectedQualityNames = $derived.by(() => {
		const names = new Map(qualities.map((quality) => [quality.qualityId, quality.name]));
		return (profile: MediaProfile) =>
			profile.qualityIds.map((id) => names.get(id) ?? id).filter(Boolean);
	});

	let canSave = $derived(form.name.trim() !== '' && form.qualityIds.length > 0);
	let qualityGroups = $derived(groupQualitiesByResolution(qualities));

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

	function toggleQuality(qualityId: string) {
		if (form.qualityIds.includes(qualityId)) {
			form = {
				...form,
				qualityIds: form.qualityIds.filter((id) => id !== qualityId)
			};
			return;
		}
		form = {
			...form,
			qualityIds: [...form.qualityIds, qualityId]
		};
	}

	function selectAllQualities() {
		form = {
			...form,
			qualityIds: qualities.map((quality) => quality.qualityId)
		};
	}

	function clearQualities() {
		form = {
			...form,
			qualityIds: []
		};
	}

	function selectQualityGroup(group: QualityResolutionGroup<QualitySizeSetting>) {
		form = {
			...form,
			qualityIds: [
				...new Set([...form.qualityIds, ...group.qualities.map((quality) => quality.qualityId)])
			]
		};
	}

	function clearQualityGroup(group: QualityResolutionGroup<QualitySizeSetting>) {
		const groupIDs = new Set(group.qualities.map((quality) => quality.qualityId));
		form = {
			...form,
			qualityIds: form.qualityIds.filter((qualityId) => !groupIDs.has(qualityId))
		};
	}

	function selectedQualityCount(group: QualityResolutionGroup<QualitySizeSetting>) {
		return group.qualities.filter((quality) => form.qualityIds.includes(quality.qualityId)).length;
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

	<div class="table-wrap">
		<table>
			<thead>
				<tr>
					<th>Name</th>
					<th>Qualities</th>
					<th>Updated</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each profiles as profile (profile.id)}
					<tr>
						<td>
							<strong>{profile.name}</strong>
						</td>
						<td>
							<div class="quality-chip-list" aria-label={`${profile.name} qualities`}>
								{#each selectedQualityNames(profile).slice(0, 6) as name (name)}
									<span>{name}</span>
								{/each}
								{#if profile.qualityIds.length > 6}
									<span>+{profile.qualityIds.length - 6}</span>
								{/if}
							</div>
						</td>
						<td>{new Date(profile.updatedAt).toLocaleDateString()}</td>
						<td class="row-actions">
							<button type="button" class="secondary" onclick={() => editProfile(profile)}>
								Edit
							</button>
							<button
								type="button"
								class="danger"
								disabled={deletingId === profile.id}
								onclick={() => onDelete(profile.id)}
							>
								{deletingId === profile.id ? 'Deleting' : 'Delete'}
							</button>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="4" class="empty">No profiles configured</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	{#if modalOpen}
		<SettingsFormModal title={form.id ? 'Edit profile' : 'Add profile'} onClose={closeModal}>
			<form class="settings-form profile-form" onsubmit={saveProfile}>
				<label class="wide">
					<span>Name</span>
					<input bind:value={form.name} type="text" maxlength="200" required />
				</label>

				<div class="wide profile-quality-header">
					<span>Qualities</span>
					<div>
						<button type="button" class="secondary" onclick={selectAllQualities}>Select all</button>
						<button type="button" class="secondary" onclick={clearQualities}>Clear</button>
					</div>
				</div>

				{#if qualityError}
					<p class="form-status error wide">{qualityError}</p>
				{/if}

				<div class="quality-group-stack wide" aria-label="Profile qualities">
					{#if loadingQualities}
						<p class="muted">Loading qualities</p>
					{:else}
						{#each qualityGroups as group (group.id)}
							<section class="quality-checkbox-group" aria-labelledby={`quality-group-${group.id}`}>
								<div class="quality-group-heading">
									<div>
										<h3 id={`quality-group-${group.id}`}>{group.label}</h3>
										<span>{selectedQualityCount(group)} / {group.qualities.length}</span>
									</div>
									<div>
										<button
											type="button"
											class="secondary"
											onclick={() => selectQualityGroup(group)}
										>
											Select
										</button>
										<button
											type="button"
											class="secondary"
											onclick={() => clearQualityGroup(group)}
										>
											Clear
										</button>
									</div>
								</div>
								<div class="quality-checkbox-grid">
									{#each group.qualities as quality (quality.qualityId)}
										<label class="quality-checkbox">
											<input
												type="checkbox"
												checked={form.qualityIds.includes(quality.qualityId)}
												onchange={() => toggleQuality(quality.qualityId)}
											/>
											<span>{quality.name}</span>
										</label>
									{/each}
								</div>
							</section>
						{/each}
					{/if}
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
