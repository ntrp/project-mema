<script lang="ts">
	import { resolve } from '$app/paths';
	import StatusPill from '$lib/components/shared/StatusPill.svelte';

	interface Props {
		keywords?: string[];
	}

	let { keywords = [] }: Props = $props();

	const interactiveChipClass =
		'transition-[background-color,color,box-shadow] group-hover/chip:bg-primary group-hover/chip:text-primary-foreground group-focus-visible/chip:bg-primary group-focus-visible/chip:text-primary-foreground group-focus-visible/chip:ring-2 group-focus-visible/chip:ring-ring';
</script>

{#if keywords.length > 0}
	<div class="grid gap-0" aria-labelledby="metadata-keywords-title">
		<h3 id="metadata-keywords-title" class="mt-3 mb-4 text-xl text-foreground">Keywords</h3>
		<div class="flex flex-wrap gap-[7px]" aria-label="Keywords">
			{#each keywords as keyword (keyword)}
				<a
					class="group/chip rounded-[3px] text-foreground no-underline outline-none"
					href={`${resolve('/discover/movies')}?keywords=${encodeURIComponent(keyword)}`}
				>
					<StatusPill class={interactiveChipClass}>{keyword}</StatusPill>
				</a>
			{/each}
		</div>
	</div>
{/if}
