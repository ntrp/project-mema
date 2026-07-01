<script lang="ts">
	/* global KeyboardEvent */
	import { resolve } from '$app/paths';
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

<div class="global-search">
	<label for="global-search">Search</label>
	<input
		id="global-search"
		bind:value={searchQuery}
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
		<div class="search-suggestions" role="listbox" aria-label="Search suggestions">
			{#if resultCount > 0}
				{#each groups as group (`${group.sourceType}:${group.sourceName}`)}
					{#if group.results.length > 0}
						<div class="search-group">
							<div class="search-group-title">{group.sourceName}</div>
							{#each group.results as result (resultKey(result))}
								{@const index = resultIndex(result)}
								<button
									type="button"
									role="option"
									aria-selected={selectedKey === resultKey(result)}
									class:active-option={index === selectedIndex}
									onpointerdown={(event) => event.preventDefault()}
									onclick={() => selectResult(result)}
								>
									<div class="search-result-thumb">
										{#if posterUrl(result.posterPath)}
											<img src={posterUrl(result.posterPath)} alt="" loading="lazy" />
										{:else}
											<span>{result.type}</span>
										{/if}
									</div>
									<div class="search-result-copy">
										<span>{result.title}</span>
										<small>{result.type}{result.year ? ` · ${result.year}` : ''}</small>
										{#if result.overview}
											<p>{result.overview}</p>
										{/if}
									</div>
								</button>
							{/each}
						</div>
					{/if}
				{/each}
			{:else if loading}
				<div class="search-status">Searching</div>
			{:else}
				<div class="search-status">No matches</div>
			{/if}
			<a
				class="advanced-search-link"
				href={trimmedQuery
					? resolve(`/search/advanced?q=${encodeURIComponent(trimmedQuery)}`)
					: resolve('/search/advanced')}
			>
				Advanced search
			</a>
		</div>
	{/if}
</div>
