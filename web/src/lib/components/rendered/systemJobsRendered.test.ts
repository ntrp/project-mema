import { describe, expect, it, vi } from 'vitest';

vi.mock('$lib/components/settings/system/cache/rowPulse.svelte', () => ({
	createRowPulse: () => ({
		update: vi.fn(),
		classFor: () => undefined
	})
}));

import SystemJobHistoryTable from '$lib/components/settings/system/jobs/SystemJobHistoryTable.svelte';
import SystemOneShotJobsTable from '$lib/components/settings/system/jobs/SystemOneShotJobsTable.svelte';
import SystemScheduledJobsTable from '$lib/components/settings/system/jobs/SystemScheduledJobsTable.svelte';
import {
	defaultHistoryIncludesExecution,
	mergeExecutions,
	upsertExecution
} from '$lib/components/settings/system/jobs/systemJobsState';
import type { SystemJobExecution, SystemJobSchedule } from '$lib/settings/types';
import { renderWithTooltip } from './renderHelpers';

describe('rendered system job dashboard components', () => {
	it('renders fixed scheduled jobs with pause and active abort controls', () => {
		const { body } = renderWithTooltip(SystemScheduledJobsTable, {
			schedules: [
				systemJobSchedule({
					id: 'rss_sync',
					name: 'RSS sync',
					activeStatus: 'running',
					activeRiverJobId: 99,
					activeProgressPercent: 50,
					activeProgressLabel: 'Checking indexers',
					historyPolicy: 'routine',
					intervalConfigurable: true
				})
			],
			onPause: vi.fn(),
			onResume: vi.fn(),
			onRun: vi.fn(),
			onSaveInterval: vi.fn(),
			onAbort: vi.fn()
		});

		expect(body).toContain('RSS sync');
		expect(body).toContain('Release search');
		expect(body).toContain('Checks indexer feeds');
		expect(body).toContain('Checking indexers');
		expect(body).toContain('routine');
		expect(body).toContain('Run schedule now');
		expect(body).toContain('Disable automatic job');
		expect(body).toContain('Abort active run');
	});

	it('renders one-shot jobs with progress and abort controls', () => {
		const { body } = renderWithTooltip(SystemOneShotJobsTable, {
			jobs: [systemJobExecution({ riverJobId: 51, status: 'scheduled', progressLabel: 'Waiting' })],
			onAbort: vi.fn()
		});

		expect(body).toContain('media.release_search');
		expect(body).toContain('Waiting');
		expect(body).toContain('Abort job');
	});

	it('renders execution history with logs action', () => {
		const { body } = renderWithTooltip(SystemJobHistoryTable, {
			executions: [
				systemJobExecution({ riverJobId: 52, status: 'completed', progressPercent: 100 })
			],
			loading: false,
			hasMore: false,
			loadingMore: false,
			onAbort: vi.fn(),
			onLogs: vi.fn(),
			onLoadMore: vi.fn()
		});

		expect(body).toContain('completed');
		expect(body).toContain('Show execution logs');
	});

	it('keeps execution history ordered by latest update across pagination and SSE', () => {
		const older = systemJobExecution({
			riverJobId: 1,
			updatedAt: '2026-07-03T01:02:03Z'
		});
		const newer = systemJobExecution({
			riverJobId: 2,
			updatedAt: '2026-07-03T01:04:03Z'
		});
		const replacement = systemJobExecution({
			riverJobId: 1,
			status: 'completed',
			updatedAt: '2026-07-03T01:05:03Z'
		});

		expect(mergeExecutions([newer], [older]).map((execution) => execution.riverJobId)).toEqual([
			2, 1
		]);
		expect(
			upsertExecution([newer, older], replacement).map((execution) => execution.riverJobId)
		).toEqual([1, 2]);
	});

	it('keeps routine failures in default history and hides routine successes', () => {
		expect(
			defaultHistoryIncludesExecution(
				systemJobExecution({ historyPolicy: 'routine', status: 'completed' })
			)
		).toBe(false);
		expect(
			defaultHistoryIncludesExecution(
				systemJobExecution({ historyPolicy: 'routine', status: 'retryable' })
			)
		).toBe(true);
	});
});

function systemJobSchedule(overrides: Partial<SystemJobSchedule>): SystemJobSchedule {
	return {
		id: 'schedule-1',
		name: 'Schedule',
		category: 'release_search',
		description: 'Checks indexer feeds',
		kind: 'media.release_search',
		queue: 'media_search',
		intervalSeconds: 900,
		intervalConfigurable: false,
		historyPolicy: 'standard',
		automatic: true,
		manualActionAvailable: true,
		enabled: true,
		paused: false,
		activeStatus: '',
		activeProgressLabel: '',
		activeInfoMessage: '',
		lastStatus: '',
		createdAt: '2026-07-03T01:02:03Z',
		updatedAt: '2026-07-03T01:02:03Z',
		...overrides
	};
}

function systemJobExecution(overrides: Partial<SystemJobExecution>): SystemJobExecution {
	return {
		riverJobId: 1,
		classification: 'one_shot',
		historyPolicy: 'standard',
		status: 'running',
		kind: 'media.release_search',
		queue: 'media_search',
		attempt: 1,
		maxAttempts: 3,
		priority: 1,
		progressLabel: '',
		args: '{}',
		metadata: '{}',
		errors: '[]',
		infoMessage: '',
		scheduledAt: '2026-07-03T01:02:03Z',
		createdAt: '2026-07-03T01:02:03Z',
		updatedAt: '2026-07-03T01:02:03Z',
		...overrides
	};
}
