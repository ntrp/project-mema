import type { SystemJobExecution, SystemJobSchedule } from '$lib/settings/types';
import { activeStatuses, finalStatuses } from './systemJobDisplay';

export function upsertExecution(list: SystemJobExecution[], execution: SystemJobExecution) {
	return mergeExecutions(list, [execution]);
}

export function mergeExecutions(list: SystemJobExecution[], executions: SystemJobExecution[]) {
	const byID = new Map(list.map((execution) => [execution.riverJobId, execution]));
	for (const execution of executions) byID.set(execution.riverJobId, execution);
	return sortExecutions(Array.from(byID.values()));
}

export function sortExecutions(executions: SystemJobExecution[]) {
	return [...executions].sort((left, right) => {
		const updated = Date.parse(right.updatedAt) - Date.parse(left.updatedAt);
		return updated === 0 ? right.riverJobId - left.riverJobId : updated;
	});
}

export function defaultHistoryIncludesExecution(execution: SystemJobExecution) {
	if (!finalStatuses.includes(execution.status)) return false;
	if (execution.historyPolicy !== 'routine') return true;
	return ['cancelled', 'discarded'].includes(execution.status);
}

export function updateOneShotJobs(list: SystemJobExecution[], execution: SystemJobExecution) {
	if (activeStatuses.includes(execution.status)) {
		return upsertExecution(list, execution);
	}
	return list.filter((job) => job.riverJobId !== execution.riverJobId);
}

export function updateScheduleFromExecution(
	schedules: SystemJobSchedule[],
	execution: SystemJobExecution
) {
	return schedules.map((schedule) => {
		if (schedule.id !== execution.scheduleId) return schedule;
		const next = { ...schedule };
		if (activeStatuses.includes(execution.status)) {
			next.activeRiverJobId = execution.riverJobId;
			next.activeStatus = execution.status;
			next.activeProgressPercent = execution.progressPercent;
			next.activeProgressLabel = execution.progressLabel;
			next.activeInfoMessage = execution.infoMessage;
		} else if (next.activeRiverJobId === execution.riverJobId) {
			next.activeRiverJobId = undefined;
			next.activeStatus = '';
			next.activeProgressPercent = undefined;
			next.activeProgressLabel = '';
			next.activeInfoMessage = '';
		}
		next.lastRiverJobId = execution.riverJobId;
		next.lastStatus = execution.status;
		next.lastCreatedAt = execution.createdAt;
		next.lastFinalizedAt = execution.finalizedAt;
		next.nextRunAt = next.paused ? undefined : nextRun(execution.createdAt, next.intervalSeconds);
		return next;
	});
}

export function matchesExecutionFilters(
	execution: SystemJobExecution,
	queues: string[],
	kinds: string[]
) {
	if (queues.length > 0 && !queues.includes(execution.queue)) return false;
	if (kinds.length > 0 && !kinds.includes(execution.kind)) return false;
	return true;
}

export function optionList(values: string[]) {
	return Array.from(new Set(values.filter(Boolean)))
		.sort()
		.map((value) => ({ value, label: value }));
}

function nextRun(createdAt: string, intervalSeconds: number) {
	return new Date(new Date(createdAt).getTime() + intervalSeconds * 1000).toISOString();
}
