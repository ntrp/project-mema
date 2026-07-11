import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import * as api from '$lib/components/settings/system/jobs/api';
import type { SystemJobExecutionFilters } from '$lib/components/settings/system/jobs/api';

export function createSystemJobsOperations() {
	const client = useQueryClient();
	const abort = createMutation(() => ({ mutationFn: api.abortSystemJob }));
	const pause = createMutation(() => ({ mutationFn: api.pauseSystemJobSchedule }));
	const resume = createMutation(() => ({ mutationFn: api.resumeSystemJobSchedule }));
	const run = createMutation(() => ({ mutationFn: api.runSystemJobSchedule }));
	const interval = createMutation(() => ({
		mutationFn: ({ id, seconds }: { id: string; seconds: number }) =>
			api.updateSystemJobScheduleInterval(id, seconds)
	}));
	const retention = createMutation(() => ({ mutationFn: api.updateSystemJobHistorySettings }));
	return {
		overview: () =>
			client.fetchQuery({
				queryKey: ['system', 'jobs', 'overview'],
				queryFn: api.getSystemJobsOverview
			}),
		history: (filters: SystemJobExecutionFilters) =>
			client.fetchQuery({
				queryKey: ['system', 'jobs', 'history', filters],
				queryFn: () => api.listSystemJobExecutions(filters)
			}),
		logs: (id: number) =>
			client.fetchQuery({
				queryKey: ['system', 'jobs', 'logs', id],
				queryFn: () => api.listSystemJobExecutionLogs(id)
			}),
		abort: (id: number) => abort.mutateAsync(id),
		pause: (id: string) => pause.mutateAsync(id),
		resume: (id: string) => resume.mutateAsync(id),
		run: (id: string) => run.mutateAsync(id),
		interval: (id: string, seconds: number) => interval.mutateAsync({ id, seconds }),
		retention: retention.mutateAsync
	};
}

export type SystemJobsOperations = ReturnType<typeof createSystemJobsOperations>;
