<script lang="ts">
	import HomeArea from '$lib/components/app/home/HomeArea.svelte';
	import type { ActivitySection, HomeSection } from '$lib/settings/types';
	import { getAppShellContext } from '$lib/features/app/appShellContext';

	interface Props {
		section: HomeSection;
		activitySection?: ActivitySection;
		selectedMediaItemId?: string;
		selectedRequestId?: string;
	}

	let {
		section,
		activitySection = 'queue',
		selectedMediaItemId,
		selectedRequestId
	}: Props = $props();
	const app = getAppShellContext();
</script>

<HomeArea
	activeSection={section}
	{activitySection}
	{selectedMediaItemId}
	{selectedRequestId}
	mediaItems={app.mediaItems}
	mediaRequests={app.mediaRequests}
	discoverSections={app.discoverSections}
	discoverBlacklist={app.discoverBlacklist}
	libraryFolders={app.libraryFolders}
	languages={app.languages}
	qualityProfiles={app.mediaProfiles}
	activities={app.activities}
	releaseBlocklist={app.releaseBlocklist}
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
	assemblingMediaItemId={app.assemblingMediaItemId}
	reviewingComponentDecisionId={app.reviewingComponentDecisionId}
	cancellingActivityId={app.cancellingActivityId}
	deletingActivityId={app.deletingActivityId}
	deletingReleaseBlocklistId={app.deletingReleaseBlocklistId}
	clearingReleaseBlocklist={app.clearingReleaseBlocklist}
	loadingActivity={app.loadingActivity}
	canManage={app.isAdmin}
	onAddMedia={app.addMedia}
	onBlacklistMedia={app.blacklistDiscoverMedia}
	onRemoveBlacklistMedia={app.removeDiscoverBlacklistItem}
	onApproveMediaRequest={app.approveMediaRequest}
	onFindReleases={app.findReleases}
	onAutoSearchMedia={app.autoSearchMedia}
	onSearchMediaSubtitle={app.searchMediaSubtitle}
	onGrabMediaSubtitle={app.grabMediaSubtitle}
	onDeleteMediaSubtitle={app.deleteMediaSubtitle}
	onUpdateMediaSubtitle={app.updateMediaSubtitle}
	onRefreshMediaMetadata={app.refreshMediaMetadata}
	onSaveMediaItemOptions={app.saveMediaItemOptions}
	onDeleteMediaFile={app.deleteMediaFile}
	onDeleteMediaFileTrack={app.deleteMediaFileTrack}
	onAssembleMediaComponents={app.assembleMediaComponents}
	onReviewComponentCompatibility={app.reviewComponentCompatibility}
	onDeleteMedia={app.deleteMediaItem}
	onGrabRelease={app.grabRelease}
	onRefreshActivity={app.loadDownloadActivity}
	onRefreshReleaseBlocklist={app.loadReleaseBlocklist}
	onCancelActivity={app.cancelActivity}
	onDeleteActivity={app.deleteActivity}
	onDeleteReleaseBlocklistItem={app.deleteReleaseBlocklistItem}
	onClearReleaseBlocklist={app.clearReleaseBlocklist}
/>
