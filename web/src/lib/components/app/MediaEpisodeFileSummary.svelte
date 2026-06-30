<script lang="ts">
	import type { MediaFileRow } from './mediaFiles';
	import type { ActivityQueueStatus } from './activityQueue';

	interface Props {
		row: MediaFileRow;
		activityStatus?: ActivityQueueStatus;
		canManage: boolean;
		searching: boolean;
		fileLabel?: string;
		missingLabel?: string;
		onInfo: (_row: MediaFileRow) => void;
		onAutoSearch: () => void;
		onManualSearch: () => void;
		onDelete: (_row: MediaFileRow) => void;
	}

	let {
		row,
		activityStatus,
		canManage,
		searching,
		fileLabel = 'Episode file',
		missingLabel = 'No matched file for this episode',
		onInfo,
		onAutoSearch,
		onManualSearch,
		onDelete
	}: Props = $props();
	const busy = $derived(
		searching ||
			activityStatus?.status === 'queued' ||
			activityStatus?.status === 'grabbed' ||
			activityStatus?.status === 'downloading'
	);
</script>

<div class:missing-file={!row.exists} class="episode-file-summary">
	<div class="episode-file-path">
		<strong>{row.exists ? row.relativePath : 'Missing file'}</strong>
		<span>{row.exists ? fileLabel : missingLabel}</span>
	</div>

	<div class="episode-file-facts" aria-label="Episode file details">
		<span>
			<strong>Quality</strong>
			{row.quality}
		</span>
		<span>
			<strong>Video</strong>
			{row.videoCodec}
		</span>
		<span>
			<strong>Audio</strong>
			{row.audioInfo}
		</span>
		<span>
			<strong>Languages</strong>
			{row.languages}
		</span>
		<span>
			<strong>Score</strong>
			{row.score}
		</span>
		<span>
			<strong>Status</strong>
			{#if activityStatus}
				<small
					class="activity-status-chip"
					class:activity-failed={activityStatus.status === 'failed'}
				>
					<span class="app-icon" aria-hidden="true">sync</span>
					{activityStatus.label}
				</small>
			{:else}
				-
			{/if}
		</span>
	</div>

	{#if row.formats.length > 0}
		<div class="format-chip-list episode-file-formats" aria-label="Matched formats">
			{#each row.formats as format (format)}
				<span>{format}</span>
			{/each}
		</div>
	{/if}

	<div class="row-actions episode-file-actions">
		{#if row.exists}
			<button
				type="button"
				class="secondary icon-button"
				aria-label="File info"
				title="File info"
				onclick={() => onInfo(row)}
			>
				<span class="app-icon" aria-hidden="true">info</span>
			</button>
		{/if}
		<button
			type="button"
			class="secondary icon-button"
			aria-label="Automatic search"
			title="Automatic search"
			disabled={!canManage || busy}
			onclick={onAutoSearch}
		>
			<span class="app-icon" aria-hidden="true">search</span>
		</button>
		<button
			type="button"
			class="secondary icon-button"
			aria-label="Manual search"
			title="Manual search"
			disabled={busy}
			onclick={onManualSearch}
		>
			<span class="app-icon" aria-hidden="true">person</span>
		</button>
		{#if row.exists}
			<button
				type="button"
				class="danger icon-button"
				aria-label="Delete file"
				title="Delete file"
				disabled={!canManage || !row.path}
				onclick={() => onDelete(row)}
			>
				<span class="app-icon" aria-hidden="true">delete</span>
			</button>
		{/if}
	</div>
</div>
