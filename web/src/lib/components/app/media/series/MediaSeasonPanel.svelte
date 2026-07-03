<script lang="ts">
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import { Card } from '$lib/components/ui/card';
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';

	interface Props {
		summary: string;
		size?: string;
		tone: 'success' | 'active' | 'missing' | 'neutral';
		title: Snippet;
		actions?: Snippet;
		children: Snippet;
	}

	let { summary, size, tone, title, actions, children }: Props = $props();
	let open = $state(false);

	function toggle() {
		open = !open;
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key !== 'Enter' && event.key !== ' ') return;
		event.preventDefault();
		toggle();
	}

	function toneClass(value: Props['tone']) {
		switch (value) {
			case 'success':
				return 'bg-green-600 text-white';
			case 'active':
				return 'bg-primary text-primary-foreground';
			case 'missing':
				return 'bg-destructive text-destructive-foreground';
			default:
				return 'bg-border text-muted-foreground';
		}
	}
</script>

<Card class="overflow-hidden bg-sidebar/80 p-0">
	<div
		role="button"
		tabindex="0"
		aria-expanded={open}
		class="grid min-h-12 cursor-pointer grid-cols-[minmax(0,1fr)_2rem_minmax(88px,auto)] items-center gap-3 bg-muted/70 px-4 py-3 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/50"
		onclick={toggle}
		onkeydown={handleKeydown}
	>
		<span class="inline-flex min-w-0 items-center gap-2.5 font-black">
			{@render title()}
			<span class={cn('rounded-md px-2 py-0.5 text-xs font-black', toneClass(tone))}>
				{summary}
			</span>
			{#if size}
				<span class="rounded-md bg-border px-2 py-0.5 text-xs font-black text-muted-foreground">
					{size}
				</span>
			{/if}
		</span>
		<span
			class="grid size-8 place-items-center self-center justify-self-center text-muted-foreground"
			aria-hidden="true"
		>
			<ChevronDownIcon
				class={open ? 'size-5 rotate-180 transition-transform' : 'size-5 transition-transform'}
			/>
		</span>
		<span class="flex min-w-0 items-center justify-end gap-2">
			{@render actions?.()}
		</span>
	</div>
	{#if open}
		<div class="border-t border-border">
			{@render children()}
		</div>
	{/if}
</Card>
