<script lang="ts">
	import type { Snippet } from 'svelte';
	import { cn } from '$lib/utils';

	interface Props {
		eyebrow: string;
		title: string;
		titleId?: string;
		description?: string;
		actions?: Snippet;
		children?: Snippet;
		class?: string;
	}

	let {
		eyebrow,
		title,
		titleId,
		description,
		actions,
		children,
		class: className
	}: Props = $props();
</script>

<div
	class={cn(
		'relative mb-5 grid gap-1',
		actions && 'grid-cols-[minmax(0,1fr)_auto] items-start gap-3',
		className
	)}
>
	<div class="min-w-0">
		<p class="m-0 text-xs font-extrabold tracking-normal text-muted-foreground uppercase">
			{eyebrow}
		</p>
		<div class="flex items-center gap-2">
			<h1 id={titleId} class="text-4xl leading-tight font-semibold text-foreground">{title}</h1>
			{@render children?.()}
		</div>
		{#if description}
			<p
				class="mt-2 max-w-none text-[15px] leading-relaxed font-normal text-muted-foreground normal-case"
			>
				{description}
			</p>
		{/if}
	</div>
	{#if actions}
		<div class="relative z-20 justify-self-end">
			{@render actions()}
		</div>
	{/if}
</div>
