<script lang="ts">
	import MediaMetadataCore from './MediaMetadataCore.svelte';
	import MediaMetadataHero from './MediaMetadataHero.svelte';
	import type { MediaMetadataDetails, MediaSearchResult } from '$lib/settings/types';

	interface Props {
		detail?: MediaMetadataDetails;
		loading: boolean;
		addingKey?: string;
		actionLabel: string;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	type MetadataAddCandidate = MediaSearchResult &
		Partial<MediaMetadataDetails> & { metadataStatus?: string };

	let { detail, loading, addingKey, actionLabel, onAdd }: Props = $props();

	function imageUrl(path?: string, size = 'w780') {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/${size}${path}`;
	}

	function candidate(details: MediaMetadataDetails): MetadataAddCandidate {
		return {
			title: details.title,
			type: details.type,
			year: details.year,
			externalProvider: details.externalProvider,
			externalId: details.externalId,
			overview: details.overview,
			posterPath: details.posterPath,
			collectionId: details.collectionId,
			collectionName: details.collectionName,
			backdropPath: details.backdropPath,
			metadataStatus: details.status,
			originalLanguage: details.originalLanguage,
			releaseDate: details.releaseDate,
			firstAirDate: details.firstAirDate,
			runtimeMinutes: details.runtimeMinutes,
			seasonCount: details.seasonCount,
			episodeCount: details.episodeCount,
			voteAverage: details.voteAverage,
			genres: details.genres,
			facts: details.facts,
			seasons: details.seasons,
			cast: details.cast
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
		<MediaMetadataHero {detail} titleId="metadata-detail-title">
			{#snippet actions()}
				<button
					type="button"
					disabled={addingKey === candidateKey(detail)}
					onclick={() => onAdd(candidate(detail))}
				>
					{addingKey === candidateKey(detail) ? 'Working' : actionLabel}
				</button>
			{/snippet}
		</MediaMetadataHero>

		<div class="metadata-body">
			<main class="metadata-main">
				<MediaMetadataCore {detail} />
			</main>
		</div>
	</section>
{/if}
