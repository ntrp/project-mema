<script lang="ts">
	import SettingsSelect from '$lib/components/settings/shared/SettingsSelect.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Label } from '$lib/components/ui/label';
	import { getFileNamingSettings } from '$lib/settings/api';
	import { mediaRootPreview } from '$lib/components/app/media/collection/mediaRootPreview';
	import type { FileNamingSettings, LibraryFolder, MediaItem } from '$lib/settings/types';

	interface Props {
		item: MediaItem;
		libraryFolders: LibraryFolder[];
		onCancel: () => void;
		onConfirm: (_libraryFolderId: string) => void;
	}

	let { item, libraryFolders, onCancel, onConfirm }: Props = $props();
	let open = $state(true);
	let libraryFolderId = $derived(item.libraryFolderId ?? '');
	let fileNamingSettings = $state<FileNamingSettings>();

	const folderOptions = $derived(
		libraryFolders.map((folder) => ({ value: folder.id, label: folder.path }))
	);
	const selectedFolder = $derived(libraryFolders.find((folder) => folder.id === libraryFolderId));
	const previewPath = $derived(mediaRootPreview(item, selectedFolder, fileNamingSettings));
	const canConfirm = $derived(libraryFolderId !== '' && libraryFolderId !== item.libraryFolderId);

	$effect(() => {
		void loadFileNamingSettings();
	});

	function handleOpenChange(nextOpen: boolean) {
		open = nextOpen;
		if (!nextOpen) onCancel();
	}

	async function loadFileNamingSettings() {
		try {
			fileNamingSettings = await getFileNamingSettings();
		} catch {
			fileNamingSettings = undefined;
		}
	}
</script>

<Dialog.Root bind:open onOpenChange={handleOpenChange}>
	<Dialog.Content class="w-[min(780px,calc(100vw-32px))]">
		<Dialog.Header>
			<Dialog.Title>Edit media root</Dialog.Title>
			<Dialog.Description>
				Choose the library root. The final folder name is generated from the naming template.
			</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-3">
			<div class="grid gap-2">
				<Label>Library root</Label>
				<SettingsSelect
					value={libraryFolderId}
					options={folderOptions}
					placeholder="Select root"
					onValueChange={(value) => (libraryFolderId = value)}
				/>
			</div>
			<div class="grid gap-1 rounded-md border bg-muted/40 px-3 py-2 text-sm">
				<strong class="text-xs font-extrabold text-muted-foreground uppercase"
					>New media root</strong
				>
				<span class="break-anywhere text-muted-foreground">
					{previewPath}
				</span>
			</div>
		</div>
		<Dialog.Footer>
			<Button type="button" variant="outline" onclick={onCancel}>Cancel</Button>
			<Button type="button" disabled={!canConfirm} onclick={() => onConfirm(libraryFolderId)}>
				Save
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
