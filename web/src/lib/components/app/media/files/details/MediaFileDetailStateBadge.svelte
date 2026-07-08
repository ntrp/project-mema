<script lang="ts">
	import CheckCircleIcon from '@lucide/svelte/icons/check-circle';
	import ClockIcon from '@lucide/svelte/icons/clock';
	import CircleDashedIcon from '@lucide/svelte/icons/circle-dashed';
	import TriangleAlertIcon from '@lucide/svelte/icons/triangle-alert';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils';
	import { unwantedMediaBadgeClass } from './mediaFileVisualClasses';
	import type { MediaFileDetailRow } from '$lib/components/app/media/files/mediaFileDetails';

	interface Props {
		row: MediaFileDetailRow;
	}

	let { row }: Props = $props();

	const stateClass = $derived(
		cn(
			'inline-flex items-center gap-1 rounded-sm border px-1.5 py-0.5 text-[11px] font-medium',
			row.visualState === 'matching' && 'border-emerald-300 bg-emerald-50 text-emerald-700',
			row.visualState === 'partial' && 'border-amber-300 bg-amber-50 text-amber-700',
			row.visualState === 'unwanted' && unwantedMediaBadgeClass,
			row.visualState === 'pending_operation' && 'border-sky-300 bg-sky-50 text-sky-700',
			row.visualState === 'missing_placeholder' &&
				'border-destructive/40 bg-destructive/10 text-destructive'
		)
	);
</script>

{#if row.visualState}
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<button {...props} type="button" class={stateClass}>
					{#if row.visualState === 'matching'}
						<CheckCircleIcon class="size-3" aria-hidden="true" />
					{:else if row.visualState === 'pending_operation'}
						<ClockIcon class="size-3" aria-hidden="true" />
					{:else if row.visualState === 'missing_placeholder'}
						<CircleDashedIcon class="size-3" aria-hidden="true" />
					{:else}
						<TriangleAlertIcon class="size-3" aria-hidden="true" />
					{/if}
					<span>{row.statusLabel}</span>
				</button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content class="max-w-96">
			<div class="grid gap-1">
				{#if row.operationLabel}
					<strong>{row.operationLabel}</strong>
				{/if}
				<ul class="list-disc space-y-1 pl-4">
					{#each row.details ?? [] as detail (detail)}
						<li>{detail}</li>
					{/each}
				</ul>
			</div>
		</Tooltip.Content>
	</Tooltip.Root>
{/if}
