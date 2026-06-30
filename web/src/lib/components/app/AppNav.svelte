<script lang="ts">
	import EventNotifications from './EventNotifications.svelte';
	import AppNavSearch from './AppNavSearch.svelte';
	import type { MediaSearchGroup, MediaSearchResult } from '$lib/settings/types';

	interface Props {
		searchQuery: string;
		groups: MediaSearchGroup[];
		loading: boolean;
		onSearch: (_query: string) => void | Promise<void>;
		onSelect: (_result: MediaSearchResult) => void;
		onAdvancedSearch: (_query: string) => void;
		onProfile: () => void;
		onLogout: () => void | Promise<void>;
		showNotifications?: boolean;
	}

	let {
		searchQuery = $bindable(),
		groups,
		loading,
		onSearch,
		onSelect,
		onAdvancedSearch,
		onProfile,
		onLogout,
		showNotifications = false
	}: Props = $props();
	let userMenuOpen = $state(false);
</script>

<header class="app-nav">
	<AppNavSearch bind:searchQuery {groups} {loading} {onSearch} {onSelect} {onAdvancedSearch} />

	<nav class="nav-actions" aria-label="Application actions">
		{#if showNotifications}
			<EventNotifications />
		{/if}
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
