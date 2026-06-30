<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/SettingsFormModal.svelte';
	import { formatDate } from '$lib/settings/dateFormat';
	import type { Tag, TagForm } from '$lib/settings/types';

	interface Props {
		tags: Tag[];
		form: TagForm;
		saving: boolean;
		deletingId?: string;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
		onEdit: (_tag: Tag) => void;
		onDelete: (_id: string) => void | Promise<void>;
	}

	let {
		tags,
		form = $bindable(),
		saving,
		deletingId,
		onSave,
		onCancel,
		onEdit,
		onDelete
	}: Props = $props();

	let tagModalOpen = $state(false);

	function openTagModal() {
		onCancel();
		tagModalOpen = true;
	}

	function editTag(tag: Tag) {
		onEdit(tag);
		tagModalOpen = true;
	}

	function closeTagModal() {
		onCancel();
		tagModalOpen = false;
	}

	async function saveTag(event: SubmitEvent) {
		await onSave(event);
		if (!form.id && form.name === '') {
			tagModalOpen = false;
		}
	}
</script>

<div class="panel" aria-label="Tags">
	<div class="section-heading">
		<button type="button" class="add-action-button" onclick={openTagModal}>
			<span class="app-icon" aria-hidden="true">add</span>
			<span>Add tag</span>
		</button>
	</div>

	<div class="table-wrap">
		<table>
			<thead>
				<tr>
					<th>Name</th>
					<th>Updated</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				{#each tags as tag (tag.id)}
					<tr>
						<td>
							<span class="tag-pill">{tag.name}</span>
						</td>
						<td>{formatDate(tag.updatedAt)}</td>
						<td class="row-actions">
							<button
								type="button"
								class="secondary icon-button"
								aria-label={`Edit ${tag.name}`}
								onclick={() => editTag(tag)}
							>
								<span class="app-icon" aria-hidden="true">edit</span>
							</button>
							<button
								type="button"
								class="danger icon-button"
								disabled={deletingId === tag.id}
								aria-label={`${deletingId === tag.id ? 'Deleting' : 'Delete'} ${tag.name}`}
								onclick={() => onDelete(tag.id)}
							>
								<span class="app-icon" aria-hidden="true">delete</span>
							</button>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="3" class="empty">No tags configured</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	{#if tagModalOpen}
		<SettingsFormModal title={form.id ? 'Edit tag' : 'Add tag'} onClose={closeTagModal}>
			<form class="settings-form compact-form" onsubmit={saveTag}>
				<label>
					<span>Name</span>
					<input bind:value={form.name} type="text" maxlength="80" required />
				</label>
				<div class="form-actions">
					<button type="button" class="secondary" onclick={closeTagModal}>Cancel</button>
					<button type="submit" disabled={saving}>
						{saving ? 'Saving' : form.id ? 'Update tag' : 'Create tag'}
					</button>
				</div>
			</form>
		</SettingsFormModal>
	{/if}
</div>
