<script lang="ts">
	import CaptionsIcon from '@lucide/svelte/icons/captions';
	import ClapperboardIcon from '@lucide/svelte/icons/clapperboard';
	import FileTextIcon from '@lucide/svelte/icons/file-text';
	import MusicIcon from '@lucide/svelte/icons/music';
	import VideoIcon from '@lucide/svelte/icons/video';
	import * as Table from '$lib/components/ui/table';
	import { fileDetailRows } from './mediaFileDetails';
	import type { MediaFileRow } from './mediaFiles';

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
			{#each rows as track (track.key)}
				<Table.Row class={track.missing ? 'bg-destructive/10 text-destructive' : ''}>
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
					<Table.Cell class="whitespace-normal">{track.description}</Table.Cell>
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell colspan={4} class="text-muted-foreground">No track details found.</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
