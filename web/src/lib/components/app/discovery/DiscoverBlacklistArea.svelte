<script lang="ts">
	import type { DiscoverBlacklistItem } from '$lib/settings/types';

	interface Props {
		items: DiscoverBlacklistItem[];
		loading: boolean;
		removingId?: string;
		onRemove: (_item: DiscoverBlacklistItem) => void;
	}

	let { items, loading, removingId, onRemove }: Props = $props();

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
	<h1 id="home-title">Blacklist</h1>
	<p>{items.length} hidden titles</p>
</div>

{#if loading}
	<div class="media-card-grid">
		{#each Array.from({ length: 8 }) as _, index (index)}
			<div class="poster-card skeleton-card" aria-hidden="true"></div>
		{/each}
	</div>
{:else if items.length === 0}
	<section class="empty-state">
		<h2>No blacklisted media</h2>
		<p>Use the hidden-eye action on discover cards to hide titles from discovery.</p>
	</section>
{:else}
	<div class="media-card-grid discover-blacklist-grid">
		{#each items as item (item.id)}
			<article class="poster-card">
				<div class="poster-frame">
					{#if posterUrl(item.posterPath)}
						<img src={posterUrl(item.posterPath)} alt="" loading="lazy" />
					{:else}
						<div class="poster-placeholder">{item.type}</div>
					{/if}
					<span class="media-badge" class:movie={item.type === 'movie'}>{item.type}</span>
					<div class="poster-hover blacklist-card-hover">
						<span class="poster-year">{item.year ?? 'Unknown'}</span>
						<h3>{item.title}</h3>
						<p>{item.overview ?? 'No overview available.'}</p>
						<button
							type="button"
							class="poster-icon-action blacklist-action"
							disabled={removingId === item.id}
							aria-label={`Remove ${item.title} from blacklist`}
							title="Remove from blacklist"
							onclick={() => onRemove(item)}
						>
							<span class="app-icon" aria-hidden="true">visibility</span>
						</button>
					</div>
				</div>
			</article>
		{/each}
	</div>
{/if}
