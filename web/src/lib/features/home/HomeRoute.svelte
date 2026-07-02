<script lang="ts">
	import HomeArea from '$lib/components/app/home/HomeArea.svelte';
	import type { HomeSection } from '$lib/settings/types';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	interface Props {
		section: HomeSection;
		selectedMediaItemId?: string;
		selectedRequestId?: string;
	}

	let { section, selectedMediaItemId, selectedRequestId }: Props = $props();
	const app = getAppShellContext();
</script>

<HomeArea
	activeSection={section}
	{selectedMediaItemId}
	{selectedRequestId}
	mediaItems={app.mediaItems}
	mediaRequests={app.mediaRequests}
	discoverSections={app.discoverSections}
	discoverBlacklist={app.discoverBlacklist}
	libraryFolders={app.libraryFolders}
	languages={app.languages}
	qualityProfiles={app.mediaProfiles}
	releaseResults={app.releaseResults}
	activities={app.activities}
	loadingDiscover={app.loadingDiscover}
	loadingBlacklist={app.loadingBlacklist}
	loadingMediaItems={app.loadingMediaItems}
	addingKey={app.addingKey}
	blacklistingKey={app.blacklistingKey}
	removingBlacklistId={app.removingBlacklistId}
	approvingRequestId={app.approvingRequestId}
	searchingItemId={app.searchingItemId}
	refreshingMetadataItemId={app.refreshingMetadataItemId}
	savingMediaItemOptionsId={app.savingMediaItemOptionsId}
	grabbingKey={app.grabbingKey}
	deletingMediaItemId={app.deletingMediaItemId}
	cancellingActivityId={app.cancellingActivityId}
	deletingActivityId={app.deletingActivityId}
	loadingActivity={app.loadingActivity}
	canManage={app.isAdmin}
	onAddMedia={app.addMedia}
	onBlacklistMedia={app.blacklistDiscoverMedia}
	onRemoveBlacklistMedia={app.removeDiscoverBlacklistItem}
	onApproveMediaRequest={app.approveMediaRequest}
	onFindReleases={app.findReleases}
	onAutoSearchMedia={app.autoSearchMedia}
	onRefreshMediaMetadata={app.refreshMediaMetadata}
	onSaveMediaItemOptions={app.saveMediaItemOptions}
	onDeleteMediaFile={app.deleteMediaFile}
	onDeleteMedia={app.deleteMediaItem}
	onGrabRelease={app.grabRelease}
	onRefreshActivity={app.loadDownloadActivity}
	onCancelActivity={app.cancelActivity}
	onDeleteActivity={app.deleteActivity}
/>
