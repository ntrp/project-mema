<script lang="ts">
	import PlusIcon from '@lucide/svelte/icons/plus';
	import LibraryFolderPicker from '$lib/components/settings/library/LibraryFolderPicker.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
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

<Card.Root aria-labelledby="library-folder-form-title">
	<Card.Header>
		<Card.Title id="library-folder-form-title">Add library folder</Card.Title>
	</Card.Header>

	<Card.Content>
		<form class="grid gap-4" onsubmit={onSave}>
			<div class="space-y-2">
				<Label for="library-folder-path">Folder path</Label>
				<Input
					id="library-folder-path"
					bind:value={form.path}
					placeholder="/data/library or downloads/complete"
					required
					maxlength={4000}
				/>
			</div>
			<div class="flex flex-wrap justify-end gap-2">
				<Button type="button" variant="outline" onclick={openPicker}>Browse</Button>
				<Button type="submit" disabled={saving}>
					<PlusIcon aria-hidden="true" />
					<span>{saving ? 'Scanning' : 'Add and scan'}</span>
				</Button>
			</div>
		</form>
	</Card.Content>

	{#if pickerOpen}
		<LibraryFolderPicker initialPath={form.path} onClose={closePicker} onUse={useFolder} />
	{/if}
</Card.Root>
