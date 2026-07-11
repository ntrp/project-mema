<script lang="ts">
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
	import * as Table from '$lib/components/ui/table';
	import { cn } from '$lib/utils';
	import MediaFileDeleteTrackButton from '$lib/components/app/media/files/MediaFileDeleteTrackButton.svelte';
	import MediaFileFulfillmentActions from '$lib/components/app/media/files/MediaFileFulfillmentActions.svelte';
	import MediaFileTrackDeleteDialog from '$lib/components/app/media/files/MediaFileTrackDeleteDialog.svelte';
	import MediaFileDetailStateBadge from '$lib/components/app/media/files/details/MediaFileDetailStateBadge.svelte';
	import { unwantedMediaRowClass } from '$lib/components/app/media/files/details/mediaFileVisualClasses';
	import MediaFileTrackProvenanceIcon from '$lib/components/app/media/files/provenance/MediaFileTrackProvenanceIcon.svelte';
	import MediaFileTrackTypeIcon from '$lib/components/app/media/files/track-icons/MediaFileTrackTypeIcon.svelte';
	import {
		fileChapterDetailRows,
		fileChapterSummaryRow,
		fileTrackDetailRows,
		type MediaFileDetailRow
	} from '$lib/components/app/media/files/mediaFileDetails';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import type {
		MediaFileTrackDeleteRequest,
		MediaFulfillmentActionRequest
	} from '$lib/settings/types';
	import MediaFileDetailsHeader from './details/MediaFileDetailsHeader.svelte';
	import { createMediaFileRowPulse } from './details/mediaFileRowPulse';

	interface Props {
		row: MediaFileRow;
		canManage?: boolean;
		pendingFulfillmentActionKeys?: string[];
		onDeleteTrack?: (
			_row: MediaFileRow,
			_request: MediaFileTrackDeleteRequest
		) => void | Promise<void>;
		onFulfillmentAction?: (
			_row: MediaFileRow,
			_request: MediaFulfillmentActionRequest
		) => void | Promise<void>;
	}

	let {
		row,
		canManage = false,
		pendingFulfillmentActionKeys = [],
		onDeleteTrack = async () => {},
		onFulfillmentAction = async () => {}
	}: Props = $props();

	let chaptersExpanded = $state(false);
	let deleteTarget = $state<MediaFileDetailRow | undefined>();

	const trackRows = $derived(fileTrackDetailRows(row));
	const chapterRows = $derived(fileChapterDetailRows(row));
	const chapterSummary = $derived(fileChapterSummaryRow(row));
	const rows = $derived.by(() => [
		...trackRows,
		...(chapterSummary ? [chapterSummary, ...(chaptersExpanded ? chapterRows : [])] : [])
	]);
	const trackPulses = createMediaFileRowPulse();
	const pulsingRows = $derived(trackPulses(rows));

	function toggleChapters() {
		chaptersExpanded = !chaptersExpanded;
	}

	function handleChapterSummaryKeydown(event: KeyboardEvent) {
		if (event.key !== 'Enter' && event.key !== ' ') return;
		event.preventDefault();
		toggleChapters();
	}

	function requestDelete(event: Event, track: MediaFileDetailRow) {
		event.stopPropagation();
		if (!track.deleteRequest) return;
		deleteTarget = track;
	}

	async function confirmDelete() {
		if (!deleteTarget?.deleteRequest) return;
		await onDeleteTrack(row, { path: row.path ?? '', ...deleteTarget.deleteRequest });
		deleteTarget = undefined;
	}
</script>

<div class="overflow-x-auto border-t border-border bg-background" aria-label="Track details">
	<Table.Root class="min-w-170 text-sm">
		<MediaFileDetailsHeader />
		<Table.Body>
			{#each rows as track, index (track.key)}
				<Table.Row
					class={cn(
						index > 0 && track.type !== rows[index - 1]?.type && 'border-t-4 border-border',
						track.missing && 'bg-destructive/10 text-destructive',
						track.unwanted && unwantedMediaRowClass,
						pulsingRows.has(track.key) && 'live-row-pulse',
						track.chapterSummary &&
							'cursor-pointer border-t-4 border-border [&>td]:border-t-4 [&>td]:border-border'
					)}
					role={track.chapterSummary ? 'button' : undefined}
					tabindex={track.chapterSummary ? 0 : undefined}
					aria-expanded={track.chapterSummary ? chaptersExpanded : undefined}
					aria-label={track.chapterSummary
						? chaptersExpanded
							? 'Collapse chapters'
							: 'Expand chapters'
						: undefined}
					onclick={track.chapterSummary ? toggleChapters : undefined}
					onkeydown={track.chapterSummary ? handleChapterSummaryKeydown : undefined}
				>
					<Table.Cell>
						{#if track.chapterSummary}
							<span class="inline-flex items-center gap-1 text-foreground">
								{#if chaptersExpanded}
									<ChevronDownIcon class="size-4" aria-hidden="true" />
								{:else}
									<ChevronRightIcon class="size-4" aria-hidden="true" />
								{/if}
								<span>{track.trackNumber}</span>
							</span>
						{:else}
							{track.trackNumber}
						{/if}
					</Table.Cell>
					<Table.Cell>
						<MediaFileTrackTypeIcon type={track.type} />
					</Table.Cell>
					<Table.Cell>{track.language}</Table.Cell>
					<Table.Cell class="whitespace-normal">
						<span class="inline-flex items-center gap-2">
							{track.description}
							<MediaFileDetailStateBadge row={track} />
						</span>
					</Table.Cell>
					<Table.Cell class="justify-end">
						{#if !track.missing && track.type !== 'chapter'}
							<MediaFileTrackProvenanceIcon provenance={track.provenance} />
						{/if}
					</Table.Cell>
					<Table.Cell class="text-right">
						<span class="inline-flex justify-end gap-1">
							<MediaFileFulfillmentActions
								row={track}
								{canManage}
								{pendingFulfillmentActionKeys}
								onFulfillmentAction={(request) => onFulfillmentAction(row, request)}
							/>
							<MediaFileDeleteTrackButton {track} {canManage} onRequestDelete={requestDelete} />
						</span>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={6} class="text-muted-foreground">No track details found.</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>

{#if deleteTarget}
	<MediaFileTrackDeleteDialog
		track={deleteTarget}
		onCancel={() => (deleteTarget = undefined)}
		onConfirm={confirmDelete}
	/>
{/if}
