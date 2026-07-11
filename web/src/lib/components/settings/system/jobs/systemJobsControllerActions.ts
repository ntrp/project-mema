import {
	abortSystemJob,
	getSystemJobsOverview,
	listSystemJobExecutionLogs,
	listSystemJobExecutions,
	pauseSystemJobSchedule,
	resumeSystemJobSchedule,
	runSystemJobSchedule,
	updateSystemJobHistorySettings,
	updateSystemJobScheduleInterval
} from './api';
import type {
	SystemJobExecution,
	SystemJobExecutionLog,
	SystemJobSchedule
} from '$lib/settings/types';
import { mergeExecutions, sortExecutions } from './systemJobsState';

const historyPageLimit = 20;

interface ControllerState {
	schedules: SystemJobSchedule[];
	oneShotJobs: SystemJobExecution[];
	history: SystemJobExecution[];
	historyHasMore: boolean;
	retentionDays: number;
	routineRetentionHours: number;
	includeRoutine: boolean;
	selectedStatuses: string[];
	query: string;
	loadingOverview: boolean;
	loadingHistory: boolean;
	loadingMore: boolean;
	savingRetention: boolean;
	errorMessage: string;
	updatingScheduleId?: string;
	updatingIntervalId?: string;
	runningScheduleId?: string;
	abortingId?: number;
	abortCandidate?: { id: number; kind: string };
	logsExecution?: SystemJobExecution;
	executionLogs: SystemJobExecutionLog[];
	loadingLogsId?: number;
}

export async function loadOverviewState(state: ControllerState) {
	state.loadingOverview = true;
	state.errorMessage = '';
	try {
		const overview = await getSystemJobsOverview();
		state.schedules = overview.schedules;
		state.oneShotJobs = overview.oneShotJobs;
		state.retentionDays = overview.historySettings.retentionDays;
		state.routineRetentionHours = overview.historySettings.routineRetentionHours;
	} catch (error) {
		state.errorMessage = error instanceof Error ? error.message : 'Could not load jobs overview';
	} finally {
		state.loadingOverview = false;
	}
}

export async function loadHistoryState(state: ControllerState, reset: boolean) {
	if (state.loadingHistory || state.loadingMore) return;
	if (reset) state.loadingHistory = true;
	else state.loadingMore = true;
	state.errorMessage = '';
	try {
		const response = await listSystemJobExecutions({
			status: state.selectedStatuses.length > 0 ? state.selectedStatuses : undefined,
			query: state.query.trim() || undefined,
			before: reset ? undefined : state.history.at(-1)?.updatedAt,
			limit: historyPageLimit,
			includeRoutine: state.includeRoutine
		});
		state.history = reset
			? sortExecutions(response.executions)
			: mergeExecutions(state.history, response.executions);
		state.historyHasMore = response.hasMore;
	} catch (error) {
		state.errorMessage =
			error instanceof Error ? error.message : 'Could not load execution history';
	} finally {
		state.loadingHistory = false;
		state.loadingMore = false;
	}
}

export async function toggleScheduleState(
	state: ControllerState,
	schedule: SystemJobSchedule,
	paused: boolean
) {
	state.updatingScheduleId = schedule.id;
	state.errorMessage = '';
	try {
		const updated = paused
			? await pauseSystemJobSchedule(schedule.id)
			: await resumeSystemJobSchedule(schedule.id);
		state.schedules = state.schedules.map((current) =>
			current.id === updated.id ? updated : current
		);
	} catch (error) {
		state.errorMessage = error instanceof Error ? error.message : 'Could not update schedule';
	} finally {
		state.updatingScheduleId = undefined;
	}
}

export async function saveScheduleIntervalState(
	state: ControllerState,
	schedule: SystemJobSchedule,
	intervalSeconds: number
) {
	state.updatingIntervalId = schedule.id;
	state.errorMessage = '';
	try {
		const updated = await updateSystemJobScheduleInterval(schedule.id, intervalSeconds);
		state.schedules = state.schedules.map((current) =>
			current.id === updated.id ? updated : current
		);
	} catch (error) {
		state.errorMessage =
			error instanceof Error ? error.message : 'Could not update schedule interval';
	} finally {
		state.updatingIntervalId = undefined;
	}
}

export async function runScheduleState(state: ControllerState, schedule: SystemJobSchedule) {
	state.runningScheduleId = schedule.id;
	state.errorMessage = '';
	try {
		const updated = await runSystemJobSchedule(schedule.id);
		state.schedules = state.schedules.map((current) =>
			current.id === updated.id ? updated : current
		);
	} catch (error) {
		state.errorMessage = error instanceof Error ? error.message : 'Could not run schedule';
	} finally {
		state.runningScheduleId = undefined;
	}
}

export async function abortJobState(state: ControllerState) {
	if (!state.abortCandidate) return;
	state.abortingId = state.abortCandidate.id;
	try {
		await abortSystemJob(state.abortCandidate.id);
		state.abortCandidate = undefined;
		void loadOverviewState(state);
		void loadHistoryState(state, true);
	} catch (error) {
		state.errorMessage = error instanceof Error ? error.message : 'Could not abort job';
	} finally {
		state.abortingId = undefined;
	}
}

export async function openLogsState(state: ControllerState, execution: SystemJobExecution) {
	state.logsExecution = execution;
	state.loadingLogsId = execution.riverJobId;
	try {
		state.executionLogs = await listSystemJobExecutionLogs(execution.riverJobId);
	} finally {
		state.loadingLogsId = undefined;
	}
}

export async function saveRetentionState(state: ControllerState) {
	state.savingRetention = true;
	try {
		const settings = await updateSystemJobHistorySettings({
			retentionDays: Number(state.retentionDays),
			routineRetentionHours: Number(state.routineRetentionHours)
		});
		state.retentionDays = settings.retentionDays;
		state.routineRetentionHours = settings.routineRetentionHours;
		void loadHistoryState(state, true);
	} finally {
		state.savingRetention = false;
	}
}
