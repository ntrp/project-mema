<script lang="ts">
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
		grabbingKey?: string;
		deletingMediaItemId?: string;
		canManage: boolean;
		onFindReleases: (_item: MediaItem) => void;
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
		grabbingKey,
		deletingMediaItemId,
		canManage,
		onFindReleases,
		onDeleteMedia,
		onGrabRelease
	}: Props = $props();

	const releaseCount = $derived(releaseResults?.releases.length ?? 0);
	const libraryFolderLabel = $derived(resolveLibraryFolderLabel(item));
	const qualityProfileLabel = $derived(resolveQualityProfileLabel(item));
	const filePaths = $derived(item?.filePaths ?? []);
	const metadataFilePaths = $derived(item?.metadataFilePaths ?? []);

	function releaseKey(mediaItem: MediaItem, release: ReleaseCandidate) {
		return `${mediaItem.id}:${release.id}`;
	}

	function sizeLabel(sizeBytes: number) {
		if (!sizeBytes) {
			return 'Unknown size';
		}
		const gib = sizeBytes / 1024 / 1024 / 1024;
		return `${gib.toFixed(gib >= 10 ? 0 : 1)} GiB`;
	}

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
		<section class="metadata-hero" aria-labelledby="home-title">
			<div class="metadata-poster">
				{#if posterUrl(item.posterPath, 'w500')}
					<img src={posterUrl(item.posterPath, 'w500')} alt="" />
				{:else}
					<div class="poster-placeholder">{mediaType === 'movie' ? 'Movie' : 'Series'}</div>
				{/if}
			</div>
			<div class="metadata-title-block">
				<h1 id="home-title">{item.title}</h1>
				<p>{mediaType === 'movie' ? 'Movie' : 'Series'}</p>
				<div class="metadata-info-bar" aria-label="Library media information">
					<span><strong>Year</strong>{item.year ?? 'Unknown'}</span>
					<span><strong>Type</strong>{item.type}</span>
					<span><strong>Status</strong>{statusLabel(item.status)}</span>
					<span><strong>Profile</strong>{qualityProfileLabel}</span>
				</div>
				{#if item.tags?.length}
					<div class="metadata-tags" aria-label="Tags">
						{#each item.tags as tag (tag)}
							<span><span class="app-icon" aria-hidden="true">sell</span>{tag}</span>
						{/each}
					</div>
				{/if}
				{#if canManage}
					<div class="metadata-actions">
						<button
							type="button"
							disabled={searchingItemId === item.id}
							onclick={() => onFindReleases(item)}
						>
							{searchingItemId === item.id ? 'Queued' : 'Find releases'}
						</button>
						<button
							type="button"
							class="danger"
							disabled={deletingMediaItemId === item.id}
							onclick={() => onDeleteMedia(item)}
						>
							{deletingMediaItemId === item.id ? 'Removing' : 'Remove'}
						</button>
					</div>
				{/if}
			</div>
		</section>

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
					</div>
				</section>

				<section aria-labelledby="media-files-title">
					<h2 id="media-files-title">Files</h2>
					{#if filePaths.length || metadataFilePaths.length}
						<div class="media-file-groups">
							{#if filePaths.length}
								<div class="media-file-group">
									<h3>Downloaded files</h3>
									<div class="media-file-list">
										{#each filePaths as path (path)}
											<code>{path}</code>
										{/each}
									</div>
								</div>
							{/if}
							{#if metadataFilePaths.length}
								<div class="media-file-group">
									<h3>Metadata files</h3>
									<div class="media-file-list">
										{#each metadataFilePaths as path (path)}
											<code>{path}</code>
										{/each}
									</div>
								</div>
							{/if}
						</div>
					{:else}
						<p class="empty">No imported files found.</p>
					{/if}
				</section>

				<section aria-labelledby="release-candidates-title">
					<div>
						<h2 id="release-candidates-title">Release Candidates</h2>
						<p>Run a release search to populate candidates.</p>
					</div>

					{#if releaseResults?.errors.length}
						<div class="inline-errors">
							{#each releaseResults.errors as searchError (searchError)}
								<p>{searchError}</p>
							{/each}
						</div>
					{/if}

					{#if releaseResults?.loaded}
						<div class="release-list standalone">
							{#each releaseResults.releases as release (release.id)}
								<div class="release-row">
									<div>
										<strong>{release.title}</strong>
										<span>
											{release.indexerName} · {sizeLabel(release.sizeBytes)}
											{release.seeders !== undefined ? ` · ${release.seeders} seeders` : ''}
										</span>
									</div>
									{#if canManage}
										<button
											type="button"
											disabled={grabbingKey === releaseKey(item, release)}
											onclick={() => onGrabRelease(item, release)}
										>
											{grabbingKey === releaseKey(item, release) ? 'Queueing' : 'Grab'}
										</button>
									{/if}
								</div>
							{:else}
								<p class="empty">No release candidates stored yet</p>
							{/each}
						</div>
					{:else}
						<p class="empty">Run a release search to populate candidates.</p>
					{/if}
				</section>
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
