<script lang="ts">
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import MediaAddButton from '$lib/components/app/media/shared/MediaAddButton.svelte';
	import MediaMetadataCore from '$lib/components/app/media/metadata/MediaMetadataCore.svelte';
	import MediaDetailSkeleton from '$lib/components/app/media/detail/MediaDetailSkeleton.svelte';
	import MediaMetadataHero from '$lib/components/app/media/metadata/MediaMetadataHero.svelte';
	import MediaMetadataShell from '$lib/components/app/media/metadata/MediaMetadataShell.svelte';
	import MediaRelatedSections from '$lib/components/app/media/posters/MediaRelatedSections.svelte';
	import type { MediaItem, MediaMetadataDetails, MediaSearchResult } from '$lib/settings/types';

	interface Props {
		detail?: MediaMetadataDetails;
		loading: boolean;
		mediaItems?: MediaItem[];
		addingKey?: string;
		actionLabel: string;
		onAdd: (_candidate: MediaSearchResult) => void;
	}

	type MetadataAddCandidate = MediaSearchResult &
		Partial<MediaMetadataDetails> & { metadataStatus?: string };

	let { detail, loading, mediaItems = [], addingKey, actionLabel, onAdd }: Props = $props();

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
			keywords: details.keywords,
			facts: details.facts,
			seasons: details.seasons,
			cast: details.cast,
			recommendations: details.recommendations,
			similar: details.similar
		};
	}

	function candidateKey(details: MediaMetadataDetails) {
		return `${details.type}:${details.title}:${details.year ?? ''}`;
	}
</script>

{#if loading}
	<MediaDetailSkeleton />
{:else if !detail}
	<EmptyState
		title="Details not available"
		description="Could not load provider metadata for this item."
	/>
{:else}
	<MediaMetadataShell labelledby="metadata-detail-title">
		<MediaMetadataHero {detail} titleId="metadata-detail-title" showMonitorBookmark={false}>
			{#snippet actions()}
				<MediaAddButton
					result={candidate(detail)}
					adding={addingKey === candidateKey(detail)}
					label={actionLabel}
					size="sm"
					class="ml-auto"
					{onAdd}
				/>
			{/snippet}
		</MediaMetadataHero>

		<div class="grid items-start gap-7">
			<main class="grid min-w-0 gap-6 [&>section]:grid [&>section]:min-w-0 [&>section]:gap-2.5">
				<MediaMetadataCore {detail} />
				<MediaRelatedSections {detail} {mediaItems} {addingKey} {actionLabel} {onAdd} />
			</main>
		</div>
	</MediaMetadataShell>
{/if}
