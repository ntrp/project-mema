<script lang="ts">
	import LibraryMediaStatusSection from './LibraryMediaStatusSection.svelte';
	import MediaFilesTable from './MediaFilesTable.svelte';
	import MediaMetadataCore from './MediaMetadataCore.svelte';
	import MediaMetadataHero from './MediaMetadataHero.svelte';
	import ReleaseCandidatesSection from './ReleaseCandidatesSection.svelte';
	import {
		imageUrl,
		libraryFolderLabel,
		mediaMetadataDetail,
		monitorModeLabel,
		qualityProfileLabel
	} from './mediaDetail';
	import type {
		LibraryFolder,
		MediaItem,
		MediaType,
		QualityProfileOption,
		ReleaseCandidate,
		ReleaseSearchState
	} from '$lib/settings/types';

	interface Props {
		mediaType: MediaType;
		item?: MediaItem;
		requestedItemId: string;
		libraryFolders: LibraryFolder[];
		qualityProfiles: QualityProfileOption[];
		releaseResults?: ReleaseSearchState;
		searchingItemId?: string;
		scanningMediaItemId?: string;
		grabbingKey?: string;
		deletingMediaItemId?: string;
		canManage: boolean;
		onFindReleases: (_item: MediaItem) => void;
		onAutoSearchMedia: (_item: MediaItem) => void;
		onRescanMediaFiles: (_item: MediaItem) => void;
		onDeleteMediaFile: (_item: MediaItem, _path: string) => void;
		onDeleteMedia: (_item: MediaItem) => void;
		onGrabRelease: (_item: MediaItem, _release: ReleaseCandidate) => void;
	}

	let {
		mediaType,
		item,
		requestedItemId,
		libraryFolders,
		qualityProfiles,
		releaseResults,
		searchingItemId,
		scanningMediaItemId,
		grabbingKey,
		deletingMediaItemId,
		canManage,
		onFindReleases,
		onAutoSearchMedia,
		onRescanMediaFiles,
		onDeleteMediaFile,
		onDeleteMedia,
		onGrabRelease
	}: Props = $props();

	const releaseCount = $derived(releaseResults?.releases.length ?? 0);
	const resolvedLibraryFolderLabel = $derived(libraryFolderLabel(item, libraryFolders));
	const resolvedQualityProfileLabel = $derived(qualityProfileLabel(item, qualityProfiles));
	const resolvedMonitorModeLabel = $derived(monitorModeLabel(item));
	const detail = $derived(item ? mediaMetadataDetail(item) : undefined);
</script>

{#if item && detail}
	<section
		class="metadata-detail media-library-detail"
		aria-labelledby="library-media-title"
		style:--backdrop-url={imageUrl(detail.backdropPath, 'original')
			? `url("${imageUrl(detail.backdropPath, 'original')}")`
			: undefined}
	>
		<MediaMetadataHero {detail} titleId="library-media-title">
			{#snippet actions()}
				{#if canManage}
					<button
						type="button"
						disabled={searchingItemId === item.id}
						onclick={() => onFindReleases(item)}
					>
						{searchingItemId === item.id ? 'Queued' : 'Find releases'}
					</button>
					<button
						type="button"
						class="secondary"
						disabled={scanningMediaItemId === item.id || !item.mediaFolderPath}
						onclick={() => onRescanMediaFiles(item)}
					>
						<span class="app-icon" aria-hidden="true">sync</span>
						<span>{scanningMediaItemId === item.id ? 'Scanning' : 'Rescan files'}</span>
					</button>
					<button
						type="button"
						class="danger"
						disabled={deletingMediaItemId === item.id}
						onclick={() => onDeleteMedia(item)}
					>
						{deletingMediaItemId === item.id ? 'Removing' : 'Remove'}
					</button>
				{/if}
			{/snippet}
		</MediaMetadataHero>

		<div class="metadata-body">
			<main class="metadata-main">
				<MediaMetadataCore {detail} />
				<LibraryMediaStatusSection
					{item}
					{releaseCount}
					qualityProfileLabel={resolvedQualityProfileLabel}
					libraryFolderLabel={resolvedLibraryFolderLabel}
					monitorModeLabel={resolvedMonitorModeLabel}
				/>
				<MediaFilesTable
					{item}
					{releaseResults}
					{searchingItemId}
					{grabbingKey}
					{canManage}
					onAutoSearch={onAutoSearchMedia}
					onManualSearch={onFindReleases}
					onDeleteFile={onDeleteMediaFile}
					{onGrabRelease}
				/>
				<ReleaseCandidatesSection
					{item}
					{releaseResults}
					{grabbingKey}
					{canManage}
					{onGrabRelease}
				/>
			</main>
		</div>
	</section>
{:else}
	<div class="detail-stack">
		<section class="panel">
			<div class="page-heading">
				<p>{mediaType === 'movie' ? 'Movies' : 'Series'}</p>
				<h1 id="home-title">Media item not found</h1>
			</div>
			<p class="empty">No monitored {mediaType} matches {requestedItemId}.</p>
		</section>
	</div>
{/if}
