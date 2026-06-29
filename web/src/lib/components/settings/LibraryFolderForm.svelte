<script lang="ts">
	import LibraryFolderPicker from '$lib/components/settings/LibraryFolderPicker.svelte';
	import type { LibraryFolderForm } from '$lib/settings/types';

	interface Props {
		form: LibraryFolderForm;
		saving: boolean;
		onSave: (_event: SubmitEvent) => void | Promise<void>;
	}

	let { form = $bindable(), saving, onSave }: Props = $props();
	let pickerOpen = $state(false);

	function openPicker() {
		pickerOpen = true;
	}

	function closePicker() {
		pickerOpen = false;
	}

	function useFolder(path: string) {
		form.path = path;
		closePicker();
	}
</script>

<div class="panel" aria-labelledby="library-folder-form-title">
	<div class="section-heading">
		<h2 id="library-folder-form-title">Add library folder</h2>
	</div>

	<form class="settings-form" onsubmit={onSave}>
		<label class="wide">
			<span>Folder path</span>
			<input
				bind:value={form.path}
				placeholder="/data/library or downloads/complete"
				required
				maxlength="4000"
			/>
		</label>
		<div class="form-actions library-folder-actions">
			<button type="button" class="secondary" onclick={openPicker}>Browse</button>
			<button type="submit" disabled={saving}>{saving ? 'Scanning' : 'Add and scan'}</button>
		</div>
	</form>

	{#if pickerOpen}
		<LibraryFolderPicker initialPath={form.path} onClose={closePicker} onUse={useFolder} />
	{/if}
</div>
