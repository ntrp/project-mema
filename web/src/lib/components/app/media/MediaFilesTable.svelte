<script lang="ts">
	import MediaEpisodeFileSummary from './MediaEpisodeFileSummary.svelte';
	import MediaFileDeleteModal from './MediaFileDeleteModal.svelte';
	import MediaFileInfoModal from './MediaFileInfoModal.svelte';
	import MediaFileSearchModal from './MediaFileSearchModal.svelte';
	import { activityForMovie } from '../activity/activityQueue';
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
	let deleteRow = $state<MediaFileRow | undefined>();
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

	function requestDelete(row: MediaFileRow) {
		if (!row.path) return;
		deleteRow = row;
	}

	function confirmDelete() {
		if (!deleteRow?.path) return;
		onDeleteFile(item, deleteRow.path);
		deleteRow = undefined;
	}
</script>

<section aria-labelledby="media-files-title">
	<h2 id="media-files-title" class="m-0 text-3xl font-semibold text-foreground">Files</h2>
	<div class="grid gap-3.5">
		{#each groups as group (group.key)}
			<div class="grid px-4.5" aria-label={group.title}>
				{#each group.rows as row (row.key)}
					<MediaEpisodeFileSummary
						{row}
						{activityStatus}
						{canManage}
						searching={busy}
						fileLabel="Movie file"
						missingLabel="No matched file for this movie"
						onInfo={(nextRow) => (detailRow = nextRow)}
						onAutoSearch={() => onAutoSearch(item)}
						onManualSearch={() => (searchOpen = true)}
						onDelete={requestDelete}
					/>
				{/each}
			</div>
		{/each}
	</div>
</section>

{#if deleteRow}
	<MediaFileDeleteModal
		row={deleteRow}
		onCancel={() => (deleteRow = undefined)}
		onConfirm={confirmDelete}
	/>
{/if}

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
