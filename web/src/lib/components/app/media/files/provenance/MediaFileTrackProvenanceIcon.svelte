<script lang="ts">
	import InfoIcon from '@lucide/svelte/icons/info';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils';
	import {
		provenanceFields,
		type MediaFileTrackProvenance
	} from '$lib/components/app/media/files/provenance/mediaFileTrackProvenance';

	interface Props {
		provenance?: MediaFileTrackProvenance;
	}

	let { provenance }: Props = $props();
	const fields = $derived(provenance ? provenanceFields(provenance) : []);
	const hasProvenance = $derived(fields.length > 0);
</script>

<Tooltip.Root>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<button
				{...props}
				type="button"
				class={cn(
					'inline-flex h-6 w-6 items-center justify-center rounded-md hover:bg-accent',
					hasProvenance ? 'text-sky-600 dark:text-sky-300' : 'text-muted-foreground'
				)}
				aria-label="Track provenance"
			>
				<InfoIcon class="size-4" aria-hidden="true" />
			</button>
		{/snippet}
	</Tooltip.Trigger>
	<Tooltip.Content class="max-h-[min(520px,calc(100vh-96px))] max-w-128 overflow-auto">
		<div class="grid gap-2 text-left">
			<span class="font-bold">Provenance</span>
			{#if hasProvenance}
				{#each fields as field (field.label)}
					<div class="grid grid-cols-[104px_minmax(0,1fr)] gap-3">
						<span class="text-muted-foreground">{field.label}</span>
						{#if field.multiline}
							<pre class="break-anywhere whitespace-pre-wrap font-mono text-xs">{field.value}</pre>
						{:else}
							<span class="break-anywhere font-mono text-xs">{field.value}</span>
						{/if}
					</div>
				{/each}
			{:else}
				<span class="text-muted-foreground">No provenance recorded for this track.</span>
			{/if}
		</div>
	</Tooltip.Content>
</Tooltip.Root>
