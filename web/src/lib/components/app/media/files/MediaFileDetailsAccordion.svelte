<script lang="ts">
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
	import TriangleAlertIcon from '@lucide/svelte/icons/triangle-alert';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils';
	import MediaFileTrackProvenanceIcon from '$lib/components/app/media/files/provenance/MediaFileTrackProvenanceIcon.svelte';
	import MediaFileTrackTypeIcon from '$lib/components/app/media/files/track-icons/MediaFileTrackTypeIcon.svelte';
	import {
		fileChapterDetailRows,
		fileChapterSummaryRow,
		fileTrackDetailRows
	} from '$lib/components/app/media/files/mediaFileDetails';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';

	interface Props {
		row: MediaFileRow;
	}

	let { row }: Props = $props();

	let chaptersExpanded = $state(false);

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
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each rows as track, index (track.key)}
				<Table.Row
					class={cn(
						index > 0 && track.type !== rows[index - 1]?.type && 'border-t-4 border-border',
						track.missing && 'bg-destructive/10 text-destructive',
						track.unwanted && 'bg-secondary/40',
						track.chapterSummary && 'cursor-pointer border-t-4 border-border'
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
										This track is not enabled for the profile and will be removed after the download
										client item is gone.
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
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={5} class="text-muted-foreground">No track details found.</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
