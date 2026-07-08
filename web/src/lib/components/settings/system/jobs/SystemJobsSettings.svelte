<script lang="ts">
	import { onMount } from 'svelte';
	import SystemJobAbortDialog from './SystemJobAbortDialog.svelte';
	import SystemJobHistorySection from './SystemJobHistorySection.svelte';
	import SystemJobLogsDialog from './SystemJobLogsDialog.svelte';
	import SystemJobsLiveSections from './SystemJobsLiveSections.svelte';
	import { SystemJobsController } from './systemJobsController.svelte';

	const controller = new SystemJobsController();
	const reloadDelayMs = 180;
	let mounted = false;
	const filterKey = $derived(
		JSON.stringify({
			statuses: controller.selectedStatuses,
			queues: controller.selectedQueues,
			kinds: controller.selectedKinds,
			query: controller.query,
			includeRoutine: controller.includeRoutine
		})
	);

	onMount(() => {
		mounted = true;
		return controller.start();
	});

	$effect(() => {
		filterKey;
		if (!mounted) return;
		const timeout = window.setTimeout(() => void controller.loadHistory(true), reloadDelayMs);
		return () => window.clearTimeout(timeout);
	});
</script>

<div class="grid gap-4">
	{#if controller.errorMessage}
		<p
			class="m-0 rounded-md border border-destructive/40 bg-destructive/10 p-3 text-sm text-destructive"
		>
			{controller.errorMessage}
		</p>
	{/if}

	<SystemJobsLiveSections
		schedules={controller.schedules}
		oneShotJobs={controller.oneShotJobs}
		loadingOverview={controller.loadingOverview}
		updatingScheduleId={controller.updatingScheduleId}
		updatingIntervalId={controller.updatingIntervalId}
		runningScheduleId={controller.runningScheduleId}
		abortingId={controller.abortingId}
		onRefresh={() => void controller.loadOverview()}
		onPause={(schedule) => void controller.toggleSchedule(schedule, true)}
		onResume={(schedule) => void controller.toggleSchedule(schedule, false)}
		onRun={(schedule) => void controller.runSchedule(schedule)}
		onSaveInterval={(schedule, intervalSeconds) =>
			void controller.saveScheduleInterval(schedule, intervalSeconds)}
		onAbort={(id, kind) => (controller.abortCandidate = { id, kind })}
	/>

	<SystemJobHistorySection
		executions={controller.visibleHistory}
		loading={controller.loadingHistory}
		hasMore={controller.historyHasMore}
		loadingMore={controller.loadingMore}
		savingRetention={controller.savingRetention}
		bind:retentionDays={controller.retentionDays}
		bind:routineRetentionHours={controller.routineRetentionHours}
		bind:includeRoutine={controller.includeRoutine}
		bind:selectedStatuses={controller.selectedStatuses}
		bind:selectedQueues={controller.selectedQueues}
		bind:selectedKinds={controller.selectedKinds}
		bind:query={controller.query}
		queueOptions={controller.queueOptions}
		kindOptions={controller.kindOptions}
		abortingId={controller.abortingId}
		loadingLogsId={controller.loadingLogsId}
		onSaveRetention={() => void controller.saveRetention()}
		onAbort={(id) =>
			(controller.abortCandidate = {
				id,
				kind: controller.history.find((job) => job.riverJobId === id)?.kind ?? 'job'
			})}
		onLogs={(execution) => void controller.openLogs(execution)}
		onLoadMore={() => void controller.loadHistory(false)}
	/>
</div>

<SystemJobAbortDialog
	job={controller.abortCandidate}
	onClose={() => (controller.abortCandidate = undefined)}
	onAbort={() => controller.abortJob()}
/>
<SystemJobLogsDialog
	execution={controller.logsExecution}
	logs={controller.executionLogs}
	loading={!!controller.loadingLogsId}
	onClose={() => (controller.logsExecution = undefined)}
/>
