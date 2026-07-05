<script lang="ts">
	import CaptionsIcon from '@lucide/svelte/icons/captions';
	import ClapperboardIcon from '@lucide/svelte/icons/clapperboard';
	import FileTextIcon from '@lucide/svelte/icons/file-text';
	import MusicIcon from '@lucide/svelte/icons/music';
	import TriangleAlertIcon from '@lucide/svelte/icons/triangle-alert';
	import VideoIcon from '@lucide/svelte/icons/video';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils';
	import { fileDetailRows } from '$lib/components/app/media/files/mediaFileDetails';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';

	interface Props {
		row: MediaFileRow;
	}

	let { row }: Props = $props();

	const rows = $derived(fileDetailRows(row));
</script>

<div class="overflow-x-auto border-t border-border bg-background" aria-label="Track details">
	<Table.Root class="min-w-170 text-sm">
		<Table.Header>
			<Table.Row>
				<Table.Head class="w-24">Track Nr.</Table.Head>
				<Table.Head class="w-20">Type</Table.Head>
				<Table.Head class="w-36">Language</Table.Head>
				<Table.Head>Track description</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each rows as track, index (track.key)}
				<Table.Row
					class={cn(
						index > 0 && track.type !== rows[index - 1]?.type && 'border-t-4 border-border',
						track.missing && 'bg-destructive/10 text-destructive',
						track.unwanted && 'bg-secondary/40'
					)}
				>
					<Table.Cell>{track.trackNumber}</Table.Cell>
					<Table.Cell>
						<span class="inline-flex items-center">
							{#if track.type === 'video'}
								<VideoIcon aria-label="Video" />
							{:else if track.type === 'audio'}
								<MusicIcon aria-label="Audio" />
							{:else if track.type === 'subtitle'}
								<CaptionsIcon aria-label="Subtitle" />
							{:else if track.type === 'chapter'}
								<ClapperboardIcon aria-label="Chapter" />
							{:else}
								<FileTextIcon aria-label="Track" />
							{/if}
						</span>
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
										This track is not enabled for the profile and will be removed after the
										download client item is gone.
									</Tooltip.Content>
								</Tooltip.Root>
							{/if}
						</span>
					</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={4} class="text-muted-foreground">No track details found.</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
