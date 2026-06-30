<script lang="ts">
	import ActivityList from './ActivityList.svelte';
	import MediaDetail from './MediaDetail.svelte';
	import MediaItemList from './MediaItemList.svelte';
	import MediaRequestArea from './MediaRequestArea.svelte';
	import MediaSearchPanel from './MediaSearchPanel.svelte';
	import WantedMediaTable from './WantedMediaTable.svelte';
	import type {
		DownloadActivity,
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
		libraryFolders: LibraryFolder[];
		qualityProfiles: QualityProfileOption[];
		releaseResults: ReleaseSearchResults;
		activities: DownloadActivity[];
		loadingDiscover: boolean;
		addingKey?: string;
		approvingRequestId?: string;
		searchingItemId?: string;
		scanningMediaItemId?: string;
		grabbingKey?: string;
		deletingMediaItemId?: string;
		cancellingActivityId?: string;
		canManage: boolean;
		loadingActivity: boolean;
		onAddMedia: (_candidate: MediaSearchResult) => void;
		onApproveMediaRequest: (_request: MediaRequest, _approval: MediaRequestApproveRequest) => void;
		onFindReleases: (_item: MediaItem) => void;
		onAutoSearchMedia: (_item: MediaItem) => void;
		onRescanMediaFiles: (_item: MediaItem) => void;
		onDeleteMediaFile: (_item: MediaItem, _path: string) => void;
		onDeleteMedia: (_item: MediaItem) => void;
		onGrabRelease: (_item: MediaItem, _release: ReleaseCandidate) => void;
		onRefreshActivity: () => void;
		onCancelActivity: (_activity: DownloadActivity) => void;
	}

	let {
		activeSection,
		selectedMediaItemId,
		selectedRequestId,
		mediaItems,
		mediaRequests,
		discoverSections,
		libraryFolders,
		qualityProfiles,
		releaseResults,
		activities,
		loadingDiscover,
		addingKey,
		approvingRequestId,
		searchingItemId,
		scanningMediaItemId,
		grabbingKey,
		deletingMediaItemId,
		cancellingActivityId,
		canManage,
		loadingActivity,
		onAddMedia,
		onApproveMediaRequest,
		onFindReleases,
		onAutoSearchMedia,
		onRescanMediaFiles,
		onDeleteMediaFile,
		onDeleteMedia,
		onGrabRelease,
		onRefreshActivity,
		onCancelActivity
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
			onAdd={onAddMedia}
			actionLabel={canManage ? 'Add' : 'Request'}
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
	{:else if activeSection === 'movies'}
		{#if selectedMediaItemId}
			<MediaDetail
				mediaType="movie"
				item={selectedMediaItem}
				requestedItemId={selectedMediaItemId}
				{libraryFolders}
				{qualityProfiles}
				releaseResults={selectedMediaItem ? releaseResults[selectedMediaItem.id] : undefined}
				{searchingItemId}
				{scanningMediaItemId}
				{grabbingKey}
				{deletingMediaItemId}
				{canManage}
				{onFindReleases}
				{onAutoSearchMedia}
				{onRescanMediaFiles}
				{onDeleteMediaFile}
				{onDeleteMedia}
				{onGrabRelease}
			/>
		{:else}
			<MediaItemList mediaType="movie" items={movies} />
		{/if}
	{:else if activeSection === 'series'}
		{#if selectedMediaItemId}
			<MediaDetail
				mediaType="series"
				item={selectedMediaItem}
				requestedItemId={selectedMediaItemId}
				{libraryFolders}
				{qualityProfiles}
				releaseResults={selectedMediaItem ? releaseResults[selectedMediaItem.id] : undefined}
				{searchingItemId}
				{scanningMediaItemId}
				{grabbingKey}
				{deletingMediaItemId}
				{canManage}
				{onFindReleases}
				{onAutoSearchMedia}
				{onRescanMediaFiles}
				{onDeleteMediaFile}
				{onDeleteMedia}
				{onGrabRelease}
			/>
		{:else}
			<MediaItemList mediaType="series" items={series} />
		{/if}
	{:else if activeSection === 'wanted'}
		<WantedMediaTable items={wanted} {searchingItemId} {canManage} {onFindReleases} />
	{:else}
		<ActivityList
			{activities}
			loading={loadingActivity}
			{canManage}
			cancellingId={cancellingActivityId}
			onRefresh={onRefreshActivity}
			onCancel={onCancelActivity}
		/>
	{/if}
</section>
