<script lang="ts">
	import InfoIcon from '@lucide/svelte/icons/info';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import SearchIcon from '@lucide/svelte/icons/search';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import UserIcon from '@lucide/svelte/icons/user';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils';
	import type { MediaFileRow } from './mediaFiles';
	import type { ActivityQueueStatus } from '../activity/activityQueue';

	interface Props {
		row: MediaFileRow;
		activityStatus?: ActivityQueueStatus;
		canManage: boolean;
		searching: boolean;
		fileLabel?: string;
		missingLabel?: string;
		onInfo: (_row: MediaFileRow) => void;
		onAutoSearch: () => void;
		onManualSearch: () => void;
		onDelete: (_row: MediaFileRow) => void;
	}

	let {
		row,
		activityStatus,
		canManage,
		searching,
		fileLabel = 'Episode file',
		missingLabel = 'No matched file for this episode',
		onInfo,
		onAutoSearch,
		onManualSearch,
		onDelete
	}: Props = $props();
	const busy = $derived(
		searching ||
			activityStatus?.status === 'queued' ||
			activityStatus?.status === 'grabbed' ||
			activityStatus?.status === 'downloading'
	);
</script>

<div
	class={cn(
		'grid gap-4 rounded-md border bg-card p-4 text-card-foreground shadow-xs',
		!row.exists && 'border-dashed bg-muted/30'
	)}
>
	<div class="grid gap-1">
		<strong class="break-all text-sm font-semibold"
			>{row.exists ? row.relativePath : 'Missing file'}</strong
		>
		<span class="text-sm text-muted-foreground">{row.exists ? fileLabel : missingLabel}</span>
	</div>

	<div class="grid gap-3 text-sm sm:grid-cols-2 lg:grid-cols-3" aria-label="Episode file details">
		<span class="grid gap-1 rounded-md border bg-background px-3 py-2">
			<strong class="text-xs font-medium uppercase text-muted-foreground">Quality</strong>
			<span>{row.quality}</span>
		</span>
		<span class="grid gap-1 rounded-md border bg-background px-3 py-2">
			<strong class="text-xs font-medium uppercase text-muted-foreground">Video</strong>
			<span>{row.videoCodec}</span>
		</span>
		<span class="grid gap-1 rounded-md border bg-background px-3 py-2">
			<strong class="text-xs font-medium uppercase text-muted-foreground">Audio</strong>
			<span>{row.audioInfo}</span>
		</span>
		<span class="grid gap-1 rounded-md border bg-background px-3 py-2">
			<strong class="text-xs font-medium uppercase text-muted-foreground">Languages</strong>
			<span>{row.languages}</span>
		</span>
		<span class="grid gap-1 rounded-md border bg-background px-3 py-2">
			<strong class="text-xs font-medium uppercase text-muted-foreground">Score</strong>
			<span>{row.score}</span>
		</span>
		<span class="grid gap-1 rounded-md border bg-background px-3 py-2">
			<strong class="text-xs font-medium uppercase text-muted-foreground">Status</strong>
			{#if activityStatus}
				<Badge
					variant={activityStatus.status === 'failed' ? 'destructive' : 'secondary'}
					class="justify-self-start"
				>
					<RefreshCwIcon aria-hidden="true" />
					{activityStatus.label}
				</Badge>
			{:else}
				-
			{/if}
		</span>
	</div>

	{#if row.formats.length > 0}
		<div class="flex flex-wrap gap-2" aria-label="Matched formats">
			{#each row.formats as format (format)}
				<Badge variant="outline">{format}</Badge>
			{/each}
		</div>
	{/if}

	<div class="flex flex-wrap justify-end gap-2">
		{#if row.exists}
			<Tooltip.Root>
				<Tooltip.Trigger>
					{#snippet child({ props })}
						<Button
							{...props}
							type="button"
							variant="outline"
							size="icon-sm"
							aria-label="File info"
							onclick={() => onInfo(row)}
						>
							<InfoIcon aria-hidden="true" />
						</Button>
					{/snippet}
				</Tooltip.Trigger>
				<Tooltip.Content>File info</Tooltip.Content>
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
</div>
