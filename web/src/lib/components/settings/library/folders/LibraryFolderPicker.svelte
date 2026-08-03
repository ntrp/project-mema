<script lang="ts">
	import { onMount } from 'svelte';
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { createFolderBrowser } from './folderBrowser.svelte';
	import LibraryFolderCreateForm from './LibraryFolderCreateForm.svelte';
	import LibraryFolderTree from './LibraryFolderTree.svelte';
	import type { LibraryFolderPickerProps } from './libraryFolderTree';

	let { initialPath, onClose, onUse }: LibraryFolderPickerProps = $props();
	// svelte-ignore state_referenced_locally
	const browser = createFolderBrowser(initialPath);

	onMount(() => {
		void browser.initialize();
	});
</script>

<SettingsFormModal title="Select library folder" modalClass="grid gap-4" {onClose}>
	<div class="grid gap-4">
		<Card class="grid gap-2 p-3">
			<span class="text-xs font-black text-muted-foreground uppercase">Selected</span>
			<strong class="break-all text-sm text-foreground"
				>{browser.selectedPath || 'No folder selected'}</strong
			>
		</Card>

		{#if browser.error}
			<p class="m-0 text-sm leading-6 text-muted-foreground">{browser.error}</p>
		{:else if browser.loading}
			<p class="m-0 text-sm leading-6 text-muted-foreground">Loading folders</p>
		{:else}
			<Button
				type="button"
				disabled={!browser.selectedPath}
				onclick={() => onUse(browser.selectedPath)}
			>
				Use selected folder
			</Button>

			<LibraryFolderCreateForm
				bind:name={browser.newFolderName}
				disabled={!browser.selectedPath}
				creating={browser.creatingFolder}
				onCreate={() => void browser.createDirectory()}
			/>
			{#if browser.createFolderError}
				<p class="m-0 text-sm font-semibold text-destructive">{browser.createFolderError}</p>
			{/if}

			<LibraryFolderTree
				visibleNodes={browser.visibleNodes}
				selectedPath={browser.selectedPath}
				onSelect={(path) => (browser.selectedPath = path)}
				onToggle={(path) => void browser.toggleNode(path)}
			/>
		{/if}
	</div>
</SettingsFormModal>
