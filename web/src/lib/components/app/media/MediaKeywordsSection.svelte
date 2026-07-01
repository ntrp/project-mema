<script lang="ts">
	import StatusPill from '$lib/components/shared/StatusPill.svelte';

	interface Props {
		keywords?: string[];
	}

	let { keywords = [] }: Props = $props();
</script>

{#if keywords.length > 0}
	<div class="grid gap-0" aria-labelledby="metadata-keywords-title">
		<h3 id="metadata-keywords-title" class="mt-3 mb-4 text-xl text-foreground">Keywords</h3>
		<div class="flex flex-wrap gap-[7px]" aria-label="Keywords">
			{#each keywords as keyword (keyword)}
				<StatusPill>{keyword}</StatusPill>
			{/each}
		</div>
	</div>
{/if}
