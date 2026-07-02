<script lang="ts">
	import SearchIcon from '@lucide/svelte/icons/search';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import UserIcon from '@lucide/svelte/icons/user';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MediaFileRow } from './mediaFiles';

	interface Props {
		row: MediaFileRow;
		canManage: boolean;
		busy: boolean;
		showSearchActions?: boolean;
		onAutoSearch: () => void;
		onManualSearch: () => void;
		onDelete: (_row: MediaFileRow) => void;
	}

	let {
		row,
		canManage,
		busy,
		showSearchActions = true,
		onAutoSearch,
		onManualSearch,
		onDelete
	}: Props = $props();

	function stopActionClick(event: Event) {
		event.stopPropagation();
	}

	function stopActionKeydown(event: KeyboardEvent) {
		event.stopPropagation();
	}
</script>

<div
	role="presentation"
	class="flex flex-wrap items-start justify-end gap-2"
	onclick={stopActionClick}
	onkeydown={stopActionKeydown}
>
	{#if showSearchActions}
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#snippet child({ props })}
					<Button
						{...props}
						type="button"
						variant="outline"
						size="icon-sm"
						aria-label="Automatic search"
						disabled={!canManage || busy}
						onclick={onAutoSearch}
					>
						<SearchIcon aria-hidden="true" />
					</Button>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content>Automatic search</Tooltip.Content>
		</Tooltip.Root>
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#snippet child({ props })}
					<Button
						{...props}
						type="button"
						variant="outline"
						size="icon-sm"
						aria-label="Manual search"
						disabled={busy}
						onclick={onManualSearch}
					>
						<UserIcon aria-hidden="true" />
					</Button>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content>Manual search</Tooltip.Content>
		</Tooltip.Root>
	{/if}
	{#if row.exists}
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#snippet child({ props })}
					<Button
						{...props}
						type="button"
						variant="destructive"
						size="icon-sm"
						aria-label="Delete file"
						disabled={!canManage || !row.path}
						onclick={() => onDelete(row)}
					>
						<TrashIcon aria-hidden="true" />
					</Button>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content>Delete file</Tooltip.Content>
		</Tooltip.Root>
	{/if}
</div>
