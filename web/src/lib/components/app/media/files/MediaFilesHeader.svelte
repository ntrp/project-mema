<script lang="ts">
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import FilePenLineIcon from '@lucide/svelte/icons/file-pen-line';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MediaItem } from '$lib/settings/types';

	interface Props {
		item: MediaItem;
		canManage: boolean;
		scanningMediaItemId?: string;
		onRename: () => void;
		onRescanMediaFiles: (_item: MediaItem) => void;
	}

	let { item, canManage, scanningMediaItemId, onRename, onRescanMediaFiles }: Props = $props();
	const scanning = $derived(scanningMediaItemId === item.id);
	const canRename = $derived(canManage && item.filePaths.length > 0 && !!item.mediaFolderPath);
</script>

<div class="flex min-w-0 items-center justify-between gap-3">
	<h2 id="media-files-title" class="m-0 text-3xl font-semibold text-foreground">Files</h2>
	{#if canManage}
		<div class="flex items-center gap-2">
			<Tooltip.Root>
				<Tooltip.Trigger>
					{#snippet child({ props })}
						<Button
							{...props}
							type="button"
							variant="outline"
							size="icon-sm"
							aria-label="Rename files"
							disabled={!canRename}
							onclick={onRename}
						>
							<FilePenLineIcon aria-hidden="true" />
						</Button>
					{/snippet}
				</Tooltip.Trigger>
				<Tooltip.Content>Rename files</Tooltip.Content>
			</Tooltip.Root>
			<Tooltip.Root>
				<Tooltip.Trigger>
					{#snippet child({ props })}
						<Button
							{...props}
							type="button"
							variant="outline"
							size="icon-sm"
							aria-label={scanning ? 'Refreshing file metadata' : 'Refresh file metadata'}
							disabled={scanning || !item.mediaFolderPath}
							onclick={() => onRescanMediaFiles(item)}
						>
							<RefreshCwIcon class={scanning ? 'animate-spin' : undefined} aria-hidden="true" />
						</Button>
					{/snippet}
				</Tooltip.Trigger>
				<Tooltip.Content
					>{scanning ? 'Refreshing file metadata' : 'Refresh file metadata'}</Tooltip.Content
				>
			</Tooltip.Root>
		</div>
	{/if}
</div>
