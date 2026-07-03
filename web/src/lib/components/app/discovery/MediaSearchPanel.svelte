<script lang="ts">
	import ArrowRightIcon from '@lucide/svelte/icons/arrow-right';
	import { resolve } from '$app/paths';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import type {
		DiscoverBlacklistItem,
		MediaDiscoverSection,
		MediaItem,
		MediaSearchResult
	} from '$lib/settings/types';
	import MediaPosterCard from '$lib/components/app/media/posters/MediaPosterCard.svelte';
	import PosterRowControls from '$lib/components/app/media/posters/PosterRowControls.svelte';
	import { createPosterRowScroller } from '$lib/components/app/media/posters/posterRowScroller.svelte';
	import { Skeleton } from '$lib/components/ui/skeleton';

	interface Props {
		sections: MediaDiscoverSection[];
		mediaItems: MediaItem[];
		loading: boolean;
		addingKey?: string;
		blacklistingKey?: string;
		actionLabel: string;
		canManage: boolean;
		blacklist: DiscoverBlacklistItem[];
		onAdd: (_candidate: MediaSearchResult) => void;
		onBlacklist: (_candidate: MediaSearchResult) => void;
	}

	let {
		sections,
		mediaItems,
		loading,
		addingKey,
		blacklistingKey,
		actionLabel,
		canManage,
		blacklist,
		onAdd,
		onBlacklist
	}: Props = $props();
	const rowScroller = createPosterRowScroller();

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

	function resultKey(result: MediaSearchResult) {
		return `${result.type}:${result.externalProvider ?? ''}:${result.externalId ?? ''}:${result.title}:${result.year ?? ''}`;
	}

	function sectionResults(section: MediaDiscoverSection) {
		return (section.results ?? []).filter(
			(result) => !isInLibrary(result) && !isBlacklisted(result)
		);
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

	function trackRow(node: HTMLDivElement, sectionId: string) {
		return rowScroller.trackRow(node, sectionId);
	}

	function scrollSection(sectionId: string, direction: -1 | 1) {
		rowScroller.scrollRow(sectionId, direction);
	}
</script>

<PageHeading eyebrow="Discover" title="Browse media from metadata providers" titleId="home-title" />

{#if loading}
	<div class="min-w-0" aria-busy="true" aria-live="polite">
		<SectionHeading title="Loading discovery" />
		<div
			class="-mx-3.5 mt-[-12px] grid auto-cols-[minmax(190px,220px)] grid-flow-col gap-5 overflow-x-auto overflow-y-hidden overscroll-x-contain snap-x snap-proximity scroll-px-3.5 px-3.5 pt-[18px] pb-5 [scrollbar-width:none] max-sm:mx-0 max-sm:auto-cols-[minmax(128px,150px)] max-sm:gap-3 max-sm:px-0 max-sm:pt-3.5 max-sm:pb-4 [&::-webkit-scrollbar]:hidden"
		>
			{#each Array.from({ length: 8 }) as _, index (index)}
				<Skeleton class="min-w-0 snap-start aspect-[2/3]" aria-hidden="true" />
			{/each}
		</div>
	</div>
{:else if sections.length === 0}
	<EmptyState
		title="No discovery sections available"
		description="Enable and configure TMDB in metadata settings to load provider-backed discovery."
	/>
{:else}
	<div class="grid gap-[22px]" aria-label="Discover media sections">
		{#each sections as section (section.id)}
			{@const results = sectionResults(section)}
			{@const edges = rowScroller.edgeState(section.id)}
			<section class="min-w-0" aria-labelledby={`discover-${section.id}`}>
				<SectionHeading
					title={section.title}
					titleId={`discover-${section.id}`}
					href={resolve('/discover/[sectionId]', { sectionId: section.id })}
				>
					<ArrowRightIcon aria-hidden="true" />
					{#snippet actions()}
						<PosterRowControls
							ariaLabel={`${section.title} carousel controls`}
							canScrollLeft={edges.canScrollLeft}
							canScrollRight={edges.canScrollRight}
							onScroll={(direction) => scrollSection(section.id, direction)}
						/>
					{/snippet}
				</SectionHeading>
				{#if results.length > 0}
					<div
						class="-mx-3.5 mt-[-12px] grid auto-cols-[minmax(190px,220px)] grid-flow-col gap-5 overflow-x-auto overflow-y-hidden overscroll-x-contain snap-x snap-proximity scroll-px-3.5 px-3.5 pt-[18px] pb-5 [scrollbar-width:none] max-sm:mx-0 max-sm:auto-cols-[minmax(128px,150px)] max-sm:gap-3 max-sm:px-0 max-sm:pt-3.5 max-sm:pb-4 [&::-webkit-scrollbar]:hidden"
						use:trackRow={section.id}
					>
						{#each results as result (resultKey(result))}
							<MediaPosterCard
								{result}
								adding={addingKey === resultKey(result)}
								blacklisting={blacklistingKey === resultKey(result)}
								{actionLabel}
								showBlacklistAction={canManage}
								{onAdd}
								{onBlacklist}
							/>
						{/each}
						<a
							class="min-w-0 snap-start text-foreground no-underline"
							href={resolve('/discover/[sectionId]', { sectionId: section.id })}
							aria-label={`Open all ${section.title}`}
						>
							<div
								class="relative grid aspect-[2/3] place-items-center gap-2 overflow-hidden rounded-md border border-dashed border-border bg-muted text-center text-sm font-black text-primary-hover"
							>
								<ArrowRightIcon aria-hidden="true" />
								<span>View all</span>
							</div>
						</a>
					</div>
				{:else}
					<div
						class="m-0 rounded-md border border-dashed border-border bg-muted p-[18px] text-sm text-muted-foreground"
					>
						No results loaded for this section.
					</div>
				{/if}
			</section>
		{/each}
	</div>
{/if}
