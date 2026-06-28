<script lang="ts">
	import type { PathnameWithSearchOrHash } from '$app/types';

	import ActivityList from './ActivityList.svelte';
	import MediaDetail from './MediaDetail.svelte';
	import MediaItemList from './MediaItemList.svelte';
	import MediaSearchPanel from './MediaSearchPanel.svelte';
	import SidebarMenu from './SidebarMenu.svelte';
	import type {
		DownloadActivity,
		HomeSection,
		MediaItem,
		MediaSearchRequest,
		MediaSearchResult,
		ReleaseCandidate,
		ReleaseSearchResults
	} from '$lib/settings/types';

	interface Props {
		activeSection: HomeSection;
		selectedMediaItemId?: string;
		mediaItems: MediaItem[];
		mediaSearchResults: MediaSearchResult[];
		releaseResults: ReleaseSearchResults;
		activities: DownloadActivity[];
		searchingMedia: boolean;
		addingKey?: string;
		searchingItemId?: string;
		grabbingKey?: string;
		deletingMediaItemId?: string;
		loadingActivity: boolean;
		onSelect: (_section: HomeSection) => void;
		onSearchMedia: (_request: MediaSearchRequest) => void;
		onAddMedia: (_candidate: MediaSearchResult) => void;
		onFindReleases: (_item: MediaItem) => void;
		onDeleteMedia: (_item: MediaItem) => void;
		onGrabRelease: (_item: MediaItem, _release: ReleaseCandidate) => void;
		onRefreshActivity: () => void;
	}

	let {
		activeSection,
		selectedMediaItemId,
		mediaItems,
		mediaSearchResults,
		releaseResults,
		activities,
		searchingMedia,
		addingKey,
		searchingItemId,
		grabbingKey,
		deletingMediaItemId,
		loadingActivity,
		onSelect,
		onSearchMedia,
		onAddMedia,
		onFindReleases,
		onDeleteMedia,
		onGrabRelease,
		onRefreshActivity
	}: Props = $props();

	const homeItems = [
		{ value: 'explore', label: 'Explore', meta: 'Digest', href: '/explore' },
		{ value: 'movies', label: 'Movies', meta: 'Anime included', href: '/movies' },
		{ value: 'series', label: 'Series', meta: 'Anime included', href: '/series' },
		{ value: 'activity', label: 'Activity', meta: 'Queue', href: '/activity' }
	] satisfies { value: HomeSection; label: string; meta: string; href: PathnameWithSearchOrHash }[];

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

<div class="workspace-layout">
	<SidebarMenu
		title="Library"
		items={homeItems}
		active={activeSection}
		onSelect={(section) => onSelect(section as HomeSection)}
	/>

	<section class="workspace-main" aria-labelledby="home-title">
		{#if activeSection === 'explore'}
			<MediaSearchPanel
				results={mediaSearchResults}
				searching={searchingMedia}
				{addingKey}
				mediaItemsCount={mediaItems.length}
				onSearch={onSearchMedia}
				onAdd={onAddMedia}
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
</div>
