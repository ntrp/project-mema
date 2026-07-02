<script lang="ts">
	import ActivityList from '../activity/ActivityList.svelte';
	import DiscoverBlacklistArea from '../discovery/DiscoverBlacklistArea.svelte';
	import MediaDetail from '../media/MediaDetail.svelte';
	import MediaItemList from './MediaItemList.svelte';
	import MediaRequestArea from '../requests/MediaRequestArea.svelte';
	import MediaSearchPanel from '../discovery/MediaSearchPanel.svelte';
	import WantedMediaTable from './WantedMediaTable.svelte';
	import type { HomeAreaProps } from './homeAreaTypes';

	let {
		activeSection,
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
		loadingDiscover,
		loadingBlacklist,
		loadingMediaItems,
		addingKey,
		blacklistingKey,
		removingBlacklistId,
		approvingRequestId,
		searchingItemId,
		refreshingMetadataItemId,
		savingMediaItemOptionsId,
		grabbingKey,
		deletingMediaItemId,
		cancellingActivityId,
		deletingActivityId,
		canManage,
		loadingActivity,
		onAddMedia,
		onBlacklistMedia,
		onRemoveBlacklistMedia,
		onApproveMediaRequest,
		onFindReleases,
		onAutoSearchMedia,
		onRefreshMediaMetadata,
		onSaveMediaItemOptions,
		onDeleteMediaFile,
		onDeleteMedia,
		onGrabRelease,
		onRefreshActivity,
		onCancelActivity,
		onDeleteActivity
	}: HomeAreaProps = $props();

	const movies = $derived(mediaItems.filter((item) => item.type === 'movie'));
	const series = $derived(mediaItems.filter((item) => item.type === 'series'));
	const wanted = $derived(mediaItems.filter((item) => item.status === 'missing'));
	const selectedMediaItem = $derived(
		selectedMediaItemId
			? mediaItems.find(
					(item) =>
						item.id === selectedMediaItemId &&
						item.type === (activeSection === 'movies' ? 'movie' : 'series')
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
				mediaType={activeSection === 'movies' ? 'movie' : 'series'}
				item={selectedMediaItem}
				loading={loadingMediaItems && !selectedMediaItem}
				{mediaItems}
				{libraryFolders}
				{languages}
				{qualityProfiles}
				requestedItemId={selectedMediaItemId}
				{activities}
				{searchingItemId}
				{refreshingMetadataItemId}
				{savingMediaItemOptionsId}
				{grabbingKey}
				{addingKey}
				{deletingMediaItemId}
				{canManage}
				actionLabel={canManage ? 'Add' : 'Request'}
				{onAutoSearchMedia}
				{onRefreshMediaMetadata}
				{onSaveMediaItemOptions}
				{onDeleteMediaFile}
				{onDeleteMedia}
				{onGrabRelease}
				{onAddMedia}
			/>
		{:else}
			<MediaItemList
				mediaType={activeSection === 'movies' ? 'movie' : 'series'}
				items={activeSection === 'movies' ? movies : series}
			/>
		{/if}
	{:else if activeSection === 'wanted'}
		<WantedMediaTable items={wanted} {searchingItemId} {canManage} {onFindReleases} />
	{:else}
		<ActivityList
			{activities}
			loading={loadingActivity}
			{canManage}
			cancellingId={cancellingActivityId}
			deletingId={deletingActivityId}
			onRefresh={onRefreshActivity}
			onCancel={onCancelActivity}
			onDelete={onDeleteActivity}
		/>
	{/if}
</section>
