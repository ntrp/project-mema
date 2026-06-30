<script lang="ts">
	import type { MediaItem, MediaItemStatus } from '$lib/settings/types';

	interface Props {
		item: MediaItem;
		releaseCount: number;
		qualityProfileLabel: string;
		libraryFolderLabel: string;
		monitorModeLabel: string;
	}

	let { item, releaseCount, qualityProfileLabel, libraryFolderLabel, monitorModeLabel }: Props =
		$props();

	function statusLabel(status: MediaItemStatus) {
		switch (status) {
			case 'downloaded':
				return 'Downloaded';
			case 'downloading':
				return 'Downloading';
			default:
				return 'Missing';
		}
	}
</script>

<section aria-labelledby="library-status-title">
	<h2 id="library-status-title">Library Status</h2>
	<div class="metadata-facts-grid" aria-label="Library status facts">
		<div>
			<strong>{statusLabel(item.status)}</strong>
			<span>Status</span>
		</div>
		<div>
			<strong>{releaseCount}</strong>
			<span>Release candidates</span>
		</div>
		<div>
			<strong>{qualityProfileLabel}</strong>
			<span>Profile</span>
		</div>
		<div>
			<strong>{libraryFolderLabel}</strong>
			<span>Media folder</span>
		</div>
		<div>
			<strong>{item.monitored ? 'Monitored' : 'Paused'}</strong>
			<span>Monitor state</span>
		</div>
		<div>
			<strong>{monitorModeLabel}</strong>
			<span>Monitor</span>
		</div>
	</div>
	{#if item.tags?.length}
		<div class="metadata-tags library-tags" aria-label="Library tags">
			{#each item.tags as tag (tag)}
				<span><span class="app-icon" aria-hidden="true">sell</span>{tag}</span>
			{/each}
		</div>
	{/if}
</section>
