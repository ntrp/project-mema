<script lang="ts">
	import { tick } from 'svelte';
	import InlineSpinner from '$lib/components/shared/InlineSpinner.svelte';
	import PosterPlaceholder from '$lib/components/app/media/posters/PosterPlaceholder.svelte';
	import { imageUrl, mediaResultKey } from '$lib/components/app/search/advancedSearchResults';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import type { LibraryScanItem, MediaSearchResult } from '$lib/settings/types';
	import {
		cleanMatchSearchTitle,
		type MatchDraft
	} from '$lib/components/settings/library/scan/libraryScanImport';

	interface Props {
		item: LibraryScanItem;
		draft: MatchDraft;
		onSearch: (_item: LibraryScanItem) => void;
		onSelect: (_item: LibraryScanItem, _result: MediaSearchResult) => void;
	}

	let { item, draft = $bindable(), onSearch, onSelect }: Props = $props();
	let editing = $state(false);
	let inputEl = $state<HTMLInputElement | null>(null);
	let queryBeforeEdit = '';

	const badgeLabel = $derived(
		draft.matched ? `${draft.matched.title}${draft.matched.year ? ` (${draft.matched.year})` : ''}` : 'No match'
	);

	async function openSearch() {
		queryBeforeEdit = draft.query;
		draft.query = cleanMatchSearchTitle(
			draft.matched?.title ?? item.matchedTitle ?? item.detectedTitle ?? item.fileName ?? draft.query
		);
		draft.results = [];
		draft.searched = false;
		editing = true;
		await tick();
		inputEl?.focus();
		inputEl?.select();
	}

	function closeSearch() {
		draft.query = '';
		onSearch(item);
		draft.query = queryBeforeEdit;
		draft.results = [];
		draft.searching = false;
		draft.searched = Boolean(draft.matched);
		editing = false;
	}

	function handleFocusOut(event: FocusEvent) {
		const nextTarget = event.relatedTarget as Node | null;
		if (nextTarget && (event.currentTarget as HTMLElement).contains(nextTarget)) return;
		closeSearch();
	}

	function handleInput() {
		onSearch(item);
	}

	function choose(result: MediaSearchResult) {
		onSelect(item, result);
		editing = false;
	}

	function subtitle(result: MediaSearchResult) {
		if (result.overview) return result.overview;
		return [result.type, result.externalProvider].filter(Boolean).join(' · ');
	}
</script>

{#if !editing}
	<Button
		type="button"
		variant="secondary"
		class={`${draft.matched ? '' : 'border-amber-500/40 bg-amber-500/10 text-amber-500 hover:bg-amber-500/15 hover:text-amber-500 '}h-9 max-w-96 justify-start truncate`}
		onclick={openSearch}
	>
		<span class="truncate">{badgeLabel}</span>
	</Button>
{:else}
	<div class="relative inline-block align-top" onfocusout={handleFocusOut}>
		<div class="flex max-w-96 items-center gap-2">
			<Input
				bind:ref={inputEl}
				bind:value={draft.query}
				placeholder="Search manually"
				oninput={handleInput}
			/>
			{#if draft.searching}
				<InlineSpinner label="Searching" />
			{/if}
		</div>
		{#if draft.results.length}
			<div
				class="absolute top-full left-0 z-50 mt-2 grid w-[36rem] max-w-[calc(100vw-4rem)] rounded-md border border-border bg-popover p-1 shadow-xl"
			>
				{#each draft.results as result (mediaResultKey(result))}
					<Button
						type="button"
						variant="ghost"
						class="h-auto min-h-22 w-full items-start justify-start gap-3 whitespace-normal rounded-sm border-0 border-b border-border/60 bg-transparent p-2 text-left text-foreground shadow-none last:border-b-0"
						onclick={() => choose(result)}
					>
						<span class="h-18 aspect-[2/3] shrink-0 overflow-hidden rounded-sm bg-muted">
							{#if imageUrl(result.posterPath)}
								<img
									class="block size-full object-cover"
									src={imageUrl(result.posterPath)}
									alt=""
									loading="lazy"
								/>
							{:else}
								<PosterPlaceholder label={result.title} class="min-h-0 rounded-none p-1 text-[10px]" />
							{/if}
						</span>
						<span class="grid min-w-0 content-start gap-0.5 pt-0.5">
							<strong class="truncate text-sm">
								{result.title}{result.year ? ` (${result.year})` : ''}
							</strong>
							<span class="line-clamp-2 text-xs font-medium text-muted-foreground">
								{subtitle(result)}
							</span>
						</span>
					</Button>
				{/each}
			</div>
		{:else if draft.searched && !draft.searching}
			<div
				class="absolute top-full left-0 z-50 mt-2 w-72 rounded-md border border-border bg-popover p-3 text-xs text-muted-foreground shadow-xl"
			>
				No results
			</div>
		{/if}
	</div>
{/if}
