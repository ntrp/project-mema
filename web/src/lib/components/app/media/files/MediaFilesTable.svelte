<script lang="ts">
	import MediaFileSummary from '$lib/components/app/media/files/MediaFileSummary.svelte';
	import MediaFileDeleteModal from '$lib/components/app/media/files/MediaFileDeleteModal.svelte';
	import MediaRenameApplyModal from '$lib/components/app/media/files/MediaRenameApplyModal.svelte';
	import MediaRenamePreviewPanel from '$lib/components/app/media/files/MediaRenamePreviewPanel.svelte';
	import MediaFileSearchModal from '$lib/components/app/media/files/MediaFileSearchModal.svelte';
	import SubtitleSearchModal from '$lib/components/app/media/subtitle-search/SubtitleSearchModal.svelte';
	import MediaRootPanel from '$lib/components/app/media/collection/MediaRootPanel.svelte';
	import { activityForMovie } from '$lib/components/app/activity/activityQueue';
	import { mediaFileGroups, type MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import type { MediaFilesTableProps as Props } from '$lib/components/app/media/file-data/mediaFileComponentTypes';
	import { applyMediaRename, previewMediaRename } from '$lib/settings/api';
	import type {
		MediaItemSubtitleSelectionRequest,
		MediaRenamePreviewRow
	} from '$lib/settings/types';

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
	let subtitleSearch = $state<{ row: MediaFileRow; languageId: string } | undefined>();
	let previewRows = $state<MediaRenamePreviewRow[]>([]);
	let previewLoading = $state(false);
	let previewApplying = $state(false);
	let previewError = $state<string | undefined>();
	let renameApplyOpen = $state(false);
	const groups = $derived(mediaFileGroups(item, qualityProfiles));
	const safeRenameCount = $derived(previewRows.filter((row) => row.status === 'safe').length);
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
			languageId: languageId ?? row.expectedSubtitleLanguages[0] ?? 'english'
		};
	}

	async function loadRenamePreview() {
		previewLoading = true;
		previewError = undefined;
		try {
			const preview = await previewMediaRename(item.id);
			previewRows = preview.rows;
		} catch (error) {
			previewError = error instanceof Error ? error.message : 'Could not preview rename';
		} finally {
			previewLoading = false;
		}
	}

	async function confirmRenameApply() {
		renameApplyOpen = false;
		previewApplying = true;
		previewError = undefined;
		try {
			const result = await applyMediaRename(item.id);
			previewRows = result.rows;
		} catch (error) {
			previewError = error instanceof Error ? error.message : 'Could not apply rename';
		} finally {
			previewApplying = false;
		}
	}
</script>

<section aria-labelledby="media-files-title">
	<h2 id="media-files-title" class="m-0 text-3xl font-semibold text-foreground">Files</h2>
	<div class="grid gap-3.5">
		<MediaRootPanel {item} {libraryFolders} {canManage} {onSaveOptions} />
		{#if canManage && item.filePaths.length > 0}
			<MediaRenamePreviewPanel
				rows={previewRows}
				loading={previewLoading}
				applying={previewApplying}
				errorMessage={previewError}
				onPreview={loadRenamePreview}
				onApply={() => (renameApplyOpen = true)}
			/>
		{/if}
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

{#if renameApplyOpen}
	<MediaRenameApplyModal
		safeCount={safeRenameCount}
		onCancel={() => (renameApplyOpen = false)}
		onConfirm={confirmRenameApply}
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
