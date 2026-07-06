<script lang="ts">
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import { Badge } from '$lib/components/ui/badge';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils';
	import MediaFileDetailsAccordion from '$lib/components/app/media/files/MediaFileDetailsAccordion.svelte';
	import MediaFileOtherFilesPanel from '$lib/components/app/media/files/other-files/MediaFileOtherFilesPanel.svelte';
	import MediaFilePreviewModal from '$lib/components/app/media/files/preview/MediaFilePreviewModal.svelte';
	import MediaFileSummaryActions from '$lib/components/app/media/files/MediaFileSummaryActions.svelte';
	import {
		audioSatisfaction,
		subtitleSatisfaction
	} from '$lib/components/app/media/files/mediaFileSummaryStatus';
	import type { MediaFileSummaryProps as Props } from '$lib/components/app/media/file-data/mediaFileComponentTypes';

	let {
		mediaItemId,
		mediaTitle,
		row,
		activityStatus,
		canManage,
		searching,
		fileLabel = 'Episode file',
		missingLabel = 'No matched file for this episode',
		showSearchActions = true,
		onAutoSearch,
		onManualSearch,
		onDelete
	}: Props = $props();
	let detailsOpen = $state(false);
	let previewOpen = $state(false);
	const busy = $derived(
		searching ||
			activityStatus?.status === 'queued' ||
			activityStatus?.status === 'grabbed' ||
			activityStatus?.status === 'downloading'
	);
	const statusLabel = $derived(activityStatus?.label ?? '-');
	const audioStatus = $derived(audioSatisfaction(row));
	const subtitleStatus = $derived(subtitleSatisfaction(row));
	const upgradeVariant = $derived(row.upgrade.state === 'upgradeable' ? 'secondary' : 'outline');

	function toggleDetails() {
		if (row.exists) {
			detailsOpen = !detailsOpen;
		}
	}

	function handleCardKeydown(event: KeyboardEvent) {
		if (!row.exists || (event.key !== 'Enter' && event.key !== ' ')) return;
		event.preventDefault();
		toggleDetails();
	}
</script>

<div
	class={cn(
		'relative rounded-md border bg-card text-card-foreground shadow-xs',
		!row.exists && 'border-dashed bg-muted/30'
	)}
>
	<div
		role="button"
		tabindex={row.exists ? 0 : -1}
		aria-disabled={!row.exists}
		aria-expanded={row.exists ? detailsOpen : undefined}
		class={cn(
			'relative grid items-start gap-3 p-4 pb-5 lg:grid-cols-[minmax(260px,1fr)_118px_118px_76px_104px_60px_116px_auto]',
			row.exists &&
				'cursor-pointer transition-[border-color,box-shadow] hover:border-primary/40 hover:shadow-sm focus-visible:border-primary/50 focus-visible:outline-none'
		)}
		onclick={toggleDetails}
		onkeydown={handleCardKeydown}
	>
		<div class="grid min-w-0 content-start gap-1">
			<strong class="text-xs font-medium text-muted-foreground uppercase">File</strong>
			<span class="break-anywhere text-sm font-semibold">
				{row.exists ? row.relativePath : 'Missing file'}
			</span>
			<span class="sr-only">{row.exists ? fileLabel : missingLabel}</span>
		</div>
		<span class="grid content-start gap-1">
			<strong class="text-xs font-medium text-muted-foreground uppercase">Audio</strong>
			<Badge
				variant={audioStatus.state === 'missing' ? 'destructive' : 'secondary'}
				class="justify-self-start"
			>
				{audioStatus.label}
			</Badge>
		</span>
		<span class="grid content-start gap-1">
			<strong class="text-xs font-medium text-muted-foreground uppercase">Subtitles</strong>
			{#if subtitleStatus.state !== 'ignored'}
				<Badge
					variant={subtitleStatus.state === 'missing' ? 'destructive' : 'secondary'}
					class="justify-self-start"
				>
					{subtitleStatus.label}
				</Badge>
			{:else}
				<span class="text-sm">{subtitleStatus.label}</span>
			{/if}
		</span>
		<span class="grid content-start gap-1">
			<strong class="text-xs font-medium text-muted-foreground uppercase">Size</strong>
			<span class="text-sm">{row.size}</span>
		</span>
		<span class="grid content-start gap-1">
			<strong class="text-xs font-medium text-muted-foreground uppercase">Quality</strong>
			<span class="text-sm">{row.quality}</span>
		</span>
		<span class="grid content-start gap-1">
			<strong class="text-xs font-medium text-muted-foreground uppercase">Score</strong>
			<span class="text-sm">{row.score}</span>
		</span>
		<span class="grid content-start gap-1">
			<strong class="text-xs font-medium text-muted-foreground uppercase">Status</strong>
			{#if activityStatus}
				<Badge
					variant={activityStatus.status === 'failed' ? 'destructive' : 'secondary'}
					class="justify-self-start"
				>
					<RefreshCwIcon aria-hidden="true" />
					{statusLabel}
				</Badge>
			{:else}
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#snippet child({ props })}
							<Badge
								{...props}
								variant={upgradeVariant}
								class="justify-self-start"
								aria-label={row.upgrade.reasons.join(' ')}
							>
								{row.upgrade.label}
							</Badge>
						{/snippet}
					</Tooltip.Trigger>
					<Tooltip.Content>{row.upgrade.reasons.join(', ')}</Tooltip.Content>
				</Tooltip.Root>
			{/if}
		</span>
		<MediaFileSummaryActions
			{row}
			{canManage}
			{busy}
			{showSearchActions}
			{onAutoSearch}
			{onManualSearch}
			onPreview={() => (previewOpen = true)}
			{onDelete}
		/>
		{#if row.exists}
			<span
				class="pointer-events-none absolute bottom-0 left-1/2 z-[2] grid h-4 w-9 -translate-x-1/2 translate-y-1/2 place-items-center rounded-md border border-border bg-card text-muted-foreground"
				aria-hidden="true"
			>
				<ChevronDownIcon class={cn('size-3.5 transition-transform', detailsOpen && 'rotate-180')} />
			</span>
		{/if}
	</div>

	{#if row.exists && detailsOpen}
		<MediaFileDetailsAccordion {row} />
	{/if}
	<MediaFileOtherFilesPanel {row} {canManage} {onDelete} />
</div>

{#if row.exists && row.path && previewOpen}
	<MediaFilePreviewModal {mediaItemId} {mediaTitle} {row} onClose={() => (previewOpen = false)} />
{/if}
