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
		'flex size-full min-h-[150px] items-center justify-center rounded-md bg-card p-3 text-center font-black text-foreground',
		className
	)}
>
	<span class="max-w-full break-words">{label}</span>
</div>
