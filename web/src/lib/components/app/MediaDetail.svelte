<script lang="ts">
	import MediaFilesTable from './MediaFilesTable.svelte';
	import MediaHero from './MediaHero.svelte';
	import ReleaseCandidatesSection from './ReleaseCandidatesSection.svelte';
	import type {
		LibraryFolder,
		MediaItem,
		MediaItemStatus,
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
		updatingMediaModeItemId?: string;
		grabbingKey?: string;
		deletingMediaItemId?: string;
		canManage: boolean;
		onFindReleases: (_item: MediaItem) => void;
		onRescanMediaFiles: (_item: MediaItem) => void;
		onUpdateMediaMode: (_item: MediaItem, _automatic: boolean) => void;
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
		updatingMediaModeItemId,
		grabbingKey,
		deletingMediaItemId,
		canManage,
		onFindReleases,
		onRescanMediaFiles,
		onUpdateMediaMode,
		onDeleteMedia,
		onGrabRelease
	}: Props = $props();

	const releaseCount = $derived(releaseResults?.releases.length ?? 0);
	const libraryFolderLabel = $derived(resolveLibraryFolderLabel(item));
	const qualityProfileLabel = $derived(resolveQualityProfileLabel(item));
	const filePaths = $derived(item?.filePaths ?? []);
	const metadataFilePaths = $derived(item?.metadataFilePaths ?? []);
	const modeLabel = $derived(item?.manual ? 'Manual' : 'Automatic');

	function posterUrl(path?: string, size = 'w780') {
		if (!path) {
			return undefined;
		}
		if (path.startsWith('http://') || path.startsWith('https://')) {
			return path;
		}
		return `https://image.tmdb.org/t/p/${size}${path}`;
	}

	function resolveLibraryFolderLabel(mediaItem?: MediaItem) {
		if (!mediaItem) {
			return 'Not set';
		}
		return (
			mediaItem.mediaFolderPath ??
			mediaItem.libraryFolderPath ??
			libraryFolders.find((folder) => folder.id === mediaItem.libraryFolderId)?.path ??
			'Not set'
		);
	}

	function resolveQualityProfileLabel(mediaItem?: MediaItem) {
		if (!mediaItem) {
			return 'Not set';
		}
		return (
			mediaItem.qualityProfileName ??
			qualityProfiles.find((profile) => profile.id === mediaItem.qualityProfileId)?.name ??
			'Not set'
		);
	}

	function statusLabel(status: MediaItemStatus) {
		switch (status) {
			case 'downloaded':
				return 'Downloaded';
			case 'downloading':
				return 'Downloading';
			default:
				return 'Missing';
		}
	}
</script>

{#if item}
	<div
		class="metadata-detail media-library-detail"
		style:--backdrop-url={posterUrl(item.posterPath, 'original')
			? `url("${posterUrl(item.posterPath, 'original')}")`
			: undefined}
	>
		<MediaHero
			{mediaType}
			{item}
			{qualityProfileLabel}
			{canManage}
			{searchingItemId}
			{scanningMediaItemId}
			{deletingMediaItemId}
			{onFindReleases}
			{onRescanMediaFiles}
			{onDeleteMedia}
		/>

		<div class="metadata-body">
			<main class="metadata-main">
				<section aria-labelledby="library-status-title">
					<h2 id="library-status-title">Library Status</h2>
					<div class="metadata-facts-grid" aria-label="Library status facts">
						<div>
							<strong>{statusLabel(item.status)}</strong>
							<span>Status</span>
						</div>
						<div>
							<strong>{releaseCount}</strong>
							<span>Release candidates</span>
						</div>
						<div>
							<strong>{item.year ?? 'Unknown'}</strong>
							<span>Year</span>
						</div>
						<div>
							<strong>{qualityProfileLabel}</strong>
							<span>Profile</span>
						</div>
						<div>
							<strong>{libraryFolderLabel}</strong>
							<span>Media folder</span>
						</div>
						<div>
							<strong>{item.monitored ? 'Monitored' : 'Paused'}</strong>
							<span>Monitor state</span>
						</div>
						<div>
							<strong>{modeLabel}</strong>
							<span>Mode</span>
							{#if canManage}
								<button
									type="button"
									class="secondary compact-action"
									disabled={updatingMediaModeItemId === item.id}
									onclick={() => onUpdateMediaMode(item, item.manual)}
								>
									{#if updatingMediaModeItemId === item.id}
										Updating
									{:else}
										Switch to {item.manual ? 'automatic' : 'manual'}
									{/if}
								</button>
							{/if}
						</div>
					</div>
				</section>

				<MediaFilesTable {filePaths} {metadataFilePaths} />

				<ReleaseCandidatesSection
					{item}
					{releaseResults}
					{grabbingKey}
					{canManage}
					{onGrabRelease}
				/>
			</main>
		</div>
	</div>
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
