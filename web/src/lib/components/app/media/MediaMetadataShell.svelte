<script lang="ts">
	import type { Snippet } from 'svelte';
	import { imageUrl } from './mediaDetail';

	interface Props {
		backdropPath?: string;
		labelledby: string;
		class?: string;
		children: Snippet;
	}

	let { backdropPath, labelledby, class: className = '', children }: Props = $props();

	const backdropUrl = $derived(imageUrl(backdropPath, 'original'));
</script>

<section
	class={`relative -mx-6 -mb-10 grid min-h-[calc(100vh-76px)] overflow-hidden bg-background px-6 pt-12 pb-10 max-[980px]:-mx-[18px] max-[980px]:-mb-10 max-[980px]:pt-9 ${className}`}
	aria-labelledby={labelledby}
>
	{#if backdropUrl}
		<img
			class="pointer-events-none absolute inset-0 z-0 size-full object-cover opacity-25"
			src={backdropUrl}
			alt=""
		/>
		<div
			class="pointer-events-none absolute inset-0 z-0 bg-linear-to-r from-background via-background/80 to-background/25"
		></div>
	{/if}
	<div class="relative z-[1]">
		{@render children()}
	</div>
</section>
