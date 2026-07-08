import {
	abortSystemJob,
	getSystemJobsOverview,
	listSystemJobExecutionLogs,
	listSystemJobExecutions,
	pauseSystemJobSchedule,
	resumeSystemJobSchedule,
	updateSystemJobHistorySettings
} from '$lib/settings/api';
import type {
	SystemJobExecution,
	SystemJobExecutionLog,
	SystemJobSchedule
} from '$lib/settings/types';
import { parseSystemEvent } from '../events/systemEventStream';
import {
	matchesExecutionFilters,
	mergeExecutions,
	optionList,
	sortExecutions,
	updateOneShotJobs,
	updateScheduleFromExecution,
	upsertExecution
} from './systemJobsState';

const historyPageLimit = 20;

export class SystemJobsController {
	schedules = $state<SystemJobSchedule[]>([]);
	oneShotJobs = $state<SystemJobExecution[]>([]);
	history = $state<SystemJobExecution[]>([]);
	historyHasMore = $state(false);
	retentionDays = $state(30);
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
		this.loadingOverview = true;
		this.errorMessage = '';
		try {
			const overview = await getSystemJobsOverview();
			this.schedules = overview.schedules;
			this.oneShotJobs = overview.oneShotJobs;
			this.retentionDays = overview.historySettings.retentionDays;
		} catch (error) {
			this.errorMessage = error instanceof Error ? error.message : 'Could not load jobs overview';
		} finally {
			this.loadingOverview = false;
		}
	}

	async loadHistory(reset: boolean) {
		if (this.loadingHistory || this.loadingMore) return;
		if (reset) {
			this.loadingHistory = true;
		} else {
			this.loadingMore = true;
		}
		this.errorMessage = '';
		try {
			const before = reset ? undefined : this.history.at(-1)?.updatedAt;
			const response = await listSystemJobExecutions({
				status: this.selectedStatuses.length > 0 ? this.selectedStatuses : undefined,
				query: this.query.trim() || undefined,
				before,
				limit: historyPageLimit
			});
			this.history = reset
				? sortExecutions(response.executions)
				: mergeExecutions(this.history, response.executions);
			this.historyHasMore = response.hasMore;
		} catch (error) {
			this.errorMessage =
				error instanceof Error ? error.message : 'Could not load execution history';
		} finally {
			this.loadingHistory = false;
			this.loadingMore = false;
		}
	}

	applyExecutionUpdate(execution: SystemJobExecution) {
		if (execution.classification === 'fixed') {
			this.schedules = updateScheduleFromExecution(this.schedules, execution);
		}
		if (execution.classification === 'one_shot') {
			this.oneShotJobs = updateOneShotJobs(this.oneShotJobs, execution);
		}
		this.history = upsertExecution(this.history, execution).slice(0, 300);
	}

	async toggleSchedule(schedule: SystemJobSchedule, paused: boolean) {
		this.updatingScheduleId = schedule.id;
		this.errorMessage = '';
		try {
			const updated = paused
				? await pauseSystemJobSchedule(schedule.id)
				: await resumeSystemJobSchedule(schedule.id);
			this.schedules = this.schedules.map((current) =>
				current.id === updated.id ? updated : current
			);
		} catch (error) {
			this.errorMessage = error instanceof Error ? error.message : 'Could not update schedule';
		} finally {
			this.updatingScheduleId = undefined;
		}
	}

	async abortJob() {
		if (!this.abortCandidate) return;
		this.abortingId = this.abortCandidate.id;
		try {
			await abortSystemJob(this.abortCandidate.id);
			this.abortCandidate = undefined;
			void this.loadOverview();
			void this.loadHistory(true);
		} catch (error) {
			this.errorMessage = error instanceof Error ? error.message : 'Could not abort job';
		} finally {
			this.abortingId = undefined;
		}
	}

	async openLogs(execution: SystemJobExecution) {
		this.logsExecution = execution;
		this.loadingLogsId = execution.riverJobId;
		try {
			this.executionLogs = await listSystemJobExecutionLogs(execution.riverJobId);
		} finally {
			this.loadingLogsId = undefined;
		}
	}

	async saveRetention() {
		this.savingRetention = true;
		try {
			const settings = await updateSystemJobHistorySettings({
				retentionDays: Number(this.retentionDays)
			});
			this.retentionDays = settings.retentionDays;
			void this.loadHistory(true);
		} finally {
			this.savingRetention = false;
		}
	}
}
