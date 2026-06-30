<script lang="ts">
	import type { MediaMetadataDetails } from '$lib/settings/types';

	interface Props {
		detail: MediaMetadataDetails;
	}

	let { detail }: Props = $props();

	const facts = $derived(detail.facts ?? []);
	const seasons = $derived(detail.seasons ?? []);
	const cast = $derived(detail.cast ?? []);

	function imageUrl(path?: string, size = 'w780') {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/${size}${path}`;
	}

	function episodeTitle(episode: { episodeNumber: number; name: string }) {
		return `${episode.episodeNumber} - ${episode.name}`;
	}
</script>

<section aria-labelledby="metadata-overview-title">
	<h2 id="metadata-overview-title">Overview</h2>
	<p>{detail.overview ?? 'No overview available.'}</p>
</section>

{#if facts.length > 0}
	<section class="metadata-facts-grid" aria-label="Crew and source facts">
		{#each facts as fact (`${fact.label}:${fact.value}`)}
			<div>
				<strong>{fact.label}</strong>
				<span>{fact.value}</span>
			</div>
		{/each}
	</section>
{/if}

{#if seasons.length > 0}
	<section aria-labelledby="metadata-seasons-title">
		<h2 id="metadata-seasons-title">Seasons</h2>
		<div class="metadata-season-list">
			{#each seasons as season, index (season.name)}
				<details class="metadata-season-panel" open={index === 0}>
					<summary>
						<span>{season.name}</span>
						<small
							>{season.episodeCount ? `${season.episodeCount} episodes` : 'Episodes unknown'}</small
						>
					</summary>
					{#if season.episodes && season.episodes.length > 0}
						<div class="metadata-episode-list">
							{#each season.episodes as episode (episode.episodeNumber)}
								<article class="metadata-episode-row">
									<div class="metadata-episode-copy">
										<h3>
											{episodeTitle(episode)}
											{#if episode.airDate}
												<span>{episode.airDate}</span>
											{/if}
										</h3>
										<p>{episode.overview ?? 'No episode overview available.'}</p>
									</div>
									{#if imageUrl(episode.stillPath, 'w300')}
										<img src={imageUrl(episode.stillPath, 'w300')} alt="" loading="lazy" />
									{/if}
								</article>
							{/each}
						</div>
					{:else}
						<p class="metadata-season-empty">No episode details available.</p>
					{/if}
				</details>
			{/each}
		</div>
	</section>
{/if}

{#if cast.length > 0}
	<section aria-labelledby="metadata-cast-title">
		<h2 id="metadata-cast-title">Cast</h2>
		<div class="metadata-cast-grid">
			{#each cast as person (`${person.name}:${person.role ?? ''}`)}
				<article class="metadata-cast-card">
					<div>
						{#if imageUrl(person.profilePath, 'w185')}
							<img src={imageUrl(person.profilePath, 'w185')} alt="" loading="lazy" />
						{:else}
							<span>{person.name.slice(0, 1)}</span>
						{/if}
					</div>
					<strong>{person.name}</strong>
					{#if person.role}
						<p>{person.role}</p>
					{/if}
				</article>
			{/each}
		</div>
	</section>
{/if}
