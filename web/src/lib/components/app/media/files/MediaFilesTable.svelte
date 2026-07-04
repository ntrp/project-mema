<script lang="ts">
	import MediaFileSummary from '$lib/components/app/media/files/MediaFileSummary.svelte';
	import MediaFileDeleteModal from '$lib/components/app/media/files/MediaFileDeleteModal.svelte';
	import MediaFileSearchModal from '$lib/components/app/media/files/MediaFileSearchModal.svelte';
	import MediaRootPanel from '$lib/components/app/media/collection/MediaRootPanel.svelte';
	import { activityForMovie } from '$lib/components/app/activity/activityQueue';
	import { mediaFileGroups, type MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import type {
		DownloadActivity,
		Language,
		LibraryFolder,
		MediaItem,
		MediaItemUpdateRequest,
		ReleaseCandidate,
		ReleaseOverrideDetails
	} from '$lib/settings/types';

	interface Props {
		item: MediaItem;
		activities: DownloadActivity[];
		searchingItemId?: string;
		grabbingKey?: string;
		canManage: boolean;
		libraryFolders: LibraryFolder[];
		languages: Language[];
		qualityProfiles: {
			id: string;
			targetLanguages?: string[];
			removeNonEnabledLanguages?: boolean;
		}[];
		onSaveOptions: (_item: MediaItem, _request: MediaItemUpdateRequest) => void;
		onAutoSearch: (_item: MediaItem) => void;
		onDeleteFile: (_item: MediaItem, _path: string) => void;
		onGrabRelease: (
			_item: MediaItem,
			_release: ReleaseCandidate,
			_overrideMatch?: boolean,
			_details?: ReleaseOverrideDetails
		) => void;
	}

	let {
		item,
		activities,
		searchingItemId,
		grabbingKey,
		canManage,
		libraryFolders,
		languages,
		qualityProfiles,
		onSaveOptions,
		onAutoSearch,
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
						mediaItemId={item.id}
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
		{languages}
		searchContext={{ type: 'title' }}
		{grabbingKey}
		{canManage}
		onGrab={onGrabRelease}
		onClose={() => (searchOpen = false)}
	/>
{/if}
