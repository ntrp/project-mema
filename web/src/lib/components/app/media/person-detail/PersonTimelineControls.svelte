<script lang="ts">
	import ChevronLeftIcon from '@lucide/svelte/icons/chevron-left';
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
	import ChevronsLeftIcon from '@lucide/svelte/icons/chevrons-left';
	import ChevronsRightIcon from '@lucide/svelte/icons/chevrons-right';
	import { Button } from '$lib/components/ui/button';

	interface Props {
		onStep: (_direction: -1 | 1) => void;
		onEdge: (_direction: -1 | 1) => void;
	}

	let { onStep, onEdge }: Props = $props();
</script>

<div class="flex justify-center gap-2">
	<Button
		type="button"
		variant="outline"
		size="icon-sm"
		aria-label="Scroll timeline to beginning"
		onclick={() => onEdge(-1)}
	>
		<ChevronsLeftIcon aria-hidden="true" />
	</Button>
	<Button
		type="button"
		variant="outline"
		size="icon-sm"
		aria-label="Scroll timeline left"
		onclick={() => onStep(-1)}
	>
		<ChevronLeftIcon aria-hidden="true" />
	</Button>
	<Button
		type="button"
		variant="outline"
		size="icon-sm"
		aria-label="Scroll timeline right"
		onclick={() => onStep(1)}
	>
		<ChevronRightIcon aria-hidden="true" />
	</Button>
	<Button
		type="button"
		variant="outline"
		size="icon-sm"
		aria-label="Scroll timeline to end"
		onclick={() => onEdge(1)}
	>
		<ChevronsRightIcon aria-hidden="true" />
	</Button>
</div>
