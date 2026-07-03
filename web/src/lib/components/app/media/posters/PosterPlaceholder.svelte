<script lang="ts">
	import { cn } from '$lib/utils';

	interface Props {
		label: string;
		class?: string;
	}

	let { label, class: className = '' }: Props = $props();
</script>

<div
	class={cn(
		'grid min-h-[150px] place-items-center rounded-md bg-card font-black text-foreground',
		className
	)}
>
	{label}
</div>
