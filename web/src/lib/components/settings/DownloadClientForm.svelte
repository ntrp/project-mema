<script lang="ts">
	import type { DownloadClientForm, DownloadClientType } from '$lib/settings/types';

	interface Props {
		form: DownloadClientForm;
		saving: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
		onCancel: () => void;
	}

	let { form = $bindable(), saving, onSave, onCancel }: Props = $props();
	const downloadClientTypes: DownloadClientType[] = ['transmission', 'sabnzbd'];
</script>

<div class="panel" aria-labelledby="download-client-form-title">
	<div class="section-heading">
		<h2 id="download-client-form-title">
			{form.id ? 'Edit download client' : 'Add download client'}
		</h2>
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
				{#each downloadClientTypes as type (type)}
					<option value={type}>{type}</option>
				{/each}
			</select>
		</label>
		<label class="wide">
			<span>Base URL</span>
			<input bind:value={form.baseUrl} placeholder="http://host:port" required />
		</label>
		<label>
			<span>Username</span>
			<input bind:value={form.username} autocomplete="off" />
		</label>
		<label>
			<span>Password</span>
			<input bind:value={form.password} autocomplete="off" type="password" />
		</label>
		<label>
			<span>API key</span>
			<input bind:value={form.apiKey} autocomplete="off" />
		</label>
		<label>
			<span>Category</span>
			<input bind:value={form.category} placeholder="movies" />
		</label>
		<label>
			<span>Priority</span>
			<input bind:value={form.priority} min="0" max="1000" type="number" />
		</label>
		<label class="toggle">
			<input bind:checked={form.enabled} type="checkbox" />
			<span>Enabled</span>
		</label>
		<button type="submit" disabled={saving}>{saving ? 'Saving' : 'Save client'}</button>
	</form>
</div>
