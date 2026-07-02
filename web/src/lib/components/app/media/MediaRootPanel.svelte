<script lang="ts">
	import PencilIcon from '@lucide/svelte/icons/pencil';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { LibraryFolder, MediaItem, MediaItemUpdateRequest } from '$lib/settings/types';
	import MediaRootEditModal from './MediaRootEditModal.svelte';

	interface Props {
		item: MediaItem;
		libraryFolders: LibraryFolder[];
		canManage: boolean;
		onSaveOptions: (_item: MediaItem, _request: MediaItemUpdateRequest) => void;
	}

	let { item, libraryFolders, canManage, onSaveOptions }: Props = $props();
	let editing = $state(false);

	const rootPath = $derived(item.mediaFolderPath ?? item.libraryFolderPath ?? '-');
	const canEdit = $derived(canManage && libraryFolders.length > 0);

	function saveRoot(libraryFolderId: string) {
		onSaveOptions(item, { libraryFolderId });
		editing = false;
	}
</script>

<div
	class="grid grid-cols-[minmax(0,1fr)_auto] items-end gap-3 rounded-md border bg-card p-4 text-card-foreground shadow-xs"
>
	<div class="grid min-w-0 gap-1">
		<strong class="text-xs font-extrabold text-muted-foreground uppercase">Media root</strong>
		<span class="break-anywhere text-sm text-foreground">{rootPath}</span>
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
		{libraryFolders}
		onCancel={() => (editing = false)}
		onConfirm={saveRoot}
	/>
{/if}
