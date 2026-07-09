<script lang="ts">
	import MediaFileSummary from '$lib/components/app/media/files/MediaFileSummary.svelte';
	import MediaFileDeleteModal from '$lib/components/app/media/files/MediaFileDeleteModal.svelte';
	import MediaFileSearchModal from '$lib/components/app/media/files/MediaFileSearchModal.svelte';
	import SubtitleSearchModal from '$lib/components/app/media/subtitle-search/SubtitleSearchModal.svelte';
	import MediaEpisodeRow from '$lib/components/app/media/series/MediaEpisodeRow.svelte';
	import MediaRootPanel from '$lib/components/app/media/collection/MediaRootPanel.svelte';
	import MediaSeasonActions from '$lib/components/app/media/series/MediaSeasonActions.svelte';
	import MediaSeasonPanel from '$lib/components/app/media/series/MediaSeasonPanel.svelte';
	import MediaSeriesMonitorBookmark from '$lib/components/app/media/series/MediaSeriesMonitorBookmark.svelte';
	import {
		monitorUpdate,
		toggledEpisodeMonitor,
		toggledSeasonMonitor
	} from '$lib/components/app/media/series/mediaMonitoring';
	import { fileRow, type MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
	import { seasonFileSummary } from '$lib/components/app/media/series/mediaSeasonSummary';
	import {
		episodeReleaseSearchContext,
		episodeTitle,
		seasonReleaseSearchContext,
		seasonEpisodeRows,
		seasonMonitored
	} from '$lib/components/app/media/series/mediaSeriesRows';
	import type { ReleaseSearchContext } from '$lib/components/app/media/release-search/releaseSearchQuery';
	import type { MediaSeriesSeasonsProps as Props } from './mediaSeriesSeasonsTypes';
	import type { MediaMetadataEpisode, MediaMetadataSeason } from '$lib/settings/types';

	let {
		item,
		activities,
		searchingItemId,
		grabbingKey,
		canManage,
		pendingFulfillmentActionKeys = [],
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
		onFulfillmentAction = () => {},
		onGrabRelease
	}: Props = $props();

	let deleteRow = $state<MediaFileRow | undefined>();
	let searchContext = $state<ReleaseSearchContext | undefined>();
	let subtitleSearch = $state<{ row: MediaFileRow; languageId: string } | undefined>();
	const seasons = $derived(item.seasons ?? []);
	const mediaRows = $derived(item.filePaths.map((path) => fileRow(item, path, qualityProfiles)));

	function saveSeasonMonitor(season: MediaMetadataSeason) {
		onSaveOptions(item, monitorUpdate(toggledSeasonMonitor(item, seasons, season)));
	}

	function saveEpisodeMonitor(season: MediaMetadataSeason, episode: MediaMetadataEpisode) {
		onSaveOptions(item, monitorUpdate(toggledEpisodeMonitor(item, seasons, season, episode)));
	}

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

	function openSubtitleSearch(row: MediaFileRow, languageId?: string) {
		subtitleSearch = {
			row,
			languageId: languageId ?? row.subtitleSatisfaction?.wantedLanguages[0] ?? 'english'
		};
	}
</script>

{#if seasons.length > 0}
	<section aria-labelledby="metadata-seasons-title">
		<h2 id="metadata-seasons-title" class="m-0 text-3xl font-semibold text-foreground">Files</h2>
		<div class="grid gap-2.5">
			<MediaRootPanel {item} {libraryFolders} {canManage} {onSaveOptions} />
			<h3 class="m-0 text-xl font-semibold text-foreground">Seasons</h3>
			{#each seasons as season, index (season.name)}
				{@const seasonRows = seasonEpisodeRows(season, index, mediaRows, activities, item.id)}
				{@const summary = seasonFileSummary(seasonRows)}
				<MediaSeasonPanel summary={summary.label} size={summary.size} tone={summary.tone}>
					{#snippet title()}
						<span class="inline-flex min-w-0 items-center gap-2.5">
							<MediaSeriesMonitorBookmark
								name={season.name}
								monitored={seasonMonitored(season)}
								target="season"
								disabled={!canManage}
								onToggle={() => saveSeasonMonitor(season)}
							/>
							<span>{season.name}</span>
						</span>
					{/snippet}
					{#snippet actions()}
						<MediaSeasonActions
							{canManage}
							busy={searchingItemId === item.id || summary.hasActive}
							onAutoSearch={() => onAutoSearch(item)}
							onManualSearch={() => (searchContext = seasonReleaseSearchContext(season, index))}
						/>
					{/snippet}
					{#if seasonRows.length > 0}
						<div class="grid px-4.5">
							{#each seasonRows as file (file.episode.episodeNumber)}
								<MediaEpisodeRow episode={file.episode} title={episodeTitle(file.episode)}>
									{#snippet beforeTitle()}
										<MediaSeriesMonitorBookmark
											name={file.episode.name}
											monitored={file.episode.monitored === true}
											target="episode"
											disabled={!canManage}
											onToggle={() => saveEpisodeMonitor(season, file.episode)}
										/>
									{/snippet}
									<MediaFileSummary
										mediaItemId={item.id}
										mediaTitle={item.title}
										row={file.row}
										activityStatus={file.activityStatus}
										{canManage}
										{pendingFulfillmentActionKeys}
										searching={searchingItemId === item.id}
										onAutoSearch={() => onAutoSearch(item)}
										onManualSearch={() => (searchContext = episodeReleaseSearchContext(file.row))}
										onSearchSubtitle={searchSubtitle}
										onManualSubtitleSearch={openSubtitleSearch}
										onDeleteSubtitle={(subtitle) => deleteSubtitle(subtitle.id)}
										onUpdateSubtitle={(subtitle, request) =>
											onUpdateSubtitle(item, subtitle.id, request)}
										onDeleteTrack={(row, request) =>
											onDeleteFileTrack(item, { ...request, path: row.path ?? '' })}
										onFulfillmentAction={(row, request) =>
											onFulfillmentAction(item, {
												...request,
												filePath: row.path ?? request.filePath
											})}
										onDelete={requestDelete}
									/>
								</MediaEpisodeRow>
							{/each}
						</div>
					{:else}
						<p class="p-4.5 text-sm text-muted-foreground">No episode details available.</p>
					{/if}
				</MediaSeasonPanel>
			{/each}
		</div>
	</section>
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

{#if deleteRow}
	<MediaFileDeleteModal
		row={deleteRow}
		onCancel={() => (deleteRow = undefined)}
		onConfirm={confirmDelete}
	/>
{/if}

{#if searchContext}
	<MediaFileSearchModal
		{item}
		{languages}
		{searchContext}
		{grabbingKey}
		{canManage}
		onGrab={onGrabRelease}
		onClose={() => (searchContext = undefined)}
	/>
{/if}
