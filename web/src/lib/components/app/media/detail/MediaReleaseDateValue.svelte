<script lang="ts">
	import ClapperboardIcon from '@lucide/svelte/icons/clapperboard';
	import DiscIcon from '@lucide/svelte/icons/disc-3';
	import MonitorPlayIcon from '@lucide/svelte/icons/monitor-play';
	import * as Tooltip from '$lib/components/ui/tooltip';

	interface Props {
		kind: 'cinema' | 'digital' | 'physical';
		label: string;
		date: string;
	}

	let { kind, label, date }: Props = $props();
	let open = $state(false);

	const tooltipLabel = $derived(`${label} release date`);
	const Icon = $derived(releaseDateIcon(kind));

	function releaseDateIcon(kind: Props['kind']) {
		switch (kind) {
			case 'cinema':
				return ClapperboardIcon;
			case 'physical':
				return DiscIcon;
			default:
				return MonitorPlayIcon;
		}
	}
</script>

<span class="inline-flex items-center justify-end gap-1.5">
	<Tooltip.Root bind:open>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<button
					{...props}
					type="button"
					class="inline-flex items-center rounded-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/50"
					aria-label={tooltipLabel}
				>
					<Icon aria-hidden="true" class="size-3.5 text-foreground" />
				</button>
			{/snippet}
		</Tooltip.Trigger>
		{#if open}
			<Tooltip.Content>{tooltipLabel}</Tooltip.Content>
		{/if}
	</Tooltip.Root>
	<span class="sr-only">{tooltipLabel}</span>
	<span>{date}</span>
</span>
