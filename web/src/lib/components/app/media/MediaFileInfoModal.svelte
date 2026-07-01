<script lang="ts">
	import SettingsFormModal from '$lib/components/settings/shared/SettingsFormModal.svelte';
	import { Card } from '$lib/components/ui/card';
	import type { MediaFileRow } from './mediaFiles';

	interface Props {
		row: MediaFileRow;
		onClose: () => void;
	}

	let { row, onClose }: Props = $props();

	const fields = $derived([
		['Relative path', row.relativePath],
		['Quality', row.quality],
		['Video codec', row.videoCodec],
		['Audio info', row.audioInfo],
		['Languages', row.languages],
		['Formats', row.formats.join(', ') || '-'],
		['Score', String(row.score)]
	]);
	const tracks = $derived([
		{
			title: 'Video',
			rows: [
				['Codec', row.videoCodec],
				['Resolution', row.quality],
				['Bitrate', '-']
			]
		},
		{
			title: 'Audio',
			rows: [
				['Codec', row.audioInfo],
				['Languages', row.languages],
				['Channels', '-']
			]
		},
		{
			title: 'Subtitles',
			rows: [
				['Languages', '-'],
				['Forced', '-'],
				['Count', '-']
			]
		}
	]);
</script>

<SettingsFormModal title="File details" modalClass="w-[min(1960px,calc(100vw-32px))]" {onClose}>
	<div class="grid gap-2.5 sm:grid-cols-2 lg:grid-cols-3">
		{#each fields as [label, value] (label)}
			<Card class="gap-1 rounded-md border-border bg-card p-3 shadow-none">
				<span class="text-xs font-extrabold text-muted-foreground uppercase">{label}</span>
				<strong class="break-anywhere text-sm text-foreground">{value}</strong>
			</Card>
		{/each}
	</div>
	<div class="grid gap-3 sm:grid-cols-3">
		{#each tracks as track (track.title)}
			<Card class="gap-2 rounded-md bg-card p-3 shadow-none">
				<h3 class="m-0 text-base font-semibold text-foreground">{track.title}</h3>
				<div class="grid gap-1">
					{#each track.rows as [label, value] (label)}
						<p class="m-0 text-sm text-muted-foreground">{label}: {value}</p>
					{/each}
				</div>
			</Card>
		{/each}
	</div>
</SettingsFormModal>
