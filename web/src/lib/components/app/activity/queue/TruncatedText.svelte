<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip';

	interface Props {
		text: string;
		class?: string;
	}

	let { text, class: className = '' }: Props = $props();
	let open = $state(false);
</script>

<Tooltip.Root bind:open>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<span {...props} class={`block min-w-0 truncate ${className}`}>{text}</span>
		{/snippet}
	</Tooltip.Trigger>
	{#if open}
		<Tooltip.Content class="max-w-160 break-all">{text}</Tooltip.Content>
	{/if}
</Tooltip.Root>
