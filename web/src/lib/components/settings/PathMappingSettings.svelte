<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/SettingsFormModal.svelte';
	import type { PathMapping, PathMappingForm } from '$lib/settings/types';

	interface Props {
		mappings: PathMapping[];
		form: PathMappingForm;
		saving: boolean;
		deletingId?: string;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let { mappings, form = $bindable(), saving, deletingId, onSave, onDelete }: Props = $props();
	let modalOpen = $state(false);

	async function save(event: SubmitEvent) {
		await onSave(event);
		if (form.clientPath.trim() === '' && form.appPath.trim() === '') {
			modalOpen = false;
		}
	}
</script>

<section class="settings-panel" aria-labelledby="path-mapping-title">
	<div class="settings-panel-header">
		<div>
			<h2 id="path-mapping-title">Path mappings</h2>
			<p>Map download client paths to paths visible by the app for hardlink imports.</p>
		</div>
		<button type="button" onclick={() => (modalOpen = true)}>Add path</button>
	</div>
	<div class="data-list compact-list">
		{#each mappings as mapping (mapping.id)}
			<div class="data-row">
				<div>
					<strong>{mapping.clientPath}</strong>
					<span>{mapping.appPath}</span>
				</div>
				<button
					type="button"
					class="danger"
					disabled={deletingId === mapping.id}
					onclick={() => onDelete(mapping.id)}
				>
					{deletingId === mapping.id ? 'Removing' : 'Delete'}
				</button>
			</div>
		{:else}
			<p class="empty">No paths have been defined.</p>
		{/each}
	</div>
	{#if modalOpen}
		<SettingsFormModal title="Add path mapping" onClose={() => (modalOpen = false)}>
			<form class="settings-grid compact" onsubmit={save}>
				<label>
					<span>Client path</span>
					<input bind:value={form.clientPath} placeholder="/downloads" required />
				</label>
				<label>
					<span>App path</span>
					<input bind:value={form.appPath} placeholder="/mnt/downloads" required />
				</label>
				<div class="form-actions inline">
					<button type="submit" disabled={saving}>
						{saving ? 'Saving' : 'Save mapping'}
					</button>
				</div>
			</form>
		</SettingsFormModal>
	{/if}
</section>
