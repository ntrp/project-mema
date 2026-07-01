<script lang="ts">
	import ArrowRightIcon from '@lucide/svelte/icons/arrow-right';
	import { resolve } from '$app/paths';
	import SectionHeading from '$lib/components/shared/SectionHeading.svelte';
	import MediaPosterCard from '../media/MediaPosterCard.svelte';
	import PosterRowControls from './PosterRowControls.svelte';
	import { createPosterRowScroller } from './posterRowScroller.svelte';
	import type { MediaItem, MediaMetadataDetails, MediaSearchResult } from '$lib/settings/types';

	interface Props {
		detail: MediaMetadataDetails;
		mediaItems?: MediaItem[];
		addingKey?: string;
		actionLabel: string;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let { detail, mediaItems = [], addingKey, actionLabel, onAdd }: Props = $props();
	const relatedScroller = createPosterRowScroller();

	const sections = $derived([
		{ title: 'Recommendations', kind: 'recommendations', results: detail.recommendations ?? [] },
		{
			title: detail.type === 'movie' ? 'Similar Movies' : 'Similar Series',
			kind: 'similar',
			results: detail.similar ?? []
		}
	]);
	const libraryExternalKeys = $derived(
		new Set(
			mediaItems.map((item) => externalKey(item)).filter((key): key is string => Boolean(key))
		)
	);
	const libraryTitleKeys = $derived(new Set(mediaItems.map((item) => titleKey(item))));
	const visibleSections = $derived(
		sections.map((section) => ({ ...section, results: section.results.slice(0, 20) }))
	);

	function resultKey(result: MediaSearchResult) {
		return `${result.type}:${result.title}:${result.year ?? ''}`;
	}

	function isInLibrary(result: MediaSearchResult) {
		const key = externalKey(result);
		return Boolean(key && libraryExternalKeys.has(key)) || libraryTitleKeys.has(titleKey(result));
	}

	function externalKey(item: MediaItem | MediaSearchResult) {
		if (!item.externalProvider || !item.externalId) {
			return undefined;
		}
		return `${item.type}:${clean(item.externalProvider)}:${clean(item.externalId)}`;
	}

	function titleKey(item: MediaItem | MediaSearchResult) {
		return `${item.type}:${clean(item.title)}:${item.year ?? ''}`;
	}

	function clean(value: string) {
		return value.trim().toLowerCase();
	}

	function sectionId(title: string) {
		return `metadata-${title.toLowerCase().replaceAll(' ', '-')}`;
	}

	function sectionHref(kind: string) {
		if (!detail.externalProvider || !detail.externalId) {
			return undefined;
		}
		const params = {
			provider: detail.externalProvider,
			type: detail.type,
			externalId: detail.externalId
		};
		return kind === 'recommendations'
			? resolve('/media/[provider]/[type]/[externalId]/recommendations', params)
			: resolve('/media/[provider]/[type]/[externalId]/similar', params);
	}

	function trackRelatedRow(node: HTMLDivElement, title: string) {
		return relatedScroller.trackRow(node, title);
	}

	function scrollRelated(title: string, direction: -1 | 1) {
		relatedScroller.scrollRow(title, direction, 160, 240);
	}
</script>

{#each visibleSections as section (section.title)}
	{#if section.results.length > 0}
		{@const titleId = sectionId(section.title)}
		{@const href = sectionHref(section.kind)}
		{@const edges = relatedScroller.edgeState(section.title)}
		<section aria-labelledby={titleId}>
			<SectionHeading title={section.title} {titleId} {href}>
				{#if href}
					<ArrowRightIcon aria-hidden="true" />
				{/if}
				{#snippet actions()}
					<span>{section.results.length}</span>
					<PosterRowControls
						ariaLabel={`${section.title} carousel controls`}
						leftLabel={`Scroll ${section.title} left`}
						rightLabel={`Scroll ${section.title} right`}
						canScrollLeft={edges.canScrollLeft}
						canScrollRight={edges.canScrollRight}
						onScroll={(direction) => scrollRelated(section.title, direction)}
					/>
				{/snippet}
			</SectionHeading>
			<div
				class="-mx-3.5 mt-[-12px] grid auto-cols-[minmax(190px,220px)] grid-flow-col gap-5 overflow-x-auto overflow-y-hidden overscroll-x-contain snap-x snap-proximity scroll-px-3.5 px-3.5 pt-[18px] pb-5 [scrollbar-width:none] max-sm:mx-0 max-sm:auto-cols-[minmax(128px,150px)] max-sm:gap-3 max-sm:px-0 max-sm:pt-3.5 max-sm:pb-4 [&::-webkit-scrollbar]:hidden"
				use:trackRelatedRow={section.title}
			>
				{#each section.results as result (resultKey(result))}
					<MediaPosterCard
						{result}
						adding={addingKey === resultKey(result)}
						inLibrary={isInLibrary(result)}
						{actionLabel}
						{onAdd}
					/>
				{/each}
				{#if href}
					<!-- eslint-disable svelte/no-navigation-without-resolve -->
					<a
						class="min-w-0 snap-start text-foreground no-underline"
						{href}
						aria-label={`Open ${section.title}`}
					>
						<span
							class="relative grid aspect-[2/3] place-items-center gap-2 overflow-hidden rounded-md border border-dashed border-border bg-muted text-center text-sm font-black text-primary-hover"
						>
							<ArrowRightIcon aria-hidden="true" />
							<span>View all</span>
						</span>
					</a>
					<!-- eslint-enable svelte/no-navigation-without-resolve -->
				{/if}
			</div>
		</section>
	{/if}
{/each}
