<script lang="ts">
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import InlineSpinner from '$lib/components/shared/InlineSpinner.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import { onMount } from 'svelte';
	import type {
		DiscoverBlacklistItem,
		MediaDiscoverSection,
		MediaItem,
		MediaSearchResult
	} from '$lib/settings/types';
	import MediaPosterCard from '$lib/components/app/media/posters/MediaPosterCard.svelte';
	import { Skeleton } from '$lib/components/ui/skeleton';

	interface Props {
		section?: MediaDiscoverSection;
		mediaItems: MediaItem[];
		loading: boolean;
		loadingMore?: boolean;
		hasMore?: boolean;
		addingKey?: string;
		blacklistingKey?: string;
		actionLabel: string;
		canManage: boolean;
		blacklist: DiscoverBlacklistItem[];
		onAdd: (_candidate: MediaSearchResult) => void;
		onBlacklist: (_candidate: MediaSearchResult) => void;
		onLoadMore: () => void;
	}

	let {
		section,
		mediaItems,
		loading,
		loadingMore = false,
		hasMore = true,
		addingKey,
		blacklistingKey,
		actionLabel,
		canManage,
		blacklist,
		onAdd,
		onBlacklist,
		onLoadMore
	}: Props = $props();

	const libraryExternalKeys = $derived(
		new Set(
			(mediaItems ?? [])
				.map((item) => externalKey(item))
				.filter((key): key is string => Boolean(key))
		)
	);
	const libraryTitleKeys = $derived(new Set((mediaItems ?? []).map((item) => titleKey(item))));
	const blacklistExternalKeys = $derived(
		new Set(
			(blacklist ?? [])
				.map((item) => externalKey(item))
				.filter((key): key is string => Boolean(key))
		)
	);
	const blacklistTitleKeys = $derived(new Set((blacklist ?? []).map((item) => titleKey(item))));
	const results = $derived((section?.results ?? []).filter((result) => !isBlacklisted(result)));

	onMount(() => {
		const handleScroll = () => {
			const distanceToEnd =
				document.documentElement.scrollHeight - window.scrollY - window.innerHeight;
			if (!loading && !loadingMore && hasMore && distanceToEnd < 700) {
				onLoadMore();
			}
		};
		window.addEventListener('scroll', handleScroll, { passive: true });
		handleScroll();
		return () => window.removeEventListener('scroll', handleScroll);
	});

	function resultKey(result: MediaSearchResult) {
		return `${result.type}:${result.externalProvider ?? ''}:${result.externalId ?? ''}:${result.title}:${result.year ?? ''}`;
	}

	function isInLibrary(result: MediaSearchResult) {
		const key = externalKey(result);
		return Boolean(key && libraryExternalKeys.has(key)) || libraryTitleKeys.has(titleKey(result));
	}

	function isBlacklisted(result: MediaSearchResult) {
		const key = externalKey(result);
		return (
			Boolean(key && blacklistExternalKeys.has(key)) || blacklistTitleKeys.has(titleKey(result))
		);
	}

	function externalKey(item: MediaItem | MediaSearchResult | DiscoverBlacklistItem) {
		if (!item.externalProvider || !item.externalId) {
			return undefined;
		}
		return `${item.type}:${clean(item.externalProvider)}:${clean(item.externalId)}`;
	}

	function titleKey(item: MediaItem | MediaSearchResult | DiscoverBlacklistItem) {
		return `${item.type}:${clean(item.title)}:${item.year ?? ''}`;
	}

	function clean(value: string) {
		return value.trim().toLowerCase();
	}
</script>

{#if loading}
	<PageHeading eyebrow="Discover" title="Loading section" />
	<div
		aria-busy="true"
		aria-live="polite"
		class="grid grid-cols-[repeat(auto-fill,minmax(132px,1fr))] gap-3 sm:grid-cols-[repeat(auto-fill,minmax(190px,1fr))] sm:gap-5"
	>
		{#each Array.from({ length: 12 }) as _, index (index)}
			<Skeleton class="min-w-0 snap-start aspect-2/3" aria-hidden="true" />
		{/each}
	</div>
{:else if !section}
	<EmptyState
		title="Discovery section not available"
		description="Could not load this discover section."
	/>
{:else}
	<PageHeading eyebrow="Discover" title={section.title} />

	{#if results.length > 0}
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
					showBlacklistAction={canManage}
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
	{:else}
		<EmptyState title="No results loaded" description="This section did not return any media." />
	{/if}
{/if}
