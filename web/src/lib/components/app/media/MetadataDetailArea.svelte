<script lang="ts">
	import PlusIcon from '@lucide/svelte/icons/plus';
	import { Button } from '$lib/components/ui/button';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import MediaMetadataCore from './MediaMetadataCore.svelte';
	import MediaMetadataHero from './MediaMetadataHero.svelte';
	import MediaMetadataShell from './MediaMetadataShell.svelte';
	import MediaRelatedSections from './MediaRelatedSections.svelte';
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
	<section class="min-h-[220px] rounded-md border border-border bg-card p-5">
		<p class="m-0 text-sm leading-6 text-muted-foreground">Loading media details</p>
	</section>
{:else if !detail}
	<EmptyState
		title="Details not available"
		description="Could not load provider metadata for this item."
	/>
{:else}
	<MediaMetadataShell backdropPath={detail.backdropPath} labelledby="metadata-detail-title">
		<MediaMetadataHero {detail} titleId="metadata-detail-title" showMonitorBookmark={false}>
			{#snippet actions()}
				<Button
					type="button"
					aria-label={addingKey === candidateKey(detail) ? 'Working' : actionLabel}
					disabled={addingKey === candidateKey(detail)}
					onclick={() => onAdd(candidate(detail))}
				>
					<PlusIcon aria-hidden="true" />
					<span>{addingKey === candidateKey(detail) ? 'Working' : actionLabel}</span>
				</Button>
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
