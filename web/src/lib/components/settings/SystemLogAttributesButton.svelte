<script lang="ts">
	import BracesIcon from '@lucide/svelte/icons/braces';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	interface Props {
		attributes: Record<string, unknown>;
	}

	let { attributes }: Props = $props();

	const formatted = $derived(JSON.stringify(attributes, null, 2));
</script>

<Tooltip.Root>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<Button type="button" variant="ghost" size="icon-xs" aria-label="Log attributes" {...props}>
				<BracesIcon aria-hidden="true" />
			</Button>
		{/snippet}
	</Tooltip.Trigger>
	<Tooltip.Content class="max-h-80 max-w-[min(560px,70vw)] overflow-auto">
		<code class="whitespace-pre text-left text-xs">{formatted}</code>
	</Tooltip.Content>
</Tooltip.Root>
