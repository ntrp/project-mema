<script lang="ts">
	import type { DownloadActivity } from '$lib/settings/types';

	interface Props {
		activities: DownloadActivity[];
		loading: boolean;
		canManage: boolean;
		cancellingId?: string;
		onRefresh: () => void;
		onCancel: (_activity: DownloadActivity) => void;
	}

	let { activities, loading, canManage, cancellingId, onRefresh, onCancel }: Props = $props();

	function createdLabel(value: string) {
		return new Intl.DateTimeFormat(undefined, {
			dateStyle: 'short',
			timeStyle: 'short'
		}).format(new Date(value));
	}

	function cancellable(activity: DownloadActivity) {
		return ['queued', 'grabbed', 'downloading'].includes(activity.status);
	}

	function showsProgress(activity: DownloadActivity) {
		return ['queued', 'grabbed', 'downloading', 'completed'].includes(activity.status);
	}

	function progressValue(activity: DownloadActivity) {
		if (activity.status === 'completed') {
			return 100;
		}
		return activity.progressPercent ?? undefined;
	}

	function progressLabel(activity: DownloadActivity) {
		const value = progressValue(activity);
		return typeof value === 'number' ? `${value}%` : 'Waiting for client progress';
	}
</script>

<div class="page-heading split-heading">
	<div>
		<p>Activity</p>
		<h1 id="home-title">Downloads and imports</h1>
	</div>
	<button type="button" class="secondary" disabled={loading} onclick={onRefresh}>
		{loading ? 'Refreshing' : 'Refresh'}
	</button>
</div>

<div class="data-list">
	{#each activities as activity (activity.id)}
		<div class="data-row activity-row">
			<div>
				<strong>{activity.releaseTitle}</strong>
				<span>{activity.mediaTitle} · {activity.mediaType}</span>
				{#if showsProgress(activity)}
					<div
						class="download-progress"
						role="progressbar"
						aria-label="Download progress"
						aria-valuemin="0"
						aria-valuemax="100"
						aria-valuenow={progressValue(activity)}
					>
						<span
							class:indeterminate={progressValue(activity) === undefined}
							style={progressValue(activity) !== undefined
								? `width: ${progressValue(activity)}%`
								: undefined}
						></span>
					</div>
					<small class="progress-label">{progressLabel(activity)}</small>
				{/if}
			</div>
			<span>{activity.downloadClientName} · {activity.indexerName}</span>
			<div class="status-stack">
				<small
					class:status-enabled={activity.status === 'grabbed' || activity.status === 'completed'}
					class:pending={activity.status === 'queued' || activity.status === 'downloading'}
					class:test-failed={activity.status === 'failed' || activity.status === 'cancelled'}
				>
					{activity.status}
				</small>
				<small>{createdLabel(activity.createdAt)}</small>
				{#if activity.error}
					<small class="test-detail">{activity.error}</small>
				{/if}
			</div>
			{#if canManage && cancellable(activity)}
				<button
					type="button"
					class="danger"
					disabled={cancellingId === activity.id}
					onclick={() => onCancel(activity)}
				>
					{cancellingId === activity.id ? 'Cancelling' : 'Cancel'}
				</button>
			{/if}
		</div>
	{:else}
		<div class="panel">
			<p class="empty">No download activity yet</p>
		</div>
	{/each}
</div>
