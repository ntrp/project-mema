<script lang="ts">
	import type { DownloadActivity } from '$lib/settings/types';

	interface Props {
		activities: DownloadActivity[];
		loading: boolean;
		onRefresh: () => void;
	}

	let { activities, loading, onRefresh }: Props = $props();

	function createdLabel(value: string) {
		return new Intl.DateTimeFormat(undefined, {
			dateStyle: 'short',
			timeStyle: 'short'
		}).format(new Date(value));
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
			</div>
			<span>{activity.downloadClientName} · {activity.indexerName}</span>
			<div class="status-stack">
				<small
					class:status-enabled={activity.status === 'grabbed'}
					class:pending={activity.status === 'queued'}
					class:test-failed={activity.status === 'failed'}
				>
					{activity.status}
				</small>
				<small>{createdLabel(activity.createdAt)}</small>
				{#if activity.error}
					<small class="test-detail">{activity.error}</small>
				{/if}
			</div>
		</div>
	{:else}
		<div class="panel">
			<p class="empty">No download activity yet</p>
		</div>
	{/each}
</div>
