<script lang="ts">
	import ActivityList from './ActivityList.svelte';
	import MediaDetail from './MediaDetail.svelte';
	import MediaItemList from './MediaItemList.svelte';
	import MediaRequestArea from './MediaRequestArea.svelte';
	import MediaSearchPanel from './MediaSearchPanel.svelte';
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
		grabbingKey?: string;
		deletingMediaItemId?: string;
		canManage: boolean;
		loadingActivity: boolean;
		onAddMedia: (_candidate: MediaSearchResult) => void;
		onApproveMediaRequest: (_request: MediaRequest, _approval: MediaRequestApproveRequest) => void;
		onFindReleases: (_item: MediaItem) => void;
		onDeleteMedia: (_item: MediaItem) => void;
		onGrabRelease: (_item: MediaItem, _release: ReleaseCandidate) => void;
		onRefreshActivity: () => void;
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
		grabbingKey,
		deletingMediaItemId,
		canManage,
		loadingActivity,
		onAddMedia,
		onApproveMediaRequest,
		onFindReleases,
		onDeleteMedia,
		onGrabRelease,
		onRefreshActivity
	}: Props = $props();

	const movies = $derived(mediaItems.filter((item) => item.type === 'movie'));
	const series = $derived(mediaItems.filter((item) => item.type === 'series'));
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
				releaseResults={selectedMediaItem ? releaseResults[selectedMediaItem.id] : undefined}
				{searchingItemId}
				{grabbingKey}
				{deletingMediaItemId}
				{canManage}
				{onFindReleases}
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
				releaseResults={selectedMediaItem ? releaseResults[selectedMediaItem.id] : undefined}
				{searchingItemId}
				{grabbingKey}
				{deletingMediaItemId}
				{canManage}
				{onFindReleases}
				{onDeleteMedia}
				{onGrabRelease}
			/>
		{:else}
			<MediaItemList mediaType="series" items={series} />
		{/if}
	{:else}
		<ActivityList {activities} loading={loadingActivity} onRefresh={onRefreshActivity} />
	{/if}
</section>
