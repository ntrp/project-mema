<script lang="ts">
	import ActivityList from '../activity/ActivityList.svelte';
	import DiscoverBlacklistArea from '../discovery/DiscoverBlacklistArea.svelte';
	import MediaDetail from '$lib/components/app/media/detail/MediaDetail.svelte';
	import MediaItemList from './MediaItemList.svelte';
	import MediaRequestArea from '../requests/MediaRequestArea.svelte';
	import MediaSearchPanel from '../discovery/MediaSearchPanel.svelte';
	import WantedMediaTable from './WantedMediaTable.svelte';
	import type { HomeAreaProps } from './homeAreaTypes';

	let {
		activeSection,
		activitySection,
		selectedMediaItemId,
		selectedRequestId,
		mediaItems,
		mediaRequests,
		discoverSections,
		discoverBlacklist,
		libraryFolders,
		languages,
		qualityProfiles,
		activities,
		releaseBlocklist,
		loadingDiscover,
		loadingBlacklist,
		loadingMediaItems,
		addingKey,
		blacklistingKey,
		removingBlacklistId,
		approvingRequestId,
		searchingItemId,
		scanningMediaItemId,
		refreshingMetadataItemId,
		savingMediaItemOptionsId,
		grabbingKey,
		deletingMediaItemId,
		assemblingMediaItemId,
		reviewingComponentDecisionId,
		cancellingActivityId,
		deletingActivityId,
		deletingReleaseBlocklistId,
		clearingReleaseBlocklist,
		canManage,
		loadingActivity,
		onAddMedia,
		onBlacklistMedia,
		onRemoveBlacklistMedia,
		onApproveMediaRequest,
		onFindReleases,
		onAutoSearchMedia,
		onRescanMediaFiles,
		onSearchMediaSubtitle = () => {},
		onGrabMediaSubtitle = () => {},
		onDeleteMediaSubtitle = () => {},
		onUpdateMediaSubtitle = () => {},
		onRefreshMediaMetadata,
		onSaveMediaItemOptions,
		onDeleteMediaFile,
		onDeleteMediaFileTrack = () => {},
		onAssembleMediaComponents = () => {},
		onReviewComponentCompatibility = () => {},
		onDeleteMedia,
		onGrabRelease,
		onRefreshActivity,
		onRefreshReleaseBlocklist,
		onCancelActivity,
		onDeleteActivity,
		onDeleteReleaseBlocklistItem,
		onClearReleaseBlocklist
	}: HomeAreaProps = $props();

	const movies = $derived(mediaItems.filter((item) => item.type === 'movie'));
	const series = $derived(mediaItems.filter((item) => item.type === 'serie'));
	const selectedMediaItem = $derived(
		selectedMediaItemId
			? mediaItems.find(
					(item) =>
						item.id === selectedMediaItemId &&
						item.type === (activeSection === 'movies' ? 'movie' : 'serie')
				)
			: undefined
	);
</script>

<section class="grid min-w-0 gap-[18px]" aria-labelledby="home-title">
	{#if activeSection === 'discover'}
		<MediaSearchPanel
			sections={discoverSections}
			{mediaItems}
			loading={loadingDiscover}
			{addingKey}
			{blacklistingKey}
			blacklist={discoverBlacklist}
			onAdd={onAddMedia}
			onBlacklist={onBlacklistMedia}
			actionLabel={canManage ? 'Add' : 'Request'}
			{canManage}
		/>
	{:else if activeSection === 'blacklist'}
		<DiscoverBlacklistArea
			items={discoverBlacklist}
			loading={loadingBlacklist}
			removingId={removingBlacklistId}
			onRemove={onRemoveBlacklistMedia}
		/>
	{:else if activeSection === 'requests'}
		<MediaRequestArea
			requests={mediaRequests}
			{selectedRequestId}
			{libraryFolders}
			{qualityProfiles}
			{canManage}
			{approvingRequestId}
			onApprove={onApproveMediaRequest}
		/>
	{:else if activeSection === 'movies' || activeSection === 'series'}
		{#if selectedMediaItemId}
			<MediaDetail
				mediaType={activeSection === 'movies' ? 'movie' : 'serie'}
				item={selectedMediaItem}
				loading={loadingMediaItems && !selectedMediaItem}
				{mediaItems}
				{libraryFolders}
				{languages}
				{qualityProfiles}
				requestedItemId={selectedMediaItemId}
				{activities}
				{searchingItemId}
				{scanningMediaItemId}
				{refreshingMetadataItemId}
				{savingMediaItemOptionsId}
				{grabbingKey}
				{addingKey}
				{deletingMediaItemId}
				{assemblingMediaItemId}
				{reviewingComponentDecisionId}
				{canManage}
				actionLabel={canManage ? 'Add' : 'Request'}
				{onAutoSearchMedia}
				{onRescanMediaFiles}
				{onSearchMediaSubtitle}
				{onGrabMediaSubtitle}
				{onDeleteMediaSubtitle}
				{onUpdateMediaSubtitle}
				{onRefreshMediaMetadata}
				{onSaveMediaItemOptions}
				{onDeleteMediaFile}
				{onDeleteMediaFileTrack}
				{onAssembleMediaComponents}
				{onReviewComponentCompatibility}
				{onDeleteMedia}
				{onGrabRelease}
				{onAddMedia}
			/>
		{:else}
			<MediaItemList
				mediaType={activeSection === 'movies' ? 'movie' : 'serie'}
				items={activeSection === 'movies' ? movies : series}
			/>
		{/if}
	{:else if activeSection === 'wanted'}
		<WantedMediaTable
			items={mediaItems}
			{languages}
			{searchingItemId}
			{grabbingKey}
			{canManage}
			{onFindReleases}
			{onGrabRelease}
		/>
	{:else}
		<ActivityList
			section={activitySection}
			{activities}
			{releaseBlocklist}
			loading={loadingActivity}
			{canManage}
			cancellingId={cancellingActivityId}
			deletingId={deletingActivityId}
			deletingBlocklistId={deletingReleaseBlocklistId}
			{clearingReleaseBlocklist}
			onRefresh={activitySection === 'blocklist' ? onRefreshReleaseBlocklist : onRefreshActivity}
			onCancel={onCancelActivity}
			onDelete={onDeleteActivity}
			{onDeleteReleaseBlocklistItem}
			{onClearReleaseBlocklist}
		/>
	{/if}
</section>
