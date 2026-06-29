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
		mediaType === 'movie' ? 'No monitored movies yet' : 'No monitored series yet'
	);
</script>

<div class="page-heading">
	<p>{label}</p>
	<h1 id="home-title">{heading}</h1>
</div>

<div class="media-card-grid">
	{#each items as item (item.id)}
		<a
			class="media-library-card"
			href={mediaType === 'movie'
				? resolve('/movies/[id]', { id: item.id })
				: resolve('/series/[id]', { id: item.id })}
			aria-label={`Open ${item.title} details`}
		>
			<div class="poster-placeholder compact">{mediaType === 'movie' ? 'Movie' : 'Series'}</div>
			<div class="media-library-card-body">
				<strong>{item.title}</strong>
				<span>{item.year ? `${item.year} · ` : ''}{item.type}</span>
				{#if item.tags?.length}
					<div class="media-tags compact-tags" aria-label="Tags">
						{#each item.tags.slice(0, 3) as tag (tag)}
							<span>{tag}</span>
						{/each}
					</div>
				{/if}
				<small class:status-enabled={item.monitored} class:status-disabled={!item.monitored}>
					{item.monitored ? 'Monitored' : 'Paused'}
				</small>
			</div>
		</a>
	{:else}
		<div class="panel">
			<p class="empty">{emptyText}</p>
		</div>
	{/each}
</div>
