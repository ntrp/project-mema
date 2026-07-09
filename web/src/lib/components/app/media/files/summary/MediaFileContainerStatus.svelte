<script lang="ts">
	import PackageIcon from '@lucide/svelte/icons/package';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import { cn } from '$lib/utils';

	interface Props {
		row: MediaFileRow;
	}

	let { row }: Props = $props();

	const status = $derived(row.requirements?.container);
	const badgeClass = $derived(
		cn(
			'inline-flex items-center gap-1 rounded-sm border px-1.5 py-0.5 text-[11px] font-medium',
			status?.state === 'satisfied' && 'border-emerald-300 bg-emerald-50 text-emerald-700',
			status?.state === 'pending' && 'border-sky-300 bg-sky-50 text-sky-700',
			status?.state === 'partial' && 'border-amber-300 bg-amber-50 text-amber-700',
			status?.state === 'missing' && 'border-destructive/40 bg-destructive/10 text-destructive',
			status?.state === 'ignored' && 'border-border bg-muted/50 text-muted-foreground'
		)
	);

	function stopClick(event: Event) {
		event.stopPropagation();
	}
</script>

{#if status && status.state !== 'satisfied' && status.state !== 'ignored'}
	<span class="inline-flex items-center gap-1">
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#snippet child({ props })}
					<button {...props} type="button" class={badgeClass} onclick={stopClick}>
						<PackageIcon class="size-3" aria-hidden="true" />
						<span>{status.label}</span>
					</button>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content class="max-w-96">
				<div class="grid gap-1">
					<strong>Container {status.label}</strong>
					<ul class="list-disc space-y-1 pl-4">
						{#each status.details as detail (detail)}
							<li>{detail}</li>
						{/each}
					</ul>
				</div>
			</Tooltip.Content>
		</Tooltip.Root>
	</span>
{/if}
