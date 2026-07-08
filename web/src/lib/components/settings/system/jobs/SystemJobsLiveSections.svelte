<script lang="ts">
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
		abortingId?: number;
		onRefresh: () => void;
		onPause: (schedule: SystemJobSchedule) => void;
		onResume: (schedule: SystemJobSchedule) => void;
		onAbort: (id: number, kind: string) => void;
	}

	let {
		schedules,
		oneShotJobs,
		loadingOverview,
		updatingScheduleId,
		abortingId,
		onRefresh,
		onPause,
		onResume,
		onAbort
	}: Props = $props();
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Fixed Scheduled Jobs</Card.Title>
		<Card.Action>
			<Button variant="outline" disabled={loadingOverview} onclick={onRefresh}>Refresh</Button>
		</Card.Action>
	</Card.Header>
	<Card.Content>
		<SystemScheduledJobsTable
			{schedules}
			updatingId={updatingScheduleId}
			{abortingId}
			{onPause}
			{onResume}
			onAbort={(id) => onAbort(id, 'fixed scheduled job')}
		/>
	</Card.Content>
</Card.Root>

<Card.Root>
	<Card.Header><Card.Title>Current And Planned One-Shot Jobs</Card.Title></Card.Header>
	<Card.Content>
		<SystemOneShotJobsTable
			jobs={oneShotJobs}
			{abortingId}
			onAbort={(id) => onAbort(id, oneShotJobs.find((job) => job.riverJobId === id)?.kind ?? 'job')}
		/>
	</Card.Content>
</Card.Root>
