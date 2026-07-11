import { beforeEach, describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
	createMutation: vi.fn((options: () => unknown) => options()),
	client: { fetchQuery: vi.fn() },
	api: {
		abortSystemJob: vi.fn(),
		pauseSystemJobSchedule: vi.fn(),
		resumeSystemJobSchedule: vi.fn(),
		runSystemJobSchedule: vi.fn(),
		updateSystemJobScheduleInterval: vi.fn(),
		updateSystemJobHistorySettings: vi.fn(),
		getSystemJobsOverview: vi.fn(),
		listSystemJobExecutions: vi.fn(),
		listSystemJobExecutionLogs: vi.fn()
	}
}));

vi.mock('@tanstack/svelte-query', () => ({
	createMutation: mocks.createMutation,
	useQueryClient: () => mocks.client
}));
vi.mock('$lib/components/settings/system/jobs/api', () => mocks.api);

import { createSystemJobsOperations } from './systemJobs.svelte';

describe('system jobs operations', () => {
	beforeEach(() => vi.clearAllMocks());

	it('routes reads through stable query keys', async () => {
		mocks.client.fetchQuery.mockImplementation(async (options) => options.queryFn());
		mocks.api.getSystemJobsOverview.mockResolvedValue({ schedules: [] });
		mocks.api.listSystemJobExecutions.mockResolvedValue({ executions: [], hasMore: false });
		const operations = createSystemJobsOperations();
		await operations.overview();
		await operations.history({ status: ['failed'] });
		expect(mocks.client.fetchQuery.mock.calls[0]?.[0].queryKey).toEqual([
			'system',
			'jobs',
			'overview'
		]);
		expect(mocks.client.fetchQuery.mock.calls[1]?.[0].queryKey).toEqual([
			'system',
			'jobs',
			'history',
			{ status: ['failed'] }
		]);
	});
});
