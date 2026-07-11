import { client } from '$lib/api/client';
import type {
	SystemJobExecutionListResponse,
	SystemJobExecutionLog,
	SystemJobHistorySettings,
	SystemJobSchedule,
	SystemJobsOverviewResponse
} from '$lib/settings/types';

export interface SystemJobFilters {
	status?: string[];
	queue?: string;
	kind?: string;
	query?: string;
	limit?: number;
}

export interface SystemJobExecutionFilters extends SystemJobFilters {
	scheduleId?: string;
	before?: string;
	includeRoutine?: boolean;
}

export async function listSystemJobs(filters: SystemJobFilters = {}) {
	const { data, error } = await client.GET('/system/jobs', { params: { query: filters } });
	if (error) throw new Error(error.message);
	return data?.jobs ?? [];
}

export async function abortSystemJob(id: number) {
	const { data, error } = await client.POST('/system/jobs/{id}/abort', {
		params: { path: { id } }
	});
	if (error) throw new Error(error.message);
	if (!data) throw new Error('Job abort did not return a job');
	return data;
}

export async function getSystemJobsOverview(): Promise<SystemJobsOverviewResponse> {
	const { data, error } = await client.GET('/system/jobs/overview');
	if (error) throw new Error(error.message);
	if (!data) throw new Error('Jobs overview request did not return a result');
	return data;
}

async function updateSchedule(
	path: 'pause' | 'resume' | 'run',
	id: string,
	emptyMessage: string
): Promise<SystemJobSchedule> {
	const { data, error } = await client.POST(`/system/job-schedules/{id}/${path}`, {
		params: { path: { id } }
	});
	if (error) throw new Error(error.message);
	if (!data) throw new Error(emptyMessage);
	return data;
}

export const pauseSystemJobSchedule = (id: string) =>
	updateSchedule('pause', id, 'Schedule pause did not return a schedule');
export const resumeSystemJobSchedule = (id: string) =>
	updateSchedule('resume', id, 'Schedule resume did not return a schedule');
export const runSystemJobSchedule = (id: string) =>
	updateSchedule('run', id, 'Schedule run did not return a schedule');

export async function updateSystemJobScheduleInterval(id: string, intervalSeconds: number) {
	const { data, error } = await client.PUT('/system/job-schedules/{id}/interval', {
		params: { path: { id } },
		body: { intervalSeconds }
	});
	if (error) throw new Error(error.message);
	if (!data) throw new Error('Schedule interval update did not return a schedule');
	return data;
}

export async function listSystemJobExecutions(
	filters: SystemJobExecutionFilters = {}
): Promise<SystemJobExecutionListResponse> {
	const { data, error } = await client.GET('/system/job-executions', {
		params: { query: filters }
	});
	if (error) throw new Error(error.message);
	return data ?? { executions: [], hasMore: false };
}

export async function listSystemJobExecutionLogs(riverJobId: number) {
	const { data, error } = await client.GET('/system/job-executions/{riverJobId}/logs', {
		params: { path: { riverJobId } }
	});
	if (error) throw new Error(error.message);
	return (data?.logs ?? []) as SystemJobExecutionLog[];
}

export async function updateSystemJobHistorySettings(request: SystemJobHistorySettings) {
	const { data, error } = await client.PUT('/system/job-history-settings', { body: request });
	if (error) throw new Error(error.message);
	if (!data) throw new Error('Job history settings update did not return a result');
	return data;
}
