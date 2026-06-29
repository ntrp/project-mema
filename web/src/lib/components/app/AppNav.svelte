<script lang="ts">
	import { resolve } from '$app/paths';
	import type { MediaSearchGroup, MediaSearchResult } from '$lib/settings/types';

	interface Props {
		searchQuery: string;
		groups: MediaSearchGroup[];
		loading: boolean;
		onSearch: (_query: string) => void | Promise<void>;
		onSelect: (_result: MediaSearchResult) => void;
		onProfile: () => void;
		onLogout: () => void | Promise<void>;
	}

	let {
		searchQuery = $bindable(),
		groups,
		loading,
		onSearch,
		onSelect,
		onProfile,
		onLogout
	}: Props = $props();
	let userMenuOpen = $state(false);
	let searchOpen = $state(false);

	const trimmedQuery = $derived(searchQuery.trim());
	const resultCount = $derived(groups.reduce((count, group) => count + group.results.length, 0));
	const showSuggestions = $derived(searchOpen && trimmedQuery.length >= 2);

	$effect(() => {
		const query = trimmedQuery;
		if (query.length < 2) {
			return;
		}
		const timeout = window.setTimeout(() => {
			void onSearch(query);
		}, 300);
		return () => window.clearTimeout(timeout);
	});

	function selectResult(result: MediaSearchResult) {
		searchQuery = result.title;
		searchOpen = false;
		onSelect(result);
	}

	function resultKey(result: MediaSearchResult) {
		return `${result.id ?? ''}:${result.type}:${result.externalProvider ?? ''}:${result.externalId ?? ''}:${result.title}:${result.year ?? ''}`;
	}
</script>

<header class="app-nav">
	<div class="global-search">
		<label for="global-search">Search</label>
		<input
			id="global-search"
			bind:value={searchQuery}
			placeholder="Search Movies & TV"
			autocomplete="off"
			onfocus={() => (searchOpen = true)}
			onblur={() => {
				window.setTimeout(() => {
					searchOpen = false;
				}, 120);
			}}
		/>
		{#if showSuggestions}
			<div class="search-suggestions" role="listbox" aria-label="Search suggestions">
				{#if loading}
					<div class="search-status">Searching</div>
				{:else if resultCount > 0}
					{#each groups as group (`${group.sourceType}:${group.sourceName}`)}
						{#if group.results.length > 0}
							<div class="search-group">
								<div class="search-group-title">{group.sourceName}</div>
								{#each group.results as result (resultKey(result))}
									<button
										type="button"
										role="option"
										aria-selected={searchQuery === result.title}
										onpointerdown={(event) => event.preventDefault()}
										onclick={() => selectResult(result)}
									>
										<span>{result.title}</span>
										<small>{result.type}{result.year ? ` · ${result.year}` : ''}</small>
									</button>
								{/each}
							</div>
						{/if}
					{/each}
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

	<nav class="nav-actions" aria-label="Application actions">
		<div class="user-menu">
			<button
				type="button"
				class="icon-button"
				aria-label="User menu"
				aria-haspopup="menu"
				aria-expanded={userMenuOpen}
				title="User"
				onclick={() => (userMenuOpen = !userMenuOpen)}
			>
				<span aria-hidden="true">A</span>
			</button>
			{#if userMenuOpen}
				<div class="user-dropdown" role="menu">
					<button
						type="button"
						role="menuitem"
						onclick={() => {
							userMenuOpen = false;
							onProfile();
						}}
					>
						Profile
					</button>
					<button
						type="button"
						role="menuitem"
						onclick={() => {
							userMenuOpen = false;
							void onLogout();
						}}
					>
						Logout
					</button>
				</div>
			{/if}
		</div>
	</nav>
</header>
