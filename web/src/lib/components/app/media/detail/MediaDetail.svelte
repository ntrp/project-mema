<script lang="ts">
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	import MediaDetailActions from '$lib/components/app/media/actions/MediaDetailActions.svelte';
	import MediaDetailSkeleton from '$lib/components/app/media/detail/MediaDetailSkeleton.svelte';
	import MediaFilesTable from '$lib/components/app/media/files/MediaFilesTable.svelte';
	import MediaComponentAssemblySection from '$lib/components/app/media/detail/assembly/MediaComponentAssemblySection.svelte';
	import MediaMetadataCore from '$lib/components/app/media/metadata/MediaMetadataCore.svelte';
	import MediaMetadataHero from '$lib/components/app/media/metadata/MediaMetadataHero.svelte';
	import MediaMetadataShell from '$lib/components/app/media/metadata/MediaMetadataShell.svelte';
	import MediaRelatedSections from '$lib/components/app/media/posters/MediaRelatedSections.svelte';
	import MediaSeriesSeasons from '$lib/components/app/media/series/MediaSeriesSeasons.svelte';
	import { resolve } from '$app/paths';
	import {
		mediaMetadataDetail,
		type MediaDetailProps as Props
	} from '$lib/components/app/media/detail/mediaDetail';
	import {
		monitorUpdate,
		titleMonitorActive,
		titleMonitorHint,
		titleMonitorStatus,
		toggledMediaMonitor
	} from '$lib/components/app/media/series/mediaMonitoring';
	import type { MediaItem } from '$lib/settings/types';

	let {
		mediaType,
		item,
		loading = false,
		mediaItems = [],
		libraryFolders,
		languages,
		qualityProfiles,
		requestedItemId,
		activities,
		searchingItemId,
		scanningMediaItemId,
		refreshingMetadataItemId,
		savingMediaItemOptionsId,
		assemblingMediaItemId,
		reviewingComponentDecisionId,
		grabbingKey,
		addingKey,
		deletingMediaItemId,
		canManage,
		actionLabel,
		onAutoSearchMedia,
		onSearchMediaSubtitle = () => {},
		onGrabMediaSubtitle = () => {},
		onDeleteMediaSubtitle = () => {},
		onUpdateMediaSubtitle = () => {},
		onRefreshMediaMetadata,
		onRescanMediaFiles,
		onSaveMediaItemOptions,
		onDeleteMediaFile,
		onDeleteMediaFileTrack = () => {},
		onAssembleMediaComponents = () => {},
		onReviewComponentCompatibility = () => {},
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

	function toggleMediaMonitor(nextItem: MediaItem) {
		onSaveMediaItemOptions(nextItem, monitorUpdate(toggledMediaMonitor(nextItem)));
	}
</script>

{#if loading}
	<MediaDetailSkeleton />
{:else if item && detail}
	<MediaMetadataShell labelledby="library-media-title">
		<MediaMetadataHero
			{detail}
			titleId="library-media-title"
			mediaStatus={item.status}
			monitorMonitored={titleMonitorActive(item)}
			monitorStatusText={titleMonitorStatus(item)}
			monitorHintText={titleMonitorHint(item)}
			monitorDisabled={!canManage}
			onToggleMonitor={() => toggleMediaMonitor(item)}
		>
			{#snippet actions()}
				{#if canManage}
					<MediaDetailActions
						{item}
						{qualityProfiles}
						refreshing={refreshingMetadataItemId === item.id}
						savingOptions={savingMediaItemOptionsId === item.id}
						deleting={deletingMediaItemId === item.id}
						onRefreshMetadata={onRefreshMediaMetadata}
						onSaveOptions={onSaveMediaItemOptions}
						onDelete={onDeleteMedia}
					/>
				{/if}
			{/snippet}
		</MediaMetadataHero>

		<div class="grid items-start gap-7">
			<main class="grid min-w-0 gap-6 [&>section]:grid [&>section]:min-w-0 [&>section]:gap-2.5">
				<MediaMetadataCore {detail} {castHref} {crewHref}>
					{#snippet seasonsContent()}
						{#if item.type === 'serie'}
							<MediaSeriesSeasons
								{item}
								activities={itemActivities}
								{libraryFolders}
								{languages}
								{qualityProfiles}
								{searchingItemId}
								{grabbingKey}
								{canManage}
								onSaveOptions={onSaveMediaItemOptions}
								onAutoSearch={onAutoSearchMedia}
								onSearchSubtitle={onSearchMediaSubtitle}
								onGrabSubtitle={onGrabMediaSubtitle}
								onDeleteSubtitle={onDeleteMediaSubtitle}
								onUpdateSubtitle={onUpdateMediaSubtitle}
								onDeleteFile={onDeleteMediaFile}
								onDeleteFileTrack={onDeleteMediaFileTrack}
								{onGrabRelease}
							/>
						{/if}
					{/snippet}
					{#snippet beforeCastContent()}
						{#if item.type === 'movie'}
							<MediaFilesTable
								{item}
								activities={itemActivities}
								{libraryFolders}
								{languages}
								{qualityProfiles}
								{searchingItemId}
								{scanningMediaItemId}
								{grabbingKey}
								{canManage}
								onSaveOptions={onSaveMediaItemOptions}
								onAutoSearch={onAutoSearchMedia}
								{onRescanMediaFiles}
								onSearchSubtitle={onSearchMediaSubtitle}
								onGrabSubtitle={onGrabMediaSubtitle}
								onDeleteSubtitle={onDeleteMediaSubtitle}
								onUpdateSubtitle={onUpdateMediaSubtitle}
								onDeleteFile={onDeleteMediaFile}
								onDeleteFileTrack={onDeleteMediaFileTrack}
								{onGrabRelease}
							/>
						{/if}
					{/snippet}
				</MediaMetadataCore>
				<MediaComponentAssemblySection
					{item}
					{canManage}
					{assemblingMediaItemId}
					{reviewingComponentDecisionId}
					onAssemble={onAssembleMediaComponents}
					onReview={onReviewComponentCompatibility}
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
