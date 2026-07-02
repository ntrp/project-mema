<script lang="ts">
	import MediaFileSummary from './MediaFileSummary.svelte';
	import MediaFileDeleteModal from './MediaFileDeleteModal.svelte';
	import MediaFileSearchModal from './MediaFileSearchModal.svelte';
	import MediaRootPanel from './MediaRootPanel.svelte';
	import { activityForMovie } from '../activity/activityQueue';
	import { mediaFileGroups, type MediaFileRow } from './mediaFiles';
	import type {
		DownloadActivity,
		LibraryFolder,
		MediaItem,
		MediaItemUpdateRequest,
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
		libraryFolders: LibraryFolder[];
		qualityProfiles: {
			id: string;
			targetLanguages?: string[];
			removeNonEnabledLanguages?: boolean;
		}[];
		onSaveOptions: (_item: MediaItem, _request: MediaItemUpdateRequest) => void;
		onAutoSearch: (_item: MediaItem) => void;
		onManualSearch: (_item: MediaItem, _query?: string) => void;
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
		libraryFolders,
		qualityProfiles,
		onSaveOptions,
		onAutoSearch,
		onManualSearch,
		onDeleteFile,
		onGrabRelease
	}: Props = $props();

	let deleteRow = $state<MediaFileRow | undefined>();
	let searchOpen = $state(false);
	const groups = $derived(mediaFileGroups(item, qualityProfiles));
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
		<MediaRootPanel {item} {libraryFolders} {canManage} {onSaveOptions} />
		{#each groups as group (group.key)}
			<div class="grid" aria-label={group.title}>
				{#each group.rows as row (row.key)}
					<MediaFileSummary
						{row}
						{activityStatus}
						{canManage}
						searching={busy}
						fileLabel="Movie file"
						missingLabel="No matched file for this movie"
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

{#if searchOpen}
	<MediaFileSearchModal
		{item}
		{releaseResults}
		searchContext={{ type: 'title' }}
		searching={searchingItemId === item.id}
		{grabbingKey}
		{canManage}
		onSearch={onManualSearch}
		onGrab={onGrabRelease}
		onClose={() => (searchOpen = false)}
	/>
{/if}
