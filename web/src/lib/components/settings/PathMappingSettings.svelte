<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
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
		<button type="button" class="add-action-button" onclick={() => (modalOpen = true)}>
			<span class="app-icon" aria-hidden="true">add</span>
			<span>Add path</span>
		</button>
	</div>
	{#if mappings.length > 0}
		<div class="table-wrap">
			<table>
				<thead>
					<tr>
						<th>Client path</th>
						<th>App path</th>
						<th class="table-action-heading">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each mappings as mapping (mapping.id)}
						<tr>
							<td>{mapping.clientPath}</td>
							<td>{mapping.appPath}</td>
							<td class="row-actions">
								<button
									type="button"
									class="danger icon-button"
									disabled={deletingId === mapping.id}
									aria-label="Delete path mapping"
									title="Delete path mapping"
									onclick={() => onDelete(mapping.id)}
								>
									<span class="app-icon" aria-hidden="true">delete</span>
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{:else}
		<p class="empty">No paths have been defined.</p>
	{/if}
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
