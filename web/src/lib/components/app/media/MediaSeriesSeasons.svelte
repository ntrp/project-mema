<script lang="ts">
	import MediaFileSummary from './MediaFileSummary.svelte';
	import MediaFileDeleteModal from './MediaFileDeleteModal.svelte';
	import MediaFileSearchModal from './MediaFileSearchModal.svelte';
	import MediaEpisodeRow from './MediaEpisodeRow.svelte';
	import MediaRootPanel from './MediaRootPanel.svelte';
	import MediaSeasonActions from './MediaSeasonActions.svelte';
	import MediaSeasonPanel from './MediaSeasonPanel.svelte';
	import MediaSeriesMonitorBookmark from './MediaSeriesMonitorBookmark.svelte';
	import { monitorUpdate, toggledEpisodeMonitor, toggledSeasonMonitor } from './mediaMonitoring';
	import { fileRow, seasonNumberFromName, type MediaFileRow } from './mediaFiles';
	import { seasonFileSummary } from './mediaSeasonSummary';
	import { episodeTitle, seasonEpisodeRows, seasonMonitored } from './mediaSeriesRows';
	import type { ReleaseSearchContext } from './releaseSearchQuery';
	import type {
		DownloadActivity,
		Language,
		LibraryFolder,
		MediaItem,
		MediaItemUpdateRequest,
		MediaMetadataEpisode,
		MediaMetadataSeason,
		ReleaseCandidate
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
		onGrabRelease: (_item: MediaItem, _release: ReleaseCandidate) => void;
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
	let searchContext = $state<ReleaseSearchContext | undefined>();
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

	function seasonSearchContext(season: MediaMetadataSeason, index: number): ReleaseSearchContext {
		return { type: 'season', seasonNumber: seasonNumberFromName(season.name) ?? index + 1 };
	}

	function episodeSearchContext(row: MediaFileRow): ReleaseSearchContext {
		return { type: 'episode', seasonNumber: row.seasonNumber, episodeNumber: row.episodeNumber };
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
							onManualSearch={() => (searchContext = seasonSearchContext(season, index))}
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
										row={file.row}
										activityStatus={file.activityStatus}
										{canManage}
										searching={searchingItemId === item.id}
										onAutoSearch={() => onAutoSearch(item)}
										onManualSearch={() => (searchContext = episodeSearchContext(file.row))}
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
