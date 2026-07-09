<script lang="ts">
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import { cn } from '$lib/utils';
	import MediaFileDetailsAccordion from '$lib/components/app/media/files/MediaFileDetailsAccordion.svelte';
	import MediaFileOtherFilesPanel from '$lib/components/app/media/files/other-files/MediaFileOtherFilesPanel.svelte';
	import MediaFilePreviewModal from '$lib/components/app/media/files/preview/MediaFilePreviewModal.svelte';
	import MediaFileRequirementIcon from '$lib/components/app/media/files/MediaFileRequirementIcon.svelte';
	import MediaFileSummaryActions from '$lib/components/app/media/files/MediaFileSummaryActions.svelte';
	import { fallbackRequirementStatus } from '$lib/components/app/media/files/mediaFileSummaryStatus';
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
		onSearchSubtitle = async () => {},
		onManualSubtitleSearch = () => {},
		onDeleteSubtitle = async () => {},
		onUpdateSubtitle = async () => {},
		onDeleteTrack = async () => {},
		onFulfillmentAction = async () => {},
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
	const videoStatus = $derived(
		row.requirements?.video ?? fallbackRequirementStatus('Video', row.exists)
	);
	const audioStatus = $derived(
		row.requirements?.audio ?? fallbackRequirementStatus('Audio', row.exists)
	);
	const subtitleStatus = $derived(
		row.requirements?.subtitles ?? fallbackRequirementStatus('Subtitles', row.exists)
	);

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
			'relative grid items-start gap-3 p-4 pb-5 lg:grid-cols-[minmax(260px,1fr)_76px_104px_60px_88px_auto]',
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
			<span class="flex items-center gap-1">
				<MediaFileRequirementIcon type="video" status={videoStatus} {row} />
				<MediaFileRequirementIcon type="audio" status={audioStatus} {row} />
				<MediaFileRequirementIcon type="subtitle" status={subtitleStatus} {row} />
			</span>
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
		<MediaFileDetailsAccordion {row} {canManage} {onDeleteTrack} {onFulfillmentAction} />
	{/if}
	<MediaFileOtherFilesPanel
		{row}
		{canManage}
		onSearch={(languageId) => onSearchSubtitle(row, languageId)}
		onManualSearch={(languageId) => onManualSubtitleSearch(row, languageId)}
		{onDeleteSubtitle}
		{onUpdateSubtitle}
		{onFulfillmentAction}
		{onDelete}
	/>
</div>

{#if row.exists && row.path && previewOpen}
	<MediaFilePreviewModal {mediaItemId} {mediaTitle} {row} onClose={() => (previewOpen = false)} />
{/if}
