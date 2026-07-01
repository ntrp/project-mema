<script lang="ts">
	import { onMount } from 'svelte';
	import type {
		DiscoverBlacklistItem,
		MediaDiscoverSection,
		MediaItem,
		MediaSearchResult
	} from '$lib/settings/types';
	import MediaPosterCard from '../media/MediaPosterCard.svelte';

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

	let showScrollTop = $state(false);

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
			showScrollTop = window.scrollY > 700;
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

	function scrollToTop() {
		window.scrollTo({ top: 0, behavior: 'smooth' });
	}
</script>

{#if loading}
	<div class="page-heading">
		<p>Discover</p>
		<h1>Loading section</h1>
	</div>
	<div class="media-card-grid">
		{#each Array.from({ length: 12 }) as _, index (index)}
			<div class="poster-card skeleton-card" aria-hidden="true"></div>
		{/each}
	</div>
{:else if !section}
	<section class="empty-state">
		<h2>Discovery section not available</h2>
		<p>Could not load this discover section.</p>
	</section>
{:else}
	<div class="page-heading">
		<p>Discover</p>
		<h1>{section.title}</h1>
		<p>{results.length} titles</p>
	</div>

	{#if results.length > 0}
		<div class="media-card-grid discover-section-grid">
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
			<div class="discover-section-loading">
				<span class="inline-spinner">Loading more</span>
			</div>
		{/if}
	{:else}
		<section class="empty-state">
			<h2>No results loaded</h2>
			<p>This section did not return any media.</p>
		</section>
	{/if}
{/if}

{#if showScrollTop}
	<button
		type="button"
		class="scroll-top-button"
		aria-label="Scroll to top"
		title="Scroll to top"
		onclick={scrollToTop}
	>
		<span class="app-icon" aria-hidden="true">keyboard_arrow_up</span>
	</button>
{/if}
