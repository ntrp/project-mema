<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		labelledby: string;
		class?: string;
		children: Snippet;
	}

	let { labelledby, class: className = '', children }: Props = $props();
</script>

<section
	class={`relative isolate -mx-6 -mb-10 grid min-h-[calc(100vh-76px)] px-6 pt-12 pb-10 max-[980px]:-mx-[18px] max-[980px]:-mb-10 max-[980px]:pt-9 ${className}`}
	aria-labelledby={labelledby}
>
	<div class="relative z-[1]">
		{@render children()}
	</div>
</section>
