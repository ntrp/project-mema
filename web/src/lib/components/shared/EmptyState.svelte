<script lang="ts">
	import { Card } from '$lib/components/ui/card';
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	interface Props extends HTMLAttributes<HTMLDivElement> {
		title?: string;
		description?: string;
		children?: Snippet;
		class?: string;
	}

	let { title, description, children, class: className = '', ...restProps }: Props = $props();
</script>

<Card
	class={cn(
		'gap-2 rounded-md border-dashed border-border bg-muted p-[18px] shadow-none',
		className
	)}
	{...restProps}
>
	{#if title}
		<h2 class="m-0 text-2xl leading-tight font-extrabold text-foreground">{title}</h2>
	{/if}
	{#if description}
		<p class={cn('m-0 text-sm text-muted-foreground', title ? 'mt-2' : '')}>{description}</p>
	{/if}
	{@render children?.()}
</Card>
