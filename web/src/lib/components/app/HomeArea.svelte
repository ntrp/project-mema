<script lang="ts">
	import ActivityList from './ActivityList.svelte';
	import DiscoverBlacklistArea from './DiscoverBlacklistArea.svelte';
	import MediaDetail from './MediaDetail.svelte';
	import MediaItemList from './MediaItemList.svelte';
	import MediaRequestArea from './MediaRequestArea.svelte';
	import MediaSearchPanel from './MediaSearchPanel.svelte';
	import WantedMediaTable from './WantedMediaTable.svelte';
	import type {
		DownloadActivity,
		DiscoverBlacklistItem,
		HomeSection,
		LibraryFolder,
		MediaDiscoverSection,
		MediaItem,
		MediaRequest,
		MediaRequestApproveRequest,
		MediaSearchResult,
		QualityProfileOption,
		ReleaseCandidate,
		ReleaseSearchResults
	} from '$lib/settings/types';

	interface Props {
		activeSection: HomeSection;
		selectedMediaItemId?: string;
		selectedRequestId?: string;
		mediaItems: MediaItem[];
		mediaRequests: MediaRequest[];
		discoverSections: MediaDiscoverSection[];
		discoverBlacklist: DiscoverBlacklistItem[];
		libraryFolders: LibraryFolder[];
		qualityProfiles: QualityProfileOption[];
		releaseResults: ReleaseSearchResults;
		activities: DownloadActivity[];
		loadingDiscover: boolean;
		loadingBlacklist: boolean;
		addingKey?: string;
		blacklistingKey?: string;
		removingBlacklistId?: string;
		approvingRequestId?: string;
		searchingItemId?: string;
		grabbingKey?: string;
		deletingMediaItemId?: string;
		cancellingActivityId?: string;
		deletingActivityId?: string;
		canManage: boolean;
		loadingActivity: boolean;
		onAddMedia: (_candidate: MediaSearchResult) => void;
		onBlacklistMedia: (_candidate: MediaSearchResult) => void;
		onRemoveBlacklistMedia: (_item: DiscoverBlacklistItem) => void;
		onApproveMediaRequest: (_request: MediaRequest, _approval: MediaRequestApproveRequest) => void;
		onFindReleases: (_item: MediaItem) => void;
		onAutoSearchMedia: (_item: MediaItem) => void;
		onDeleteMediaFile: (_item: MediaItem, _path: string) => void;
		onDeleteMedia: (_item: MediaItem) => void;
		onGrabRelease: (_item: MediaItem, _release: ReleaseCandidate) => void;
		onRefreshActivity: () => void;
		onCancelActivity: (_activity: DownloadActivity) => void;
		onDeleteActivity: (_activity: DownloadActivity) => void;
	}

	let {
		activeSection,
		selectedMediaItemId,
		selectedRequestId,
		mediaItems,
		mediaRequests,
		discoverSections,
		discoverBlacklist,
		libraryFolders,
		qualityProfiles,
		releaseResults,
		activities,
		loadingDiscover,
		loadingBlacklist,
		addingKey,
		blacklistingKey,
		removingBlacklistId,
		approvingRequestId,
		searchingItemId,
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
		onDeleteMediaFile,
		onDeleteMedia,
		onGrabRelease,
		onRefreshActivity,
		onCancelActivity,
		onDeleteActivity
	}: Props = $props();

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

<section class="workspace-main" aria-labelledby="home-title">
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
				{mediaItems}
				requestedItemId={selectedMediaItemId}
				releaseResults={selectedMediaItem ? releaseResults[selectedMediaItem.id] : undefined}
				{activities}
				{searchingItemId}
				{grabbingKey}
				{addingKey}
				{deletingMediaItemId}
				{canManage}
				actionLabel={canManage ? 'Add' : 'Request'}
				{onFindReleases}
				{onAutoSearchMedia}
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
