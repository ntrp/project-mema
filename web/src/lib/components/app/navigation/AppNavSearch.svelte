<script lang="ts">
	import SearchIcon from '@lucide/svelte/icons/search';
	import { resolve } from '$app/paths';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { cn } from '$lib/utils';
	import type { MediaSearchGroup, MediaSearchResult } from '$lib/settings/types';

	interface Props {
		searchQuery: string;
		groups: MediaSearchGroup[];
		loading: boolean;
		onSearch: (_query: string) => void | Promise<void>;
		onSelect: (_result: MediaSearchResult) => void;
		onAdvancedSearch: (_query: string) => void;
	}

	let {
		searchQuery = $bindable(),
		groups,
		loading,
		onSearch,
		onSelect,
		onAdvancedSearch
	}: Props = $props();
	let searchOpen = $state(false);
	let selectedIndex = $state(-1);

	const trimmedQuery = $derived(searchQuery.trim());
	const flatResults = $derived(groups.flatMap((group) => group.results));
	const resultCount = $derived(groups.reduce((count, group) => count + group.results.length, 0));
	const showSuggestions = $derived(searchOpen && trimmedQuery.length >= 2);
	const selectedResult = $derived(selectedIndex >= 0 ? flatResults[selectedIndex] : undefined);
	const selectedKey = $derived(selectedResult ? resultKey(selectedResult) : undefined);

	function selectResult(result: MediaSearchResult) {
		searchQuery = result.title;
		searchOpen = false;
		onSelect(result);
	}

	function resultKey(result: MediaSearchResult) {
		return `${result.id ?? ''}:${result.type}:${result.externalProvider ?? ''}:${result.externalId ?? ''}:${result.title}:${result.year ?? ''}`;
	}

	function resultIndex(result: MediaSearchResult) {
		const key = resultKey(result);
		return flatResults.findIndex((item) => resultKey(item) === key);
	}

	function handleSearchInput(event: Event) {
		selectedIndex = -1;
		const query = ((event.currentTarget as { value?: string } | null)?.value ?? '').trim();
		if (query.length >= 2) {
			void onSearch(query);
		}
	}

	function handleSearchKeydown(event: KeyboardEvent) {
		if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
			if (resultCount === 0) return;
			event.preventDefault();
			searchOpen = true;
			selectedIndex =
				event.key === 'ArrowDown'
					? Math.min(selectedIndex + 1, resultCount - 1)
					: Math.max(selectedIndex - 1, 0);
			return;
		}
		if (event.key === 'Enter') {
			handleEnter(event);
		} else if (event.key === 'Escape') {
			searchOpen = false;
			selectedIndex = -1;
		}
	}

	function handleEnter(event: KeyboardEvent) {
		if (trimmedQuery.length === 0) return;
		event.preventDefault();
		if (selectedResult) {
			selectResult(selectedResult);
			return;
		}
		searchOpen = false;
		onAdvancedSearch(trimmedQuery);
	}

	function posterUrl(path?: string) {
		if (!path) return undefined;
		if (path.startsWith('http://') || path.startsWith('https://')) return path;
		return `https://image.tmdb.org/t/p/w92${path}`;
	}
</script>

<div class="relative min-w-0">
	<label class="sr-only" for="global-search">Search</label>
	<SearchIcon
		aria-hidden="true"
		class="pointer-events-none absolute top-1/2 left-3 z-10 size-4 -translate-y-1/2 text-muted-foreground"
	/>
	<Input
		id="global-search"
		bind:value={searchQuery}
		class="bg-background pl-9"
		placeholder="Search Movies & TV"
		autocomplete="off"
		onfocus={() => (searchOpen = true)}
		oninput={handleSearchInput}
		onkeydown={handleSearchKeydown}
		onblur={() => {
			window.setTimeout(() => {
				searchOpen = false;
			}, 120);
		}}
	/>
	{#if showSuggestions}
		<div
			class="absolute inset-x-0 top-[calc(100%+6px)] z-20 max-h-[min(560px,calc(100vh-84px))] overflow-auto rounded-md border border-border bg-popover p-1.5 text-popover-foreground shadow-md"
			role="listbox"
			aria-label="Search suggestions"
		>
			{#if resultCount > 0}
				{#each groups as group (`${group.sourceType}:${group.sourceName}`)}
					{#if group.results.length > 0}
						<div class="grid gap-1 border-t border-border py-1 first:border-t-0">
							<div class="px-2 py-1 text-xs font-bold text-muted-foreground uppercase">
								{group.sourceName}
							</div>
							{#each group.results as result (resultKey(result))}
								{@const index = resultIndex(result)}
								<Button
									type="button"
									variant="ghost"
									role="option"
									aria-selected={selectedKey === resultKey(result)}
									class={cn(
										'grid min-h-14 w-full grid-cols-[34px_minmax(0,1fr)] items-center gap-2.5 rounded-sm border-0 bg-transparent px-2 py-1.5 text-left text-popover-foreground hover:bg-accent hover:text-accent-foreground',
										index === selectedIndex && 'bg-accent text-accent-foreground'
									)}
									onpointerdown={(event) => event.preventDefault()}
									onclick={() => selectResult(result)}
								>
									<div
										class="grid h-11 w-[34px] place-items-center overflow-hidden rounded-md border border-border bg-background text-[9px] font-bold text-muted-foreground uppercase"
									>
										{#if posterUrl(result.posterPath)}
											<img
												class="size-full object-cover"
												src={posterUrl(result.posterPath)}
												alt=""
												loading="lazy"
											/>
										{:else}
											<span>{result.type}</span>
										{/if}
									</div>
									<div class="grid min-w-0 gap-px">
										<span class="truncate">{result.title}</span>
										<small class="truncate text-xs font-bold text-muted-foreground">
											{result.type}{result.year ? ` · ${result.year}` : ''}
										</small>
										{#if result.overview}
											<p class="m-0 truncate text-xs leading-snug text-muted-foreground">
												{result.overview}
											</p>
										{/if}
									</div>
								</Button>
							{/each}
						</div>
					{/if}
				{/each}
			{:else if loading}
				<div class="px-2 py-1 text-xs font-bold text-muted-foreground uppercase">Searching</div>
			{:else}
				<div class="px-2 py-1 text-xs font-bold text-muted-foreground uppercase">No matches</div>
			{/if}
			<a
				class="block border-t border-border px-2 py-2 text-sm font-bold text-primary no-underline hover:bg-accent hover:text-accent-foreground"
				href={trimmedQuery
					? resolve(`/search/advanced?q=${encodeURIComponent(trimmedQuery)}`)
					: resolve('/search/advanced')}
			>
				Advanced search
			</a>
		</div>
	{/if}
</div>
