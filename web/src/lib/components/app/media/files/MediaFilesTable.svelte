<script lang="ts">
	import MediaFileSummary from '$lib/components/app/media/files/MediaFileSummary.svelte';
	import MediaFileDeleteModal from '$lib/components/app/media/files/MediaFileDeleteModal.svelte';
	import MediaFilesHeader from '$lib/components/app/media/files/MediaFilesHeader.svelte';
	import MediaRenameModal from '$lib/components/app/media/files/MediaRenameModal.svelte';
	import MediaFileSearchModal from '$lib/components/app/media/files/MediaFileSearchModal.svelte';
	import SubtitleSearchModal from '$lib/components/app/media/subtitle-search/SubtitleSearchModal.svelte';
	import MediaRootPanel from '$lib/components/app/media/collection/MediaRootPanel.svelte';
	import { activityForMovie } from '$lib/components/app/activity/activityQueue';
	import { mediaFileGroups, type MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import type { MediaFilesTableProps as Props } from '$lib/components/app/media/file-data/mediaFileComponentTypes';
	import type { MediaItemSubtitleSelectionRequest } from '$lib/settings/types';

	let {
		item,
		activities,
		searchingItemId,
		scanningMediaItemId,
		grabbingKey,
		canManage,
		libraryFolders,
		languages,
		qualityProfiles,
		onSaveOptions,
		onAutoSearch,
		onRescanMediaFiles,
		onSearchSubtitle,
		onGrabSubtitle,
		onDeleteSubtitle,
		onUpdateSubtitle,
		onDeleteFile,
		onDeleteFileTrack,
		onGrabRelease
	}: Props = $props();

	let deleteRow = $state<MediaFileRow | undefined>();
	let searchOpen = $state(false);
	let renameOpen = $state(false);
	let subtitleSearch = $state<{ row: MediaFileRow; languageId: string } | undefined>();
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

	async function searchSubtitle(row: MediaFileRow, languageId?: string) {
		if (!row.path) return;
		await onSearchSubtitle(item, { languageId, filePath: row.path });
	}

	function deleteSubtitle(subtitleId: string) {
		return onDeleteSubtitle(item, subtitleId);
	}

	function updateSubtitle(subtitleId: string, request: MediaItemSubtitleSelectionRequest) {
		return onUpdateSubtitle(item, subtitleId, request);
	}

	function openSubtitleSearch(row: MediaFileRow, languageId?: string) {
		subtitleSearch = {
			row,
			languageId: languageId ?? row.subtitleSatisfaction?.wantedLanguages[0] ?? 'english'
		};
	}

	function renameApplied() {
		onRescanMediaFiles(item);
	}
</script>

<section aria-labelledby="media-files-title">
	<MediaFilesHeader
		{item}
		{canManage}
		{scanningMediaItemId}
		onRename={() => (renameOpen = true)}
		{onRescanMediaFiles}
	/>
	<div class="grid gap-3.5">
		<MediaRootPanel {item} {libraryFolders} {canManage} {onSaveOptions} />
		{#each groups as group (group.key)}
			<div class="grid" aria-label={group.title}>
				{#each group.rows as row (row.key)}
					<MediaFileSummary
						mediaItemId={item.id}
						mediaTitle={item.title}
						{row}
						{activityStatus}
						{canManage}
						searching={busy}
						fileLabel="Movie file"
						missingLabel="No matched file for this movie"
						onAutoSearch={() => onAutoSearch(item)}
						onManualSearch={() => (searchOpen = true)}
						onSearchSubtitle={searchSubtitle}
						onManualSubtitleSearch={openSubtitleSearch}
						onDeleteSubtitle={(subtitle) => deleteSubtitle(subtitle.id)}
						onUpdateSubtitle={(subtitle, request) => updateSubtitle(subtitle.id, request)}
						onDeleteTrack={(row, request) =>
							onDeleteFileTrack(item, { ...request, path: row.path ?? '' })}
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

{#if subtitleSearch}
	<SubtitleSearchModal
		{item}
		row={subtitleSearch.row}
		languageId={subtitleSearch.languageId}
		{canManage}
		onGrab={onGrabSubtitle}
		onClose={() => (subtitleSearch = undefined)}
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

{#if renameOpen}
	<MediaRenameModal {item} onClose={() => (renameOpen = false)} onApplied={renameApplied} />
{/if}
