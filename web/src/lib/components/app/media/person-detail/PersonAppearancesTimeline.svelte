<script lang="ts">
	import MediaPosterCard from '$lib/components/app/media/posters/MediaPosterCard.svelte';
	import type { MediaItem, MediaSearchResult, PersonAppearance } from '$lib/settings/types';
	import { appearanceTimelineData, type AppearanceTimelineItem } from './personTimeline';
	import { nextTimelineCardScrollLeft, nextTimelineYearScrollLeft } from './personTimelineScroll';
	import { dragScroll } from './personTimelineDrag';
	import PersonTimelineControls from './PersonTimelineControls.svelte';
	import { appearanceResult, resultKey } from './personDetail';

	interface Props {
		appearances: PersonAppearance[];
		mediaItems?: MediaItem[];
		addingKey?: string;
		actionLabel: string;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let { appearances, mediaItems = [], addingKey, actionLabel, onAdd }: Props = $props();

	let scroller = $state<HTMLDivElement>();
	const cardWidth = 151;
	const cardCircleOffset = Math.round(cardWidth * 1.65);
	const timelineY = 350;
	const emptyYearWidth = 96;
	const yearEntryWidth = 120;
	const paddingX = 80;
	const timeline = $derived(
		appearanceTimelineData(appearances, {
			cardWidth,
			cardGap: 32,
			cardTopY: timelineY - cardCircleOffset - 70,
			cardBottomY: timelineY + 70,
			emptyYearWidth,
			yearEntryWidth,
			paddingX
		})
	);
	const libraryKeys = $derived(new Set(mediaItems.map((item) => resultKey(item))));

	function scrollTimeline(direction: -1 | 1) {
		if (!scroller) return;
		const target = nextTimelineCardScrollLeft(
			timeline.items,
			scrollMetrics(scroller),
			cardWidth,
			direction
		);
		scroller.scrollTo({ left: target, behavior: 'smooth' });
	}

	function scrollTimelineEdge(direction: -1 | 1) {
		if (!scroller) return;
		scroller.scrollTo({
			left: direction < 0 ? 0 : scroller.scrollWidth - scroller.clientWidth,
			behavior: 'smooth'
		});
	}

	function wheelTimeline(event: WheelEvent) {
		if (!scroller) return;
		const delta = Math.abs(event.deltaY) >= Math.abs(event.deltaX) ? event.deltaY : event.deltaX;
		if (Math.abs(delta) < 1) return;
		event.preventDefault();
		const target = nextTimelineYearScrollLeft(
			timeline.years,
			scrollMetrics(scroller),
			delta > 0 ? 1 : -1
		);
		scroller.scrollTo({ left: target, behavior: 'smooth' });
	}

	function scrollMetrics(node: HTMLDivElement) {
		return {
			scrollLeft: node.scrollLeft,
			clientWidth: node.clientWidth,
			scrollWidth: node.scrollWidth
		};
	}

	function cardCircleY(item: AppearanceTimelineItem) {
		return item.top ? item.cardY + cardCircleOffset : item.cardY;
	}
	function connectorPath(item: AppearanceTimelineItem) {
		const x = item.cardX + cardWidth / 2;
		const y = cardCircleY(item);
		const midY = (y + timelineY) / 2;
		return `M ${x} ${y} C ${x} ${midY}, ${item.markerX} ${midY}, ${item.markerX} ${timelineY}`;
	}
</script>

<div
	use:dragScroll
	bind:this={scroller}
	class="timeline-scroll overflow-x-auto pb-3"
	onwheel={wheelTimeline}
>
	<div class="relative h-[640px] min-w-max py-4" style={`width: ${timeline.contentWidth}px`}>
		<div class="absolute right-0 left-0 h-px bg-border" style={`top: ${timelineY}px`}></div>

		<svg
			class="pointer-events-none absolute inset-0"
			width={timeline.contentWidth}
			height="500"
			aria-hidden="true"
		>
			{#each timeline.items.filter((item) => !item.unreleased) as item (`line:${item.appearance.type}:${item.appearance.externalProvider}:${item.appearance.externalId}`)}
				<path d={connectorPath(item)} fill="none" stroke="rgb(59 130 246)" stroke-width="2.5" />
			{/each}
		</svg>

		{#each timeline.years as year (year.year)}
			<div
				class="absolute z-20 -translate-x-1/2 -translate-y-1/2 rounded-[3px] border border-white bg-background px-1.5 text-m font-semibold text-muted-foreground shadow-xs"
				style={`left: ${year.x}px; top: ${timelineY}px`}
			>
				{year.year}
			</div>
		{/each}

		{#each timeline.items as item (`dot:${item.appearance.type}:${item.appearance.externalProvider}:${item.appearance.externalId}`)}
			{#if !item.unreleased}
				<span
					class="absolute size-2.5 -translate-x-1/2 -translate-y-1/2 rounded-full border-2 border-blue-400 bg-background shadow-xs"
					style={`left: ${item.markerX}px; top: ${timelineY}px`}
					aria-hidden="true"
				></span>
			{/if}
			<div
				class="absolute"
				style={`left: ${item.cardX}px; top: ${item.cardY}px; width: ${cardWidth}px`}
			>
				{@render TimelineCard(item)}
			</div>
		{/each}
	</div>
</div>

<PersonTimelineControls onStep={scrollTimeline} onEdge={scrollTimelineEdge} />

{#snippet TimelineCard(item: AppearanceTimelineItem)}
	{@const result = appearanceResult(item.appearance)}
	<div class="relative grid w-full shrink-0 gap-1.5">
		{#if item.top}
			{@render RoleLabel(item)}
		{/if}
		{#if !item.unreleased && item.top}
			<span
				class="absolute top-[calc(100%-0.375rem)] left-1/2 z-0 size-3 -translate-x-1/2 rounded-full border-2 border-background bg-purple-500 shadow-md"
				aria-hidden="true"
			></span>
		{:else if !item.unreleased}
			<span
				class="absolute bottom-[calc(100%-0.375rem)] left-1/2 z-0 size-3 -translate-x-1/2 rounded-full border-2 border-background bg-purple-500 shadow-md"
				aria-hidden="true"
			></span>
		{/if}
		<div class="relative z-10 rounded-md">
			<MediaPosterCard
				{result}
				adding={addingKey === resultKey(result)}
				{actionLabel}
				inLibrary={libraryKeys.has(resultKey(result))}
				unreleased={item.unreleased}
				{onAdd}
			/>
		</div>
		{#if !item.top}
			{@render RoleLabel(item)}
		{/if}
	</div>
{/snippet}

{#snippet RoleLabel(item: AppearanceTimelineItem)}
	<p class="m-0 min-h-4 truncate text-center text-xs font-medium text-muted-foreground">
		{#if item.appearance.role}
			as {item.appearance.role}
		{:else}
			&nbsp;
		{/if}
	</p>
{/snippet}

<style>
	.timeline-scroll {
		scrollbar-width: none;
	}

	.timeline-scroll::-webkit-scrollbar {
		display: none;
	}
</style>
