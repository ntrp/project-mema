<script lang="ts">
	import HomeArea from '$lib/components/app/HomeArea.svelte';
	import type { createAppShellController } from '$lib/components/app/appShellController.svelte';

	type AppShellController = ReturnType<typeof createAppShellController>;

	interface Props {
		app: AppShellController;
	}

	let { app }: Props = $props();
</script>

<HomeArea
	activeSection={app.activeHomeSection}
	selectedMediaItemId={app.selectedMediaItemId}
	selectedRequestId={app.selectedRequestId}
	mediaItems={app.mediaItems}
	mediaRequests={app.mediaRequests}
	discoverSections={app.discoverSections}
	discoverBlacklist={app.discoverBlacklist}
	libraryFolders={app.libraryFolders}
	qualityProfiles={app.mediaProfiles}
	releaseResults={app.releaseResults}
	activities={app.activities}
	loadingDiscover={app.loadingDiscover}
	loadingBlacklist={app.loadingBlacklist}
	addingKey={app.addingKey}
	blacklistingKey={app.blacklistingKey}
	removingBlacklistId={app.removingBlacklistId}
	approvingRequestId={app.approvingRequestId}
	searchingItemId={app.searchingItemId}
	scanningMediaItemId={app.scanningMediaItemId}
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
	onRescanMediaFiles={app.rescanMediaFiles}
	onDeleteMediaFile={app.deleteMediaFile}
	onDeleteMedia={app.deleteMediaItem}
	onGrabRelease={app.grabRelease}
	onRefreshActivity={app.loadDownloadActivity}
	onCancelActivity={app.cancelActivity}
	onDeleteActivity={app.deleteActivity}
/>
