<script lang="ts">
	import { resolve } from '$app/paths';
	import type {
		MediaItem,
		MediaType,
		ReleaseCandidate,
		ReleaseSearchState
	} from '$lib/settings/types';

	interface Props {
		mediaType: MediaType;
		item?: MediaItem;
		requestedItemId: string;
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
		releaseResults,
		searchingItemId,
		grabbingKey,
		deletingMediaItemId,
		canManage,
		onFindReleases,
		onDeleteMedia,
		onGrabRelease
	}: Props = $props();

	const sectionLabel = $derived(mediaType === 'movie' ? 'Movies' : 'Series');
	const releaseCount = $derived(releaseResults?.releases.length ?? 0);

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
</script>

{#if item}
	<div class="detail-stack">
		<a class="back-link" href={resolve(mediaType === 'movie' ? '/movies' : '/series')}>
			Back to {sectionLabel}
		</a>

		<section class="media-detail-hero" aria-labelledby="home-title">
			<div class="poster-placeholder detail-poster">
				{mediaType === 'movie' ? 'Movie' : 'Series'}
			</div>
			<div class="media-detail-main">
				<div class="page-heading">
					<p>{sectionLabel}</p>
					<h1 id="home-title">{item.title}</h1>
				</div>
				<div class="detail-meta">
					<span>{item.year ?? 'Unknown year'}</span>
					<span>{item.type}</span>
					<span>{item.monitored ? 'Monitored' : 'Paused'}</span>
				</div>
				{#if item.tags?.length}
					<div class="media-tags" aria-label="Tags">
						{#each item.tags as tag (tag)}
							<span>{tag}</span>
						{/each}
					</div>
				{/if}
				{#if canManage}
					<div class="detail-actions">
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

		<div class="detail-grid">
			<section class="panel">
				<div class="section-heading">
					<div>
						<p class="section-kicker">Overview</p>
						<h2>Library status</h2>
					</div>
				</div>
				<div class="status-grid">
					<div>
						<strong>{item.monitored ? 'Monitored' : 'Paused'}</strong>
						<span>Monitor state</span>
					</div>
					<div>
						<strong>{releaseCount}</strong>
						<span>Release candidates</span>
					</div>
					<div>
						<strong>{item.year ?? 'Unknown'}</strong>
						<span>Year</span>
					</div>
				</div>
			</section>

			<section class="panel">
				<div class="section-heading">
					<div>
						<p class="section-kicker">Search</p>
						<h2>Release candidates</h2>
					</div>
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
		</div>
	</div>
{:else}
	<div class="detail-stack">
		<a class="back-link" href={resolve(mediaType === 'movie' ? '/movies' : '/series')}>
			Back to {sectionLabel}
		</a>
		<section class="panel">
			<div class="page-heading">
				<p>{sectionLabel}</p>
				<h1 id="home-title">Media item not found</h1>
			</div>
			<p class="empty">No monitored {mediaType} matches {requestedItemId}.</p>
		</section>
	</div>
{/if}
