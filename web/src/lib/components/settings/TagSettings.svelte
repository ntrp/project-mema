<script lang="ts">
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
</script>

<div class="panel" aria-labelledby="tag-settings-title">
	<div class="section-heading">
		<div>
			<p class="section-kicker">Organization</p>
			<h2 id="tag-settings-title">Tags</h2>
		</div>
	</div>

	<form class="settings-form compact-form" onsubmit={onSave}>
		<label>
			<span>Name</span>
			<input bind:value={form.name} type="text" maxlength="80" required />
		</label>
		<div class="form-actions">
			<button type="button" class="secondary" onclick={onCancel}>Cancel</button>
			<button type="submit" disabled={saving}>
				{saving ? 'Saving' : form.id ? 'Update tag' : 'Create tag'}
			</button>
		</div>
	</form>

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
						<td>{new Date(tag.updatedAt).toLocaleDateString()}</td>
						<td class="row-actions">
							<button type="button" class="secondary" onclick={() => onEdit(tag)}>Edit</button>
							<button
								type="button"
								class="danger"
								disabled={deletingId === tag.id}
								onclick={() => onDelete(tag.id)}
							>
								{deletingId === tag.id ? 'Deleting' : 'Delete'}
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
</div>
