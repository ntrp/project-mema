<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import CirclePlayIcon from '@lucide/svelte/icons/circle-play';
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';
	import PackageIcon from '@lucide/svelte/icons/package';
	import SearchIcon from '@lucide/svelte/icons/search';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import UserIcon from '@lucide/svelte/icons/user';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import { mediaFulfillmentActionKey } from '$lib/settings/mediaFulfillmentActionKey';
	import type { MediaFulfillmentActionRequest } from '$lib/settings/types';

	interface Props {
		row: MediaFileRow;
		canManage: boolean;
		busy: boolean;
		pendingFulfillmentActionKeys?: string[];
		showSearchActions?: boolean;
		onAutoSearch: () => void;
		onManualSearch: () => void;
		onPreview: () => void;
		onFulfillmentAction: (_row: MediaFileRow, _request: MediaFulfillmentActionRequest) => void | Promise<void>;
		onDelete: (_row: MediaFileRow) => void;
	}

	let {
		row,
		canManage,
		busy,
		pendingFulfillmentActionKeys = [],
		showSearchActions = true,
		onAutoSearch,
		onManualSearch,
		onPreview,
		onFulfillmentAction,
		onDelete
	}: Props = $props();

	let pendingAction = $state<string | undefined>();
	const remuxRequest = $derived<MediaFulfillmentActionRequest>({
		operation: 'container_remux',
		filePath: row.path,
		targetType: 'video'
	});
	const remuxKey = $derived(mediaFulfillmentActionKey(remuxRequest));
	const remuxAvailable = $derived(row.exists && row.path && row.requirements?.container?.state === 'pending');
	const remuxPending = $derived(
		pendingAction === remuxKey || pendingFulfillmentActionKeys.includes(remuxKey)
	);

	function stopActionClick(event: Event) {
		event.stopPropagation();
	}

	function stopActionKeydown(event: KeyboardEvent) {
		event.stopPropagation();
	}

	async function remuxContainer() {
		if (remuxPending) {
			await goto(resolve('/system/jobs'));
			return;
		}
		pendingAction = remuxKey;
		try {
			await onFulfillmentAction(row, remuxRequest);
		} finally {
			if (pendingAction === remuxKey) pendingAction = undefined;
		}
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
		{#if remuxAvailable}
			<Tooltip.Root>
				<Tooltip.Trigger>
					{#snippet child({ props })}
						<Button
							{...props}
							type="button"
							variant="outline"
							size="icon-sm"
							aria-label="Remux container"
							aria-busy={remuxPending}
							disabled={!canManage || !row.path}
							onclick={remuxContainer}
						>
							{#if remuxPending}
								<LoaderCircleIcon class="animate-spin" aria-hidden="true" />
							{:else}
								<PackageIcon aria-hidden="true" />
							{/if}
						</Button>
					{/snippet}
				</Tooltip.Trigger>
				<Tooltip.Content>{remuxPending ? 'View remux job' : 'Remux container'}</Tooltip.Content>
			</Tooltip.Root>
		{/if}
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#snippet child({ props })}
					<Button
						{...props}
						type="button"
						variant="outline"
						size="icon-sm"
						aria-label="Preview file"
						disabled={!row.path}
						onclick={onPreview}
					>
						<CirclePlayIcon aria-hidden="true" />
					</Button>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content>Preview file</Tooltip.Content>
		</Tooltip.Root>
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
