<script lang="ts">
	import ChevronLeftIcon from '@lucide/svelte/icons/chevron-left';
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
	import { Button } from '$lib/components/ui/button';

	interface Props {
		ariaLabel: string;
		leftLabel?: string;
		rightLabel?: string;
		onScroll: (_direction: -1 | 1) => void;
	}

	let {
		ariaLabel,
		leftLabel = 'Scroll left',
		rightLabel = 'Scroll right',
		onScroll
	}: Props = $props();
</script>

<div class="hidden items-center gap-2 sm:inline-flex" aria-label={ariaLabel}>
	<Button
		type="button"
		variant="ghost"
		size="icon-sm"
		class="size-[34px] min-h-[34px] border border-border bg-card p-0 text-muted-foreground hover:border-primary/50 hover:bg-muted hover:text-primary-foreground"
		aria-label={leftLabel}
		onclick={() => onScroll(-1)}
	>
		<ChevronLeftIcon aria-hidden="true" />
	</Button>
	<Button
		type="button"
		variant="ghost"
		size="icon-sm"
		class="size-[34px] min-h-[34px] border border-border bg-card p-0 text-muted-foreground hover:border-primary/50 hover:bg-muted hover:text-primary-foreground"
		aria-label={rightLabel}
		onclick={() => onScroll(1)}
	>
		<ChevronRightIcon aria-hidden="true" />
	</Button>
</div>
