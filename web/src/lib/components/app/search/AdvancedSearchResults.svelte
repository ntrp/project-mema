<script lang="ts">
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import type { MediaSearchGroup, MediaSearchResult } from '$lib/settings/types';
	import AdvancedMediaResultItem from './AdvancedMediaResultItem.svelte';
	import AdvancedPersonResultItem from './AdvancedPersonResultItem.svelte';
	import { mediaResultKey, personResultKey } from './advancedSearchResults';

	interface Props {
		groups: MediaSearchGroup[];
		addingKey?: string;
		actionLabel: string;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let { groups, addingKey, actionLabel, onAdd }: Props = $props();

	function groupDomId(group: MediaSearchGroup) {
		return `advanced-${group.sourceType}-${group.sourceName.toLowerCase().replace(/[^a-z0-9]+/g, '-')}`;
	}
</script>

<div class="grid gap-[22px]" aria-label="Advanced search results">
	{#each groups as group (`${group.sourceType}:${group.sourceName}`)}
		{#if group.results.length > 0 || (group.people?.length ?? 0) > 0}
			{@const headingId = groupDomId(group)}
			<section aria-labelledby={headingId}>
				<SectionHeading title={group.sourceName} titleId={headingId}>
					{#snippet actions()}
						<span>{group.sourceType}</span>
					{/snippet}
				</SectionHeading>
				<div class="grid gap-2.5">
					{#each group.results as result (mediaResultKey(result))}
						<AdvancedMediaResultItem
							{result}
							inLibrary={group.sourceType === 'library'}
							{addingKey}
							{actionLabel}
							{onAdd}
						/>
					{/each}
					{#each group.people ?? [] as person (personResultKey(person))}
						<AdvancedPersonResultItem {person} />
					{/each}
				</div>
			</section>
		{/if}
	{/each}
</div>
