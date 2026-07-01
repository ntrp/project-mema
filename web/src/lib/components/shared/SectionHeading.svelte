<script lang="ts">
	import type { Snippet } from 'svelte';
	import { cn } from '$lib/utils';

	interface Props {
		title?: string;
		titleId?: string;
		kicker?: string;
		href?: string;
		titleContent?: Snippet;
		actions?: Snippet;
		children?: Snippet;
		class?: string;
	}

	let {
		title,
		titleId,
		kicker,
		href,
		titleContent,
		actions,
		children,
		class: className
	}: Props = $props();
</script>

<div class={cn('mb-4 flex items-center justify-between gap-3', className)}>
	<div class="min-w-0">
		{#if titleContent}
			{@render titleContent()}
		{:else if title}
			{#if kicker}
				<p class="m-0 mb-1 text-xs font-extrabold tracking-normal text-muted-foreground uppercase">
					{kicker}
				</p>
			{/if}
			{#if href}
				<!-- eslint-disable svelte/no-navigation-without-resolve -->
				<a
					class="inline-flex items-center gap-2 text-foreground no-underline hover:text-primary-hover focus-visible:text-primary-hover focus-visible:outline-none"
					{href}
				>
					<h2 id={titleId} class="m-0 text-3xl font-semibold text-foreground">{title}</h2>
					{@render children?.()}
				</a>
				<!-- eslint-enable svelte/no-navigation-without-resolve -->
			{:else}
				<h2 id={titleId} class="m-0 text-3xl font-semibold text-foreground">{title}</h2>
			{/if}
		{/if}
	</div>
	{#if actions}
		<div class="inline-flex items-center gap-2 text-xs font-black text-muted-foreground">
			{@render actions()}
		</div>
	{/if}
</div>
