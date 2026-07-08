import type {
	SystemJobExecution,
	SystemJobExecutionLog,
	SystemJobSchedule
} from '$lib/settings/types';
import { parseSystemEvent } from '../events/systemEventStream';
import {
	abortJobState,
	loadHistoryState,
	loadOverviewState,
	openLogsState,
	saveRetentionState,
	saveScheduleIntervalState,
	toggleScheduleState
} from './systemJobsControllerActions';
import {
	defaultHistoryIncludesExecution,
	matchesExecutionFilters,
	optionList,
	updateOneShotJobs,
	updateScheduleFromExecution,
	upsertExecution
} from './systemJobsState';

export class SystemJobsController {
	schedules = $state<SystemJobSchedule[]>([]);
	oneShotJobs = $state<SystemJobExecution[]>([]);
	history = $state<SystemJobExecution[]>([]);
	historyHasMore = $state(false);
	retentionDays = $state(30);
	routineRetentionHours = $state(24);
	includeRoutine = $state(false);
	selectedStatuses = $state<string[]>([]);
	selectedQueues = $state<string[]>([]);
	selectedKinds = $state<string[]>([]);
	query = $state('');
	loadingOverview = $state(false);
	loadingHistory = $state(false);
	loadingMore = $state(false);
	savingRetention = $state(false);
	errorMessage = $state('');
	updatingScheduleId = $state<string | undefined>();
	updatingIntervalId = $state<string | undefined>();
	abortingId = $state<number | undefined>();
	abortCandidate = $state<{ id: number; kind: string } | undefined>();
	logsExecution = $state<SystemJobExecution | undefined>();
	executionLogs = $state<SystemJobExecutionLog[]>([]);
	loadingLogsId = $state<number | undefined>();
	private source?: EventSource;

	get visibleHistory() {
		return this.history.filter((execution) =>
			matchesExecutionFilters(execution, this.selectedQueues, this.selectedKinds)
		);
	}

	get queueOptions() {
		return optionList(
			[...this.schedules, ...this.oneShotJobs, ...this.history].map((job) => job.queue)
		);
	}

	get kindOptions() {
		return optionList(
			[...this.schedules, ...this.oneShotJobs, ...this.history].map((job) => job.kind)
		);
	}

	start() {
		void this.loadOverview();
		void this.loadHistory(true);
		this.source = new EventSource('/api/events', { withCredentials: true });
		this.source.addEventListener('system.job.updated', () => void this.loadOverview());
		this.source.addEventListener('system.job.execution.updated', (event) => {
			const execution = parseSystemEvent<SystemJobExecution>(event);
			if (execution) this.applyExecutionUpdate(execution);
		});
		this.source.addEventListener('error', () => {
			this.errorMessage = this.errorMessage || 'Job event stream disconnected';
		});
		this.source.addEventListener('open', () => {
			if (this.errorMessage === 'Job event stream disconnected') this.errorMessage = '';
		});
		return () => this.source?.close();
	}

	async loadOverview() {
		await loadOverviewState(this);
	}

	async loadHistory(reset: boolean) {
		await loadHistoryState(this, reset);
	}

	applyExecutionUpdate(execution: SystemJobExecution) {
		if (execution.classification === 'fixed') {
			this.schedules = updateScheduleFromExecution(this.schedules, execution);
		}
		if (execution.classification === 'one_shot') {
			this.oneShotJobs = updateOneShotJobs(this.oneShotJobs, execution);
		}
		if (!this.includeRoutine && !defaultHistoryIncludesExecution(execution)) {
			this.history = this.history.filter((job) => job.riverJobId !== execution.riverJobId);
			return;
		}
		this.history = upsertExecution(this.history, execution).slice(0, 300);
	}

	async toggleSchedule(schedule: SystemJobSchedule, paused: boolean) {
		await toggleScheduleState(this, schedule, paused);
	}

	async saveScheduleInterval(schedule: SystemJobSchedule, intervalSeconds: number) {
		await saveScheduleIntervalState(this, schedule, intervalSeconds);
	}

	async abortJob() {
		await abortJobState(this);
	}

	async openLogs(execution: SystemJobExecution) {
		await openLogsState(this, execution);
	}

	async saveRetention() {
		await saveRetentionState(this);
	}
}
