<script lang="ts">
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import MediaDetailSkeleton from './MediaDetailSkeleton.svelte';
	import MediaFilesTable from './MediaFilesTable.svelte';
	import MediaMetadataCore from './MediaMetadataCore.svelte';
	import MediaMetadataHero from './MediaMetadataHero.svelte';
	import MediaMetadataShell from './MediaMetadataShell.svelte';
	import MediaRelatedSections from './MediaRelatedSections.svelte';
	import MediaSeriesSeasons from './MediaSeriesSeasons.svelte';
	import ReleaseCandidatesSection from './ReleaseCandidatesSection.svelte';
	import { resolve } from '$app/paths';
	import { mediaMetadataDetail } from './mediaDetail';
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
		loading?: boolean;
		mediaItems?: MediaItem[];
		requestedItemId: string;
		releaseResults?: ReleaseSearchState;
		activities: DownloadActivity[];
		searchingItemId?: string;
		refreshingMetadataItemId?: string;
		grabbingKey?: string;
		addingKey?: string;
		deletingMediaItemId?: string;
		canManage: boolean;
		actionLabel: string;
		onFindReleases: (_item: MediaItem) => void;
		onAutoSearchMedia: (_item: MediaItem) => void;
		onRefreshMediaMetadata: (_item: MediaItem) => void;
		onDeleteMediaFile: (_item: MediaItem, _path: string) => void;
		onDeleteMedia: (_item: MediaItem) => void;
		onGrabRelease: (_item: MediaItem, _release: ReleaseCandidate) => void;
		onAddMedia: (_candidate: MediaSearchResult) => void;
	}

	let {
		mediaType,
		item,
		loading = false,
		mediaItems = [],
		requestedItemId,
		releaseResults,
		activities,
		searchingItemId,
		refreshingMetadataItemId,
		grabbingKey,
		addingKey,
		deletingMediaItemId,
		canManage,
		actionLabel,
		onFindReleases,
		onAutoSearchMedia,
		onRefreshMediaMetadata,
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

{#if loading}
	<MediaDetailSkeleton />
{:else if item && detail}
	<MediaMetadataShell labelledby="library-media-title">
		<MediaMetadataHero {detail} titleId="library-media-title">
			{#snippet actions()}
				{#if canManage}
					<Tooltip.Root>
						<Tooltip.Trigger>
							{#snippet child({ props })}
								<Button
									{...props}
									type="button"
									variant="outline"
									size="icon-sm"
									class="ml-auto"
									aria-label="Refresh metadata"
									disabled={refreshingMetadataItemId === item.id}
									onclick={() => onRefreshMediaMetadata(item)}
								>
									<RefreshCwIcon
										class={refreshingMetadataItemId === item.id ? 'animate-spin' : undefined}
										aria-hidden="true"
									/>
								</Button>
							{/snippet}
						</Tooltip.Trigger>
						<Tooltip.Content>Refresh metadata</Tooltip.Content>
					</Tooltip.Root>
					<Tooltip.Root>
						<Tooltip.Trigger>
							{#snippet child({ props })}
								<Button
									{...props}
									type="button"
									variant="destructive"
									size="icon-sm"
									aria-label="Delete media"
									disabled={deletingMediaItemId === item.id}
									onclick={() => onDeleteMedia(item)}
								>
									<TrashIcon aria-hidden="true" />
								</Button>
							{/snippet}
						</Tooltip.Trigger>
						<Tooltip.Content>Delete media</Tooltip.Content>
					</Tooltip.Root>
				{/if}
			{/snippet}
		</MediaMetadataHero>

		<div class="grid items-start gap-7">
			<main class="grid min-w-0 gap-6 [&>section]:grid [&>section]:min-w-0 [&>section]:gap-2.5">
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
	</MediaMetadataShell>
{:else}
	<EmptyState
		title="Media item not found"
		description={`No monitored ${mediaType} matches ${requestedItemId}.`}
	/>
{/if}
