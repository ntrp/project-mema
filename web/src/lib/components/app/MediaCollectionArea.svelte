<script lang="ts">
	import { resolve } from '$app/paths';
	import type { MediaCollection, MediaItem, MediaSearchResult } from '$lib/settings/types';

	interface Props {
		collection?: MediaCollection;
		mediaItems: MediaItem[];
		loading: boolean;
		addingKey?: string;
		actionLabel: string;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let { collection, mediaItems, loading, addingKey, actionLabel, onAdd }: Props = $props();

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

{#if loading}
	<section class="empty-state">
		<h2>Loading collection</h2>
		<p>Fetching collection media from the metadata provider.</p>
	</section>
{:else if !collection}
	<section class="empty-state">
		<h2>Collection not available</h2>
		<p>Could not load this collection.</p>
	</section>
{:else}
	<div class="page-heading">
		<p>{collection.provider}</p>
		<h1>{collection.name}</h1>
		{#if collection.overview}
			<p>{collection.overview}</p>
		{/if}
	</div>

	<section class="discover-section" aria-labelledby="collection-results-title">
		<div class="section-heading">
			<h2 id="collection-results-title">Collection media</h2>
			<span>{collection.results.length} titles</span>
		</div>
		<div class="poster-row">
			{#each collection.results as result (resultKey(result))}
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
						<span class="media-badge">{isInLibrary(result) ? 'In library' : result.type}</span>
						<div class="poster-hover">
							<span class="poster-year">{result.year ?? 'Unknown'}</span>
							<h3>{result.title}</h3>
							<p>{result.overview ?? 'No overview available.'}</p>
							{#if isInLibrary(result)}
								<span class="status-pill">In library</span>
							{:else}
								<button
									type="button"
									disabled={addingKey === resultKey(result)}
									onclick={() => onAdd(result)}
								>
									{addingKey === resultKey(result) ? 'Working' : actionLabel}
								</button>
							{/if}
						</div>
					</div>
				</article>
			{/each}
		</div>
	</section>
{/if}
