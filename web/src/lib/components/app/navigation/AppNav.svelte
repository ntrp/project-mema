<script lang="ts">
	import EventNotifications from './EventNotifications.svelte';
	import AppNavSearch from './AppNavSearch.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import type { MediaSearchGroup, MediaSearchResult, UserSummary } from '$lib/settings/types';

	interface Props {
		searchQuery: string;
		currentUser?: UserSummary;
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
		currentUser,
		groups,
		loading,
		onSearch,
		onSelect,
		onAdvancedSearch,
		onProfile,
		onLogout,
		showNotifications = false
	}: Props = $props();

	const profileName = $derived(currentUser?.displayName || currentUser?.username || 'User');
	const profileInitial = $derived(profileName.slice(0, 1).toUpperCase());
</script>

<header
	class="sticky top-0 z-10 grid grid-cols-[auto_minmax(0,1fr)_auto] items-center gap-3 border-b border-border bg-card px-3 py-2 min-[641px]:gap-4 min-[641px]:px-[18px] min-[981px]:py-1.5"
>
	<Sidebar.Trigger class="shrink-0" />
	<AppNavSearch bind:searchQuery {groups} {loading} {onSearch} {onSelect} {onAdvancedSearch} />

	<nav class="flex w-full justify-end gap-2.5 justify-self-end" aria-label="Application actions">
		{#if showNotifications}
			<EventNotifications />
		{/if}
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Button variant="outline" size="icon" aria-label="User menu" {...props}>
						{#if currentUser?.pictureUrl}
							<img
								class="size-full rounded-md object-cover"
								src={currentUser.pictureUrl}
								alt=""
								aria-hidden="true"
							/>
						{:else}
							<span aria-hidden="true">{profileInitial}</span>
						{/if}
					</Button>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content align="end" class="w-44">
				<DropdownMenu.Item onclick={onProfile}>Profile</DropdownMenu.Item>
				<DropdownMenu.Item onclick={() => void onLogout()}>Logout</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</nav>
</header>
