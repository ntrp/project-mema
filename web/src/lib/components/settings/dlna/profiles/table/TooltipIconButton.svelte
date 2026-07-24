<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';

	interface Props {
		label: string;
		tooltip: string;
		onclick: () => void;
		variant?: 'default' | 'outline';
		children: Snippet;
	}

	let { label, tooltip, onclick, variant = 'outline', children }: Props = $props();
</script>

<Tooltip.Root>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<Button {...props} type="button" {variant} size="icon-sm" aria-label={label} {onclick}>
				{@render children()}
			</Button>
		{/snippet}
	</Tooltip.Trigger>
	<Tooltip.Content>{tooltip}</Tooltip.Content>
</Tooltip.Root>
