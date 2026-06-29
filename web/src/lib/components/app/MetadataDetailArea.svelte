<script lang="ts">
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
			details.type,
			details.seasonCount ? `${details.seasonCount} seasons` : undefined,
			details.runtimeMinutes ? `${details.runtimeMinutes} min` : undefined,
			genres.slice(0, 3).join(', ')
		].filter(Boolean);
		return parts.join(' | ');
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
				<div class="metadata-actions">
					{#if detail.voteAverage}
						<span class="status-pill">{Math.round(detail.voteAverage * 10)}%</span>
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

				{#if genres.length > 0}
					<div class="metadata-tags" aria-label="Genres">
						{#each genres as genre (genre)}
							<span>{genre}</span>
						{/each}
					</div>
				{/if}

				{#if seasons.length > 0}
					<section aria-labelledby="metadata-seasons-title">
						<h2 id="metadata-seasons-title">Seasons</h2>
						<div class="metadata-season-list">
							{#each seasons as season (season.name)}
								<div class="metadata-season-row">
									<strong>{season.name}</strong>
									<span
										>{season.episodeCount
											? `${season.episodeCount} episodes`
											: 'Episodes unknown'}</span
									>
								</div>
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

			<aside class="metadata-sidebar" aria-label="Metadata facts">
				{#if detail.status}
					<div><strong>Status</strong><span>{detail.status}</span></div>
				{/if}
				{#if detail.releaseDate || detail.firstAirDate}
					<div>
						<strong>{detail.type === 'movie' ? 'Release Date' : 'First Air Date'}</strong>
						<span>{detail.releaseDate ?? detail.firstAirDate}</span>
					</div>
				{/if}
				{#if detail.originalLanguage}
					<div><strong>Original Language</strong><span>{detail.originalLanguage}</span></div>
				{/if}
				{#if detail.episodeCount}
					<div><strong>Episodes</strong><span>{detail.episodeCount}</span></div>
				{/if}
				<div><strong>Provider</strong><span>{detail.externalProvider}</span></div>
			</aside>
		</div>
	</section>
{/if}
