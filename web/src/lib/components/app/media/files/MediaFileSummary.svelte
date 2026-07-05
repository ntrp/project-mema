<script lang="ts">
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import { Badge } from '$lib/components/ui/badge';
	import { cn } from '$lib/utils';
	import MediaFileDetailsAccordion from '$lib/components/app/media/files/MediaFileDetailsAccordion.svelte';
	import MediaFilePreviewModal from '$lib/components/app/media/files/preview/MediaFilePreviewModal.svelte';
	import MediaFileSummaryActions from '$lib/components/app/media/files/MediaFileSummaryActions.svelte';
	import { displayLanguage } from '$lib/settings/languageDisplay';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import type { ActivityQueueStatus } from '$lib/components/app/activity/activityQueue';

	interface Props {
		mediaItemId: string;
		mediaTitle: string;
		row: MediaFileRow;
		activityStatus?: ActivityQueueStatus;
		canManage: boolean;
		searching: boolean;
		fileLabel?: string;
		missingLabel?: string;
		showSearchActions?: boolean;
		onAutoSearch: () => void;
		onManualSearch: () => void;
		onDelete: (_row: MediaFileRow) => void;
	}

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
	const subtitleState = $derived(row.subtitleSatisfaction?.state ?? 'ignored');
	const subtitleLabel = $derived(subtitleSatisfactionLabel(row));

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

	function subtitleSatisfactionLabel(row: MediaFileRow) {
		const satisfaction = row.subtitleSatisfaction;
		if (!satisfaction || satisfaction.state === 'ignored') return 'Ignored';
		if (satisfaction.state === 'satisfied') {
			return `Satisfied: ${languageList(satisfaction.matchedLanguages)}`;
		}
		return `Missing: ${languageList(satisfaction.missingLanguages)}`;
	}

	function languageList(values: string[]) {
		return values.map(displayLanguage).join(', ') || '-';
	}
</script>

<div
	class={cn(
		'relative overflow-hidden rounded-md border bg-card text-card-foreground shadow-xs',
		!row.exists && 'border-dashed bg-muted/30'
	)}
>
	<div
		role="button"
		tabindex={row.exists ? 0 : -1}
		aria-disabled={!row.exists}
		aria-expanded={row.exists ? detailsOpen : undefined}
		class={cn(
			'relative grid gap-3 p-4 pb-5 lg:grid-cols-[minmax(180px,1.2fr)_repeat(5,minmax(84px,0.55fr))_minmax(120px,0.8fr)_auto]',
			row.exists &&
				'cursor-pointer transition-[border-color,box-shadow] hover:border-primary/40 hover:shadow-sm focus-visible:border-primary/50 focus-visible:outline-none'
		)}
		onclick={toggleDetails}
		onkeydown={handleCardKeydown}
	>
		<div class="grid min-w-0 gap-1">
			<strong class="break-anywhere text-sm font-semibold">
				{row.exists ? row.relativePath : 'Missing file'}
			</strong>
			<span class="text-sm text-muted-foreground">{row.exists ? fileLabel : missingLabel}</span>
		</div>
		<span class="grid gap-1">
			<strong class="text-xs font-medium text-muted-foreground uppercase">Size</strong>
			<span class="text-sm">{row.size}</span>
		</span>
		<span class="grid gap-1">
			<strong class="text-xs font-medium text-muted-foreground uppercase">Quality</strong>
			<span class="text-sm">{row.quality}</span>
		</span>
		<span class="grid gap-1">
			<strong class="text-xs font-medium text-muted-foreground uppercase">Formats</strong>
			<span class="flex flex-wrap gap-1">
				{#if row.formats.length > 0}
					{#each row.formats as format (format)}
						<Badge variant="outline">{format}</Badge>
					{/each}
				{:else}
					<span class="text-sm">-</span>
				{/if}
			</span>
		</span>
		<span class="grid gap-1">
			<strong class="text-xs font-medium text-muted-foreground uppercase">Score</strong>
			<span class="text-sm">{row.score}</span>
		</span>
		<span class="grid gap-1">
			<strong class="text-xs font-medium text-muted-foreground uppercase">Subtitles</strong>
			{#if subtitleState !== 'ignored'}
				<Badge
					variant={subtitleState === 'missing' ? 'destructive' : 'secondary'}
					class="justify-self-start"
				>
					{subtitleLabel}
				</Badge>
			{:else}
				<span class="text-sm">{subtitleLabel}</span>
			{/if}
		</span>
		<span class="grid gap-1">
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
				<span class="text-sm">-</span>
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
</div>

{#if row.exists && row.path && previewOpen}
	<MediaFilePreviewModal {mediaItemId} {mediaTitle} {row} onClose={() => (previewOpen = false)} />
{/if}
