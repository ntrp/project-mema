<script lang="ts">
	import MediaEpisodeFileSummary from './MediaEpisodeFileSummary.svelte';
	import MediaFileDeleteModal from './MediaFileDeleteModal.svelte';
	import MediaFileInfoModal from './MediaFileInfoModal.svelte';
	import MediaFileSearchModal from './MediaFileSearchModal.svelte';
	import MediaEpisodeRow from './MediaEpisodeRow.svelte';
	import MediaMonitorBookmark from './MediaMonitorBookmark.svelte';
	import MediaSeasonPanel from './MediaSeasonPanel.svelte';
	import { activityForEpisode } from '../activity/activityQueue';
	import {
		episodeKey,
		fileRow,
		missingRow,
		seasonNumberFromName,
		type MediaFileRow
	} from './mediaFiles';
	import type {
		DownloadActivity,
		MediaItem,
		MediaMetadataEpisode,
		MediaMetadataSeason,
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
	const seasons = $derived(item.seasons ?? []);
	const mediaRows = $derived(item.filePaths.map((path) => fileRow(item, path)));

	function episodeTitle(episode: MediaMetadataEpisode) {
		return `${episode.episodeNumber} - ${episode.name}`;
	}

	function episodeFileRow(
		season: MediaMetadataSeason,
		seasonIndex: number,
		episode: MediaMetadataEpisode
	) {
		const seasonNumber = seasonNumberFromName(season.name) ?? seasonIndex + 1;
		return (
			mediaRows.find(
				(row) =>
					episodeKey(row.seasonNumber, row.episodeNumber) ===
					episodeKey(seasonNumber, episode.episodeNumber)
			) ??
			missingRow(
				`s${seasonNumber}e${episode.episodeNumber}`,
				episode.name,
				seasonNumber,
				episode.episodeNumber
			)
		);
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
</script>

{#if seasons.length > 0}
	<section aria-labelledby="metadata-seasons-title">
		<h2 id="metadata-seasons-title" class="m-0 text-3xl font-semibold text-foreground">Seasons</h2>
		<div class="grid gap-2.5">
			{#each seasons as season, index (season.name)}
				<MediaSeasonPanel
					meta={season.episodeCount ? `${season.episodeCount} episodes` : 'Episodes unknown'}
				>
					{#snippet title()}
						<span class="inline-flex min-w-0 items-center gap-2.5">
							<MediaMonitorBookmark
								monitored={season.monitored === true}
								label={`${season.name} ${season.monitored ? 'monitored' : 'not monitored'}`}
							/>
							<span>{season.name}</span>
						</span>
					{/snippet}
					{#if season.episodes && season.episodes.length > 0}
						<div class="grid px-4.5">
							{#each season.episodes as episode (episode.episodeNumber)}
								{@const row = episodeFileRow(season, index, episode)}
								{@const activityStatus = activityForEpisode(
									activities,
									item.id,
									row.seasonNumber,
									row.episodeNumber
								)}
								<MediaEpisodeRow {episode} title={episodeTitle(episode)}>
									{#snippet beforeTitle()}
										<MediaMonitorBookmark
											monitored={episode.monitored === true}
											label={`${episode.name} ${episode.monitored ? 'monitored' : 'not monitored'}`}
										/>
									{/snippet}
									<MediaEpisodeFileSummary
										{row}
										{activityStatus}
										{canManage}
										searching={searchingItemId === item.id}
										onInfo={(nextRow) => (detailRow = nextRow)}
										onAutoSearch={() => onAutoSearch(item)}
										onManualSearch={() => (searchOpen = true)}
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
