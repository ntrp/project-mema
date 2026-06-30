<script lang="ts">
	import type { LibraryScanItem, MediaSearchResult } from '$lib/settings/types';
	import type { MatchDraft } from '$lib/components/settings/libraryScanImport';

	interface Props {
		item: LibraryScanItem;
		draft: MatchDraft;
		onSearch: (_item: LibraryScanItem) => void;
		onSelect: (_item: LibraryScanItem, _result: MediaSearchResult) => void;
	}

	let { item, draft, onSearch, onSelect }: Props = $props();
</script>

{#if draft.searching}
	<span class="inline-spinner">Searching</span>
{:else if draft.matched}
	<button type="button" class="match-pill" onclick={() => (draft.matched = undefined)}>
		{draft.matched.title}{draft.matched.year ? ` (${draft.matched.year})` : ''}
	</button>
{:else}
	<p class="error">No match found!</p>
	<input bind:value={draft.query} oninput={() => onSearch(item)} />
	{#if draft.results.length}
		<div class="autocomplete-list compact">
			{#each draft.results as result (`${result.type}:${result.title}:${result.year ?? ''}`)}
				<button type="button" onclick={() => onSelect(item, result)}>
					<strong>{result.title}</strong>
					<span>{result.type}{result.year ? ` · ${result.year}` : ''}</span>
				</button>
			{/each}
		</div>
	{/if}
{/if}
