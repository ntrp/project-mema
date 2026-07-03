<script lang="ts">
	import ChevronLeftIcon from '@lucide/svelte/icons/chevron-left';
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
	import { Button } from '$lib/components/ui/button';

	interface Props {
		ariaLabel: string;
		leftLabel?: string;
		rightLabel?: string;
		canScrollLeft?: boolean;
		canScrollRight?: boolean;
		onScroll: (_direction: -1 | 1) => void;
	}

	let {
		ariaLabel,
		leftLabel = 'Scroll left',
		rightLabel = 'Scroll right',
		canScrollLeft = true,
		canScrollRight = true,
		onScroll
	}: Props = $props();
</script>

<div class="hidden items-center gap-2 sm:inline-flex" aria-label={ariaLabel}>
	<Button
		type="button"
		variant="ghost"
		size="icon-sm"
		class="size-[34px] min-h-[34px] border border-border bg-card p-0 text-muted-foreground hover:border-primary/50 hover:bg-muted hover:text-primary-foreground disabled:cursor-default disabled:text-muted-foreground/30 disabled:hover:border-border disabled:hover:bg-card disabled:hover:text-muted-foreground/30"
		aria-label={leftLabel}
		disabled={!canScrollLeft}
		onclick={() => onScroll(-1)}
	>
		<ChevronLeftIcon aria-hidden="true" />
	</Button>
	<Button
		type="button"
		variant="ghost"
		size="icon-sm"
		class="size-[34px] min-h-[34px] border border-border bg-card p-0 text-muted-foreground hover:border-primary/50 hover:bg-muted hover:text-primary-foreground disabled:cursor-default disabled:text-muted-foreground/30 disabled:hover:border-border disabled:hover:bg-card disabled:hover:text-muted-foreground/30"
		aria-label={rightLabel}
		disabled={!canScrollRight}
		onclick={() => onScroll(1)}
	>
		<ChevronRightIcon aria-hidden="true" />
	</Button>
</div>
