<script lang="ts">
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import TriangleAlertIcon from '@lucide/svelte/icons/triangle-alert';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils';
	import MediaFileTrackProvenanceIcon from '$lib/components/app/media/files/provenance/MediaFileTrackProvenanceIcon.svelte';
	import MediaFileTrackTypeIcon from '$lib/components/app/media/files/track-icons/MediaFileTrackTypeIcon.svelte';
	import {
		fileChapterDetailRows,
		fileChapterSummaryRow,
		fileTrackDetailRows,
		type MediaFileDetailRow
	} from '$lib/components/app/media/files/mediaFileDetails';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import type { MediaFileTrackDeleteRequest } from '$lib/settings/types';

	interface Props {
		row: MediaFileRow;
		canManage?: boolean;
		onDeleteTrack?: (
			_row: MediaFileRow,
			_request: MediaFileTrackDeleteRequest
		) => void | Promise<void>;
	}

	let { row, canManage = false, onDeleteTrack = async () => {} }: Props = $props();

	let chaptersExpanded = $state(false);
	let deleteTarget = $state<MediaFileDetailRow | undefined>();

	const trackRows = $derived(fileTrackDetailRows(row));
	const chapterRows = $derived(fileChapterDetailRows(row));
	const chapterSummary = $derived(fileChapterSummaryRow(row));
	const rows = $derived.by(() => [
		...trackRows,
		...(chapterSummary ? [chapterSummary, ...(chaptersExpanded ? chapterRows : [])] : [])
	]);

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

	function deleteDescription(track: MediaFileDetailRow) {
		if (track.chapterSummary) return 'Delete all chapters from this file?';
		if (track.type === 'chapter') return `Delete chapter ${track.trackNumber} from this file?`;
		return `Delete ${track.type} track ${track.trackNumber} from this file?`;
	}
</script>

<div class="overflow-x-auto border-t border-border bg-background" aria-label="Track details">
	<Table.Root class="min-w-170 text-sm">
		<Table.Header>
			<Table.Row>
				<Table.Head class="w-24">Track Nr.</Table.Head>
				<Table.Head class="w-20">Type</Table.Head>
				<Table.Head class="w-36">Language</Table.Head>
				<Table.Head>Track description</Table.Head>
				<Table.Head class="w-24">Provenance</Table.Head>
				<Table.Head class="w-20 text-right">Actions</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each rows as track, index (track.key)}
				<Table.Row
					class={cn(
						index > 0 && track.type !== rows[index - 1]?.type && 'border-t-4 border-border',
						track.missing && 'bg-destructive/10 text-destructive',
						track.unwanted && 'bg-secondary/40',
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
							{#if track.unwanted}
								<Tooltip.Root>
									<Tooltip.Trigger>
										{#snippet child({ props })}
											<button
												{...props}
												type="button"
												class="inline-flex border-0 bg-transparent p-0 text-secondary-foreground"
											>
												<TriangleAlertIcon class="size-4" aria-label="Not wanted" />
											</button>
										{/snippet}
									</Tooltip.Trigger>
									<Tooltip.Content>
										This track does not match the subtitle mode or enabled profile targets.
									</Tooltip.Content>
								</Tooltip.Root>
							{/if}
						</span>
					</Table.Cell>
					<Table.Cell class="justify-end">
						{#if !track.missing && track.type !== 'chapter'}
							<MediaFileTrackProvenanceIcon provenance={track.provenance} />
						{/if}
					</Table.Cell>
					<Table.Cell class="text-right">
						{#if track.deleteRequest}
							<Tooltip.Root>
								<Tooltip.Trigger>
									{#snippet child({ props })}
										<Button
											{...props}
											type="button"
											variant="destructive"
											size="icon-sm"
											aria-label="Delete embedded track"
											disabled={!canManage}
											onclick={(event) => requestDelete(event, track)}
											onkeydown={(event) => event.stopPropagation()}
										>
											<TrashIcon aria-hidden="true" />
										</Button>
									{/snippet}
								</Tooltip.Trigger>
								<Tooltip.Content>Delete embedded track</Tooltip.Content>
							</Tooltip.Root>
						{/if}
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

<Dialog.Root open={!!deleteTarget} onOpenChange={(open) => !open && (deleteTarget = undefined)}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Delete embedded track</Dialog.Title>
			<Dialog.Description>
				{deleteTarget ? deleteDescription(deleteTarget) : ''}
			</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer>
			<Button type="button" variant="outline" onclick={() => (deleteTarget = undefined)}>
				Cancel
			</Button>
			<Button type="button" variant="destructive" onclick={confirmDelete}>Delete</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
