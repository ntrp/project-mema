<script lang="ts">
	import BookmarkIcon from '@lucide/svelte/icons/bookmark';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { cn } from '$lib/utils';

	interface Props {
		monitored?: boolean;
		status: string;
		hint: string;
		disabled?: boolean;
		size?: number;
		onToggle: () => void;
	}

	let { monitored = false, status, hint, disabled = false, size = 32, onToggle }: Props = $props();

	function handleClick(event: globalThis.MouseEvent) {
		event.stopPropagation();
		onToggle();
	}

	function handleKeydown(event: globalThis.KeyboardEvent) {
		if (event.key !== 'Enter' && event.key !== ' ') return;
		event.stopPropagation();
	}
</script>

<Tooltip.Root>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<button
				{...props}
				type="button"
				class={cn(
					'inline-flex shrink-0 items-center justify-center rounded-md leading-none transition-colors hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-60',
					monitored ? 'text-secondary-foreground' : 'text-muted-foreground'
				)}
				aria-label={`${status}. ${hint}`}
				{disabled}
				onclick={handleClick}
				onkeydown={handleKeydown}
			>
				<BookmarkIcon aria-hidden="true" {size} class={monitored ? 'fill-current' : undefined} />
			</button>
		{/snippet}
	</Tooltip.Trigger>
	<Tooltip.Content>
		<span class="grid gap-0.5">
			<strong>{status}</strong>
			<span>{hint}</span>
		</span>
	</Tooltip.Content>
</Tooltip.Root>
