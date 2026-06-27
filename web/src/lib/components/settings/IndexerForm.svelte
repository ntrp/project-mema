<script lang="ts">
	import type { IndexerForm, IndexerType } from '$lib/settings/types';

	interface Props {
		form: IndexerForm;
		saving: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
	}

	let { form = $bindable(), saving, onSave, onCancel }: Props = $props();
	const indexerTypes: IndexerType[] = ['torznab', 'newznab', 'rss'];
</script>

<div class="panel" aria-labelledby="indexer-form-title">
	<div class="section-heading">
		<h2 id="indexer-form-title">{form.id ? 'Edit indexer' : 'Add indexer'}</h2>
		{#if form.id}
			<button type="button" class="secondary" onclick={onCancel}>Cancel</button>
		{/if}
	</div>

	<form class="settings-form" onsubmit={onSave}>
		<label>
			<span>Name</span>
			<input bind:value={form.name} required maxlength="200" />
		</label>
		<label>
			<span>Type</span>
			<select bind:value={form.type}>
				{#each indexerTypes as type (type)}
					<option value={type}>{type}</option>
				{/each}
			</select>
		</label>
		<label class="wide">
			<span>Base URL</span>
			<input bind:value={form.baseUrl} placeholder="https://indexer.example" required />
		</label>
		<label>
			<span>API key</span>
			<input bind:value={form.apiKey} autocomplete="off" />
		</label>
		<label>
			<span>Categories</span>
			<input bind:value={form.categoriesText} placeholder="2000, 5000" />
		</label>
		<label>
			<span>Priority</span>
			<input bind:value={form.priority} min="0" max="1000" type="number" />
		</label>
		<label class="toggle">
			<input bind:checked={form.enabled} type="checkbox" />
			<span>Enabled</span>
		</label>
		<button type="submit" disabled={saving}>{saving ? 'Saving' : 'Save indexer'}</button>
	</form>
</div>
