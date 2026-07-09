<script lang="ts">
	import LivePulseDot from '$lib/components/shared/LivePulseDot.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import type { SystemJobExecution, SystemJobSchedule } from '$lib/settings/types';
	import SystemOneShotJobsTable from './SystemOneShotJobsTable.svelte';
	import SystemScheduledJobsTable from './SystemScheduledJobsTable.svelte';

	interface Props {
		schedules: SystemJobSchedule[];
		oneShotJobs: SystemJobExecution[];
		loadingOverview: boolean;
		updatingScheduleId?: string;
		updatingIntervalId?: string;
		runningScheduleId?: string;
		abortingId?: number;
		loadingLogsId?: number;
		onRefresh: () => void;
		onPause: (schedule: SystemJobSchedule) => void;
		onResume: (schedule: SystemJobSchedule) => void;
		onRun: (schedule: SystemJobSchedule) => void;
		onSaveInterval: (schedule: SystemJobSchedule, intervalSeconds: number) => void;
		onAbort: (id: number, kind: string) => void;
		onLogs: (execution: SystemJobExecution) => void;
	}

	let {
		schedules,
		oneShotJobs,
		loadingOverview,
		updatingScheduleId,
		updatingIntervalId,
		runningScheduleId,
		abortingId,
		loadingLogsId,
		onRefresh,
		onPause,
		onResume,
		onRun,
		onSaveInterval,
		onAbort,
		onLogs
	}: Props = $props();
</script>

<Card.Root>
	<Card.Header>
		<Card.Description class="flex items-center gap-2">
			<LivePulseDot /><span>Live</span>
		</Card.Description>
		<Card.Title>Fixed Scheduled Jobs</Card.Title>
		<Card.Action>
			<Button variant="outline" disabled={loadingOverview} onclick={onRefresh}>Refresh</Button>
		</Card.Action>
	</Card.Header>
	<Card.Content>
		<SystemScheduledJobsTable
			{schedules}
			updatingId={updatingScheduleId}
			{updatingIntervalId}
			runningId={runningScheduleId}
			{abortingId}
			{onPause}
			{onResume}
			{onRun}
			{onSaveInterval}
			onAbort={(id) => onAbort(id, 'fixed scheduled job')}
		/>
	</Card.Content>
</Card.Root>

<Card.Root>
	<Card.Header>
		<Card.Description class="flex items-center gap-2">
			<LivePulseDot /><span>Live</span>
		</Card.Description>
		<Card.Title>Current And Planned One-Shot Jobs</Card.Title>
	</Card.Header>
	<Card.Content>
		<SystemOneShotJobsTable
			jobs={oneShotJobs}
			{abortingId}
			{loadingLogsId}
			onAbort={(id) => onAbort(id, oneShotJobs.find((job) => job.riverJobId === id)?.kind ?? 'job')}
			{onLogs}
		/>
	</Card.Content>
</Card.Root>
