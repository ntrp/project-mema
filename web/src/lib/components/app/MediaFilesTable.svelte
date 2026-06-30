<script lang="ts">
	import MediaFileInfoModal from './MediaFileInfoModal.svelte';
	import MediaFileSearchModal from './MediaFileSearchModal.svelte';
	import { activityForMovie } from './activityQueue';
	import { mediaFileGroups, type MediaFileRow } from './mediaFiles';
	import type {
		DownloadActivity,
		MediaItem,
		ReleaseCandidate,
		ReleaseSearchState
	} from '$lib/settings/types';

	interface Props {
		item: MediaItem;
		releaseResults?: ReleaseSearchState;
		activities: DownloadActivity[];
		searchingItemId?: string;
		grabbingKey?: string;
		canManage: boolean;
		onAutoSearch: (_item: MediaItem) => void;
		onManualSearch: (_item: MediaItem) => void;
		onDeleteFile: (_item: MediaItem, _path: string) => void;
		onGrabRelease: (_item: MediaItem, _release: ReleaseCandidate) => void;
	}

	let {
		item,
		releaseResults,
		activities,
		searchingItemId,
		grabbingKey,
		canManage,
		onAutoSearch,
		onManualSearch,
		onDeleteFile,
		onGrabRelease
	}: Props = $props();

	let detailRow = $state<MediaFileRow | undefined>();
	let searchOpen = $state(false);
	const groups = $derived(mediaFileGroups(item));
	const activityStatus = $derived(
		item.type === 'movie' ? activityForMovie(activities, item.id) : undefined
	);
	const busy = $derived(
		searchingItemId === item.id ||
			activityStatus?.status === 'queued' ||
			activityStatus?.status === 'grabbed' ||
			activityStatus?.status === 'downloading'
	);

	function confirmDelete(row: MediaFileRow) {
		if (!row.path) return;
		if (window.confirm(`Delete ${row.relativePath}?`)) {
			onDeleteFile(item, row.path);
		}
	}
</script>

<section aria-labelledby="media-files-title">
	<h2 id="media-files-title">Files</h2>
	<div class="file-section-stack">
		{#each groups as group (group.key)}
			<section class="panel media-file-panel" aria-labelledby={`files-${group.key}`}>
				<h3 id={`files-${group.key}`}>{group.title}</h3>
				<div class="table-wrap media-files-table">
					<table>
						<thead>
							<tr>
								<th>Relative Path</th>
								<th>Video Codec</th>
								<th>Audio Info</th>
								<th>Size</th>
								<th>Languages</th>
								<th>Quality</th>
								<th>Status</th>
								<th>Formats</th>
								<th>Score</th>
								<th>Actions</th>
							</tr>
						</thead>
						<tbody>
							{#each group.rows as row (row.key)}
								<tr class:missing-file={!row.exists}>
									<td>
										<strong>{row.relativePath}</strong>
										{#if row.episodeNumber}
											<small>S{row.seasonNumber}E{row.episodeNumber} {row.episodeTitle ?? ''}</small
											>
										{/if}
									</td>
									<td>{row.videoCodec}</td>
									<td>{row.audioInfo}</td>
									<td>{row.size}</td>
									<td>{row.languages}</td>
									<td>{row.quality}</td>
									<td>
										{#if activityStatus}
											<span
												class="activity-status-chip"
												class:activity-failed={activityStatus.status === 'failed'}
											>
												<span class="app-icon" aria-hidden="true">sync</span>
												{activityStatus.label}
											</span>
										{:else}
											-
										{/if}
									</td>
									<td>
										{#if row.formats.length}
											<div class="format-chip-list">
												{#each row.formats as format (format)}
													<span>{format}</span>
												{/each}
											</div>
										{:else}
											-
										{/if}
									</td>
									<td>{row.score}</td>
									<td class="row-actions media-file-actions">
										{#if row.exists}
											<button
												type="button"
												class="secondary icon-button"
												aria-label="File info"
												onclick={() => (detailRow = row)}
											>
												<span class="app-icon" aria-hidden="true">info</span>
											</button>
										{/if}
										<button
											type="button"
											class="secondary icon-button"
											aria-label="Automatic search"
											title="Automatic search"
											disabled={!canManage || busy}
											onclick={() => onAutoSearch(item)}
										>
											<span class="app-icon" aria-hidden="true">search</span>
										</button>
										<button
											type="button"
											class="secondary icon-button"
											aria-label="Manual search"
											title="Manual search"
											disabled={busy}
											onclick={() => (searchOpen = true)}
										>
											<span class="app-icon" aria-hidden="true">person</span>
										</button>
										{#if row.exists}
											<button
												type="button"
												class="danger icon-button"
												aria-label="Delete file"
												disabled={!canManage || !row.path}
												onclick={() => confirmDelete(row)}
											>
												<span class="app-icon" aria-hidden="true">delete</span>
											</button>
										{/if}
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</section>
		{/each}
	</div>
</section>

{#if detailRow}
	<MediaFileInfoModal row={detailRow} onClose={() => (detailRow = undefined)} />
{/if}

{#if searchOpen}
	<MediaFileSearchModal
		{item}
		{releaseResults}
		searching={searchingItemId === item.id}
		{grabbingKey}
		{canManage}
		onSearch={onManualSearch}
		onGrab={onGrabRelease}
		onClose={() => (searchOpen = false)}
	/>
{/if}
