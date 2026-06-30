<script lang="ts">
	import ActivityManualImportModal from './ActivityManualImportModal.svelte';
	import ActivityActions from './ActivityActions.svelte';
	import { activityDisplay, createdLabel } from './activityDisplay';
	import { manualImportDownloadActivity } from '$lib/settings/api';
	import type { DownloadActivity, ManualImportRequest } from '$lib/settings/types';

	interface Props {
		activities: DownloadActivity[];
		loading: boolean;
		canManage: boolean;
		cancellingId?: string;
		deletingId?: string;
		onRefresh: () => void;
		onCancel: (_activity: DownloadActivity) => void;
		onDelete: (_activity: DownloadActivity) => void;
	}

	let {
		activities,
		loading,
		canManage,
		cancellingId,
		deletingId,
		onRefresh,
		onCancel,
		onDelete
	}: Props = $props();
	let manualImportActivity = $state<DownloadActivity | undefined>();
	let importingId = $state<string | undefined>();
	let importError = $state<string | undefined>();

	function showsProgress(activity: DownloadActivity) {
		return ['queued', 'grabbed', 'downloading', 'completed'].includes(activity.status);
	}

	function openManualImport(activity: DownloadActivity) {
		manualImportActivity = activity;
		importError = undefined;
	}

	async function submitManualImport(request: ManualImportRequest) {
		if (!manualImportActivity) return;
		importingId = manualImportActivity.id;
		importError = undefined;
		try {
			await manualImportDownloadActivity(manualImportActivity.id, request);
			manualImportActivity = undefined;
			await onRefresh();
		} catch (error) {
			importError = error instanceof Error ? error.message : 'Manual import failed';
		} finally {
			importingId = undefined;
		}
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

{#if activities.length > 0}
	<div class="table-wrap activity-table">
		<table>
			<thead>
				<tr>
					<th><span class="sr-only">Select</span></th>
					<th>Media</th>
					<th>Year</th>
					<th>Languages</th>
					<th>Quality</th>
					<th>Formats</th>
					<th>Time left</th>
					<th>Progress</th>
					<th class="table-action-heading">Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each activities as activity (activity.id)}
					{@const display = activityDisplay(activity)}
					<tr>
						<td class="activity-check-cell">
							<input type="checkbox" aria-label={`Select ${activity.releaseTitle}`} />
						</td>
						<td class="activity-media-cell">
							<strong>{activity.mediaTitle}</strong>
							<small>{activity.downloadClientName} · {activity.indexerName}</small>
							{#if activity.error}
								<small class="test-detail">{activity.error}</small>
							{/if}
						</td>
						<td>{display.year}</td>
						<td>{display.languages.length ? display.languages.join(', ') : '-'}</td>
						<td>{display.quality}</td>
						<td>
							{#if display.formats.length}
								<div class="format-chip-list">
									{#each display.formats as format (format)}
										<span>{format}</span>
									{/each}
								</div>
							{:else}
								-
							{/if}
						</td>
						<td>{display.timeLeft}</td>
						<td>
							<div class="activity-progress-stack">
								<small
									class:status-enabled={activity.status === 'grabbed' ||
										activity.status === 'completed'}
									class:pending={activity.status === 'queued' || activity.status === 'downloading'}
									class:test-failed={activity.status === 'failed' ||
										activity.status === 'cancelled'}
								>
									{activity.status}
								</small>
								{#if showsProgress(activity)}
									<div
										class="download-progress"
										role="progressbar"
										aria-label="Download progress"
										aria-valuemin="0"
										aria-valuemax="100"
										aria-valuenow={display.progressValue}
									>
										<span
											class:indeterminate={display.progressValue === undefined}
											style={display.progressValue !== undefined
												? `width: ${display.progressValue}%`
												: undefined}
										></span>
									</div>
									<small>{display.progressLabel}</small>
								{/if}
								<small>{createdLabel(activity.createdAt)}</small>
							</div>
						</td>
						<td class="row-actions activity-actions">
							<ActivityActions
								{activity}
								{canManage}
								{cancellingId}
								{deletingId}
								onManualImport={openManualImport}
								{onCancel}
								{onDelete}
							/>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{:else}
	<div class="panel">
		<p class="empty">No download activity yet</p>
	</div>
{/if}

{#if manualImportActivity}
	<ActivityManualImportModal
		activity={manualImportActivity}
		importing={importingId === manualImportActivity.id}
		error={importError}
		onImport={submitManualImport}
		onClose={() => (manualImportActivity = undefined)}
	/>
{/if}
