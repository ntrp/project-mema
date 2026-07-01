<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';

	type Tone = 'success' | 'muted' | 'pending' | 'error';

	interface Props {
		tone?: Tone;
		children: Snippet;
		class?: string;
	}

	let { tone = 'muted', children, class: className = '' }: Props = $props();

	const toneClass = $derived(
		tone === 'success'
			? 'bg-primary/10 text-primary'
			: tone === 'pending'
				? 'bg-secondary text-secondary-foreground'
				: tone === 'error'
					? 'bg-destructive/10 text-destructive'
					: 'bg-muted text-muted-foreground'
	);
</script>

<Badge
	variant="secondary"
	class={cn('h-auto rounded-md px-2 py-1 leading-tight font-extrabold', toneClass, className)}
>
	{@render children()}
</Badge>
