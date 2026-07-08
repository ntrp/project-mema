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
					activeProgressLabel: 'Checking indexers'
				})
			],
			onPause: vi.fn(),
			onResume: vi.fn(),
			onAbort: vi.fn()
		});

		expect(body).toContain('RSS sync');
		expect(body).toContain('Checking indexers');
		expect(body).toContain('Pause schedule');
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
});

function systemJobSchedule(overrides: Partial<SystemJobSchedule>): SystemJobSchedule {
	return {
		id: 'schedule-1',
		name: 'Schedule',
		kind: 'media.release_search',
		queue: 'media_search',
		intervalSeconds: 900,
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
