<script lang="ts">
	import { providerDisplayName, providerPageUrl } from '$lib/settings/providerLinks';
	import type { MediaMetadataDetails, MediaSearchResult } from '$lib/settings/types';

	interface Props {
		detail?: MediaMetadataDetails;
		loading: boolean;
		addingKey?: string;
		actionLabel: string;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	let { detail, loading, addingKey, actionLabel, onAdd }: Props = $props();

	const genres = $derived(detail?.genres ?? []);
	const facts = $derived(detail?.facts ?? []);
	const seasons = $derived(detail?.seasons ?? []);
	const cast = $derived(detail?.cast ?? []);

	function imageUrl(path?: string, size = 'w780') {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/${size}${path}`;
	}

	function titleWithYear(details: MediaMetadataDetails) {
		return `${details.title}${details.year ? ` (${details.year})` : ''}`;
	}

	function subtitle(details: MediaMetadataDetails) {
		const parts = [
			details.type === 'movie' ? 'Movie' : 'Series',
			details.year,
			details.voteAverage ? `${Math.round(details.voteAverage * 10)}%` : undefined
		].filter(Boolean);
		return parts.join(' | ');
	}

	function topInfo(details: MediaMetadataDetails) {
		return [
			details.status ? ['Status', details.status] : undefined,
			details.releaseDate || details.firstAirDate
				? [
						details.type === 'movie' ? 'Release' : 'First aired',
						details.releaseDate ?? details.firstAirDate
					]
				: undefined,
			details.runtimeMinutes ? ['Runtime', `${details.runtimeMinutes} min`] : undefined,
			details.seasonCount ? ['Seasons', `${details.seasonCount}`] : undefined,
			details.episodeCount ? ['Episodes', `${details.episodeCount}`] : undefined,
			details.originalLanguage ? ['Language', details.originalLanguage.toUpperCase()] : undefined
		].filter((item): item is [string, string] => Boolean(item));
	}

	function candidate(details: MediaMetadataDetails): MediaSearchResult {
		return {
			title: details.title,
			type: details.type,
			year: details.year,
			externalProvider: details.externalProvider,
			externalId: details.externalId,
			overview: details.overview,
			posterPath: details.posterPath
		};
	}

	function candidateKey(details: MediaMetadataDetails) {
		return `${details.type}:${details.title}:${details.year ?? ''}`;
	}

	function externalUrl(details: MediaMetadataDetails) {
		return providerPageUrl(details.externalProvider, details.type, details.externalId);
	}

	function externalLabel(details: MediaMetadataDetails) {
		return providerDisplayName(details.externalProvider);
	}

	function episodeTitle(episode: { episodeNumber: number; name: string }) {
		return `${episode.episodeNumber} - ${episode.name}`;
	}
</script>

{#if loading}
	<section class="metadata-detail-loading panel">
		<p class="muted">Loading media details</p>
	</section>
{:else if !detail}
	<section class="empty-state">
		<h2>Details not available</h2>
		<p>Could not load provider metadata for this item.</p>
	</section>
{:else}
	<section
		class="metadata-detail"
		aria-labelledby="metadata-detail-title"
		style:--backdrop-url={imageUrl(detail.backdropPath, 'original')
			? `url("${imageUrl(detail.backdropPath, 'original')}")`
			: undefined}
	>
		<div class="metadata-hero">
			<div class="metadata-poster">
				{#if imageUrl(detail.posterPath, 'w342')}
					<img src={imageUrl(detail.posterPath, 'w342')} alt="" />
				{:else}
					<div class="poster-placeholder">{detail.type}</div>
				{/if}
			</div>
			<div class="metadata-title-block">
				<h1 id="metadata-detail-title">{titleWithYear(detail)}</h1>
				<p>{subtitle(detail)}</p>
				{#if topInfo(detail).length > 0}
					<div class="metadata-info-bar" aria-label="Media information">
						{#each topInfo(detail) as [label, value] (`${label}:${value}`)}
							<span><strong>{label}</strong>{value}</span>
						{/each}
					</div>
				{/if}
				{#if genres.length > 0}
					<div class="metadata-tags" aria-label="Genres">
						{#each genres as genre (genre)}
							<span><span class="app-icon" aria-hidden="true">sell</span>{genre}</span>
						{/each}
					</div>
				{/if}
				<div class="metadata-actions">
					{#if externalUrl(detail)}
						<!-- eslint-disable svelte/no-navigation-without-resolve -->
						<a
							class="external-link"
							href={externalUrl(detail)}
							target="_blank"
							rel="noreferrer"
							aria-label={`Open ${externalLabel(detail)} page in a new tab`}
						>
							<span class="app-icon" aria-hidden="true">open_in_new</span>
							<span>{externalLabel(detail)}</span>
						</a>
						<!-- eslint-enable svelte/no-navigation-without-resolve -->
					{/if}
					<button
						type="button"
						disabled={addingKey === candidateKey(detail)}
						onclick={() => onAdd(candidate(detail))}
					>
						{addingKey === candidateKey(detail) ? 'Working' : actionLabel}
					</button>
				</div>
			</div>
		</div>

		<div class="metadata-body">
			<main class="metadata-main">
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
											>{season.episodeCount
												? `${season.episodeCount} episodes`
												: 'Episodes unknown'}</small
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
			</main>
		</div>
	</section>
{/if}
