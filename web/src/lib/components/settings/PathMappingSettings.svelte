<script lang="ts">
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
</script>

<section class="settings-panel" aria-labelledby="path-mapping-title">
	<div class="settings-panel-header">
		<div>
			<h2 id="path-mapping-title">Path mappings</h2>
			<p>Map download client paths to paths visible by the app for hardlink imports.</p>
		</div>
	</div>
	<form class="settings-grid compact" onsubmit={onSave}>
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
			<p class="empty">No path mappings configured.</p>
		{/each}
	</div>
</section>
