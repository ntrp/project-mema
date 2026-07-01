<script lang="ts">
	import MediaEpisodeFileSummary from './MediaEpisodeFileSummary.svelte';
	import MediaFileInfoModal from './MediaFileInfoModal.svelte';
	import MediaFileSearchModal from './MediaFileSearchModal.svelte';
	import MediaMonitorBookmark from './MediaMonitorBookmark.svelte';
	import { activityForEpisode } from '../activity/activityQueue';
	import {
		episodeKey,
		fileRow,
		missingRow,
		seasonNumberFromName,
		type MediaFileRow
	} from './mediaFiles';
	import { imageUrl } from './mediaDetail';
	import { formatDate } from '$lib/settings/dateFormat';
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

	function confirmDelete(row: MediaFileRow) {
		if (!row.path) return;
		if (window.confirm(`Delete ${row.relativePath}?`)) {
			onDeleteFile(item, row.path);
		}
	}
</script>

{#if seasons.length > 0}
	<section aria-labelledby="metadata-seasons-title">
		<h2 id="metadata-seasons-title">Seasons</h2>
		<div class="metadata-season-list">
			{#each seasons as season, index (season.name)}
				<details class="metadata-season-panel">
					<summary>
						<span class="metadata-season-title">
							<MediaMonitorBookmark
								monitored={season.monitored === true}
								label={`${season.name} ${season.monitored ? 'monitored' : 'not monitored'}`}
							/>
							<span>{season.name}</span>
						</span>
						<small
							>{season.episodeCount ? `${season.episodeCount} episodes` : 'Episodes unknown'}</small
						>
					</summary>
					{#if season.episodes && season.episodes.length > 0}
						<div class="metadata-episode-list">
							{#each season.episodes as episode (episode.episodeNumber)}
								{@const row = episodeFileRow(season, index, episode)}
								{@const activityStatus = activityForEpisode(
									activities,
									item.id,
									row.seasonNumber,
									row.episodeNumber
								)}
								<article class="metadata-episode-row">
									<div class="metadata-episode-copy">
										<h3>
											<MediaMonitorBookmark
												monitored={episode.monitored === true}
												label={`${episode.name} ${episode.monitored ? 'monitored' : 'not monitored'}`}
											/>
											<span class="metadata-episode-title">{episodeTitle(episode)}</span>
											{#if episode.airDate}
												<span class="metadata-episode-date">{formatDate(episode.airDate)}</span>
											{/if}
										</h3>
										<p>{episode.overview ?? 'No episode overview available.'}</p>
									</div>
									{#if imageUrl(episode.stillPath, 'w300')}
										<img src={imageUrl(episode.stillPath, 'w300')} alt="" loading="lazy" />
									{/if}
									<MediaEpisodeFileSummary
										{row}
										{activityStatus}
										{canManage}
										searching={searchingItemId === item.id}
										onInfo={(nextRow) => (detailRow = nextRow)}
										onAutoSearch={() => onAutoSearch(item)}
										onManualSearch={() => (searchOpen = true)}
										onDelete={confirmDelete}
									/>
								</article>
							{/each}
						</div>
					{:else}
						<p class="metadata-season-empty">No episode details available.</p>
					{/if}
				</details>
			{/each}
		</div>
	</section>
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
