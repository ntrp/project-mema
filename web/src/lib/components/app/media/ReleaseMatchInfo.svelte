<script lang="ts">
	import CircleAlertIcon from '@lucide/svelte/icons/circle-alert';
	import CircleXIcon from '@lucide/svelte/icons/circle-x';
	import InfoIcon from '@lucide/svelte/icons/info';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MatchInfo } from './releaseCandidateDisplay';

	interface Props {
		info: MatchInfo;
	}

	let { info }: Props = $props();

	const label = $derived(
		info.severity === 'error'
			? 'Release mismatch'
			: info.severity === 'warning'
				? 'Release warning'
				: 'Release match'
	);
	const iconClass = $derived(
		info.severity === 'error'
			? 'text-destructive'
			: info.severity === 'warning'
				? 'text-amber-500'
				: 'text-sky-500'
	);
</script>

<Tooltip.Root>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<button
				{...props}
				type="button"
				class="inline-flex h-8 w-8 items-center justify-center rounded-md hover:bg-accent"
				aria-label={label}
			>
				{#if info.severity === 'error'}
					<CircleXIcon class={iconClass} aria-hidden="true" />
				{:else if info.severity === 'warning'}
					<CircleAlertIcon class={iconClass} aria-hidden="true" />
				{:else}
					<InfoIcon class={iconClass} aria-hidden="true" />
				{/if}
			</button>
		{/snippet}
	</Tooltip.Trigger>
	<Tooltip.Content class="max-w-80">
		<div class="grid gap-1 text-left">
			{#each info.details as detail (detail)}
				<span>{detail}</span>
			{/each}
		</div>
	</Tooltip.Content>
</Tooltip.Root>
