<script lang="ts">
	import MediaFilesTable from './MediaFilesTable.svelte';
	import MediaMetadataCore from './MediaMetadataCore.svelte';
	import MediaMetadataHero from './MediaMetadataHero.svelte';
	import MediaRelatedSections from './MediaRelatedSections.svelte';
	import MediaSeriesSeasons from './MediaSeriesSeasons.svelte';
	import ReleaseCandidatesSection from './ReleaseCandidatesSection.svelte';
	import { resolve } from '$app/paths';
	import { imageUrl, mediaMetadataDetail } from './mediaDetail';
	import type {
		DownloadActivity,
		MediaItem,
		MediaSearchResult,
		MediaType,
		ReleaseCandidate,
		ReleaseSearchState
	} from '$lib/settings/types';

	interface Props {
		mediaType: MediaType;
		item?: MediaItem;
		mediaItems?: MediaItem[];
		requestedItemId: string;
		releaseResults?: ReleaseSearchState;
		activities: DownloadActivity[];
		searchingItemId?: string;
		grabbingKey?: string;
		addingKey?: string;
		deletingMediaItemId?: string;
		canManage: boolean;
		actionLabel: string;
		onFindReleases: (_item: MediaItem) => void;
		onAutoSearchMedia: (_item: MediaItem) => void;
		onDeleteMediaFile: (_item: MediaItem, _path: string) => void;
		onDeleteMedia: (_item: MediaItem) => void;
		onGrabRelease: (_item: MediaItem, _release: ReleaseCandidate) => void;
		onAddMedia: (_candidate: MediaSearchResult) => void;
	}

	let {
		mediaType,
		item,
		mediaItems = [],
		requestedItemId,
		releaseResults,
		activities,
		searchingItemId,
		grabbingKey,
		addingKey,
		deletingMediaItemId,
		canManage,
		actionLabel,
		onFindReleases,
		onAutoSearchMedia,
		onDeleteMediaFile,
		onDeleteMedia,
		onGrabRelease,
		onAddMedia
	}: Props = $props();

	const detail = $derived(item ? mediaMetadataDetail(item) : undefined);
	const itemActivities = $derived(
		item ? activities.filter((activity) => activity.mediaItemId === item.id) : []
	);
	const castHref = $derived(
		item
			? resolve(item.type === 'movie' ? '/movies/[id]/cast' : '/series/[id]/cast', {
					id: item.id
				})
			: undefined
	);
	const crewHref = $derived(
		item
			? resolve(item.type === 'movie' ? '/movies/[id]/crew' : '/series/[id]/crew', {
					id: item.id
				})
			: undefined
	);
</script>

{#if item && detail}
	<section
		class="metadata-detail media-library-detail"
		aria-labelledby="library-media-title"
		style:--backdrop-url={imageUrl(detail.backdropPath, 'original')
			? `url("${imageUrl(detail.backdropPath, 'original')}")`
			: undefined}
	>
		<MediaMetadataHero {detail} titleId="library-media-title">
			{#snippet actions()}
				{#if canManage}
					<button
						type="button"
						class="danger icon-button metadata-delete-action"
						aria-label="Delete media"
						title="Delete media"
						disabled={deletingMediaItemId === item.id}
						onclick={() => onDeleteMedia(item)}
					>
						<span class="app-icon" aria-hidden="true">delete</span>
					</button>
				{/if}
			{/snippet}
		</MediaMetadataHero>

		<div class="metadata-body">
			<main class="metadata-main">
				<MediaMetadataCore {detail} {castHref} {crewHref}>
					{#snippet seasonsContent()}
						{#if item.type === 'series'}
							<MediaSeriesSeasons
								{item}
								{releaseResults}
								activities={itemActivities}
								{searchingItemId}
								{grabbingKey}
								{canManage}
								onAutoSearch={onAutoSearchMedia}
								onManualSearch={onFindReleases}
								onDeleteFile={onDeleteMediaFile}
								{onGrabRelease}
							/>
						{/if}
					{/snippet}
					{#snippet beforeCastContent()}
						{#if item.type === 'movie'}
							<MediaFilesTable
								{item}
								{releaseResults}
								activities={itemActivities}
								{searchingItemId}
								{grabbingKey}
								{canManage}
								onAutoSearch={onAutoSearchMedia}
								onManualSearch={onFindReleases}
								onDeleteFile={onDeleteMediaFile}
								{onGrabRelease}
							/>
						{/if}
					{/snippet}
				</MediaMetadataCore>
				<ReleaseCandidatesSection
					{item}
					{releaseResults}
					{grabbingKey}
					{canManage}
					{onGrabRelease}
				/>
				<MediaRelatedSections {detail} {mediaItems} {addingKey} {actionLabel} onAdd={onAddMedia} />
			</main>
		</div>
	</section>
{:else}
	<div class="detail-stack">
		<section class="panel">
			<div class="page-heading">
				<p>{mediaType === 'movie' ? 'Movies' : 'Series'}</p>
				<h1 id="home-title">Media item not found</h1>
			</div>
			<p class="empty">No monitored {mediaType} matches {requestedItemId}.</p>
		</section>
	</div>
{/if}
