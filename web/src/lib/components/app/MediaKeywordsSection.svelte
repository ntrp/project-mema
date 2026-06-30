<script lang="ts">
	interface Props {
		keywords?: string[];
	}

	let { keywords = [] }: Props = $props();
</script>

{#if keywords.length > 0}
	<div class="metadata-keywords-block" aria-labelledby="metadata-keywords-title">
		<h3 id="metadata-keywords-title" class="metadata-keywords-title">Keywords</h3>
		<div class="metadata-tags" aria-label="Keywords">
			{#each keywords as keyword (keyword)}
				<span>{keyword}</span>
			{/each}
		</div>
	</div>
{/if}
