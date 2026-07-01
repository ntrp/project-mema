<script lang="ts">
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import type { Snippet } from 'svelte';

	interface Props {
		meta: string;
		title: Snippet;
		children: Snippet;
	}

	let { meta, title, children }: Props = $props();
	let open = $state(false);
</script>

<Card class="overflow-hidden bg-sidebar/80 p-0">
	<Button
		type="button"
		variant="ghost"
		class="flex h-auto min-h-12 w-full items-center justify-between gap-3 rounded-none border-0 bg-muted/70 px-4 py-3 text-left font-black text-muted-foreground hover:bg-muted hover:text-foreground"
		aria-expanded={open}
		onclick={() => (open = !open)}
	>
		<span class="inline-flex min-w-0 items-center gap-2.5">{@render title()}</span>
		<span
			class="inline-flex items-center gap-2 rounded-md bg-border px-2 py-0.5 text-xs font-black text-muted-foreground"
		>
			{meta}
			<ChevronRightIcon
				aria-hidden="true"
				class={open ? 'rotate-90 transition-transform' : 'transition-transform'}
			/>
		</span>
	</Button>
	{#if open}
		<div class="border-t border-border">
			{@render children()}
		</div>
	{/if}
</Card>
