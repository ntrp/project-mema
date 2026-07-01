<script lang="ts">
	/* global HTMLDivElement */
	import { resolve } from '$app/paths';
	import type {
		DiscoverBlacklistItem,
		MediaDiscoverSection,
		MediaItem,
		MediaSearchResult
	} from '$lib/settings/types';
	import MediaPosterCard from '../media/MediaPosterCard.svelte';

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
	let rows = $state<Record<string, HTMLDivElement>>({});

	const safeSections = $derived(sections ?? []);
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
		if (key && libraryExternalKeys.has(key)) {
			return true;
		}
		return libraryTitleKeys.has(titleKey(result));
	}

	function isBlacklisted(result: MediaSearchResult) {
		const key = externalKey(result);
		if (key && blacklistExternalKeys.has(key)) {
			return true;
		}
		return blacklistTitleKeys.has(titleKey(result));
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
		rows[sectionId] = node;
		return {
			destroy() {
				delete rows[sectionId];
			}
		};
	}

	function scrollSection(sectionId: string, direction: number) {
		const row = rows[sectionId];
		if (!row) {
			return;
		}
		row.scrollBy({
			left: direction * Math.max(row.clientWidth - 140, 220),
			behavior: 'smooth'
		});
	}
</script>

<div class="page-heading">
	<p>Discover</p>
	<h1 id="home-title">Browse media from metadata providers</h1>
</div>

{#if loading}
	<div class="discover-section skeleton-section">
		<div class="section-heading">
			<h2>Loading discovery</h2>
		</div>
		<div class="poster-row">
			{#each Array.from({ length: 8 }) as _, index (index)}
				<div class="poster-card skeleton-card" aria-hidden="true"></div>
			{/each}
		</div>
	</div>
{:else if safeSections.length === 0}
	<section class="empty-state">
		<h2>No discovery sections available</h2>
		<p>Enable and configure TMDB in metadata settings to load provider-backed discovery.</p>
	</section>
{:else}
	<div class="discover-sections" aria-label="Discover media sections">
		{#each safeSections as section (section.id)}
			{@const results = sectionResults(section)}
			<section class="discover-section" aria-labelledby={`discover-${section.id}`}>
				<div class="section-heading">
					<a
						class="section-title-link"
						href={resolve('/discover/[sectionId]', { sectionId: section.id })}
					>
						<h2 id={`discover-${section.id}`}>{section.title}</h2>
						<span class="app-icon" aria-hidden="true">arrow_forward</span>
					</a>
					<div class="poster-row-controls" aria-label={`${section.title} carousel controls`}>
						<button
							type="button"
							aria-label="Scroll left"
							onclick={() => scrollSection(section.id, -1)}
						>
							<span class="app-icon" aria-hidden="true">chevron_left</span>
						</button>
						<button
							type="button"
							aria-label="Scroll right"
							onclick={() => scrollSection(section.id, 1)}
						>
							<span class="app-icon" aria-hidden="true">chevron_right</span>
						</button>
					</div>
				</div>
				{#if results.length > 0}
					<div class="poster-row" use:trackRow={section.id}>
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
							class="poster-card poster-more-card"
							href={resolve('/discover/[sectionId]', { sectionId: section.id })}
							aria-label={`Open all ${section.title}`}
						>
							<div class="poster-frame poster-more-frame">
								<span class="app-icon" aria-hidden="true">arrow_forward</span>
								<span>View all</span>
							</div>
						</a>
					</div>
				{:else}
					<div class="section-empty">No results loaded for this section.</div>
				{/if}
			</section>
		{/each}
	</div>
{/if}
