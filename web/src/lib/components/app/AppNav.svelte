<script lang="ts">
	import type { AppView } from '$lib/settings/types';

	interface Props {
		activeView: AppView;
		searchQuery: string;
		onHome: () => void;
		onSettings: () => void;
		onProfile: () => void;
		onLogout: () => void | Promise<void>;
	}

	let {
		activeView,
		searchQuery = $bindable(),
		onHome,
		onSettings,
		onProfile,
		onLogout
	}: Props = $props();
	let userMenuOpen = $state(false);
	let searchOpen = $state(false);

	const suggestions = [
		"Frieren: Beyond Journey's End",
		'Dune: Part Two',
		'The Apothecary Diaries',
		'The Last of Us',
		'Cyberpunk: Edgerunners',
		'Foundation',
		'Chainsaw Man'
	];
	const filteredSuggestions = $derived(
		searchQuery.trim()
			? suggestions
					.filter((item) => item.toLowerCase().includes(searchQuery.trim().toLowerCase()))
					.slice(0, 5)
			: suggestions.slice(0, 4)
	);
	const showSuggestions = $derived(searchOpen && filteredSuggestions.length > 0);

	function selectSuggestion(value: string) {
		searchQuery = value;
		searchOpen = false;
	}
</script>

<header class="app-nav">
	<button type="button" class="brand-button" aria-label="Open dashboard" onclick={onHome}>
		<span class="brand-mark" aria-hidden="true">M</span>
		<span class="brand-name">mema</span>
	</button>

	<div class="global-search">
		<label for="global-search">Search</label>
		<input
			id="global-search"
			bind:value={searchQuery}
			placeholder="Search movies, series, anime, music"
			autocomplete="off"
			onfocus={() => (searchOpen = true)}
			onblur={() => {
				window.setTimeout(() => {
					searchOpen = false;
				}, 100);
			}}
		/>
		{#if showSuggestions}
			<div class="search-suggestions" role="listbox" aria-label="Search suggestions">
				{#each filteredSuggestions as suggestion (suggestion)}
					<button
						type="button"
						role="option"
						aria-selected={searchQuery === suggestion}
						onpointerdown={(event) => event.preventDefault()}
						onclick={() => selectSuggestion(suggestion)}
					>
						{suggestion}
					</button>
				{/each}
			</div>
		{/if}
	</div>

	<nav class="nav-actions" aria-label="Application actions">
		<button
			type="button"
			class:active-icon={activeView === 'settings'}
			class="icon-button"
			aria-label="Settings"
			title="Settings"
			onclick={onSettings}
		>
			<span aria-hidden="true">⚙</span>
		</button>
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
				<span aria-hidden="true">👤</span>
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
