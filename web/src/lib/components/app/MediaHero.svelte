<script lang="ts">
	import type { MediaItem, MediaItemStatus, MediaType } from '$lib/settings/types';

	interface Props {
		mediaType: MediaType;
		item: MediaItem;
		qualityProfileLabel: string;
		canManage: boolean;
		searchingItemId?: string;
		deletingMediaItemId?: string;
		onFindReleases: (_item: MediaItem) => void;
		onDeleteMedia: (_item: MediaItem) => void;
	}

	let {
		mediaType,
		item,
		qualityProfileLabel,
		canManage,
		searchingItemId,
		deletingMediaItemId,
		onFindReleases,
		onDeleteMedia
	}: Props = $props();

	function posterUrl(path?: string, size = 'w780') {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/${size}${path}`;
	}

	function statusLabel(status: MediaItemStatus) {
		switch (status) {
			case 'downloaded':
				return 'Downloaded';
			case 'downloading':
				return 'Downloading';
			default:
				return 'Missing';
		}
	}
</script>

<section class="metadata-hero" aria-labelledby="home-title">
	<div class="metadata-poster">
		{#if posterUrl(item.posterPath, 'w500')}
			<img src={posterUrl(item.posterPath, 'w500')} alt="" />
		{:else}
			<div class="poster-placeholder">{mediaType === 'movie' ? 'Movie' : 'Series'}</div>
		{/if}
	</div>
	<div class="metadata-title-block">
		<h1 id="home-title">{item.title}</h1>
		<p>{mediaType === 'movie' ? 'Movie' : 'Series'}</p>
		<div class="metadata-info-bar" aria-label="Library media information">
			<span><strong>Year</strong>{item.year ?? 'Unknown'}</span>
			<span><strong>Type</strong>{item.type}</span>
			<span><strong>Status</strong>{statusLabel(item.status)}</span>
			<span><strong>Profile</strong>{qualityProfileLabel}</span>
		</div>
		{#if item.tags?.length}
			<div class="metadata-tags" aria-label="Tags">
				{#each item.tags as tag (tag)}
					<span><span class="app-icon" aria-hidden="true">sell</span>{tag}</span>
				{/each}
			</div>
		{/if}
		{#if canManage}
			<div class="metadata-actions">
				<button
					type="button"
					disabled={searchingItemId === item.id}
					onclick={() => onFindReleases(item)}
				>
					{searchingItemId === item.id ? 'Queued' : 'Find releases'}
				</button>
				<button
					type="button"
					class="danger"
					disabled={deletingMediaItemId === item.id}
					onclick={() => onDeleteMedia(item)}
				>
					{deletingMediaItemId === item.id ? 'Removing' : 'Remove'}
				</button>
			</div>
		{/if}
	</div>
</section>
