import type {
	SystemJobExecution,
	SystemJobExecutionLog,
	SystemJobSchedule
} from '$lib/settings/types';
import {
	subscribeToAppEvent,
	subscribeToAppEventSourceStatus
} from '$lib/app/realtime/appEventSource';
import {
	abortJobState,
	loadHistoryState,
	loadOverviewState,
	openLogsState,
	runScheduleState,
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
import type { SystemJobsOperations } from '$lib/features/settings/resources/systemJobs.svelte';

export class SystemJobsController {
	constructor(public readonly operations: SystemJobsOperations) {}
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
	runningScheduleId = $state<string | undefined>();
	abortingId = $state<number | undefined>();
	abortCandidate = $state<{ id: number; kind: string } | undefined>();
	logsExecution = $state<SystemJobExecution | undefined>();
	executionLogs = $state<SystemJobExecutionLog[]>([]);
	loadingLogsId = $state<number | undefined>();

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
		const unsubscribeExecution = subscribeToAppEvent<SystemJobExecution>(
			'system.job.execution.updated',
			({ data: execution }) => {
				if (execution) this.applyExecutionUpdate(execution);
			}
		);
		const unsubscribeStatus = subscribeToAppEventSourceStatus((status) => {
			if (status === 'error') {
				this.errorMessage = this.errorMessage || 'Job event stream disconnected';
			}
			if (status === 'open' && this.errorMessage === 'Job event stream disconnected') {
				this.errorMessage = '';
			}
		});
		return () => {
			unsubscribeExecution();
			unsubscribeStatus();
		};
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

	async runSchedule(schedule: SystemJobSchedule) {
		await runScheduleState(this, schedule);
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
