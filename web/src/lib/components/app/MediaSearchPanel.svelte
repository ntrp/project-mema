<script lang="ts">
	import { resolve } from '$app/paths';
	import type { MediaDiscoverSection, MediaItem, MediaSearchResult } from '$lib/settings/types';

	interface Props {
		sections: MediaDiscoverSection[];
		mediaItems: MediaItem[];
		loading: boolean;
		addingKey?: string;
		actionLabel: string;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let { sections, mediaItems, loading, addingKey, actionLabel, onAdd }: Props = $props();

	const safeSections = $derived(sections ?? []);
	const libraryExternalKeys = $derived(
		new Set(
			(mediaItems ?? [])
				.map((item) => externalKey(item))
				.filter((key): key is string => Boolean(key))
		)
	);
	const libraryTitleKeys = $derived(new Set((mediaItems ?? []).map((item) => titleKey(item))));

	function resultKey(result: MediaSearchResult) {
		return `${result.type}:${result.externalProvider ?? ''}:${result.externalId ?? ''}:${result.title}:${result.year ?? ''}`;
	}

	function sectionResults(section: MediaDiscoverSection) {
		return (section.results ?? []).filter((result) => !isInLibrary(result));
	}

	function isInLibrary(result: MediaSearchResult) {
		const key = externalKey(result);
		if (key && libraryExternalKeys.has(key)) {
			return true;
		}
		return libraryTitleKeys.has(titleKey(result));
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

	function posterUrl(path?: string) {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/w342${path}`;
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
					<h2 id={`discover-${section.id}`}>{section.title}</h2>
					<span>{section.providerName}</span>
				</div>
				{#if results.length > 0}
					<div class="poster-row">
						{#each results as result (resultKey(result))}
							<article class="poster-card">
								<div class="poster-frame">
									{#if posterUrl(result.posterPath)}
										<img src={posterUrl(result.posterPath)} alt="" loading="lazy" />
									{:else}
										<div class="poster-placeholder">{result.type}</div>
									{/if}
									{#if result.externalProvider && result.externalId}
										<a
											class="poster-detail-link"
											href={resolve('/media/[provider]/[type]/[externalId]', {
												provider: result.externalProvider,
												type: result.type,
												externalId: result.externalId
											})}
											aria-label={`Open ${result.title} details`}
										></a>
									{/if}
									<span class="media-badge">{result.type}</span>
									<div class="poster-hover">
										<span class="poster-year">{result.year ?? 'Unknown'}</span>
										<h3>
											{result.title}
										</h3>
										<p>{result.overview ?? 'No overview available.'}</p>
										<button
											type="button"
											disabled={addingKey === resultKey(result)}
											onclick={(event) => {
												event.stopPropagation();
												onAdd(result);
											}}
										>
											{addingKey === resultKey(result) ? 'Working' : actionLabel}
										</button>
									</div>
								</div>
							</article>
						{/each}
					</div>
				{:else}
					<div class="section-empty">No results loaded for this section.</div>
				{/if}
			</section>
		{/each}
	</div>
{/if}
