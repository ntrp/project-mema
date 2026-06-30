<script lang="ts">
	/* global HTMLDivElement */
	import { resolve } from '$app/paths';
	import MediaPosterCard from './MediaPosterCard.svelte';
	import type { MediaItem, MediaMetadataDetails, MediaSearchResult } from '$lib/settings/types';

	interface Props {
		detail: MediaMetadataDetails;
		mediaItems?: MediaItem[];
		addingKey?: string;
		actionLabel: string;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let { detail, mediaItems = [], addingKey, actionLabel, onAdd }: Props = $props();
	let relatedRows = $state<Record<string, HTMLDivElement | undefined>>({});

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
		relatedRows[title] = node;
		return {
			destroy() {
				relatedRows[title] = undefined;
			}
		};
	}

	function scrollRelated(title: string, direction: number) {
		const row = relatedRows[title];
		if (!row) {
			return;
		}
		row.scrollBy({
			left: direction * Math.max(row.clientWidth - 160, 240),
			behavior: 'smooth'
		});
	}
</script>

{#each visibleSections as section (section.title)}
	{#if section.results.length > 0}
		{@const titleId = sectionId(section.title)}
		{@const href = sectionHref(section.kind)}
		<section class="metadata-related-section" aria-labelledby={titleId}>
			<div class="section-heading">
				{#if href}
					<!-- eslint-disable svelte/no-navigation-without-resolve -->
					<a class="section-title-link" {href}>
						<h2 id={titleId}>{section.title}</h2>
						<span class="app-icon" aria-hidden="true">arrow_forward</span>
					</a>
					<!-- eslint-enable svelte/no-navigation-without-resolve -->
				{:else}
					<h2 id={titleId}>{section.title}</h2>
				{/if}
				<div class="section-heading-actions">
					<span>{section.results.length}</span>
					<div class="poster-row-controls" aria-label={`${section.title} carousel controls`}>
						<button
							type="button"
							aria-label={`Scroll ${section.title} left`}
							onclick={() => scrollRelated(section.title, -1)}
						>
							<span class="app-icon" aria-hidden="true">chevron_left</span>
						</button>
						<button
							type="button"
							aria-label={`Scroll ${section.title} right`}
							onclick={() => scrollRelated(section.title, 1)}
						>
							<span class="app-icon" aria-hidden="true">chevron_right</span>
						</button>
					</div>
				</div>
			</div>
			<div class="poster-row" use:trackRelatedRow={section.title}>
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
					<a class="poster-card poster-more-card" {href} aria-label={`Open ${section.title}`}>
						<span class="poster-frame poster-more-frame">
							<span class="app-icon" aria-hidden="true">arrow_forward</span>
							<span>View all</span>
						</span>
					</a>
					<!-- eslint-enable svelte/no-navigation-without-resolve -->
				{/if}
			</div>
		</section>
	{/if}
{/each}
