<script lang="ts">
	import { onMount } from 'svelte';
	import PencilIcon from '@lucide/svelte/icons/pencil';
	import TriangleAlertIcon from '@lucide/svelte/icons/triangle-alert';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { getFileNamingSettings } from '$lib/settings/api';
	import { mediaRootWarning } from '$lib/components/app/media/collection/mediaRootPreview';
	import { matchingLibraryFolders } from '$lib/components/app/media/actions/mediaActionDefaults';
	import type {
		FileNamingSettings,
		LibraryFolder,
		MediaItem,
		MediaItemUpdateRequest
	} from '$lib/settings/types';
	import MediaRootEditModal from '$lib/components/app/media/collection/MediaRootEditModal.svelte';

	interface Props {
		item: MediaItem;
		libraryFolders: LibraryFolder[];
		canManage: boolean;
		onSaveOptions: (_item: MediaItem, _request: MediaItemUpdateRequest) => void;
	}

	let { item, libraryFolders, canManage, onSaveOptions }: Props = $props();
	let editing = $state(false);
	let fileNamingSettings = $state<FileNamingSettings>();

	const rootPath = $derived(item.mediaFolderPath ?? '-');
	const matchingFolders = $derived(matchingLibraryFolders(item.type, libraryFolders));
	const selectedFolder = $derived(
		libraryFolders.find((folder) => folder.id === item.libraryFolderId) ??
			libraryFolders.find((folder) => folder.path === item.libraryFolderPath)
	);
	const warning = $derived(mediaRootWarning(item, selectedFolder, fileNamingSettings));
	const canEdit = $derived(canManage && matchingFolders.length > 0);

	onMount(() => {
		void loadFileNamingSettings();
	});

	function saveRoot(libraryFolderId: string) {
		onSaveOptions(item, { libraryFolderId });
		editing = false;
	}

	async function loadFileNamingSettings() {
		try {
			fileNamingSettings = await getFileNamingSettings();
		} catch {
			fileNamingSettings = undefined;
		}
	}
</script>

<div
	class="grid grid-cols-[minmax(0,1fr)_auto] items-end gap-3 rounded-md border bg-card p-4 text-card-foreground shadow-xs"
>
	<div class="grid min-w-0 gap-1">
		<strong class="text-xs font-extrabold text-muted-foreground uppercase">Media root</strong>
		<span class="flex min-w-0 items-start gap-2 text-sm text-foreground">
			<span class="break-anywhere min-w-0">{rootPath}</span>
			{#if warning}
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<TriangleAlertIcon
								{...props}
								class="mt-0.5 size-4 shrink-0 text-amber-500"
								aria-label="Media root does not match template"
							/>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content class="grid max-w-[360px] gap-1">
						<span>Media root does not match the naming template.</span>
						<span class="break-anywhere">Expected: {warning.expected}</span>
					</Tooltip.Content>
				</Tooltip.Root>
			{/if}
		</span>
	</div>
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					variant="outline"
					size="icon-sm"
					aria-label="Edit media root"
					disabled={!canEdit}
					onclick={() => (editing = true)}
				>
					<PencilIcon aria-hidden="true" />
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>Edit media root</Tooltip.Content>
	</Tooltip.Root>
</div>

{#if editing}
	<MediaRootEditModal
		{item}
		libraryFolders={matchingFolders}
		onCancel={() => (editing = false)}
		onConfirm={saveRoot}
	/>
{/if}
