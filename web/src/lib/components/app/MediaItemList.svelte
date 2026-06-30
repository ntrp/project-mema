<script lang="ts">
	import { resolve } from '$app/paths';
	import type { MediaItem, MediaType } from '$lib/settings/types';

	interface Props {
		mediaType: MediaType;
		items: MediaItem[];
	}

	let { mediaType, items }: Props = $props();

	const heading = $derived(mediaType === 'movie' ? 'Added movies' : 'Added series');
	const label = $derived(mediaType === 'movie' ? 'Movies' : 'Series');
	const emptyText = $derived(
		mediaType === 'movie' ? 'No movies added yet.' : 'No series added yet.'
	);

	function posterUrl(path?: string) {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/w500${path}`;
	}

	function statusLabel(status: MediaItem['status']) {
		switch (status) {
			case 'downloaded':
				return 'Downloaded';
			case 'downloading':
				return 'Downloading';
			default:
				return 'Missing';
		}
	}

	function typeLabel(type: MediaType) {
		return type === 'movie' ? 'Movie' : 'Series';
	}
</script>

<div class="page-heading">
	<p>{label}</p>
	<h1 id="home-title">{heading}</h1>
</div>

{#if items.length === 0}
	<section class="empty-state media-library-empty">
		<p>{emptyText}</p>
		<a class="add-action-button" href={resolve('/discover')}>
			<span class="app-icon" aria-hidden="true">travel_explore</span>
			<span>Discover</span>
		</a>
	</section>
{:else}
	<div class="media-card-grid">
		{#each items as item (item.id)}
			<a
				class="media-library-card"
				href={mediaType === 'movie'
					? resolve('/movies/[id]', { id: item.id })
					: resolve('/series/[id]', { id: item.id })}
				aria-label={`Open ${item.title} details`}
				style:--library-card-bg={posterUrl(item.posterPath)
					? `url("${posterUrl(item.posterPath)}")`
					: undefined}
			>
				<div class="library-cover">
					{#if posterUrl(item.posterPath)}
						<img src={posterUrl(item.posterPath)} alt="" loading="lazy" />
					{:else}
						<div class="poster-placeholder compact">{typeLabel(mediaType)}</div>
					{/if}
				</div>
				<div class="poster-hover">
					<span class="poster-year">{item.year ?? 'Unknown'}</span>
					<h3>{item.title}</h3>
					<p>{item.overview ?? 'No overview available.'}</p>
				</div>
				<span
					class="media-card-status-line"
					class:downloaded={item.status === 'downloaded'}
					class:downloading={item.status === 'downloading'}
					class:missing={item.status === 'missing'}
					aria-label={statusLabel(item.status)}
				></span>
			</a>
		{/each}
	</div>
{/if}
