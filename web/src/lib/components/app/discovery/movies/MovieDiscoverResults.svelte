<script lang="ts">
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import InlineSpinner from '$lib/components/shared/InlineSpinner.svelte';
	import MediaPosterCard from '$lib/components/app/media/posters/MediaPosterCard.svelte';
	import type { DiscoverBlacklistItem, MediaItem, MediaSearchResult } from '$lib/settings/types';

	interface Props {
		results: MediaSearchResult[];
		mediaItems: MediaItem[];
		blacklist: DiscoverBlacklistItem[];
		loading: boolean;
		loadingMore: boolean;
		hasSearched: boolean;
		addingKey?: string;
		blacklistingKey?: string;
		actionLabel: string;
		canManage: boolean;
		onAdd: (_candidate: MediaSearchResult) => void;
		onBlacklist: (_candidate: MediaSearchResult) => void;
	}

	let {
		results,
		mediaItems,
		blacklist,
		loading,
		loadingMore,
		hasSearched,
		addingKey,
		blacklistingKey,
		actionLabel,
		canManage,
		onAdd,
		onBlacklist
	}: Props = $props();

	const libraryKeys = $derived(new Set(mediaItems.map(externalKey).filter(Boolean)));
	const blacklistKeys = $derived(new Set(blacklist.map(externalKey).filter(Boolean)));

	function resultKey(result: MediaSearchResult) {
		return `${result.type}:${result.externalProvider ?? ''}:${result.externalId ?? ''}:${result.title}`;
	}

	function isInLibrary(result: MediaSearchResult) {
		return libraryKeys.has(externalKey(result));
	}

	function isBlacklisted(result: MediaSearchResult) {
		return blacklistKeys.has(externalKey(result));
	}

	function externalKey(item: MediaItem | MediaSearchResult | DiscoverBlacklistItem) {
		return `${item.type}:${item.externalProvider ?? ''}:${item.externalId ?? ''}`.toLowerCase();
	}
</script>

{#if loading && results.length === 0}
	<div
		class="grid grid-cols-[repeat(auto-fill,minmax(132px,1fr))] items-start gap-3 sm:grid-cols-[repeat(auto-fill,minmax(190px,1fr))] sm:gap-5"
	>
		{#each Array.from({ length: 12 }) as _, index (index)}
			<div class="min-w-0 aspect-[2/3] rounded-md bg-card" aria-hidden="true"></div>
		{/each}
	</div>
{:else if results.length > 0}
	<div
		class="grid grid-cols-[repeat(auto-fill,minmax(132px,1fr))] items-start gap-3 sm:grid-cols-[repeat(auto-fill,minmax(190px,1fr))] sm:gap-5"
	>
		{#each results as result (resultKey(result))}
			<MediaPosterCard
				{result}
				adding={addingKey === resultKey(result)}
				blacklisting={blacklistingKey === resultKey(result)}
				inLibrary={isInLibrary(result)}
				{actionLabel}
				showBlacklistAction={canManage && !isBlacklisted(result)}
				{onAdd}
				{onBlacklist}
			/>
		{/each}
	</div>
	{#if loadingMore}
		<div class="flex justify-center pt-6 pb-2">
			<InlineSpinner label="Loading more" />
		</div>
	{/if}
{:else if hasSearched}
	<EmptyState title="No movies found" description="Adjust the filters and search again." />
{:else}
	<EmptyState
		title="Search movies"
		description="Use filters to browse movies from metadata providers."
	/>
{/if}
