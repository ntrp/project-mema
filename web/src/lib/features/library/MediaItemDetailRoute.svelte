<script lang="ts">
	import MediaDetail from '$lib/components/app/media/detail/MediaDetail.svelte';
	import { getAppShellContext } from '$lib/features/app/appShellContext';
	import type { MediaItem, MediaType } from '$lib/settings/types';

	interface Props {
		mediaType: MediaType;
		id: string;
	}

	let { mediaType, id }: Props = $props();
	const app = getAppShellContext();
	const item = $derived(
		app.mediaItems.find((entry: MediaItem) => entry.id === id && entry.type === mediaType)
	);
</script>

<MediaDetail
	{mediaType}
	{item}
	loading={app.loadingMediaItems && !item}
	mediaItems={app.mediaItems}
	libraryFolders={app.libraryFolders}
	languages={app.languages}
	qualityProfiles={app.mediaProfiles}
	requestedItemId={id}
	activities={app.activities}
	searchingItemId={app.searchingItemId}
	scanningMediaItemId={app.scanningMediaItemId}
	refreshingMetadataItemId={app.refreshingMetadataItemId}
	savingMediaItemOptionsId={app.savingMediaItemOptionsId}
	grabbingKey={app.grabbingKey}
	addingKey={app.addingKey}
	deletingMediaItemId={app.deletingMediaItemId}
	assemblingMediaItemId={app.assemblingMediaItemId}
	reviewingComponentDecisionId={app.reviewingComponentDecisionId}
	pendingFulfillmentActionKeys={Object.keys(app.pendingFulfillmentActions ?? {})}
	canManage={app.isAdmin}
	actionLabel={app.isAdmin ? 'Add' : 'Request'}
	onAutoSearchMedia={app.autoSearchMedia}
	onRescanMediaFiles={app.rescanMediaFiles}
	onSearchMediaSubtitle={app.searchMediaSubtitle}
	onGrabMediaSubtitle={app.grabMediaSubtitle}
	onDeleteMediaSubtitle={app.deleteMediaSubtitle}
	onUpdateMediaSubtitle={app.updateMediaSubtitle}
	onRefreshMediaMetadata={app.refreshMediaMetadata}
	onSaveMediaItemOptions={app.saveMediaItemOptions}
	onDeleteMediaFile={app.deleteMediaFile}
	onDeleteMediaFileTrack={app.deleteMediaFileTrack}
	onFulfillmentAction={app.enqueueMediaFulfillment}
	onAssembleMediaComponents={app.assembleMediaComponents}
	onReviewComponentCompatibility={app.reviewComponentCompatibility}
	onDeleteMedia={app.deleteMediaItem}
	onGrabRelease={app.grabRelease}
	onAddMedia={app.addMedia}
/>
