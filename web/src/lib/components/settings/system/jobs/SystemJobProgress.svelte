<script lang="ts">
	import { progressStyle } from './systemJobDisplay';

	interface Props {
		value?: number;
		label: string;
	}

	let { value, label }: Props = $props();
</script>

<div class="grid min-w-36 gap-1">
	<div class="h-1.5 overflow-hidden rounded-full bg-muted">
		{#if value !== undefined}
			<div class="h-full bg-primary transition-all" style={progressStyle(value)}></div>
		{:else}
			<div class="h-full w-1/3 animate-pulse bg-primary/60"></div>
		{/if}
	</div>
	<span class="block truncate text-xs text-muted-foreground">{label}</span>
</div>
