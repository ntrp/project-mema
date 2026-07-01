<script lang="ts">
	import InlineSpinner from '$lib/components/shared/InlineSpinner.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
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
	<InlineSpinner label="Searching" />
{:else if draft.matched}
	<Button
		type="button"
		variant="secondary"
		class="h-auto whitespace-normal"
		onclick={() => (draft.matched = undefined)}
	>
		{draft.matched.title}{draft.matched.year ? ` (${draft.matched.year})` : ''}
	</Button>
{:else}
	<p class="m-0 font-bold text-destructive">No match found!</p>
	<Input bind:value={draft.query} oninput={() => onSearch(item)} />
	{#if draft.results.length}
		<div class="grid max-w-[520px] gap-1.5">
			{#each draft.results as result (`${result.type}:${result.title}:${result.year ?? ''}`)}
				<Button
					type="button"
					variant="ghost"
					class="h-auto w-full justify-between gap-3 border border-border bg-card text-left text-foreground"
					onclick={() => onSelect(item, result)}
				>
					<strong>{result.title}</strong>
					<span class="font-bold text-muted-foreground">
						{result.type}{result.year ? ` · ${result.year}` : ''}
					</span>
				</Button>
			{/each}
		</div>
	{/if}
{/if}
