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
	libraryFolders={app.libraryFolders}
	qualityProfiles={app.mediaProfiles}
	releaseResults={app.releaseResults}
	activities={app.activities}
	loadingDiscover={app.loadingDiscover}
	addingKey={app.addingKey}
	approvingRequestId={app.approvingRequestId}
	searchingItemId={app.searchingItemId}
	scanningMediaItemId={app.scanningMediaItemId}
	updatingMediaModeItemId={app.updatingMediaModeItemId}
	grabbingKey={app.grabbingKey}
	deletingMediaItemId={app.deletingMediaItemId}
	cancellingActivityId={app.cancellingActivityId}
	loadingActivity={app.loadingActivity}
	canManage={app.isAdmin}
	onAddMedia={app.addMedia}
	onApproveMediaRequest={app.approveMediaRequest}
	onFindReleases={app.findReleases}
	onRescanMediaFiles={app.rescanMediaFiles}
	onUpdateMediaMode={app.updateMediaMode}
	onDeleteMedia={app.deleteMediaItem}
	onGrabRelease={app.grabRelease}
	onRefreshActivity={app.loadDownloadActivity}
	onCancelActivity={app.cancelActivity}
/>
