<script lang="ts">
	import { onMount } from 'svelte';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils';
	import type { CustomFormatSpec } from '$lib/settings/types';
	import { customFormatSpecTone, customFormatSpecToneClass } from './customFormatSpecTone';

	interface Props {
		spec: CustomFormatSpec;
		negated?: boolean;
		class?: string;
	}

	let { spec, negated = false, class: className = '' }: Props = $props();
	let chipElement = $state<globalThis.HTMLSpanElement>();
	let tooltipReady = $state(false);
	const label = $derived(`${spec.type}: ${spec.value}`);
	const toneClass = $derived(
		customFormatSpecToneClass(customFormatSpecTone(spec, negated || className.includes('exclude')))
	);
	const chipClass = $derived(
		cn(
			'inline-flex min-h-5 items-center rounded-md px-1.5 text-[11px] font-extrabold',
			toneClass,
			className
		)
	);

	onMount(() => {
		if (!chipElement || !('IntersectionObserver' in globalThis)) {
			tooltipReady = true;
			return;
		}

		const observer = new globalThis.IntersectionObserver(
			(entries) => {
				if (entries.some((entry) => entry.isIntersecting)) {
					tooltipReady = true;
					observer.disconnect();
				}
			},
			{ rootMargin: '80px' }
		);
		observer.observe(chipElement);

		return () => observer.disconnect();
	});
</script>

{#if tooltipReady}
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<span bind:this={chipElement} class={chipClass} {...props}>
					{spec.name}
				</span>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>{label}</Tooltip.Content>
	</Tooltip.Root>
{:else}
	<span bind:this={chipElement} class={chipClass}>
		{spec.name}
	</span>
{/if}
